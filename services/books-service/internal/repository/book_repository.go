package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/youngermaster/bookstore/services/books-service/internal/domain"
)

// BookRepository defines the interface for book data access
type BookRepository interface {
	Create(ctx context.Context, book *domain.Book) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Book, error)
	FindByISBN(ctx context.Context, isbn string) (*domain.Book, error)
	FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]domain.Book, int64, error)
	Update(ctx context.Context, book *domain.Book) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error
}

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	FindBySlug(ctx context.Context, slug string) (*domain.Category, error)
	FindAll(ctx context.Context) ([]domain.Category, error)
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// AuthorRepository defines the interface for author data access
type AuthorRepository interface {
	Create(ctx context.Context, author *domain.Author) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Author, error)
	FindAll(ctx context.Context, limit, offset int) ([]domain.Author, int64, error)
	Update(ctx context.Context, author *domain.Author) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// PublisherRepository defines the interface for publisher data access
type PublisherRepository interface {
	Create(ctx context.Context, publisher *domain.Publisher) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Publisher, error)
	FindAll(ctx context.Context, limit, offset int) ([]domain.Publisher, int64, error)
	Update(ctx context.Context, publisher *domain.Publisher) error
	Delete(ctx context.Context, id uuid.UUID) error
}
