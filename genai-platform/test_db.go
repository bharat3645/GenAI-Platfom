package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Connection string
	connStr := "host=localhost port=5432 user=postgres password=password dbname=genai_platform sslmode=disable"

	fmt.Println("Attempting to connect to database...")
	fmt.Println("Connection string:", connStr)

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open database connection:", err)
	}
	defer db.Close()

	// Test connection
	fmt.Println("Testing connection...")
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Successfully connected to database!")

	// Try to create a test table
	fmt.Println("Creating test table...")
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS test_table (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100)
		)
	`)
	if err != nil {
		log.Fatal("Failed to create test table:", err)
	}

	fmt.Println("Test table created successfully!")
}
