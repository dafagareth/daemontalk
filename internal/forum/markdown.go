package forum

import (
	"bytes"
	"html/template"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

var sanitizer = func() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").OnElements("pre", "code", "span", "div")
	p.AllowAttrs("tabindex").OnElements("pre")
	return p
}()

var forumMD = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		extension.Typographer,
		highlighting.NewHighlighting(
			highlighting.WithFormatOptions(
				chromahtml.WithClasses(true),
				chromahtml.TabWidth(2),
			),
		),
	),
	goldmark.WithRendererOptions(
		html.WithHardWraps(),
		html.WithXHTML(),
	),
)

// RenderMarkdown converts a raw markdown string into template.HTML.
func RenderMarkdown(input string) template.HTML {
	var buf bytes.Buffer
	if err := forumMD.Convert([]byte(input), &buf); err != nil {
		return template.HTML("<p>" + template.HTMLEscapeString(input) + "</p>")
	}
	clean := sanitizer.SanitizeBytes(buf.Bytes())
	return template.HTML(clean)
}
