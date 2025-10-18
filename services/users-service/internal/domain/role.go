package domain

import (
	"github.com/google/uuid"
)

// Role represents a user role
type Role struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"size:50;uniqueIndex;not null"` // e.g., "customer", "admin"
	Permissions string    `json:"permissions" gorm:"type:jsonb"`             // JSON array of permissions
}

// TableName specifies the table name for Role
func (Role) TableName() string {
	return "roles"
}

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	UserID uuid.UUID `gorm:"type:uuid;primaryKey"`
	RoleID uuid.UUID `gorm:"type:uuid;primaryKey"`
}

// TableName specifies the table name for UserRole
func (UserRole) TableName() string {
	return "user_roles"
}
