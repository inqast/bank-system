package crypto

import (
	"crypto"
	"os"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"golang.org/x/crypto/openpgp/armor"
)

func (s *Service) generatePGPKeys() error {
	if s.isExist() {
		return nil
	}

	config := &packet.Config{
		RSABits:     4096,
		DefaultHash: crypto.SHA256,
		Time:        time.Now,
	}

	entity, err := openpgp.NewEntity("bank-app", "", "example@mail.ru", config)
	if err != nil {
		return err
	}

	err = s.generatePublic(entity)
	if err != nil {
		return err
	}

	return s.generatePrivate(entity)
}

func (s *Service) isExist() bool {
	_, err := os.Stat(s.cfg.PGPPublicPath)
	if err != nil {
		return false
	}

	_, err = os.Stat(s.cfg.PGPPrivatePath)
	if err != nil {
		return false
	}

	return true
}

func (s *Service) generatePublic(entity *openpgp.Entity) error {
	pubFile, err := os.Create(s.cfg.PGPPublicPath)
	if err != nil {
		return err
	}
	defer pubFile.Close()

	pubArmor, err := armor.Encode(pubFile, openpgp.PublicKeyType, nil)
	if err != nil {
		return err
	}
	err = entity.Serialize(pubArmor)
	if err != nil {
		return err
	}
	pubArmor.Close()

	return nil
}

func (s *Service) generatePrivate(entity *openpgp.Entity) error {
	privFile, err := os.Create(s.cfg.PGPPrivatePath)
	if err != nil {
		return err
	}
	defer privFile.Close()

	privArmor, err := armor.Encode(privFile, openpgp.PrivateKeyType, nil)
	if err != nil {
		return err
	}
	err = entity.SerializePrivate(privArmor, nil)
	if err != nil {
		return err
	}
	privArmor.Close()

	return nil
}
