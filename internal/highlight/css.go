package highlight

import (
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
)

// GenerateCSS returns Chroma CSS for light (github), dark (github-dark), and sepia (solarized-light) themes,
// with proper scoping so rules from one theme never leak into or corrupt another theme.
func GenerateCSS() string {
	var buf strings.Builder

	// Base defaults: ensure identifiers, variables, punctuation inherit readable text color
	buf.WriteString(`/* Base Chroma reset & variable inheritance */
.chroma {
  color: var(--c-text);
  background-color: var(--c-surface);
}
.chroma .nx,
.chroma .p,
.chroma .nn,
.chroma .n {
  color: inherit;
}

`)

	// 1. Light theme (github) — active on light theme or default when no dark/sepia preference
	writeCSSWithPrefixes(&buf, "github", []string{
		`[data-theme="light"]`,
		`html:not([data-theme="dark"]):not([data-theme="sepia"])`,
	})

	// 2. Dark theme (github-dark) — manual dark toggle
	buf.WriteString("\n/* Dark Theme (GitHub Dark) */\n")
	writeCSSWithPrefixes(&buf, "github-dark", []string{
		`[data-theme="dark"]`,
	})
	buf.WriteString("[data-theme=\"dark\"] .chroma { color: #e6edf3; }\n")
	buf.WriteString("[data-theme=\"dark\"] .chroma .nx, [data-theme=\"dark\"] .chroma .p, [data-theme=\"dark\"] .chroma .nn, [data-theme=\"dark\"] .chroma .n { color: #e6edf3; }\n")

	// 3. System preference dark (when theme is not explicitly set to light or sepia)
	buf.WriteString("\n@media (prefers-color-scheme: dark) {\n")
	writeCSSWithPrefixes(&buf, "github-dark", []string{
		`html:not([data-theme="light"]):not([data-theme="sepia"])`,
	})
	buf.WriteString("  html:not([data-theme=\"light\"]):not([data-theme=\"sepia\"]) .chroma { color: #e6edf3; }\n")
	buf.WriteString("  html:not([data-theme=\"light\"]):not([data-theme=\"sepia\"]) .chroma .nx, html:not([data-theme=\"light\"]):not([data-theme=\"sepia\"]) .chroma .p, html:not([data-theme=\"light\"]):not([data-theme=\"sepia\"]) .chroma .nn, html:not([data-theme=\"light\"]):not([data-theme=\"sepia\"]) .chroma .n { color: #e6edf3; }\n")
	buf.WriteString("}\n")

	// 4. Sepia theme (solarized-light) — reading mode
	buf.WriteString("\n/* Sepia Theme (Solarized Light) */\n")
	writeCSSWithPrefixes(&buf, "solarized-light", []string{
		`[data-theme="sepia"]`,
	})
	buf.WriteString("[data-theme=\"sepia\"] .chroma { color: #586e75; }\n")
	buf.WriteString("[data-theme=\"sepia\"] .chroma .nx, [data-theme=\"sepia\"] .chroma .p, [data-theme=\"sepia\"] .chroma .nn, [data-theme=\"sepia\"] .chroma .n { color: #586e75; }\n")

	return buf.String()
}

func writeCSSWithPrefixes(buf *strings.Builder, styleName string, prefixes []string) {
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

		var combined []string
		for _, prefix := range prefixes {
			if prefix == "" {
				combined = append(combined, selector)
			} else {
				combined = append(combined, prefix+" "+selector)
			}
		}
		buf.WriteString(comment + strings.Join(combined, ", ") + " " + rest + "\n")
	}
}
