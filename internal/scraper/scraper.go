package scraper

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/3milly4ever/fighter-application/internal/database"
	"github.com/3milly4ever/fighter-application/internal/parser"
	"github.com/sirupsen/logrus"
)

// Struct to parse the sitemap
type SitemapIndex struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []URL    `xml:"url"`
}

type URL struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod"`
}

// Get the Sherdog sitemap URLs
func getSitemapURLs() []string {
	sitemapURLs := []string{}
	for i := 1; i <= 10; i++ {
		url := fmt.Sprintf("https://www.sherdog.com/sitemap-fighters%d.xml", i)
		sitemapURLs = append(sitemapURLs, url)
	}
	return sitemapURLs
}

// Get the fighter URLs from all sitemaps
func GetFighterURLs() ([]string, error) {
	sitemapURLs := getSitemapURLs()
	var fighterURLs []string
	for _, sitemapURL := range sitemapURLs {
		logrus.Infof("Fetching sitemap from: %s", sitemapURL)
		urls, err := fetchSitemap(sitemapURL)
		if err != nil {
			logrus.Errorf("Error fetching sitemap from %s: %v", sitemapURL, err)
			continue
		}
		fighterURLs = append(fighterURLs, urls...)
		logrus.Infof("Fetched %d fighter URLs from sitemap: %s", len(urls), sitemapURL)
	}
	return fighterURLs, nil
}

// Fetch a single sitemap and extract the URLs
func fetchSitemap(sitemapURL string) ([]string, error) {
	resp, err := http.Get(sitemapURL)
	if err != nil {
		logrus.Errorf("Error making HTTP request to %s: %v", sitemapURL, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Failed to fetch sitemap from %s: %s", sitemapURL, resp.Status)
		return nil, fmt.Errorf("failed to fetch sitemap: %s", resp.Status)
	}

	var sitemap SitemapIndex
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&sitemap)
	if err != nil && err != io.EOF {
		logrus.Errorf("Error decoding XML from sitemap: %v", err)
		return nil, err
	}

	var fighterURLs []string
	for _, url := range sitemap.URLs {
		logrus.Infof("Found fighter URL: %s", url.Loc)
		fighterURLs = append(fighterURLs, url.Loc)
	}

	return fighterURLs, nil
}

// ScrapeFighters scrapes the data from each fighter URL and stores it in the DB
func ScrapeFighters() {
	fighterURLs, err := GetFighterURLs()
	if err != nil {
		logrus.Fatal("Failed to get fighter URLs:", err)
	}

	logrus.Infof("Starting to scrape %d fighter profiles", len(fighterURLs))

	for _, url := range fighterURLs {
		logrus.Infof("Scraping fighter profile from: %s", url)

		// Scrape data from each fighter profile using the parser package
		fighterData, err := parser.ScrapeFighterProfile(url)
		if err != nil {
			logrus.Errorf("Error scraping fighter profile from %s: %v", url, err)
			continue
		}

		// Insert the fighter data into the database
		err = database.InsertFighter(fighterData)
		if err != nil {
			logrus.Errorf("Error inserting fighter %s into the database: %v", fighterData.Name, err)
		} else {
			logrus.Infof("Successfully inserted fighter %s into the database", fighterData.Name)
		}
	}

	logrus.Info("Finished scraping fighter profiles")
}
