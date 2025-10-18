package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/youngermaster/my-distributed-bookstore/services/users-service/internal/domain"
	"github.com/youngermaster/my-distributed-bookstore/services/users-service/internal/repository"
	"gorm.io/gorm"
)

type wishlistRepositoryImpl struct {
	db *gorm.DB
}

// NewWishlistRepository creates a new instance of WishlistRepository
func NewWishlistRepository(db *gorm.DB) repository.WishlistRepository {
	return &wishlistRepositoryImpl{db: db}
}

// GetByUserID retrieves all wishlist items for a user
func (r *wishlistRepositoryImpl) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.WishlistItem, error) {
	var items []domain.WishlistItem
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// Add adds a book to user's wishlist
func (r *wishlistRepositoryImpl) Add(ctx context.Context, item *domain.WishlistItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

// Remove removes a book from user's wishlist
func (r *wishlistRepositoryImpl) Remove(ctx context.Context, userID, bookID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND book_id = ?", userID, bookID).
		Delete(&domain.WishlistItem{}).Error
}

// Exists checks if a book is in user's wishlist
func (r *wishlistRepositoryImpl) Exists(ctx context.Context, userID, bookID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.WishlistItem{}).
		Where("user_id = ? AND book_id = ?", userID, bookID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
