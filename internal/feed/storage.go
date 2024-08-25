package feed

import (
	"fmt"
	"sync"
)

// ErrFeedNotFound возвращается, когда RSS-лента не найдена.
type ErrFeedNotFound struct {
	URL     string
	Message string
}

func (e ErrFeedNotFound) Error() string {
	return fmt.Sprintf("feed not found: %s", e.URL)
}

// FeedStorage хранит данные RSS-лент.
type Storage struct {
	// Здесь мы будем хранить данные RSS-ленты.
	feeds map[string][]Item

	// Мьютекс для безопасного доступа к данным RSS-лент.
	m sync.Mutex
}

// NewFeedStorage создает новое хранилище RSS-лент.
func NewFeedStorage(feeds []string) *Storage {
	f := make(map[string][]Item, len(feeds))
	for _, feed := range feeds {
		// Заполняем feeds данными RSS-лент.
		f[feed] = []Item{}
	}
	return &Storage{
		feeds: f,
	}
}

// SetFeed устанавливает данные RSS-ленты по URL.
func (fs *Storage) SetFeed(url string, feed []Item) {
	fs.m.Lock()
	defer fs.m.Unlock()
	fs.feeds[url] = feed
}

// GetFeedLinks возвращает список URL RSS-лент.
func (fs *Storage) GetLinks() []string {
	fs.m.Lock()
	defer fs.m.Unlock()
	links := make([]string, 0, len(fs.feeds))
	for link := range fs.feeds {
		links = append(links, link)
	}
	return links
}

// GetFeed возвращает данные RSS-ленты по URL.
func (fs *Storage) GetFeed(url string) ([]Item, error) {
	fs.m.Lock()
	defer fs.m.Unlock()
	if _, ok := fs.feeds[url]; !ok {
		return nil, ErrFeedNotFound{URL: url}
	}
	return fs.feeds[url], nil
}

// GetFeeds возвращает данные всех RSS-лент.
func (fs *Storage) GetFeeds() []Item {
	fs.m.Lock()
	defer fs.m.Unlock()
	result := make([]Item, 0, len(fs.feeds))
	for _, feed := range fs.feeds {
		if feed != nil {
			result = append(result, feed...)
		}
	}
	return result
}
