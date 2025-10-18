package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"not null"` // - means don't include in JSON
	FullName     string    `json:"full_name" gorm:"size:255;not null"`
	Roles        []Role    `json:"roles,omitempty" gorm:"many2many:user_roles;"`
	Addresses    []Address `json:"addresses,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}
