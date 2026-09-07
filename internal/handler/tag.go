package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"daemontalk/internal/i18n"
	"daemontalk/internal/post"
	"daemontalk/web/templates"
	"github.com/go-chi/chi/v5"
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
	}, templates.TagPage(ui, tag, filtered, lang, viewCounts)))
}

func (h *Handler) TagPostsPartial(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)
	tag := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("tag")))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 8 {
		offset = 8
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

func (h *Handler) RedirectTag(w http.ResponseWriter, r *http.Request) {
	tag := strings.TrimSpace(r.URL.Query().Get("tag"))
	lang := langFromRequest(r)
	prefix := ""
	if lang == "id" {
		prefix = "/id"
	}

	if tag != "" {
		http.Redirect(w, r, prefix+"/blog/tag/"+url.PathEscape(strings.ToLower(tag)), http.StatusMovedPermanently)
		return
	}

	if prefix == "" {
		http.Redirect(w, r, "/blog", http.StatusMovedPermanently)
	} else {
		http.Redirect(w, r, prefix, http.StatusMovedPermanently)
	}
}
