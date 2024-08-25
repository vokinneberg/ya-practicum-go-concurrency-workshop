package feed

type Image struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

type Author struct {
	Name string `json:"name"`
}

type Source struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

type Item struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Link        string  `json:"link"`
	Image       *Image  `json:"image,omitempty"`
	Author      *Author `json:"author,omitempty"`
	Published   string  `json:"published"`
	Source      *Source `json:"source,omitempty"`
	SourceLink  string  `json:"sourceLink"`
}
