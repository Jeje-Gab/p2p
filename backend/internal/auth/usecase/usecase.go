package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/cs2-p2p-skins/backend/internal/auth"
	"github.com/cs2-p2p-skins/backend/internal/entity"
	pkgAuth "github.com/cs2-p2p-skins/backend/pkg/auth"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrAccountLocked      = errors.New("account temporarily locked due to too many failed attempts")
	ErrInvalid2FACode     = errors.New("invalid 2FA code")
	Err2FANotEnabled      = errors.New("2FA not enabled for this user")
	Err2FAAlreadyEnabled  = errors.New("2FA already enabled")
)

const (
	// Security: Brute force protection settings
	maxFailedAttempts = 5
	lockoutWindow     = 15 // minutes
)

type usecase struct {
	repo        auth.Repository
	jwtManager  *pkgAuth.JWTManager
	totpManager *pkgAuth.TOTPManager
	steamAuth   *pkgAuth.SteamOAuthManager
}

// New creates a new auth usecase
func New(
	repo auth.Repository,
	jwtManager *pkgAuth.JWTManager,
	totpManager *pkgAuth.TOTPManager,
	steamAuth *pkgAuth.SteamOAuthManager,
) *usecase {
	return &usecase{
		repo:        repo,
		jwtManager:  jwtManager,
		totpManager: totpManager,
		steamAuth:   steamAuth,
	}
}

// Register creates a new user account
// Security: Password is hashed with bcrypt before storage
func (u *usecase) Register(ctx context.Context, email, password string) (*entity.User, error) {
	// Check if user already exists
	existingUser, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	if existingUser != nil {
		return nil, ErrUserExists
	}

	// Hash password - NEVER store plain passwords!
	hashedPassword, err := pkgAuth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &entity.User{
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         entity.RoleUser,
		TwoFAEnabled: false,
	}

	if err := u.repo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login authenticates a user (first factor: email + password)
// Security: Implements brute force protection by tracking failed attempts
func (u *usecase) Login(ctx context.Context, email, password, ip string) (requiresTwoFA bool, token string, user *entity.User, err error) {
	// Brute force protection: Check failed attempts by email
	failedAttempts, err := u.repo.CountFailedAttempts(ctx, email, lockoutWindow)
	if err != nil {
		return false, "", nil, fmt.Errorf("failed to check failed attempts: %w", err)
	}

	if failedAttempts >= maxFailedAttempts {
		// Account is temporarily locked
		return false, "", nil, ErrAccountLocked
	}

	// Brute force protection: Check failed attempts by IP
	failedAttemptsIP, err := u.repo.CountFailedAttemptsByIP(ctx, ip, lockoutWindow)
	if err != nil {
		return false, "", nil, fmt.Errorf("failed to check failed attempts by IP: %w", err)
	}

	if failedAttemptsIP >= maxFailedAttempts*2 {
		// IP is temporarily blocked
		return false, "", nil, ErrAccountLocked
	}

	// Get user
	loggedUser, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return false, "", nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Record failed attempt if user not found
	if loggedUser == nil {
		_ = u.repo.CreateLoginAttempt(ctx, &entity.LoginAttempt{
			IPAddress: ip,
			Success:   false,
			CreatedAt: time.Now(),
		})
		return false, "", nil, ErrInvalidCredentials
	}

	// Verify password
	if !pkgAuth.VerifyPassword(loggedUser.PasswordHash, password) {
		// Record failed attempt
		_ = u.repo.CreateLoginAttempt(ctx, &entity.LoginAttempt{
			UserID:    &loggedUser.ID,
			IPAddress: ip,
			Success:   false,
			CreatedAt: time.Now(),
		})
		return false, "", nil, ErrInvalidCredentials
	}

	// If 2FA is enabled, don't generate token yet
	if loggedUser.TwoFAEnabled {
		return true, "", nil, nil
	}

	// Record successful attempt
	_ = u.repo.CreateLoginAttempt(ctx, &entity.LoginAttempt{
		UserID:    &loggedUser.ID,
		IPAddress: ip,
		Success:   true,
		CreatedAt: time.Now(),
	})

	// Generate JWT token
	token, err = u.jwtManager.GenerateToken(loggedUser.ID, loggedUser.Email, loggedUser.Role)
	if err != nil {
		return false, "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return false, token, loggedUser, nil
}

// Verify2FA verifies the 2FA code and generates JWT token
// Security: Second factor authentication significantly increases account security
func (u *usecase) Verify2FA(ctx context.Context, email, code, ip string) (string, *entity.User, error) {
	// Get user
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return "", nil, ErrUserNotFound
	}

	if !user.TwoFAEnabled || user.TwoFASecret == nil {
		return "", nil, Err2FANotEnabled
	}

	// Validate TOTP code
	if !u.totpManager.ValidateCode(*user.TwoFASecret, code) {
		// Record failed attempt
		_ = u.repo.CreateLoginAttempt(ctx, &entity.LoginAttempt{
			UserID:    &user.ID,
			IPAddress: ip,
			Success:   false,
			CreatedAt: time.Now(),
		})
		return "", nil, ErrInvalid2FACode
	}

	// Record successful attempt
	_ = u.repo.CreateLoginAttempt(ctx, &entity.LoginAttempt{
		UserID:    &user.ID,
		IPAddress: ip,
		Success:   true,
		CreatedAt: time.Now(),
	})

	// Generate JWT token
	token, err := u.jwtManager.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, nil
}

// GetCurrentUser retrieves the current authenticated user
func (u *usecase) GetCurrentUser(ctx context.Context, userID int64) (*entity.User, error) {
	user, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// Setup2FA generates a new 2FA secret for a user
// Returns the secret and QR code URL for Google Authenticator
func (u *usecase) Setup2FA(ctx context.Context, userID int64) (string, string, error) {
	user, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return "", "", ErrUserNotFound
	}

	if user.TwoFAEnabled {
		return "", "", Err2FAAlreadyEnabled
	}

	// Generate TOTP secret
	secret, qrURL, err := u.totpManager.GenerateSecret(user.Email)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate 2FA secret: %w", err)
	}

	return secret, qrURL, nil
}

// Enable2FA enables 2FA for a user after verifying the code
func (u *usecase) Enable2FA(ctx context.Context, userID int64, code string) error {
	user, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return ErrUserNotFound
	}

	if user.TwoFAEnabled {
		return Err2FAAlreadyEnabled
	}

	// Generate and validate secret
	secret, _, err := u.totpManager.GenerateSecret(user.Email)
	if err != nil {
		return fmt.Errorf("failed to generate secret: %w", err)
	}

	// Validate code before enabling
	if !u.totpManager.ValidateCode(secret, code) {
		return ErrInvalid2FACode
	}

	// Enable 2FA
	if err := u.repo.Enable2FA(ctx, userID, secret); err != nil {
		return fmt.Errorf("failed to enable 2FA: %w", err)
	}

	return nil
}

// Disable2FA disables 2FA for a user after verifying the code
func (u *usecase) Disable2FA(ctx context.Context, userID int64, code string) error {
	user, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return ErrUserNotFound
	}

	if !user.TwoFAEnabled || user.TwoFASecret == nil {
		return Err2FANotEnabled
	}

	// Validate code before disabling
	if !u.totpManager.ValidateCode(*user.TwoFASecret, code) {
		return ErrInvalid2FACode
	}

	// Disable 2FA
	if err := u.repo.Disable2FA(ctx, userID); err != nil {
		return fmt.Errorf("failed to disable 2FA: %w", err)
	}

	return nil
}

// GenerateSteamAuthURL generates the Steam OAuth login URL
// Security: Uses state parameter for CSRF protection
func (u *usecase) GenerateSteamAuthURL(ctx context.Context) (string, string, error) {
	authURL, state, err := u.steamAuth.GenerateAuthURL()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate Steam auth URL: %w", err)
	}

	return authURL, state, nil
}

// HandleSteamCallback processes the Steam OAuth callback
// Creates user if doesn't exist, generates JWT token
func (u *usecase) HandleSteamCallback(ctx context.Context, values map[string]string) (string, error) {
	// Convert map to url.Values
	urlValues := url.Values{}
	for k, v := range values {
		urlValues.Set(k, v)
	}

	// Validate Steam response
	steamUser, err := u.steamAuth.ValidateCallback(ctx, urlValues)
	if err != nil {
		return "", fmt.Errorf("failed to validate Steam callback: %w", err)
	}

	// Check if user exists
	user, err := u.repo.GetUserBySteamID(ctx, steamUser.SteamID)
	if err != nil {
		return "", fmt.Errorf("failed to get user by Steam ID: %w", err)
	}

	// Create user if doesn't exist
	if user == nil {
		user = &entity.User{
			Email:        fmt.Sprintf("steam_%s@steamuser.local", steamUser.SteamID),
			PasswordHash: "", // No password for OAuth users
			SteamID:      &steamUser.SteamID,
			Role:         entity.RoleUser,
			TwoFAEnabled: false,
		}

		if err := u.repo.CreateUser(ctx, user); err != nil {
			return "", fmt.Errorf("failed to create Steam user: %w", err)
		}
	}

	// Generate JWT token
	token, err := u.jwtManager.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}
