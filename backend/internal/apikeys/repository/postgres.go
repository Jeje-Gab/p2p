package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/cs2-p2p-skins/backend/internal/apikeys"
	"github.com/jmoiron/sqlx"
)

type postgresRepository struct {
	db *sqlx.DB
}

// New creates a new PostgreSQL API keys repository
func New(db *sqlx.DB) apikeys.Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, apiKey *apikeys.APIKey) error {
	query := `
		INSERT INTO api_keys (key_hash, key_prefix, user_id, name, description, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		apiKey.KeyHash,
		apiKey.KeyPrefix,
		apiKey.UserID,
		apiKey.Name,
		apiKey.Description,
		apiKey.ExpiresAt,
	).Scan(&apiKey.ID, &apiKey.CreatedAt, &apiKey.UpdatedAt)
}

func (r *postgresRepository) GetByHash(ctx context.Context, keyHash string) (*apikeys.APIKey, error) {
	var apiKey apikeys.APIKey
	query := `
		SELECT id, key_hash, key_prefix, user_id, name, description, is_active,
		       last_used_at, expires_at, created_at, updated_at
		FROM api_keys
		WHERE key_hash = $1
	`

	err := r.db.GetContext(ctx, &apiKey, query, keyHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("API key not found")
		}
		return nil, err
	}

	return &apiKey, nil
}

func (r *postgresRepository) GetByUserID(ctx context.Context, userID int64) ([]*apikeys.APIKey, error) {
	var keys []*apikeys.APIKey
	query := `
		SELECT id, key_hash, key_prefix, user_id, name, description, is_active,
		       last_used_at, expires_at, created_at, updated_at
		FROM api_keys
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &keys, query, userID)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (*apikeys.APIKey, error) {
	var apiKey apikeys.APIKey
	query := `
		SELECT id, key_hash, key_prefix, user_id, name, description, is_active,
		       last_used_at, expires_at, created_at, updated_at
		FROM api_keys
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &apiKey, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("API key not found")
		}
		return nil, err
	}

	return &apiKey, nil
}

func (r *postgresRepository) UpdateLastUsed(ctx context.Context, id int64) error {
	query := `
		UPDATE api_keys
		SET last_used_at = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (r *postgresRepository) Revoke(ctx context.Context, id int64, userID int64) error {
	query := `
		UPDATE api_keys
		SET is_active = false, updated_at = $1
		WHERE id = $2 AND user_id = $3
	`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("API key not found or unauthorized")
	}

	return nil
}

func (r *postgresRepository) Delete(ctx context.Context, id int64, userID int64) error {
	query := `
		DELETE FROM api_keys
		WHERE id = $1 AND user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("API key not found or unauthorized")
	}

	return nil
}
