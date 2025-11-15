package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/cs2-p2p-skins/backend/internal/entity"
	"github.com/cs2-p2p-skins/backend/internal/skins"
)

var (
	ErrSkinNotFound = errors.New("skin not found")
)

type usecase struct {
	repo skins.Repository
}

func New(repo skins.Repository) *usecase {
	return &usecase{repo: repo}
}

func (u *usecase) ListSkins(ctx context.Context, limit, offset int) ([]entity.Skin, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	skins, err := u.repo.ListSkins(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list skins: %w", err)
	}

	return skins, nil
}

func (u *usecase) GetSkin(ctx context.Context, id int64) (*entity.Skin, error) {
	skin, err := u.repo.GetSkinByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get skin: %w", err)
	}

	if skin == nil {
		return nil, ErrSkinNotFound
	}

	return skin, nil
}

func (u *usecase) GetUserSkins(ctx context.Context, userID int64) ([]entity.UserSkinWithDetails, error) {
	skins, err := u.repo.GetUserSkins(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user skins: %w", err)
	}

	return skins, nil
}

func (u *usecase) AddSkinToUser(ctx context.Context, userID, skinID int64, quantity int) error {
	// Verify skin exists
	skin, err := u.repo.GetSkinByID(ctx, skinID)
	if err != nil {
		return fmt.Errorf("failed to verify skin: %w", err)
	}

	if skin == nil {
		return ErrSkinNotFound
	}

	if err := u.repo.AddSkinToUser(ctx, userID, skinID, quantity); err != nil {
		return fmt.Errorf("failed to add skin: %w", err)
	}

	return nil
}

func (u *usecase) RemoveSkinFromUser(ctx context.Context, userID, skinID int64, quantity int) error {
	// Verify skin exists
	skin, err := u.repo.GetSkinByID(ctx, skinID)
	if err != nil {
		return fmt.Errorf("failed to verify skin: %w", err)
	}

	if skin == nil {
		return ErrSkinNotFound
	}

	if err := u.repo.RemoveSkinFromUser(ctx, userID, skinID, quantity); err != nil {
		return fmt.Errorf("failed to remove skin: %w", err)
	}

	return nil
}
