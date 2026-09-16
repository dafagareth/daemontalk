package handler

import (
	"net/http"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"daemontalk/internal/auth"
	"daemontalk/internal/comment"
	"daemontalk/internal/forum"
	"daemontalk/internal/handler/common"
	"daemontalk/internal/handler/distribution"
	"daemontalk/internal/handler/pages"
	"daemontalk/internal/post"
	"github.com/a-h/templ"
)

type Handler struct {
	ContentDir   string
	FilePosts    []post.Post
	filePostsMu  sync.RWMutex
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

	pagesOnce    sync.Once
	pagesHandler *pages.Handler
	distOnce     sync.Once
	distHandler  *distribution.Handler
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
	postsDir := common.GetContentPath(h.ContentDir, "posts")
	if fps, err := post.LoadAllWithDrafts(postsDir); err == nil {
		h.filePostsMu.Lock()
		h.FilePosts = fps
		h.filePostsMu.Unlock()
	}
}

func (h *Handler) GetFilePosts() []post.Post {
	h.filePostsMu.RLock()
	defer h.filePostsMu.RUnlock()
	out := make([]post.Post, len(h.FilePosts))
	copy(out, h.FilePosts)
	return out
}

func (h *Handler) RefreshPosts() {
	h.filePostsMu.RLock()
	merged := make([]post.Post, len(h.FilePosts))
	copy(merged, h.FilePosts)
	h.filePostsMu.RUnlock()

	sort.SliceStable(merged, func(i, j int) bool {
		return merged[i].Date.After(merged[j].Date)
	})
	h.merged.Store(&merged)
}

func (h *Handler) isAdmin(r *http.Request) bool {
	return common.IsAdmin(h.AdminToken, r)
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

func (h *Handler) Render(w http.ResponseWriter, r *http.Request, c templ.Component) {
	common.Render(w, r, c)
}

func (h *Handler) AbsoluteURL(r *http.Request, path string) string {
	return common.AbsoluteURL(r, path)
}

func stripCRLF(s string) string {
	return common.StripCRLF(s)
}
