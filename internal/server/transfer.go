package server

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"bank/pkg/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Transfer(ctx context.Context, req *api.TransferRequest) (*emptypb.Empty, error) {
	userID, err := getUserId(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if userID == uuid.Nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	accountFrom, err := uuid.Parse(req.GetIdFrom())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	accountTo, err := uuid.Parse(req.GetIdTo())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	amount := moneyToDecimal(req.Amount)
	if amount.IsNegative() {
		return nil, status.Error(codes.InvalidArgument, "below zero")
	}

	err = s.service.Transfer(ctx, userID, accountFrom, accountTo, amount)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
