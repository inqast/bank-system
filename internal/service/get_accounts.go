package service

import (
	"bank/internal/model"
	"context"
	"github.com/google/uuid"
)

func (s *Service) GetAccounts(ctx context.Context, userID uuid.UUID) ([]*model.Account, error) {
	return s.accountRepo.GetByUserID(ctx, userID)
}
