package post

import (
	"bytes"
	"fmt"
	"html"
	"strings"
)

type refItem struct {
	Title     string
	Author    string
	Year      string
	Publisher string
	URL       string
	Note      string
}

func renderReferencesHTML(rawContent string) string {
	lines := strings.Split(rawContent, "\n")
	var items []refItem
	var current refItem

	flushCurrent := func() {
		if current.Title != "" || current.URL != "" || current.Author != "" {
			items = append(items, current)
			current = refItem{}
		}
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			flushCurrent()
			trimmed = strings.TrimSpace(trimmed[2:])
		}

		if strings.Contains(trimmed, "|") && !strings.Contains(trimmed, ": ") {
			flushCurrent()
			parts := strings.Split(trimmed, "|")
			var it refItem
			if len(parts) >= 1 {
				it.Author = strings.TrimSpace(parts[0])
			}
			if len(parts) >= 2 {
				it.Title = strings.TrimSpace(parts[1])
			}
			if len(parts) >= 3 {
				p := strings.TrimSpace(parts[2])
				if strings.HasPrefix(p, "http") {
					it.URL = p
				} else {
					it.Publisher = p
				}
			}
			if len(parts) >= 4 && it.URL == "" {
				it.URL = strings.TrimSpace(parts[3])
			}
			items = append(items, it)
			continue
		}

		parts := strings.SplitN(trimmed, ":", 2)
		if len(parts) == 2 {
			k := strings.ToLower(strings.TrimSpace(parts[0]))
			v := strings.TrimSpace(parts[1])
			switch k {
			case "title":
				current.Title = v
			case "author", "authors":
				current.Author = v
			case "year", "date":
				current.Year = v
			case "publisher", "source", "venue":
				current.Publisher = v
			case "url", "link", "doi":
				current.URL = v
			case "note":
				current.Note = v
			}
		}
	}
	flushCurrent()

	if len(items) == 0 {
		return ""
	}

	var buf bytes.Buffer
	buf.WriteString("\n<div class=\"post-references my-8 not-prose\">\n")
	buf.WriteString("  <ol class=\"space-y-2.5 my-4 p-0 list-decimal list-inside text-[0.92em] text-muted\">\n")

	for _, it := range items {
		titleEsc := html.EscapeString(it.Title)
		if titleEsc == "" {
			titleEsc = html.EscapeString(it.URL)
		}

		buf.WriteString("    <li class=\"leading-relaxed pl-1\">\n")

		if it.Author != "" {
			buf.WriteString(fmt.Sprintf("      <strong class=\"text-text font-semibold\">%s</strong>", html.EscapeString(it.Author)))
			if it.Year != "" {
				buf.WriteString(fmt.Sprintf(" (%s). ", html.EscapeString(it.Year)))
			} else {
				buf.WriteString(". ")
			}
		}

		if it.Title != "" {
			buf.WriteString(fmt.Sprintf("<em class=\"italic text-text\">%s</em>. ", titleEsc))
		}

		if it.Publisher != "" {
			buf.WriteString(fmt.Sprintf("%s. ", html.EscapeString(it.Publisher)))
		}

		if it.Note != "" {
			buf.WriteString(fmt.Sprintf("<span class=\"text-xs text-text/80 italic\">(%s) </span>", html.EscapeString(it.Note)))
		}

		if it.URL != "" {
			displayURL := it.URL
			if strings.HasPrefix(displayURL, "https://") {
				displayURL = strings.TrimPrefix(displayURL, "https://")
			} else if strings.HasPrefix(displayURL, "http://") {
				displayURL = strings.TrimPrefix(displayURL, "http://")
			}
			buf.WriteString(fmt.Sprintf("<a href=\"%s\" target=\"_blank\" rel=\"noopener noreferrer\" class=\"text-link underline hover:text-text font-mono text-[0.9em]\">%s</a>\n", html.EscapeString(it.URL), html.EscapeString(displayURL)))
		}

		buf.WriteString("    </li>\n")
	}

	buf.WriteString("  </ol>\n</div>\n")
	return buf.String()
}
