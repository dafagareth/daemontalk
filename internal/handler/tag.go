package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"daemontalk/internal/i18n"
	"daemontalk/internal/post"
	"daemontalk/web/templates"
)

func (h *Handler) TagIndex(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)
	tag := chi.URLParam(r, "tag")
	isAdmin := h.isAdmin(r)

	var filtered []post.Post
	for _, p := range h.VisiblePosts(isAdmin) {
		for _, t := range p.Tags {
			if strings.EqualFold(t, tag) {
				filtered = append(filtered, p)
				break
			}
		}
	}

	if len(filtered) == 0 {
		h.NotFound(w, r)
		return
	}

	var viewCounts map[string]int
	if h.Comments != nil {
		if vc, err := h.Comments.AllViewCounts(); err != nil {
			log.Printf("tag view counts: %v", err)
		} else {
			viewCounts = vc
		}
	}

	if err := templates.Layout(ui, lang, "#"+tag, r.URL.Path, templates.PageMeta{
		Description: fmt.Sprintf("Posts tagged #%s on daemontalk.com", tag),
	}, templates.TagPage(ui, tag, filtered, lang, viewCounts)).Render(r.Context(), w); err != nil {
		log.Printf("render error: %v", err)
	}
}
