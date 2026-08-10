package models

import (
	"time"
)

// SystemMetric tracks system performance metrics
type SystemMetric struct {
	ID          int                    `json:"id" db:"id"`
	MetricName  string                 `json:"metric_name" db:"metric_name"`
	MetricValue float64                `json:"metric_value" db:"metric_value"`
	MetricType  *string                `json:"metric_type" db:"metric_type"`
	ServiceName *string                `json:"service_name" db:"service_name"`
	Metadata    map[string]interface{} `json:"metadata" db:"metadata"`
	RecordedAt  time.Time              `json:"recorded_at" db:"recorded_at"`
}

// UserActivity tracks user interactions
type UserActivity struct {
	ID             int                    `json:"id" db:"id"`
	UserID         *int                   `json:"user_id" db:"user_id"`
	ActivityType   string                 `json:"activity_type" db:"activity_type"`
	ActivityData   map[string]interface{} `json:"activity_data" db:"activity_data"`
	ResponseTimeMs *int                   `json:"response_time_ms" db:"response_time_ms"`
	Success        bool                   `json:"success" db:"success"`
	ErrorMessage   *string                `json:"error_message" db:"error_message"`
	CreatedAt      time.Time              `json:"created_at" db:"created_at"`
}

// ModelInferenceLog tracks AI model usage
type ModelInferenceLog struct {
	ID              int                    `json:"id" db:"id"`
	ModelName       string                 `json:"model_name" db:"model_name"`
	ModelVersion    *string                `json:"model_version" db:"model_version"`
	InputTokens     *int                   `json:"input_tokens" db:"input_tokens"`
	OutputTokens    *int                   `json:"output_tokens" db:"output_tokens"`
	InferenceTimeMs *int                   `json:"inference_time_ms" db:"inference_time_ms"`
	Success         bool                   `json:"success" db:"success"`
	ErrorMessage    *string                `json:"error_message" db:"error_message"`
	Metadata        map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt       time.Time              `json:"created_at" db:"created_at"`
}

// SQLSchemaMetadata stores database schema information
type SQLSchemaMetadata struct {
	ID               int       `json:"id" db:"id"`
	TableName        string    `json:"table_name" db:"table_name"`
	ColumnName       string    `json:"column_name" db:"column_name"`
	ColumnType       *string   `json:"column_type" db:"column_type"`
	IsPrimaryKey     bool      `json:"is_primary_key" db:"is_primary_key"`
	IsForeignKey     bool      `json:"is_foreign_key" db:"is_foreign_key"`
	ReferencesTable  *string   `json:"references_table" db:"references_table"`
	ReferencesColumn *string   `json:"references_column" db:"references_column"`
	SampleValues     []string  `json:"sample_values" db:"sample_values"`
	Description      *string   `json:"description" db:"description"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// SQLValidationResult stores SQL validation checks
type SQLValidationResult struct {
	ID                int       `json:"id" db:"id"`
	QueryID           int       `json:"query_id" db:"query_id"`
	ValidationLayer   string    `json:"validation_layer" db:"validation_layer"`
	ValidationStatus  *string   `json:"validation_status" db:"validation_status"`
	ValidationMessage *string   `json:"validation_message" db:"validation_message"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}
