package domain

import (
	"time"

	"github.com/google/uuid"
)

// Address represents a user's shipping/billing address
type Address struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Street     string    `json:"street" gorm:"size:255;not null"`
	City       string    `json:"city" gorm:"size:100;not null"`
	State      string    `json:"state" gorm:"size:100"`
	PostalCode string    `json:"postal_code" gorm:"size:20;not null"`
	Country    string    `json:"country" gorm:"size:100;not null"`
	IsDefault  bool      `json:"is_default" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TableName specifies the table name for Address
func (Address) TableName() string {
	return "addresses"
}
