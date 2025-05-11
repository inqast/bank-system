package service

import (
	"bank/internal/config"
	"bank/internal/model"
	"context"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type (
	accountRepo interface {
		Create(ctx context.Context, account *model.Account) (uuid.UUID, error)
		GetByID(ctx context.Context, id uuid.UUID) (*model.Account, error)
		GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Account, error)
	}

	cardRepo interface {
		Create(ctx context.Context, card *model.EncryptedCard) error
		GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.EncryptedCard, error)
	}

	creditRepo interface {
		CreateCredit(ctx context.Context, credit *model.Credit) (uuid.UUID, error)
		GetByID(ctx context.Context, id uuid.UUID) (*model.Credit, error)
		GetPaymentSchedule(ctx context.Context, creditID uuid.UUID) ([]*model.Payment, error)
		ProcessPayments(ctx context.Context) ([]*model.Payment, error)
	}

	transactionRepo interface {
		Create(ctx context.Context, transaction *model.Transaction) (uuid.UUID, error)
		Transfer(ctx context.Context, from, to uuid.UUID, amount decimal.Decimal) error
		GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Transaction, error)
	}

	userRepo interface {
		CreateUser(ctx context.Context, user *model.User) (uuid.UUID, error)
		GetUserByEmail(ctx context.Context, email string) (*model.User, error)
		GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error)
	}

	cryptoService interface {
		DecryptCard(
			card *model.EncryptedCard,
		) (*model.DecryptedCard, error)
		EncryptCard(
			card *model.DecryptedCard,
		) (*model.EncryptedCard, error)
		GenerateJWT(userID uuid.UUID) (string, error)
	}

	centralBankClient interface {
		GetCentralBankRate() (decimal.Decimal, error)
	}

	mailer interface {
		SendPaymentNotification(userEmail string, payment *model.Payment) error
	}
)

type Service struct {
	cfg               config.Business
	accountRepo       accountRepo
	cardRepo          cardRepo
	creditRepo        creditRepo
	transactionRepo   transactionRepo
	userRepo          userRepo
	cryptoService     cryptoService
	centralBankClient centralBankClient
	mailer            mailer
}

func New(
	accountRepo accountRepo,
	cardRepo cardRepo,
	creditRepo creditRepo,
	transactionRepo transactionRepo,
	userRepo userRepo,
	cryptoService cryptoService,
	centralBankClient centralBankClient,
	mailer mailer,
) *Service {
	return &Service{
		accountRepo:       accountRepo,
		cardRepo:          cardRepo,
		creditRepo:        creditRepo,
		transactionRepo:   transactionRepo,
		userRepo:          userRepo,
		cryptoService:     cryptoService,
		centralBankClient: centralBankClient,
		mailer:            mailer,
	}
}
