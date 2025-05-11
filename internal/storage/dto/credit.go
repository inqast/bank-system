package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Credit struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	OwnerID   uuid.UUID `gorm:"type:uuid;not null;index"`
	AccountID uuid.UUID `gorm:"type:uuid;not null;index"`

	Amount       float64 `gorm:"type:decimal(15,2);not null"`
	InterestRate float64 `gorm:"type:decimal(5,2);not null"`
	Month        int     `gorm:"not null"`
	Status       string  `gorm:"type:varchar(20);default:'active'"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Credit) BeforeCreate(_ *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}

	c.CreatedAt = time.Now()
	return
}
