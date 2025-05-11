package service

import (
	"context"

	"bank/internal/model"
	"github.com/google/uuid"
)

func (s *Service) GetTransactions(ctx context.Context, userID uuid.UUID) ([]*model.Transaction, error) {
	return s.transactionRepo.GetByUserID(ctx, userID)
}
