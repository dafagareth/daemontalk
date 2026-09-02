package post

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAllPostsFrontmatterIntegrity(t *testing.T) {
	postsDir := "../../content/posts"
	files, err := os.ReadDir(postsDir)
	if err != nil {
		t.Skipf("skipping posts integrity test: %v", err)
		return
	}

	slugs := make(map[string]string)

	for _, fi := range files {
		if fi.IsDir() || !strings.HasSuffix(fi.Name(), ".md") {
			continue
		}

		fullPath := filepath.Join(postsDir, fi.Name())
		p, err := ParseFile(fullPath)
		if err != nil {
			t.Errorf("file %s failed to parse frontmatter: %v", fi.Name(), err)
			continue
		}

		if p.Title == "" {
			t.Errorf("file %s missing title", fi.Name())
		}
		if p.Slug == "" {
			t.Errorf("file %s missing slug", fi.Name())
		}

		// Slug uniqueness check per language
		key := p.Slug + ":" + p.Lang
		if existingFile, exists := slugs[key]; exists {
			t.Errorf("duplicate slug %q for lang %q in %s and %s", p.Slug, p.Lang, fi.Name(), existingFile)
		} else {
			slugs[key] = fi.Name()
		}

		if p.Date.IsZero() {
			t.Errorf("file %s has invalid date", fi.Name())
		}
	}
}
