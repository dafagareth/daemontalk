package handler

import (
	"log"
	"net/http"

	"daemontalk/internal/i18n"
	"daemontalk/web/templates"
)

func (h *Handler) Resume(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)

	meta := templates.PageMeta{Description: "Dafa — Software Engineer. Resume and CV."}
	err := templates.Layout(ui, lang, "resume", r.URL.Path, meta,
		templates.ResumePage()).Render(r.Context(), w)
	if err != nil {
		log.Printf("render error: %v", err)
	}
}
