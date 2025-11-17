package auth

import (
	"context"

	"github.com/cs2-p2p-skins/backend/internal/entity"
)

// UseCase defines the auth business logic interface
type UseCase interface {
	// Authentication
	Register(ctx context.Context, email, password string) (*entity.User, error)
	Login(ctx context.Context, email, password, ip string) (requiresCode bool, token string, user *entity.User, err error)
	Verify2FA(ctx context.Context, email, code, ip string) (token string, user *entity.User, err error)
	GetCurrentUser(ctx context.Context, userID int64) (*entity.User, error)

	// 2FA Management
	Setup2FA(ctx context.Context, userID int64) (secret, qrURL string, err error)
	Enable2FA(ctx context.Context, userID int64, code string) (backupCodes []string, err error)
	Disable2FA(ctx context.Context, userID int64, code string) error
	Verify2FAWithBackupCode(ctx context.Context, email, code, ip string) (token string, user *entity.User, err error)
	Verify2FAForOAuth(ctx context.Context, email, code, ip string) (token string, user *entity.User, err error)

	// Steam OAuth
	GenerateSteamAuthURL(ctx context.Context) (authURL, state string, err error)
	HandleSteamCallback(ctx context.Context, values map[string]string) (requires2FA bool, token string, user *entity.User, err error)
}
