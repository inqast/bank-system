package server

import (
	"bank/pkg/api"
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateCredit(ctx context.Context, req *api.CreateCreditRequest) (*api.CreateCreditResponse, error) {
	userID, err := getUserId(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if userID == uuid.Nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	accountID, err := uuid.Parse(req.GetAccountId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	amount := moneyToDecimal(req.Amount)
	if amount.IsNegative() {
		return nil, status.Error(codes.InvalidArgument, "below zero")
	}

	credit, err := s.service.CreateCredit(
		ctx,
		userID, accountID,
		moneyToDecimal(req.Amount),
		int(req.GetMonth()),
	)
	if err != nil {
		return nil, err
	}

	return &api.CreateCreditResponse{
		Id:           credit.ID.String(),
		InterestRate: credit.InterestRate.IntPart(),
		Months:       int64(credit.Month),
	}, nil
}
