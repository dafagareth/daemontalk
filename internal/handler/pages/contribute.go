package pages

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"daemontalk/internal/handler/common"
	"daemontalk/internal/i18n"
	"daemontalk/internal/post"
	"daemontalk/web/templates/layouts"
	pagestmpl "daemontalk/web/templates/pages"
	"daemontalk/web/templates/shared"
)

func (h *Handler) loadContributeSection(section, lang string) template.HTML {
	filename := common.GetContentPath(h.ContentDir, filepath.Join("contribute", section+".md"))
	body, err := post.LoadBody(filename)
	if err != nil {
		defaultFile := common.GetContentPath(h.ContentDir, "contribute.md")
		body, _ = post.LoadBody(defaultFile)
	}
	return body
}

func (h *Handler) Contribute(w http.ResponseWriter, r *http.Request) {
	lang := common.LangFromRequest(r)
	ui := i18n.Get(lang)

	sections := pagestmpl.ContributeSections{
		Dispatches:   h.loadContributeSection("dispatches", lang),
		Engine:       h.loadContributeSection("engine", lang),
		Corrections:  h.loadContributeSection("corrections", lang),
		I18n:         h.loadContributeSection("i18n", lang),
		Contributors: post.GetAllContributors(h.GetAllPosts()),
	}

	title := "Contributor Guide · daemontalk"

	common.Render(w, r, layouts.Layout(ui, lang, title, r.URL.Path, shared.PageMeta{
		Description: "Editorial standards and submission guide for daemontalk technical writers and contributors.",
	}, pagestmpl.ContributePage(ui, sections, lang)))
}

func (h *Handler) DownloadTemplate(w http.ResponseWriter, r *http.Request) {
	filename := common.GetContentPath(h.ContentDir, "daemontalk-template.md")
	data, err := os.ReadFile(filename)
	if err != nil {
		filename = common.GetContentPath(h.ContentDir, "template.md")
		data, err = os.ReadFile(filename)
		if err != nil {
			http.Error(w, "Template file not found", http.StatusNotFound)
			return
		}
	}

	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=\"daemontalk-template.md\"")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
