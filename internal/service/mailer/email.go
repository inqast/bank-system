package mailer

import (
	"encoding/json"

	"bank/internal/config"
	"bank/internal/model"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	dialer   *gomail.Dialer
	from     string
	fromName string
}

func NewEmailService(cfg *config.SMTP) *EmailService {
	return &EmailService{
		dialer:   gomail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Pass),
		from:     cfg.Email,
		fromName: cfg.Name,
	}
}

func (s *EmailService) SendPaymentNotification(userEmail string, payment *model.Payment) error {
	msg, err := json.Marshal(payment)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", s.from, s.fromName)
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", "Payment Notification")
	m.SetBody("text/html", string(msg))

	s.dialer.DialAndSend(m)

	return nil
}
