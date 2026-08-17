package handler

import (
	"net/http"

	"daemontalk/internal/github"
	"daemontalk/internal/i18n"
	"daemontalk/internal/project"
	"daemontalk/web/templates"
)

func (h *Handler) Behind(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)

	featured := make([]project.Project, 0, 4)
	for _, p := range h.AllProjects {
		if p.Featured {
			featured = append(featured, p)
		}
	}

	ghStats := github.Fetch("dafagareth", h.GitHubToken)

	meta := templates.PageMeta{
		Description: "Who writes this blog, what I'm building, and how this website works.",
	}
	h.Render(w, r, templates.Layout(ui, lang, "behind", r.URL.Path, meta, templates.Behind(ui, lang, featured, ghStats)))
}

