package service

import (
	"context"

	"bank/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Register(ctx context.Context, email string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &model.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	_, err = s.userRepo.CreateUser(ctx, user)

	return err
}
