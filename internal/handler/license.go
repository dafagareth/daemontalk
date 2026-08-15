package handler

import (
	"net/http"
	"daemontalk/web/templates"
)

func (h *Handler) License(w http.ResponseWriter, r *http.Request) {
	h.renderMarkdownPage(w, r, "license", "license", templates.PageMeta{
		Description: "License",
	}, templates.LicensePage)
}