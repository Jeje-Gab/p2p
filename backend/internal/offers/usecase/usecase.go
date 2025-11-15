package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/cs2-p2p-skins/backend/internal/entity"
	"github.com/cs2-p2p-skins/backend/internal/offers"
	"github.com/cs2-p2p-skins/backend/internal/skins"
	"github.com/cs2-p2p-skins/backend/internal/trades"
)

var (
	ErrOfferNotFound  = errors.New("offer not found")
	ErrNotAuthorized  = errors.New("not authorized")
	ErrOfferNotOpen   = errors.New("offer is not open")
	ErrInsufficientSkins = errors.New("insufficient skins")
)

type usecase struct {
	offersRepo offers.Repository
	skinsRepo  skins.Repository
	tradesRepo trades.Repository
}

func New(offersRepo offers.Repository, skinsRepo skins.Repository, tradesRepo trades.Repository) *usecase {
	return &usecase{
		offersRepo: offersRepo,
		skinsRepo:  skinsRepo,
		tradesRepo: tradesRepo,
	}
}

func (u *usecase) CreateOffer(ctx context.Context, ownerID, skinOfferedID, skinRequestedID int64) (*entity.Offer, error) {
	// Verify owner has the skin
	hasSkin, err := u.skinsRepo.HasSkin(ctx, ownerID, skinOfferedID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify skin ownership: %w", err)
	}

	if !hasSkin {
		return nil, ErrInsufficientSkins
	}

	offer := &entity.Offer{
		OwnerID:         ownerID,
		SkinOfferedID:   skinOfferedID,
		SkinRequestedID: skinRequestedID,
		Status:          entity.OfferStatusOpen,
	}

	if err := u.offersRepo.Create(ctx, offer); err != nil {
		return nil, fmt.Errorf("failed to create offer: %w", err)
	}

	return offer, nil
}

func (u *usecase) GetOffer(ctx context.Context, id int64) (*entity.OfferWithDetails, error) {
	offer, err := u.offersRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get offer: %w", err)
	}

	if offer == nil {
		return nil, ErrOfferNotFound
	}

	return offer, nil
}

func (u *usecase) ListOffers(ctx context.Context, status string, limit, offset int) ([]entity.OfferWithDetails, error) {
	if limit <= 0 {
		limit = 50
	}

	offers, err := u.offersRepo.List(ctx, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list offers: %w", err)
	}

	return offers, nil
}

func (u *usecase) ListMyOffers(ctx context.Context, userID int64, limit, offset int) ([]entity.OfferWithDetails, error) {
	if limit <= 0 {
		limit = 50
	}

	offers, err := u.offersRepo.ListByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list user offers: %w", err)
	}

	return offers, nil
}

func (u *usecase) AcceptOffer(ctx context.Context, offerID, accepterID int64) error {
	offer, err := u.offersRepo.GetByID(ctx, offerID)
	if err != nil {
		return fmt.Errorf("failed to get offer: %w", err)
	}

	if offer == nil {
		return ErrOfferNotFound
	}

	if offer.Status != entity.OfferStatusOpen {
		return ErrOfferNotOpen
	}

	// Verify accepter has the requested skin
	hasSkin, err := u.skinsRepo.HasSkin(ctx, accepterID, offer.SkinRequestedID)
	if err != nil {
		return fmt.Errorf("failed to verify skin: %w", err)
	}

	if !hasSkin {
		return ErrInsufficientSkins
	}

	// Execute trade: swap skins
	if err := u.skinsRepo.RemoveSkinFromUser(ctx, offer.OwnerID, offer.SkinOfferedID, 1); err != nil {
		return fmt.Errorf("failed to remove skin from owner: %w", err)
	}

	if err := u.skinsRepo.RemoveSkinFromUser(ctx, accepterID, offer.SkinRequestedID, 1); err != nil {
		return fmt.Errorf("failed to remove skin from accepter: %w", err)
	}

	if err := u.skinsRepo.AddSkinToUser(ctx, accepterID, offer.SkinOfferedID, 1); err != nil {
		return fmt.Errorf("failed to add skin to accepter: %w", err)
	}

	if err := u.skinsRepo.AddSkinToUser(ctx, offer.OwnerID, offer.SkinRequestedID, 1); err != nil {
		return fmt.Errorf("failed to add skin to owner: %w", err)
	}

	// Update offer status
	if err := u.offersRepo.UpdateStatus(ctx, offerID, entity.OfferStatusAccepted); err != nil {
		return fmt.Errorf("failed to update offer: %w", err)
	}

	// Record trade
	trade := &entity.Trade{
		OfferID:    offerID,
		FromUserID: offer.OwnerID,
		ToUserID:   accepterID,
	}

	if err := u.tradesRepo.Create(ctx, trade); err != nil {
		return fmt.Errorf("failed to record trade: %w", err)
	}

	return nil
}

func (u *usecase) CancelOffer(ctx context.Context, offerID, userID int64) error {
	offer, err := u.offersRepo.GetByID(ctx, offerID)
	if err != nil {
		return fmt.Errorf("failed to get offer: %w", err)
	}

	if offer == nil {
		return ErrOfferNotFound
	}

	if offer.OwnerID != userID {
		return ErrNotAuthorized
	}

	if offer.Status != entity.OfferStatusOpen {
		return ErrOfferNotOpen
	}

	if err := u.offersRepo.UpdateStatus(ctx, offerID, entity.OfferStatusCancelled); err != nil {
		return fmt.Errorf("failed to cancel offer: %w", err)
	}

	return nil
}
