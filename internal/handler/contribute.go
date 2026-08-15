package handler

import (
	"log"
	"fmt"
	"net/http"
	"os"

	"daemontalk/internal/i18n"
	"daemontalk/internal/post"
	"daemontalk/web/templates"
)

// Contribute handles GET /contribute for both Web and CLI clients.
func (h *Handler) Contribute(w http.ResponseWriter, r *http.Request) {
	if IsCLIRequest(r) {
		h.CLIContribute(w, r)
		return
	}

	lang := langFromRequest(r)
	ui := i18n.Get(lang)

	filename := h.getContentPath("contribute.md")
	if lang == "id" {
		filename = h.getContentPath("contribute.id.md")
	}
	body, _ := post.LoadBody(filename)

	title := "Contributor Guide · daemontalk"
	if lang == "id" {
		title = "Panduan Kontribusi · daemontalk"
	}

	err := templates.Layout(ui, lang, title, r.URL.Path, templates.PageMeta{
		Description: "Editorial standards and submission guide for daemontalk technical writers and contributors.",
	}, templates.ContributePage(body, lang)).Render(r.Context(), w)
	if err != nil {
		log.Printf("render error: %v", err)
	}
}

// DownloadTemplate handles GET /download/template.md and /template.md to serve the markdown template file.
func (h *Handler) DownloadTemplate(w http.ResponseWriter, r *http.Request) {
	filename := h.getContentPath("template.md")
	data, err := os.ReadFile(filename)
	if err != nil {
		http.Error(w, "Template file not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=\"template.md\"")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

// CLIContribute outputs the contributor guide in ANSI plain-text.
func (h *Handler) CLIContribute(w http.ResponseWriter, r *http.Request) {
	color := r.URL.Query().Get("nocolor") != "1"

	data, err := os.ReadFile(h.getContentPath("contribute.md"))
	if err != nil {
		http.Error(w, "Contributor guide not found", http.StatusNotFound)
		return
	}

	var b stringsBuilderWrapper
	b.WriteString(fmt.Sprintf("%s%s[ DAEMONTALK COMMUNITY CONTRIBUTOR GUIDE ]%s\n", ansiGreen, ansiBold, ansiReset))
	b.WriteString(fmt.Sprintf("%s══════════════════════════════════════════════════════════════════════════%s\n\n", ansiDim, ansiReset))
	b.WriteString(string(data))
	b.WriteString(fmt.Sprintf("\n\n%sSubmit PRs at: https://github.com/dafagareth/daemontalk%s\n", ansiCyan, ansiReset))

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(applyColors(b.String(), color)))
}

type stringsBuilderWrapper struct {
	content string
}

func (s *stringsBuilderWrapper) WriteString(str string) {
	s.content += str
}

func (s *stringsBuilderWrapper) String() string {
	return s.content
}
