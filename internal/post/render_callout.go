package post

import (
	"fmt"
	"html"
	"strings"
)

func renderCalloutHTML(kind, rawBody string) string {
	kind = strings.ToUpper(strings.TrimSpace(kind))
	if kind == "" {
		kind = "NOTE"
	}

	var borderClass, textClass, iconSVG string

	switch kind {
	case "TIP":
		borderClass = "border-l-[var(--c-ok,#10b981)]"
		textClass = "text-[var(--c-ok,#10b981)]"
		iconSVG = `<svg class="w-4 h-4 ` + textClass + ` shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 2v4m0 12v4M4.93 4.93l2.83 2.83m8.48 8.48 2.83 2.83M2 12h4m12 0h4M4.93 19.07l2.83-2.83m8.48-8.48 2.83-2.83"/></svg>`

	case "WARNING":
		borderClass = "border-l-amber-500"
		textClass = "text-amber-500"
		iconSVG = `<svg class="w-4 h-4 ` + textClass + ` shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"></path><line x1="12" y1="9" x2="12" y2="13"></line><line x1="12" y1="17" x2="12.01" y2="17"></line></svg>`

	case "IMPORTANT":
		borderClass = "border-l-cyan-500"
		textClass = "text-cyan-500"
		iconSVG = `<svg class="w-4 h-4 ` + textClass + ` shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>`

	case "CAUTION":
		borderClass = "border-l-red-500"
		textClass = "text-red-500"
		iconSVG = `<svg class="w-4 h-4 ` + textClass + ` shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="7.86 2 16.14 2 22 7.86 22 16.14 16.14 22 7.86 22 2 16.14 2 7.86 7.86 2"></polygon><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>`

	default: // NOTE
		kind = "NOTE"
		borderClass = "border-l-[var(--c-link,#3b82f6)]"
		textClass = "text-[var(--c-link,#3b82f6)]"
		iconSVG = `<svg class="w-4 h-4 ` + textClass + ` shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>`
	}

	bodyEsc := formatSimpleInline(html.EscapeString(strings.TrimSpace(rawBody)))
	bodyLines := strings.Split(bodyEsc, "\n")
	var formattedBody []string
	for _, l := range bodyLines {
		t := strings.TrimSpace(l)
		if t != "" {
			formattedBody = append(formattedBody, fmt.Sprintf("<p class=\"m-0 leading-relaxed\">%s</p>", t))
		}
	}

	return fmt.Sprintf(`
<div class="post-callout my-6 p-4 border-l-4 border-t border-r border-b border-border bg-surface not-prose rounded-none %s">
  <div class="flex items-center gap-2 mb-2">
    %s
    <span class="callout-title font-mono text-[0.8em] font-bold uppercase tracking-wider %s">%s</span>
  </div>
  <div class="callout-body text-[0.92em] text-text/90 space-y-2">
    %s
  </div>
</div>
`, borderClass, iconSVG, textClass, kind, strings.Join(formattedBody, "\n"))
}
