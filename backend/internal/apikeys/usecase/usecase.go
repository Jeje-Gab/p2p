package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/cs2-p2p-skins/backend/internal/apikeys"
)

type useCase struct {
	repo apikeys.Repository
}

// New creates a new API keys use case
func New(repo apikeys.Repository) apikeys.UseCase {
	return &useCase{repo: repo}
}

// generateSecureKey generates a cryptographically secure random API key
// Format: sk_live_<32 hex characters>
func generateSecureKey() (string, error) {
	// Generate 32 random bytes (256 bits)
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Convert to hex string
	hexStr := hex.EncodeToString(bytes)

	// Format with prefix
	return fmt.Sprintf("sk_live_%s", hexStr), nil
}

// hashKey creates a SHA256 hash of the API key for storage
func hashKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

// getKeyPrefix extracts the first 12 characters for display
// Example: sk_live_abc12345...
func getKeyPrefix(key string) string {
	if len(key) > 20 {
		return key[:20] + "..."
	}
	return key
}

func (uc *useCase) Create(ctx context.Context, userID int64, req apikeys.CreateAPIKeyRequest) (*apikeys.APIKeyWithSecret, error) {
	// Validate request
	if req.Name == "" {
		return nil, errors.New("API key name is required")
	}

	// Check if expiration is in the future
	if req.ExpiresAt != nil && req.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("expiration date must be in the future")
	}

	// Generate secure API key
	key, err := generateSecureKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	// Hash the key for storage
	keyHash := hashKey(key)
	keyPrefix := getKeyPrefix(key)

	// Create API key entity
	apiKey := &apikeys.APIKey{
		KeyHash:     keyHash,
		KeyPrefix:   keyPrefix,
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
		ExpiresAt:   req.ExpiresAt,
	}

	// Save to database
	if err := uc.repo.Create(ctx, apiKey); err != nil {
		return nil, fmt.Errorf("failed to create API key: %w", err)
	}

	// Return with the plain key (shown only once!)
	return &apikeys.APIKeyWithSecret{
		APIKey: *apiKey,
		Key:    key,
	}, nil
}

func (uc *useCase) List(ctx context.Context, userID int64) ([]*apikeys.APIKey, error) {
	keys, err := uc.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list API keys: %w", err)
	}

	return keys, nil
}

func (uc *useCase) Revoke(ctx context.Context, userID int64, keyID int64) error {
	if err := uc.repo.Revoke(ctx, keyID, userID); err != nil {
		return fmt.Errorf("failed to revoke API key: %w", err)
	}

	return nil
}

func (uc *useCase) Delete(ctx context.Context, userID int64, keyID int64) error {
	if err := uc.repo.Delete(ctx, keyID, userID); err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}

	return nil
}

func (uc *useCase) Validate(ctx context.Context, key string) (*apikeys.APIKey, error) {
	// Hash the provided key
	keyHash := hashKey(key)

	// Look up in database
	apiKey, err := uc.repo.GetByHash(ctx, keyHash)
	if err != nil {
		return nil, errors.New("invalid API key")
	}

	// Check if active
	if !apiKey.IsActive {
		return nil, errors.New("API key has been revoked")
	}

	// Check expiration
	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("API key has expired")
	}

	// Update last used timestamp (async - don't wait)
	go func() {
		_ = uc.repo.UpdateLastUsed(context.Background(), apiKey.ID)
	}()

	return apiKey, nil
}
