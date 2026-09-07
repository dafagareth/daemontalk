package handler

import (
	"net/http"

	"daemontalk/internal/graph"
	"daemontalk/internal/i18n"
	"daemontalk/web/templates"
)

func (h *Handler) Graph(w http.ResponseWriter, r *http.Request) {
	if !h.IsRadarEnabled() {
		h.NotFound(w, r)
		return
	}

	lang := langFromRequest(r)
	ui := i18n.Get(lang)
	posts := h.VisiblePosts(false)

	graphData := graph.Build(posts)

	title := "Knowledge Graph · daemontalk"
	if lang == "id" {
		title = "Peta Konsep & Tech Stack · daemontalk"
	}

	h.Render(w, r, templates.Layout(ui, lang, title, r.URL.Path, templates.PageMeta{
		Description: "Interactive systems knowledge graph connecting Linux kernel architectures, language runtimes, memory models, and distributed storage engines.",
	}, templates.GraphPage(ui, lang, posts, graphData.ToJSON(), graphData.Stats)))
}

func (h *Handler) GraphDataAPI(w http.ResponseWriter, r *http.Request) {
	if !h.IsRadarEnabled() {
		h.NotFound(w, r)
		return
	}
	posts := h.VisiblePosts(false)
	graphData := graph.Build(posts)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(graphData.ToJSON()))
}
