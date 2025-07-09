# GenAI Platform - Deployment Guide

## Overview

The GenAI Platform is a comprehensive AI-powered platform that provides multiple advanced features including PDF chat, GraphRAG, research assistance, resume feedback, and text-to-SQL conversion. This guide covers deployment and configuration.

## Architecture

### Backend (Go)
- **Framework**: Go with Gin web framework
- **Database**: PostgreSQL
- **Authentication**: JWT-based authentication
- **AI Integration**: Python bridge for LLM operations
- **File Storage**: Local filesystem with configurable upload directory

### Frontend (React)
- **Framework**: React with Vite
- **UI Library**: Tailwind CSS + shadcn/ui components
- **Routing**: React Router
- **State Management**: React Context API

### AI Services (Python)
- **LLM Integration**: OpenAI GPT / Google Gemini (configurable)
- **Vector Database**: FAISS for embeddings
- **Text Processing**: LangChain for document processing
- **File Processing**: PyPDF2 for PDF extraction, python-docx for Word documents

## Prerequisites

### System Requirements
- Ubuntu 22.04 or compatible Linux distribution
- Go 1.18 or higher
- Node.js 20.x or higher
- Python 3.11 or higher
- PostgreSQL 14 or higher

### Required Packages
```bash
# System packages
sudo apt-get update
sudo apt-get install -y golang-go nodejs npm postgresql postgresql-contrib python3 python3-pip

# Python packages
pip3 install langchain langchain-openai langchain-google-genai faiss-cpu pypdf2 python-docx
```

## Environment Configuration

### Backend Environment Variables
Create a `.env` file in the backend directory:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=genai_platform

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here

# Server Configuration
PORT=8080
UPLOAD_DIR=/path/to/uploads

# AI Service Configuration (Optional)
OPENAI_API_KEY=your-openai-api-key
GEMINI_API_KEY=your-gemini-api-key
```

### Frontend Environment Variables
Create a `.env` file in the frontend directory:

```env
VITE_API_URL=http://localhost:8080
```

## Database Setup

### PostgreSQL Installation and Configuration
```bash
# Install PostgreSQL
sudo apt-get install -y postgresql postgresql-contrib

# Start and enable PostgreSQL
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Create database and user
sudo -u postgres createdb genai_platform
sudo -u postgres psql -c "ALTER USER postgres PASSWORD 'password';"
```

### Database Schema
The application automatically creates the required tables on startup. The schema includes:

- `users` - User accounts and authentication
- `documents` - Uploaded PDF documents
- `chat_sessions` - PDF chat conversations
- `chat_messages` - Individual chat messages
- `research_tasks` - Research assistant tasks
- `resume_analyses` - Resume feedback results
- `sql_queries` - Text-to-SQL query history

## Deployment Steps

### 1. Backend Deployment

```bash
# Clone or copy the backend code
cd /path/to/genai-platform

# Install Go dependencies
go mod tidy

# Build the application
go build -o bin/server ./cmd/server

# Create necessary directories
mkdir -p uploads logs

# Set permissions
chmod +x bin/server
chmod +x ai_bridge.py

# Start the server
./bin/server
```

### 2. Frontend Deployment

```bash
# Navigate to frontend directory
cd /path/to/genai-frontend

# Install dependencies
npm install

# Build for production
npm run build

# Serve the built files (using a web server like nginx)
# Or for development
npm run dev
```

### 3. AI Service Configuration

The AI service supports both OpenAI and Google Gemini APIs. If no API keys are provided, it falls back to mock implementations for testing.

```bash
# Set environment variables for AI services
export OPENAI_API_KEY="your-openai-api-key"
export GEMINI_API_KEY="your-gemini-api-key"

# Test the AI service
cd /path/to/genai-platform
python3 ai_service.py
```

## Production Deployment

### Using systemd (Recommended)

Create a systemd service file for the backend:

```ini
# /etc/systemd/system/genai-platform.service
[Unit]
Description=GenAI Platform Backend
After=network.target postgresql.service

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/path/to/genai-platform
ExecStart=/path/to/genai-platform/bin/server
Restart=always
RestartSec=5
Environment=DB_HOST=localhost
Environment=DB_PORT=5432
Environment=DB_USER=postgres
Environment=DB_PASSWORD=password
Environment=DB_NAME=genai_platform
Environment=JWT_SECRET=your-super-secret-jwt-key-here
Environment=PORT=8080
Environment=OPENAI_API_KEY=your-openai-api-key

[Install]
WantedBy=multi-user.target
```

Enable and start the service:
```bash
sudo systemctl daemon-reload
sudo systemctl enable genai-platform
sudo systemctl start genai-platform
```

### Using Docker (Alternative)

Create a Dockerfile for the backend:

```dockerfile
FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bin/server ./cmd/server

FROM python:3.11-alpine
RUN apk add --no-cache postgresql-client
WORKDIR /app
COPY --from=builder /app/bin/server .
COPY --from=builder /app/ai_bridge.py .
COPY --from=builder /app/ai_service.py .
RUN pip install langchain langchain-openai langchain-google-genai faiss-cpu pypdf2 python-docx
EXPOSE 8080
CMD ["./server"]
```

### Nginx Configuration

Configure nginx as a reverse proxy:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # Frontend
    location / {
        root /path/to/genai-frontend/dist;
        try_files $uri $uri/ /index.html;
    }

    # Backend API
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # File uploads
    client_max_body_size 100M;
}
```

## Security Considerations

### Authentication
- JWT tokens are used for authentication
- Passwords are hashed using bcrypt
- Implement proper session management

### File Upload Security
- Validate file types and sizes
- Scan uploaded files for malware
- Store files outside the web root
- Implement proper access controls

### Database Security
- Use strong passwords
- Enable SSL connections
- Implement proper backup strategies
- Regular security updates

### API Security
- Rate limiting
- Input validation
- CORS configuration
- API key management

## Monitoring and Logging

### Application Logs
The application logs to stdout by default. Configure log rotation:

```bash
# Using logrotate
sudo nano /etc/logrotate.d/genai-platform
```

### Health Checks
Implement health check endpoints:
- `/health` - Basic health check
- `/api/v1/health` - API health check
- Database connectivity check

### Metrics
Consider implementing metrics collection:
- Request/response times
- Error rates
- Resource usage
- AI service performance

## Backup and Recovery

### Database Backup
```bash
# Create backup
pg_dump -h localhost -U postgres genai_platform > backup.sql

# Restore backup
psql -h localhost -U postgres genai_platform < backup.sql
```

### File Backup
```bash
# Backup uploaded files
tar -czf uploads-backup.tar.gz uploads/
```

## Troubleshooting

### Common Issues

1. **Database Connection Issues**
   - Check PostgreSQL service status
   - Verify connection parameters
   - Check firewall settings

2. **AI Service Errors**
   - Verify API keys are set correctly
   - Check Python dependencies
   - Monitor API rate limits

3. **File Upload Issues**
   - Check file permissions
   - Verify upload directory exists
   - Check disk space

4. **Frontend Build Issues**
   - Clear node_modules and reinstall
   - Check Node.js version compatibility
   - Verify environment variables

### Log Analysis
```bash
# View backend logs
journalctl -u genai-platform -f

# View nginx logs
tail -f /var/log/nginx/access.log
tail -f /var/log/nginx/error.log
```

## Performance Optimization

### Backend Optimization
- Database indexing
- Connection pooling
- Caching strategies
- Async processing for AI tasks

### Frontend Optimization
- Code splitting
- Asset optimization
- CDN usage
- Caching strategies

### AI Service Optimization
- Batch processing
- Model caching
- Vector index optimization
- Async task processing

## Scaling Considerations

### Horizontal Scaling
- Load balancer configuration
- Database replication
- Shared file storage
- Session management

### Vertical Scaling
- Resource monitoring
- Performance profiling
- Bottleneck identification
- Hardware upgrades

## Support and Maintenance

### Regular Maintenance Tasks
- Security updates
- Database maintenance
- Log rotation
- Backup verification
- Performance monitoring

### Update Procedures
1. Backup current deployment
2. Test updates in staging environment
3. Deploy during maintenance window
4. Verify functionality
5. Monitor for issues

## API Documentation

The platform provides RESTful APIs for all functionality:

### Authentication Endpoints
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/logout` - User logout

### Document Management
- `POST /api/v1/documents/upload` - Upload documents
- `GET /api/v1/documents` - List user documents
- `DELETE /api/v1/documents/:id` - Delete document

### Chat Functionality
- `POST /api/v1/chat/sessions` - Create chat session
- `POST /api/v1/chat/message` - Send chat message
- `GET /api/v1/chat/sessions` - List chat sessions

### Research Assistant
- `POST /api/v1/research/tasks` - Submit research task
- `GET /api/v1/research/tasks` - List research tasks
- `GET /api/v1/research/tasks/:id` - Get task result

### Resume Analysis
- `POST /api/v1/resume/upload` - Upload resume for analysis
- `GET /api/v1/resume/feedback/:id` - Get analysis result

### Text-to-SQL
- `POST /api/v1/sql/query` - Execute natural language query
- `GET /api/v1/sql/queries` - List query history

This deployment guide provides comprehensive instructions for setting up and maintaining the GenAI Platform in production environments.

