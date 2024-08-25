package crawler

import (
	"fmt"
	"log"
	"time"

	"main/internal/feed"

	"github.com/go-resty/resty/v2"
	"github.com/mmcdole/gofeed"
)

type Crawler struct {
	httpClient *resty.Client
	feedParser *gofeed.Parser
	feeds      *feed.Storage
}

// NewCrawler creates a new Crawler
func New(feeds *feed.Storage, httpClient *resty.Client, feedParser *gofeed.Parser) *Crawler {
	return &Crawler{
		feeds:      feeds,
		httpClient: httpClient,
		feedParser: feedParser,
	}
}

// Start starts the crawler
func (c *Crawler) Start() {
	// Звпускаем краулер, который будет получать данные RSS-ленты периодически.
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()
	for {
		select {
		//  time.Ticker.C возвращает канал, который будет отправлять событие каждый раз, когда таймер истекает.
		case <-t.C:
			// Запускаем горутины для каждой RSS-ленты.
			// Создаем WaitGroup, чтобы дождаться завершения всех горутин.
			log.Println("fetching feeds...")
			// Получаем все RSS-ленты.
			links := c.feeds.GetLinks()
			for _, link := range links {
				// Здесь мы используем анонимную функцию, чтобы передать feed внутрь неё и напечатать его.
				go func(l string) {
					if rssData, err := c.fetchFeedData(l); err != nil {
						log.Printf("failed to fetch feed data: %v", err)
					} else {
						// Парсим данные RSS.
						feed, err := c.feedParser.ParseString(rssData)
						if err != nil {
							log.Printf("failed to parse feed data: %v", err)
						}
						// Печатаем заголовок.
						log.Printf("fetched feed: %s\n", feed.Title)
						c.feeds.SetFeed(l, mapFeed(feed))
					}
				}(link)
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

func mapFeed(fetchedData *gofeed.Feed) []feed.Item {
	if fetchedData == nil {
		return nil
	}

	mapItem := func(fetchedItem *gofeed.Item) *feed.Item {
		item := &feed.Item{
			Title:       fetchedItem.Title,
			Description: fetchedItem.Description,
			Link:        fetchedItem.Link,
			Published:   fetchedItem.Published,
			Source: &feed.Source{
				Title: fetchedData.Title,
				Link:  fetchedData.Link,
			},
		}

		if fetchedItem.Author != nil {
			item.Author = &feed.Author{
				Name: fetchedItem.Author.Name,
			}
		}

		if fetchedItem.Image != nil {
			item.Image = &feed.Image{
				Title: fetchedItem.Image.Title,
				Link:  fetchedItem.Image.URL,
			}
		}

		return item
	}

	items := make([]feed.Item, 0, len(fetchedData.Items))
	for _, fetchedItem := range fetchedData.Items {
		mappedItem := mapItem(fetchedItem)
		items = append(items, *mappedItem)
	}
	return items
}
