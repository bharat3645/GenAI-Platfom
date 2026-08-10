package models

import (
	"time"
)

// ResearchSubtask represents a decomposed subtask
type ResearchSubtask struct {
	ID              int        `json:"id" db:"id"`
	TaskID          int        `json:"task_id" db:"task_id"`
	ParentSubtaskID *int       `json:"parent_subtask_id" db:"parent_subtask_id"`
	SubtaskQuery    string     `json:"subtask_query" db:"subtask_query"`
	SubtaskType     *string    `json:"subtask_type" db:"subtask_type"`
	Priority        int        `json:"priority" db:"priority"`
	Status          string     `json:"status" db:"status"`
	Dependencies    []int      `json:"dependencies" db:"dependencies"`
	Result          *string    `json:"result" db:"result"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	CompletedAt     *time.Time `json:"completed_at" db:"completed_at"`
}

// ResearchSource represents a retrieved information source
type ResearchSource struct {
	ID               int                    `json:"id" db:"id"`
	TaskID           int                    `json:"task_id" db:"task_id"`
	SubtaskID        *int                   `json:"subtask_id" db:"subtask_id"`
	SourceType       *string                `json:"source_type" db:"source_type"`
	SourceURL        *string                `json:"source_url" db:"source_url"`
	Title            *string                `json:"title" db:"title"`
	Authors          []string               `json:"authors" db:"authors"`
	PublicationDate  *time.Time             `json:"publication_date" db:"publication_date"`
	RelevanceScore   float64                `json:"relevance_score" db:"relevance_score"`
	CredibilityScore float64                `json:"credibility_score" db:"credibility_score"`
	ContentSummary   *string                `json:"content_summary" db:"content_summary"`
	FullContent      *string                `json:"full_content" db:"full_content"`
	Metadata         map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt        time.Time              `json:"created_at" db:"created_at"`
}

// FactVerification represents a fact-checking result
type FactVerification struct {
	ID                     int       `json:"id" db:"id"`
	TaskID                 int       `json:"task_id" db:"task_id"`
	SourceID               *int      `json:"source_id" db:"source_id"`
	ClaimText              string    `json:"claim_text" db:"claim_text"`
	VerificationStatus     *string   `json:"verification_status" db:"verification_status"`
	SupportingSources      []int     `json:"supporting_sources" db:"supporting_sources"`
	ContradictingSources   []int     `json:"contradicting_sources" db:"contradicting_sources"`
	ConfidenceScore        float64   `json:"confidence_score" db:"confidence_score"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
}

// ResearchCitation represents a citation
type ResearchCitation struct {
	ID             int        `json:"id" db:"id"`
	TaskID         int        `json:"task_id" db:"task_id"`
	SourceID       *int       `json:"source_id" db:"source_id"`
	CitationStyle  string     `json:"citation_style" db:"citation_style"`
	CitationText   string     `json:"citation_text" db:"citation_text"`
	PageReference  *string    `json:"page_reference" db:"page_reference"`
	QuotationText  *string    `json:"quotation_text" db:"quotation_text"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
}
