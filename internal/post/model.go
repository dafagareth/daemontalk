package post

import (
	"html/template"
	"time"
)

type TOCEntry struct {
	ID    string
	Title string
	Level int
}

type Post struct {
	Title          string
	Slug           string
	Aliases        []string
	Date           time.Time
	Tags           []string
	Lang           string
	Draft          bool
	Status         string
	PublishAt      time.Time
	Type           string
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
