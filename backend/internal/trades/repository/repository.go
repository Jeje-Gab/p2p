package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cs2-p2p-skins/backend/internal/entity"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, trade *entity.Trade) error {
	query := `
		INSERT INTO trades (offer_id, from_user_id, to_user_id, executed_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	now := time.Now()
	err := r.db.QueryRowContext(
		ctx, query,
		trade.OfferID, trade.FromUserID, trade.ToUserID, now,
	).Scan(&trade.ID)

	if err != nil {
		return fmt.Errorf("failed to create trade: %w", err)
	}

	trade.ExecutedAt = now
	return nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*entity.TradeWithDetails, error) {
	query := `
		SELECT
			t.id, t.offer_id, t.from_user_id, t.to_user_id, t.executed_at,
			u1.email as from_user_email,
			u2.email as to_user_email
		FROM trades t
		JOIN users u1 ON t.from_user_id = u1.id
		JOIN users u2 ON t.to_user_id = u2.id
		WHERE t.id = $1
	`

	var trade entity.TradeWithDetails
	if err := r.db.GetContext(ctx, &trade, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get trade: %w", err)
	}

	return &trade, nil
}

func (r *repository) ListByUser(ctx context.Context, userID int64, limit, offset int) ([]entity.TradeWithDetails, error) {
	query := `
		SELECT
			t.id, t.offer_id, t.from_user_id, t.to_user_id, t.executed_at,
			u1.email as from_user_email,
			u2.email as to_user_email
		FROM trades t
		JOIN users u1 ON t.from_user_id = u1.id
		JOIN users u2 ON t.to_user_id = u2.id
		WHERE t.from_user_id = $1 OR t.to_user_id = $1
		ORDER BY t.executed_at DESC
		LIMIT $2 OFFSET $3
	`

	var trades []entity.TradeWithDetails
	if err := r.db.SelectContext(ctx, &trades, query, userID, limit, offset); err != nil {
		return nil, fmt.Errorf("failed to list user trades: %w", err)
	}

	return trades, nil
}

func (r *repository) List(ctx context.Context, limit, offset int) ([]entity.TradeWithDetails, error) {
	query := `
		SELECT
			t.id, t.offer_id, t.from_user_id, t.to_user_id, t.executed_at,
			u1.email as from_user_email,
			u2.email as to_user_email
		FROM trades t
		JOIN users u1 ON t.from_user_id = u1.id
		JOIN users u2 ON t.to_user_id = u2.id
		ORDER BY t.executed_at DESC
		LIMIT $1 OFFSET $2
	`

	var trades []entity.TradeWithDetails
	if err := r.db.SelectContext(ctx, &trades, query, limit, offset); err != nil {
		return nil, fmt.Errorf("failed to list trades: %w", err)
	}

	return trades, nil
}
