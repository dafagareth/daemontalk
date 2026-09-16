package blog

import (
	"net/http"

	"daemontalk/internal/comment"
	"daemontalk/internal/handler/common"
	"daemontalk/internal/post"
	"github.com/a-h/templ"
)

type Handler struct {
	AllPosts     func() []post.Post
	VisiblePosts func(isAdmin bool) []post.Post
	Comments     *comment.Store
	AdminToken   string
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPass     string
	SMTPTo       string
}

func (h *Handler) getAllPosts() []post.Post {
	if h.AllPosts != nil {
		return h.AllPosts()
	}
	return nil
}

func (h *Handler) getVisiblePosts(isAdmin bool) []post.Post {
	if h.VisiblePosts != nil {
		return h.VisiblePosts(isAdmin)
	}
	var out []post.Post
	for _, p := range h.getAllPosts() {
		if !p.Draft || isAdmin {
			out = append(out, p)
		}
	}
	return out
}

func (h *Handler) isAdmin(r *http.Request) bool {
	return common.IsAdmin(h.AdminToken, r)
}

func (h *Handler) render(w http.ResponseWriter, r *http.Request, c templ.Component) {
	common.Render(w, r, c)
}

func (h *Handler) notFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func (h *Handler) absoluteURL(r *http.Request, path string) string {
	return common.AbsoluteURL(r, path)
}
