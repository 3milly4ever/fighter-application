package database

import (
	"context"
	"time"

	"github.com/3milly4ever/fighter-application/internal/model"
	config "github.com/3milly4ever/fighter-application/pkg"
	"github.com/sirupsen/logrus"
)

// InsertFighter inserts a fighter into the PostgreSQL database
func InsertFighter(fighter *model.Fighter) error {
	// Define the SQL query with ON CONFLICT for upsert behavior
	query := `
    INSERT INTO fighters (
        name, age, height_cm, height_in, weight_kg, weight_lb,
        association, wins, losses, ko_wins, sub_wins, dec_wins,
        ko_losses, sub_losses, dec_losses
    ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
    ON CONFLICT (name) DO UPDATE SET
        age = EXCLUDED.age,
        height_cm = EXCLUDED.height_cm,
        height_in = EXCLUDED.height_in,
        weight_kg = EXCLUDED.weight_kg,
        weight_lb = EXCLUDED.weight_lb,
        association = EXCLUDED.association,
        wins = EXCLUDED.wins,
        losses = EXCLUDED.losses,
        ko_wins = EXCLUDED.ko_wins,
        sub_wins = EXCLUDED.sub_wins,
        dec_wins = EXCLUDED.dec_wins,
        ko_losses = EXCLUDED.ko_losses,
        sub_losses = EXCLUDED.sub_losses,
        dec_losses = EXCLUDED.dec_losses;
    `

	// Create a context with a timeout for the SQL execution
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Execute the SQL query
	_, err := config.DB.ExecContext(ctx, query,
		fighter.Name, fighter.Age, fighter.HeightCm, fighter.HeightIn,
		fighter.WeightKg, fighter.WeightLb, fighter.Association,
		fighter.Wins, fighter.Losses, fighter.KOWins, fighter.SubWins,
		fighter.DecWins, fighter.KOLosses, fighter.SubLosses, fighter.DecLosses)

	if err != nil {
		logrus.Errorf("Error inserting/updating fighter %s: %v", fighter.Name, err)
		return err
	}

	logrus.Infof("Successfully inserted/updated fighter %s", fighter.Name)
	return nil
}
