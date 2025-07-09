package config

import (
	"os"
)

type Config struct {
	DatabaseURL    string
	JWTSecret      string
	OpenAIAPIKey   string
	GeminiAPIKey   string
	SendGridAPIKey string
	Port           string
}

func Load() *Config {
	return &Config{
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:password@localhost/genai_platform?sslmode=disable"),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		OpenAIAPIKey:   getEnv("OPENAI_API_KEY", ""),
		GeminiAPIKey:   getEnv("GEMINI_API_KEY", ""),
		SendGridAPIKey: getEnv("SENDGRID_API_KEY", ""),
		Port:           getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

