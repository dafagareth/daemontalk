package post

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yuin/goldmark/parser"
)

func LoadAll(dir string) ([]Post, error) {
	return loadDir(dir, false)
}

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

func LoadBodyWithTOC(path string) (template.HTML, []TOCEntry, error) {
	body, err := LoadBody(path)
	if err != nil {
		return "", nil, err
	}
	return body, extractTOC(string(body)), nil
}
