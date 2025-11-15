package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cs2-p2p-skins/backend/internal/entity"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) ListSkins(ctx context.Context, limit, offset int) ([]entity.Skin, error) {
	query := `
		SELECT id, name, weapon_type, rarity, float_value, image_url, created_at, updated_at
		FROM skins
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	var skins []entity.Skin
	if err := r.db.SelectContext(ctx, &skins, query, limit, offset); err != nil {
		return nil, fmt.Errorf("failed to list skins: %w", err)
	}

	return skins, nil
}

func (r *repository) GetSkinByID(ctx context.Context, id int64) (*entity.Skin, error) {
	query := `
		SELECT id, name, weapon_type, rarity, float_value, image_url, created_at, updated_at
		FROM skins
		WHERE id = $1 AND deleted_at IS NULL
	`

	var skin entity.Skin
	if err := r.db.GetContext(ctx, &skin, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get skin: %w", err)
	}

	return &skin, nil
}

func (r *repository) GetUserSkins(ctx context.Context, userID int64) ([]entity.UserSkinWithDetails, error) {
	query := `
		SELECT
			us.id, us.user_id, us.skin_id, us.quantity, us.created_at, us.updated_at,
			s.id as "skin.id", s.name as "skin.name", s.weapon_type as "skin.weapon_type",
			s.rarity as "skin.rarity", s.float_value as "skin.float_value", s.image_url as "skin.image_url",
			s.created_at as "skin.created_at", s.updated_at as "skin.updated_at"
		FROM user_skins us
		JOIN skins s ON us.skin_id = s.id
		WHERE us.user_id = $1 AND s.deleted_at IS NULL
		ORDER BY us.created_at DESC
	`

	var userSkins []entity.UserSkinWithDetails
	if err := r.db.SelectContext(ctx, &userSkins, query, userID); err != nil {
		return nil, fmt.Errorf("failed to get user skins: %w", err)
	}

	return userSkins, nil
}

func (r *repository) AddSkinToUser(ctx context.Context, userID, skinID int64, quantity int) error {
	query := `
		INSERT INTO user_skins (user_id, skin_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, skin_id)
		DO UPDATE SET quantity = user_skins.quantity + $3, updated_at = NOW()
	`

	if _, err := r.db.ExecContext(ctx, query, userID, skinID, quantity); err != nil {
		return fmt.Errorf("failed to add skin to user: %w", err)
	}

	return nil
}

func (r *repository) RemoveSkinFromUser(ctx context.Context, userID, skinID int64, quantity int) error {
	query := `
		UPDATE user_skins
		SET quantity = quantity - $3, updated_at = NOW()
		WHERE user_id = $1 AND skin_id = $2 AND quantity >= $3
	`

	result, err := r.db.ExecContext(ctx, query, userID, skinID, quantity)
	if err != nil {
		return fmt.Errorf("failed to remove skin: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("insufficient quantity or skin not found")
	}

	// Delete if quantity becomes 0
	deleteQuery := `DELETE FROM user_skins WHERE user_id = $1 AND skin_id = $2 AND quantity <= 0`
	_, _ = r.db.ExecContext(ctx, deleteQuery, userID, skinID)

	return nil
}

func (r *repository) HasSkin(ctx context.Context, userID, skinID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_skins WHERE user_id = $1 AND skin_id = $2 AND quantity > 0)`

	var exists bool
	if err := r.db.GetContext(ctx, &exists, query, userID, skinID); err != nil {
		return false, fmt.Errorf("failed to check skin ownership: %w", err)
	}

	return exists, nil
}
