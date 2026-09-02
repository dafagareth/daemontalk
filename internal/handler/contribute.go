package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"daemontalk/internal/i18n"
	"daemontalk/internal/post"
	"daemontalk/web/templates"
)

func (h *Handler) loadContributeSection(section, lang string) template.HTML {
	var filename string
	if lang == "id" {
		filename = h.getContentPath(filepath.Join("contribute", section+".id.md"))
	} else {
		filename = h.getContentPath(filepath.Join("contribute", section+".md"))
	}

	body, err := post.LoadBody(filename)
	if err != nil {
		defaultFile := h.getContentPath("contribute.md")
		if lang == "id" {
			defaultFile = h.getContentPath("contribute.id.md")
		}
		body, _ = post.LoadBody(defaultFile)
	}
	return body
}

// Contribute handles GET /contribute for both Web and CLI clients.
func (h *Handler) Contribute(w http.ResponseWriter, r *http.Request) {
	if IsCLIRequest(r) {
		h.cliContribute(w, r)
		return
	}

	lang := langFromRequest(r)
	ui := i18n.Get(lang)

	sections := templates.ContributeSections{
		Dispatches:   h.loadContributeSection("dispatches", lang),
		Engine:       h.loadContributeSection("engine", lang),
		Corrections:  h.loadContributeSection("corrections", lang),
		I18n:         h.loadContributeSection("i18n", lang),
		Contributors: post.GetAllContributors(h.AllPosts()),
	}

	title := "Contributor Guide · daemontalk"
	if lang == "id" {
		title = "Panduan Kontribusi · daemontalk"
	}

	h.Render(w, r, templates.Layout(ui, lang, title, r.URL.Path, templates.PageMeta{
		Description: "Editorial standards and submission guide for daemontalk technical writers and contributors.",
	}, templates.ContributePage(ui, sections, lang)))
}

// DownloadTemplate handles GET /download/daemontalk-template.md, /daemontalk-template.md, and legacy /template.md.
func (h *Handler) DownloadTemplate(w http.ResponseWriter, r *http.Request) {
	filename := h.getContentPath("daemontalk-template.md")
	data, err := os.ReadFile(filename)
	if err != nil {
		filename = h.getContentPath("template.md")
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

// cliContribute outputs the contributor guide in ANSI plain-text.
func (h *Handler) cliContribute(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	filename := h.getContentPath(filepath.Join("contribute", "dispatches.md"))
	if lang == "id" {
		filename = h.getContentPath(filepath.Join("contribute", "dispatches.id.md"))
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		defaultFile := h.getContentPath("contribute.md")
		data, err = os.ReadFile(defaultFile)
		if err != nil {
			http.Error(w, "Contributor guide not found", http.StatusNotFound)
			return
		}
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "\033[1;36m=== DAEMONTALK CONTRIBUTOR GUIDE ===\033[0m")
	fmt.Fprintln(w)
	fmt.Fprint(w, string(data))
	fmt.Fprintln(w)
	fmt.Fprintln(w, "\033[90m------------------------------------------------------------\033[0m")
	fmt.Fprintln(w, "\033[1;33mDownload template: \033[0mcurl -s https://daemontalk.com/daemontalk-template.md -o daemontalk-template.md")
	fmt.Fprintln(w, "\033[1;33mSubmit PR:         \033[0mhttps://github.com/dafagareth/daemontalk")
}
