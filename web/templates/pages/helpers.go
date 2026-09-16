package pages

import (
	"fmt"
	"html"
	"html/template"
	"sort"
	"strings"
	"time"

	"daemontalk/internal/post"
)

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

func fmtNumber(n int) string {
	if n >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(n)/1000000)
	}
	if n >= 1000 {
		return fmt.Sprintf("%.1fk", float64(n)/1000)
	}
	return fmt.Sprintf("%d", n)
}

func getDailyBriefingDate(articles []post.Post) time.Time {
	now := time.Now()
	if len(articles) > 0 && articles[0].Date.After(now) {
		return articles[0].Date
	}
	return now
}

func extractTagsFromPosts(posts []post.Post) []string {
	tagSet := make(map[string]bool)
	for _, p := range posts {
		if !p.Draft {
			for _, t := range p.Tags {
				clean := strings.ToLower(strings.TrimSpace(t))
				if clean != "" {
					tagSet[clean] = true
				}
			}
		}
	}
	tags := make([]string, 0, len(tagSet))
	for t := range tagSet {
		tags = append(tags, t)
	}
	sort.Strings(tags)
	return tags
}

func tagPostCount(posts []post.Post, tag string) int {
	cnt := 0
	for _, p := range posts {
		if !p.Draft {
			for _, t := range p.Tags {
				if strings.EqualFold(t, tag) {
					cnt++
					break
				}
			}
		}
	}
	return cnt
}
