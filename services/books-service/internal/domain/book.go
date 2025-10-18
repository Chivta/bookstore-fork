package domain

import (
	"time"

	"github.com/google/uuid"
)

// Book represents a book in the catalog
type Book struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ISBN            string     `json:"isbn" gorm:"uniqueIndex;not null"`
	Title           string     `json:"title" gorm:"size:500;not null"`
	Description     string     `json:"description" gorm:"type:text"`
	PublisherID     *uuid.UUID `json:"publisher_id" gorm:"type:uuid"`
	Publisher       *Publisher `json:"publisher,omitempty" gorm:"foreignKey:PublisherID"`
	PublicationDate *time.Time `json:"publication_date"`
	Language        string     `json:"language" gorm:"size:10;default:'en'"`
	Pages           int        `json:"pages"`
	Format          string     `json:"format" gorm:"size:50"` // hardcover, paperback, ebook
	Price           float64    `json:"price" gorm:"not null;check:price >= 0"`
	StockQuantity   int        `json:"stock_quantity" gorm:"not null;default:0;check:stock_quantity >= 0"`
	CoverImageURL   string     `json:"cover_image_url" gorm:"type:text"`
	Metadata        string     `json:"metadata" gorm:"type:jsonb"` // flexible additional data
	Authors         []Author   `json:"authors,omitempty" gorm:"many2many:book_authors;"`
	Categories      []Category `json:"categories,omitempty" gorm:"many2many:book_categories;"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// TableName specifies the table name for Book
func (Book) TableName() string {
	return "books"
}

// BookAuthor represents the many-to-many relationship between books and authors
type BookAuthor struct {
	BookID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	AuthorID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	AuthorOrder int       `gorm:"default:1"`
}

// TableName specifies the table name for BookAuthor
func (BookAuthor) TableName() string {
	return "book_authors"
}

// BookCategory represents the many-to-many relationship between books and categories
type BookCategory struct {
	BookID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	CategoryID uuid.UUID `gorm:"type:uuid;primaryKey"`
}

// TableName specifies the table name for BookCategory
func (BookCategory) TableName() string {
	return "book_categories"
}
