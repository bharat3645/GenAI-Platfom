#!/bin/bash

# GenAI Platform Setup Script
# This script automates the setup process for the GenAI Platform

set -e

echo "🚀 GenAI Platform Setup Script"
echo "================================"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check prerequisites
check_prerequisites() {
    echo "📋 Checking prerequisites..."
    
    # Check Go
    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ Go is not installed. Please install Go 1.21+${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Go installed:$(go version)${NC}"
    
    # Check Node.js
    if ! command -v node &> /dev/null; then
        echo -e "${RED}❌ Node.js is not installed. Please install Node.js 18+${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Node.js installed: $(node --version)${NC}"
    
    # Check npm
    if ! command -v npm &> /dev/null; then
        echo -e "${RED}❌ npm is not installed${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ npm installed: $(npm --version)${NC}"
    
    # Check PostgreSQL (optional)
    if command -v psql &> /dev/null; then
        echo -e "${GREEN}✓ PostgreSQL installed${NC}"
    else
        echo -e "${YELLOW}⚠ PostgreSQL not found locally (will use Neon/Supabase)${NC}"
    fi
    
    echo ""
}

# Setup backend
setup_backend() {
    echo "🔧 Setting up backend..."
    cd genai-platform
    
    # Install Go dependencies
    echo "📦 Installing Go dependencies..."
    go mod download
    go mod tidy
    
    # Create .env if it doesn't exist
    if [ ! -f .env ]; then
        echo "📝 Creating .env file from template..."
        cp .env.example .env
        echo -e "${YELLOW}⚠ Please edit genai-platform/.env with your configuration${NC}"
    fi
    
    # Create uploads directory
    mkdir -p uploads
    mkdir -p uploads/resumes
    
    echo -e "${GREEN}✓ Backend setup complete${NC}"
    cd ..
    echo ""
}

# Setup frontend
setup_frontend() {
    echo "🎨 Setting up frontend..."
    cd genai-frontend
    
    # Install npm dependencies
    echo "📦 Installing npm dependencies..."
    npm install --legacy-peer-deps
    
    echo -e "${GREEN}✓ Frontend setup complete${NC}"
    cd ..
    echo ""
}

# Setup database
setup_database() {
    echo "🗄️ Database setup..."
    
    if [ -z "$DATABASE_URL" ]; then
        echo -e "${YELLOW}⚠ DATABASE_URL not set in environment${NC}"
        echo "Options:"
        echo "1. Set DATABASE_URL environment variable"
        echo "2. Update genai-platform/.env file"
        echo "3. Run migrations manually: psql -f genai-platform/migrations/001_initial_schema.sql"
    else
        echo -e "${GREEN}✓ DATABASE_URL configured${NC}"
        echo "To run migrations, execute:"
        echo "  psql \$DATABASE_URL -f genai-platform/migrations/001_initial_schema.sql"
    fi
    echo ""
}

# Build backend
build_backend() {
    echo "🔨 Building backend..."
    cd genai-platform
    
    go build -o bin/server cmd/server/main.go
    
    echo -e "${GREEN}✓ Backend built successfully${NC}"
    cd ..
    echo ""
}

# Main setup
main() {
    check_prerequisites
    setup_backend
    setup_frontend
    setup_database
    
    echo "✨ Setup complete!"
    echo ""
    echo "Next steps:"
    echo "1. Configure genai-platform/.env with your database credentials"
    echo "2. Run database migrations (see instructions above)"
    echo "3. Start the backend: cd genai-platform && go run cmd/server/main.go"
    echo "4. Start the frontend: cd genai-frontend && npm run dev"
    echo ""
    echo "Or use Docker:"
    echo "  docker-compose up --build"
    echo ""
    echo "Documentation: README.md"
}

main
