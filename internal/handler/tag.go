package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"daemontalk/internal/i18n"
	"daemontalk/internal/post"
	"daemontalk/web/templates"
)

func (h *Handler) TagIndex(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)
	tag := chi.URLParam(r, "tag")
	isAdmin := h.isAdmin(r)

	var filtered []post.Post
	for _, p := range h.VisiblePosts(isAdmin) {
		for _, t := range p.Tags {
			if strings.EqualFold(t, tag) {
				filtered = append(filtered, p)
				break
			}
		}
	}

	if len(filtered) == 0 {
		h.NotFound(w, r)
		return
	}

	pageStr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	pageSize := 12
	total := len(filtered)
	totalPages := (total + pageSize - 1) / pageSize
	if totalPages < 1 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}
	pagePosts := filtered[start:end]

	var viewCounts map[string]int
	if h.Comments != nil {
		if vc, err := h.Comments.AllViewCounts(); err != nil {
			slog.Error("tag view counts query failed", "tag", tag, "error", err)
		} else {
			viewCounts = vc
		}
	}

	h.Render(w, r, templates.Layout(ui, lang, "#"+tag, r.URL.Path, templates.PageMeta{
		Description: fmt.Sprintf("Posts tagged #%s on daemontalk.com", tag),
	}, templates.TagPage(ui, tag, filtered, pagePosts, lang, viewCounts, page, totalPages)))
}

func (h *Handler) TagPostsPartial(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)
	tag := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("tag")))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 6 {
		offset = 6
	}

	isAdmin := h.isAdmin(r)
	var filtered []post.Post
	for _, p := range h.VisiblePosts(isAdmin) {
		for _, t := range p.Tags {
			if strings.EqualFold(t, tag) {
				filtered = append(filtered, p)
				break
			}
		}
	}

	total := len(filtered)
	if offset >= total {
		h.Render(w, r, templates.TagRiverItems(ui, nil, lang, offset, 0, 0, tag))
		return
	}

	batchSize := 12
	end := offset + batchSize
	if end > total {
		end = total
	}

	pagePosts := filtered[offset:end]
	nextOffset := end
	remaining := total - end

	h.Render(w, r, templates.TagRiverItems(ui, pagePosts, lang, offset, nextOffset, remaining, tag))
}

