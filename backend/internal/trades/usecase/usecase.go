package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/cs2-p2p-skins/backend/internal/entity"
	"github.com/cs2-p2p-skins/backend/internal/trades"
)

var ErrTradeNotFound = errors.New("trade not found")

type usecase struct {
	repo trades.Repository
}

func New(repo trades.Repository) *usecase {
	return &usecase{repo: repo}
}

func (u *usecase) GetTrade(ctx context.Context, id int64) (*entity.TradeWithDetails, error) {
	trade, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get trade: %w", err)
	}

	if trade == nil {
		return nil, ErrTradeNotFound
	}

	return trade, nil
}

func (u *usecase) ListMyTrades(ctx context.Context, userID int64, limit, offset int) ([]entity.TradeWithDetails, error) {
	if limit <= 0 {
		limit = 50
	}

	trades, err := u.repo.ListByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list user trades: %w", err)
	}

	return trades, nil
}

func (u *usecase) ListAllTrades(ctx context.Context, limit, offset int) ([]entity.TradeWithDetails, error) {
	if limit <= 0 {
		limit = 50
	}

	trades, err := u.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list trades: %w", err)
	}

	return trades, nil
}
