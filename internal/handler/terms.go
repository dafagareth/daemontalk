package handler

import (
	"net/http"

	"daemontalk/web/templates"
)

func (h *Handler) Terms(w http.ResponseWriter, r *http.Request) {
	h.renderMarkdownPage(w, r, "terms", "terms", templates.PageMeta{
		Description: "Terms of Service",
	}, templates.TermsPage)
}
