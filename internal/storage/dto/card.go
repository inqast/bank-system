package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Card struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	OwnerID         uuid.UUID `gorm:"type:uuid;not null;index"`
	AccountID       uuid.UUID `gorm:"type:uuid;not null;index"`
	Status          string    `gorm:"type:varchar(20);default:'active'"`
	EncryptedNumber []byte    `gorm:"type:bytea;not null"`
	NumberHMAC      []byte    `gorm:"type:bytea;not null"`
	EncryptedExpiry []byte    `gorm:"type:bytea;not null"`
	ExpiryHMAC      []byte    `gorm:"type:bytea;not null"`
	EncryptedCVV    []byte    `gorm:"type:bytea;not null"`
	CVVHMAC         []byte    `gorm:"type:bytea;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Card) BeforeCreate(_ *gorm.DB) (err error) {
	c.ID = uuid.New()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return
}
