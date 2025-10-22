package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/youngermaster/my-distributed-bookstore/services/users-service/internal/domain"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	AssignRole(ctx context.Context, userID, roleID uuid.UUID) error
}

// AddressRepository defines the interface for address data access
type AddressRepository interface {
	Create(ctx context.Context, address *domain.Address) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Address, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Address, error)
	Update(ctx context.Context, address *domain.Address) error
	Delete(ctx context.Context, id uuid.UUID) error
	SetDefault(ctx context.Context, userID, addressID uuid.UUID) error
}

// SessionRepository defines the interface for session data access
type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*domain.Session, error)
	DeleteExpired(ctx context.Context) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
