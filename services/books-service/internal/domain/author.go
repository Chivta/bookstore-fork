package domain

import (
	"time"

	"github.com/google/uuid"
)

// Author represents a book author
type Author struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string     `json:"name" gorm:"size:255;not null"`
	Bio       string     `json:"bio" gorm:"type:text"`
	BirthDate *time.Time `json:"birth_date"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// TableName specifies the table name for Author
func (Author) TableName() string {
	return "authors"
}
