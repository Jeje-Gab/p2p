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
	ID              int64     `json:"id" db:"id"`
	OwnerID         int64     `json:"owner_id" db:"owner_id"`
	SkinOfferedID   int64     `json:"skin_offered_id" db:"skin_offered_id"`
	SkinRequestedID int64     `json:"skin_requested_id" db:"skin_requested_id"`
	Status          string    `json:"status" db:"status"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	OwnerEmail      string    `json:"owner_email" db:"owner_email"`

	// Flattened skin fields for DB scanning
	SkinOfferedIDField       int64  `json:"-" db:"skin_offered.id"`
	SkinOfferedName          string `json:"-" db:"skin_offered.name"`
	SkinOfferedWeaponType    string `json:"-" db:"skin_offered.weapon_type"`
	SkinOfferedRarity        string `json:"-" db:"skin_offered.rarity"`
	SkinOfferedImageURL      string `json:"-" db:"skin_offered.image_url"`

	SkinRequestedIDField     int64  `json:"-" db:"skin_requested.id"`
	SkinRequestedName        string `json:"-" db:"skin_requested.name"`
	SkinRequestedWeaponType  string `json:"-" db:"skin_requested.weapon_type"`
	SkinRequestedRarity      string `json:"-" db:"skin_requested.rarity"`
	SkinRequestedImageURL    string `json:"-" db:"skin_requested.image_url"`

	// JSON output fields
	SkinOffered   Skin `json:"skin_offered" db:"-"`
	SkinRequested Skin `json:"skin_requested" db:"-"`
}

// PopulateSkins populates the nested Skin structs from flattened fields
func (o *OfferWithDetails) PopulateSkins() {
	o.SkinOffered = Skin{
		ID:         o.SkinOfferedIDField,
		Name:       o.SkinOfferedName,
		WeaponType: o.SkinOfferedWeaponType,
		Rarity:     o.SkinOfferedRarity,
		ImageURL:   o.SkinOfferedImageURL,
	}
	o.SkinRequested = Skin{
		ID:         o.SkinRequestedIDField,
		Name:       o.SkinRequestedName,
		WeaponType: o.SkinRequestedWeaponType,
		Rarity:     o.SkinRequestedRarity,
		ImageURL:   o.SkinRequestedImageURL,
	}
}
