package handler

import (
	"encoding/json"

	"daemontalk/internal/post"
	"daemontalk/web/templates"
)

const (
	siteName = "daemontalk"
)

// articleJSONLD builds a schema.org Article block for a blog post.
func articleJSONLD(p post.Post, imageURL string) string {
	author := p.Author
	if author == "" {
		author = siteName
	}
	data := map[string]any{
		"@context":      "https://schema.org",
		"@type":         "BlogPosting",
		"headline":      p.Title,
		"url":           templates.AbsoluteURL("/blog/" + p.Slug),
		"datePublished": p.Date.Format("2006-01-02"),
		"author": map[string]any{
			"@type": "Person",
			"name":  author,
		},
		"publisher": map[string]any{
			"@type": "Organization",
			"name":  siteName,
			"url":   templates.AbsoluteURL("/"),
		},
	}
	if p.Description != "" {
		data["description"] = p.Description
	}
	if imageURL != "" {
		data["image"] = imageURL
	}
	if len(p.Tags) > 0 {
		data["keywords"] = joinTags(p.Tags)
	}
	return marshalJSONLD(data)
}

// siteJSONLD builds a WebSite block for the home page.
func siteJSONLD() string {
	data := map[string]any{
		"@context": "https://schema.org",
		"@type":    "WebSite",
		"name":     siteName,
		"url":      templates.AbsoluteURL("/"),
		"publisher": map[string]any{
			"@type": "Organization",
			"name":  siteName,
			"url":   templates.AbsoluteURL("/"),
		},
	}
	return marshalJSONLD(data)
}

func marshalJSONLD(data map[string]any) string {
	b, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(b)
}

func joinTags(tags []string) string {
	out := ""
	for i, t := range tags {
		if i > 0 {
			out += ", "
		}
		out += t
	}
	return out
}
