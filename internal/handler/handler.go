package handler

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/a-h/templ"
	"daemontalk/internal/comment"
	"daemontalk/internal/i18n"
	"daemontalk/internal/post"
	"daemontalk/internal/postdb"
	"daemontalk/internal/project"
	"daemontalk/web/templates"
)

type Handler struct {
	ContentDir  string // path to content directory (defaults to "content")
	AllProjects []project.Project
	FilePosts   []post.Post   // post markdown dari content/posts, dimuat sekali saat startup
	PostDB      *postdb.Store // post buatan editor web (opsional)
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

	// GitHub token for higher API rate limits (optional)
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

// RefreshPosts merender ulang post dari DB, menggabungkannya dengan post file,
// dan mengganti snapshot. Dipanggil saat startup dan setiap kali editor
// menyimpan/menghapus post — post baru langsung tampil tanpa restart.
func (h *Handler) RefreshPosts() {
	merged := make([]post.Post, 0, len(h.FilePosts)+8)
	merged = append(merged, h.FilePosts...)

	if h.PostDB != nil {
		webPosts, err := h.PostDB.List()
		if err != nil {
			log.Printf("refresh posts: list db: %v", err)
		}
		for _, wp := range webPosts {
			p, err := post.Parse(wp.ToMarkdown())
			if err != nil {
				log.Printf("refresh posts: render %q: %v", wp.Slug, err)
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

// isAdmin reports whether the request carries a valid admin token cookie.
func (h *Handler) isAdmin(r *http.Request) bool {
	if h.AdminToken == "" {
		return false
	}
	c, err := r.Cookie("admin_token")
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
		log.Printf("warn: load %s: %v", filename, err)
	}

	if err := templates.Layout(ui, lang, title, r.URL.Path, meta,
		render(ui, body, lang)).Render(r.Context(), w); err != nil {
		log.Printf("render error: %v", err)
	}
}
