package models

import (
	"time"
)

type User struct {
	ID           int       `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Document struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Filename  string    `json:"filename" db:"filename"`
	FilePath  string    `json:"file_path" db:"file_path"`
	FileType  string    `json:"file_type" db:"file_type"`
	FileSize  int       `json:"file_size" db:"file_size"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ChatSession struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	DocumentIDs []int     `json:"document_ids" db:"document_ids"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type ChatMessage struct {
	ID        int                    `json:"id" db:"id"`
	SessionID int                    `json:"session_id" db:"session_id"`
	Role      string                 `json:"role" db:"role"`
	Content   string                 `json:"content" db:"content"`
	Metadata  map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
}

type ResearchTask struct {
	ID          int                    `json:"id" db:"id"`
	UserID      int                    `json:"user_id" db:"user_id"`
	Query       string                 `json:"query" db:"query"`
	Status      string                 `json:"status" db:"status"`
	Result      string                 `json:"result" db:"result"`
	Metadata    map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	CompletedAt *time.Time             `json:"completed_at" db:"completed_at"`
}

type ResumeAnalysis struct {
	ID             int        `json:"id" db:"id"`
	UserID         int        `json:"user_id" db:"user_id"`
	ResumePath     string     `json:"resume_path" db:"resume_path"`
	JobDescription string     `json:"job_description" db:"job_description"`
	Feedback       string     `json:"feedback" db:"feedback"`
	Score          int        `json:"score" db:"score"`
	Status         string     `json:"status" db:"status"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	CompletedAt    *time.Time `json:"completed_at" db:"completed_at"`
}

type SQLQuery struct {
	ID           int                    `json:"id" db:"id"`
	UserID       int                    `json:"user_id" db:"user_id"`
	NaturalQuery string                 `json:"natural_query" db:"natural_query"`
	GeneratedSQL string                 `json:"generated_sql" db:"generated_sql"`
	ResultData   map[string]interface{} `json:"result_data" db:"result_data"`
	Status       string                 `json:"status" db:"status"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
}

