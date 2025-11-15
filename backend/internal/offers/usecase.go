package offers

import (
	"context"

	"github.com/cs2-p2p-skins/backend/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, offer *entity.Offer) error
	GetByID(ctx context.Context, id int64) (*entity.OfferWithDetails, error)
	List(ctx context.Context, status string, limit, offset int) ([]entity.OfferWithDetails, error)
	ListByUser(ctx context.Context, userID int64, limit, offset int) ([]entity.OfferWithDetails, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
}

type UseCase interface {
	CreateOffer(ctx context.Context, ownerID, skinOfferedID, skinRequestedID int64) (*entity.Offer, error)
	GetOffer(ctx context.Context, id int64) (*entity.OfferWithDetails, error)
	ListOffers(ctx context.Context, status string, limit, offset int) ([]entity.OfferWithDetails, error)
	ListMyOffers(ctx context.Context, userID int64, limit, offset int) ([]entity.OfferWithDetails, error)
	AcceptOffer(ctx context.Context, offerID, accepterID int64) error
	CancelOffer(ctx context.Context, offerID, userID int64) error
}
