package auth

import (
	authstore "daemontalk/internal/auth"
	"daemontalk/internal/comment"
	"daemontalk/internal/forum"
	"daemontalk/internal/post"
)

type Handler struct {
	Auth         *authstore.Store
	GitHubOAuth  *authstore.GitHubOAuth
	Comments     *comment.Store
	Forum        *forum.Store
	AllPosts     func() []post.Post
	IsProduction bool
}

func (h *Handler) GetAllPosts() []post.Post {
	if h.AllPosts != nil {
		return h.AllPosts()
	}
	return nil
}
