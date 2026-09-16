package distribution

import (
	"net/http"

	"daemontalk/internal/handler/common"
	"daemontalk/internal/post"
	"github.com/a-h/templ"
)

const (
	SeoBaseURL = "https://daemontalk.com"
)

type Handler struct {
	AllPosts        func() []post.Post
	ReloadFilePosts func()
	RefreshPosts    func()
}

func (h *Handler) getAllPosts() []post.Post {
	if h.AllPosts != nil {
		return h.AllPosts()
	}
	return nil
}

func (h *Handler) render(w http.ResponseWriter, r *http.Request, c templ.Component) {
	common.Render(w, r, c)
}

func (h *Handler) absoluteURL(r *http.Request, path string) string {
	return common.AbsoluteURL(r, path)
}
