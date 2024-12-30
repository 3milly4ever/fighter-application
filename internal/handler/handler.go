package handler

import (
	"github.com/3milly4ever/fighter-application/internal/database"
	"github.com/3milly4ever/fighter-application/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// GetFighterHandler handles retrieving a fighter by name
func GetFighterHandler(c *fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Fighter name is required",
		})
	}

	fighter, err := database.GetFighter(name)
	if err != nil {
		logrus.Errorf("Error fetching fighter: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Fighter not found",
		})
	}

	return c.JSON(fighter)
}

// ListFightersHandler handles retrieving all fighters
func ListFightersHandler(c *fiber.Ctx) error {
	fighters, err := database.ListFighters()
	if err != nil {
		logrus.Errorf("Error listing fighters: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve fighters",
		})
	}

	return c.JSON(fighters)
}

// AddFighterHandler handles adding a new fighter
func AddFighterHandler(c *fiber.Ctx) error {
	var fighter model.Fighter
	if err := c.BodyParser(&fighter); err != nil {
		logrus.Errorf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err := database.InsertFighter(&fighter)
	if err != nil {
		logrus.Errorf("Error inserting fighter: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add fighter",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Fighter added successfully",
	})
}

// UpdateFighterHandler handles updating a fighter
func UpdateFighterHandler(c *fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Fighter name is required",
		})
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		logrus.Errorf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err := database.UpdateFighter(name, updates)
	if err != nil {
		logrus.Errorf("Error updating fighter: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update fighter",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Fighter updated successfully",
	})
}

// DeleteFighterHandler handles deleting a fighter
func DeleteFighterHandler(c *fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Fighter name is required",
		})
	}

	err := database.DeleteFighter(name)
	if err != nil {
		logrus.Errorf("Error deleting fighter: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete fighter",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Fighter deleted successfully",
	})
}
