package model

import (
	"github.com/shopspring/decimal"
	"time"

	"github.com/google/uuid"
)

type PaymentStatus uint8

const (
	PaymentStatusUnknown PaymentStatus = iota
	PaymentStatusPending
	PaymentStatusOverdue
	PaymentStatusPaid
)

type Payment struct {
	ID       uuid.UUID
	CreditID uuid.UUID
	DueDate  time.Time
	Amount   decimal.Decimal
	Penalty  decimal.Decimal
	Status   PaymentStatus
	PaidAt   *time.Time
}
