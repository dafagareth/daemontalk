package handler

import (
	"log/slog"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"daemontalk/internal/auth"
	"daemontalk/internal/comment"
	"daemontalk/internal/i18n"
	"daemontalk/internal/post"
	"daemontalk/web/templates"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) BlogIndex(w http.ResponseWriter, r *http.Request) {
	if IsCLIRequest(r) {
		if tag := r.URL.Query().Get("tag"); tag != "" {
			h.CLITag(w, r)
			return
		}
		h.CLIMain(w, r)
		return
	}

	lang := langFromRequest(r)
	ui := i18n.Get(lang)
	isAdmin := h.isAdmin(r)

	if tagFilter := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("tag"))); tagFilter != "" {
		prefix := ""
		if lang == "id" {
			prefix = "/id"
		}
		http.Redirect(w, r, prefix+"/blog/tag/"+url.PathEscape(tagFilter), http.StatusMovedPermanently)
		return
	}

	visible := h.VisiblePosts(isAdmin)

	tagCounts := make(map[string]int)
	for _, p := range visible {
		for _, t := range p.Tags {
			tagCounts[strings.ToLower(t)]++
		}
	}

	filtered := visible

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize := DefaultPostsPerPage
	total := len(filtered)
	totalPages := (total + pageSize - 1) / pageSize
	if totalPages < 1 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}

	var pagePosts []post.Post
	if total == 0 {
		pagePosts = nil
	} else {
		start := (page - 1) * pageSize
		end := start + pageSize
		if end > total {
			end = total
		}
		pagePosts = filtered[start:end]
	}

	var viewCounts map[string]int
	if h.Comments != nil {
		if vc, err := h.Comments.AllViewCounts(); err != nil {
			slog.Error("blog index view counts query failed", "error", err)
		} else {
			viewCounts = vc
		}
	}

	pageName := "blog"
	meta := templates.PageMeta{}
	if r.URL.Path == "/" || r.URL.Path == "/id" || r.URL.Path == "/id/" {
		pageName = "home"
		meta.JSONLD = siteJSONLD()
	}
	h.Render(w, r, templates.Layout(ui, lang, pageName, r.URL.Path, meta, templates.BlogIndex(ui, filtered, pagePosts, lang, page, totalPages, viewCounts, tagCounts, "")))
}

func (h *Handler) BlogPost(w http.ResponseWriter, r *http.Request) {
	if IsCLIRequest(r) {
		h.CLIPost(w, r)
		return
	}

	lang := langFromRequest(r)
	ui := i18n.Get(lang)
	slug := chi.URLParam(r, "slug")
	isAdmin := h.isAdmin(r)

	p, ok := post.FindBySlug(h.AllPosts(), slug)
	if !ok {
		h.NotFound(w, r)
		return
	}

	if p.Slug != slug {
		prefix := ""
		if lang == "id" {
			prefix = "/id"
		}
		target := prefix + "/blog/" + p.Slug
		if r.URL.RawQuery != "" {
			target += "?" + r.URL.RawQuery
		}
		http.Redirect(w, r, target, http.StatusMovedPermanently)
		return
	}

	if p.Draft && !isAdmin {
		h.NotFound(w, r)
		return
	}
	if !p.PublishAt.IsZero() && p.PublishAt.After(time.Now()) && !isAdmin {
		h.NotFound(w, r)
		return
	}

	var nav templates.PostNav
	visible := h.VisiblePosts(isAdmin)
	for i, vp := range visible {
		if vp.Slug == slug {
			if i+1 < len(visible) {
				nav.HasPrev = true
				nav.Prev = visible[i+1]
			}
			if i > 0 {
				nav.HasNext = true
				nav.Next = visible[i-1]
			}
			break
		}
	}

	var seriesParts []post.Post
	if p.Series != "" {
		for _, sp := range h.AllPosts() {
			if sp.Series == p.Series && !sp.Draft {
				seriesParts = append(seriesParts, sp)
			}
		}
		sort.Slice(seriesParts, func(i, j int) bool {
			return seriesParts[i].SeriesPart < seriesParts[j].SeriesPart
		})
	}

	related := relatedPosts(h.VisiblePosts(isAdmin), p)

	author := p.Author
	if author == "" {
		author = "daemontalk"
	}
	meta := templates.PageMeta{
		Description:   p.Description,
		Type:          "article",
		PublishedTime: p.Date.Format("2006-01-02T15:04:05Z07:00"),
		Author:        author,
		Image:         h.AbsoluteURL(r, "/blog/"+slug+"/og.png"),
		URL:           h.AbsoluteURL(r, r.URL.Path),
	}
	meta.JSONLD = articleJSONLD(p, meta.Image)

	authUser := auth.GetUser(r.Context())
	var comments []comment.Comment
	views := 0
	var reactions map[string]int
	if h.Comments != nil {
		if cs, err := h.Comments.ListBySlug(slug); err != nil {
			slog.Error("load comments query failed", "slug", slug, "error", err)
		} else {
			comments = cs
		}
		if isAdmin || IsCLIRequest(r) || isBot(r) {
			views, _ = h.Comments.ViewCount(slug)
		} else {
			cookieKey := CookieViewCooldownPrefix + slug
			if _, err := r.Cookie(cookieKey); err == nil {
				views, _ = h.Comments.ViewCount(slug)
			} else {
				viewerKey := GetViewerKey(w, r, authUser)
				if n, recorded, err := h.Comments.RecordPostView(slug, viewerKey); err != nil {
					slog.Error("record post view failed", "slug", slug, "error", err)
					views, _ = h.Comments.ViewCount(slug)
				} else {
					views = n
					if recorded {
						http.SetCookie(w, &http.Cookie{
							Name:     cookieKey,
							Value:    "1",
							Path:     "/",
							MaxAge:   CookieViewCooldownMaxAge,
							HttpOnly: true,
							SameSite: http.SameSiteLaxMode,
						})
					}
				}
			}
		}
		if rx, err := h.Comments.GetReactions(slug); err != nil {
			slog.Error("load reactions query failed", "slug", slug, "error", err)
		} else {
			reactions = rx
		}
	}

	var userReaction string
	if cookie, err := r.Cookie(CookieReactedPrefix + slug); err == nil && cookie.Value != "" {
		userReaction, _ = url.QueryUnescape(cookie.Value)
	}

	visitorName := GetVisitorIdentity(w, r)
	if isAdmin {
		visitorName = "daemontalk"
	}

	h.Render(w, r, templates.Layout(ui, lang, p.Title+" · daemontalk", r.URL.Path, meta,
		templates.BlogPostPage(ui, p, related, comments, views, isAdmin, lang, reactions, seriesParts, nav, userReaction, visitorName, authUser),
	))
}
