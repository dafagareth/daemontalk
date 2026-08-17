package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"daemontalk/internal/i18n"
	"daemontalk/web/templates"
)

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)
	isAdmin := h.isAdmin(r)

	query := strings.TrimSpace(r.URL.Query().Get("q"))

	var viewCounts map[string]int
	var results []templates.SearchResult

	if query != "" {
		q := strings.ToLower(query)
		for _, p := range h.VisiblePosts(isAdmin) {
			haystack := strings.ToLower(p.Title + " " + p.Description + " " + strings.Join(p.Tags, " "))
			if strings.Contains(haystack, q) {
				results = append(results, templates.SearchResult{Post: p})
			}
		}
		if h.Comments != nil {
			if vc, err := h.Comments.AllViewCounts(); err == nil {
				viewCounts = vc
			} else {
				slog.Error("search all view counts failed", "error", err)
			}
		}
	}

	desc := "Search posts on daemontalk.com"
	if query != "" {
		desc = fmt.Sprintf("Search results for \"%s\" on daemontalk.com", query)
	}
	meta := templates.PageMeta{Description: desc}

	h.Render(w, r, templates.Layout(ui, lang, "search", r.URL.Path, meta,
		templates.SearchPage(ui, query, results, lang, viewCounts)))
}

