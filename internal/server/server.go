package server

import (
	"context"
	"google.golang.org/genproto/googleapis/type/money"

	"bank/internal/model"
	"bank/pkg/api"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type service interface {
	CreateAccount(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	GetAccounts(ctx context.Context, userID uuid.UUID) ([]*model.Account, error)
	CreateCard(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) (*model.DecryptedCard, error)
	CreateCredit(
		ctx context.Context,
		userID uuid.UUID,
		accountID uuid.UUID,
		amount decimal.Decimal,
		month int,
	) (*model.Credit, error)
	GetCards(ctx context.Context, userID uuid.UUID) ([]*model.DecryptedCard, error)
	GetPaymentSchedule(ctx context.Context, userID uuid.UUID, creditID uuid.UUID) ([]*model.Payment, error)
	GetTransactions(ctx context.Context, userID uuid.UUID) ([]*model.Transaction, error)
	Transfer(
		ctx context.Context,
		userID uuid.UUID,
		from uuid.UUID,
		to uuid.UUID,
		amount decimal.Decimal,
	) error

	Login(ctx context.Context, email string, password string) (string, error)
	Register(ctx context.Context, email string, password string) error
}

type Server struct {
	api.UnimplementedBankServer
	service service
}

func New(service service) *Server {
	return &Server{
		service: service,
	}
}

func getUserId(ctx context.Context) (uuid.UUID, error) {
	userID := ctx.Value("userID")

	if userID == nil {
		return uuid.Nil, nil
	}

	return uuid.Parse(userID.(string))
}

func moneyToDecimal(m *money.Money) decimal.Decimal {
	d := decimal.NewFromInt(m.Units)

	nanoDecimal := decimal.NewFromInt(int64(m.Nanos)).Div(decimal.NewFromInt(1_000_000_000))

	return d.Add(nanoDecimal)
}

func decimalToMoney(d decimal.Decimal) *money.Money {
	units := d.Truncate(0).IntPart()
	nanosDecimal := d.Sub(decimal.NewFromInt(units))
	nanos := nanosDecimal.Mul(decimal.NewFromInt(1_000_000_000)).IntPart()

	return &money.Money{
		CurrencyCode: "RUB",
		Units:        units,
		Nanos:        int32(nanos),
	}
}
