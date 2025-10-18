package domain

import (
	"time"

	"github.com/google/uuid"
)

// Category represents a book category (hierarchical)
type Category struct {
	ID        uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string      `json:"name" gorm:"size:100;not null"`
	Slug      string      `json:"slug" gorm:"size:100;uniqueIndex;not null"`
	ParentID  *uuid.UUID  `json:"parent_id" gorm:"type:uuid"`
	Parent    *Category   `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children  []Category  `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// TableName specifies the table name for Category
func (Category) TableName() string {
	return "categories"
}
