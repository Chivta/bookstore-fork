package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/youngermaster/bookstore/services/logging-service/internal/domain"
	"gorm.io/gorm"
)

// LogService defines the interface for logging business logic
type LogService interface {
	CreateLog(ctx context.Context, log *domain.Log) error
	GetLogs(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]domain.Log, int64, error)
	DeleteOldLogs(ctx context.Context, olderThan time.Duration) error
}

type logService struct {
	db *gorm.DB
}

// NewLogService creates a new instance of LogService
func NewLogService(db *gorm.DB) LogService {
	return &logService{db: db}
}

func (s *logService) CreateLog(ctx context.Context, log *domain.Log) error {
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now().UTC()
	}
	return s.db.WithContext(ctx).Create(log).Error
}

func (s *logService) GetLogs(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]domain.Log, int64, error) {
	var logs []domain.Log
	var total int64

	query := s.db.WithContext(ctx).Model(&domain.Log{})

	// Apply filters
	if serviceName, ok := filters["service_name"]; ok {
		query = query.Where("service_name = ?", serviceName)
	}
	if level, ok := filters["level"]; ok {
		query = query.Where("level = ?", level)
	}
	if traceID, ok := filters["trace_id"]; ok {
		query = query.Where("trace_id = ?", traceID)
	}
	if userID, ok := filters["user_id"]; ok {
		if uid, ok := userID.(uuid.UUID); ok {
			query = query.Where("user_id = ?", uid)
		}
	}
	if startTime, ok := filters["start_time"]; ok {
		if t, ok := startTime.(time.Time); ok {
			query = query.Where("timestamp >= ?", t)
		}
	}
	if endTime, ok := filters["end_time"]; ok {
		if t, ok := endTime.(time.Time); ok {
			query = query.Where("timestamp <= ?", t)
		}
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	err := query.
		Order("timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error

	return logs, total, err
}

func (s *logService) DeleteOldLogs(ctx context.Context, olderThan time.Duration) error {
	cutoffTime := time.Now().UTC().Add(-olderThan)
	result := s.db.WithContext(ctx).
		Where("timestamp < ?", cutoffTime).
		Delete(&domain.Log{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete old logs: %w", result.Error)
	}

	return nil
}
