package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	From uuid.UUID `gorm:"type:uuid"`
	To   uuid.UUID `gorm:"type:uuid;not null"`

	Amount float64 `gorm:"type:decimal(15,2);not null"`
	Status string  `gorm:"type:varchar(20);not null"`

	CreatedAt time.Time
}

func (t *Transaction) BeforeCreate(_ *gorm.DB) (err error) {
	t.ID = uuid.New()
	t.CreatedAt = time.Now()
	return
}
