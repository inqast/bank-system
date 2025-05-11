package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email         string    `gorm:"uniqueIndex;not null;size=255"`
	PasswordHash  string    `gorm:"not null"`
	EmailVerified bool      `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	u.ID = uuid.New()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return
}
