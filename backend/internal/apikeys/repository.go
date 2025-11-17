package apikeys

import "context"

// Repository defines the interface for API key data access
type Repository interface {
	// Create creates a new API key
	Create(ctx context.Context, apiKey *APIKey) error

	// GetByHash retrieves an API key by its hash
	GetByHash(ctx context.Context, keyHash string) (*APIKey, error)

	// GetByUserID retrieves all API keys for a user
	GetByUserID(ctx context.Context, userID int64) ([]*APIKey, error)

	// GetByID retrieves an API key by ID
	GetByID(ctx context.Context, id int64) (*APIKey, error)

	// UpdateLastUsed updates the last used timestamp
	UpdateLastUsed(ctx context.Context, id int64) error

	// Revoke deactivates an API key
	Revoke(ctx context.Context, id int64, userID int64) error

	// Delete permanently deletes an API key
	Delete(ctx context.Context, id int64, userID int64) error
}
