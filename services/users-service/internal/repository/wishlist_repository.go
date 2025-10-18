package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/youngermaster/my-distributed-bookstore/services/users-service/internal/domain"
)

// WishlistRepository defines methods for wishlist data access
type WishlistRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.WishlistItem, error)
	Add(ctx context.Context, item *domain.WishlistItem) error
	Remove(ctx context.Context, userID, bookID uuid.UUID) error
	Exists(ctx context.Context, userID, bookID uuid.UUID) (bool, error)
}
