package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
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

func (cfg *apiConfig) scraperMain() {
	fmt.Println("test")

}

// gets N feeds from DB
func (cfg *apiConfig) nextFeedGetter(n int) {
	//fetchedFeed, err := cfg.DB.GetNextFeedsToFetch()

}

// marks if feeds been fetched
func (cfg *apiConfig) feedMarker() {

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

func FeedFetchWorker(n int) {
	for {
		time.Tick(60)
		// NextFeedGet get from DB
		// Call feedMarker to mark as fetched
		// Call fetchDataFromFeed to get feed data.
		// Use sync.WaitGroup to spawn multiple goroutines
	}
}
