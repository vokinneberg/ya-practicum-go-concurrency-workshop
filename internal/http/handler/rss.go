package handler

import "net/http"

func RSSFeed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Здесь мы будем возвращать RSS-ленту.
	}
}
