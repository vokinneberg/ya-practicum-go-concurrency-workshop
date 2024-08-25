package crawler

import (
	"fmt"
	"log"
	"time"

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
	// Звпускаем краулер, который будет получать данные RSS-ленты периодически.
	t := time.NewTicker(30 * time.Second)
	defer t.Stop()
	for {
		select {
		//  time.Ticker.C возвращает канал, который будет отправлять событие каждый раз, когда таймер истекает.
		case <-t.C:
			// Запускаем горутины для каждой RSS-ленты.
			// Создаем WaitGroup, чтобы дождаться завершения всех горутин.
			log.Println("fetching feeds...")
			for _, feed := range c.feeds {
				// Здесь мы используем анонимную функцию, чтобы передать feed внутрь неё и напечатать его.
				go func(f string) {
					if rssData, err := c.fetchFeedData(f); err != nil {
						log.Printf("failed to fetch feed data: %v", err)
					} else {
						// Парсим данные RSS.
						feed, err := c.feedParser.ParseString(rssData)
						if err != nil {
							log.Printf("failed to parse feed data: %v", err)
						}
						// Печатаем заголовок.
						log.Printf("fetched feed: %s\n", feed.Title)
					}
				}(feed)
			}
		}
	}
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
