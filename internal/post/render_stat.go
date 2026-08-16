package post

import (
	"bytes"
	"fmt"
	"html"
	"strings"
)

type statItem struct {
	Value       string
	Label       string
	Description string
}

func renderStatHTML(rawContent string) string {
	lines := strings.Split(rawContent, "\n")
	var items []statItem
	var current statItem

	flush := func() {
		if current.Value != "" || current.Label != "" {
			items = append(items, current)
			current = statItem{}
		}
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			flush()
			trimmed = strings.TrimSpace(trimmed[2:])
		}

		if strings.Contains(trimmed, "|") && !strings.Contains(trimmed, ": ") {
			flush()
			parts := strings.Split(trimmed, "|")
			var it statItem
			if len(parts) >= 1 {
				it.Value = strings.TrimSpace(parts[0])
			}
			if len(parts) >= 2 {
				it.Label = strings.TrimSpace(parts[1])
			}
			if len(parts) >= 3 {
				it.Description = strings.TrimSpace(parts[2])
			}
			items = append(items, it)
			continue
		}

		parts := strings.SplitN(trimmed, ":", 2)
		if len(parts) == 2 {
			k := strings.ToLower(strings.TrimSpace(parts[0]))
			v := strings.TrimSpace(parts[1])
			v = strings.Trim(v, `"'`)
			switch k {
			case "value", "stat", "metric", "num", "number":
				current.Value = v
			case "label", "title", "name":
				current.Label = v
			case "description", "desc", "sub", "detail":
				current.Description = v
			}
		}
	}
	flush()

	if len(items) == 0 {
		return ""
	}

	cols := "grid-cols-1"
	if len(items) == 2 {
		cols = "grid-cols-1 sm:grid-cols-2"
	} else if len(items) >= 3 {
		cols = "grid-cols-1 sm:grid-cols-2 lg:grid-cols-3"
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("\n<div class=\"post-stat-grid my-8 not-prose grid %s gap-4\">\n", cols))
	for _, it := range items {
		buf.WriteString(fmt.Sprintf(`  <div class="p-4 bg-surface border border-border flex flex-col justify-between rounded-none">
    <div class="stat-value font-mono text-[1.85em] font-black text-link tracking-tight">%s</div>
    <div class="stat-label mt-2 font-mono text-[0.8em] font-bold uppercase tracking-wider text-text">%s</div>
`, html.EscapeString(it.Value), html.EscapeString(it.Label)))
		if it.Description != "" {
			buf.WriteString(fmt.Sprintf(`    <div class="stat-desc text-[0.8em] text-muted mt-0.5">%s</div>
`, html.EscapeString(it.Description)))
		}
		buf.WriteString("  </div>\n")
	}
	buf.WriteString("</div>\n")
	return buf.String()
}
