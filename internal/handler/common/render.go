package common

import (
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"

	"daemontalk/internal/i18n"
	"daemontalk/internal/post"
	"daemontalk/web/templates/layouts"
	"daemontalk/web/templates/shared"
	"github.com/a-h/templ"
)

func LangFromRequest(r *http.Request) string {
	return "en"
}

func URLPrefix(lang string) string {
	return ""
}

func Render(w http.ResponseWriter, r *http.Request, c templ.Component) {
	if err := c.Render(r.Context(), w); err != nil {
		slog.Error("render component failed", "error", err, "path", r.URL.Path, "method", r.Method)
	}
}

func GetContentPath(contentDir, subpath string) string {
	if contentDir == "" {
		contentDir = "content"
	}
	return filepath.Join(contentDir, subpath)
}

func RenderMarkdownPage(
	w http.ResponseWriter, r *http.Request,
	contentDir, contentKey, title string, meta shared.PageMeta,
	render func(i18n.UI, template.HTML, string) templ.Component,
) {
	lang := "en"
	ui := i18n.Get(lang)

	filename := GetContentPath(contentDir, contentKey+".md")
	body, err := post.LoadBody(filename)
	if err != nil {
		slog.Warn("load markdown page failed", "file", filename, "error", err)
	}

	Render(w, r, layouts.Layout(ui, lang, title, r.URL.Path, meta, render(ui, body, lang)))
}

func AbsoluteURL(r *http.Request, path string) string {
	if path == "" {
		return ""
	}
	if len(path) >= 4 && path[:4] == "http" {
		return path
	}
	if path[0] != '/' {
		path = "/" + path
	}

	base := shared.SiteBaseURL
	if r != nil {
		proto := "http"
		if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
			proto = "https"
		}
		host := r.Header.Get("X-Forwarded-Host")
		if host == "" {
			host = r.Host
		}
		if host != "" {
			base = proto + "://" + host
		}
	}
	return strings.TrimSuffix(base, "/") + path
}

func StripCRLF(s string) string {
	return strings.NewReplacer("\r", "", "\n", "").Replace(s)
}
