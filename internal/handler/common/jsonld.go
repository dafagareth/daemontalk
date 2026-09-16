package common

import (
	"encoding/json"

	"daemontalk/internal/post"
	"daemontalk/web/templates/shared"
)

const (
	SiteName = "daemontalk"
)

func ArticleJSONLD(p post.Post, imageURL string) string {
	author := p.Author
	if author == "" {
		author = SiteName
	}
	data := map[string]any{
		"@context":      "https://schema.org",
		"@type":         "BlogPosting",
		"headline":      p.Title,
		"url":           shared.AbsoluteURL("/blog/" + p.Slug),
		"datePublished": p.Date.Format("2006-01-02"),
		"author": map[string]any{
			"@type": "Person",
			"name":  author,
		},
		"publisher": map[string]any{
			"@type": "Organization",
			"name":  SiteName,
			"url":   shared.AbsoluteURL("/"),
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
	return MarshalJSONLD(data)
}

func SiteJSONLD() string {
	data := map[string]any{
		"@context": "https://schema.org",
		"@type":    "WebSite",
		"name":     SiteName,
		"url":      shared.AbsoluteURL("/"),
		"publisher": map[string]any{
			"@type": "Organization",
			"name":  SiteName,
			"url":   shared.AbsoluteURL("/"),
		},
	}
	return MarshalJSONLD(data)
}

func MarshalJSONLD(data map[string]any) string {
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
