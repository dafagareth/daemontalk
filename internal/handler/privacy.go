package handler

import (
	"net/http"

	"daemontalk/web/templates"
)

func (h *Handler) Privacy(w http.ResponseWriter, r *http.Request) {
	h.renderMarkdownPage(w, r, "privacy", "privacy", templates.PageMeta{
		Description: "Privacy Policy",
	}, templates.PrivacyPage)
}
