package model

import (
	"time"

	"github.com/google/uuid"
)

type Fighter struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string
	Age         int
	HeightCm    float64
	HeightIn    float64
	WeightKg    float64
	WeightLb    float64
	ReachIn     float64
	ReachCm     float64
	Association string
	Record      string
	Location    string
	Wins        int
	Losses      int
	KOWins      int
	SubWins     int
	DecWins     int
	KOLosses    int
	SubLosses   int
	DecLosses   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
