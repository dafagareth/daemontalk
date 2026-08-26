package post

import (
	"bytes"
	"fmt"
	"html"
	"strings"
)

type faqItem struct {
	Question string
	Answer   string
}

func renderFAQHTML(rawContent string) string {
	lines := strings.Split(rawContent, "\n")
	var items []faqItem
	var currentQ, currentA string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "Q:") || strings.HasPrefix(trimmed, "q:") {
			if currentQ != "" {
				items = append(items, faqItem{
					Question: currentQ,
					Answer:   strings.TrimSpace(currentA),
				})
				currentA = ""
			}
			currentQ = strings.TrimSpace(trimmed[2:])
		} else if strings.HasPrefix(trimmed, "A:") || strings.HasPrefix(trimmed, "a:") {
			currentA = strings.TrimSpace(trimmed[2:])
		} else if currentQ != "" {
			if currentA != "" {
				currentA += " " + trimmed
			} else {
				currentA = trimmed
			}
		}
	}

	if currentQ != "" {
		items = append(items, faqItem{
			Question: currentQ,
			Answer:   strings.TrimSpace(currentA),
		})
	}

	if len(items) == 0 {
		return ""
	}

	var buf bytes.Buffer
	buf.WriteString("\n<div class=\"post-faq-wrap my-8 not-prose space-y-3\">\n")

	for _, item := range items {
		qEsc := html.EscapeString(item.Question)
		aEsc := html.EscapeString(item.Answer)

		// Basic inline formatting support (`code` and **bold**)
		aEsc = formatSimpleInline(aEsc)

		buf.WriteString(fmt.Sprintf(`  <details class="border border-border bg-surface rounded-none overflow-hidden group">
    <summary class="p-4 text-[0.95em] font-semibold text-text cursor-pointer select-none list-none flex items-center justify-between gap-4 hover:bg-hover transition-colors">
      <span>%s</span>
      `+GetIcon(IconChevronDown, "w-4 h-4 text-muted faq-chevron transition-transform shrink-0")+`
    </summary>
    <div class="px-4 pb-4 pt-1 text-[0.9em] text-muted leading-relaxed border-t border-border/40 bg-bg/40">
      %s
    </div>
  </details>
`, qEsc, aEsc))
	}

	buf.WriteString("</div>\n")
	return buf.String()
}
