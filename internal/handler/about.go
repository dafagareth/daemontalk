package handler

import (
	"log"
	"net/http"

	"daemontalk/internal/i18n"
	"daemontalk/web/templates"
)

func (h *Handler) About(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)

	title := "About Daemontalk"
	if lang == "id" {
		title = "Tentang Daemontalk"
	}

	err := templates.Layout(ui, lang, title, r.URL.Path, templates.PageMeta{
		Description: "About Daemontalk philosophy and editorial standards.",
	}, templates.AboutPage(ui, lang)).Render(r.Context(), w)
	if err != nil {
		log.Printf("render error: %v", err)
	}
}
