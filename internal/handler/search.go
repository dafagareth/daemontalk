package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"daemontalk/internal/i18n"
	"daemontalk/web/templates"
)

var reHTMLTag = regexp.MustCompile(`<[^>]+>`)

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
			cleanBody := reHTMLTag.ReplaceAllString(string(p.Body), " ")
			haystack := strings.ToLower(p.Title + " " + p.Description + " " + strings.Join(p.Tags, " ") + " " + cleanBody)
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

	isHTMX := r.Header.Get("HX-Request") == "true"
	if isHTMX {
		h.Render(w, r, templates.SearchResultsList(ui, query, results, lang, viewCounts))
		return
	}

	desc := "Search posts on daemontalk.com"
	if query != "" {
		desc = fmt.Sprintf("Search results for \"%s\" on daemontalk.com", query)
	}
	meta := templates.PageMeta{Description: desc}

	h.Render(w, r, templates.Layout(ui, lang, "search", r.URL.Path, meta,
		templates.SearchPage(ui, query, results, lang, viewCounts)))
}

