package server

import (
	"bank/internal/model"
	"bank/pkg/api"
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) GetCards(ctx context.Context, _ *emptypb.Empty) (*api.GetCardsResponse, error) {
	userID, err := getUserId(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if userID == uuid.Nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	cards, err := s.service.GetCards(
		ctx,
		userID)
	if err != nil {
		return nil, err
	}

	return &api.GetCardsResponse{
		Cards: mapApiCards(cards),
	}, nil
}

func mapApiCards(cards []*model.DecryptedCard) []*api.Card {
	apiCards := make([]*api.Card, 0, len(cards))

	for _, card := range cards {
		apiCards = append(apiCards, mapApiCard(card))
	}

	return apiCards
}

func mapApiCard(card *model.DecryptedCard) *api.Card {
	return &api.Card{
		Number:  card.Number,
		ExpDate: card.Expiry,
		CVV:     card.CVV,
	}
}
