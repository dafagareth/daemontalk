package handler

import (
	"net/http"

	"daemontalk/internal/i18n"
	"daemontalk/internal/links"
	"daemontalk/web/templates"
)

func (h *Handler) Links(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)

	grouped := links.ByCategory(links.All)

	h.Render(w, r, templates.Layout(ui, lang, "Links · daemontalk", r.URL.Path, templates.PageMeta{
		Description: "A curated list of tools, articles, and resources I find useful.",
	}, templates.LinksPage(ui, grouped, links.CategoryOrder)))
}
