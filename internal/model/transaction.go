package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionStatus uint8

const (
	TransactionStatusUnknown TransactionStatus = iota
	TransactionStatusPending
	TransactionStatusCompleted
)

type Transaction struct {
	ID     uuid.UUID
	From   uuid.UUID
	To     uuid.UUID
	Amount decimal.Decimal
	Status TransactionStatus
}
