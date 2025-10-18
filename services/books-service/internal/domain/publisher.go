package domain

import (
	"time"

	"github.com/google/uuid"
)

// Publisher represents a book publisher
type Publisher struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"size:255;not null"`
	Country   string    `json:"country" gorm:"size:100"`
	Website   string    `json:"website" gorm:"size:255"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for Publisher
func (Publisher) TableName() string {
	return "publishers"
}
