package post

import (
	"html/template"
	"time"
)

type TOCEntry struct {
	ID    string
	Title string
	Level int // 2 = h2, 3 = h3
}

type Post struct {
	Title          string
	Slug           string
	Aliases        []string
	Date           time.Time
	Tags           []string
	Lang           string
	Draft          bool
	Status         string    // e.g., "published", "draft", "archived"
	PublishAt      time.Time // future date = scheduled post
	Type           string    // Post type (e.g., standard)
	ReadTime       int
	Cover          string
	CoverCaption   string
	CoverSource    string
	Body           template.HTML
	Description    string
	TOC            []TOCEntry
	Series         string
	SeriesPart     int
	Author         string
	AuthorAvatar   string
	AuthorGitHub   string
	Contributors   []string
	SearchHaystack string
}
