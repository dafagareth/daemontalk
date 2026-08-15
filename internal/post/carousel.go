package post

import (
	"bytes"
	"fmt"
	"html"
	"regexp"
	"strings"
)

var (
	reCarouselFenced = regexp.MustCompile("(?s)```(?:carousel|slider)\\s*\\n(.*?)\\n```")
	reGalleryFenced  = regexp.MustCompile("(?s)```gallery\\s*\\n(.*?)\\n```")
	reFAQFenced      = regexp.MustCompile("(?s)```faq\\s*\\n(.*?)\\n```")
	reAuthorFenced   = regexp.MustCompile("(?s)```author\\s*\\n(.*?)\\n```")
	reMarkdownImage  = regexp.MustCompile(`!\[([^\]]*)\]\(([^)"\s]+)(?:\s+(?:"([^"]*)"|'([^']*)'))?\)`)
)

type mediaItem struct {
	Alt     string
	URL     string
	Caption string
}

func extractMediaItems(content string) []mediaItem {
	var items []mediaItem
	matches := reMarkdownImage.FindAllStringSubmatch(content, -1)
	for _, m := range matches {
		alt := strings.TrimSpace(m[1])
		url := strings.TrimSpace(m[2])
		caption := ""
		if len(m) > 3 && m[3] != "" {
			caption = strings.TrimSpace(m[3])
		} else if len(m) > 4 && m[4] != "" {
			caption = strings.TrimSpace(m[4])
		}
		if url != "" {
			items = append(items, mediaItem{
				Alt:     alt,
				URL:     url,
				Caption: caption,
			})
		}
	}
	return items
}

func renderCarouselHTML(items []mediaItem) string {
	if len(items) == 0 {
		return ""
	}

	total := len(items)
	var buf bytes.Buffer

	buf.WriteString("\n<div class=\"post-carousel-wrap my-8 not-prose\">\n")
	buf.WriteString("  <div class=\"carousel-track flex gap-6 overflow-x-auto snap-x snap-mandatory pb-2 scrollbar-none scroll-smooth\">\n")

	for i, it := range items {
		altEsc := html.EscapeString(it.Alt)
		urlEsc := html.EscapeString(it.URL)
		capEsc := html.EscapeString(it.Caption)

		captionText := capEsc
		if captionText == "" && altEsc != "" {
			captionText = altEsc
		}

		buf.WriteString(fmt.Sprintf(`    <figure class="snap-start shrink-0 w-full sm:w-[90%%] md:w-[85%%] flex flex-col items-center m-0 bg-transparent">
      <div class="relative w-full aspect-[16/9] overflow-hidden flex items-center justify-center bg-transparent border-0">
        <img src="%s" alt="%s" class="w-full h-full object-contain rounded-none border-0 !border-none !m-0 select-none" loading="lazy" />
        <span class="absolute top-2.5 right-2.5 bg-black/80 text-white font-mono text-[11px] px-2 py-0.5 border border-white/20 select-none z-10">%d / %d</span>
      </div>`, urlEsc, altEsc, i+1, total))

		if captionText != "" {
			buf.WriteString(fmt.Sprintf(`
      <figcaption class="mt-2 text-xs font-mono text-muted text-center leading-relaxed">
        %s
      </figcaption>`, captionText))
		}

		buf.WriteString("\n    </figure>\n")
	}

	buf.WriteString("  </div>\n")
	buf.WriteString(fmt.Sprintf(`  <div class="flex items-center justify-between text-xs font-mono text-muted mt-2 px-0.5 select-none">
    <span>%d Gambar</span>
    <div class="flex items-center gap-1.5">
      <button
        type="button"
        onclick="var t=this.closest('.post-carousel-wrap').querySelector('.carousel-track'); t.scrollBy({left: -t.clientWidth*0.9, behavior: 'smooth'});"
        aria-label="Previous slide"
        class="w-8 h-8 flex items-center justify-center border border-border bg-surface text-text hover:bg-hover hover:border-[var(--c-link)] transition-colors cursor-pointer"
      >
        <svg class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" /></svg>
      </button>
      <button
        type="button"
        onclick="var t=this.closest('.post-carousel-wrap').querySelector('.carousel-track'); t.scrollBy({left: t.clientWidth*0.9, behavior: 'smooth'});"
        aria-label="Next slide"
        class="w-8 h-8 flex items-center justify-center border border-border bg-surface text-text hover:bg-hover hover:border-[var(--c-link)] transition-colors cursor-pointer"
      >
        <svg class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" /></svg>
      </button>
    </div>
  </div>
</div>
`, total))

	return buf.String()
}

func renderGalleryHTML(items []mediaItem) string {
	if len(items) == 0 {
		return ""
	}

	total := len(items)
	gridCols := "grid-cols-1 sm:grid-cols-2"
	if total == 3 {
		gridCols = "grid-cols-1 sm:grid-cols-3"
	} else if total >= 4 {
		gridCols = "grid-cols-1 sm:grid-cols-2 lg:grid-cols-4"
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("\n<div class=\"post-gallery-wrap my-8 not-prose grid %s gap-6 items-start\">\n", gridCols))

	for _, it := range items {
		altEsc := html.EscapeString(it.Alt)
		urlEsc := html.EscapeString(it.URL)
		capEsc := html.EscapeString(it.Caption)

		captionText := capEsc
		if captionText == "" && altEsc != "" {
			captionText = altEsc
		}

		buf.WriteString(fmt.Sprintf(`  <figure class="flex flex-col items-center m-0 bg-transparent">
    <div class="relative w-full aspect-[16/9] overflow-hidden flex items-center justify-center bg-transparent border-0">
      <img src="%s" alt="%s" class="w-full h-full object-contain rounded-none border-0 !border-none !m-0" loading="lazy" />
    </div>`, urlEsc, altEsc))

		if captionText != "" {
			buf.WriteString(fmt.Sprintf(`
    <figcaption class="mt-2 text-xs font-mono text-muted text-center leading-relaxed">
      %s
    </figcaption>`, captionText))
		}

		buf.WriteString("\n  </figure>\n")
	}

	buf.WriteString("</div>\n")
	return buf.String()
}

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
    <summary class="p-4 text-sm font-semibold text-text cursor-pointer select-none list-none flex items-center justify-between gap-4 hover:bg-hover transition-colors">
      <span>%s</span>
      <span class="text-base font-mono text-muted group-open:rotate-45 transition-transform shrink-0">+</span>
    </summary>
    <div class="px-4 pb-4 pt-1 text-sm text-muted leading-relaxed border-t border-border/40 bg-bg/40">
      %s
    </div>
  </details>
`, qEsc, aEsc))
	}

	buf.WriteString("</div>\n")
	return buf.String()
}

func renderAuthorHTML(rawContent string) string {
	meta := make(map[string]string)
	lines := strings.Split(rawContent, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			k := strings.ToLower(strings.TrimSpace(parts[0]))
			v := strings.TrimSpace(parts[1])
			meta[k] = v
		}
	}

	name := meta["name"]
	if name == "" {
		name = "Dafa Gareth"
	}
	role := meta["role"]
	if role == "" {
		role = "Software Engineer"
	}
	avatar := meta["avatar"]
	if avatar == "" {
		avatar = "/static/logo/logo-dark.png"
	}
	bio := meta["bio"]
	if bio == "" {
		bio = "Software Engineer yang berfokus pada sistem terdistribusi, rekayasa kernel Linux, dan optimasi performa backend Go/Rust."
	}
	github := meta["github"]
	email := meta["email"]
	website := meta["website"]

	var buf bytes.Buffer
	buf.WriteString("\n<div class=\"post-author-card my-8 not-prose p-5 border border-border bg-surface rounded-none flex flex-col sm:flex-row gap-5 items-start\">\n")
	buf.WriteString(fmt.Sprintf(`  <img src="%s" alt="%s" class="w-14 h-14 rounded-none object-cover border border-border bg-black/10 shrink-0" loading="lazy" />
  <div class="flex-1 min-w-0">
    <div class="flex flex-wrap items-baseline gap-2 mb-1.5">
      <h3 class="text-sm font-bold text-text uppercase tracking-wider m-0">%s</h3>
      <span class="text-xs font-mono text-muted">%s</span>
    </div>
    <p class="text-xs text-muted leading-relaxed mb-3 m-0">%s</p>
    <div class="flex flex-wrap items-center gap-4 text-xs font-mono">
`, html.EscapeString(avatar), html.EscapeString(name), html.EscapeString(name), html.EscapeString(role), html.EscapeString(bio)))

	if github != "" {
		buf.WriteString(fmt.Sprintf(`      <a href="%s" target="_blank" rel="noopener" class="text-link hover:underline inline-flex items-center gap-1">
        <span>GitHub</span>
      </a>
`, html.EscapeString(github)))
	}
	if email != "" {
		buf.WriteString(fmt.Sprintf(`      <a href="mailto:%s" class="text-muted hover:text-text transition-colors">
        <span>%s</span>
      </a>
`, html.EscapeString(email), html.EscapeString(email)))
	}
	if website != "" {
		buf.WriteString(fmt.Sprintf(`      <a href="%s" target="_blank" rel="noopener" class="text-link hover:underline">
        <span>Website</span>
      </a>
`, html.EscapeString(website)))
	}

	buf.WriteString(`    </div>
  </div>
</div>
`)
	return buf.String()
}

func formatSimpleInline(s string) string {
	// Format inline `code`
	reCode := regexp.MustCompile("`([^`]+)`")
	s = reCode.ReplaceAllString(s, `<code class="font-mono text-xs px-1.5 py-0.5 bg-chip border border-border text-text">$1</code>`)

	// Format inline **bold**
	reBold := regexp.MustCompile(`\*\*([^*]+)\*\*`)
	s = reBold.ReplaceAllString(s, `<strong class="font-semibold text-text">$1</strong>`)

	// Format inline *italic*
	reItalic := regexp.MustCompile(`\*([^*]+)\*`)
	s = reItalic.ReplaceAllString(s, `<em class="italic text-muted">$1</em>`)

	return s
}

// preprocessMarkdown handles custom media blocks such as ```carousel, ```gallery, ```faq, and ```author.
func preprocessMarkdown(src []byte) []byte {
	srcStr := string(src)

	// Replace ```carousel ... ```
	srcStr = reCarouselFenced.ReplaceAllStringFunc(srcStr, func(m string) string {
		sub := reCarouselFenced.FindStringSubmatch(m)
		if len(sub) < 2 {
			return m
		}
		items := extractMediaItems(sub[1])
		if len(items) == 0 {
			return m
		}
		return renderCarouselHTML(items)
	})

	// Replace ```gallery ... ```
	srcStr = reGalleryFenced.ReplaceAllStringFunc(srcStr, func(m string) string {
		sub := reGalleryFenced.FindStringSubmatch(m)
		if len(sub) < 2 {
			return m
		}
		items := extractMediaItems(sub[1])
		if len(items) == 0 {
			return m
		}
		return renderGalleryHTML(items)
	})

	// Replace ```faq ... ```
	srcStr = reFAQFenced.ReplaceAllStringFunc(srcStr, func(m string) string {
		sub := reFAQFenced.FindStringSubmatch(m)
		if len(sub) < 2 {
			return m
		}
		return renderFAQHTML(sub[1])
	})

	// Replace ```author ... ```
	srcStr = reAuthorFenced.ReplaceAllStringFunc(srcStr, func(m string) string {
		sub := reAuthorFenced.FindStringSubmatch(m)
		if len(sub) < 2 {
			return m
		}
		return renderAuthorHTML(sub[1])
	})

	return []byte(srcStr)
}
