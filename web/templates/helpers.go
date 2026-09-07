package templates

import (
	"fmt"
	"html"
	"html/template"
	"regexp"
	"strings"

	"daemontalk/internal/post"
)

var (
	cmBoldRe = regexp.MustCompile(`\*\*([^*\n]+)\*\*`)
	cmCodeRe = regexp.MustCompile("`([^`\n]+)`")
)

func renderCommentBody(body string) template.HTML {
	escaped := html.EscapeString(body)
	escaped = cmBoldRe.ReplaceAllString(escaped, "<strong>$1</strong>")
	escaped = cmCodeRe.ReplaceAllString(escaped, `<code class="font-mono text-xs bg-[var(--c-chip)] px-1 rounded">$1</code>`)
	escaped = strings.ReplaceAll(escaped, "\n", "<br>")
	return template.HTML(escaped)
}

func highlightHTML(text, query string) template.HTML {
	if query == "" || text == "" || len(query) < 2 {
		return template.HTML(html.EscapeString(text))
	}
	escaped := html.EscapeString(text)
	escapedQ := html.EscapeString(query)
	lowerEscaped := strings.ToLower(escaped)
	lowerQ := strings.ToLower(escapedQ)

	qLen := len(lowerQ)
	if qLen == 0 {
		return template.HTML(escaped)
	}

	idx := strings.Index(lowerEscaped, lowerQ)
	if idx == -1 {
		return template.HTML(escaped)
	}

	var b strings.Builder
	b.Grow(len(escaped) + 64)
	lastIdx := 0

	for idx != -1 {
		b.WriteString(escaped[lastIdx:idx])
		b.WriteString(`<mark class="bg-yellow-100 dark:bg-yellow-900/40 text-[var(--c-text)] rounded px-0.5">`)
		b.WriteString(escaped[idx : idx+qLen])
		b.WriteString(`</mark>`)
		lastIdx = idx + qLen
		nextMatch := strings.Index(lowerEscaped[lastIdx:], lowerQ)
		if nextMatch == -1 {
			break
		}
		idx = lastIdx + nextMatch
	}
	b.WriteString(escaped[lastIdx:])
	return template.HTML(b.String())
}

func generateDigestText(posts []post.Post) string {
	var sb strings.Builder
	sb.WriteString("# ⚡ DaemonTalk Weekly Tech Digest\n\n")
	sb.WriteString("> Curated deep-dives in modern computing, AI, and universal technology.\n\n")
	sb.WriteString("## 📚 Featured Tech Dispatches This Week\n\n")

	base := strings.TrimSuffix(SiteBaseURL, "/")
	hostDomain := strings.TrimPrefix(strings.TrimPrefix(base, "https://"), "http://")
	count := 0
	for _, p := range posts {
		if p.Draft {
			continue
		}
		sb.WriteString("### [" + p.Title + "](" + base + "/blog/" + p.Slug + ")\n")
		sb.WriteString(fmt.Sprintf("**Date:** %s  •  **Read Time:** %d min read\n\n", p.Date.Format("02 Jan 2006"), p.ReadTime))
		if p.Description != "" {
			sb.WriteString(p.Description + "\n\n")
		}
		sb.WriteString("[Read full dispatch →](" + base + "/blog/" + p.Slug + ")\n\n---\n\n")
		count++
		if count >= 5 {
			break
		}
	}

	sb.WriteString("## 💻 Read Terminal-First\n\n")
	sb.WriteString("Access DaemonTalk directly in your terminal:\n")
	sb.WriteString("```bash\n# Instant TUI over SSH\nssh ssh.daemontalk.com -p 2222\n\n# Daily dispatch stream\ncurl -sL " + hostDomain + "/daily\n```\n\n")
	sb.WriteString("---\n*Published by [DaemonTalk](" + base + ") · Independent Systems Journalism.*")

	return sb.String()
}

func coverSourceName(rawURL string) string {
	if rawURL == "" {
		return ""
	}
	lower := strings.ToLower(rawURL)
	switch {
	case strings.Contains(lower, "unsplash.com"):
		return "Unsplash"
	case strings.Contains(lower, "github.com"):
		return "GitHub"
	case strings.Contains(lower, "wikimedia.org") || strings.Contains(lower, "wikipedia.org"):
		return "Wikimedia"
	case strings.Contains(lower, "pexels.com"):
		return "Pexels"
	}
	s := strings.TrimPrefix(rawURL, "https://")
	s = strings.TrimPrefix(s, "http://")
	s = strings.TrimPrefix(s, "www.")
	if idx := strings.Index(s, "/"); idx != -1 {
		s = s[:idx]
	}
	if s != "" {
		return s
	}
	return "Source"
}

func bibtexFromMeta(slug, title, date string) string {
	year := "2026"
	if len(date) >= 4 {
		year = date[:4]
	}
	return fmt.Sprintf(`@article{daemontalk_%s,
  title = "{%s}",
  author = {Daemontalk Team},
  year = "%s",
  url = "https://daemontalk.com/blog/%s"
}`, strings.ReplaceAll(slug, "-", "_"), title, year, slug)
}

func demoLabel(url string) string {
	if strings.Contains(url, "/releases") {
		return "Releases"
	}
	return "Demo"
}
