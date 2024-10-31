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
	// Load configuration
	config.LoadConfig()

	// Initialize the logger
	log.InitLogger()

	// Create a new Fiber app
	app := fiber.New()

	// Apply the CORS middleware from the middleware package
	app.Use(middleware.CORS())

	// Optionally, add Fiber's built-in logging middleware to log HTTP requests
	app.Use(logger.New()) // Uncomment if you'd like to log requests

	// Set up routes
	// routes.Setup(app)

	// Log server start
	logrus.Infof("Starting server on %s:%s", config.AppConfigInstance.ServerIP, config.AppConfigInstance.ServerPort)

	// Start server in a goroutine so that it doesn't block
	go func() {
		if err := app.Listen(config.AppConfigInstance.ServerIP + ":" + config.AppConfigInstance.ServerPort); err != nil {
			logrus.Fatal("Server failed to start: ", err)
		}
	}()

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
