package repo

import (
	"context"
	"errors"
	"math"
	"time"

	"bank/internal/model"
	"bank/internal/storage/dto"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Credit struct {
	db *gorm.DB
}

var (
	ErrRecordNotFound = errors.New("record not found")
)

func NewCreditRepository(db *gorm.DB) *Credit {
	return &Credit{db: db}
}

func (r *Credit) CreateCredit(_ context.Context, credit *model.Credit) (uuid.UUID, error) {
	creditDto := mapCreditToDTO(credit)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		schedule, err := r.calculatePaymentSchedule(creditDto)
		if err != nil {
			return err
		}

		if err = tx.Create(creditDto).Error; err != nil {
			return err
		}

		for _, payment := range schedule {
			payment.CreditID = credit.ID
			if err = tx.Create(&payment).Error; err != nil {
				return err
			}
		}

		account := &dto.Account{}
		if err = tx.First(account, credit.AccountID).Error; err != nil {
			return err
		}

		if err := tx.Model(&account).Update("amount", gorm.Expr("amount + ?", creditDto.Amount)).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	return creditDto.ID, nil
}

func (r *Credit) calculatePaymentSchedule(credit *dto.Credit) ([]*dto.Payment, error) {
	monthlyRate := credit.InterestRate / 100 / 12
	annuity := (credit.Amount * monthlyRate) / (1 - math.Pow(1+monthlyRate, float64(-credit.Month)))

	schedule := make([]*dto.Payment, credit.Month)
	date := time.Now().AddDate(0, 1, 0)

	for i := 0; i < credit.Month; i++ {
		schedule[i] = &dto.Payment{
			DueDate: date.AddDate(0, i, 0),
			Amount:  math.Round(annuity*100) / 100,
			Status:  "pending",
		}
	}

	return schedule, nil
}

func (r *Credit) GetByID(_ context.Context, id uuid.UUID) (*model.Credit, error) {
	credit := &dto.Credit{}
	err := r.db.
		First(credit, "id = ?", id).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}

	return mapCreditToDomain(credit), err
}

func (r *Credit) GetPaymentSchedule(_ context.Context, creditID uuid.UUID) ([]*model.Payment, error) {
	var payments []*dto.Payment

	err := r.db.Model(&dto.Payment{}).
		Where("credit_id = ?", creditID).
		Order("due_date ASC").
		Find(&payments).
		Error
	if err != nil {
		return nil, err
	}

	return mapPaymentsToDomain(payments), nil
}

func (r *Credit) ProcessPayments(_ context.Context) ([]*model.Payment, error) {
	payments, err := r.findUnpaid()
	if err != nil {
		return nil, err
	}

	processedPayments := make([]*model.Payment, 0, len(payments))
	for _, payment := range payments {
		err = r.db.Transaction(func(tx *gorm.DB) error {
			credit := &dto.Credit{}
			if err = tx.First(credit, payment.CreditID).Error; err != nil {
				return err
			}

			account := &dto.Account{}
			if err = tx.First(account, credit.AccountID).Error; err != nil {
				return err
			}

			if err = processPayment(tx, payment, account); err != nil {
				return err
			}

			if err = tx.Save(&payment).Error; err != nil {
				return err
			}

			return nil
		})
		if err == nil {
			processedPayments = append(processedPayments, mapPaymentToDomain(payment))
		}
	}

	return processedPayments, nil
}

func (r *Credit) findUnpaid() ([]*dto.Payment, error) {
	var payments []*dto.Payment
	now := time.Now().UTC().Truncate(24 * time.Hour)

	if err := r.db.Debug().Where("status in ('pending', 'overdue') AND due_date < ?", now).Find(&payments).Error; err != nil {
		return nil, err
	}

	return payments, nil
}

func processPayment(
	tx *gorm.DB,
	payment *dto.Payment,
	account *dto.Account,
) error {
	now := time.Now()

	amountToPay := payment.Amount + payment.Penalty

	if account.Amount < amountToPay {
		payment.Penalty += payment.Amount * 0.10
		payment.Status = "overdue"

		return nil
	}

	if err := tx.Model(&account).Update("amount", gorm.Expr("amount - ?", amountToPay)).Error; err != nil {
		return err
	}

	payment.Status = "paid"
	payment.PaidAt = &now

	return nil
}

func mapCreditToDTO(credit *model.Credit) *dto.Credit {
	amount, _ := credit.Amount.Float64()
	interestRate, _ := credit.InterestRate.Float64()

	return &dto.Credit{
		ID:           credit.ID,
		OwnerID:      credit.OwnerID,
		AccountID:    credit.AccountID,
		Amount:       amount,
		InterestRate: interestRate,
		Month:        credit.Month,
		Status:       mapCreditStatusToDTO(credit.Status),
	}
}

func mapCreditToDomain(credit *dto.Credit) *model.Credit {
	return &model.Credit{
		ID:           credit.ID,
		OwnerID:      credit.OwnerID,
		AccountID:    credit.AccountID,
		Amount:       decimal.NewFromFloat(credit.Amount),
		InterestRate: decimal.NewFromFloat(credit.InterestRate),
		Month:        credit.Month,
		Status:       mapCreditStatusToDomain(credit.Status),
	}
}

func mapCreditStatusToDomain(
	status string,
) model.CreditStatus {
	switch status {
	case "active":
		return model.CreditStatusActive
	case "paid":
		return model.CreditStatusPaid
	default:
		return model.CreditStatusUnknown
	}
}

func mapCreditStatusToDTO(
	status model.CreditStatus,
) string {
	switch status {
	case model.CreditStatusActive:
		return "active"
	case model.CreditStatusPaid:
		return "paid"
	default:
		return "unknown"
	}
}

func mapPaymentsToDomain(
	paymentDTOs []*dto.Payment,
) []*model.Payment {
	payments := make([]*model.Payment, 0, len(paymentDTOs))

	for _, paymentDTO := range paymentDTOs {
		payments = append(payments, mapPaymentToDomain(paymentDTO))
	}

	return payments
}

func mapPaymentToDomain(payment *dto.Payment) *model.Payment {
	return &model.Payment{
		ID:       payment.ID,
		CreditID: payment.CreditID,
		DueDate:  payment.DueDate,
		Amount:   decimal.NewFromFloat(payment.Amount),
		Penalty:  decimal.NewFromFloat(payment.Penalty),
		Status:   mapPaymentStatusToDomain(payment.Status),
		PaidAt:   payment.PaidAt,
	}
}

func mapPaymentStatusToDomain(
	status string,
) model.PaymentStatus {
	switch status {
	case "pending":
		return model.PaymentStatusPending
	case "overdue":
		return model.PaymentStatusOverdue
	case "paid":
		return model.PaymentStatusPaid
	default:
		return model.PaymentStatusUnknown
	}
}
