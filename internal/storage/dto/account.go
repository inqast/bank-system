package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Amount    float64   `gorm:"type:decimal(15,2);default:0.0"`
	Owner     uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt time.Time
}

func (a *Account) BeforeCreate(_ *gorm.DB) (err error) {
	a.ID = uuid.New()
	a.CreatedAt = time.Now()
	return
}
