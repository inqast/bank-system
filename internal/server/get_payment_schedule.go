package server

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"bank/internal/model"
	"bank/pkg/api"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) GetPaymentSchedule(ctx context.Context, req *api.GetPaymentScheduleRequest) (*api.GetPaymentScheduleResponse, error) {
	userID, err := getUserId(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if userID == uuid.Nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	creditID, err := uuid.Parse(req.GetCreditId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	payments, err := s.service.GetPaymentSchedule(ctx, userID, creditID)
	if err != nil {
		return nil, err
	}

	return &api.GetPaymentScheduleResponse{
		Payments: mapApiPayments(payments),
	}, nil
}

func mapApiPayments(payments []*model.Payment) []*api.GetPaymentScheduleResponse_Payment {
	apiPayments := make([]*api.GetPaymentScheduleResponse_Payment, 0, len(payments))

	for _, payment := range payments {
		apiPayments = append(apiPayments, mapApiPayment(payment))
	}

	return apiPayments
}

func mapApiPayment(payment *model.Payment) *api.GetPaymentScheduleResponse_Payment {
	apiPayment := &api.GetPaymentScheduleResponse_Payment{
		ID:       payment.ID.String(),
		CreditID: payment.CreditID.String(),
		DueDate:  timestamppb.New(payment.DueDate),
		Amount:   decimalToMoney(payment.Amount),
		Penalty:  decimalToMoney(payment.Penalty),
		Status:   mapPaymentStatus(payment.Status),
	}

	if payment.PaidAt != nil {
		apiPayment.PaidAt = timestamppb.New(*payment.PaidAt)
	}

	return apiPayment
}

func mapPaymentStatus(status model.PaymentStatus) api.GetPaymentScheduleResponse_Payment_Status {
	switch status {
	case model.PaymentStatusPending:
		return api.GetPaymentScheduleResponse_Payment_PENDING
	case model.PaymentStatusOverdue:
		return api.GetPaymentScheduleResponse_Payment_OVERDUE
	case model.PaymentStatusPaid:
		return api.GetPaymentScheduleResponse_Payment_PAID
	default:
		return api.GetPaymentScheduleResponse_Payment_UNKNOWN
	}
}
