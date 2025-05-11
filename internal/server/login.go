package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"

	"bank/pkg/api"
)

func (s *Server) Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	if !validateEmail(req.GetEmail()) || !validatePassword(req.GetPassword()) {
		return nil, status.Error(codes.InvalidArgument, "invalid data")
	}

	token, err := s.service.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &api.LoginResponse{
		AuthToken: token,
	}, nil
}

func validateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func validatePassword(password string) bool {
	return len(password) >= 8
}
