package pages

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"daemontalk/internal/handler/common"
	"daemontalk/internal/i18n"
	"daemontalk/web/templates/layouts"
	pagestmpl "daemontalk/web/templates/pages"
	"daemontalk/web/templates/shared"
)

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	lang := common.LangFromRequest(r)
	ui := i18n.Get(lang)
	isAdmin := h.IsAdmin(r)

	query := strings.TrimSpace(r.URL.Query().Get("q"))

	visiblePosts := h.GetVisiblePosts(isAdmin)
	var results []pagestmpl.SearchResult

	if query != "" {
		q := strings.ToLower(query)
		for _, p := range visiblePosts {
			if strings.Contains(p.SearchHaystack, q) {
				results = append(results, pagestmpl.SearchResult{Post: p})
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
			common.Render(w, r, pagestmpl.SearchDropdownResults(ui, query, results, lang, viewCounts))
			return
		}

		if r.Header.Get("HX-Target") == "search-results" {
			common.Render(w, r, pagestmpl.SearchResultsList(ui, query, results, lang, viewCounts))
			return
		}
	}

	desc := "Search posts on daemontalk.com"
	if query != "" {
		desc = fmt.Sprintf("Search results for \"%s\" on daemontalk.com", query)
	}
	meta := shared.PageMeta{Description: desc}

	common.Render(w, r, layouts.Layout(ui, lang, "search", r.URL.Path, meta,
		pagestmpl.SearchPage(ui, query, results, lang, viewCounts),
	))
}
