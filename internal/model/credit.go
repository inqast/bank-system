package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreditStatus uint8

const (
	CreditStatusUnknown CreditStatus = iota
	CreditStatusActive
	CreditStatusPaid
)

type Credit struct {
	ID           uuid.UUID
	OwnerID      uuid.UUID
	AccountID    uuid.UUID
	Amount       decimal.Decimal
	InterestRate decimal.Decimal
	Month        int
	Status       CreditStatus
}
