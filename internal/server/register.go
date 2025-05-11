package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"bank/pkg/api"
)

func (s *Server) Register(ctx context.Context, req *api.RegisterRequest) (*emptypb.Empty, error) {
	if !validatePassword(req.GetPassword()) {
		return nil, status.Error(codes.InvalidArgument, "invalid data")
	}

	err := s.service.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
