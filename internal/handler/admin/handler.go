package admin

import (
	"net/http"

	"daemontalk/internal/comment"
	"daemontalk/internal/handler/common"
	"daemontalk/internal/post"
)

type Handler struct {
	AdminToken      string
	Comments        *comment.Store
	ContentDir      string
	FilePosts       func() []post.Post
	ReloadFilePosts func()
	VisiblePosts    func(isAdmin bool) []post.Post
	AllPosts        func() []post.Post
	RefreshPosts    func()
	NotFound        http.HandlerFunc
}

func (h *Handler) IsAdmin(r *http.Request) bool {
	return common.IsAdmin(h.AdminToken, r)
}

func (h *Handler) HandleNotFound(w http.ResponseWriter, r *http.Request) {
	if h.NotFound != nil {
		h.NotFound(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func (h *Handler) GetContentPath(subpath string) string {
	return common.GetContentPath(h.ContentDir, subpath)
}

func (h *Handler) GetFilePosts() []post.Post {
	if h.FilePosts != nil {
		return h.FilePosts()
	}
	return nil
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

func (h *Handler) DoRefreshPosts() {
	if h.RefreshPosts != nil {
		h.RefreshPosts()
	}
}

func (h *Handler) DoReloadFilePosts() {
	if h.ReloadFilePosts != nil {
		h.ReloadFilePosts()
	}
}
