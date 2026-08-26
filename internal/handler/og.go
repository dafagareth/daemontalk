package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"daemontalk/internal/og"
	"daemontalk/internal/post"
	"github.com/go-chi/chi/v5"
)

// ogCache memoizes rendered OG PNGs by slug. Cards only change when post
// metadata changes (i.e. a redeploy), so an unbounded in-memory cache is fine.
var (
	ogMu    sync.RWMutex
	ogCache = map[string][]byte{}
)

// OGImage renders (and caches) a 1200×630 share image for a post.
func (h *Handler) OGImage(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	p, ok := post.FindBySlug(h.AllPosts(), slug)
	if !ok {
		http.NotFound(w, r)
		return
	}

	ogMu.RLock()
	cached, hit := ogCache[slug]
	ogMu.RUnlock()

	if !hit {
		var buf bytes.Buffer
		card := og.Card{
			Title:    p.Title,
			Subtitle: ogSubtitle(p),
			Site:     "daemontalk.com",
		}
		if err := og.Render(&buf, card); err != nil {
			http.Error(w, "render error", http.StatusInternalServerError)
			return
		}
		cached = buf.Bytes()
		ogMu.Lock()
		ogCache[slug] = cached
		ogMu.Unlock()
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	_, _ = w.Write(cached)
}

// SiteOGImage renders the default share card used for non-post pages.
func (h *Handler) SiteOGImage(w http.ResponseWriter, r *http.Request) {
	const key = "__site_default__"

	ogMu.RLock()
	cached, hit := ogCache[key]
	ogMu.RUnlock()

	if !hit {
		var buf bytes.Buffer
		card := og.Card{
			Title:    "Dafa — developer tools in Go",
			Subtitle: "Go · CLI · systems",
			Site:     "daemontalk.com",
		}
		if err := og.Render(&buf, card); err != nil {
			http.Error(w, "render error", http.StatusInternalServerError)
			return
		}
		cached = buf.Bytes()
		ogMu.Lock()
		ogCache[key] = cached
		ogMu.Unlock()
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	_, _ = w.Write(cached)
}

func ogSubtitle(p post.Post) string {
	parts := []string{fmt.Sprintf("%d min read", p.ReadTime)}
	if len(p.Tags) > 0 {
		parts = append(parts, strings.Join(p.Tags, ", "))
	}
	return strings.Join(parts, " · ")
}
