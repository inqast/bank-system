package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
)

func (s *Service) generateHMAC(data []byte) []byte {
	h := hmac.New(sha256.New, []byte(s.cfg.HMACSecret))
	h.Write(data)
	return h.Sum(nil)
}

func (s *Service) verifyHMAC(data []byte, receivedHMAC []byte) bool {
	expectedHMAC := s.generateHMAC(data)

	return hmac.Equal(receivedHMAC, expectedHMAC)
}
