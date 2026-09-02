package handler

import (
	"net/http"

	"daemontalk/internal/i18n"
	"daemontalk/web/templates"
)

// NotFound renders the custom 404 Not Found error page.
func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)
	w.WriteHeader(http.StatusNotFound)
	_ = templates.Layout(ui, lang, "404", r.URL.Path, templates.PageMeta{
		Description: ui.NotFound_Body,
	}, templates.NotFound(ui, lang)).Render(r.Context(), w)
}
