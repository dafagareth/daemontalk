package highlight

import (
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
)

// GenerateCSS returns Chroma CSS for light (github) and dark (dracula) themes,
// with dark rules scoped under [data-theme="dark"] and the system media query.
func GenerateCSS() string {
	var buf strings.Builder

	writeCSSWithPrefix(&buf, "github", "")

	writeCSSWithPrefix(&buf, "dracula", `[data-theme="dark"]`)

	buf.WriteString("\n@media (prefers-color-scheme: dark) {\n")
	writeCSSWithPrefix(&buf, "dracula", `html:not([data-theme="light"])`)
	buf.WriteString("}\n")

	return buf.String()
}

func writeCSSWithPrefix(buf *strings.Builder, styleName, prefix string) {
	style := styles.Get(styleName)
	if style == nil {
		style = styles.Fallback
	}

	var raw strings.Builder
	formatter := chromahtml.New(chromahtml.WithClasses(true))
	_ = formatter.WriteCSS(&raw, style)

	for _, line := range strings.Split(raw.String(), "\n") {
		if strings.TrimSpace(line) == "" {
			buf.WriteString("\n")
			continue
		}
		if prefix == "" {
			buf.WriteString(line + "\n")
			continue
		}
		braceIdx := strings.Index(line, "{")
		if braceIdx == -1 {
			buf.WriteString(line + "\n")
			continue
		}
		beforeBrace := line[:braceIdx]
		selectorStart := 0
		if idx := strings.LastIndex(beforeBrace, "*/"); idx != -1 {
			selectorStart = idx + 2
		}
		comment := line[:selectorStart]
		selector := strings.TrimSpace(beforeBrace[selectorStart:])
		rest := line[braceIdx:]
		buf.WriteString(comment + prefix + " " + selector + " " + rest + "\n")
	}
}
