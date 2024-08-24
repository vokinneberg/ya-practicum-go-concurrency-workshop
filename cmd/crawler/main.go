package main

import (
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/mmcdole/gofeed"

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

	// Create a new gofeed parser
	fp := gofeed.NewParser()

	// Create a new Crawler
	c := crawler.New(httpClient, fp)

	// Add RSS feeds to the crawler
	for _, rssFeed := range rssFeedList {
		c.AddFeed(rssFeed)
	}

	// Start the crawler
	c.Start()
}
