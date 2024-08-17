package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/Atviksord/BlogAggregator/internal/database"
)

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Atom    string   `xml:"atom,attr"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Text  string `xml:",chardata"`
		Title string `xml:"title"`
		Link  struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Description   string `xml:"description"`
		Generator     string `xml:"generator"`
		Language      string `xml:"language"`
		LastBuildDate string `xml:"lastBuildDate"`
		Item          []struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			PubDate     string `xml:"pubDate"`
			Guid        string `xml:"guid"`
			Description string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

// gets N feeds from DB
func (cfg *apiConfig) nextFeedGetter(n int32) ([]database.Feed, error) {
	fetchedFeed, err := cfg.DB.GetNextFeedsToFetch(context.Background(), n)
	if err != nil {
		return fetchedFeed, fmt.Errorf("Failed to Get next feed from DB %v", err)
	}

	return fetchedFeed, nil

}

// marks if feeds been fetched
func (cfg *apiConfig) feedMarker(feed []database.Feed) error {

	// Loop over the feed slice and edit timings
	for _, c := range feed {
		c.LastFetchedAt = sql.NullTime{Time: time.Now().UTC(), Valid: true}
		c.UpdatedAt = time.Now().UTC()
		err := cfg.DB.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
			LastFetchedAt: c.LastFetchedAt,
			UpdatedAt:     c.UpdatedAt,
			ID:            c.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to mark feed with id %d as fetched: %w,", c.ID, err)
		}

	}
	return nil

}

func fetchDataFromFeed(urlz string) {
	r, err := http.Get(urlz)
	if err != nil {
		fmt.Printf("Failed to get URL %v", err)
	}

	response := Rss{}
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Failed to read requestbody %v", err)
	}

	xml.Unmarshal(requestBody, &response)
}

func (cfg *apiConfig) FeedFetchWorker(n int32) {
	for {
		time.Sleep(60 * time.Second)
		// NextFeedGet get from DB
		feed, err := cfg.nextFeedGetter(n)
		if err != nil {
			fmt.Printf("Failed to get Feed from DB %v", err)
		}

		// Call feedMarker to mark as fetched
		err = cfg.feedMarker(feed)
		if err != nil {
			fmt.Printf("Failed to mark feed as fetched %v", err)
		}
		// Call fetchDataFromFeed to get feed data.
		for i := range feed {
			fetchDataFromFeed(feed[i].Url)
		}

		// Use sync.WaitGroup to spawn multiple goroutines

		var wg sync.WaitGroup
		// Placeholder for urls
		urls := []string{"url1", "url2", "url3"}

		for _, url := range urls {
			wg.Add(1)
			go func(u string) {
				defer wg.Done()
				// Simulate fetching data
				fmt.Println("Fetching", u)
			}(url)
		}

		wg.Wait()
		fmt.Println("All goroutines complete.")
	}
}
