package crawler

import (
	"fmt"
	"log"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/mmcdole/gofeed"
)

type Crawler struct {
	httpClient *resty.Client
	feedParser *gofeed.Parser
	feeds      []string
}

// NewCrawler creates a new Crawler
func New(httpClient *resty.Client, feedParser *gofeed.Parser) *Crawler {
	return &Crawler{
		httpClient: httpClient,
		feedParser: feedParser,
	}
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
			if rssData, err := c.fetchFeedData(f); err != nil {
				log.Printf("failed to fetch feed data: %v", err)
			} else {
				feed, err := c.feedParser.ParseString(rssData)
				if err != nil {
					log.Printf("failed to parse feed data: %v", err)
				}
				fmt.Printf("fetched feed: %s\n", feed.Title)
			}
		}(feed)
	}
	// Ждем завершения всех горутин.
	wg.Wait()
}

func (c *Crawler) fetchFeedData(feed string) (string, error) {
	// Fetch the feed using resty http client
	resp, err := c.httpClient.R().Get(feed)
	if err != nil {
		return "", fmt.Errorf("failed to fetch feed: %w", err)
	}

	if resp.IsError() {
		return "", fmt.Errorf("failed to fetch feed: %s", resp.Status())
	}

	return string(resp.Body()), nil
}
