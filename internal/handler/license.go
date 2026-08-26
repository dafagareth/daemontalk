package handler

import (
	"daemontalk/web/templates"
	"net/http"
)

func (h *Handler) License(w http.ResponseWriter, r *http.Request) {
	h.renderMarkdownPage(w, r, "license", "license", templates.PageMeta{
		Description: "License",
	}, templates.LicensePage)
}
