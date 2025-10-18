package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/youngermaster/bookstore/services/books-service/internal/domain"
	"github.com/youngermaster/bookstore/services/books-service/internal/repository"
	"gorm.io/gorm"
)

type bookRepository struct {
	db *gorm.DB
}

// NewBookRepository creates a new instance of BookRepository
func NewBookRepository(db *gorm.DB) repository.BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(ctx context.Context, book *domain.Book) error {
	return r.db.WithContext(ctx).Create(book).Error
}

func (r *bookRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Book, error) {
	var book domain.Book
	err := r.db.WithContext(ctx).
		Preload("Authors").
		Preload("Categories").
		Preload("Publisher").
		First(&book, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) FindByISBN(ctx context.Context, isbn string) (*domain.Book, error) {
	var book domain.Book
	err := r.db.WithContext(ctx).
		Preload("Authors").
		Preload("Categories").
		Preload("Publisher").
		First(&book, "isbn = ?", isbn).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]domain.Book, int64, error) {
	var books []domain.Book
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Book{})

	// Apply filters
	if categoryID, ok := filters["category_id"]; ok {
		query = query.Joins("JOIN book_categories ON book_categories.book_id = books.id").
			Where("book_categories.category_id = ?", categoryID)
	}

	if authorID, ok := filters["author_id"]; ok {
		query = query.Joins("JOIN book_authors ON book_authors.book_id = books.id").
			Where("book_authors.author_id = ?", authorID)
	}

	if title, ok := filters["title"]; ok {
		query = query.Where("title ILIKE ?", fmt.Sprintf("%%%s%%", title))
	}

	if minPrice, ok := filters["min_price"]; ok {
		query = query.Where("price >= ?", minPrice)
	}

	if maxPrice, ok := filters["max_price"]; ok {
		query = query.Where("price <= ?", maxPrice)
	}

	// Count total matching records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated results with preloaded relationships
	err := query.
		Preload("Authors").
		Preload("Categories").
		Preload("Publisher").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&books).Error

	return books, total, err
}

func (r *bookRepository) Update(ctx context.Context, book *domain.Book) error {
	return r.db.WithContext(ctx).Save(book).Error
}

func (r *bookRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Book{}, "id = ?", id).Error
}

func (r *bookRepository) UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error {
	return r.db.WithContext(ctx).Model(&domain.Book{}).
		Where("id = ?", id).
		Update("stock_quantity", gorm.Expr("stock_quantity + ?", quantity)).Error
}
