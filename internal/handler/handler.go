package handler

import (
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"daemontalk/internal/auth"
	"daemontalk/internal/comment"
	"daemontalk/internal/forum"
	"daemontalk/internal/i18n"
	"daemontalk/internal/post"
	"daemontalk/internal/postdb"
	"daemontalk/internal/project"
	"daemontalk/web/templates"
	"github.com/a-h/templ"
)

type Handler struct {
	ContentDir   string
	AllProjects  []project.Project
	FilePosts    []post.Post
	filePostsMu  sync.RWMutex
	PostDB       *postdb.Store
	Comments     *comment.Store
	Auth         *auth.Store
	GitHubOAuth  *auth.GitHubOAuth
	Forum        *forum.Store
	IsProduction bool
	AdminToken   string

	merged atomic.Pointer[[]post.Post]

	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string
	SMTPTo   string

	GitHubToken string
}

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

func (h *Handler) ReloadFilePosts() {
	postsDir := h.getContentPath("posts")
	if fps, err := post.LoadAllWithDrafts(postsDir); err == nil {
		h.filePostsMu.Lock()
		h.FilePosts = fps
		h.filePostsMu.Unlock()
	}
}

func (h *Handler) getFilePosts() []post.Post {
	h.filePostsMu.RLock()
	defer h.filePostsMu.RUnlock()
	out := make([]post.Post, len(h.FilePosts))
	copy(out, h.FilePosts)
	return out
}

func (h *Handler) RefreshPosts() {
	h.filePostsMu.RLock()
	merged := make([]post.Post, 0, len(h.FilePosts)+8)
	merged = append(merged, h.FilePosts...)
	h.filePostsMu.RUnlock()

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

func (h *Handler) IsRadarEnabled() bool {
	return templates.IsRadarEnabled()
}

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

func (h *Handler) getContentPath(subpath string) string {
	dir := h.ContentDir
	if dir == "" {
		dir = "content"
	}
	return filepath.Join(dir, subpath)
}

func (h *Handler) Render(w http.ResponseWriter, r *http.Request, c templ.Component) {
	if err := c.Render(r.Context(), w); err != nil {
		slog.Error("render component failed", "error", err, "path", r.URL.Path, "method", r.Method)
	}
}

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

func urlPrefix(lang string) string {
	if lang == "id" {
		return "/id"
	}
	return ""
}

func (h *Handler) AbsoluteURL(r *http.Request, path string) string {
	if path == "" {
		return ""
	}
	if len(path) >= 4 && path[:4] == "http" {
		return path
	}
	if path[0] != '/' {
		path = "/" + path
	}

	base := templates.SiteBaseURL
	if r != nil {
		proto := "http"
		if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
			proto = "https"
		}
		host := r.Header.Get("X-Forwarded-Host")
		if host == "" {
			host = r.Host
		}
		if host != "" {
			base = proto + "://" + host
		}
	}
	return strings.TrimSuffix(base, "/") + path
}
