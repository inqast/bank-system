package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreditID uuid.UUID `gorm:"type:uuid;not null;index"`

	Amount  float64 `gorm:"type:decimal(15,2);not null"`
	Penalty float64 `gorm:"type:decimal(15,2);default:0.0"`
	Status  string  `gorm:"type:varchar(20);default:'pending'"`

	DueDate time.Time `gorm:"not null"`
	PaidAt  *time.Time
}

func (p *Payment) BeforeCreate(_ *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
