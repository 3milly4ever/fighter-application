package scraper

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
)

type SitemapIndex struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []URL    `xml:"url"`
}

type URL struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod"`
}

func getSitemapURLs() []string {
	sitemapURLs := []string{}
	for i := 1; i <= 10; i++ {
		url := fmt.Sprintf("https://www.sherdog.com/sitemap-fighters%d.xml", i)
		sitemapURLs = append(sitemapURLs, url)
	}
	return sitemapURLs
}

func getFighterURLs() ([]string, error) {
	sitemapURLs := getSitemapURLs()
	var fighterURLs []string
	for _, sitemapURL := range sitemapURLs {
		urls, err := fetchSitemap(sitemapURL)
		if err != nil {
			log.Println("Error fetching sitemap:", err)
			continue
		}
		fighterURLs = append(fighterURLs, urls...)
	}
	return fighterURLs, nil
}

func fetchSitemap(sitemapURL string) ([]string, error) {
	resp, err := http.Get(sitemapURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch sitemap: %s", resp.Status)
	}

	var sitemap SitemapIndex
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&sitemap)
	if err != nil && err != io.EOF {
		return nil, err
	}

	var fighterURLs []string
	for _, url := range sitemap.URLs {
		fighterURLs = append(fighterURLs, url.Loc)
	}

	return fighterURLs, nil
}
