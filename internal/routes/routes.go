package routes

import (
	"github.com/3milly4ever/fighter-application/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Fighter routes
	app.Get("/fighters", handler.ListFightersHandler)           // Get all fighters
	app.Get("/fighters/:name", handler.GetFighterHandler)       // Get a fighter by name
	app.Post("/fighters", handler.AddFighterHandler)            // Add a new fighter
	app.Put("/fighters/:name", handler.UpdateFighterHandler)    // Update a fighter
	app.Delete("/fighters/:name", handler.DeleteFighterHandler) // Delete a fighter
}
