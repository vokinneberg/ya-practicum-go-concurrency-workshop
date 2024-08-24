package main

import (
	"time"

	"github.com/go-resty/resty/v2"

	"main/internal/crawler"
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

	// Create a new Crawler
	c := crawler.New(httpClient)

	// Add RSS feeds to the crawler
	for _, rssFeed := range rssFeedList {
		c.AddFeed(rssFeed)
	}

	// Start the crawler
	c.Start()
}
