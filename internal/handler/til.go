package handler

import (
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
	h.Render(w, r, templates.Layout(ui, lang, "til", r.URL.Path, meta, templates.TILPage(ui, posts, lang)))
}

