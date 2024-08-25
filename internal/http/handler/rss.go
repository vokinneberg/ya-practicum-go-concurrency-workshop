package handler

import (
	"main/internal/feed"
	"net/http"
	"text/template"
)

// HTML-шаблон для RSS-ленты.
const rssFeedTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>RSS Feed Items</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; padding: 20px; }
        .main-description { color: #666; font-size: 0.9em; margin-bottom: 20px; }
        .item { margin-bottom: 20px; display: flex; }
        .item-content { flex: 1; }
        .item h3 { margin-bottom: 5px; }
        .item-image { width: 100px; height: 100px; object-fit: cover; margin-right: 15px; }
        .tease { color: #666; }
        .meta { font-size: 0.8em; color: #888; margin-top: 5px; }
        .source { font-style: italic; color: #444; }
        .no-data { text-align: center; color: #666; font-size: 1.2em; margin-top: 50px; }
    </style>
</head>
<body>
    <p class="main-description">Latest updates from RSS feeds</p>
    {{if .Items}}
        {{range .Items}}
        <div class="item">
            {{if .Image}}
            <img src="{{.Image.Link}}" alt="{{.Image.Title}}" class="item-image">
            {{end}}
            <div class="item-content">
                <h3><a href="{{.Link}}">{{.Title}}</a></h3>
                <p class="tease">{{if gt (len .Description) 400}}{{slice .Description 0 397}}...{{else}}{{.Description}}{{end}}</p>
                <p class="meta">
                    {{if .Author}}By {{.Author.Name}} · {{end}}
                    Published on {{.Published}} ·
                    <span class="source">Source: 
                        {{if .Source}}
                            <a href="{{.Source.Link}}">{{.Source.Title}}</a>
                        {{else}}
                            <a href="{{.SourceLink}}">{{.SourceLink}}</a>
                        {{end}}
                    </span>
                </p>
            </div>
        </div>
        {{end}}
    {{else}}
        <p class="no-data">No feed data available yet. Please check back later.</p>
    {{end}}
</body>
</html>
`

var rssFeedTmpl = template.Must(template.New("rssFeed").Parse(rssFeedTemplate))

func RSSFeed(feeds *feed.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Здесь мы будем возвращать RSS-ленту.
		// Получаем все RSS-ленты.
		feedsItems := feeds.GetFeeds()
		// Отправляем данные в шаблон.
		if err := rssFeedTmpl.Execute(w, struct {
			Items []feed.Item
		}{
			Items: feedsItems,
		}); err != nil {
			http.Error(w, "failed to render RSS feed", http.StatusInternalServerError)
		}
	}
}
