package crawler

import "fmt"

type Crawler struct {
	feeds []string
}

// NewCrawler creates a new Crawler
func New() *Crawler {
	return &Crawler{}
}

// AddFeed adds a new RSS feed to the crawler
func (c *Crawler) AddFeed(feed string) {
	c.feeds = append(c.feeds, feed)
}

// Start starts the crawler
func (c *Crawler) Start() {
	for _, feed := range c.feeds {
		fmt.Printf("Crawling %s\n", feed)
	}
}
