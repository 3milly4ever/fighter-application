package database

import (
	"log"

	"github.com/3milly4ever/fighter-application/config"
)

type Fighter struct {
	Name        string
	Age         int
	HeightCm    float64
	HeightIn    float64
	WeightKg    float64
	WeightLb    float64
	Association string
	Wins        int
	Losses      int
	KOWins      int
	SubWins     int
	DecWins     int
	KOLosses    int
	SubLosses   int
	DecLosses   int
}

// InsertFighter inserts a fighter into the PostgreSQL database
func InsertFighter(fighter *Fighter) error {
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
	_, err := config.DB.Exec(query,
		fighter.Name, fighter.Age, fighter.HeightCm, fighter.HeightIn,
		fighter.WeightKg, fighter.WeightLb, fighter.Association,
		fighter.Wins, fighter.Losses, fighter.KOWins, fighter.SubWins,
		fighter.DecWins, fighter.KOLosses, fighter.SubLosses, fighter.DecLosses)
	if err != nil {
		log.Printf("Error inserting fighter %s: %v", fighter.Name, err)
	}
	return err
}
