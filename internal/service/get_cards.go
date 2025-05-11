package service

import (
	"context"

	"bank/internal/model"
	"github.com/google/uuid"
)

func (s *Service) GetCards(ctx context.Context, userID uuid.UUID) ([]*model.DecryptedCard, error) {
	encryptedCards, err := s.cardRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.getDecryptedCards(encryptedCards)
}

func (s *Service) getDecryptedCards(
	encryptedCards []*model.EncryptedCard,
) ([]*model.DecryptedCard, error) {
	decryptedCards := make([]*model.DecryptedCard, 0, len(encryptedCards))

	for _, encryptedCard := range encryptedCards {
		decryptedCard, err := s.cryptoService.DecryptCard(encryptedCard)
		if err != nil {
			return nil, err
		}

		decryptedCards = append(decryptedCards, decryptedCard)
	}

	return decryptedCards, nil
}
