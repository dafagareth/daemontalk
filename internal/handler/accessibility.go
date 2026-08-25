package handler

import (
	"net/http"

	"daemontalk/web/templates"
)

func (h *Handler) Accessibility(w http.ResponseWriter, r *http.Request) {
	h.renderMarkdownPage(w, r, "accessibility", "accessibility", templates.PageMeta{
		Description: "Accessibility Statement",
	}, templates.AccessibilityPage)
}
