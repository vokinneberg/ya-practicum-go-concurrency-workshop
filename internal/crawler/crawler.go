package crawler

import (
	"fmt"
	"sync"
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
	// Создаем WaitGroup, чтобы дождаться завершения всех горутин.
	wg := sync.WaitGroup{}
	for _, feed := range c.feeds {
		wg.Add(1)
		// Здесь мы используем анонимную функцию, чтобы передать feed внутрь неё и напечатать его.
		go func(f string) {
			defer wg.Done()
			fmt.Printf("Crawling %s\n", f)
		}(feed)
	}
	// Ждем завершения всех горутин.
	wg.Wait()
}
