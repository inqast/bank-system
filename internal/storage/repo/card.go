package repo

import (
	"bank/internal/model"
	"bank/internal/storage/dto"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Card struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) *Card {
	return &Card{db: db}
}

func (r *Card) Create(_ context.Context, card *model.EncryptedCard) error {
	dto, err := r.mapCardToDTO(card)
	if err != nil {
		return err
	}

	return r.db.Create(dto).Error
}

func (r *Card) GetByUserID(_ context.Context, userID uuid.UUID) ([]*model.EncryptedCard, error) {
	var cards []*dto.Card

	err := r.db.Model(&dto.Card{}).
		Where("owner_id = ?", userID).
		Find(&cards).
		Error
	if err != nil {
		return nil, err
	}

	return r.mapCardsToDomain(cards)
}

func (r *Card) mapCardsToDomain(
	cardsDTO []*dto.Card,
) ([]*model.EncryptedCard, error) {
	cards := make([]*model.EncryptedCard, 0, len(cardsDTO))

	for _, cardDTO := range cardsDTO {
		card, err := r.mapCardToDomain(cardDTO)
		if err != nil {
			return nil, err
		}

		cards = append(cards, card)
	}

	return cards, nil
}

func (r *Card) mapCardToDomain(
	card *dto.Card,
) (*model.EncryptedCard, error) {

	return &model.EncryptedCard{
		ID:              card.ID,
		OwnerID:         card.OwnerID,
		AccountID:       card.AccountID,
		Status:          mapCardStatusToDomain(card.Status),
		EncryptedNumber: card.EncryptedNumber,
		NumberHMAC:      card.NumberHMAC,
		EncryptedExpiry: card.EncryptedExpiry,
		ExpiryHMAC:      card.ExpiryHMAC,
		EncryptedCVV:    card.EncryptedCVV,
		CVVHMAC:         card.CVVHMAC,
	}, nil
}

func (r *Card) mapCardToDTO(
	card *model.EncryptedCard,
) (*dto.Card, error) {
	return &dto.Card{
		ID:              card.ID,
		OwnerID:         card.OwnerID,
		AccountID:       card.AccountID,
		Status:          mapCardStatusToDTO(card.Status),
		EncryptedNumber: card.EncryptedNumber,
		NumberHMAC:      card.NumberHMAC,
		EncryptedExpiry: card.EncryptedExpiry,
		ExpiryHMAC:      card.ExpiryHMAC,
		EncryptedCVV:    card.EncryptedCVV,
		CVVHMAC:         card.CVVHMAC,
	}, nil
}

func mapCardStatusToDomain(
	status string,
) model.CardStatus {
	switch status {
	case "active":
		return model.CardStatusActive
	case "blocked":
		return model.CardStatusBlocked
	default:
		return model.CardStatusUnknown
	}
}

func mapCardStatusToDTO(
	status model.CardStatus,
) string {
	switch status {
	case model.CardStatusActive:
		return "active"
	case model.CardStatusBlocked:
		return "blocked"
	default:
		return "unknown"
	}
}
