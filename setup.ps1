# GenAI Platform Setup Script for Windows
# Run this script with PowerShell

Write-Host "🚀 GenAI Platform Setup Script" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

# Check prerequisites
Write-Host "📋 Checking prerequisites..." -ForegroundColor Yellow

# Check Go
if (Get-Command go -ErrorAction SilentlyContinue) {
    $goVersion = go version
    Write-Host "✓ Go installed: $goVersion" -ForegroundColor Green
} else {
    Write-Host "❌ Go is not installed. Please install Go 1.21+" -ForegroundColor Red
    exit 1
}

# Check Node.js
if (Get-Command node -ErrorAction SilentlyContinue) {
    $nodeVersion = node --version
    Write-Host "✓ Node.js installed: $nodeVersion" -ForegroundColor Green
} else {
    Write-Host "❌ Node.js is not installed. Please install Node.js 18+" -ForegroundColor Red
    exit 1
}

# Check npm
if (Get-Command npm -ErrorAction SilentlyContinue) {
    $npmVersion = npm --version
    Write-Host "✓ npm installed: $npmVersion" -ForegroundColor Green
} else {
    Write-Host "❌ npm is not installed" -ForegroundColor Red
    exit 1
}

Write-Host ""

# Setup backend
Write-Host "🔧 Setting up backend..." -ForegroundColor Yellow
Set-Location genai-platform

Write-Host "📦 Installing Go dependencies..." -ForegroundColor Cyan
go mod download
go mod tidy

# Create .env if it doesn't exist
if (!(Test-Path .env)) {
    Write-Host "📝 Creating .env file from template..." -ForegroundColor Cyan
    Copy-Item .env.example .env
    Write-Host "⚠ Please edit genai-platform\.env with your configuration" -ForegroundColor Yellow
}

# Create uploads directory
New-Item -ItemType Directory -Force -Path uploads | Out-Null
New-Item -ItemType Directory -Force -Path uploads\resumes | Out-Null

Write-Host "✓ Backend setup complete" -ForegroundColor Green
Set-Location ..
Write-Host ""

# Setup frontend
Write-Host "🎨 Setting up frontend..." -ForegroundColor Yellow
Set-Location genai-frontend

Write-Host "📦 Installing npm dependencies..." -ForegroundColor Cyan
npm install --legacy-peer-deps

Write-Host "✓ Frontend setup complete" -ForegroundColor Green
Set-Location ..
Write-Host ""

# Database setup info
Write-Host "🗄️ Database setup..." -ForegroundColor Yellow
if ($env:DATABASE_URL) {
    Write-Host "✓ DATABASE_URL configured" -ForegroundColor Green
    Write-Host "To run migrations, execute:" -ForegroundColor Cyan
    Write-Host "  psql `$env:DATABASE_URL -f genai-platform\migrations\001_initial_schema.sql"
} else {
    Write-Host "⚠ DATABASE_URL not set in environment" -ForegroundColor Yellow
    Write-Host "Options:"
    Write-Host "1. Set DATABASE_URL environment variable"
    Write-Host "2. Update genai-platform\.env file"
    Write-Host "3. Run migrations manually"
}
Write-Host ""

# Build backend
Write-Host "🔨 Building backend..." -ForegroundColor Yellow
Set-Location genai-platform

go build -o bin\server.exe cmd\server\main.go

Write-Host "✓ Backend built successfully" -ForegroundColor Green
Set-Location ..
Write-Host ""

# Final instructions
Write-Host "✨ Setup complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "1. Configure genai-platform\.env with your database credentials"
Write-Host "2. Run database migrations (see instructions above)"
Write-Host "3. Start the backend: cd genai-platform; go run cmd\server\main.go"
Write-Host "4. Start the frontend: cd genai-frontend; npm run dev"
Write-Host ""
Write-Host "Or use Docker:" -ForegroundColor Cyan
Write-Host "  docker-compose up --build"
Write-Host ""
Write-Host "Documentation: README.md" -ForegroundColor Cyan
