package skins

import (
	"context"

	"github.com/cs2-p2p-skins/backend/internal/entity"
)

type UseCase interface {
	ListSkins(ctx context.Context, limit, offset int) ([]entity.Skin, error)
	GetSkin(ctx context.Context, id int64) (*entity.Skin, error)
	GetUserSkins(ctx context.Context, userID int64) ([]entity.UserSkinWithDetails, error)
	AddSkinToUser(ctx context.Context, userID, skinID int64, quantity int) error
}
