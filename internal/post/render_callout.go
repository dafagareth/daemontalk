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
		iconSVG = GetIcon(IconCalloutTip, "w-4 h-4 "+textClass)

	case "WARNING":
		borderClass = "border-l-amber-500"
		textClass = "text-amber-500"
		iconSVG = GetIcon(IconCalloutWarning, "w-4 h-4 "+textClass)

	case "IMPORTANT":
		borderClass = "border-l-cyan-500"
		textClass = "text-cyan-500"
		iconSVG = GetIcon(IconCalloutImportant, "w-4 h-4 "+textClass)

	case "CAUTION":
		borderClass = "border-l-red-500"
		textClass = "text-red-500"
		iconSVG = GetIcon(IconCalloutCaution, "w-4 h-4 "+textClass)

	default: // NOTE
		kind = "NOTE"
		borderClass = "border-l-[var(--c-link,#3b82f6)]"
		textClass = "text-[var(--c-link,#3b82f6)]"
		iconSVG = GetIcon(IconCalloutNote, "w-4 h-4 "+textClass)
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
