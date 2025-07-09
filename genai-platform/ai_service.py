import os
import sys
import json
import logging
from typing import List, Dict, Any
from pathlib import Path

# Add the project root to Python path
sys.path.append('/home/ubuntu/genai-platform')

from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain_openai import OpenAIEmbeddings, ChatOpenAI
from langchain_google_genai import GoogleGenerativeAIEmbeddings, ChatGoogleGenerativeAI
import faiss
import numpy as np
import PyPDF2
from docx import Document

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class AIService:
    def __init__(self):
        self.openai_api_key = os.getenv('OPENAI_API_KEY', '')
        self.gemini_api_key = os.getenv('GEMINI_API_KEY', '')
        
        # Initialize embeddings
        if self.openai_api_key:
            self.embeddings = OpenAIEmbeddings(openai_api_key=self.openai_api_key)
            self.llm = ChatOpenAI(openai_api_key=self.openai_api_key, model="gpt-3.5-turbo")
        elif self.gemini_api_key:
            self.embeddings = GoogleGenerativeAIEmbeddings(
                model="models/embedding-001",
                google_api_key=self.gemini_api_key
            )
            self.llm = ChatGoogleGenerativeAI(
                model="gemini-pro",
                google_api_key=self.gemini_api_key
            )
        else:
            logger.warning("No API keys found. Using mock implementations.")
            self.embeddings = None
            self.llm = None
        
        # Initialize text splitter
        self.text_splitter = RecursiveCharacterTextSplitter(
            chunk_size=1000,
            chunk_overlap=200,
            length_function=len,
        )
        
        # Initialize FAISS index
        self.vector_store = None
        self.document_chunks = []
        self.chunk_metadata = []

    def extract_text_from_pdf(self, file_path: str) -> str:
        """Extract text from PDF file."""
        try:
            with open(file_path, 'rb') as file:
                pdf_reader = PyPDF2.PdfReader(file)
                text = ""
                for page in pdf_reader.pages:
                    text += page.extract_text() + "\n"
                return text
        except Exception as e:
            logger.error(f"Error extracting text from PDF {file_path}: {e}")
            return ""

    def extract_text_from_docx(self, file_path: str) -> str:
        """Extract text from DOCX file."""
        try:
            doc = Document(file_path)
            text = ""
            for paragraph in doc.paragraphs:
                text += paragraph.text + "\n"
            return text
        except Exception as e:
            logger.error(f"Error extracting text from DOCX {file_path}: {e}")
            return ""

    def process_document(self, file_path: str, document_id: int) -> bool:
        """Process a document and add it to the vector store."""
        try:
            # Extract text based on file extension
            file_ext = Path(file_path).suffix.lower()
            if file_ext == '.pdf':
                text = self.extract_text_from_pdf(file_path)
            elif file_ext in ['.docx', '.doc']:
                text = self.extract_text_from_docx(file_path)
            else:
                logger.error(f"Unsupported file type: {file_ext}")
                return False

            if not text.strip():
                logger.error(f"No text extracted from {file_path}")
                return False

            # Split text into chunks
            chunks = self.text_splitter.split_text(text)
            
            if not self.embeddings:
                # Mock implementation
                logger.info(f"Mock processing: {len(chunks)} chunks for document {document_id}")
                return True

            # Generate embeddings
            embeddings = self.embeddings.embed_documents(chunks)
            
            # Initialize FAISS index if not exists
            if self.vector_store is None:
                dimension = len(embeddings[0])
                self.vector_store = faiss.IndexFlatL2(dimension)
            
            # Add to FAISS index
            embeddings_array = np.array(embeddings).astype('float32')
            self.vector_store.add(embeddings_array)
            
            # Store chunks and metadata
            for i, chunk in enumerate(chunks):
                self.document_chunks.append(chunk)
                self.chunk_metadata.append({
                    'document_id': document_id,
                    'chunk_index': i,
                    'file_path': file_path
                })
            
            logger.info(f"Successfully processed document {document_id} with {len(chunks)} chunks")
            return True
            
        except Exception as e:
            logger.error(f"Error processing document {file_path}: {e}")
            return False

    def search_similar_chunks(self, query: str, document_ids: List[int], k: int = 5) -> List[str]:
        """Search for similar chunks based on query."""
        try:
            if not self.embeddings or not self.vector_store:
                # Mock implementation
                mock_context = f"Mock context for query: '{query}' from documents {document_ids}"
                return [mock_context]

            # Generate query embedding
            query_embedding = self.embeddings.embed_query(query)
            query_array = np.array([query_embedding]).astype('float32')
            
            # Search FAISS index
            distances, indices = self.vector_store.search(query_array, k)
            
            # Filter by document IDs and return relevant chunks
            relevant_chunks = []
            for idx in indices[0]:
                if idx < len(self.chunk_metadata):
                    metadata = self.chunk_metadata[idx]
                    if metadata['document_id'] in document_ids:
                        relevant_chunks.append(self.document_chunks[idx])
            
            return relevant_chunks[:k]
            
        except Exception as e:
            logger.error(f"Error searching similar chunks: {e}")
            return [f"Error retrieving context for query: {query}"]

    def generate_chat_response(self, query: str, context: List[str]) -> str:
        """Generate a chat response using LLM."""
        try:
            if not self.llm:
                # Mock implementation
                return f"Mock response: Based on the provided context, here's information about '{query}'. Context: {' '.join(context[:200])}..."

            # Prepare prompt
            context_text = "\n".join(context)
            prompt = f"""Based on the following context, please answer the user's question.

Context:
{context_text}

Question: {query}

Please provide a helpful and accurate answer based on the context provided."""

            # Generate response
            response = self.llm.invoke(prompt)
            return response.content
            
        except Exception as e:
            logger.error(f"Error generating chat response: {e}")
            return f"I apologize, but I encountered an error while processing your question: {query}"

    def analyze_resume(self, resume_path: str, job_description: str = "") -> Dict[str, Any]:
        """Analyze resume and provide feedback."""
        try:
            # Extract text from resume
            file_ext = Path(resume_path).suffix.lower()
            if file_ext == '.pdf':
                resume_text = self.extract_text_from_pdf(resume_path)
            elif file_ext in ['.docx', '.doc']:
                resume_text = self.extract_text_from_docx(resume_path)
            else:
                return {"error": "Unsupported file format"}

            if not self.llm:
                # Mock implementation
                return {
                    "feedback": f"Mock feedback for resume analysis. Resume contains {len(resume_text)} characters. Job description: {job_description[:100]}...",
                    "score": 75
                }

            # Prepare analysis prompt
            prompt = f"""Please analyze the following resume and provide detailed feedback.

Resume:
{resume_text}

Job Description (if provided):
{job_description}

Please provide:
1. Overall assessment and ATS score (0-100)
2. Strengths and weaknesses
3. Specific recommendations for improvement
4. Keyword optimization suggestions

Format your response as constructive feedback."""

            # Generate analysis
            response = self.llm.invoke(prompt)
            
            # Extract score (simple heuristic)
            score = 75  # Default score
            response_text = response.content.lower()
            if "excellent" in response_text or "outstanding" in response_text:
                score = 85
            elif "good" in response_text or "strong" in response_text:
                score = 75
            elif "needs improvement" in response_text or "weak" in response_text:
                score = 60
            
            return {
                "feedback": response.content,
                "score": score
            }
            
        except Exception as e:
            logger.error(f"Error analyzing resume: {e}")
            return {"error": f"Failed to analyze resume: {str(e)}"}

    def generate_sql_from_natural_language(self, natural_query: str) -> str:
        """Generate SQL from natural language query."""
        try:
            if not self.llm:
                # Mock implementation
                return f"-- Generated SQL for: {natural_query}\nSELECT * FROM users WHERE created_at >= NOW() - INTERVAL '30 days';"

            # Database schema context
            schema_context = """
Database Schema:
- users (id, email, password_hash, created_at, updated_at)
- documents (id, user_id, filename, file_path, file_type, file_size, status, created_at)
- chat_sessions (id, user_id, document_ids, created_at)
- chat_messages (id, session_id, role, content, metadata, created_at)
- research_tasks (id, user_id, query, status, result, metadata, created_at, completed_at)
- resume_analyses (id, user_id, resume_path, job_description, feedback, score, status, created_at, completed_at)
- sql_queries (id, user_id, natural_query, generated_sql, result_data, status, created_at)
"""

            prompt = f"""{schema_context}

Convert the following natural language query to SQL:
"{natural_query}"

Please provide only the SQL query without explanations."""

            # Generate SQL
            response = self.llm.invoke(prompt)
            return response.content.strip()
            
        except Exception as e:
            logger.error(f"Error generating SQL: {e}")
            return f"-- Error generating SQL for: {natural_query}\n-- {str(e)}"

    def conduct_research(self, research_query: str) -> str:
        """Conduct research on a given topic."""
        try:
            if not self.llm:
                # Mock implementation
                return f"""Research Report: {research_query}

This is a mock research report. In a real implementation, this would involve:
1. Breaking down the research query into subtasks
2. Searching multiple information sources
3. Synthesizing findings from various sources
4. Generating a comprehensive report

Key findings would be presented here with proper citations and analysis."""

            prompt = f"""Please conduct research on the following topic and provide a comprehensive report:

Topic: {research_query}

Please provide:
1. Executive summary
2. Key findings
3. Analysis and insights
4. Conclusions and recommendations

Structure your response as a professional research report."""

            # Generate research report
            response = self.llm.invoke(prompt)
            return response.content
            
        except Exception as e:
            logger.error(f"Error conducting research: {e}")
            return f"Failed to conduct research on: {research_query}. Error: {str(e)}"

# Global AI service instance
ai_service = AIService()

if __name__ == "__main__":
    # Test the AI service
    print("AI Service initialized successfully!")
    print(f"OpenAI API Key available: {bool(ai_service.openai_api_key)}")
    print(f"Gemini API Key available: {bool(ai_service.gemini_api_key)}")
    print(f"LLM available: {ai_service.llm is not None}")
    print(f"Embeddings available: {ai_service.embeddings is not None}")

