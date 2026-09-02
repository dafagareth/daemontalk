package handler

import (
	"net/http"

	"daemontalk/web/templates"
)

func (h *Handler) About(w http.ResponseWriter, r *http.Request) {
	h.renderMarkdownPage(w, r, "about", "about", templates.PageMeta{
		Description: "About Daemontalk philosophy, editorial standards, and systems research notebook.",
	}, templates.AboutPage)
}
