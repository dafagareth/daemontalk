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

// formatSimpleInline provides basic bold, italic, code formatting for custom blocks.
func formatSimpleInline(s string) string {
	s = reInlineCode.ReplaceAllString(s, `<code class="font-mono text-xs px-1.5 py-0.5 bg-chip border border-border text-text">$1</code>`)
	s = reInlineBold.ReplaceAllString(s, `<strong class="font-semibold text-text">$1</strong>`)
	s = reInlineItalic.ReplaceAllString(s, `<em class="italic text-text">$1</em>`)
	return s
}

// safeURL validates that a URL begins with an allowed scheme (http, https, mailto, or relative /)
// to prevent javascript: or data: XSS vectors in markdown link cards and references.
func safeURL(raw string) string {
	raw = strings.TrimSpace(raw)
	lower := strings.ToLower(raw)
	if strings.HasPrefix(lower, "https://") || strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "mailto:") || strings.HasPrefix(raw, "/") {
		return raw
	}
	return "#"
}
