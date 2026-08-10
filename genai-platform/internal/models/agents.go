package models

import (
	"time"
)

// AgentExecution tracks the execution of a single agent
type AgentExecution struct {
	ID            int                    `json:"id" db:"id"`
	ExecutionType string                 `json:"execution_type" db:"execution_type"`
	UserID        *int                   `json:"user_id" db:"user_id"`
	ParentTaskID  *int                   `json:"parent_task_id" db:"parent_task_id"`
	AgentName     string                 `json:"agent_name" db:"agent_name"`
	Status        string                 `json:"status" db:"status"`
	InputData     map[string]interface{} `json:"input_data" db:"input_data"`
	OutputData    map[string]interface{} `json:"output_data" db:"output_data"`
	ErrorMessage  *string                `json:"error_message" db:"error_message"`
	StartedAt     *time.Time             `json:"started_at" db:"started_at"`
	CompletedAt   *time.Time             `json:"completed_at" db:"completed_at"`
	DurationMs    *int                   `json:"duration_ms" db:"duration_ms"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
}

// AgentMemory represents shared memory between agents
type AgentMemory struct {
	ID          int                    `json:"id" db:"id"`
	ExecutionID int                    `json:"execution_id" db:"execution_id"`
	MemoryKey   string                 `json:"memory_key" db:"memory_key"`
	MemoryValue map[string]interface{} `json:"memory_value" db:"memory_value"`
	AgentName   *string                `json:"agent_name" db:"agent_name"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

// ResumeAnalysisDetail stores per-agent analysis results
type ResumeAnalysisDetail struct {
	ID              int                    `json:"id" db:"id"`
	AnalysisID      int                    `json:"analysis_id" db:"analysis_id"`
	AgentName       string                 `json:"agent_name" db:"agent_name"`
	AgentOutput     map[string]interface{} `json:"agent_output" db:"agent_output"`
	ExecutionTimeMs *int                   `json:"execution_time_ms" db:"execution_time_ms"`
	CreatedAt       time.Time              `json:"created_at" db:"created_at"`
}

// ResumeKeyword represents a keyword analysis result
type ResumeKeyword struct {
	ID               int      `json:"id" db:"id"`
	AnalysisID       int      `json:"analysis_id" db:"analysis_id"`
	Keyword          string   `json:"keyword" db:"keyword"`
	KeywordType      *string  `json:"keyword_type" db:"keyword_type"`
	FoundInResume    bool     `json:"found_in_resume" db:"found_in_resume"`
	FoundInJD        bool     `json:"found_in_jd" db:"found_in_jd"`
	RelevanceScore   float64  `json:"relevance_score" db:"relevance_score"`
	DensityScore     float64  `json:"density_score" db:"density_score"`
	PlacementQuality *string  `json:"placement_quality" db:"placement_quality"`
	Synonyms         []string `json:"synonyms" db:"synonyms"`
	Recommendations  *string  `json:"recommendations" db:"recommendations"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}
