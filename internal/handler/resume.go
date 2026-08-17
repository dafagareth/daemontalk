package handler

import (
	"net/http"

	"daemontalk/internal/i18n"
	"daemontalk/web/templates"
)

func (h *Handler) Resume(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)

	meta := templates.PageMeta{Description: "Dafa — Software Engineer. Resume and CV."}
	h.Render(w, r, templates.Layout(ui, lang, "resume", r.URL.Path, meta, templates.ResumePage()))
}

