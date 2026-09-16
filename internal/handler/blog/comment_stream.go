package blog

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

var (
	commentSubscribers   = make(map[chan string]string)
	commentSubscribersMu sync.RWMutex
)

func broadcastNewComment(slug string) {
	commentSubscribersMu.RLock()
	defer commentSubscribersMu.RUnlock()
	for ch, s := range commentSubscribers {
		if s == slug {
			select {
			case ch <- slug:
			default:
			}
		}
	}
}

func (h *Handler) StreamComments(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "text/html") && !strings.Contains(accept, "text/event-stream") {
		http.Redirect(w, r, "/blog/"+slug+"#comments", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache, no-transform")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, ": connected\n\n")
	flusher.Flush()

	ch := make(chan string, 1)
	commentSubscribersMu.Lock()
	commentSubscribers[ch] = slug
	commentSubscribersMu.Unlock()

	defer func() {
		commentSubscribersMu.Lock()
		delete(commentSubscribers, ch)
		commentSubscribersMu.Unlock()
	}()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	notify := r.Context().Done()
	for {
		select {
		case <-notify:
			return
		case <-ticker.C:
			fmt.Fprintf(w, ": ping\n\n")
			flusher.Flush()
		case s := <-ch:
			fmt.Fprintf(w, "event: new_comment\ndata: %s\n\n", s)
			flusher.Flush()
		}
	}
}
