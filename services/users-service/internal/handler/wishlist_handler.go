package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/youngermaster/my-distributed-bookstore/services/users-service/internal/service"
)

// WishlistHandler handles wishlist-related HTTP requests
type WishlistHandler struct {
	wishlistService *service.WishlistService
}

// NewWishlistHandler creates a new WishlistHandler
func NewWishlistHandler(wishlistService *service.WishlistService) *WishlistHandler {
	return &WishlistHandler{
		wishlistService: wishlistService,
	}
}

// GetWishlist retrieves the user's wishlist
// @Summary Get user wishlist
// @Tags wishlist
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/me/wishlist [get]
func (h *WishlistHandler) GetWishlist(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	uid, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	items, err := h.wishlistService.GetUserWishlist(c.Context(), uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch wishlist",
		})
	}

	return c.JSON(fiber.Map{
		"data":  items,
		"total": len(items),
	})
}

// AddToWishlist adds a book to the user's wishlist
// @Summary Add book to wishlist
// @Tags wishlist
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "Book ID"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/me/wishlist [post]
func (h *WishlistHandler) AddToWishlist(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	uid, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req struct {
		BookID string `json:"book_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	bookID, err := uuid.Parse(req.BookID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid book ID",
		})
	}

	item, err := h.wishlistService.AddToWishlist(c.Context(), uid, bookID)
	if err != nil {
		if err.Error() == "book already in wishlist" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add to wishlist",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": item,
	})
}

// RemoveFromWishlist removes a book from the user's wishlist
// @Summary Remove book from wishlist
// @Tags wishlist
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param book_id path string true "Book ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/me/wishlist/{book_id} [delete]
func (h *WishlistHandler) RemoveFromWishlist(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	uid, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	bookID, err := uuid.Parse(c.Params("book_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid book ID",
		})
	}

	if err := h.wishlistService.RemoveFromWishlist(c.Context(), uid, bookID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove from wishlist",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Book removed from wishlist",
	})
}
