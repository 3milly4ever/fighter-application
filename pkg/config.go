package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// DB is a global variable to hold the database connection
var DB *sql.DB

// InitDB initializes the PostgreSQL connection
func InitDB() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Get connection info from environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Build connection string
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open connection to the database
	var errOpen error
	DB, errOpen = sql.Open("postgres", psqlInfo)
	if errOpen != nil {
		log.Fatalf("Error connecting to database: %v", errOpen)
	}

	// Ping to check if the connection is alive
	errPing := DB.Ping()
	if errPing != nil {
		log.Fatalf("Error pinging database: %v", errPing)
	}

	log.Println("Database connected successfully!")
}
