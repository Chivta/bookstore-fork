package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/youngermaster/bookstore/services/books-service/internal/domain"
	"github.com/youngermaster/bookstore/services/books-service/internal/service"
)

// BookHandler handles HTTP requests for books
type BookHandler struct {
	bookService service.BookService
}

// NewBookHandler creates a new instance of BookHandler
func NewBookHandler(bookService service.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

// CreateBook handles POST /api/v1/books
func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
	var book domain.Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.bookService.CreateBook(c.Context(), &book); err != nil {
		if errors.Is(err, service.ErrBookAlreadyExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if errors.Is(err, service.ErrInvalidInput) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create book",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(book)
}

// GetBook handles GET /api/v1/books/:id
func (h *BookHandler) GetBook(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid book ID",
		})
	}

	book, err := h.bookService.GetBook(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrBookNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Book not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get book",
		})
	}

	return c.JSON(book)
}

// ListBooks handles GET /api/v1/books
func (h *BookHandler) ListBooks(c *fiber.Ctx) error {
	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	// Parse filters
	filters := make(map[string]interface{})
	if categoryID := c.Query("category_id"); categoryID != "" {
		if id, err := uuid.Parse(categoryID); err == nil {
			filters["category_id"] = id
		}
	}
	if authorID := c.Query("author_id"); authorID != "" {
		if id, err := uuid.Parse(authorID); err == nil {
			filters["author_id"] = id
		}
	}
	if title := c.Query("title"); title != "" {
		filters["title"] = title
	}
	if minPrice := c.Query("min_price"); minPrice != "" {
		if price, err := strconv.ParseFloat(minPrice, 64); err == nil {
			filters["min_price"] = price
		}
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if price, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filters["max_price"] = price
		}
	}

	books, total, err := h.bookService.ListBooks(c.Context(), limit, offset, filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list books",
		})
	}

	return c.JSON(fiber.Map{
		"data":   books,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// UpdateBook handles PUT /api/v1/books/:id
func (h *BookHandler) UpdateBook(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid book ID",
		})
	}

	var book domain.Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	book.ID = id

	if err := h.bookService.UpdateBook(c.Context(), &book); err != nil {
		if errors.Is(err, service.ErrBookNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Book not found",
			})
		}
		if errors.Is(err, service.ErrBookAlreadyExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update book",
		})
	}

	return c.JSON(book)
}

// DeleteBook handles DELETE /api/v1/books/:id
func (h *BookHandler) DeleteBook(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid book ID",
		})
	}

	if err := h.bookService.DeleteBook(c.Context(), id); err != nil {
		if errors.Is(err, service.ErrBookNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Book not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete book",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// UpdateStock handles PATCH /api/v1/books/:id/stock
func (h *BookHandler) UpdateStock(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid book ID",
		})
	}

	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.bookService.UpdateBookStock(c.Context(), id, req.Quantity); err != nil {
		if errors.Is(err, service.ErrBookNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Book not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update stock",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Stock updated successfully",
	})
}
