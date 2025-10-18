package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/youngermaster/bookstore/services/books-service/internal/domain"
	"github.com/youngermaster/bookstore/services/books-service/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrBookNotFound      = errors.New("book not found")
	ErrBookAlreadyExists = errors.New("book with this ISBN already exists")
	ErrInvalidInput      = errors.New("invalid input")
)

// BookService defines the interface for book business logic
type BookService interface {
	CreateBook(ctx context.Context, book *domain.Book) error
	GetBook(ctx context.Context, id uuid.UUID) (*domain.Book, error)
	GetBookByISBN(ctx context.Context, isbn string) (*domain.Book, error)
	ListBooks(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]domain.Book, int64, error)
	UpdateBook(ctx context.Context, book *domain.Book) error
	DeleteBook(ctx context.Context, id uuid.UUID) error
	UpdateBookStock(ctx context.Context, id uuid.UUID, quantity int) error
}

type bookService struct {
	bookRepo repository.BookRepository
}

// NewBookService creates a new instance of BookService
func NewBookService(bookRepo repository.BookRepository) BookService {
	return &bookService{
		bookRepo: bookRepo,
	}
}

func (s *bookService) CreateBook(ctx context.Context, book *domain.Book) error {
	if book == nil {
		return ErrInvalidInput
	}

	// Validate required fields
	if book.ISBN == "" || book.Title == "" || book.Price < 0 {
		return ErrInvalidInput
	}

	// Check if book with same ISBN already exists
	existing, err := s.bookRepo.FindByISBN(ctx, book.ISBN)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check existing book: %w", err)
	}
	if existing != nil {
		return ErrBookAlreadyExists
	}

	// Create the book
	if err := s.bookRepo.Create(ctx, book); err != nil {
		return fmt.Errorf("failed to create book: %w", err)
	}

	return nil
}

func (s *bookService) GetBook(ctx context.Context, id uuid.UUID) (*domain.Book, error) {
	book, err := s.bookRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBookNotFound
		}
		return nil, fmt.Errorf("failed to get book: %w", err)
	}
	return book, nil
}

func (s *bookService) GetBookByISBN(ctx context.Context, isbn string) (*domain.Book, error) {
	book, err := s.bookRepo.FindByISBN(ctx, isbn)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBookNotFound
		}
		return nil, fmt.Errorf("failed to get book by ISBN: %w", err)
	}
	return book, nil
}

func (s *bookService) ListBooks(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]domain.Book, int64, error) {
	// Set default pagination
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	books, total, err := s.bookRepo.FindAll(ctx, limit, offset, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list books: %w", err)
	}

	return books, total, nil
}

func (s *bookService) UpdateBook(ctx context.Context, book *domain.Book) error {
	if book == nil || book.ID == uuid.Nil {
		return ErrInvalidInput
	}

	// Check if book exists
	existing, err := s.bookRepo.FindByID(ctx, book.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBookNotFound
		}
		return fmt.Errorf("failed to check existing book: %w", err)
	}

	// If ISBN is being changed, check it's not already used
	if existing.ISBN != book.ISBN {
		existingISBN, err := s.bookRepo.FindByISBN(ctx, book.ISBN)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to check ISBN: %w", err)
		}
		if existingISBN != nil {
			return ErrBookAlreadyExists
		}
	}

	if err := s.bookRepo.Update(ctx, book); err != nil {
		return fmt.Errorf("failed to update book: %w", err)
	}

	return nil
}

func (s *bookService) DeleteBook(ctx context.Context, id uuid.UUID) error {
	// Check if book exists
	_, err := s.bookRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBookNotFound
		}
		return fmt.Errorf("failed to check existing book: %w", err)
	}

	if err := s.bookRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}

	return nil
}

func (s *bookService) UpdateBookStock(ctx context.Context, id uuid.UUID, quantity int) error {
	// Check if book exists
	book, err := s.bookRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBookNotFound
		}
		return fmt.Errorf("failed to check existing book: %w", err)
	}

	// Validate stock won't go negative
	if book.StockQuantity+quantity < 0 {
		return errors.New("insufficient stock")
	}

	if err := s.bookRepo.UpdateStock(ctx, id, quantity); err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	return nil
}
