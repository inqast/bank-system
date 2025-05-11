package crypto

import (
	"bytes"
	"errors"
	"io"

	"bank/internal/model"
	"golang.org/x/crypto/openpgp"
)

func (s *Service) DecryptCard(
	card *model.EncryptedCard,
) (*model.DecryptedCard, error) {
	if !s.verifyHMAC(card.EncryptedNumber, card.NumberHMAC) {
		return nil, errors.New("invalid HMAC for card number")
	}
	decryptedNum, err := s.decryptPGP(card.EncryptedNumber)
	if err != nil {
		return nil, err
	}

	if !s.verifyHMAC(card.EncryptedExpiry, card.ExpiryHMAC) {
		return nil, errors.New("invalid HMAC for expiry date")
	}
	decryptedExp, err := s.decryptPGP(card.EncryptedExpiry)
	if err != nil {
		return nil, err
	}

	if !s.verifyHMAC(card.EncryptedCVV, card.CVVHMAC) {
		return nil, errors.New("invalid HMAC for cvv")
	}
	decryptedCVV, err := s.decryptPGP(card.EncryptedCVV)
	if err != nil {
		return nil, err
	}

	return &model.DecryptedCard{
		ID:        card.ID,
		OwnerID:   card.OwnerID,
		AccountID: card.AccountID,
		Status:    card.Status,
		Number:    string(decryptedNum),
		Expiry:    string(decryptedExp),
		CVV:       string(decryptedCVV),
	}, nil
}

func (s *Service) EncryptCard(
	card *model.DecryptedCard,
) (*model.EncryptedCard, error) {
	encryptedNumber, err := s.encryptPGP([]byte(card.Number))
	if err != nil {
		return nil, err
	}
	encryptedExpiry, err := s.encryptPGP([]byte(card.Expiry))
	if err != nil {
		return nil, err
	}
	encryptedCVV, err := s.encryptPGP([]byte(card.CVV))
	if err != nil {
		return nil, err
	}

	numberHMAC := s.generateHMAC(encryptedNumber)
	expiryHMAC := s.generateHMAC(encryptedExpiry)
	cvvHMAC := s.generateHMAC(encryptedCVV)

	return &model.EncryptedCard{
		ID:              card.ID,
		OwnerID:         card.OwnerID,
		AccountID:       card.AccountID,
		Status:          card.Status,
		EncryptedNumber: encryptedNumber,
		NumberHMAC:      numberHMAC,
		EncryptedExpiry: encryptedExpiry,
		ExpiryHMAC:      expiryHMAC,
		EncryptedCVV:    encryptedCVV,
		CVVHMAC:         cvvHMAC,
	}, nil
}

func (s *Service) encryptPGP(data []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	w, err := openpgp.Encrypt(buf, []*openpgp.Entity{s.pgpPublic}, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *Service) decryptPGP(data []byte) ([]byte, error) {
	md, err := openpgp.ReadMessage(bytes.NewReader(data), openpgp.EntityList{s.pgpPrivate}, nil, nil)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(md.UnverifiedBody)
}
