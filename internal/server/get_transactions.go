package server

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"bank/internal/model"
	"bank/pkg/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) GetTransactions(ctx context.Context, _ *emptypb.Empty) (*api.GetTransactionsResponse, error) {
	userID, err := getUserId(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if userID == uuid.Nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	transactions, err := s.service.GetTransactions(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &api.GetTransactionsResponse{
		Transactions: mapApiTransactions(transactions),
	}, nil
}

func mapApiTransactions(transactions []*model.Transaction) []*api.GetTransactionsResponse_Transaction {
	apiTransactions := make([]*api.GetTransactionsResponse_Transaction, 0, len(transactions))

	for _, transaction := range transactions {
		apiTransactions = append(apiTransactions, mapApiTransaction(transaction))
	}

	return apiTransactions
}

func mapApiTransaction(transaction *model.Transaction) *api.GetTransactionsResponse_Transaction {
	apiPayment := &api.GetTransactionsResponse_Transaction{
		Id:     transaction.ID.String(),
		IdFrom: transaction.From.String(),
		IdTo:   transaction.To.String(),
		Amount: decimalToMoney(transaction.Amount),
		Status: mapTransactionStatus(transaction.Status),
	}

	return apiPayment
}

func mapTransactionStatus(status model.TransactionStatus) api.GetTransactionsResponse_Status {
	switch status {
	case model.TransactionStatusPending:
		return api.GetTransactionsResponse_PENDING
	case model.TransactionStatusCompleted:
		return api.GetTransactionsResponse_COMPLETED
	default:
		return api.GetTransactionsResponse_UNKNOWN
	}
}
