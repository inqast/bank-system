package repo

import (
	"bank/internal/model"
	"bank/internal/storage/dto"
	"context"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Account struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *Account {
	return &Account{db: db}
}

func (r *Account) Create(_ context.Context, account *model.Account) (uuid.UUID, error) {
	dto := mapAccountToDTO(account)

	if err := r.db.Create(dto).Error; err != nil {
		return uuid.UUID{}, err
	}

	return dto.ID, nil
}

func (r *Account) GetByID(_ context.Context, id uuid.UUID) (*model.Account, error) {
	account := &dto.Account{}
	err := r.db.First(account, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return mapAccountToDomain(account), nil
}

func (r *Account) GetByUserID(_ context.Context, userID uuid.UUID) ([]*model.Account, error) {
	var accounts []*dto.Account

	err := r.db.Model(&dto.Account{}).
		Where("owner = ?", userID).
		Find(&accounts).
		Error
	if err != nil {
		return nil, err
	}

	return mapAccountsToDomain(accounts)
}

func mapAccountToDTO(acc *model.Account) *dto.Account {
	amount, _ := acc.Balance.Float64()

	return &dto.Account{
		ID:     acc.ID,
		Owner:  acc.Owner,
		Amount: amount,
	}
}

func mapAccountsToDomain(
	accounts []*dto.Account,
) ([]*model.Account, error) {
	cards := make([]*model.Account, 0, len(accounts))

	for _, account := range accounts {
		cards = append(cards, mapAccountToDomain(account))
	}

	return cards, nil
}

func mapAccountToDomain(acc *dto.Account) *model.Account {
	return &model.Account{
		ID:      acc.ID,
		Owner:   acc.Owner,
		Balance: decimal.NewFromFloat(acc.Amount),
	}
}
