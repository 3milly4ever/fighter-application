package database

import (
	"github.com/3milly4ever/fighter-application/internal/model"
	config "github.com/3milly4ever/fighter-application/pkg"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

// GetFighter retrieves a fighter by name from the PostgreSQL database
func GetFighter(name string) (*model.Fighter, error) {
	var fighter model.Fighter
	err := config.DB.Where("name = ?", name).First(&fighter).Error
	if err != nil {
		logrus.Errorf("Error retrieving fighter %s: %v", name, err)
		return nil, err
	}

	logrus.Infof("Successfully retrieved fighter %s", name)
	return &fighter, nil
}

// ListFighters retrieves all fighters from the PostgreSQL database
func ListFighters() ([]model.Fighter, error) {
	var fighters []model.Fighter
	err := config.DB.Find(&fighters).Error
	if err != nil {
		logrus.Errorf("Error listing fighters: %v", err)
		return nil, err
	}

	logrus.Infof("Successfully retrieved %d fighters", len(fighters))
	return fighters, nil
}

// DeleteFighter deletes a fighter by name from the PostgreSQL database
func DeleteFighter(name string) error {
	err := config.DB.Where("name = ?", name).Delete(&model.Fighter{}).Error
	if err != nil {
		logrus.Errorf("Error deleting fighter %s: %v", name, err)
		return err
	}

	logrus.Infof("Successfully deleted fighter %s", name)
	return nil
}

// UpdateFighter updates specific fields of a fighter in the PostgreSQL database
func UpdateFighter(name string, updates map[string]interface{}) error {
	err := config.DB.Model(&model.Fighter{}).Where("name = ?", name).Updates(updates).Error
	if err != nil {
		logrus.Errorf("Error updating fighter %s: %v", name, err)
		return err
	}

	logrus.Infof("Successfully updated fighter %s", name)
	return nil
}

// BulkInsertFighters inserts or updates multiple fighters in the PostgreSQL database
func BulkInsertFighters(fighters []model.Fighter) error {
	for _, fighter := range fighters {
		err := InsertFighter(&fighter)
		if err != nil {
			logrus.Errorf("Error inserting/updating fighter %s in bulk operation: %v", fighter.Name, err)
			return err
		}
	}

	logrus.Infof("Successfully inserted/updated %d fighters", len(fighters))
	return nil
}

// InsertFighter inserts or updates a fighter in the PostgreSQL database using GORM
func InsertFighter(fighter *model.Fighter) error {
	// Use GORM's Clauses method with OnConflict to define upsert behavior
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
