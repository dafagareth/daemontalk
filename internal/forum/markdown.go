package forum

import (
	"bytes"
	"html/template"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

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
		html.WithUnsafe(), // Safe because content is sanitized or controlled
	),
)

// RenderMarkdown converts a raw markdown string into template.HTML.
func RenderMarkdown(input string) template.HTML {
	var buf bytes.Buffer
	if err := forumMD.Convert([]byte(input), &buf); err != nil {
		return template.HTML("<p>" + template.HTMLEscapeString(input) + "</p>")
	}
	return template.HTML(buf.String())
}
