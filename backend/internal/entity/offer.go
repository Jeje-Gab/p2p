package entity

import "time"

// Offer status types
const (
	OfferStatusOpen      = "open"
	OfferStatusAccepted  = "accepted"
	OfferStatusCancelled = "cancelled"
)

// Offer represents a P2P trade offer
type Offer struct {
	ID              int64      `json:"id" db:"id"`
	OwnerID         int64      `json:"owner_id" db:"owner_id"`
	SkinOfferedID   int64      `json:"skin_offered_id" db:"skin_offered_id"`
	SkinRequestedID int64      `json:"skin_requested_id" db:"skin_requested_id"`
	Status          string     `json:"status" db:"status"` // open, accepted, cancelled
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// OfferWithDetails extends Offer with full skin and user information
type OfferWithDetails struct {
	Offer
	OwnerEmail   string `json:"owner_email"`
	SkinOffered  Skin   `json:"skin_offered"`
	SkinRequested Skin  `json:"skin_requested"`
}
