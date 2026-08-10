-- Enhanced schema for GraphRAG, Multi-Agent systems, and Analytics
-- Migration 002: Add advanced features

-- ============================================================================
-- GRAPHRAG TABLES
-- ============================================================================

-- Document embeddings for vector search
CREATE TABLE IF NOT EXISTS document_embeddings (
    id SERIAL PRIMARY KEY,
    document_id INTEGER REFERENCES documents(id) ON DELETE CASCADE,
    chunk_text TEXT NOT NULL,
    chunk_index INTEGER NOT NULL,
    embedding_vector FLOAT8[] NOT NULL,  -- Store as array, migrate to pgvector later
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(document_id, chunk_index)
);

CREATE INDEX IF NOT EXISTS idx_document_embeddings_doc_id ON document_embeddings(document_id);
CREATE INDEX IF NOT EXISTS idx_document_embeddings_metadata ON document_embeddings USING GIN(metadata);

-- Entities extracted from documents (50 entity types)
CREATE TABLE IF NOT EXISTS entities (
    id SERIAL PRIMARY KEY,
    document_id INTEGER REFERENCES documents(id) ON DELETE CASCADE,
    entity_type VARCHAR(100) NOT NULL,  -- person, organization, technology, concept, etc.
    entity_name VARCHAR(500) NOT NULL,
    normalized_name VARCHAR(500) NOT NULL,  -- Canonical form after coreference resolution
    confidence_score FLOAT DEFAULT 1.0,
    occurrence_count INTEGER DEFAULT 1,
    positions JSONB DEFAULT '[]',  -- [{chunk_id, start_pos, end_pos, context}]
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_entities_document ON entities(document_id);
CREATE INDEX IF NOT EXISTS idx_entities_type ON entities(entity_type);
CREATE INDEX IF NOT EXISTS idx_entities_normalized ON entities(normalized_name);
CREATE INDEX IF NOT EXISTS idx_entities_metadata ON entities USING GIN(metadata);

-- Relationships between entities (200 relationship types)
CREATE TABLE IF NOT EXISTS entity_relationships (
    id SERIAL PRIMARY KEY,
    source_entity_id INTEGER REFERENCES entities(id) ON DELETE CASCADE,
    target_entity_id INTEGER REFERENCES entities(id) ON DELETE CASCADE,
    relationship_type VARCHAR(100) NOT NULL,  -- works-for, located-in, uses-technology, etc.
    confidence_score FLOAT DEFAULT 1.0,
    document_id INTEGER REFERENCES documents(id) ON DELETE CASCADE,
    context TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(source_entity_id, target_entity_id, relationship_type, document_id)
);

CREATE INDEX IF NOT EXISTS idx_relationships_source ON entity_relationships(source_entity_id);
CREATE INDEX IF NOT EXISTS idx_relationships_target ON entity_relationships(target_entity_id);
CREATE INDEX IF NOT EXISTS idx_relationships_type ON entity_relationships(relationship_type);
CREATE INDEX IF NOT EXISTS idx_relationships_document ON entity_relationships(document_id);

-- Knowledge graph nodes (aggregated entities across documents)
CREATE TABLE IF NOT EXISTS graph_nodes (
    id SERIAL PRIMARY KEY,
    node_type VARCHAR(100) NOT NULL,
    canonical_name VARCHAR(500) NOT NULL UNIQUE,
    entity_ids INTEGER[] DEFAULT '{}',  -- References to entities table
    document_ids INTEGER[] DEFAULT '{}',
    centrality_score FLOAT DEFAULT 0.0,
    page_rank FLOAT DEFAULT 0.0,
    community_id INTEGER,
    occurrence_frequency INTEGER DEFAULT 0,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_graph_nodes_type ON graph_nodes(node_type);
CREATE INDEX IF NOT EXISTS idx_graph_nodes_community ON graph_nodes(community_id);
CREATE INDEX IF NOT EXISTS idx_graph_nodes_centrality ON graph_nodes(centrality_score DESC);

-- Graph edges (aggregated relationships)
CREATE TABLE IF NOT EXISTS graph_edges (
    id SERIAL PRIMARY KEY,
    source_node_id INTEGER REFERENCES graph_nodes(id) ON DELETE CASCADE,
    target_node_id INTEGER REFERENCES graph_nodes(id) ON DELETE CASCADE,
    edge_type VARCHAR(100) NOT NULL,
    weight FLOAT DEFAULT 1.0,
    relationship_ids INTEGER[] DEFAULT '{}',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(source_node_id, target_node_id, edge_type)
);

CREATE INDEX IF NOT EXISTS idx_graph_edges_source ON graph_edges(source_node_id);
CREATE INDEX IF NOT EXISTS idx_graph_edges_target ON graph_edges(target_node_id);
CREATE INDEX IF NOT EXISTS idx_graph_edges_type ON graph_edges(edge_type);

-- ============================================================================
-- MULTI-AGENT SYSTEM TABLES
-- ============================================================================

-- Agent execution states
CREATE TABLE IF NOT EXISTS agent_executions (
    id SERIAL PRIMARY KEY,
    execution_type VARCHAR(50) NOT NULL,  -- 'resume_analysis', 'research', 'graphrag_query'
    user_id INTEGER REFERENCES users(id),
    parent_task_id INTEGER,  -- For tracking related tasks
    agent_name VARCHAR(100) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',  -- pending, running, completed, failed
    input_data JSONB NOT NULL,
    output_data JSONB DEFAULT '{}',
    error_message TEXT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    duration_ms INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_agent_executions_type ON agent_executions(execution_type);
CREATE INDEX IF NOT EXISTS idx_agent_executions_user ON agent_executions(user_id);
CREATE INDEX IF NOT EXISTS idx_agent_executions_parent ON agent_executions(parent_task_id);
CREATE INDEX IF NOT EXISTS idx_agent_executions_status ON agent_executions(status);

-- Agent collaboration memory (shared working memory)
CREATE TABLE IF NOT EXISTS agent_memory (
    id SERIAL PRIMARY KEY,
    execution_id INTEGER REFERENCES agent_executions(id) ON DELETE CASCADE,
    memory_key VARCHAR(255) NOT NULL,
    memory_value JSONB NOT NULL,
    agent_name VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_agent_memory_execution ON agent_memory(execution_id);
CREATE INDEX IF NOT EXISTS idx_agent_memory_key ON agent_memory(memory_key);

-- Enhanced resume analyses with agent-specific outputs
CREATE TABLE IF NOT EXISTS resume_analysis_details (
    id SERIAL PRIMARY KEY,
    analysis_id INTEGER REFERENCES resume_analyses(id) ON DELETE CASCADE,
    agent_name VARCHAR(100) NOT NULL,
    agent_output JSONB NOT NULL,
    execution_time_ms INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_resume_details_analysis ON resume_analysis_details(analysis_id);
CREATE INDEX IF NOT EXISTS idx_resume_details_agent ON resume_analysis_details(agent_name);

-- Keywords extracted and analyzed (for ATS keyword agent)
CREATE TABLE IF NOT EXISTS resume_keywords (
    id SERIAL PRIMARY KEY,
    analysis_id INTEGER REFERENCES resume_analyses(id) ON DELETE CASCADE,
    keyword VARCHAR(255) NOT NULL,
    keyword_type VARCHAR(50),  -- 'skill', 'technology', 'soft_skill', 'domain', 'certification'
    found_in_resume BOOLEAN DEFAULT FALSE,
    found_in_jd BOOLEAN DEFAULT FALSE,
    relevance_score FLOAT DEFAULT 0.0,
    density_score FLOAT DEFAULT 0.0,
    placement_quality VARCHAR(50),  -- 'excellent', 'good', 'poor', 'missing'
    synonyms TEXT[] DEFAULT '{}',
    recommendations TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_resume_keywords_analysis ON resume_keywords(analysis_id);
CREATE INDEX IF NOT EXISTS idx_resume_keywords_keyword ON resume_keywords(keyword);

-- ============================================================================
-- RESEARCH ASSISTANT TABLES
-- ============================================================================

-- Research subtasks (HTN decomposition)
CREATE TABLE IF NOT EXISTS research_subtasks (
    id SERIAL PRIMARY KEY,
    task_id INTEGER REFERENCES research_tasks(id) ON DELETE CASCADE,
    parent_subtask_id INTEGER REFERENCES research_subtasks(id),
    subtask_query TEXT NOT NULL,
    subtask_type VARCHAR(50),  -- 'literature_review', 'survey', 'analysis', 'synthesis'
    priority INTEGER DEFAULT 5,
    status VARCHAR(50) DEFAULT 'pending',
    dependencies INTEGER[] DEFAULT '{}',
    result TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_research_subtasks_task ON research_subtasks(task_id);
CREATE INDEX IF NOT EXISTS idx_research_subtasks_parent ON research_subtasks(parent_subtask_id);
CREATE INDEX IF NOT EXISTS idx_research_subtasks_status ON research_subtasks(status);

-- Retrieved sources (multi-source intelligence)
CREATE TABLE IF NOT EXISTS research_sources (
    id SERIAL PRIMARY KEY,
    task_id INTEGER REFERENCES research_tasks(id) ON DELETE CASCADE,
    subtask_id INTEGER REFERENCES research_subtasks(id),
    source_type VARCHAR(50),  -- 'web', 'arxiv', 'pubmed', 'scholar'
    source_url TEXT,
    title TEXT,
    authors TEXT[] DEFAULT '{}',
    publication_date DATE,
    relevance_score FLOAT DEFAULT 0.0,
    credibility_score FLOAT DEFAULT 0.0,
    content_summary TEXT,
    full_content TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_research_sources_task ON research_sources(task_id);
CREATE INDEX IF NOT EXISTS idx_research_sources_subtask ON research_sources(subtask_id);
CREATE INDEX IF NOT EXISTS idx_research_sources_relevance ON research_sources(relevance_score DESC);

-- Fact verification results
CREATE TABLE IF NOT EXISTS fact_verifications (
    id SERIAL PRIMARY KEY,
    task_id INTEGER REFERENCES research_tasks(id) ON DELETE CASCADE,
    source_id INTEGER REFERENCES research_sources(id),
    claim_text TEXT NOT NULL,
    verification_status VARCHAR(50),  -- 'supported', 'contradicted', 'neutral', 'insufficient'
    supporting_sources INTEGER[] DEFAULT '{}',
    contradicting_sources INTEGER[] DEFAULT '{}',
    confidence_score FLOAT DEFAULT 0.0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_fact_verifications_task ON fact_verifications(task_id);
CREATE INDEX IF NOT EXISTS idx_fact_verifications_source ON fact_verifications(source_id);

-- Citations management
CREATE TABLE IF NOT EXISTS research_citations (
    id SERIAL PRIMARY KEY,
    task_id INTEGER REFERENCES research_tasks(id) ON DELETE CASCADE,
    source_id INTEGER REFERENCES research_sources(id),
    citation_style VARCHAR(50) DEFAULT 'IEEE',  -- IEEE, APA, MLA
    citation_text TEXT NOT NULL,
    page_reference VARCHAR(50),
    quotation_text TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_research_citations_task ON research_citations(task_id);
CREATE INDEX IF NOT EXISTS idx_research_citations_source ON research_citations(source_id);

-- ============================================================================
-- TEXT-TO-SQL ENHANCEMENTS
-- ============================================================================

-- SQL schema metadata (for schema-aware generation)
CREATE TABLE IF NOT EXISTS sql_schema_metadata (
    id SERIAL PRIMARY KEY,
    table_name VARCHAR(255) NOT NULL,
    column_name VARCHAR(255) NOT NULL,
    column_type VARCHAR(100),
    is_primary_key BOOLEAN DEFAULT FALSE,
    is_foreign_key BOOLEAN DEFAULT FALSE,
    references_table VARCHAR(255),
    references_column VARCHAR(255),
    sample_values TEXT[] DEFAULT '{}',
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_sql_schema_table ON sql_schema_metadata(table_name);

-- SQL validation results (triple-layer safety)
CREATE TABLE IF NOT EXISTS sql_validation_results (
    id SERIAL PRIMARY KEY,
    query_id INTEGER REFERENCES sql_queries(id) ON DELETE CASCADE,
    validation_layer VARCHAR(50) NOT NULL,  -- 'parsing', 'allowlist', 'manual'
    validation_status VARCHAR(50),  -- 'passed', 'failed', 'warning'
    validation_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_sql_validation_query ON sql_validation_results(query_id);

-- ============================================================================
-- ANALYTICS AND MONITORING
-- ============================================================================

-- System metrics
CREATE TABLE IF NOT EXISTS system_metrics (
    id SERIAL PRIMARY KEY,
    metric_name VARCHAR(255) NOT NULL,
    metric_value FLOAT NOT NULL,
    metric_type VARCHAR(50),  -- 'latency', 'throughput', 'accuracy', 'uptime'
    service_name VARCHAR(100),
    metadata JSONB DEFAULT '{}',
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_system_metrics_name ON system_metrics(metric_name);
CREATE INDEX IF NOT EXISTS idx_system_metrics_service ON system_metrics(service_name);
CREATE INDEX IF NOT EXISTS idx_system_metrics_recorded ON system_metrics(recorded_at DESC);

-- User activity tracking
CREATE TABLE IF NOT EXISTS user_activities (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    activity_type VARCHAR(100) NOT NULL,  -- 'pdf_upload', 'chat_query', 'resume_analysis', etc.
    activity_data JSONB DEFAULT '{}',
    response_time_ms INTEGER,
    success BOOLEAN DEFAULT TRUE,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_activities_user ON user_activities(user_id);
CREATE INDEX IF NOT EXISTS idx_user_activities_type ON user_activities(activity_type);
CREATE INDEX IF NOT EXISTS idx_user_activities_created ON user_activities(created_at DESC);

-- Model inference logs
CREATE TABLE IF NOT EXISTS model_inference_logs (
    id SERIAL PRIMARY KEY,
    model_name VARCHAR(255) NOT NULL,
    model_version VARCHAR(50),
    input_tokens INTEGER,
    output_tokens INTEGER,
    inference_time_ms INTEGER,
    success BOOLEAN DEFAULT TRUE,
    error_message TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_model_logs_model ON model_inference_logs(model_name);
CREATE INDEX IF NOT EXISTS idx_model_logs_created ON model_inference_logs(created_at DESC);

-- ============================================================================
-- PERFORMANCE OPTIMIZATIONS
-- ============================================================================

-- Function to update graph node metrics (centrality, pagerank)
CREATE OR REPLACE FUNCTION update_graph_metrics()
RETURNS TRIGGER AS $$
BEGIN
    -- Update occurrence frequency when new entities are linked
    UPDATE graph_nodes
    SET occurrence_frequency = array_length(entity_ids, 1),
        updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_graph_metrics
AFTER UPDATE OF entity_ids ON graph_nodes
FOR EACH ROW
EXECUTE FUNCTION update_graph_metrics();

-- Function to log user activities automatically
CREATE OR REPLACE FUNCTION log_user_activity()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO user_activities (user_id, activity_type, activity_data)
    VALUES (
        NEW.user_id,
        TG_TABLE_NAME,
        jsonb_build_object('id', NEW.id, 'created_at', NEW.created_at)
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply activity logging to key tables (optional, can be enabled per table)
-- Example: CREATE TRIGGER log_document_upload AFTER INSERT ON documents FOR EACH ROW EXECUTE FUNCTION log_user_activity();
