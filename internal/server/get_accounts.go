package server

import (
	"bank/internal/model"
	"bank/pkg/api"
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) GetAccounts(ctx context.Context, _ *emptypb.Empty) (*api.GetAccountsResponse, error) {
	userID, err := getUserId(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if userID == uuid.Nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	accounts, err := s.service.GetAccounts(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &api.GetAccountsResponse{
		Accounts: mapApiAccounts(accounts),
	}, nil
}

func mapApiAccounts(accounts []*model.Account) []*api.GetAccountsResponse_Account {
	accs := make([]*api.GetAccountsResponse_Account, 0, len(accounts))

	for _, acc := range accounts {
		accs = append(accs, &api.GetAccountsResponse_Account{
			Id:     acc.ID.String(),
			Amount: decimalToMoney(acc.Balance),
		})
	}

	return accs
}
