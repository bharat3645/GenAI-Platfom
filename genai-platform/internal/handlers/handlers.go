package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"genai-platform/internal/auth"
	"genai-platform/internal/models"
	"genai-platform/internal/services"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/lib/pq" // Import the pq library for array handling
)

type Handler struct {
	db              *sql.DB
	llmService      *services.LLMService
	fileService     *services.FileService
	graphRAGService *services.GraphRAGService
	atsService      *services.ATSAgenticService
	researchService *services.ResearchAgentService
	sqlService      *services.TextToSQLService
}

func New(db *sql.DB) *Handler {
	return &Handler{
		db:              db,
		llmService:      services.NewLLMService(),
		fileService:     services.NewFileService(),
		graphRAGService: services.NewGraphRAGService(db),
		atsService:      services.NewATSAgenticService(db),
		researchService: services.NewResearchAgentService(db),
		sqlService:      services.NewTextToSQLService(db),
	}
}

// Auth handlers
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Insert user
	var userID int
	if err := h.db.QueryRow(
		"INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id",
		req.Email, string(hashedPassword),
	).Scan(&userID); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Generate token
	token, err := auth.GenerateToken(userID, req.Email, "your-secret-key-change-in-production")
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":   token,
		"user_id": userID,
		"email":   req.Email,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := h.db.QueryRow(
		"SELECT id, email, password_hash FROM users WHERE email = $1",
		req.Email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID, user.Email, "your-secret-key-change-in-production")
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":   token,
		"user_id": user.ID,
		"email":   user.Email,
	})
}

// PDF Chat handlers
func (h *Handler) UploadPDF(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	// Parse multipart form
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	if filepath.Ext(header.Filename) != ".pdf" {
		http.Error(w, "Only PDF files are allowed", http.StatusBadRequest)
		return
	}

	// Save file
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("%d_%d_%s", userID, time.Now().Unix(), header.Filename)
	filePath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	size, err := io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Save to database
	var docID int
	if err := h.db.QueryRow(
		`INSERT INTO documents (user_id, filename, file_path, file_type, file_size) 
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		userID, header.Filename, filePath, "pdf", size,
	).Scan(&docID); err != nil {
		http.Error(w, "Failed to save document info", http.StatusInternalServerError)
		return
	}

	// Process PDF for GraphRAG (async)
	go func() {
		if err := h.graphRAGService.ProcessDocumentForGraphRAG(docID, filePath); err != nil {
			fmt.Printf("GraphRAG processing failed for document %d: %v\n", docID, err)
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"document_id": docID,
		"filename":    header.Filename,
		"status":      "processing",
		"message":     "Document uploaded. GraphRAG processing started.",
	})
}

func (h *Handler) ChatQuery(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	var req struct {
		Query       string `json:"query"`
		DocumentIDs []int  `json:"document_ids"`
		SessionID   *int   `json:"session_id,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create or get session
	sessionID := req.SessionID
	if sessionID == nil {
		var newSessionID int
		if err := h.db.QueryRow(
			"INSERT INTO chat_sessions (user_id, document_ids) VALUES ($1, $2) RETURNING id",
			userID, pq.Array(req.DocumentIDs),
		).Scan(&newSessionID); err != nil {
			http.Error(w, fmt.Sprintf("Failed to create session: %v", err), http.StatusInternalServerError)
			return
		}
		sessionID = &newSessionID
	}

	// Save user message
	if _, err := h.db.Exec(
		"INSERT INTO chat_messages (session_id, role, content) VALUES ($1, $2, $3)",
		*sessionID, "user", req.Query,
	); err != nil {
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	// Get relevant context from documents using GraphRAG hybrid retrieval
	var context string
	var err error
	if len(req.DocumentIDs) > 0 {
		// Use GraphRAG hybrid retrieval
		context, err = h.graphRAGService.HybridRetrievalQuery(req.Query, req.DocumentIDs, 5)
		if err != nil {
			// Fallback to traditional retrieval
			context, err = h.fileService.GetRelevantContext(req.DocumentIDs, req.Query)
			if err != nil {
				http.Error(w, "Failed to get context", http.StatusInternalServerError)
				return
			}
		}
	} else {
		context = "No documents selected for context."
	}

	// Generate response using LLM
	response, err := h.llmService.GenerateResponse(req.Query, context)
	if err != nil {
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
		return
	}

	// Save assistant message
	if _, err := h.db.Exec(
		"INSERT INTO chat_messages (session_id, role, content) VALUES ($1, $2, $3)",
		*sessionID, "assistant", response,
	); err != nil {
		http.Error(w, "Failed to save response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"session_id": *sessionID,
		"response":   response,
		"context":    context,
	})
}

// Graph RAG handlers
func (h *Handler) GraphUpload(w http.ResponseWriter, r *http.Request) {
	// Placeholder for GraphRAG upload
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "GraphRAG upload endpoint - to be implemented",
	})
}

func (h *Handler) GraphQuery(w http.ResponseWriter, r *http.Request) {
	// Placeholder for GraphRAG query
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "GraphRAG query endpoint - to be implemented",
	})
}

// Research Assistant handlers
func (h *Handler) ResearchAgent(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	var req struct {
		Query string `json:"query"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create research task
	var taskID int
	if err := h.db.QueryRow(
		"INSERT INTO research_tasks (user_id, query) VALUES ($1, $2) RETURNING id",
		userID, req.Query,
	).Scan(&taskID); err != nil {
		http.Error(w, "Failed to create research task", http.StatusInternalServerError)
		return
	}

	// Start research process using Multi-Agent Research Service (async)
	go func() {
		if err := h.researchService.ConductResearch(taskID, req.Query); err != nil {
			fmt.Printf("Research task failed for task %d: %v\n", taskID, err)
			// Update status to failed
			h.db.Exec(
				`UPDATE research_tasks SET status = 'failed', completed_at = NOW() WHERE id = $1`,
				taskID,
			)
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"task_id": taskID,
		"status":  "started",
		"message": "Research task created. 7-agent workflow initiated.",
	})
}

func (h *Handler) GetResearchResult(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	taskIDStr := chi.URLParam(r, "id")

	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var task models.ResearchTask
	if err := h.db.QueryRow(
		"SELECT id, query, status, result, created_at, completed_at FROM research_tasks WHERE id = $1 AND user_id = $2",
		taskID, userID,
	).Scan(&task.ID, &task.Query, &task.Status, &task.Result, &task.CreatedAt, &task.CompletedAt); err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// Resume Feedback handlers
func (h *Handler) ResumeUpload(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	// Parse multipart form
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("resume")
	if err != nil {
		http.Error(w, "Failed to get resume file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	jobDescription := r.FormValue("job_description")

	// Save file
	uploadDir := "./uploads/resumes"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("%d_%d_%s", userID, time.Now().Unix(), header.Filename)
	filePath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Save to database
	var analysisID int
	if err := h.db.QueryRow(
		`INSERT INTO resume_analyses (user_id, resume_path, job_description) 
		 VALUES ($1, $2, $3) RETURNING id`,
		userID, filePath, jobDescription,
	).Scan(&analysisID); err != nil {
		http.Error(w, "Failed to save analysis info", http.StatusInternalServerError)
		return
	}

	// Process resume using Multi-Agent ATS (async)
	go func() {
		if err := h.atsService.AnalyzeResumeAgentic(analysisID, filePath, jobDescription); err != nil {
			fmt.Printf("ATS analysis failed for resume %d: %v\n", analysisID, err)
			// Update status to failed
			h.db.Exec(
				`UPDATE resume_analyses SET status = 'failed', completed_at = NOW() WHERE id = $1`,
				analysisID,
			)
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"analysis_id": analysisID,
		"status":      "processing",
		"message":     "Resume uploaded. Multi-agent ATS analysis started.",
	})
}

func (h *Handler) GetResumeFeedback(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	analysisIDStr := chi.URLParam(r, "id")

	analysisID, err := strconv.Atoi(analysisIDStr)
	if err != nil {
		http.Error(w, "Invalid analysis ID", http.StatusBadRequest)
		return
	}

	var analysis models.ResumeAnalysis
	if err := h.db.QueryRow(
		`SELECT id, resume_path, job_description, feedback, score, status, created_at, completed_at 
		 FROM resume_analyses WHERE id = $1 AND user_id = $2`,
		analysisID, userID,
	).Scan(&analysis.ID, &analysis.ResumePath, &analysis.JobDescription,
		&analysis.Feedback, &analysis.Score, &analysis.Status, &analysis.CreatedAt, &analysis.CompletedAt); err != nil {
		http.Error(w, "Analysis not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analysis)
}

// Text-to-SQL handlers
func (h *Handler) SQLQuery(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	var req struct {
		Query string `json:"query"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// First, save the query to get an ID
	var queryID int
	if err := h.db.QueryRow(
		`INSERT INTO sql_queries (user_id, natural_query, status) 
		 VALUES ($1, $2, 'processing') RETURNING id`,
		userID, req.Query,
	).Scan(&queryID); err != nil {
		http.Error(w, "Failed to create query record", http.StatusInternalServerError)
		return
	}

	// Generate SQL using schema-aware Text-to-SQL service with triple-layer safety
	generatedSQL, genErr, validationWarnings := h.sqlService.GenerateSQLWithSafety(queryID, req.Query)
	if genErr != nil {
		// Update query with error
		h.db.Exec(
			`UPDATE sql_queries SET generated_sql = $1, status = 'failed', result_data = $2 WHERE id = $3`,
			generatedSQL, fmt.Sprintf(`{"error": "%s"}`, genErr.Error()), queryID,
		)
		http.Error(w, "Failed to generate SQL: "+genErr.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the safe SQL query
	results, execErr := h.sqlService.ExecuteSafeSQL(generatedSQL, queryID)
	if execErr != nil {
		// Update query with execution error
		h.db.Exec(
			`UPDATE sql_queries SET generated_sql = $1, status = 'failed', result_data = $2 WHERE id = $3`,
			generatedSQL, fmt.Sprintf(`{"error": "%s"}`, execErr.Error()), queryID,
		)
		http.Error(w, "SQL execution failed: "+execErr.Error(), http.StatusBadRequest)
		return
	}

	// Get explanation
	explanation := h.sqlService.ExplainSQL(generatedSQL)

	// Marshal results to JSON
	resultDataJSON, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "Failed to marshal result data", http.StatusInternalServerError)
		return
	}

	// Update query with success
	h.db.Exec(
		`UPDATE sql_queries SET generated_sql = $1, result_data = $2, status = 'success' WHERE id = $3`,
		generatedSQL, resultDataJSON, queryID,
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"query_id":    queryID,
		"sql":         generatedSQL,
		"explanation": explanation,
		"results":     results,
		"row_count":   len(results),
		"warnings":    validationWarnings,
		"safety":      "triple-layer-validated",
	})
}
