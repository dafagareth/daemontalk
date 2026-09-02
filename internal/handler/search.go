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

	visiblePosts := h.VisiblePosts(isAdmin)
	var results []templates.SearchResult

	if query != "" {
		q := strings.ToLower(query)
		for _, p := range visiblePosts {
			if strings.Contains(p.SearchHaystack, q) {
				results = append(results, templates.SearchResult{Post: p})
			}
		}
	}

	var viewCounts map[string]int
	if h.Comments != nil {
		if vc, err := h.Comments.AllViewCounts(); err == nil {
			viewCounts = vc
		} else {
			slog.Error("search all view counts failed", "error", err)
		}
	}

	isHTMX := r.Header.Get("HX-Request") == "true"
	if isHTMX {
		if r.URL.Query().Get("type") == "quick" || r.Header.Get("HX-Target") == "search-dropdown-results" {
			h.Render(w, r, templates.SearchDropdownResults(ui, query, results, lang, viewCounts))
			return
		}

		if r.Header.Get("HX-Target") == "search-results" {
			h.Render(w, r, templates.SearchResultsList(ui, query, results, lang, viewCounts))
			return
		}

		// If it's a boosted form submission or link, it will have a different target (or no explicit target).
		// We let it fall through to render the full page Layout, which hx-boost handles seamlessly.
	}

	desc := "Search posts on daemontalk.com"
	if query != "" {
		desc = fmt.Sprintf("Search results for \"%s\" on daemontalk.com", query)
	}
	meta := templates.PageMeta{Description: desc}

	h.Render(w, r, templates.Layout(ui, lang, "search", r.URL.Path, meta,
		templates.SearchPage(ui, query, results, lang, viewCounts)))
}
