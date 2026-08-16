package handler

import (
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"daemontalk/internal/comment"
	"daemontalk/internal/i18n"
	"daemontalk/internal/post"
	"daemontalk/web/templates"

	"github.com/go-chi/chi/v5"
)

const postsPerPage = 14

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

	// Collect visible posts (filter drafts and scheduled for non-admin)
	visible := h.VisiblePosts(h.isAdmin(r))

	// Compute tag counts from all visible posts (not just the current page).
	tagCounts := make(map[string]int)
	for _, p := range visible {
		for _, t := range p.Tags {
			tagCounts[t]++
		}
	}

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
	totalPages := (total + postsPerPage - 1) / postsPerPage
	if totalPages < 1 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}
	var pagePosts []post.Post
	if tagFilter == "" && page == 1 {
		pagePosts = filtered
	} else {
		start := (page - 1) * postsPerPage
		end := start + postsPerPage
		if end > total {
			end = total
		}
		pagePosts = filtered[start:end]
	}

	var viewCounts map[string]int
	if h.Comments != nil {
		if vc, err := h.Comments.AllViewCounts(); err != nil {
			log.Printf("blog index view counts: %v", err)
		} else {
			viewCounts = vc
		}
	}

	// Setelah pivot blog-first, BlogIndex juga melayani "/" dan "/id" sebagai
	// halaman utama: judul tab "daemontalk" + JSON-LD situs.
	pageName := "blog"
	meta := templates.PageMeta{}
	if r.URL.Path == "/" || r.URL.Path == "/id" || r.URL.Path == "/id/" {
		pageName = "home"
		meta.JSONLD = siteJSONLD()
	}
	if err := templates.Layout(ui, lang, pageName, r.URL.Path, meta, templates.BlogIndex(ui, pagePosts, lang, page, totalPages, viewCounts, tagCounts, tagFilter)).Render(r.Context(), w); err != nil {
		log.Printf("render error: %v", err)
	}
}

// BlogPostsPartial returns a partial HTML response (list items + updated Load More button)
// for HTMX infinite-style loading. It does not render the full layout.

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

	// If accessed via alias (old slug), redirect 301 to canonical Short ID URL.
	if p.Slug != slug {
		target := "/blog/" + p.Slug
		if r.URL.RawQuery != "" {
			target += "?" + r.URL.RawQuery
		}
		http.Redirect(w, r, target, http.StatusMovedPermanently)
		return
	}

	// Draft and scheduled posts are only visible to admins.
	if p.Draft && !isAdmin {
		h.NotFound(w, r)
		return
	}
	if !p.PublishAt.IsZero() && p.PublishAt.After(time.Now()) && !isAdmin {
		h.NotFound(w, r)
		return
	}

	// Admin login via ?admin=TOKEN sets a cookie, then redirect to clean URL.
	if h.AdminToken != "" {
		if tok := r.URL.Query().Get("admin"); tok != "" {
			if tok == h.AdminToken {
				http.SetCookie(w, &http.Cookie{
					Name:     "admin_token",
					Value:    tok,
					Path:     "/",
					HttpOnly: true,
					SameSite: http.SameSiteLaxMode,
					MaxAge:   60 * 60 * 24 * 30,
				})
			}
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
			return
		}
	}

	related := relatedPosts(h.AllPosts(), p)

	// Series: collect all non-draft posts in the same series, sorted by part.
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

	// Compute prev/next navigation from the visible post list.
	// h.AllPosts() is sorted newest-first, so index i+1 = older, i-1 = newer.
	var nav templates.PostNav
	visible := h.VisiblePosts(isAdmin)
	for i, vp := range visible {
		if vp.Slug == slug {
			if i+1 < len(visible) {
				nav.HasPrev = true
				nav.Prev = visible[i+1] // older post
			}
			if i > 0 {
				nav.HasNext = true
				nav.Next = visible[i-1] // newer post
			}
			break
		}
	}

	// Per-post OG image: use the cover if set, else the auto-generated card.
	meta := templates.PageMeta{
		Description:   p.Description,
		Type:          "article",
		PublishedTime: p.Date.Format("2006-01-02T15:04:05Z07:00"),
		Author:        "Dafa Gareth",
	}
	if p.Cover != "" {
		meta.Image = templates.AbsoluteURL(p.Cover)
	} else {
		meta.Image = templates.AbsoluteURL("/blog/" + slug + "/og.png")
	}
	meta.JSONLD = articleJSONLD(p, meta.Image)

	var comments []comment.Comment
	views := 0
	var reactions map[string]int
	if h.Comments != nil {
		if cs, err := h.Comments.ListBySlug(slug); err != nil {
			log.Printf("load comments for %s: %v", slug, err)
		} else {
			comments = cs
		}
		if isAdmin {
			views, _ = h.Comments.ViewCount(slug)
		} else {
			cookieKey := "v_post_" + slug
			if _, err := r.Cookie(cookieKey); err == nil {
				views, _ = h.Comments.ViewCount(slug)
			} else {
				if n, err := h.Comments.IncrementView(slug); err != nil {
					log.Printf("increment view for %s: %v", slug, err)
				} else {
					views = n
				}
				http.SetCookie(w, &http.Cookie{
					Name:     cookieKey,
					Value:    "1",
					Path:     "/",
					MaxAge:   3600 * 12, // 12-hour cooldown per post
					HttpOnly: true,
					SameSite: http.SameSiteLaxMode,
				})
			}
		}
		if rx, err := h.Comments.GetReactions(slug); err != nil {
			log.Printf("load reactions for %s: %v", slug, err)
		} else {
			reactions = rx
		}
	}

	var userReaction string
	if cookie, err := r.Cookie("reacted_" + slug); err == nil && cookie.Value != "" {
		userReaction, _ = url.QueryUnescape(cookie.Value)
	}

	visitorName := GetVisitorIdentity(w, r)

	if err := templates.Layout(ui, lang, p.Title+" · daemontalk", r.URL.Path, meta,
		templates.BlogPostPage(ui, p, related, comments, views, isAdmin, lang, reactions, seriesParts, nav, userReaction, visitorName),
	).Render(r.Context(), w); err != nil {
		log.Printf("render error: %v", err)
	}
}

// DeleteComment removes a comment (admin only) and returns the refreshed list.

// PostComment accepts a new comment (HTMX form POST) and returns the refreshed
// comment list to swap in place.
