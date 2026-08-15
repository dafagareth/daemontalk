package post

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
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

func LoadAll(dir string) ([]Post, error) {
	return loadDir(dir, false)
}

// LoadAllWithDrafts loads all posts including drafts, sorted by date descending.
func LoadAllWithDrafts(dir string) ([]Post, error) {
	return loadDir(dir, true)
}

func loadDir(dir string, includeDrafts bool) ([]Post, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read dir: %w", err)
	}

	var posts []Post
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		p, err := parseFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", e.Name(), err)
		}
		if p.Draft && !includeDrafts {
			continue
		}
		posts = append(posts, p)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})
	return posts, nil
}

func parseFile(path string) (Post, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return Post{}, fmt.Errorf("read file: %w", err)
	}

	p, err := Parse(src)
	if err != nil {
		return Post{}, err
	}
	if p.Slug == "" {
		base := filepath.Base(path)
		p.Slug = strings.TrimSuffix(base, ".md")
	}
	return p, nil
}

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
	p.Description = extractDescription(rendered)
	return p, nil
}

func extractTOC(htmlStr string) []TOCEntry {
	var toc []TOCEntry
	for _, m := range reHeading.FindAllStringSubmatch(htmlStr, -1) {
		level, _ := strconv.Atoi(m[1])
		id := m[2]
		title := strings.TrimSpace(reTag.ReplaceAllString(m[3], ""))
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

func FindBySlug(posts []Post, slug string) (Post, bool) {
	for _, p := range posts {
		if p.Slug == slug {
			return p, true
		}
	}
	for _, p := range posts {
		for _, a := range p.Aliases {
			if a == slug {
				return p, true
			}
		}
	}
	return Post{}, false
}

// LoadBody renders a markdown file and returns its HTML body.
// Unlike parseFile, frontmatter is not extracted (only body content is returned).
func LoadBody(path string) (template.HTML, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	processedSrc := preprocessMarkdown(src)
	ctx := parser.NewContext()
	var buf bytes.Buffer
	if err := md.Convert(processedSrc, &buf, parser.WithContext(ctx)); err != nil {
		return "", fmt.Errorf("render markdown: %w", err)
	}
	rawHTML := strings.ReplaceAll(buf.String(), "<img ", `<img loading="lazy" `)
	return template.HTML(rawHTML), nil
}

// LoadBodyWithTOC renders a markdown file and returns its HTML body along with
// a table of contents extracted from its headings.
func LoadBodyWithTOC(path string) (template.HTML, []TOCEntry, error) {
	body, err := LoadBody(path)
	if err != nil {
		return "", nil, err
	}
	return body, extractTOC(string(body)), nil
}
