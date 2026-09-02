package handler

import (
	"net/http"
	"os"
	"regexp"
	"strings"

	"daemontalk/internal/i18n"
	"daemontalk/web/templates"
)

var (
	reReleaseHeader = regexp.MustCompile(`^(v[\d\.]+)\s*\((.*?)\)\s*[·-]\s*(.*)`)
	reItem          = regexp.MustCompile(`^-\s*(?:\*\*(.*?)\*\*:?\s*)?(.*)`)
)

func parseChangelogMarkdown(content string) []templates.ChangelogRelease {
	var releases []templates.ChangelogRelease
	sections := strings.Split(content, "\n### ")
	for _, sec := range sections {
		sec = strings.TrimPrefix(sec, "### ")
		sec = strings.TrimSpace(sec)
		if sec == "" {
			continue
		}
		lines := strings.Split(sec, "\n")
		headerLine := strings.TrimSpace(lines[0])
		m := reReleaseHeader.FindStringSubmatch(headerLine)
		var rel templates.ChangelogRelease
		if len(m) >= 4 {
			rel.Version = m[1]
			rel.Date = m[2]
			rel.Title = m[3]
		} else {
			rel.Version = headerLine
		}

		for _, line := range lines[1:] {
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "-") {
				continue
			}
			im := reItem.FindStringSubmatch(line)
			if len(im) >= 3 {
				scope := strings.TrimSuffix(strings.TrimSpace(im[1]), ":")
				desc := strings.TrimSpace(im[2])
				if scope == "" && desc == "" {
					desc = strings.TrimPrefix(line, "- ")
				}
				rel.Items = append(rel.Items, templates.ChangelogItem{
					Scope: scope,
					Desc:  desc,
				})
			} else {
				rel.Items = append(rel.Items, templates.ChangelogItem{
					Desc: strings.TrimPrefix(line, "- "),
				})
			}
		}
		if rel.Version != "" {
			releases = append(releases, rel)
		}
	}
	return releases
}

func (h *Handler) Changelog(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)

	filename := h.getContentPath("changelog.md")
	if lang == "id" {
		filename = h.getContentPath("changelog.id.md")
	}

	b, err := os.ReadFile(filename)
	var releases []templates.ChangelogRelease
	if err == nil {
		releases = parseChangelogMarkdown(string(b))
	}

	h.Render(w, r, templates.Layout(ui, lang, "changelog", r.URL.Path, templates.PageMeta{
		Description: "A running log of features and changes shipped to daemontalk.com.",
	}, templates.ChangelogPage(ui, lang, releases)))
}
