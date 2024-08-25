package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"main/internal/crawler"
	"main/internal/feed"
	"main/internal/http/handler"

	"github.com/go-resty/resty/v2"
	"github.com/mmcdole/gofeed"
)

var rssFeedList = []string{
	"http://feeds.feedburner.com/TechCrunch",
	"https://www.wired.com/feed/rss",
	"http://feeds.arstechnica.com/arstechnica/index",
	"https://news.ycombinator.com/rss",
	"https://www.smashingmagazine.com/feed",
}

func main() {
	// Create feed storage
	fs := feed.NewFeedStorage(rssFeedList)

	// Create a new resty http client
	hc := resty.New()
	hc.SetTimeout(time.Duration(5 * time.Second))

	// Create a new gofeed parser
	fp := gofeed.NewParser()

	// Create a new Crawler
	c := crawler.New(fs, hc, fp)

	wg := sync.WaitGroup{}
	wg.Add(2)
	// Start the crawler
	go func() {
		defer wg.Done()
		c.Start()
	}()

	// Set up and start web server
	mux := http.NewServeMux()
	mux.HandleFunc("GET /rss", handler.RSSFeed(fs))
	go func() {
		defer wg.Done()
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Panic(err)
		}
	}()

	wg.Wait()
}
