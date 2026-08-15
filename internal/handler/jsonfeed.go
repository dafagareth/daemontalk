package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) JSONFeed(w http.ResponseWriter, r *http.Request) {
	type jsonFeedItem struct {
		ID            string   `json:"id"`
		URL           string   `json:"url"`
		Title         string   `json:"title"`
		DatePublished string   `json:"date_published"`
		Summary       string   `json:"summary,omitempty"`
		Tags          []string `json:"tags,omitempty"`
	}

	type jsonFeed struct {
		Version     string         `json:"version"`
		Title       string         `json:"title"`
		HomePageURL string         `json:"home_page_url"`
		FeedURL     string         `json:"feed_url"`
		Description string         `json:"description"`
		Items       []jsonFeedItem `json:"items"`
	}

	items := make([]jsonFeedItem, 0, len(h.AllPosts()))
	for _, p := range h.AllPosts() {
		if p.Draft {
			continue
		}
		items = append(items, jsonFeedItem{
			ID:            seoBaseURL + "/blog/" + p.Slug,
			URL:           seoBaseURL + "/blog/" + p.Slug,
			Title:         p.Title,
			DatePublished: p.Date.Format("2006-01-02T00:00:00Z"),
			Summary:       p.Description,
			Tags:          p.Tags,
		})
	}

	feed := jsonFeed{
		Version:     "https://jsonfeed.org/version/1.1",
		Title:       "daemontalk",
		HomePageURL: seoBaseURL,
		FeedURL:     seoBaseURL + "/feed.json",
		Description: "Software developer building developer tools in Go.",
		Items:       items,
	}

	w.Header().Set("Content-Type", "application/feed+json; charset=utf-8")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(feed)
}
