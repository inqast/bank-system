package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Account struct {
	ID      uuid.UUID
	Owner   uuid.UUID
	Balance decimal.Decimal
}
