package project

import (
	"net/url"
	"testing"
)

func TestProjectsData(t *testing.T) {
	if len(All) == 0 {
		t.Fatalf("expected non-empty projects list")
	}

	slugs := make(map[string]bool)

	for i, p := range All {
		if p.Name == "" {
			t.Errorf("project at index %d has empty Name", i)
		}
		if p.Slug == "" {
			t.Errorf("project at index %d has empty Slug", i)
		}
		if slugs[p.Slug] {
			t.Errorf("duplicate project slug: %s", p.Slug)
		}
		slugs[p.Slug] = true

		if p.Description == "" {
			t.Errorf("project %q has empty Description", p.Slug)
		}
		if p.DescriptionID == "" {
			t.Errorf("project %q has empty DescriptionID", p.Slug)
		}
		if len(p.TechStack) == 0 {
			t.Errorf("project %q has empty TechStack", p.Slug)
		}
		if p.RepoURL != "" {
			u, err := url.Parse(p.RepoURL)
			if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
				t.Errorf("project %q has invalid RepoURL: %s", p.Slug, p.RepoURL)
			}
		}
		if p.Status != StatusActive && p.Status != StatusCompleted && p.Status != StatusArchived {
			t.Errorf("project %q has invalid Status: %s", p.Slug, p.Status)
		}
	}
}
