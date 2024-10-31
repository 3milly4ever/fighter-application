package config

import (
	"database/sql"
	"fmt"
	"os"

	log "github.com/3milly4ever/fighter-application/internal/log"
	"github.com/3milly4ever/fighter-application/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// DB is a global variable to hold the database connection
var DB *sql.DB
var App *fiber.App

type AppConfig struct {
	ServerIP   string
	ServerPort string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

// Global variable to store the configuration
var AppConfigInstance *AppConfig

// InitServer initializes and returns a new Fiber app with middleware and configuration
func InitServer() *fiber.App {
	// Load configuration
	LoadConfig()

	// Initialize the logger
	log.InitLogger()

	// Create a new Fiber app
	app := fiber.New()

	// Apply CORS middleware
	app.Use(middleware.CORS())

	return app
}

// InitDB initializes the PostgreSQL connection
func InitDB() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("Error loading .env file")
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
		logrus.Fatalf("Error connecting to database: %v", errOpen)
	}

	// Ping to check if the connection is alive
	errPing := DB.Ping()
	if errPing != nil {
		logrus.Fatalf("Error pinging database: %v", errPing)
	}

	logrus.Println("Database connected successfully!")
}

// LoadConfig loads environment variables and stores them in AppConfigInstance
func LoadConfig() {
	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}

	// Assign environment variables to the AppConfig struct
	AppConfigInstance = &AppConfig{
		ServerIP:   getEnv("SERVER_IP", "127.0.0.1"), // Default to localhost if not set
		ServerPort: getEnv("SERVER_PORT", "8080"),    // Default to port 8080
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "fighter_app"),
	}
}

// Helper function to fetch environment variables with a default fallback
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
