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

func (r *repository) Create(ctx context.Context, offer *entity.Offer) error {
	query := `
		INSERT INTO offers (owner_id, skin_offered_id, skin_requested_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	now := time.Now()
	err := r.db.QueryRowContext(
		ctx, query,
		offer.OwnerID, offer.SkinOfferedID, offer.SkinRequestedID,
		offer.Status, now, now,
	).Scan(&offer.ID)

	if err != nil {
		return fmt.Errorf("failed to create offer: %w", err)
	}

	offer.CreatedAt = now
	offer.UpdatedAt = now
	return nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*entity.OfferWithDetails, error) {
	query := `
		SELECT
			o.id, o.owner_id, o.skin_offered_id, o.skin_requested_id, o.status, o.created_at, o.updated_at,
			u.email as owner_email,
			s1.id as "skin_offered.id", s1.name as "skin_offered.name",
			s1.weapon_type as "skin_offered.weapon_type", s1.rarity as "skin_offered.rarity",
			s1.image_url as "skin_offered.image_url",
			s2.id as "skin_requested.id", s2.name as "skin_requested.name",
			s2.weapon_type as "skin_requested.weapon_type", s2.rarity as "skin_requested.rarity",
			s2.image_url as "skin_requested.image_url"
		FROM offers o
		JOIN users u ON o.owner_id = u.id
		JOIN skins s1 ON o.skin_offered_id = s1.id
		JOIN skins s2 ON o.skin_requested_id = s2.id
		WHERE o.id = $1 AND o.deleted_at IS NULL
	`

	var offer entity.OfferWithDetails
	if err := r.db.GetContext(ctx, &offer, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get offer: %w", err)
	}

	return &offer, nil
}

func (r *repository) List(ctx context.Context, status string, limit, offset int) ([]entity.OfferWithDetails, error) {
	query := `
		SELECT
			o.id, o.owner_id, o.skin_offered_id, o.skin_requested_id, o.status, o.created_at, o.updated_at,
			u.email as owner_email,
			s1.id as "skin_offered.id", s1.name as "skin_offered.name",
			s1.weapon_type as "skin_offered.weapon_type", s1.rarity as "skin_offered.rarity",
			s1.image_url as "skin_offered.image_url",
			s2.id as "skin_requested.id", s2.name as "skin_requested.name",
			s2.weapon_type as "skin_requested.weapon_type", s2.rarity as "skin_requested.rarity",
			s2.image_url as "skin_requested.image_url"
		FROM offers o
		JOIN users u ON o.owner_id = u.id
		JOIN skins s1 ON o.skin_offered_id = s1.id
		JOIN skins s2 ON o.skin_requested_id = s2.id
		WHERE o.deleted_at IS NULL
	`

	if status != "" {
		query += " AND o.status = $1"
	}

	query += " ORDER BY o.created_at DESC LIMIT $2 OFFSET $3"

	var offers []entity.OfferWithDetails
	var err error

	if status != "" {
		err = r.db.SelectContext(ctx, &offers, query, status, limit, offset)
	} else {
		err = r.db.SelectContext(ctx, &offers, query, limit, offset)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list offers: %w", err)
	}

	return offers, nil
}

func (r *repository) ListByUser(ctx context.Context, userID int64, limit, offset int) ([]entity.OfferWithDetails, error) {
	query := `
		SELECT
			o.id, o.owner_id, o.skin_offered_id, o.skin_requested_id, o.status, o.created_at, o.updated_at,
			u.email as owner_email,
			s1.id as "skin_offered.id", s1.name as "skin_offered.name",
			s1.weapon_type as "skin_offered.weapon_type", s1.rarity as "skin_offered.rarity",
			s1.image_url as "skin_offered.image_url",
			s2.id as "skin_requested.id", s2.name as "skin_requested.name",
			s2.weapon_type as "skin_requested.weapon_type", s2.rarity as "skin_requested.rarity",
			s2.image_url as "skin_requested.image_url"
		FROM offers o
		JOIN users u ON o.owner_id = u.id
		JOIN skins s1 ON o.skin_offered_id = s1.id
		JOIN skins s2 ON o.skin_requested_id = s2.id
		WHERE o.owner_id = $1 AND o.deleted_at IS NULL
		ORDER BY o.created_at DESC
		LIMIT $2 OFFSET $3
	`

	var offers []entity.OfferWithDetails
	if err := r.db.SelectContext(ctx, &offers, query, userID, limit, offset); err != nil {
		return nil, fmt.Errorf("failed to list user offers: %w", err)
	}

	return offers, nil
}

func (r *repository) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE offers SET status = $1, updated_at = $2 WHERE id = $3 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update offer status: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("offer not found")
	}

	return nil
}
