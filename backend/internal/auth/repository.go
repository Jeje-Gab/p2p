package auth

import (
	"context"

	"github.com/cs2-p2p-skins/backend/internal/entity"
)

// Repository defines the interface for auth data access
type Repository interface {
	// User operations
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
	GetUserBySteamID(ctx context.Context, steamID string) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	Enable2FA(ctx context.Context, userID int64, secret string) error
	Disable2FA(ctx context.Context, userID int64) error

	// Login attempt operations (for brute force protection)
	CreateLoginAttempt(ctx context.Context, attempt *entity.LoginAttempt) error
	CountFailedAttempts(ctx context.Context, email string, since int) (int, error)
	CountFailedAttemptsByIP(ctx context.Context, ip string, since int) (int, error)

	// Backup code operations (for 2FA recovery)
	CreateBackupCodes(ctx context.Context, userID int64, codeHashes []string) error
	GetUnusedBackupCodes(ctx context.Context, userID int64) ([]*entity.BackupCode, error)
	MarkBackupCodeAsUsed(ctx context.Context, codeID int64) error
	DeleteBackupCodes(ctx context.Context, userID int64) error

	// 2FA setup - save secret without enabling
	SaveTwoFASecret(ctx context.Context, userID int64, secret string) error
}
