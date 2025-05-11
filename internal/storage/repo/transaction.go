package repo

import (
	"context"
	"errors"
	"gorm.io/gorm/clause"

	"bank/internal/model"
	"bank/internal/storage/dto"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Transaction struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *Transaction {
	return &Transaction{db: db}
}

func (r *Transaction) Create(_ context.Context, transaction *model.Transaction) (uuid.UUID, error) {
	dto := mapTransactionToDTO(transaction)

	if err := r.db.Create(dto).Error; err != nil {
		return uuid.UUID{}, err
	}

	return dto.ID, nil
}

func (r *Transaction) Transfer(_ context.Context, from, to uuid.UUID, amount decimal.Decimal) error {
	amountFloat, _ := amount.Float64()

	return r.db.Transaction(func(tx *gorm.DB) error {
		fromAccount := &dto.Account{}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(fromAccount, "id = ?", from).Error; err != nil {
			return err
		}

		if fromAccount.Amount < amountFloat {
			return errors.New("insufficient funds")
		}

		if err := tx.Model(&fromAccount).Update("amount", gorm.Expr("amount - ?", amount)).Error; err != nil {
			return err
		}

		if err := tx.Model(&dto.Account{}).
			Where("id = ?", to).
			Update("amount", gorm.Expr("amount + ?", amount)).Error; err != nil {
			return err
		}

		return tx.Create(&dto.Transaction{
			From:   from,
			To:     to,
			Amount: amountFloat,
			Status: "completed",
		}).Error
	})
}

func (r *Transaction) GetByUserID(_ context.Context, userID uuid.UUID) ([]*model.Transaction, error) {
	var transactions []*dto.Transaction

	var accountIDs []uuid.UUID
	if err := r.db.Model(&dto.Account{}).
		Where("owner = ?", userID).
		Pluck("id", &accountIDs).Error; err != nil {
		return nil, err
	}

	query := r.db.Model(&dto.Transaction{}).
		Where("(\"to\" IN (?) OR \"from\" IN (?))", accountIDs, accountIDs).
		Order("created_at DESC")

	if err := query.Find(&transactions).Error; err != nil {
		return nil, err
	}

	return mapTransactionsToDomain(transactions), nil
}

func mapTransactionToDTO(transaction *model.Transaction) *dto.Transaction {
	amount, _ := transaction.Amount.Float64()

	return &dto.Transaction{
		ID:     transaction.ID,
		From:   transaction.From,
		To:     transaction.To,
		Amount: amount,
		Status: mapTransactionStatusToDTO(transaction.Status),
	}
}

func mapTransactionsToDomain(
	transactionDTOs []*dto.Transaction,
) []*model.Transaction {
	transactions := make([]*model.Transaction, 0, len(transactionDTOs))

	for _, transactionDTO := range transactionDTOs {
		transactions = append(transactions, mapTransactionToDomain(transactionDTO))
	}

	return transactions
}

func mapTransactionToDomain(transaction *dto.Transaction) *model.Transaction {
	return &model.Transaction{
		ID:     transaction.ID,
		From:   transaction.From,
		To:     transaction.To,
		Amount: decimal.NewFromFloat(transaction.Amount),
		Status: mapTransactionStatusToDomain(transaction.Status),
	}
}

func mapTransactionStatusToDomain(
	status string,
) model.TransactionStatus {
	switch status {
	case "pending":
		return model.TransactionStatusPending
	case "completed":
		return model.TransactionStatusCompleted
	default:
		return model.TransactionStatusUnknown
	}
}

func mapTransactionStatusToDTO(
	status model.TransactionStatus,
) string {
	switch status {
	case model.TransactionStatusPending:
		return "active"
	case model.TransactionStatusCompleted:
		return "completed"
	default:
		return "unknown"
	}
}
