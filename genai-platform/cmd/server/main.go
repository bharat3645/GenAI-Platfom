package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"genai-platform/internal/auth"
	"genai-platform/internal/database"
	"genai-platform/internal/handlers"
	"genai-platform/pkg/config"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Initialize handlers
	h := handlers.New(db)

	// Routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Post("/auth/login", h.Login)
		r.Post("/auth/register", h.Register)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(auth.JWTMiddleware)

			// PDF Chat routes
			r.Post("/pdf/upload", h.UploadPDF)
			r.Post("/chat/query", h.ChatQuery)

			// Graph RAG routes
			r.Post("/graph/upload", h.GraphUpload)
			r.Post("/graph/query", h.GraphQuery)

			// Research Assistant routes
			r.Post("/agent/research", h.ResearchAgent)
			r.Get("/agent/research/{id}", h.GetResearchResult)

			// Resume Feedback routes
			r.Post("/resume/upload", h.ResumeUpload)
			r.Get("/resume/feedback/{id}", h.GetResumeFeedback)

			// Text-to-SQL routes
			r.Post("/sql/query", h.SQLQuery)
		})
	})

	// Static file serving
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, r))
}

