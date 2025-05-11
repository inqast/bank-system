package service

import (
	"context"

	"bank/internal/model"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (s *Service) CreateAccount(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	account := &model.Account{
		Owner:   userID,
		Balance: decimal.Zero,
	}

	return s.accountRepo.Create(ctx, account)
}
