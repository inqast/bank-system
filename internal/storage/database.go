package storage

import (
	"errors"
	"fmt"

	"bank/internal/config"
	"bank/internal/storage/dto"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg *config.PgConfig) (*gorm.DB, error) {
	if cfg.Host == "" || cfg.Port == "" || cfg.Username == "" {
		return nil, errors.New("invalid database configuration")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto")
	err = db.AutoMigrate(
		&dto.User{},
		&dto.Account{},
		&dto.Card{},
		&dto.Transaction{},
		&dto.Credit{},
		&dto.Payment{},
	)

	return db, err
}
