package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"genai-platform/internal/auth"
	"genai-platform/internal/database"
	"genai-platform/internal/handlers"
	mongoclient "genai-platform/internal/mongo"
	"genai-platform/pkg/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Optionally connect to MongoDB if MONGO_URI is provided. This is useful for
	// using MongoDB for vector stores or other services while keeping Postgres
	// as primary relational DB. Full migration to Mongo requires refactoring
	// data access to use Mongo collections.
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI != "" {
		log.Println("Attempting to connect to MongoDB...")
		mc, err := mongoclient.Initialize(mongoURI)
		if err != nil {
			log.Println("Warning: failed to connect to MongoDB:", err)
		} else {
			// keep the client alive for the lifetime of the server
			defer func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := mc.Disconnect(ctx); err != nil {
					log.Println("Warning: failed to disconnect MongoDB client:", err)
				}
			}()
			log.Println("Connected to MongoDB")
		}
	}

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
		// Health endpoint for quick checks
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})
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
