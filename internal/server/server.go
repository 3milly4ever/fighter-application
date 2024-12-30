package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/3milly4ever/fighter-application/internal/log"
	"github.com/3milly4ever/fighter-application/internal/middleware"
	config "github.com/3milly4ever/fighter-application/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
)

func SetupAndRun() {
	// Load configuration and logger
	config.LoadConfig()
	log.InitLogger()

	// Initialize the Fiber app
	app := initializeApp()

	// Start the server in a goroutine
	startServer(app)

	// Graceful Shutdown
	gracefulShutdown(app)
}

// Graceful Shutdown handles termination signals
func gracefulShutdown(app *fiber.App) {
	// Listen for termination signals (e.g., SIGINT, SIGTERM)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-c

	logrus.Info("Shutting down server...")

	// Create a context with a timeout to ensure ongoing requests are completed
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logrus.Fatal("Server shutdown failed: ", err)
	}

	logrus.Info("Server successfully shut down")
}

// initializeApp sets up the Fiber app with middleware and routes
func initializeApp() *fiber.App {
	// Create a new Fiber app
	app := fiber.New()

	// Apply global middleware
	app.Use(middleware.CORS())
	app.Use(logger.New()) // Logs all HTTP requests

	// Set up routes (assuming a separate routes package)
	//routes.Setup(app)

	return app
}

// startServer starts the Fiber server
func startServer(app *fiber.App) {
	go func() {
		addr := config.AppConfigInstance.ServerIP + ":" + config.AppConfigInstance.ServerPort
		logrus.Infof("Starting server on %s", addr)

		if err := app.Listen(addr); err != nil {
			logrus.Fatalf("Server failed to start: %v", err)
		}
	}()
}
