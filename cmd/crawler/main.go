package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"main/internal/crawler"
	"main/internal/feed"
	"main/internal/http/handler"

	"github.com/go-resty/resty/v2"
	"github.com/mmcdole/gofeed"
	"golang.org/x/sync/errgroup"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	// Start the crawler
	g.Go(func() error {
		if err := c.Start(ctx, 5, 3); err != nil {
			return fmt.Errorf("failed to start crawler: %w", err)
		}
		return nil
	})

	// Set up and start web server
	mux := http.NewServeMux()
	mux.HandleFunc("GET /rss", handler.RSSFeed(fs))

	// Set up and start web server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	g.Go(func() error {
		if err := srv.ListenAndServe(); err != nil {
			return fmt.Errorf("failed to start web server: %w", err)
		}
		return nil
	})

	// Handle shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	g.Go(func() error {
		select {
		case <-sigChan:
			log.Println("Received interrupt signal. Shutting down...")
			cancel()
		case <-ctx.Done():
		}

		// Shutdown the HTTP server
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("failed to shutdown web server: %w", err)
		}
		return nil
	})

	// Wait for all goroutines to complete
	if err := g.Wait(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	log.Println("Graceful shutdown completed")
}
