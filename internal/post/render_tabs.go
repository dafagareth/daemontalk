package post

import (
	"bytes"
	"fmt"
	"html"
	"strings"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

func highlightCode(code, lang, filename string) string {
	var lexer chroma.Lexer
	if lang != "" {
		lexer = lexers.Get(lang)
	}
	if lexer == nil && filename != "" {
		lexer = lexers.Match(filename)
	}
	if lexer == nil {
		lexer = lexers.Analyse(code)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return "<pre class=\"chroma\"><code>" + html.EscapeString(code) + "</code></pre>"
	}

	formatter := chromahtml.New(chromahtml.WithClasses(true))
	var buf bytes.Buffer
	if err := formatter.Format(&buf, styles.Fallback, iterator); err != nil {
		return "<pre class=\"chroma\"><code>" + html.EscapeString(code) + "</code></pre>"
	}
	return buf.String()
}

type codeTab struct {
	Name    string
	Lang    string
	Content string
}

func renderTabsHTML(rawContent string) string {
	lines := strings.Split(rawContent, "\n")
	var tabs []codeTab
	var current codeTab
	var contentLines []string

	flush := func() {
		if current.Name != "" || len(contentLines) > 0 {
			current.Content = strings.TrimSpace(strings.Join(contentLines, "\n"))
			tabs = append(tabs, current)
			current = codeTab{}
			contentLines = nil
		}
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "===") {
			flush()
			header := strings.TrimSpace(strings.TrimPrefix(trimmed, "==="))
			if idx := strings.Index(header, "["); idx != -1 && strings.HasSuffix(header, "]") {
				current.Name = strings.TrimSpace(header[:idx])
				current.Lang = strings.TrimSpace(strings.TrimSuffix(header[idx+1:], "]"))
			} else {
				current.Name = header
			}
			if current.Name == "" {
				current.Name = "Code"
			}
			continue
		}
		contentLines = append(contentLines, line)
	}
	flush()

	if len(tabs) == 0 {
		return ""
	}

	var buf bytes.Buffer
	buf.WriteString("\n<div class=\"code-tabs-wrap my-6 border border-border bg-surface not-prose rounded-none overflow-hidden relative\" data-code-tabs>\n")
	buf.WriteString("  <div class=\"code-tabs-header flex items-center justify-between bg-surface border-b border-border px-3 sm:px-4 gap-2\">\n")
	buf.WriteString("    <div class=\"tabs-nav-track flex items-center gap-1 overflow-x-auto scrollbar-none snap-x touch-pan-x whitespace-nowrap min-w-0 flex-1 -mb-px\" role=\"tablist\">\n")

	for i, t := range tabs {
		activeClass := ""
		borderClass := "border-transparent text-muted hover:text-text"
		if i == 0 {
			activeClass = " active"
			borderClass = "border-link text-text bg-bg/80 font-bold"
		}
		buf.WriteString(fmt.Sprintf(`      <button type="button" class="tab-btn px-3 py-1.5 text-[0.85em] font-mono border-b-2 %s shrink-0 snap-start transition-colors cursor-pointer select-none%s" data-tab-index="%d">%s</button>
`, borderClass, activeClass, i, html.EscapeString(t.Name)))
	}

	buf.WriteString(`    </div>
  </div>
  <div class="code-tabs-content">
`)

	for i, t := range tabs {
		hiddenClass := " hidden"
		activePaneClass := ""
		if i == 0 {
			hiddenClass = ""
			activePaneClass = " active"
		}
		highlighted := highlightCode(t.Content, t.Lang, t.Name)
		buf.WriteString(fmt.Sprintf(`    <div class="tab-pane%s%s relative overflow-x-auto" data-tab-pane="%d">
      <button type="button" class="copy-tab-code copy-btn absolute top-2 right-3 z-10" title="Copy code">copy</button>
      %s
    </div>
`, hiddenClass, activePaneClass, i, highlighted))
	}

	buf.WriteString("  </div>\n</div>\n")
	return buf.String()
}
