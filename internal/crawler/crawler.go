package crawler

import (
	"context"
	"fmt"
	"log"
	"time"

	"main/internal/feed"

	"github.com/go-resty/resty/v2"
	"github.com/mmcdole/gofeed"
	"golang.org/x/sync/errgroup"
)

type rss struct {
	link string
	data string
}

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

// Запускает краулер
func (c *Crawler) Start(ctx context.Context, numWorkers int, numConsumers int) error {
	// Создаем каналы для передачи данных между горутинами
	// jobs - канал для передачи URL RSS-лент между горутинами
	// results - канал для передачи данных RSS-лент между горутинами
	jobs := make(chan string, numWorkers)
	results := make(chan *rss, numWorkers)

	g := new(errgroup.Group)
	// Запускаем пул воркеров
	for i := 0; i < numWorkers; i++ {
		g.Go(func() error {
			if err := c.worker(ctx, jobs, results); err != nil {
				return fmt.Errorf("failed to run worker: %w", err)
			}
			return nil
		})
	}

	// Запускаем пул консьюмеров
	for i := 0; i < numConsumers; i++ {
		g.Go(func() error {
			if err := c.consumer(ctx, results); err != nil {
				return fmt.Errorf("failed to run consumer: %w", err)
			}
			return nil
		})
	}

	// Запускаем продюсера
	g.Go(func() error {
		if err := c.producer(ctx, jobs); err != nil {
			return fmt.Errorf("failed to run producer: %w", err)
		}
		return nil
	})

	// Ожидаем завершения всех горутин и возвращаем ошибку, если таковая возникла
	return g.Wait()
}

func (c *Crawler) worker(ctx context.Context, jobs <-chan string, results chan<- *rss) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled %w", ctx.Err())
		case link, ok := <-jobs:
			if !ok {
				return nil
			}
			rssData, err := c.fetchFeedData(link)
			if err != nil {
				return fmt.Errorf("failed to fetch feed data: %w", err)
			}
			results <- &rss{
				link: link,
				data: rssData,
			}
		}
	}
}

func (c *Crawler) producer(ctx context.Context, jobs chan<- string) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled %w", ctx.Err())
		case <-ticker.C:
			links := c.feeds.GetLinks()
			for _, link := range links {
				jobs <- link
			}
		}
	}
}

func (c *Crawler) consumer(ctx context.Context, results <-chan *rss) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled %w", ctx.Err())
		case rss, ok := <-results:
			if !ok {
				return nil
			}
			if rss == nil {
				log.Printf("received nil rss data")
				continue
			}
			feed, err := c.feedParser.ParseString(rss.data)
			if err != nil {
				log.Printf("failed to parse feed data: %v", err)
				return fmt.Errorf("failed to parse feed data: %w", err)
			}
			log.Printf("fetched feed: %s\n", feed.Title)
			c.feeds.SetFeed(rss.link, mapFeed(feed))
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
