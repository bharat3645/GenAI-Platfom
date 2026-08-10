package models

import (
	"time"
)

// DocumentEmbedding represents a chunk of text with its vector embedding
type DocumentEmbedding struct {
	ID              int                    `json:"id" db:"id"`
	DocumentID      int                    `json:"document_id" db:"document_id"`
	ChunkText       string                 `json:"chunk_text" db:"chunk_text"`
	ChunkIndex      int                    `json:"chunk_index" db:"chunk_index"`
	EmbeddingVector []float64              `json:"embedding_vector" db:"embedding_vector"`
	Metadata        map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt       time.Time              `json:"created_at" db:"created_at"`
}

// Entity represents an extracted named entity from documents
type Entity struct {
	ID              int                    `json:"id" db:"id"`
	DocumentID      int                    `json:"document_id" db:"document_id"`
	EntityType      string                 `json:"entity_type" db:"entity_type"`
	EntityName      string                 `json:"entity_name" db:"entity_name"`
	NormalizedName  string                 `json:"normalized_name" db:"normalized_name"`
	ConfidenceScore float64                `json:"confidence_score" db:"confidence_score"`
	OccurrenceCount int                    `json:"occurrence_count" db:"occurrence_count"`
	Positions       []EntityPosition       `json:"positions" db:"positions"`
	Metadata        map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt       time.Time              `json:"created_at" db:"created_at"`
}

// EntityPosition tracks where an entity appears in the document
type EntityPosition struct {
	ChunkID  int    `json:"chunk_id"`
	StartPos int    `json:"start_pos"`
	EndPos   int    `json:"end_pos"`
	Context  string `json:"context"`
}

// EntityRelationship represents a relationship between two entities
type EntityRelationship struct {
	ID               int                    `json:"id" db:"id"`
	SourceEntityID   int                    `json:"source_entity_id" db:"source_entity_id"`
	TargetEntityID   int                    `json:"target_entity_id" db:"target_entity_id"`
	RelationshipType string                 `json:"relationship_type" db:"relationship_type"`
	ConfidenceScore  float64                `json:"confidence_score" db:"confidence_score"`
	DocumentID       int                    `json:"document_id" db:"document_id"`
	Context          string                 `json:"context" db:"context"`
	Metadata         map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt        time.Time              `json:"created_at" db:"created_at"`
}

// GraphNode represents an aggregated entity across multiple documents
type GraphNode struct {
	ID                  int                    `json:"id" db:"id"`
	NodeType            string                 `json:"node_type" db:"node_type"`
	CanonicalName       string                 `json:"canonical_name" db:"canonical_name"`
	EntityIDs           []int                  `json:"entity_ids" db:"entity_ids"`
	DocumentIDs         []int                  `json:"document_ids" db:"document_ids"`
	CentralityScore     float64                `json:"centrality_score" db:"centrality_score"`
	PageRank            float64                `json:"page_rank" db:"page_rank"`
	CommunityID         *int                   `json:"community_id" db:"community_id"`
	OccurrenceFrequency int                    `json:"occurrence_frequency" db:"occurrence_frequency"`
	Metadata            map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt           time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at" db:"updated_at"`
}

// GraphEdge represents an aggregated relationship in the knowledge graph
type GraphEdge struct {
	ID              int                    `json:"id" db:"id"`
	SourceNodeID    int                    `json:"source_node_id" db:"source_node_id"`
	TargetNodeID    int                    `json:"target_node_id" db:"target_node_id"`
	EdgeType        string                 `json:"edge_type" db:"edge_type"`
	Weight          float64                `json:"weight" db:"weight"`
	RelationshipIDs []int                  `json:"relationship_ids" db:"relationship_ids"`
	Metadata        map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt       time.Time              `json:"created_at" db:"created_at"`
}
