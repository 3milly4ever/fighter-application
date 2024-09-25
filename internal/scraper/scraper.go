package scraper

import (
	"encoding/xml"
	"fmt"
	"io"
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
