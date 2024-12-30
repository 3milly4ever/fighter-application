package scraper

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

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

var httpClient = &http.Client{
	Timeout: 10 * time.Second, // Set a timeout for HTTP requests
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

// Get the fighter URLs from all sitemaps concurrently
func GetFighterURLs() ([]string, error) {
	sitemapURLs := getSitemapURLs()
	var fighterURLs []string
	var mu sync.Mutex // Protects shared slice
	var wg sync.WaitGroup
	var errors []error
	var errorMu sync.Mutex // Protects errors slice

	for _, sitemapURL := range sitemapURLs {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			logrus.Infof("Fetching sitemap from: %s", url)
			urls, err := fetchSitemapWithRetry(url, 3)
			if err != nil {
				logrus.Errorf("Error fetching sitemap from %s: %v", url, err)
				errorMu.Lock()
				errors = append(errors, err)
				errorMu.Unlock()
				return
			}
			mu.Lock()
			fighterURLs = append(fighterURLs, urls...)
			mu.Unlock()
			logrus.Infof("Fetched %d fighter URLs from sitemap: %s", len(urls), url)
		}(sitemapURL)
	}

	wg.Wait()
	if len(errors) > 0 {
		logrus.Warnf("Encountered %d errors during sitemap fetching", len(errors))
		return fighterURLs, fmt.Errorf("encountered %d errors during sitemap fetching", len(errors))
	}
	return fighterURLs, nil
}

func fetchSitemapWithRetry(sitemapURL string, retries int) ([]string, error) {
	for i := 0; i < retries; i++ {
		resp, err := httpClient.Get(sitemapURL)
		if err != nil {
			logrus.Errorf("Error making HTTP request to %s: %v", sitemapURL, err)
			return nil, err
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			logrus.Warnf("Rate limited. Retrying %s in %d seconds...", sitemapURL, (i+1)*2)
			time.Sleep(time.Duration((i+1)*2) * time.Second) // Exponential backoff
			continue
		}

		if resp.StatusCode != http.StatusOK {
			logrus.Errorf("Failed to fetch sitemap from %s: %s", sitemapURL, resp.Status)
			return nil, fmt.Errorf("failed to fetch sitemap: %s", resp.Status)
		}

		var sitemap SitemapIndex
		decoder := xml.NewDecoder(resp.Body)
		err = decoder.Decode(&sitemap)
		resp.Body.Close()
		if err != nil && err != io.EOF {
			logrus.Errorf("Error decoding XML from sitemap: %v", err)
			return nil, err
		}

		var fighterURLs []string
		for _, url := range sitemap.URLs {
			fighterURLs = append(fighterURLs, url.Loc)
		}

		return fighterURLs, nil
	}

	return nil, fmt.Errorf("failed to fetch sitemap after %d retries", retries)
}

func ScrapeFighters() {
	fighterURLs, err := GetFighterURLs()
	if err != nil {
		logrus.Warnf("Some sitemaps failed to fetch: %v", err)
	}

	if len(fighterURLs) == 0 {
		logrus.Fatal("No fighter URLs fetched. Exiting...")
		return
	}

	logrus.Infof("Starting to scrape %d fighter profiles", len(fighterURLs))

	var wg sync.WaitGroup
	for _, url := range fighterURLs {
		wg.Add(1)
		go func(fighterURL string) {
			defer wg.Done()
			logrus.Infof("Scraping fighter profile from: %s", fighterURL)

			// Scrape data from each fighter profile using the parser package
			fighterData, err := parser.ScrapeFighterProfile(fighterURL)
			if err != nil {
				logrus.Errorf("Error scraping fighter profile from %s: %v", fighterURL, err)
				return
			}

			// Insert the fighter data into the database
			err = database.InsertFighter(fighterData)
			if err != nil {
				logrus.Errorf("Error inserting fighter %s into the database: %v", fighterData.Name, err)
			} else {
				logrus.Infof("Successfully inserted fighter %s into the database", fighterData.Name)
			}
		}(url)
	}

	wg.Wait()
	logrus.Info("Finished scraping fighter profiles")
}
