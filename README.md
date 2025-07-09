# GenAI Platform

## Overview

GenAI Platform is a comprehensive, AI-powered system for document processing, research, and data analysis. It features PDF chat, GraphRAG, research assistant, resume feedback, and text-to-SQL conversion, all powered by modern LLMs and advanced backend services.

---

## Features
- **Multi-PDF Chat**: Upload and chat with multiple PDFs using RAG and semantic search.
- **GraphRAG**: Extract entities/relationships and build knowledge graphs for advanced Q&A.
- **Research Assistant**: Autonomous AI agent for research and synthesis.
- **Resume Feedback**: ATS scoring and AI feedback for resumes.
- **Text-to-SQL**: Convert natural language to SQL and run queries.
- **User Authentication**: JWT-based, secure.
- **Admin Dashboard**: Usage stats, user management, and more.

---

## Architecture
- **Frontend**: React (Vite), Tailwind CSS, shadcn/ui, React Router
- **Backend**: Go (Gin/Chi), JWT auth, PostgreSQL, file storage, API gateway
- **AI Service**: Python (LangChain, OpenAI, Gemini, FAISS, PyPDF2, python-docx)
- **Databases**: PostgreSQL (metadata), FAISS (vectors), Neo4j (knowledge graphs)
- **File Storage**: Local (dev), AWS S3 (prod)
- **Messaging/Email**: SendGrid
- **Scheduler**: Cron jobs
- **Agent Orchestrator**: Go FSM

### Data Flow
1. Frontend sends requests to Go backend API gateway
2. API gateway routes to microservices (PDF, GraphRAG, etc.)
3. Microservices interact with databases and LLM APIs
4. File uploads handled by file service
5. Async tasks via Go routines/message queue
6. LLM interactions via centralized service

---

## Local Setup

### Prerequisites
- Go 1.18+
- Node.js 14+/npm 6+
- Python 3.11+
- PostgreSQL 12+
- Git
- jq, pandoc (optional)

### 1. Clone the Repository
```bash
git clone https://github.com/bharat3645/GenAI.git
cd genai-platform-local
```

### 2. Backend Setup (Go)
- Ensure PostgreSQL is running and create the database/user:
```bash
sudo -u postgres psql -c "ALTER USER postgres PASSWORD 'password';"
sudo -u postgres createdb genai_platform
```
- Configure `.env` in `genai-platform/`:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=genai_platform
JWT_SECRET=your-super-secret-jwt-key-here
PORT=8080
UPLOAD_DIR=uploads
OPENAI_API_KEY=your-openai-api-key
GEMINI_API_KEY=your-gemini-api-key
```
- Build and run:
```bash
cd genai-platform
go mod tidy
go build -o bin/server_local ./cmd/server
./bin/server_local
```

### 3. Frontend Setup (React)
- Configure `.env` in `genai-frontend/`:
```
VITE_API_BASE_URL=http://localhost:8080/api/v1
```
- Install and run:
```bash
cd genai-frontend
npm install
npm run dev
```

### 4. AI Service (Python)
- Install dependencies:
```bash
cd genai-platform
pip3 install langchain langchain-openai langchain-google-genai faiss-cpu pypdf2 python-docx
```
- Set API keys as env vars if using real LLMs:
```bash
export OPENAI_API_KEY=your_openai_api_key
export GEMINI_API_KEY=your_gemini_api_key
```
- Run for testing:
```bash
python3 ai_service.py
```

---

## Production Deployment
- See `DEPLOYMENT_GUIDE.md` for full details (systemd, Docker, Nginx, etc.)
- Example: Build Go backend, serve frontend with Nginx, run Python AI service, configure environment variables, set up PostgreSQL.

---

## Usage Guide

### 1. Register/Login
- Visit frontend URL, sign up, and log in.

### 2. PDF Chat
- Go to "PDF Chat", upload PDFs, chat with the AI about their content.

### 3. GraphRAG
- Go to "GraphRAG", upload docs, explore knowledge graphs and relationships.

### 4. Research Assistant
- Enter a research question, start a task, and review the report.

### 5. Resume Feedback
- Upload a resume (and job description), get feedback and ATS score.

### 6. Text-to-SQL
- Enter a natural language query, view generated SQL and results.

---

## API Endpoints (Summary)
- `POST /api/v1/auth/register` - Register
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/logout` - Logout
- `POST /api/v1/documents/upload` - Upload document
- `GET /api/v1/documents` - List documents
- `DELETE /api/v1/documents/:id` - Delete document
- `POST /api/v1/chat/sessions` - Create chat session
- `POST /api/v1/chat/message` - Send chat message
- `GET /api/v1/chat/sessions` - List chat sessions
- `POST /api/v1/research/tasks` - Submit research task
- `GET /api/v1/research/tasks` - List research tasks
- `GET /api/v1/research/tasks/:id` - Get task result
- `POST /api/v1/resume/upload` - Upload resume
- `GET /api/v1/resume/feedback/:id` - Get resume feedback
- `POST /api/v1/sql/query` - Execute SQL query
- `GET /api/v1/sql/queries` - List SQL queries

---

## Troubleshooting
- **Database issues**: Check PostgreSQL status, credentials, firewall.
- **AI service errors**: Check Python dependencies, API keys, rate limits.
- **File upload issues**: Check permissions, directory existence, disk space.
- **Frontend issues**: Clear node_modules, check Node.js version, verify env vars.
- **Logs**: Backend logs to stdout; use `journalctl` or logrotate for production.

---

## Security & Best Practices
- JWT auth, bcrypt password hashing
- File validation, malware scanning, access controls
- Strong DB passwords, SSL, backups
- Rate limiting, input validation, CORS, API key management
- Regular updates, log rotation, health checks

---

## Scaling & Performance
- Horizontal scaling of microservices
- DB replication, vector DB distribution, Neo4j clustering
- Async processing, caching (Redis), load balancing
- Metrics: request times, error rates, resource usage

---

## Support & Contribution
- For help, see this README, `USER_GUIDE.md`, or open an issue.
- Contributions welcome! Fork, branch, and submit PRs.
- For feedback, use the platform or contact maintainers.

---

**Platform Version**: 1.0  
**Maintainers**: [Your Name/Org]  
**License**: [Your License]



