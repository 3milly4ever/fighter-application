package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/3milly4ever/fighter-application/internal/model"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

// ScrapeFighterProfile scrapes fighter data from a given Sherdog fighter profile URL
func ScrapeFighterProfile(url string) (*model.Fighter, error) {
	logrus.Infof("Starting to scrape fighter profile: %s", url)

	// Create a new instance of the Fighter struct from the model package
	fighter := &model.Fighter{}

	// Create a new collector (Colly)
	c := colly.NewCollector()

	// Scrape fighter's name
	c.OnHTML("span.fn", func(e *colly.HTMLElement) {
		fighter.Name = e.Text
		logrus.Infof("Scraped fighter name: %s", fighter.Name)
	})

	// Scrape fighter's physical attributes and other information
	c.OnHTML("div.bio li", func(e *colly.HTMLElement) {
		key := e.ChildText("strong")
		value := strings.TrimSpace(e.Text[len(key):]) // Remove the key part and trim whitespace

		switch key {
		case "Age:":
			fighter.Age = parseToInt(value)
			logrus.Infof("Scraped fighter age: %d", fighter.Age)
		case "Height:":
			fighter.HeightIn, fighter.HeightCm = parseHeight(value)
			logrus.Infof("Scraped fighter height: %f in / %f cm", fighter.HeightIn, fighter.HeightCm)
		case "Weight:":
			fighter.WeightLb, fighter.WeightKg = parseWeight(value)
			logrus.Infof("Scraped fighter weight: %f lbs / %f kg", fighter.WeightLb, fighter.WeightKg)
		case "Association:":
			fighter.Association = value
			logrus.Infof("Scraped fighter association: %s", fighter.Association)
		}
	})

	// Scrape win/loss records
	c.OnHTML("div.left_side div.record", func(e *colly.HTMLElement) {
		fighter.Wins = parseToInt(e.ChildText("span.win"))
		fighter.Losses = parseToInt(e.ChildText("span.loss"))
		logrus.Infof("Scraped fighter record: %d wins / %d losses", fighter.Wins, fighter.Losses)
	})

	// Scrape details on method of wins and losses (KO, submission, decision)
	c.OnHTML("div.module.fight_history tr", func(e *colly.HTMLElement) {
		result := e.ChildText("td.result")
		method := e.ChildText("td.method")

		if result == "win" {
			parseWinMethod(method, fighter)
			logrus.Infof("Updated win method: %s", method)
		} else if result == "loss" {
			parseLossMethod(method, fighter)
			logrus.Infof("Updated loss method: %s", method)
		}
	})

	// Visit the fighter profile URL
	err := c.Visit(url)
	if err != nil {
		logrus.Errorf("Error visiting fighter URL: %v", err)
		return nil, err
	}

	// Return the populated fighter struct
	logrus.Infof("Successfully scraped fighter profile: %s", fighter.Name)
	return fighter, nil
}

// Helper function to parse the height string into inches and centimeters
func parseHeight(value string) (float64, float64) {
	// Example: "5'11\" (180cm)"
	re := regexp.MustCompile(`(\d+)'(\d+)" \((\d+)cm\)`)
	matches := re.FindStringSubmatch(value)
	if len(matches) == 4 {
		feet, _ := strconv.Atoi(matches[1])
		inches, _ := strconv.Atoi(matches[2])
		heightIn := float64(feet*12 + inches)
		heightCm, _ := strconv.ParseFloat(matches[3], 64)
		return heightIn, heightCm
	}
	logrus.Warnf("Could not parse height: %s", value)
	return 0, 0
}

// Helper function to parse the weight string into pounds and kilograms
func parseWeight(value string) (float64, float64) {
	// Example: "170 lbs (77 kg)"
	re := regexp.MustCompile(`(\d+) lbs \((\d+) kg\)`)
	matches := re.FindStringSubmatch(value)
	if len(matches) == 3 {
		weightLb, _ := strconv.ParseFloat(matches[1], 64)
		weightKg, _ := strconv.ParseFloat(matches[2], 64)
		return weightLb, weightKg
	}
	logrus.Warnf("Could not parse weight: %s", value)
	return 0, 0
}

// Helper function to update the fighter's win statistics
func parseWinMethod(method string, fighter *model.Fighter) {
	if strings.Contains(method, "KO") || strings.Contains(method, "TKO") {
		fighter.KOWins++
	} else if strings.Contains(method, "Submission") {
		fighter.SubWins++
	} else if strings.Contains(method, "Decision") {
		fighter.DecWins++
	}
}

// Helper function to update the fighter's loss statistics
func parseLossMethod(method string, fighter *model.Fighter) {
	if strings.Contains(method, "KO") || strings.Contains(method, "TKO") {
		fighter.KOLosses++
	} else if strings.Contains(method, "Submission") {
		fighter.SubLosses++
	} else if strings.Contains(method, "Decision") {
		fighter.DecLosses++
	}
}

// Helper function to convert a string to an integer
func parseToInt(value string) int {
	num, err := strconv.Atoi(value)
	if err != nil {
		logrus.Warnf("Could not parse integer from value: %s", value)
		return 0
	}
	return num
}
