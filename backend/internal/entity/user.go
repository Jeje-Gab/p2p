package entity

import "time"

// User role types
const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

// User represents a registered user in the system
type User struct {
	ID           int64      `json:"id" db:"id"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"` // Never expose password hash in JSON
	SteamID      *string    `json:"steam_id,omitempty" db:"steam_id"`
	TwoFASecret  *string    `json:"-" db:"twofa_secret"` // Never expose 2FA secret
	TwoFAEnabled bool       `json:"twofa_enabled" db:"twofa_enabled"`
	Role         string     `json:"role" db:"role"` // user or admin
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// LoginAttempt tracks login attempts for brute force protection
type LoginAttempt struct {
	ID        int64     `json:"id" db:"id"`
	UserID    *int64    `json:"user_id,omitempty" db:"user_id"` // Nullable if user doesn't exist
	IPAddress string    `json:"ip_address" db:"ip_address"`
	Success   bool      `json:"success" db:"success"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// BackupCode represents a one-time use 2FA recovery code
type BackupCode struct {
	ID        int64      `json:"id" db:"id"`
	UserID    int64      `json:"user_id" db:"user_id"`
	CodeHash  string     `json:"-" db:"code_hash"` // Never expose hash
	Used      bool       `json:"used" db:"used"`
	UsedAt    *time.Time `json:"used_at,omitempty" db:"used_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}
