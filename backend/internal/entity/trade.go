package entity

import "time"

// Trade represents a completed trade between two users
type Trade struct {
	ID         int64     `json:"id" db:"id"`
	OfferID    int64     `json:"offer_id" db:"offer_id"`
	FromUserID int64     `json:"from_user_id" db:"from_user_id"`
	ToUserID   int64     `json:"to_user_id" db:"to_user_id"`
	ExecutedAt time.Time `json:"executed_at" db:"executed_at"`
}

// TradeWithDetails extends Trade with offer and user information
type TradeWithDetails struct {
	Trade
	Offer         OfferWithDetails `json:"offer"`
	FromUserEmail string           `json:"from_user_email" db:"from_user_email"`
	ToUserEmail   string           `json:"to_user_email" db:"to_user_email"`
}
