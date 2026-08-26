package handler

import (
	"net/http"

	"daemontalk/internal/i18n"
	"daemontalk/web/templates"
)

func (h *Handler) Saved(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)
	h.Render(w, r, templates.Layout(ui, lang, "saved", r.URL.Path, templates.PageMeta{
		Description: "Your saved posts reading list.",
	}, templates.SavedPage(ui, lang)))
}
