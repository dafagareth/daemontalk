package post

import (
	"fmt"
	"html"
	"strings"
)

type linkCard struct {
	URL         string
	Title       string
	Description string
	Site        string
	Author      string
}

func renderLinkHTML(rawContent string) string {
	lines := strings.Split(rawContent, "\n")
	var card linkCard
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed == "" {
			continue
		}
		parts := strings.SplitN(trimmed, ":", 2)
		if len(parts) == 2 {
			k := strings.ToLower(strings.TrimSpace(parts[0]))
			v := strings.TrimSpace(parts[1])
			switch k {
			case "url", "link", "href":
				card.URL = v
			case "title", "name":
				card.Title = v
			case "description", "desc", "summary":
				card.Description = v
			case "site", "domain", "source":
				card.Site = v
			case "author":
				card.Author = v
			}
		}
	}

	if card.URL == "" && card.Title == "" {
		return ""
	}
	if card.Title == "" {
		card.Title = card.URL
	}
	if card.Site == "" && strings.Contains(card.URL, "://") {
		u := strings.Split(card.URL, "://")[1]
		card.Site = strings.Split(u, "/")[0]
	}

	descHTML := ""
	if card.Description != "" {
		descHTML = fmt.Sprintf(`<div class="text-[0.88em] text-muted mt-1 leading-relaxed">%s</div>`, html.EscapeString(card.Description))
	}

	siteHTML := ""
	if card.Site != "" {
		siteHTML = fmt.Sprintf(`<div class="text-[0.8em] font-mono text-muted/70 mt-1.5">%s</div>`, html.EscapeString(card.Site))
	}

	return fmt.Sprintf(`
<div class="post-link-preview my-5 not-prose border border-border p-4 rounded-none hover:bg-hover transition-colors">
  <a href="%s" target="_blank" rel="noopener noreferrer" class="group block no-underline text-inherit hover:no-underline">
    <div class="link-title font-semibold text-text group-hover:text-link text-[0.95em] leading-snug no-underline flex items-center gap-2">
      %s
      `+GetIcon(IconExternalLink, "w-4 h-4 text-muted group-hover:text-link transition-colors")+`
    </div>
    %s
    %s
  </a>
</div>
`, html.EscapeString(safeURL(card.URL)), html.EscapeString(card.Title), descHTML, siteHTML)
}
