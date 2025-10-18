package domain

import (
	"time"

	"github.com/google/uuid"
)

// Log represents a structured log entry
type Log struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ServiceName string    `json:"service_name" gorm:"size:100;not null;index"`
	Level       string    `json:"level" gorm:"size:20;not null;index"` // DEBUG, INFO, WARN, ERROR
	Message     string    `json:"message" gorm:"type:text;not null"`
	TraceID     string    `json:"trace_id" gorm:"size:100;index"`
	SpanID      string    `json:"span_id" gorm:"size:100"`
	UserID      *uuid.UUID `json:"user_id" gorm:"type:uuid;index"`
	Metadata    string    `json:"metadata" gorm:"type:jsonb"` // Additional context as JSON
	Timestamp   time.Time `json:"timestamp" gorm:"not null;index"`
	CreatedAt   time.Time `json:"created_at"`
}

// TableName specifies the table name for Log
func (Log) TableName() string {
	return "logs"
}
