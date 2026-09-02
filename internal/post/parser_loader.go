package post

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark/parser"
)

func LoadAll(dir string) ([]Post, error) {
	return loadDir(dir, false)
}

func LoadAllWithDrafts(dir string) ([]Post, error) {
	return loadDir(dir, true)
}

// LoadArchived loads all .md.archive files from the directory.
func LoadArchived(dir string) ([]Post, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read dir: %w", err)
	}

	var posts []Post
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md.archive") {
			continue
		}
		p, err := parseFile(filepath.Join(dir, e.Name()))
		if err != nil {
			slog.Warn("parse archived post failed", "file", e.Name(), "error", err)
			continue
		}
		if p.Slug == "" {
			p.Slug = strings.TrimSuffix(e.Name(), ".md.archive")
		}
		posts = append(posts, p)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})
	return posts, nil
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

func getGitCreationDate(path string) (time.Time, error) {
	cmd := exec.Command("git", "log", "--reverse", "--format=%aI", "--", path)
	out, err := cmd.Output()
	if err != nil {
		return time.Time{}, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return time.Time{}, fmt.Errorf("no git history")
	}
	return time.Parse(time.RFC3339, lines[0])
}

func ParseFile(path string) (Post, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return Post{}, fmt.Errorf("stat file: %w", err)
	}

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
	// Fallback if frontmatter does not define date
	if p.Date.IsZero() {
		if gitDate, err := getGitCreationDate(path); err == nil && !gitDate.IsZero() {
			p.Date = gitDate
		} else {
			// Untracked files fallback to file modification time
			p.Date = fi.ModTime()
		}
	}
	return p, nil
}

func parseFile(path string) (Post, error) {
	return ParseFile(path)
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
