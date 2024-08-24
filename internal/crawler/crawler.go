package crawler

import (
	"fmt"
)

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
		// Здесь мы используем анонимную функцию, чтобы передать feed внутрь неё и напечатать его.
		// Почему програма ничего не печатает?
		go func(f string) {
			fmt.Printf("Crawling %s\n", f)
		}(feed)
	}
}
