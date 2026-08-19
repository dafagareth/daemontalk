package post

import (
	"regexp"
)

// formatSimpleInline provides basic bold, italic, code formatting for custom blocks.
func formatSimpleInline(s string) string {
	// Format inline `code`
	reCode := regexp.MustCompile("`([^`]+)`")
	s = reCode.ReplaceAllString(s, `<code class="font-mono text-xs px-1.5 py-0.5 bg-chip border border-border text-text">$1</code>`)

	// Format inline **bold**
	reBold := regexp.MustCompile(`\*\*([^*]+)\*\*`)
	s = reBold.ReplaceAllString(s, `<strong class="font-semibold text-text">$1</strong>`)

	// Format inline *italic*
	reItalic := regexp.MustCompile(`\*([^*]+)\*`)
	s = reItalic.ReplaceAllString(s, `<em class="italic text-text">$1</em>`)

	return s
}
