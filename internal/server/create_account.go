package server

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"bank/pkg/api"
)

func (s *Server) CreateAccount(ctx context.Context, _ *emptypb.Empty) (*api.CreateAccountResponse, error) {
	userID, err := getUserId(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if userID == uuid.Nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	accountID, err := s.service.CreateAccount(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &api.CreateAccountResponse{
		Id: accountID.String(),
	}, nil
}
