package post

import (
	"regexp"
	"strings"
)

var (
	reInlineCode   = regexp.MustCompile("`([^`]+)`")
	reInlineBold   = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	reInlineItalic = regexp.MustCompile(`\*([^*]+)\*`)
)

func formatSimpleInline(s string) string {
	s = reInlineCode.ReplaceAllString(s, `<code class="font-mono text-xs px-1.5 py-0.5 bg-chip border border-border text-text">$1</code>`)
	s = reInlineBold.ReplaceAllString(s, `<strong class="font-semibold text-text">$1</strong>`)
	s = reInlineItalic.ReplaceAllString(s, `<em class="italic text-text">$1</em>`)
	return s
}

func safeURL(raw string) string {
	raw = strings.TrimSpace(raw)
	lower := strings.ToLower(raw)
	if strings.HasPrefix(lower, "https://") || strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "mailto:") || strings.HasPrefix(raw, "/") {
		return raw
	}
	return "#"
}
