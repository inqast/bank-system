package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	centralbank "bank/internal/client/central_bank"
	"bank/internal/config"
	"bank/internal/middleware"
	"bank/internal/scheduler"
	"bank/internal/server"
	"bank/internal/service"
	"bank/internal/service/crypto"
	"bank/internal/service/mailer"
	"bank/internal/storage"
	"bank/internal/storage/repo"
	bank "bank/pkg/api"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	localConfigPath = "./config/local.yaml"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load(localConfigPath)
	if err != nil {
		log.Fatal("load config error: ", err)
	}

	db, err := storage.NewDB(cfg.PgConfig)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.DebugLevel)

	userRepo := repo.NewUserRepository(db)
	accountRepo := repo.NewAccountRepository(db)
	transactionRepo := repo.NewTransactionRepository(db)
	creditRepo := repo.NewCreditRepository(db)
	cardRepo := repo.NewCardRepository(db)

	cryptoService, err := crypto.New(cfg.Crypto)
	if err != nil {
		logger.Fatalf("failed to init crypto service: %v", err)
	}

	cbClient := centralbank.New()
	emailService := mailer.NewEmailService(cfg.SMTP)

	bankService := service.New(
		accountRepo,
		cardRepo,
		creditRepo,
		transactionRepo,
		userRepo,
		cryptoService,
		cbClient,
		emailService,
	)

	bankServer := server.New(bankService)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			middleware.NewJWTInterceptor(cryptoService),
		),
	)

	bank.RegisterBankServer(grpcServer, bankServer)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	mux := runtime.NewServeMux()
	if err := bank.RegisterBankHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%d", 9000), opts); err != nil {
		logger.Fatalf("failed to register gateway: %v", err)
	}

	creditScheduler := scheduler.New(bankService)
	go creditScheduler.Start()

	serverErrors := make(chan error, 1)
	go func() {
		if err := http.ListenAndServe(":8080", mux); err != nil {
			serverErrors <- err
		}
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			logger.Fatalf("failed to serve: %v", err)

			serverErrors <- err
		}
	}()

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Infof("srv started")

	select {
	case err := <-serverErrors:
		logger.Fatalf("srv error: %v", err)
	case <-ctx.Done():
		grpcServer.GracefulStop()

		logger.Println("server stopped gracefully")
	}
}
