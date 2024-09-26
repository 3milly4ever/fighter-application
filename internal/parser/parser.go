package parser

import (
	"regexp"
	"strconv"
)

func parseHeight(value string, fighter *Fighter) {
	// Example value: "5'11\" (180cm)"
	re := regexp.MustCompile(`(\d+)'(\d+)" \((\d+)cm\)`)
	matches := re.FindStringSubmatch(value)
	if len(matches) == 4 {
		feet, _ := strconv.ParseFloat(matches[1], 64)
		inches, _ := strconv.ParseFloat(matches[2], 64)
		cm, _ := strconv.ParseFloat(matches[3], 64)
		fighter.HeightIn = feet*12 + inches
		fighter.HeightCm = cm
	}
}

func parseWeight(value string, fighter *Fighter) {
	// Example value: "170 lbs (77 kg)"
	re := regexp.MustCompile(`(\d+) lbs \((\d+) kg\)`)
	matches := re.FindStringSubmatch(value)
	if len(matches) == 3 {
		lbs, _ := strconv.ParseFloat(matches[1], 64)
		kg, _ := strconv.ParseFloat(matches[2], 64)
		fighter.WeightLb = lbs
		fighter.WeightKg = kg
	}
}
