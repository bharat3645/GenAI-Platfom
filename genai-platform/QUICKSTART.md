Quick start — run the full stack locally with Docker Compose

This guide gets you a working local environment for development. It starts
Postgres (with migrations), Mongo (optional vector store), the backend and a
static frontend image. You'll run these locally via Docker Compose.

Prerequisites
- Docker Desktop (or another Docker engine) installed and running
- Git and PowerShell (you're on Windows)

1) Ensure environment
- Confirm `genai-platform/.env` is present (we've set sensible defaults pointing
  to the compose services).

2) Start services
Open PowerShell in the repository root (where `docker-compose.yml` lives) and run:

```powershell
# Build and start services in the background
docker-compose up -d --build

# Show running containers
docker ps --filter "name=genai"
```

3) Check logs and health

```powershell
# Tail backend logs
docker logs -f genai-backend

# Check DB readiness
docker logs genai-db
```

4) Frontend
- The Vite dev server runs at http://localhost:5173 (when running dev locally)
- The Dockerized frontend (if you used compose) will be served at http://localhost:5173

5) Test a simple API call

```powershell
# Check the backend is reachable
Invoke-WebRequest -UseBasicParsing http://localhost:8080/api/v1/health
```

Notes
- If you prefer not to use Docker, you can run the backend locally (requires Go 1.21+):

```powershell
# From genai-platform folder
go run ./cmd/server
```

- If you use external DBs (Neon/Supabase), replace `DATABASE_URL` in `.env` with the provided connection string.

If anything fails, paste the backend logs here and I'll help diagnose the issue.