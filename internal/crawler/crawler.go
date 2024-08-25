package crawler

import (
	"fmt"
	"log"
	"time"

	"main/internal/feed"

	"github.com/go-resty/resty/v2"
	"github.com/mmcdole/gofeed"
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
func (c *Crawler) Start(numWorkers int, numConsumers int) {
	jobs := make(chan string, numWorkers)
	results := make(chan *rss, numWorkers)
	done := make(chan struct{}, numConsumers)

	// Start worker pool
	for i := 0; i < numWorkers; i++ {
		go c.worker(jobs, results)
	}

	// Start consumer pool
	for i := 0; i < numConsumers; i++ {
		go c.consumer(results, done)
	}

	// Start job producer
	go c.producer(jobs)

	// Wait for all consumers to finish
	for i := 0; i < numConsumers; i++ {
		<-done
	}
}

func (c *Crawler) worker(jobs <-chan string, results chan<- *rss) {
	for link := range jobs {
		rssData, err := c.fetchFeedData(link)
		if err != nil {
			log.Printf("failed to fetch feed data: %v", err)
			continue
		}
		results <- &rss{
			link: link,
			data: rssData,
		}
	}
}

func (c *Crawler) producer(jobs chan<- string) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			links := c.feeds.GetLinks()
			for _, link := range links {
				jobs <- link
			}
		}
	}
}

func (c *Crawler) consumer(results <-chan *rss, done chan<- struct{}) {
	for rss := range results {
		if rss == nil {
			continue
		}

		feed, err := c.feedParser.ParseString(rss.data)
		if err != nil {
			log.Printf("failed to parse feed data: %v", err)
			continue
		}
		log.Printf("fetched feed: %s\n", feed.Title)
		c.feeds.SetFeed(rss.link, mapFeed(feed))
	}
	done <- struct{}{}
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
