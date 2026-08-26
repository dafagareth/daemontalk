package handler

import (
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"daemontalk/internal/comment"
	"daemontalk/internal/i18n"
	"daemontalk/internal/post"
	"daemontalk/internal/postdb"
	"daemontalk/internal/project"
	"daemontalk/web/templates"
	"github.com/a-h/templ"
)

type Handler struct {
	ContentDir  string // path to content directory (defaults to "content")
	AllProjects []project.Project
	FilePosts   []post.Post   // markdown posts from content/posts, loaded once at startup
	PostDB      *postdb.Store // posts created by web editor (optional)
	Comments    *comment.Store
	AdminToken  string // when set, enables comment moderation

	// merged = FilePosts + render PostDB, urut tanggal desc. Di-swap utuh oleh
	// RefreshPosts sehingga pembaca tidak perlu lock.
	merged atomic.Pointer[[]post.Post]

	// SMTP config for contact form (all optional; form still accepts without them)
	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string
	SMTPTo   string

	// token for higher API rate limits (optional)
	GitHubToken string
}

// AllPosts mengembalikan snapshot gabungan post file + post DB.
func (h *Handler) AllPosts() []post.Post {
	if p := h.merged.Load(); p != nil {
		return *p
	}
	h.RefreshPosts()
	if p := h.merged.Load(); p != nil {
		return *p
	}
	return nil
}

// ReloadFilePosts re-reads markdown posts from the content/posts directory.
func (h *Handler) ReloadFilePosts() {
	postsDir := h.getContentPath("posts")
	if fps, err := post.LoadAllWithDrafts(postsDir); err == nil {
		h.FilePosts = fps
	}
}

// RefreshPosts re-renders posts from DB, merges them with file posts,
// and replaces snapshot. Called at startup and every time the editor
// saving/deleting post — new post directly appears without restart.
func (h *Handler) RefreshPosts() {
	merged := make([]post.Post, 0, len(h.FilePosts)+8)
	merged = append(merged, h.FilePosts...)

	if h.PostDB != nil {
		webPosts, err := h.PostDB.List()
		if err != nil {
			slog.Error("refresh posts list db failed", "error", err)
		}
		for _, wp := range webPosts {
			p, err := post.Parse(wp.ToMarkdown())
			if err != nil {
				slog.Error("refresh posts render failed", "slug", wp.Slug, "error", err)
				continue
			}
			if wp.Description != "" {
				p.Description = wp.Description
			}
			merged = append(merged, p)
		}
	}

	sort.SliceStable(merged, func(i, j int) bool {
		return merged[i].Date.After(merged[j].Date)
	})
	h.merged.Store(&merged)
}

func langFromRequest(r *http.Request) string {
	if strings.HasPrefix(r.URL.Path, "/id") {
		return "id"
	}
	return "en"
}

// IsRadarEnabled reports whether the systems knowledge graph/radar feature is enabled.
func (h *Handler) IsRadarEnabled() bool {
	return templates.IsRadarEnabled()
}

// isAdmin reports whether the request carries a valid admin token cookie.
func (h *Handler) isAdmin(r *http.Request) bool {
	if h.AdminToken == "" {
		return false
	}
	c, err := r.Cookie(CookieAdminToken)
	if err != nil {
		return false
	}
	return c.Value == h.AdminToken
}

// VisiblePosts returns all posts that should be visible to the current user.
func (h *Handler) VisiblePosts(isAdmin bool) []post.Post {
	var out []post.Post
	for _, p := range h.AllPosts() {
		if p.Draft && !isAdmin {
			continue
		}
		if !p.PublishAt.IsZero() && p.PublishAt.After(time.Now()) && !isAdmin {
			continue
		}
		out = append(out, p)
	}
	return out
}

// getContentPath resolves a subpath relative to h.ContentDir (defaulting to "content").
func (h *Handler) getContentPath(subpath string) string {
	dir := h.ContentDir
	if dir == "" {
		dir = "content"
	}
	return filepath.Join(dir, subpath)
}

// Render executes a templ component and logs any rendering errors with structured metadata.
func (h *Handler) Render(w http.ResponseWriter, r *http.Request, c templ.Component) {
	if err := c.Render(r.Context(), w); err != nil {
		slog.Error("render component failed", "error", err, "path", r.URL.Path, "method", r.Method)
	}
}

// renderMarkdownPage renders a static markdown page with the given contentKey.
func (h *Handler) renderMarkdownPage(w http.ResponseWriter, r *http.Request,
	contentKey, title string, meta templates.PageMeta,
	render func(i18n.UI, template.HTML, string) templ.Component,
) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)

	filename := h.getContentPath(contentKey + ".md")
	if lang == "id" {
		filename = h.getContentPath(contentKey + ".id.md")
	}

	body, err := post.LoadBody(filename)
	if err != nil {
		slog.Warn("load markdown page failed", "file", filename, "error", err)
	}

	h.Render(w, r, templates.Layout(ui, lang, title, r.URL.Path, meta, render(ui, body, lang)))
}
