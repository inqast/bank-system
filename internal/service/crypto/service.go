package crypto

import (
	"bank/internal/config"
	"errors"
	"golang.org/x/crypto/openpgp"
	"os"
)

type Service struct {
	cfg        *config.Crypto
	pgpPublic  *openpgp.Entity
	pgpPrivate *openpgp.Entity
}

func New(cfg *config.Crypto) (*Service, error) {
	service := &Service{
		cfg: cfg,
	}

	err := service.generatePGPKeys()
	if err != nil {
		return nil, err
	}

	err = service.loadPGPKeys()
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (s *Service) loadPGPKeys() error {
	err := s.loadPublicKey()
	if err != nil {
		return err
	}

	err = s.loadPrivateKey()
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) loadPublicKey() error {
	entity, err := loadKey(s.cfg.PGPPublicPath)
	if err != nil {
		return err
	}

	s.pgpPublic = entity

	return nil
}

func (s *Service) loadPrivateKey() error {
	entity, err := loadKey(s.cfg.PGPPrivatePath)
	if err != nil {
		return err
	}

	s.pgpPrivate = entity

	return nil
}

func loadKey(path string) (*openpgp.Entity, error) {
	keyFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer keyFile.Close()

	entityList, err := openpgp.ReadArmoredKeyRing(keyFile)
	if err != nil {
		return nil, err
	}

	if len(entityList) == 0 {
		return nil, errors.New("no keys found in keyring")
	}

	return entityList[0], nil
}
