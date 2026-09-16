package distribution

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (h *Handler) Manifest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/manifest+json")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	fmt.Fprint(w, `{
  "id": "/",
  "name": "DaemonTalk",
  "short_name": "DaemonTalk",
  "description": "An independent tech publication and learning space exploring modern software, systems, and the digital world.",
  "start_url": "/",
  "scope": "/",
  "display": "standalone",
  "orientation": "any",
  "background_color": "#18181b",
  "theme_color": "#18181b",
  "categories": [
    "education",
    "technology",
    "news"
  ],
  "lang": "en",
  "icons": [
    {
      "src": "/static/logo/icon-48x48.png",
      "sizes": "48x48",
      "type": "image/png"
    },
    {
      "src": "/static/logo/icon-96x96.png",
      "sizes": "96x96",
      "type": "image/png"
    },
    {
      "src": "/static/logo/icon-192x192.png",
      "sizes": "192x192",
      "type": "image/png",
      "purpose": "any maskable"
    },
    {
      "src": "/static/images/icon-192.png",
      "sizes": "192x192",
      "type": "image/png",
      "purpose": "any"
    },
    {
      "src": "/static/images/icon-512.png",
      "sizes": "512x512",
      "type": "image/png",
      "purpose": "any maskable"
    }
  ],
  "shortcuts": [
    {
      "name": "Dispatches",
      "short_name": "Dispatches",
      "description": "Browse technical dispatches and deep dives",
      "url": "/blog",
      "icons": [{ "src": "/static/logo/icon-96x96.png", "sizes": "96x96" }]
    },
    {
      "name": "Socket Forum",
      "short_name": "Socket",
      "description": "Architectural discussions and open systems Q&A",
      "url": "/socket",
      "icons": [{ "src": "/static/logo/icon-96x96.png", "sizes": "96x96" }]
    },
    {
      "name": "System Stats",
      "short_name": "Stats",
      "description": "Real-time publication and engine metrics",
      "url": "/stats",
      "icons": [{ "src": "/static/logo/icon-96x96.png", "sizes": "96x96" }]
    }
  ]
}`)
}

func (h *Handler) Robots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	robotsTxt := `User-agent: *
Disallow: /admin
Disallow: /admin/
Allow: /

# Block AI Scrapers and Data Mining Bots
User-agent: GPTBot
User-agent: ChatGPT-User
User-agent: Google-Extended
User-agent: CCBot
User-agent: Anthropic-ai
User-agent: Claude-Web
User-agent: Omgili
User-agent: Omgilibot
User-agent: FacebookBot
User-agent: Diffbot
User-agent: Bytespider
User-agent: PerplexityBot
User-agent: cohere-ai
Disallow: /

Sitemap: %s/sitemap.xml
`
	fmt.Fprintf(w, robotsTxt, SeoBaseURL)
}

func (h *Handler) Sitemap(w http.ResponseWriter, r *http.Request) {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")

	writeURL := func(path, lastmod string) {
		b.WriteString("  <url>\n")
		b.WriteString("    <loc>" + SeoBaseURL + path + "</loc>\n")
		if lastmod != "" {
			b.WriteString("    <lastmod>" + lastmod + "</lastmod>\n")
		}
		b.WriteString("  </url>\n")
	}

	for _, p := range []string{"/", "/colophon", "/stats", "/socket"} {
		writeURL(p, "")
	}

	posts := h.getAllPosts()
	for _, post := range posts {
		if post.Draft {
			continue
		}
		if !post.PublishAt.IsZero() && post.PublishAt.After(time.Now()) {
			continue
		}
		lastmod := post.Date.Format("2006-01-02")
		writeURL("/blog/"+post.Slug, lastmod)
	}

	tagSeen := make(map[string]bool)
	for _, post := range posts {
		for _, t := range post.Tags {
			if !tagSeen[t] {
				tagSeen[t] = true
				writeURL("/blog/tag/"+t, "")
			}
		}
	}

	b.WriteString(`</urlset>`)

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write([]byte(b.String()))
}

func (h *Handler) SitemapIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/sitemap.xml", http.StatusMovedPermanently)
}

func (h *Handler) SitemapEN(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/sitemap.xml", http.StatusMovedPermanently)
}

func (h *Handler) SitemapID(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/sitemap.xml", http.StatusMovedPermanently)
}
