package blog

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
	"daemontalk/internal/handler/common"
	"daemontalk/internal/i18n"
	"daemontalk/internal/post"
	blogtmpl "daemontalk/web/templates/blog"
	"daemontalk/web/templates/layouts"
	"daemontalk/web/templates/shared"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) BlogIndex(w http.ResponseWriter, r *http.Request) {
	lang := common.LangFromRequest(r)
	ui := i18n.Get(lang)
	isAdmin := h.isAdmin(r)

	if tagFilter := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("tag"))); tagFilter != "" {
		http.Redirect(w, r, "/blog/tag/"+url.PathEscape(tagFilter), http.StatusMovedPermanently)
		return
	}

	visible := h.getVisiblePosts(isAdmin)

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
	pageSize := common.DefaultPostsPerPage
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
	meta := shared.PageMeta{}
	if r.URL.Path == "/" {
		pageName = "home"
		meta.JSONLD = common.SiteJSONLD()
	}
	h.render(w, r, layouts.Layout(ui, lang, pageName, r.URL.Path, meta, blogtmpl.BlogIndex(ui, filtered, pagePosts, lang, page, totalPages, viewCounts, tagCounts, "")))
}

func (h *Handler) BlogPost(w http.ResponseWriter, r *http.Request) {
	lang := common.LangFromRequest(r)
	ui := i18n.Get(lang)
	slug := chi.URLParam(r, "slug")
	isAdmin := h.isAdmin(r)

	p, ok := post.FindBySlug(h.getAllPosts(), slug)
	if !ok {
		h.notFound(w, r)
		return
	}

	if p.Slug != slug {
		target := "/blog/" + p.Slug
		if r.URL.RawQuery != "" {
			target += "?" + r.URL.RawQuery
		}
		http.Redirect(w, r, target, http.StatusMovedPermanently)
		return
	}

	if p.Draft && !isAdmin {
		h.notFound(w, r)
		return
	}
	if !p.PublishAt.IsZero() && p.PublishAt.After(time.Now()) && !isAdmin {
		h.notFound(w, r)
		return
	}

	var nav blogtmpl.PostNav
	visible := h.getVisiblePosts(isAdmin)
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
		for _, sp := range h.getAllPosts() {
			if sp.Series == p.Series && !sp.Draft {
				seriesParts = append(seriesParts, sp)
			}
		}
		sort.Slice(seriesParts, func(i, j int) bool {
			return seriesParts[i].SeriesPart < seriesParts[j].SeriesPart
		})
	}

	related := relatedPosts(h.getVisiblePosts(isAdmin), p)

	author := p.Author
	if author == "" {
		author = "daemontalk"
	}
	meta := shared.PageMeta{
		Description:   p.Description,
		Type:          "article",
		PublishedTime: p.Date.Format("2006-01-02T15:04:05Z07:00"),
		Author:        author,
		Image:         h.absoluteURL(r, "/blog/"+slug+"/og.png"),
		URL:           h.absoluteURL(r, r.URL.Path),
	}
	meta.JSONLD = common.ArticleJSONLD(p, meta.Image)

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
		if isAdmin || common.IsBot(r) {
			views, _ = h.Comments.ViewCount(slug)
		} else {
			cookieKey := common.CookieViewCooldownPrefix + slug
			if _, err := r.Cookie(cookieKey); err == nil {
				views, _ = h.Comments.ViewCount(slug)
			} else {
				viewerKey := common.GetViewerKey(w, r, authUser)
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
							MaxAge:   common.CookieViewCooldownMaxAge,
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
	if cookie, err := r.Cookie(common.CookieReactedPrefix + slug); err == nil && cookie.Value != "" {
		userReaction, _ = url.QueryUnescape(cookie.Value)
	}

	visitorName := common.GetVisitorIdentity(w, r)
	if isAdmin {
		visitorName = "daemontalk"
	}

	h.render(w, r, layouts.Layout(ui, lang, p.Title+" · daemontalk", r.URL.Path, meta,
		blogtmpl.BlogPostPage(ui, p, related, comments, views, isAdmin, lang, reactions, seriesParts, nav, userReaction, visitorName, authUser),
	))
}
