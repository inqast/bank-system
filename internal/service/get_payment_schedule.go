package service

import (
	"context"
	"errors"

	"bank/internal/model"
	"github.com/google/uuid"
)

func (s *Service) GetPaymentSchedule(ctx context.Context, userID uuid.UUID, creditID uuid.UUID) ([]*model.Payment, error) {
	credit, err := s.creditRepo.GetByID(ctx, creditID)
	if err != nil {
		return nil, err
	}
	if credit.OwnerID != userID {
		return nil, errors.New("not authorized to get payment schedule")
	}

	return s.creditRepo.GetPaymentSchedule(ctx, creditID)
}
