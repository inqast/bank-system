package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", err
	}

	token, err := s.cryptoService.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}

	return token, err
}
