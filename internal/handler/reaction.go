package handler

import (
	"log/slog"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"daemontalk/internal/i18n"
	"daemontalk/web/templates"
)

var validEmojis = map[string]bool{
	"like":       true,
	"heart":      true,
	"fire":       true,
	"mindblown":  true,
	"insightful": true,
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
	cookieName := CookieReactedPrefix + slug
	if cookie, err := r.Cookie(cookieName); err == nil && cookie.Value != "" {
		oldEmoji, _ := url.QueryUnescape(cookie.Value)

		if oldEmoji == emoji {
			// Undo reaction
			reactions, err := h.Comments.DecrementReaction(slug, emoji)
			if err != nil {
				slog.Error("decrement reaction failed", "slug", slug, "emoji", emoji, "error", err)
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
			h.Render(w, r, templates.ReactionsBar(ui, reactions, slug, lang, ""))
			return
		}

		// Switch reaction
		if validEmojis[oldEmoji] {
			_, _ = h.Comments.DecrementReaction(slug, oldEmoji)
		}
		reactions, err := h.Comments.IncrementReaction(slug, emoji)
		if err != nil {
			slog.Error("increment reaction failed", "slug", slug, "emoji", emoji, "error", err)
			reactions = map[string]int{}
		}
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Value:    url.QueryEscape(emoji),
			Path:     "/",
			MaxAge:   CookieReactionMaxAge,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
		h.Render(w, r, templates.ReactionsBar(ui, reactions, slug, lang, emoji))
		return
	}

	reactions, err := h.Comments.IncrementReaction(slug, emoji)
	if err != nil {
		slog.Error("increment reaction failed", "slug", slug, "emoji", emoji, "error", err)
		reactions = map[string]int{}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    url.QueryEscape(emoji),
		Path:     "/",
		MaxAge:   CookieReactionMaxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	h.Render(w, r, templates.ReactionsBar(ui, reactions, slug, lang, emoji))
}

