package trades

import (
	"context"

	"github.com/cs2-p2p-skins/backend/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, trade *entity.Trade) error
	GetByID(ctx context.Context, id int64) (*entity.TradeWithDetails, error)
	ListByUser(ctx context.Context, userID int64, limit, offset int) ([]entity.TradeWithDetails, error)
	List(ctx context.Context, limit, offset int) ([]entity.TradeWithDetails, error)
}

type UseCase interface {
	GetTrade(ctx context.Context, id int64) (*entity.TradeWithDetails, error)
	ListMyTrades(ctx context.Context, userID int64, limit, offset int) ([]entity.TradeWithDetails, error)
	ListAllTrades(ctx context.Context, limit, offset int) ([]entity.TradeWithDetails, error)
}
