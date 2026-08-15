package handler

import (
	"log"
	"net/http"

	"daemontalk/internal/i18n"
	"daemontalk/internal/post"
	"daemontalk/web/templates"
)

func (h *Handler) Changelog(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)

	filename := h.getContentPath("changelog.md")
	if lang == "id" {
		filename = h.getContentPath("changelog.id.md")
	}
	body, _ := post.LoadBody(filename)

	err := templates.Layout(ui, lang, "Changelog · daemontalk", r.URL.Path, templates.PageMeta{
		Description: "A running log of features and changes shipped to daemontalk.com.",
	}, templates.ChangelogPage(body, lang)).Render(r.Context(), w)
	if err != nil {
		log.Printf("render error: %v", err)
	}
}
