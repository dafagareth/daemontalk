package pages

import (
	"net/http"

	"daemontalk/internal/auth"
	"daemontalk/internal/comment"
	"daemontalk/internal/forum"
	"daemontalk/internal/handler/common"
	"daemontalk/internal/post"
)

type Handler struct {
	ContentDir   string
	GitHubToken  string
	Comments     *comment.Store
	Forum        *forum.Store
	Auth         *auth.Store
	AllPosts     func() []post.Post
	VisiblePosts func(isAdmin bool) []post.Post
	AdminToken   string
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPass     string
	SMTPTo       string
}

func (h *Handler) IsAdmin(r *http.Request) bool {
	return common.IsAdmin(h.AdminToken, r)
}

func (h *Handler) GetAllPosts() []post.Post {
	if h.AllPosts != nil {
		return h.AllPosts()
	}
	return nil
}

func (h *Handler) GetVisiblePosts(isAdmin bool) []post.Post {
	if h.VisiblePosts != nil {
		return h.VisiblePosts(isAdmin)
	}
	return nil
}
