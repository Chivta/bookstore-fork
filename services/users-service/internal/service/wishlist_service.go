package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/youngermaster/my-distributed-bookstore/services/users-service/internal/domain"
	"github.com/youngermaster/my-distributed-bookstore/services/users-service/internal/repository"
)

// WishlistService handles wishlist business logic
type WishlistService struct {
	wishlistRepo repository.WishlistRepository
}

// NewWishlistService creates a new WishlistService
func NewWishlistService(wishlistRepo repository.WishlistRepository) *WishlistService {
	return &WishlistService{
		wishlistRepo: wishlistRepo,
	}
}

// GetUserWishlist retrieves all wishlist items for a user
func (s *WishlistService) GetUserWishlist(ctx context.Context, userID uuid.UUID) ([]domain.WishlistItem, error) {
	return s.wishlistRepo.GetByUserID(ctx, userID)
}

// AddToWishlist adds a book to user's wishlist
func (s *WishlistService) AddToWishlist(ctx context.Context, userID, bookID uuid.UUID) (*domain.WishlistItem, error) {
	// Check if already in wishlist
	exists, err := s.wishlistRepo.Exists(ctx, userID, bookID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("book already in wishlist")
	}

	item := &domain.WishlistItem{
		UserID: userID,
		BookID: bookID,
	}

	if err := s.wishlistRepo.Add(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

// RemoveFromWishlist removes a book from user's wishlist
func (s *WishlistService) RemoveFromWishlist(ctx context.Context, userID, bookID uuid.UUID) error {
	return s.wishlistRepo.Remove(ctx, userID, bookID)
}
