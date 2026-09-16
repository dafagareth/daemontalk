package distribution

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (h *Handler) RSS(w http.ResponseWriter, r *http.Request) {
	base := strings.TrimSuffix(h.absoluteURL(r, "/"), "/")

	var items strings.Builder
	for _, p := range h.getAllPosts() {
		if p.Draft {
			continue
		}
		if !p.PublishAt.IsZero() && p.PublishAt.After(time.Now()) {
			continue
		}
		url := fmt.Sprintf("%s/blog/%s", base, p.Slug)
		pubDate := p.Date.UTC().Format(time.RFC1123Z)
		items.WriteString(fmt.Sprintf(`
		<item>
			<title><![CDATA[%s]]></title>
			<link>%s</link>
			<guid>%s</guid>
			<pubDate>%s</pubDate>
			<description><![CDATA[%s]]></description>
		</item>`, p.Title, url, url, pubDate, string(p.Body)))
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
	<channel>
		<title>daemontalk</title>
		<link>%s</link>
		<description>Writing about Go, systems, and developer tools.</description>
		<language>en</language>
		<atom:link href="%s/rss.xml" rel="self" type="application/rss+xml"/>
		%s
	</channel>
</rss>`, base, base, items.String())
}

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

	posts := h.getAllPosts()
	items := make([]jsonFeedItem, 0, len(posts))
	for _, p := range posts {
		if p.Draft {
			continue
		}
		if !p.PublishAt.IsZero() && p.PublishAt.After(time.Now()) {
			continue
		}
		items = append(items, jsonFeedItem{
			ID:            SeoBaseURL + "/blog/" + p.Slug,
			URL:           SeoBaseURL + "/blog/" + p.Slug,
			Title:         p.Title,
			DatePublished: p.Date.Format("2006-01-02T00:00:00Z"),
			Summary:       p.Description,
			Tags:          p.Tags,
		})
	}

	feed := jsonFeed{
		Version:     "https://jsonfeed.org/version/1.1",
		Title:       "daemontalk",
		HomePageURL: SeoBaseURL,
		FeedURL:     SeoBaseURL + "/feed.json",
		Description: "Technology publication & community portal.",
		Items:       items,
	}

	w.Header().Set("Content-Type", "application/feed+json; charset=utf-8")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(feed)
}
