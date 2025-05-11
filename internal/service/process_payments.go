package service

import "context"

func (s *Service) ProcessPayments(ctx context.Context) error {
	processedPayments, err := s.creditRepo.ProcessPayments(ctx)
	if err != nil {
		return err
	}

	if len(processedPayments) == 0 {
		return nil
	}

	for _, payment := range processedPayments {
		credit, err := s.creditRepo.GetByID(ctx, payment.CreditID)
		if err != nil {
			continue
		}
		user, err := s.userRepo.GetUserByID(ctx, credit.OwnerID)
		if err != nil {
			continue
		}

		s.mailer.SendPaymentNotification(user.Email, payment)
	}

	return nil
}
