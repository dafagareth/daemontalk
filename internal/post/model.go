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
	Title       string
	Slug        string
	Aliases     []string
	Date        time.Time
	Tags        []string
	Lang        string
	Draft       bool
	PublishAt   time.Time // future date = scheduled post
	Type        string    // "til" or "" for regular blog posts
	ReadTime    int
	Cover       string
	Body        template.HTML
	Description string
	TOC         []TOCEntry
	Series      string
	SeriesPart  int
	Author      string
}
