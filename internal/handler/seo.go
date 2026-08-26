package handler

import (
	"fmt"
	"net/http"
	"strings"
)

// Manifest serves the PWA web app manifest.
func (h *Handler) Manifest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/manifest+json")
	fmt.Fprint(w, `{
  "name": "daemontalk",
  "short_name": "daemontalk",
  "description": "Notes on Linux, Go, and things I learn along the way.",
  "start_url": "/",
  "display": "standalone",
  "background_color": "#f7f7f5",
  "theme_color": "#1a73e8",
  "icons": [
    {"src": "/static/images/icon-192.png", "sizes": "192x192", "type": "image/png"},
    {"src": "/static/images/icon-512.png", "sizes": "512x512", "type": "image/png"}
  ]
}`)
}

const seoBaseURL = "https://daemontalk.com"

// Robots serves robots.txt pointing crawlers at the sitemap and blocking AI scrapers.
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
Sitemap: %s/sitemap-index.xml
`
	fmt.Fprintf(w, robotsTxt, seoBaseURL, seoBaseURL)
}

// Sitemap serves an XML sitemap covering static pages and every blog post,
// in both the default and Indonesian (/id) variants.
func (h *Handler) Sitemap(w http.ResponseWriter, r *http.Request) {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")

	writeURL := func(path, lastmod string) {
		b.WriteString("  <url>\n")
		b.WriteString("    <loc>" + seoBaseURL + path + "</loc>\n")
		if lastmod != "" {
			b.WriteString("    <lastmod>" + lastmod + "</lastmod>\n")
		}
		b.WriteString("  </url>\n")
	}

	// Static pages (en + id).
	for _, p := range []string{"/", "/behind", "/stats", "/guestbook", "/resume"} {
		writeURL(p, "")
		writeURL("/id"+strings.TrimSuffix(p, "/"), "")
	}

	// Blog posts (en + id).
	for _, post := range h.AllPosts() {
		lastmod := post.Date.Format("2006-01-02")
		writeURL("/blog/"+post.Slug, lastmod)
		writeURL("/id/blog/"+post.Slug, lastmod)
	}

	// Tag pages: collect unique tags across all posts.
	tagSeen := make(map[string]bool)
	for _, post := range h.AllPosts() {
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

// SitemapIndex serves a sitemap index referencing per-language sitemaps.
func (h *Handler) SitemapIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <sitemap><loc>%s/sitemap-en.xml</loc></sitemap>
  <sitemap><loc>%s/sitemap-id.xml</loc></sitemap>
</sitemapindex>`, seoBaseURL, seoBaseURL)
}

// SitemapEN serves an English-only sitemap.
func (h *Handler) SitemapEN(w http.ResponseWriter, r *http.Request) {
	h.langSitemap(w, "en")
}

// SitemapID serves an Indonesian-only sitemap.
func (h *Handler) SitemapID(w http.ResponseWriter, r *http.Request) {
	h.langSitemap(w, "id")
}

func (h *Handler) langSitemap(w http.ResponseWriter, lang string) {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")

	writeURL := func(path, lastmod string) {
		b.WriteString("  <url>\n")
		b.WriteString("    <loc>" + seoBaseURL + path + "</loc>\n")
		if lastmod != "" {
			b.WriteString("    <lastmod>" + lastmod + "</lastmod>\n")
		}
		b.WriteString("  </url>\n")
	}

	pfx := ""
	if lang == "id" {
		pfx = "/id"
	}

	// Static pages
	for _, p := range []string{"/", "/behind", "/stats", "/guestbook", "/resume", "/changelog", "/links"} {
		path := pfx + p
		if lang != "id" && p == "/" {
			path = "/"
		}
		if lang == "id" && p == "/" {
			path = "/id"
		}
		writeURL(path, "")
	}

	// Blog posts
	for _, post := range h.AllPosts() {
		if post.Draft {
			continue
		}
		lastmod := post.Date.Format("2006-01-02")
		writeURL(pfx+"/blog/"+post.Slug, lastmod)
	}

	b.WriteString(`</urlset>`)
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write([]byte(b.String()))
}
