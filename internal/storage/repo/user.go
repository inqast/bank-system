package repo

import (
	"context"
	"errors"

	"bank/internal/model"
	"bank/internal/storage/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *User {
	return &User{db: db}
}

var (
	ErrEmailExists        = errors.New("email already registered")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrWeakPassword       = errors.New("password must be at least 8 characters")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

func (r *User) CreateUser(_ context.Context, user *model.User) (uuid.UUID, error) {
	userDTO := mapUserToDTO(user)

	var count int64
	r.db.Model(&dto.User{}).Where("email = ?", userDTO.Email).Count(&count)
	if count > 0 {
		return uuid.UUID{}, ErrEmailExists
	}

	if err := r.db.Create(userDTO).Error; err != nil {
		return uuid.UUID{}, err
	}

	return userDTO.ID, nil
}

func (r *User) GetUserByEmail(_ context.Context, email string) (*model.User, error) {
	user := &dto.User{}
	err := r.db.Where("email = ?", email).First(user).Error

	return mapUserToDomain(user), err
}

func (r *User) GetUserByID(_ context.Context, userID uuid.UUID) (*model.User, error) {
	user := &dto.User{}
	err := r.db.Where("id = ?", userID.String()).First(user).Error

	return mapUserToDomain(user), err
}

func mapUserToDTO(user *model.User) *dto.User {
	return &dto.User{
		ID:            user.ID,
		Email:         user.Email,
		PasswordHash:  user.PasswordHash,
		EmailVerified: user.EmailVerified,
	}
}

func mapUserToDomain(user *dto.User) *model.User {
	return &model.User{
		ID:            user.ID,
		Email:         user.Email,
		PasswordHash:  user.PasswordHash,
		EmailVerified: user.EmailVerified,
	}
}
