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
        'process_document': process_document,
        'search_similar_chunks': search_similar_chunks,
        'generate_chat_response': generate_chat_response,  
        'analyze_resume': analyze_resume,
        'generate_sql_from_natural_language': generate_sql_from_natural_language,
        'conduct_research': conduct_research,
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

