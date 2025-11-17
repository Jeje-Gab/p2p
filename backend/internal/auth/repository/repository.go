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

// New creates a new auth repository
func New(db *sqlx.DB) *repository {
	return &repository{db: db}
}

// CreateUser inserts a new user into the database
func (r *repository) CreateUser(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO p2p.users (email, password_hash, steam_id, role, twofa_enabled, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	now := time.Now()
	err := r.db.QueryRowContext(
		ctx, query,
		user.Email,
		user.PasswordHash,
		user.SteamID,
		user.Role,
		user.TwoFAEnabled,
		now,
		now,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

// GetUserByEmail retrieves a user by email
func (r *repository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	query := `
		SELECT id, email, password_hash, steam_id, twofa_secret, twofa_enabled, role, created_at, updated_at
		FROM p2p.users
		WHERE email = $1 AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID
func (r *repository) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	query := `
		SELECT id, email, password_hash, steam_id, twofa_secret, twofa_enabled, role, created_at, updated_at
		FROM p2p.users
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

// GetUserBySteamID retrieves a user by Steam ID
func (r *repository) GetUserBySteamID(ctx context.Context, steamID string) (*entity.User, error) {
	var user entity.User
	query := `
		SELECT id, email, password_hash, steam_id, twofa_secret, twofa_enabled, role, created_at, updated_at
		FROM p2p.users
		WHERE steam_id = $1 AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, &user, query, steamID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by Steam ID: %w", err)
	}

	return &user, nil
}

// UpdateUser updates user information
func (r *repository) UpdateUser(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE p2p.users
		SET email = $1, password_hash = $2, steam_id = $3, twofa_secret = $4,
		    twofa_enabled = $5, role = $6, updated_at = $7
		WHERE id = $8 AND deleted_at IS NULL
	`

	user.UpdatedAt = time.Now()
	result, err := r.db.ExecContext(
		ctx, query,
		user.Email,
		user.PasswordHash,
		user.SteamID,
		user.TwoFASecret,
		user.TwoFAEnabled,
		user.Role,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}

// Enable2FA enables 2FA for a user
func (r *repository) Enable2FA(ctx context.Context, userID int64, secret string) error {
	query := `
		UPDATE p2p.users
		SET twofa_secret = $1, twofa_enabled = true, updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, secret, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to enable 2FA: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}

// Disable2FA disables 2FA for a user
func (r *repository) Disable2FA(ctx context.Context, userID int64) error {
	query := `
		UPDATE p2p.users
		SET twofa_secret = NULL, twofa_enabled = false, updated_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to disable 2FA: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}

// CreateLoginAttempt records a login attempt (for brute force tracking)
func (r *repository) CreateLoginAttempt(ctx context.Context, attempt *entity.LoginAttempt) error {
	query := `
		INSERT INTO p2p.login_attempts (user_id, ip_address, success, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	now := time.Now()
	err := r.db.QueryRowContext(
		ctx, query,
		attempt.UserID,
		attempt.IPAddress,
		attempt.Success,
		now,
	).Scan(&attempt.ID)

	if err != nil {
		return fmt.Errorf("failed to create login attempt: %w", err)
	}

	attempt.CreatedAt = now
	return nil
}

// CountFailedAttempts counts failed login attempts for a user within the last N minutes
// Security: Used to detect brute force attacks
func (r *repository) CountFailedAttempts(ctx context.Context, email string, sinceMinutes int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM p2p.login_attempts la
		JOIN p2p.users u ON la.user_id = u.id
		WHERE u.email = $1
		  AND la.success = false
		  AND la.created_at > NOW() - INTERVAL '1 minute' * $2
	`

	var count int
	err := r.db.GetContext(ctx, &count, query, email, sinceMinutes)
	if err != nil {
		return 0, fmt.Errorf("failed to count failed attempts: %w", err)
	}

	return count, nil
}

// CountFailedAttemptsByIP counts failed login attempts from an IP within the last N minutes
// Security: Prevents distributed brute force attacks from same IP
func (r *repository) CountFailedAttemptsByIP(ctx context.Context, ip string, sinceMinutes int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM p2p.login_attempts
		WHERE ip_address = $1
		  AND success = false
		  AND created_at > NOW() - INTERVAL '1 minute' * $2
	`

	var count int
	err := r.db.GetContext(ctx, &count, query, ip, sinceMinutes)
	if err != nil {
		return 0, fmt.Errorf("failed to count failed attempts by IP: %w", err)
	}

	return count, nil
}

// CreateBackupCodes creates backup codes for a user
// Security: Codes are stored hashed and can only be used once
func (r *repository) CreateBackupCodes(ctx context.Context, userID int64, codeHashes []string) error {
	// Use a transaction to ensure all codes are created atomically
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO p2p.backup_codes (user_id, code_hash, created_at)
		VALUES ($1, $2, $3)
	`

	now := time.Now()
	for _, hash := range codeHashes {
		_, err := tx.ExecContext(ctx, query, userID, hash, now)
		if err != nil {
			return fmt.Errorf("failed to create backup code: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetUnusedBackupCodes retrieves all unused backup codes for a user
func (r *repository) GetUnusedBackupCodes(ctx context.Context, userID int64) ([]*entity.BackupCode, error) {
	query := `
		SELECT id, user_id, code_hash, used, used_at, created_at
		FROM p2p.backup_codes
		WHERE user_id = $1 AND used = false
		ORDER BY created_at ASC
	`

	var codes []*entity.BackupCode
	err := r.db.SelectContext(ctx, &codes, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get unused backup codes: %w", err)
	}

	return codes, nil
}

// MarkBackupCodeAsUsed marks a backup code as used
// Security: Once used, codes cannot be reused (prevents replay attacks)
func (r *repository) MarkBackupCodeAsUsed(ctx context.Context, codeID int64) error {
	query := `
		UPDATE p2p.backup_codes
		SET used = true, used_at = $1
		WHERE id = $2 AND used = false
	`

	result, err := r.db.ExecContext(ctx, query, time.Now(), codeID)
	if err != nil {
		return fmt.Errorf("failed to mark backup code as used: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return errors.New("backup code not found or already used")
	}

	return nil
}

// DeleteBackupCodes deletes all backup codes for a user
func (r *repository) DeleteBackupCodes(ctx context.Context, userID int64) error {
	query := `
		DELETE FROM p2p.backup_codes
		WHERE user_id = $1
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete backup codes: %w", err)
	}

	return nil
}

// SaveTwoFASecret saves the 2FA secret without enabling 2FA
// This is used during setup phase before user verifies the code
func (r *repository) SaveTwoFASecret(ctx context.Context, userID int64, secret string) error {
	query := `
		UPDATE p2p.users
		SET twofa_secret = $1, updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, secret, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to save 2FA secret: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}
