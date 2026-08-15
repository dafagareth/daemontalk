package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (h *Handler) RSS(w http.ResponseWriter, r *http.Request) {
	scheme := "https"
	if r.TLS == nil && r.Header.Get("X-Forwarded-Proto") != "https" && r.Host == "localhost:8080" {
		scheme = "http"
	}
	base := fmt.Sprintf("%s://%s", scheme, r.Host)

	var items strings.Builder
	for _, p := range h.AllPosts() {
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
