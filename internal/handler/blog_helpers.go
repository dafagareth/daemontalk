package handler

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"daemontalk/internal/i18n"
	"daemontalk/internal/post"
	"daemontalk/web/templates"
)

func (h *Handler) BlogPostsPartial(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("lang")
	if lang != "id" {
		lang = "en"
	}
	ui := i18n.Get(lang)

	visible := h.VisiblePosts(h.isAdmin(r))

	tagFilter := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("tag")))
	filtered := visible
	if tagFilter != "" {
		filtered = make([]post.Post, 0, len(visible))
		for _, p := range visible {
			for _, t := range p.Tags {
				if strings.ToLower(t) == tagFilter {
					filtered = append(filtered, p)
					break
				}
			}
		}
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	total := len(filtered)
	totalPages := (total + DefaultPostsPerPage - 1) / DefaultPostsPerPage
	if totalPages < 1 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}
	start := (page - 1) * DefaultPostsPerPage
	end := start + DefaultPostsPerPage
	if end > total {
		end = total
	}
	pagePosts := filtered[start:end]

	var viewCounts map[string]int
	if h.Comments != nil {
		if vc, err := h.Comments.AllViewCounts(); err != nil {
			slog.Error("blog partial view counts query failed", "error", err)
		} else {
			viewCounts = vc
		}
	}

	hasMore := page < totalPages
	h.Render(w, r, templates.BlogPostItems(ui, pagePosts, lang, viewCounts, page+1, hasMore, tagFilter))
}

func relatedPosts(posts []post.Post, current post.Post) []post.Post {
	if len(posts) == 0 || len(current.Tags) == 0 {
		return nil
	}
	tagSet := make(map[string]bool, len(current.Tags))
	for _, t := range current.Tags {
		norm := strings.ToLower(strings.TrimSpace(t))
		if norm != "" {
			tagSet[norm] = true
		}
	}
	if len(tagSet) == 0 {
		return nil
	}

	var sameLang []post.Post
	var otherLang []post.Post
	seen := make(map[string]bool)
	seen[current.Slug] = true

	for _, p := range posts {
		if p.Draft || seen[p.Slug] {
			continue
		}
		matches := false
		for _, t := range p.Tags {
			if tagSet[strings.ToLower(strings.TrimSpace(t))] {
				matches = true
				break
			}
		}
		if matches {
			seen[p.Slug] = true
			if p.Lang == current.Lang {
				sameLang = append(sameLang, p)
			} else {
				otherLang = append(otherLang, p)
			}
		}
	}

	// Combine same language first, fallback to other languages if needed
	var out []post.Post
	out = append(out, sameLang...)
	if len(out) < 4 {
		out = append(out, otherLang...)
	}

	// Cap at 4 posts max for clean 4-column grid
	if len(out) > 4 {
		out = out[:4]
	}
	return out
}
