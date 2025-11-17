package apikeys

import "context"

// UseCase defines the business logic interface for API keys
type UseCase interface {
	// Create generates a new API key for a user
	Create(ctx context.Context, userID int64, req CreateAPIKeyRequest) (*APIKeyWithSecret, error)

	// List returns all API keys for a user
	List(ctx context.Context, userID int64) ([]*APIKey, error)

	// Revoke deactivates an API key
	Revoke(ctx context.Context, userID int64, keyID int64) error

	// Delete permanently removes an API key
	Delete(ctx context.Context, userID int64, keyID int64) error

	// Validate checks if an API key is valid and active
	Validate(ctx context.Context, key string) (*APIKey, error)
}
