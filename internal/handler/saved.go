package handler

import (
	"log"
	"net/http"

	"daemontalk/internal/i18n"
	"daemontalk/web/templates"
)

func (h *Handler) Saved(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)
	err := templates.Layout(ui, lang, "saved", r.URL.Path, templates.PageMeta{
		Description: "Your saved posts reading list.",
	}, templates.SavedPage(ui, lang)).Render(r.Context(), w)
	if err != nil {
		log.Printf("render error: %v", err)
	}
}
