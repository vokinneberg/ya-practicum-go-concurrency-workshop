package feed

import (
	"fmt"
	"sort"
	"sync"
	"time"
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
	// Мьютекс позволяет только одной горутине получить доступ к данным RSS-лент.
	// RWLock позволяет нескольким горутинам получить доступ к данным RSS-лент на чтение. Но только одной горутине на запись.
	m sync.RWMutex
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
	fs.m.RLock()
	defer fs.m.RUnlock()
	links := make([]string, 0, len(fs.feeds))
	for link := range fs.feeds {
		links = append(links, link)
	}
	return links
}

// GetFeed возвращает данные RSS-ленты по URL.
func (fs *Storage) GetFeed(url string) ([]Item, error) {
	fs.m.RLock()
	defer fs.m.RUnlock()
	if _, ok := fs.feeds[url]; !ok {
		return nil, ErrFeedNotFound{URL: url}
	}
	return fs.feeds[url], nil
}

// GetFeeds возвращает данные всех RSS-лент.
func (fs *Storage) GetFeeds() []Item {
	fs.m.RLock()
	defer fs.m.RUnlock()
	result := make([]Item, 0, len(fs.feeds))
	for _, feed := range fs.feeds {
		result = append(result, feed...)
	}
	sort.Slice(result, func(i, j int) bool {
		timeI, err := time.Parse(time.RFC3339, result[i].Published)
		if err != nil {
			return false
		}
		timeJ, err := time.Parse(time.RFC3339, result[j].Published)
		if err != nil {
			return false
		}
		return timeI.After(timeJ)
	})
	return result
}
