#!/usr/bin/env python3
"""
AI Bridge - Python script to handle AI operations for the Go backend
"""
import sys
import json
import logging
from ai_service import ai_service

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

def process_document(args):
    """Process a document and add it to vector store."""
    document_id = args.get('document_id')
    file_path = args.get('file_path')
    
    success = ai_service.process_document(file_path, document_id)
    return {"success": success}

def search_similar_chunks(args):
    """Search for similar chunks."""
    query = args.get('query', '')
    document_ids = args.get('document_ids', [])
    
    chunks = ai_service.search_similar_chunks(query, document_ids)
    return chunks

def generate_chat_response(args):
    """Generate chat response."""
    query = args.get('query', '')
    context = args.get('context', [])
    
    try:
        response = ai_service.generate_chat_response(query, context)
        return {"response": response}
    except Exception as e:
        return {"response": "", "error": str(e)}

def analyze_resume(args):
    """Analyze resume."""
    resume_path = args.get('resume_path', '')
    job_description = args.get('job_description', '')
    
    result = ai_service.analyze_resume(resume_path, job_description)
    return result

def generate_sql_from_natural_language(args):
    """Generate SQL from natural language."""
    natural_query = args.get('natural_query', '')
    
    try:
        sql = ai_service.generate_sql_from_natural_language(natural_query)
        return {"response": sql}
    except Exception as e:
        return {"response": "", "error": str(e)}

def conduct_research(args):
    """Conduct research."""
    research_query = args.get('research_query', '')
    
    try:
        result = ai_service.conduct_research(research_query)
        return {"response": result}
    except Exception as e:
        return {"response": "", "error": str(e)}

# ==================== GraphRAG Methods ====================

def extract_entities(args):
    """Extract named entities from text for GraphRAG."""
    text = args.get('text', '')
    entity_types = args.get('entity_types', [])
    
    try:
        entities = ai_service.extract_entities(text, entity_types)
        return {"entities": entities}
    except Exception as e:
        return {"entities": [], "error": str(e)}

def extract_relationships(args):
    """Extract relationships between entities."""
    text = args.get('text', '')
    entities = args.get('entities', [])
    
    try:
        relationships = ai_service.extract_relationships(text, entities)
        return {"relationships": relationships}
    except Exception as e:
        return {"relationships": [], "error": str(e)}

def generate_embedding(args):
    """Generate embedding vector for text."""
    text = args.get('text', '')
    
    try:
        embedding = ai_service.generate_embedding(text)
        return {"embedding": embedding}
    except Exception as e:
        return {"embedding": [], "error": str(e)}

# ==================== Multi-Agent ATS Methods ====================

def intelligent_keyword_matching(args):
    """Intelligent keyword matching with semantic equivalence."""
    resume_text = args.get('resume_text', '')
    keywords = args.get('keywords', [])
    
    try:
        result = ai_service.intelligent_keyword_matching(resume_text, keywords)
        return result
    except Exception as e:
        return {"matched": [], "missing": keywords, "error": str(e)}

# ==================== Research Agent Methods ====================

def decompose_research_query(args):
    """Decompose research query into HTN subtasks."""
    query = args.get('query', '')
    
    try:
        subtasks = ai_service.decompose_research_query(query)
        return {"subtasks": subtasks}
    except Exception as e:
        return {"subtasks": [], "error": str(e)}

def summarize_research_source(args):
    """Summarize research source."""
    source_text = args.get('source_text', '')
    max_length = args.get('max_length', 300)
    
    try:
        summary = ai_service.summarize_research_source(source_text, max_length)
        return {"summary": summary}
    except Exception as e:
        return {"summary": "", "error": str(e)}

def verify_fact(args):
    """Verify a factual claim against sources."""
    claim = args.get('claim', '')
    sources = args.get('sources', [])
    
    try:
        result = ai_service.verify_fact(claim, sources)
        return result
    except Exception as e:
        return {"verdict": "error", "confidence": 0.0, "error": str(e)}

# ==================== Text-to-SQL Methods ====================

def generate_sql_with_schema(args):
    """Generate SQL with full schema awareness."""
    natural_query = args.get('natural_query', '')
    schema_context = args.get('schema_context', '')
    
    try:
        sql = ai_service.generate_sql_with_schema(natural_query, schema_context)
        return {"sql": sql}
    except Exception as e:
        return {"sql": "", "error": str(e)}

def explain_sql_query(args):
    """Explain SQL query in natural language."""
    sql = args.get('sql', '')
    
    try:
        explanation = ai_service.explain_sql_query(sql)
        return {"explanation": explanation}
    except Exception as e:
        return {"explanation": "", "error": str(e)}

def main():
    if len(sys.argv) != 3:
        print(json.dumps({"error": "Usage: ai_bridge.py <method> <args_json>"}))
        sys.exit(1)
    
    method = sys.argv[1]
    args_json = sys.argv[2]
    
    try:
        args = json.loads(args_json)
    except json.JSONDecodeError as e:
        print(json.dumps({"error": f"Invalid JSON: {e}"}))
        sys.exit(1)
    
    # Route to appropriate function
    functions = {
        # Original methods
        'process_document': process_document,
        'search_similar_chunks': search_similar_chunks,
        'generate_chat_response': generate_chat_response,  
        'analyze_resume': analyze_resume,
        'generate_sql_from_natural_language': generate_sql_from_natural_language,
        'conduct_research': conduct_research,
        # GraphRAG methods
        'extract_entities': extract_entities,
        'extract_relationships': extract_relationships,
        'generate_embedding': generate_embedding,
        # Multi-Agent ATS methods
        'intelligent_keyword_matching': intelligent_keyword_matching,
        # Research Agent methods
        'decompose_research_query': decompose_research_query,
        'summarize_research_source': summarize_research_source,
        'verify_fact': verify_fact,
        # Text-to-SQL methods
        'generate_sql_with_schema': generate_sql_with_schema,
        'explain_sql_query': explain_sql_query,
    }
    
    if method not in functions:
        print(json.dumps({"error": f"Unknown method: {method}"}))
        sys.exit(1)
    
    try:
        result = functions[method](args)
        print(json.dumps(result))
    except Exception as e:
        logger.error(f"Error in {method}: {e}")
        print(json.dumps({"error": str(e)}))
        sys.exit(1)

if __name__ == "__main__":
    main()

