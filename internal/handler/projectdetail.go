package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) ProjectDetail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	for _, p := range h.AllProjects {
		if p.Slug == slug {
			if p.RepoURL != "" {
				http.Redirect(w, r, p.RepoURL, http.StatusMovedPermanently)
				return
			}
		}
	}

	http.Redirect(w, r, "/colophon#projects", http.StatusMovedPermanently)
}
