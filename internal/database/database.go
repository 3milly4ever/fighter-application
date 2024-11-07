package database

import (
	"github.com/3milly4ever/fighter-application/internal/model"
	config "github.com/3milly4ever/fighter-application/pkg"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

// InsertFighter inserts or updates a fighter in the PostgreSQL database using GORM
func InsertFighter(fighter *model.Fighter) error {
	// Use GORM's `Clauses` method with `OnConflict` to define upsert behavior
	err := config.DB.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "name"}}, // Specify "name" as the unique key for conflict
			DoUpdates: clause.AssignmentColumns([]string{
				"age", "height_cm", "height_in", "weight_kg", "weight_lb",
				"association", "wins", "losses", "ko_wins", "sub_wins",
				"dec_wins", "ko_losses", "sub_losses", "dec_losses",
			}),
		},
	).Create(fighter).Error

	if err != nil {
		logrus.Errorf("Error inserting/updating fighter %s: %v", fighter.Name, err)
		return err
	}

	logrus.Infof("Successfully inserted/updated fighter %s", fighter.Name)
	return nil
}
