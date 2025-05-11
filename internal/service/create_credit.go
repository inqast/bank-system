package service

import (
	"context"
	"errors"

	"bank/internal/model"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (s *Service) CreateCredit(
	ctx context.Context,
	userID uuid.UUID,
	accountID uuid.UUID,
	amount decimal.Decimal,
	month int,
) (*model.Credit, error) {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if account.Owner != userID {
		return nil, errors.New("internal error")
	}

	interestRate, err := s.centralBankClient.GetCentralBankRate()
	if err != nil {
		return nil, err
	}

	margin := decimal.NewFromInt(int64(s.cfg.Margin))

	credit := &model.Credit{
		ID:           uuid.New(),
		OwnerID:      userID,
		AccountID:    account.ID,
		Amount:       amount,
		InterestRate: interestRate.Add(margin),
		Month:        month,
	}

	if _, err := s.creditRepo.CreateCredit(ctx, credit); err != nil {
		return nil, err
	}

	return credit, nil
}
