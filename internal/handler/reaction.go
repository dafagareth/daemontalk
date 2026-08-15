package handler

import (
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"daemontalk/internal/i18n"
	"daemontalk/web/templates"
)

var validEmojis = map[string]bool{
	"👍": true,
	"❤️": true,
	"🔥": true,
	"🤯": true,
	"💡": true,
}

func (h *Handler) PostReaction(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)
	slug := chi.URLParam(r, "slug")
	emoji := chi.URLParam(r, "emoji")

	if !validEmojis[emoji] {
		http.Error(w, "invalid emoji", http.StatusBadRequest)
		return
	}
	if h.Comments == nil {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
		return
	}

	// Prevent multiple reactions from the same user for this post.
	cookieName := "reacted_" + slug
	if cookie, err := r.Cookie(cookieName); err == nil && cookie.Value != "" {
		oldEmoji, _ := url.QueryUnescape(cookie.Value)

		if oldEmoji == emoji {
			// Undo reaction
			reactions, err := h.Comments.DecrementReaction(slug, emoji)
			if err != nil {
				log.Printf("decrement reaction %s/%s: %v", slug, emoji, err)
				reactions = map[string]int{}
			}
			http.SetCookie(w, &http.Cookie{
				Name:     cookieName,
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			})
			if err := templates.ReactionsBar(ui, reactions, slug, lang, "").Render(r.Context(), w); err != nil {
		log.Printf("render error: %v", err)
			}
			return
		}

		// Switch reaction
		if validEmojis[oldEmoji] {
			_, _ = h.Comments.DecrementReaction(slug, oldEmoji)
		}
		reactions, err := h.Comments.IncrementReaction(slug, emoji)
		if err != nil {
			log.Printf("increment reaction %s/%s: %v", slug, emoji, err)
			reactions = map[string]int{}
		}
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Value:    url.QueryEscape(emoji),
			Path:     "/",
			MaxAge:   86400 * 365,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
		if err := templates.ReactionsBar(ui, reactions, slug, lang, emoji).Render(r.Context(), w); err != nil {
		log.Printf("render error: %v", err)
		}
		return
	}

	reactions, err := h.Comments.IncrementReaction(slug, emoji)
	if err != nil {
		log.Printf("increment reaction %s/%s: %v", slug, emoji, err)
		reactions = map[string]int{}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    url.QueryEscape(emoji),
		Path:     "/",
		MaxAge:   86400 * 365,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	if err := templates.ReactionsBar(ui, reactions, slug, lang, emoji).Render(r.Context(), w); err != nil {
		log.Printf("render error: %v", err)
	}
}
