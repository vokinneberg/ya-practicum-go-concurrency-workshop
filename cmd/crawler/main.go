package main

import "main/internal/crawler"

var rssFeedList = []string{
	"https://www.reddit.com/r/golang/.rss",
	"https://www.reddit.com/r/programming/.rss",
	"https://www.reddit.com/r/golang/.rss",
	"https://www.reddit.com/r/golang/.rss",
	"https://www.reddit.com/r/golang/.rss",
}

func main() {
	// Create a new Crawler
	c := crawler.New()

	// Add RSS feeds to the crawler
	for _, rssFeed := range rssFeedList {
		c.AddFeed(rssFeed)
	}

	// Start the crawler
	c.Start()
}
