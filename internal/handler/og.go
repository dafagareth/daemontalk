package handler

import (
	"bytes"
	"net/http"
	"strings"
	"sync"
	"time"

	"daemontalk/internal/og"
	"daemontalk/internal/post"
	"github.com/go-chi/chi/v5"
)

var (
	ogMu    sync.RWMutex
	ogCache = map[string][]byte{}
)

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
			Tags:     p.Tags,
			ReadTime: p.ReadTime,
			Date:     strings.ToUpper(p.Date.Format("02 Jan 2006")),
			Site:     "daemontalk.com",
			Cover:    p.Cover,
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

func (h *Handler) SiteOGImage(w http.ResponseWriter, r *http.Request) {
	const key = "__site_default__"

	ogMu.RLock()
	cached, hit := ogCache[key]
	ogMu.RUnlock()

	if !hit {
		var buf bytes.Buffer
		card := og.Card{
			Title:    "DaemonTalk — Engineering & Systems Exploration",
			Tags:     []string{"SYSTEMS", "SOFTWARE", "LINUX"},
			ReadTime: 5,
			Date:     strings.ToUpper(time.Now().Format("02 Jan 2006")),
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
