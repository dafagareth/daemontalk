package post

import (
	"bytes"
	"fmt"
	stdhtml "html"
	"html/template"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	mathjax "github.com/litao91/goldmark-mathjax"
)

var (
	reHeading = regexp.MustCompile(`(?i)<h([23])[^>]*\bid="([^"]+)"[^>]*>([\s\S]*?)</h[23]>`)
	reTag     = regexp.MustCompile(`<[^>]+>`)
	reFirstP  = regexp.MustCompile(`(?i)<p>([\s\S]*?)</p>`)
)

type externalLinkTransformer struct{}

func (t *externalLinkTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	_ = ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if link, ok := n.(*ast.Link); ok {
			dest := string(link.Destination)
			if strings.HasPrefix(dest, "http://") || strings.HasPrefix(dest, "https://") {
				link.SetAttribute([]byte("target"), []byte("_blank"))
				link.SetAttribute([]byte("rel"), []byte("noopener noreferrer"))
			}
		} else if autoLink, ok := n.(*ast.AutoLink); ok {
			dest := string(autoLink.URL(reader.Source()))
			if strings.HasPrefix(dest, "http://") || strings.HasPrefix(dest, "https://") {
				autoLink.SetAttribute([]byte("target"), []byte("_blank"))
				autoLink.SetAttribute([]byte("rel"), []byte("noopener noreferrer"))
			}
		}
		return ast.WalkContinue, nil
	})
}

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		extension.Footnote,
		meta.Meta,
		mathjax.MathJax,
		highlighting.NewHighlighting(
			highlighting.WithFormatOptions(
				chromahtml.WithClasses(true),
			),
		),
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
		parser.WithASTTransformers(
			util.Prioritized(&externalLinkTransformer{}, 500),
		),
	),
	goldmark.WithRendererOptions(
		html.WithUnsafe(),
	),
)

// LoadAllWithDrafts loads all posts including drafts, sorted by date descending.

// Parse renders a markdown document (frontmatter + body) into a Post.
// It is the shared pipeline for file-based posts and web-authored posts
// stored in the database, so both render identically.
func Parse(src []byte) (Post, error) {
	processedSrc := preprocessMarkdown(src)
	ctx := parser.NewContext()
	var buf bytes.Buffer
	if err := md.Convert(processedSrc, &buf, parser.WithContext(ctx)); err != nil {
		return Post{}, fmt.Errorf("render markdown: %w", err)
	}

	fm := meta.Get(ctx)

	// Inject lazy loading for all images in the rendered HTML.
	rawHTML := strings.ReplaceAll(buf.String(), "<img ", `<img loading="lazy" `)

	// Restore protected math characters
	rawHTML = strings.ReplaceAll(rawHTML, "xDTESCAPEDUSCOREx", `\_`)
	rawHTML = strings.ReplaceAll(rawHTML, "xDTUSCOREx", "_")
	rawHTML = strings.ReplaceAll(rawHTML, "xDTASTx", "*")

	p := Post{
		Body: template.HTML(rawHTML),
	}

	if v, ok := fm["title"].(string); ok {
		p.Title = v
	}
	if v, ok := fm["slug"].(string); ok {
		p.Slug = v
	}
	if v, ok := fm["lang"].(string); ok {
		p.Lang = v
	}
	if v, ok := fm["draft"].(bool); ok {
		p.Draft = v
	}
	if v, ok := fm["date"].(string); ok {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			p.Date = t
		}
	}
	if v, ok := fm["cover"].(string); ok {
		p.Cover = v
	}
	if v, ok := fm["coverCaption"].(string); ok {
		p.CoverCaption = v
	} else if v, ok := fm["cover_caption"].(string); ok {
		p.CoverCaption = v
	}
	if v, ok := fm["coverSource"].(string); ok {
		p.CoverSource = v
	} else if v, ok := fm["cover_source"].(string); ok {
		p.CoverSource = v
	}
	if v, ok := fm["series"].(string); ok {
		p.Series = v
	}
	if v, ok := fm["series_part"].(int); ok {
		p.SeriesPart = v
	}
	if v, ok := fm["publish_at"].(string); ok {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			p.PublishAt = t
		}
	}
	if v, ok := fm["type"].(string); ok {
		p.Type = v
	}
	if v, ok := fm["author"].(string); ok {
		p.Author = v
	}
	if v, ok := fm["description"].(string); ok && strings.TrimSpace(v) != "" {
		p.Description = stdhtml.UnescapeString(strings.TrimSpace(v))
	}

	if raw, ok := fm["tags"]; ok {
		switch v := raw.(type) {
		case []interface{}:
			for _, t := range v {
				if s, ok := t.(string); ok {
					p.Tags = append(p.Tags, s)
				}
			}
		case []string:
			p.Tags = v
		}
	}

	if raw, ok := fm["aliases"]; ok {
		switch v := raw.(type) {
		case []interface{}:
			for _, t := range v {
				if s, ok := t.(string); ok && s != "" {
					p.Aliases = append(p.Aliases, s)
				}
			}
		case []string:
			p.Aliases = v
		}
	}
	if v, ok := fm["alias"].(string); ok && v != "" {
		p.Aliases = append(p.Aliases, v)
	}
	if p.Lang == "en" {
		rawHTML = strings.ReplaceAll(rawHTML, `<h2 id="referensi">Referensi</h2>`, `<h2 id="references">References</h2>`)
		rawHTML = strings.ReplaceAll(rawHTML, `<h3 id="referensi">Referensi</h3>`, `<h3 id="references">References</h3>`)
		rawHTML = strings.ReplaceAll(rawHTML, `<h4 id="referensi">Referensi</h4>`, `<h4 id="references">References</h4>`)
		rawHTML = strings.ReplaceAll(rawHTML, `<h2 id="bibliografi-dan-referensi">Bibliografi dan Referensi</h2>`, `<h2 id="bibliography-and-references">Bibliography &amp; References</h2>`)
		rawHTML = strings.ReplaceAll(rawHTML, `<h2 id="daftar-pustaka">Daftar Pustaka</h2>`, `<h2 id="references">References</h2>`)
	}

	rendered := rawHTML
	p.Body = template.HTML(rawHTML)
	p.ReadTime = readTime(rendered)
	p.TOC = extractTOC(rendered)
	if p.Description == "" {
		p.Description = extractDescription(rendered)
	}
	return p, nil
}

func extractTOC(htmlStr string) []TOCEntry {
	var toc []TOCEntry
	for _, m := range reHeading.FindAllStringSubmatch(htmlStr, -1) {
		level, _ := strconv.Atoi(m[1])
		id := m[2]
		title := stdhtml.UnescapeString(strings.TrimSpace(reTag.ReplaceAllString(m[3], "")))
		if title != "" {
			toc = append(toc, TOCEntry{ID: id, Title: title, Level: level})
		}
	}
	return toc
}

func extractDescription(htmlStr string) string {
	m := reFirstP.FindStringSubmatch(htmlStr)
	if m == nil {
		return ""
	}
	text := strings.TrimSpace(reTag.ReplaceAllString(m[1], ""))
	text = stdhtml.UnescapeString(text)
	if len([]rune(text)) > 160 {
		runes := []rune(text)
		return string(runes[:157]) + "..."
	}
	return text
}

func readTime(html string) int {
	text := strings.Map(func(r rune) rune {
		if r == '<' {
			return ' '
		}
		if unicode.IsLetter(r) || unicode.IsSpace(r) {
			return r
		}
		return ' '
	}, html)

	words := 0
	inWord := false
	for _, r := range text {
		if unicode.IsLetter(r) {
			if !inWord {
				words++
				inWord = true
			}
		} else {
			inWord = false
		}
	}

	minutes := words / 200
	if minutes < 1 {
		return 1
	}
	return minutes
}

