package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/mmcdole/gofeed"

	"main/internal/crawler"
	"main/internal/http/handler"
)

var rssFeedList = []string{
	"http://feeds.feedburner.com/TechCrunch",
	"https://www.wired.com/feed/rss",
	"http://feeds.arstechnica.com/arstechnica/index",
	"https://news.ycombinator.com/rss",
	"https://www.smashingmagazine.com/feed",
}

func main() {
	// Create a new resty http client
	httpClient := resty.New()
	httpClient.SetTimeout(time.Duration(5 * time.Second))

	// Create a new gofeed parser
	fp := gofeed.NewParser()

	// Create a new Crawler
	c := crawler.New(httpClient, fp)

	// Add RSS feeds to the crawler
	for _, rssFeed := range rssFeedList {
		c.AddFeed(rssFeed)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	// Start the crawler
	go func() {
		defer wg.Done()
		c.Start()
	}()

	// Set up and start web server
	mux := http.NewServeMux()
	mux.HandleFunc("GET /rss", handler.RSSFeed())
	go func() {
		defer wg.Done()
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Panic(err)
		}
	}()

	wg.Wait()
}
