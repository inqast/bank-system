package model

import (
	"github.com/google/uuid"
)

type CardStatus uint8

const (
	CardStatusUnknown CardStatus = iota
	CardStatusActive
	CardStatusBlocked
)

type EncryptedCard struct {
	ID        uuid.UUID
	OwnerID   uuid.UUID
	AccountID uuid.UUID
	Status    CardStatus

	EncryptedNumber []byte
	NumberHMAC      []byte
	EncryptedExpiry []byte
	ExpiryHMAC      []byte
	CVVHash         []byte
	EncryptedCVV    []byte
	CVVHMAC         []byte
}

type DecryptedCard struct {
	ID        uuid.UUID
	OwnerID   uuid.UUID
	AccountID uuid.UUID
	Status    CardStatus

	Number string
	Expiry string
	CVV    string
}
