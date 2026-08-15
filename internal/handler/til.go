package handler

import (
	"log"
	"net/http"

	"daemontalk/internal/i18n"
	"daemontalk/internal/post"
	"daemontalk/web/templates"
)

func (h *Handler) TIL(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)
	isAdmin := h.isAdmin(r)

	var posts []post.Post
	for _, p := range h.VisiblePosts(isAdmin) {
		if p.Type != "til" {
			continue
		}
		posts = append(posts, p)
	}

	meta := templates.PageMeta{Description: "Short notes on things I've learned — TIL posts by Dafa."}
	if err := templates.Layout(ui, lang, "til", r.URL.Path, meta,
		templates.TILPage(ui, posts, lang)).Render(r.Context(), w); err != nil {
		log.Printf("render error: %v", err)
	}
}
