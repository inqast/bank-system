package server

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"bank/pkg/api"
)

func (s *Server) CreateCard(ctx context.Context, req *api.CreateCardRequest) (*api.CreateCardResponse, error) {
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

	card, err := s.service.CreateCard(ctx, userID, accountID)
	if err != nil {
		return nil, err
	}

	return &api.CreateCardResponse{
		Card: mapApiCard(card),
	}, nil
}
