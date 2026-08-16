package post

import (
	"bytes"
	"fmt"
	"html"
	"strings"
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

	for _, it := range items {
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
      </div>`, urlEsc, altEsc))

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
    <span class="fig-count">%d</span>
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
