package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"daemontalk/internal/i18n"
	"daemontalk/internal/post"
	"daemontalk/internal/project"
	"daemontalk/web/templates"
)

func (h *Handler) ProjectDetail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	lang := langFromRequest(r)
	ui := i18n.Get(lang)

	var found *project.Project
	for i := range h.AllProjects {
		if h.AllProjects[i].Slug == slug {
			cp := h.AllProjects[i]
			found = &cp
			break
		}
	}
	if found == nil {
		h.NotFound(w, r)
		return
	}

	// Load markdown body from content/projects/{slug}.md (best-effort).
	// For Indonesian, prefer {slug}.id.md and fall back to the English body.
	var toc []post.TOCEntry
	if lang == "id" {
		if body, t, err := post.LoadBodyWithTOC(h.getContentPath("projects/" + slug + ".id.md")); err == nil {
			found.Body = body
			toc = t
		}
	}
	if found.Body == "" {
		if body, t, err := post.LoadBodyWithTOC(h.getContentPath("projects/" + slug + ".md")); err == nil {
			found.Body = body
			toc = t
		}
	}

	desc := found.Description
	if lang == "id" && found.DescriptionID != "" {
		desc = found.DescriptionID
	}

	h.Render(w, r, templates.Layout(ui, lang, found.Name+" · daemontalk", r.URL.Path, templates.PageMeta{
		Description: desc,
	}, templates.ProjectDetailPage(ui, *found, lang, toc)))
}

