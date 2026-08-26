package handler

import (
	"net/http"

	"daemontalk/internal/i18n"
	"daemontalk/web/templates"
)

// Graph renders the interactive systems knowledge graph page
func (h *Handler) Graph(w http.ResponseWriter, r *http.Request) {
	if !h.IsRadarEnabled() {
		h.NotFound(w, r)
		return
	}

	lang := langFromRequest(r)
	ui := i18n.Get(lang)
	posts := h.VisiblePosts(false)

	title := "Knowledge Graph · daemontalk"
	if lang == "id" {
		title = "Peta Konsep & Tech Stack · daemontalk"
	}

	h.Render(w, r, templates.Layout(ui, lang, title, r.URL.Path, templates.PageMeta{
		Description: "Interactive knowledge graph connecting Linux kernel architectures, language runtimes, memory models, and distributed storage engines.",
	}, templates.GraphPage(ui, lang, posts)))
}
