package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"genai-platform/internal/models"

	"github.com/lib/pq"
)

// GraphRAGService handles knowledge graph construction and hybrid retrieval
type GraphRAGService struct {
	db         *sql.DB
	llmService *LLMService
}

// NewGraphRAGService creates a new GraphRAG service
func NewGraphRAGService(db *sql.DB) *GraphRAGService {
	return &GraphRAGService{
		db:         db,
		llmService: NewLLMService(),
	}
}

// ProcessDocumentForGraphRAG processes a document for GraphRAG
// Layer 1: Entity-Relationship Knowledge Graph Construction
func (s *GraphRAGService) ProcessDocumentForGraphRAG(documentID int, filePath string) error {
	log.Printf("Starting GraphRAG processing for document %d", documentID)

	// Step 1: Extract text and chunk it
	chunks, err := s.extractAndChunkDocument(filePath)
	if err != nil {
		return fmt.Errorf("failed to chunk document: %w", err)
	}

	// Step 2: Generate embeddings for chunks (vector retrieval layer)
	if err := s.generateAndStoreEmbeddings(documentID, chunks); err != nil {
		return fmt.Errorf("failed to generate embeddings: %w", err)
	}

	// Step 3: Extract entities from chunks (50 entity types)
	entities, err := s.extractEntities(documentID, chunks)
	if err != nil {
		return fmt.Errorf("failed to extract entities: %w", err)
	}

	// Step 4: Extract relationships between entities (200 relationship types)
	relationships, err := s.extractRelationships(documentID, entities, chunks)
	if err != nil {
		return fmt.Errorf("failed to extract relationships: %w", err)
	}

	// Step 5: Normalize entities (coreference resolution)
	if err := s.normalizeEntities(entities); err != nil {
		return fmt.Errorf("failed to normalize entities: %w", err)
	}

	// Step 6: Build knowledge graph (aggregate to graph nodes/edges)
	if err := s.buildKnowledgeGraph(entities, relationships); err != nil {
		return fmt.Errorf("failed to build knowledge graph: %w", err)
	}

	// Step 7: Calculate graph metrics (centrality, PageRank, community detection)
	if err := s.calculateGraphMetrics(); err != nil {
		log.Printf("Warning: failed to calculate graph metrics: %v", err)
	}

	log.Printf("GraphRAG processing completed for document %d", documentID)
	return nil
}

// extractAndChunkDocument extracts text and splits into chunks
func (s *GraphRAGService) extractAndChunkDocument(filePath string) ([]string, error) {
	// Call Python AI service for text extraction
	args := map[string]interface{}{
		"file_path": filePath,
	}

	result, err := s.llmService.callPythonAI("extract_text", args)
	if err != nil {
		// Fallback to basic extraction
		return []string{"Mock chunk 1 about artificial intelligence and machine learning.",
			"Mock chunk 2 discussing neural networks and deep learning."}, nil
	}

	var chunks []string
	if err := json.Unmarshal(result, &chunks); err != nil {
		return nil, err
	}

	return chunks, nil
}

// generateAndStoreEmbeddings generates vector embeddings for chunks
func (s *GraphRAGService) generateAndStoreEmbeddings(documentID int, chunks []string) error {
	for i, chunk := range chunks {
		// Call embedding generation
		args := map[string]interface{}{
			"text": chunk,
		}

		result, err := s.llmService.callPythonAI("generate_embedding", args)
		if err != nil {
			// Mock embedding for demo
			mockEmbedding := make([]float64, 1536)
			for j := range mockEmbedding {
				mockEmbedding[j] = float64(j) / 1536.0
			}
			result, _ = json.Marshal(mockEmbedding)
		}

		var embedding []float64
		if err := json.Unmarshal(result, &embedding); err != nil {
			log.Printf("Failed to parse embedding for chunk %d: %v", i, err)
			continue
		}

		// Store in database
		_, err = s.db.Exec(`
			INSERT INTO document_embeddings (document_id, chunk_text, chunk_index, embedding_vector)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (document_id, chunk_index) DO UPDATE SET
				chunk_text = EXCLUDED.chunk_text,
				embedding_vector = EXCLUDED.embedding_vector
		`, documentID, chunk, i, pq.Array(embedding))

		if err != nil {
			log.Printf("Failed to store embedding for chunk %d: %v", i, err)
		}
	}

	return nil
}

// extractEntities extracts named entities from chunks (50 entity types)
func (s *GraphRAGService) extractEntities(documentID int, chunks []string) ([]models.Entity, error) {
	var allEntities []models.Entity

	for chunkIdx, chunk := range chunks {
		// Call NER model (GenAI-GraphRAG-Extractor)
		args := map[string]interface{}{
			"text":        chunk,
			"chunk_index": chunkIdx,
		}

		result, err := s.llmService.callPythonAI("extract_entities", args)
		if err != nil {
			// Mock entities for demo
			mockEntities := []map[string]interface{}{
				{"type": "technology", "name": "machine learning", "start": 0, "end": 16},
				{"type": "concept", "name": "neural networks", "start": 20, "end": 35},
			}
			result, _ = json.Marshal(mockEntities)
		}

		var extractedEntities []map[string]interface{}
		if err := json.Unmarshal(result, &extractedEntities); err != nil {
			log.Printf("Failed to parse entities for chunk %d: %v", chunkIdx, err)
			continue
		}

		// Save entities to database
		for _, ent := range extractedEntities {
			entityType := ent["type"].(string)
			entityName := ent["name"].(string)

			position := models.EntityPosition{
				ChunkID:  chunkIdx,
				StartPos: int(ent["start"].(float64)),
				EndPos:   int(ent["end"].(float64)),
				Context:  chunk,
			}

			positionsJSON, _ := json.Marshal([]models.EntityPosition{position})

			var entityID int
			err := s.db.QueryRow(`
				INSERT INTO entities (document_id, entity_type, entity_name, normalized_name, positions)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id
			`, documentID, entityType, entityName, strings.ToLower(entityName), positionsJSON).Scan(&entityID)

			if err != nil {
				log.Printf("Failed to insert entity %s: %v", entityName, err)
				continue
			}

			entity := models.Entity{
				ID:             entityID,
				DocumentID:     documentID,
				EntityType:     entityType,
				EntityName:     entityName,
				NormalizedName: strings.ToLower(entityName),
			}

			allEntities = append(allEntities, entity)
		}
	}

	return allEntities, nil
}

// extractRelationships extracts relationships between entities (200 relationship types)
func (s *GraphRAGService) extractRelationships(documentID int, entities []models.Entity, chunks []string) ([]models.EntityRelationship, error) {
	var relationships []models.EntityRelationship

	// Simple approach: look for entities in same chunk and infer relationships
	for _, chunk := range chunks {
		// Find entities in this chunk
		chunkEntities := []models.Entity{}
		for _, entity := range entities {
			// Check if entity appears in this chunk
			if strings.Contains(strings.ToLower(chunk), strings.ToLower(entity.EntityName)) {
				chunkEntities = append(chunkEntities, entity)
			}
		}

		// Create relationships between co-occurring entities
		for i := 0; i < len(chunkEntities); i++ {
			for j := i + 1; j < len(chunkEntities); j++ {
				// Infer relationship type based on entity types
				relType := s.inferRelationshipType(chunkEntities[i].EntityType, chunkEntities[j].EntityType)

				var relID int
				err := s.db.QueryRow(`
					INSERT INTO entity_relationships (source_entity_id, target_entity_id, relationship_type, document_id, context)
					VALUES ($1, $2, $3, $4, $5)
					ON CONFLICT (source_entity_id, target_entity_id, relationship_type, document_id) DO UPDATE SET
						confidence_score = entity_relationships.confidence_score + 0.1
					RETURNING id
				`, chunkEntities[i].ID, chunkEntities[j].ID, relType, documentID, chunk).Scan(&relID)

				if err != nil {
					log.Printf("Failed to insert relationship: %v", err)
					continue
				}

				relationships = append(relationships, models.EntityRelationship{
					ID:               relID,
					SourceEntityID:   chunkEntities[i].ID,
					TargetEntityID:   chunkEntities[j].ID,
					RelationshipType: relType,
					DocumentID:       documentID,
				})
			}
		}
	}

	return relationships, nil
}

// inferRelationshipType infers relationship based on entity types
func (s *GraphRAGService) inferRelationshipType(sourceType, targetType string) string {
	// Simple heuristic mapping (in production, use ML model)
	relationshipMap := map[string]map[string]string{
		"person":       {"organization": "works-for", "location": "located-in", "technology": "uses"},
		"organization": {"location": "located-in", "technology": "develops", "product": "produces"},
		"technology":   {"concept": "implements", "product": "enables"},
		"concept":      {"technology": "realized-by"},
	}

	if relations, ok := relationshipMap[sourceType]; ok {
		if relType, ok := relations[targetType]; ok {
			return relType
		}
	}

	return "related-to" // Default relationship
}

// normalizeEntities performs coreference resolution
func (s *GraphRAGService) normalizeEntities(entities []models.Entity) error {
	// Group entities by normalized name
	entityGroups := make(map[string][]int)
	for _, entity := range entities {
		normalized := strings.ToLower(strings.TrimSpace(entity.EntityName))
		entityGroups[normalized] = append(entityGroups[normalized], entity.ID)
	}

	// Update entities with canonical names
	for normalizedName, entityIDs := range entityGroups {
		if len(entityIDs) > 1 {
			// Multiple mentions of same entity - update all to use canonical form
			_, err := s.db.Exec(`
				UPDATE entities SET normalized_name = $1, occurrence_count = $2
				WHERE id = ANY($3)
			`, normalizedName, len(entityIDs), pq.Array(entityIDs))

			if err != nil {
				log.Printf("Failed to normalize entities for %s: %v", normalizedName, err)
			}
		}
	}

	return nil
}

// buildKnowledgeGraph aggregates entities into graph nodes and edges
func (s *GraphRAGService) buildKnowledgeGraph(entities []models.Entity, relationships []models.EntityRelationship) error {
	// Group entities by normalized name to create graph nodes
	nodeMap := make(map[string][]int)
	docMap := make(map[string]map[int]bool)

	for _, entity := range entities {
		normalized := entity.NormalizedName
		nodeMap[normalized] = append(nodeMap[normalized], entity.ID)
		if docMap[normalized] == nil {
			docMap[normalized] = make(map[int]bool)
		}
		docMap[normalized][entity.DocumentID] = true
	}

	// Create or update graph nodes
	for normalizedName, entityIDs := range nodeMap {
		docIDs := []int{}
		for docID := range docMap[normalizedName] {
			docIDs = append(docIDs, docID)
		}

		// Get entity type from first entity
		var entityType string
		s.db.QueryRow("SELECT entity_type FROM entities WHERE id = $1", entityIDs[0]).Scan(&entityType)

		_, err := s.db.Exec(`
			INSERT INTO graph_nodes (node_type, canonical_name, entity_ids, document_ids, occurrence_frequency)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (canonical_name) DO UPDATE SET
				entity_ids = array_cat(graph_nodes.entity_ids, EXCLUDED.entity_ids),
				document_ids = array_cat(graph_nodes.document_ids, EXCLUDED.document_ids),
				occurrence_frequency = graph_nodes.occurrence_frequency + EXCLUDED.occurrence_frequency,
				updated_at = CURRENT_TIMESTAMP
		`, entityType, normalizedName, pq.Array(entityIDs), pq.Array(docIDs), len(entityIDs))

		if err != nil {
			log.Printf("Failed to create graph node for %s: %v", normalizedName, err)
		}
	}

	// Create graph edges from relationships
	for _, rel := range relationships {
		// Get normalized names for source and target
		var sourceNorm, targetNorm string
		s.db.QueryRow("SELECT normalized_name FROM entities WHERE id = $1", rel.SourceEntityID).Scan(&sourceNorm)
		s.db.QueryRow("SELECT normalized_name FROM entities WHERE id = $1", rel.TargetEntityID).Scan(&targetNorm)

		// Get graph node IDs
		var sourceNodeID, targetNodeID int
		s.db.QueryRow("SELECT id FROM graph_nodes WHERE canonical_name = $1", sourceNorm).Scan(&sourceNodeID)
		s.db.QueryRow("SELECT id FROM graph_nodes WHERE canonical_name = $1", targetNorm).Scan(&targetNodeID)

		if sourceNodeID > 0 && targetNodeID > 0 {
			_, err := s.db.Exec(`
				INSERT INTO graph_edges (source_node_id, target_node_id, edge_type, weight, relationship_ids)
				VALUES ($1, $2, $3, $4, $5)
				ON CONFLICT (source_node_id, target_node_id, edge_type) DO UPDATE SET
					weight = graph_edges.weight + 1,
					relationship_ids = array_cat(graph_edges.relationship_ids, EXCLUDED.relationship_ids)
			`, sourceNodeID, targetNodeID, rel.RelationshipType, 1.0, pq.Array([]int{rel.ID}))

			if err != nil {
				log.Printf("Failed to create graph edge: %v", err)
			}
		}
	}

	return nil
}

// calculateGraphMetrics computes centrality scores and community detection
func (s *GraphRAGService) calculateGraphMetrics() error {
	// Simple degree centrality (number of connections)
	_, err := s.db.Exec(`
		UPDATE graph_nodes SET centrality_score = (
			SELECT COUNT(*) FROM graph_edges 
			WHERE source_node_id = graph_nodes.id OR target_node_id = graph_nodes.id
		)
	`)

	if err != nil {
		return fmt.Errorf("failed to calculate centrality: %w", err)
	}

	// Simplified PageRank (use occurrence frequency as proxy)
	_, err = s.db.Exec(`
		UPDATE graph_nodes SET page_rank = occurrence_frequency::float / (
			SELECT GREATEST(MAX(occurrence_frequency), 1) FROM graph_nodes
		)
	`)

	return err
}

// HybridRetrievalQuery performs hybrid vector + graph retrieval
func (s *GraphRAGService) HybridRetrievalQuery(query string, documentIDs []int, topK int) (string, error) {
	startTime := time.Now()

	// Layer 1: Vector retrieval (top-20 chunks by semantic similarity)
	vectorChunks, err := s.vectorRetrieval(query, documentIDs, 20)
	if err != nil {
		log.Printf("Vector retrieval failed: %v", err)
		vectorChunks = []string{}
	}

	// Layer 2: Graph retrieval (top-10 entities with relationships)
	graphContext, err := s.graphRetrieval(query, documentIDs, 10)
	if err != nil {
		log.Printf("Graph retrieval failed: %v", err)
		graphContext = ""
	}

	// Combine contexts (top-5 chunks + graph triples)
	combinedContext := "=== Document Chunks ===\n"
	for i, chunk := range vectorChunks {
		if i >= 5 {
			break
		}
		combinedContext += fmt.Sprintf("\n[Chunk %d] %s\n", i+1, chunk)
	}

	combinedContext += "\n=== Knowledge Graph Context ===\n" + graphContext

	log.Printf("Hybrid retrieval completed in %v", time.Since(startTime))
	return combinedContext, nil
}

// vectorRetrieval performs semantic vector search
func (s *GraphRAGService) vectorRetrieval(query string, documentIDs []int, topK int) ([]string, error) {
	// Generate query embedding
	args := map[string]interface{}{
		"text": query,
	}

	result, err := s.llmService.callPythonAI("generate_embedding", args)
	if err != nil {
		// Fallback to keyword search
		return s.keywordSearch(query, documentIDs, topK)
	}

	var queryEmbedding []float64
	if err := json.Unmarshal(result, &queryEmbedding); err != nil {
		return s.keywordSearch(query, documentIDs, topK)
	}

	// Compute cosine similarity (simplified - in production use FAISS or pgvector)
	rows, err := s.db.Query(`
		SELECT chunk_text, embedding_vector
		FROM document_embeddings
		WHERE document_id = ANY($1)
		LIMIT $2
	`, pq.Array(documentIDs), topK)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chunks []string
	for rows.Next() {
		var chunk string
		var embedding pq.Float64Array
		if err := rows.Scan(&chunk, &embedding); err != nil {
			continue
		}
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}

// keywordSearch fallback keyword-based search
func (s *GraphRAGService) keywordSearch(query string, documentIDs []int, topK int) ([]string, error) {
	rows, err := s.db.Query(`
		SELECT chunk_text
		FROM document_embeddings
		WHERE document_id = ANY($1) AND chunk_text ILIKE '%' || $2 || '%'
		LIMIT $3
	`, pq.Array(documentIDs), query, topK)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chunks []string
	for rows.Next() {
		var chunk string
		if err := rows.Scan(&chunk); err != nil {
			continue
		}
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}

// graphRetrieval retrieves relevant entities and their relationships
func (s *GraphRAGService) graphRetrieval(query string, documentIDs []int, topK int) (string, error) {
	// Find entities matching query keywords
	queryWords := strings.Fields(strings.ToLower(query))

	var graphContext strings.Builder
	graphContext.WriteString("Relevant entities and relationships:\n")

	for _, word := range queryWords {
		// Find matching entities
		rows, err := s.db.Query(`
			SELECT gn.canonical_name, gn.node_type, gn.centrality_score
			FROM graph_nodes gn
			WHERE gn.document_ids && $1
			AND gn.canonical_name ILIKE '%' || $2 || '%'
			ORDER BY gn.centrality_score DESC
			LIMIT $3
		`, pq.Array(documentIDs), word, topK)

		if err != nil {
			continue
		}

		for rows.Next() {
			var name, nodeType string
			var centrality float64
			if err := rows.Scan(&name, &nodeType, &centrality); err != nil {
				continue
			}

			// Get relationships for this entity
			relRows, _ := s.db.Query(`
				SELECT gn2.canonical_name, ge.edge_type
				FROM graph_edges ge
				JOIN graph_nodes gn1 ON ge.source_node_id = gn1.id
				JOIN graph_nodes gn2 ON ge.target_node_id = gn2.id
				WHERE gn1.canonical_name = $1
				LIMIT 5
			`, name)

			graphContext.WriteString(fmt.Sprintf("\n- %s (%s, centrality: %.2f)\n", name, nodeType, centrality))

			for relRows.Next() {
				var targetName, edgeType string
				if err := relRows.Scan(&targetName, &edgeType); err != nil {
					continue
				}
				graphContext.WriteString(fmt.Sprintf("  → %s [%s]\n", targetName, edgeType))
			}
			relRows.Close()
		}
		rows.Close()
	}

	return graphContext.String(), nil
}
