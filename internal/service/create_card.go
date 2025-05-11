package service

import (
	"context"
	"errors"

	"bank/internal/model"
	generator "bank/internal/utils/card"
	"github.com/google/uuid"
)

func (s *Service) CreateCard(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) (*model.DecryptedCard, error) {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if account.Owner != userID {
		return nil, errors.New("internal error")
	}

	decryptedCard := &model.DecryptedCard{
		OwnerID:   userID,
		AccountID: accountID,
		Status:    model.CardStatusActive,
		Number:    generator.GenerateCardNumber(),
		Expiry:    generator.GenerateExpiryDate(),
		CVV:       generator.GenerateCVV(),
	}

	encryptedCard, err := s.cryptoService.EncryptCard(decryptedCard)
	if err != nil {
		return nil, err
	}

	if err = s.cardRepo.Create(ctx, encryptedCard); err != nil {
		return nil, err
	}

	return decryptedCard, nil
}
