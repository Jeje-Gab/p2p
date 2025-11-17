package apikeys

import "time"

// APIKey represents an external API key for third-party integrations
type APIKey struct {
	ID          int64      `json:"id" db:"id"`
	KeyHash     string     `json:"-" db:"key_hash"` // Never expose the hash
	KeyPrefix   string     `json:"key_prefix" db:"key_prefix"`
	UserID      int64      `json:"user_id" db:"user_id"`
	Name        string     `json:"name" db:"name"`
	Description *string    `json:"description,omitempty" db:"description"`
	IsActive    bool       `json:"is_active" db:"is_active"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// APIKeyWithSecret is returned only when creating a new key
// The plain key is shown ONCE and never stored
type APIKeyWithSecret struct {
	APIKey
	Key string `json:"key"` // Plain text key - show only once!
}

// CreateAPIKeyRequest represents the request to create a new API key
type CreateAPIKeyRequest struct {
	Name        string     `json:"name" validate:"required,min=3,max=255"`
	Description *string    `json:"description,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}
