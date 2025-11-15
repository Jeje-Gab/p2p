package entity

import "time"

// Skin rarity levels
const (
	RarityConsumerGrade = "Consumer Grade"
	RarityIndustrialGrade = "Industrial Grade"
	RarityMilSpecGrade    = "Mil-Spec Grade"
	RarityRestricted      = "Restricted"
	RarityClassified      = "Classified"
	RarityCovert          = "Covert"
	RarityExceedinglyRare = "Exceedingly Rare"
)

// Skin represents a CS2 weapon skin
type Skin struct {
	ID         int64      `json:"id" db:"id"`
	Name       string     `json:"name" db:"name"`
	WeaponType string     `json:"weapon_type" db:"weapon_type"` // e.g., "AK-47", "AWP", "M4A4"
	Rarity     string     `json:"rarity" db:"rarity"`
	FloatValue *float64   `json:"float_value,omitempty" db:"float_value"` // Wear value (0.0-1.0)
	ImageURL   string     `json:"image_url" db:"image_url"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// UserSkin represents the relationship between users and their skins
type UserSkin struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	SkinID    int64     `json:"skin_id" db:"skin_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserSkinWithDetails extends UserSkin with full skin information
type UserSkinWithDetails struct {
	UserSkin
	Skin Skin `json:"skin"`
}
