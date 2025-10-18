package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/youngermaster/bookstore/services/logging-service/internal/domain"
	"github.com/youngermaster/bookstore/services/logging-service/internal/service"
)

// LogHandler handles HTTP requests for logs
type LogHandler struct {
	logService service.LogService
}

// NewLogHandler creates a new instance of LogHandler
func NewLogHandler(logService service.LogService) *LogHandler {
	return &LogHandler{
		logService: logService,
	}
}

// CreateLog handles POST /api/v1/logs
func (h *LogHandler) CreateLog(c *fiber.Ctx) error {
	var log domain.Log
	if err := c.BodyParser(&log); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.logService.CreateLog(c.Context(), &log); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create log",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(log)
}

// GetLogs handles GET /api/v1/logs
func (h *LogHandler) GetLogs(c *fiber.Ctx) error {
	// Parse pagination
	limit, _ := strconv.Atoi(c.Query("limit", "100"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	// Parse filters
	filters := make(map[string]interface{})
	if serviceName := c.Query("service_name"); serviceName != "" {
		filters["service_name"] = serviceName
	}
	if level := c.Query("level"); level != "" {
		filters["level"] = level
	}
	if traceID := c.Query("trace_id"); traceID != "" {
		filters["trace_id"] = traceID
	}
	if userID := c.Query("user_id"); userID != "" {
		if uid, err := uuid.Parse(userID); err == nil {
			filters["user_id"] = uid
		}
	}
	if startTime := c.Query("start_time"); startTime != "" {
		if t, err := time.Parse(time.RFC3339, startTime); err == nil {
			filters["start_time"] = t
		}
	}
	if endTime := c.Query("end_time"); endTime != "" {
		if t, err := time.Parse(time.RFC3339, endTime); err == nil {
			filters["end_time"] = t
		}
	}

	logs, total, err := h.logService.GetLogs(c.Context(), filters, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get logs",
		})
	}

	return c.JSON(fiber.Map{
		"data":   logs,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}
