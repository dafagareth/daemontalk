package handler

import (
	"log"
	"net/http"

	"daemontalk/internal/i18n"
	"daemontalk/internal/links"
	"daemontalk/web/templates"
)

func (h *Handler) Links(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)

	grouped := links.ByCategory(links.All)

	err := templates.Layout(ui, lang, "Links · daemontalk", r.URL.Path, templates.PageMeta{
		Description: "A curated list of tools, articles, and resources I find useful.",
	}, templates.LinksPage(ui, grouped, links.CategoryOrder)).Render(r.Context(), w)
	if err != nil {
		log.Printf("render error: %v", err)
	}
}
