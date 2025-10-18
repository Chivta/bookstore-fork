package domain

import (
	"time"

	"github.com/google/uuid"
)

// WishlistItem represents a book in a user's wishlist
type WishlistItem struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	BookID    uuid.UUID `json:"book_id" gorm:"type:uuid;not null"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName specifies the table name for WishlistItem
func (WishlistItem) TableName() string {
	return "wishlist_items"
}
