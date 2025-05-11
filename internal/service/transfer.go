package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (s *Service) Transfer(
	ctx context.Context,
	userID uuid.UUID,
	from uuid.UUID, to uuid.UUID,
	amount decimal.Decimal,
) error {
	account, err := s.accountRepo.GetByID(ctx, from)
	if err != nil {
		return err
	}
	if account.Owner != userID {
		return errors.New("internal error")
	}

	return s.transactionRepo.Transfer(
		ctx, from, to, amount,
	)
}
