package config

import (
	"fmt"
	"os"
	"time"

	"github.com/3milly4ever/fighter-application/internal/log"
	"github.com/3milly4ever/fighter-application/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is a global variable to hold the database connection
var DB *gorm.DB
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

// InitDB initializes the PostgreSQL connection using GORM
func InitDB() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}

	// Build DSN (Data Source Name) for GORM
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		AppConfigInstance.DBHost,
		AppConfigInstance.DBPort,
		AppConfigInstance.DBUser,
		AppConfigInstance.DBPassword,
		AppConfigInstance.DBName,
	)

	// Connect to the database with GORM
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Optional: Set GORM logging level
	})
	if err != nil {
		logrus.Fatalf("Error connecting to database: %v", err)
	}

	// Verify the database connection
	sqlDB, err := DB.DB()
	if err != nil {
		logrus.Fatalf("Error retrieving database instance: %v", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

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
