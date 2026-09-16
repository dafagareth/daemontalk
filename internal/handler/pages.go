package handler

import (
	"net/http"

	"daemontalk/internal/handler/pages"
)

func (h *Handler) PagesHandler() *pages.Handler {
	h.pagesOnce.Do(func() {
		h.pagesHandler = &pages.Handler{
			ContentDir:   h.ContentDir,
			GitHubToken:  h.GitHubToken,
			Comments:     h.Comments,
			Forum:        h.Forum,
			Auth:         h.Auth,
			AllPosts:     h.AllPosts,
			VisiblePosts: h.VisiblePosts,
			AdminToken:   h.AdminToken,
			SMTPHost:     h.SMTPHost,
			SMTPPort:     h.SMTPPort,
			SMTPUser:     h.SMTPUser,
			SMTPPass:     h.SMTPPass,
			SMTPTo:       h.SMTPTo,
		}
	})
	return h.pagesHandler
}

func (h *Handler) About(w http.ResponseWriter, r *http.Request) {
	h.PagesHandler().About(w, r)
}

func (h *Handler) Accessibility(w http.ResponseWriter, r *http.Request) {
	h.PagesHandler().Accessibility(w, r)
}

func (h *Handler) License(w http.ResponseWriter, r *http.Request) {
	h.PagesHandler().License(w, r)
}

func (h *Handler) Privacy(w http.ResponseWriter, r *http.Request) {
	h.PagesHandler().Privacy(w, r)
}

func (h *Handler) Terms(w http.ResponseWriter, r *http.Request) {
	h.PagesHandler().Terms(w, r)
}

func (h *Handler) Saved(w http.ResponseWriter, r *http.Request) {
	h.PagesHandler().Saved(w, r)
}

func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	h.PagesHandler().NotFound(w, r)
}

func (h *Handler) Colophon(w http.ResponseWriter, r *http.Request) {
	h.PagesHandler().Colophon(w, r)
}

func (h *Handler) Contribute(w http.ResponseWriter, r *http.Request) {
	h.PagesHandler().Contribute(w, r)
}

func (h *Handler) DownloadTemplate(w http.ResponseWriter, r *http.Request) {
	h.PagesHandler().DownloadTemplate(w, r)
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	h.PagesHandler().Search(w, r)
}

func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	h.PagesHandler().Stats(w, r)
}

func (h *Handler) Contact(w http.ResponseWriter, r *http.Request) {
	h.PagesHandler().Contact(w, r)
}

func validEmail(s string) bool {
	return pages.ValidEmail(s)
}

func (h *Handler) ContactPage(w http.ResponseWriter, r *http.Request) {
	h.PagesHandler().ContactPage(w, r)
}
