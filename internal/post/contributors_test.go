package post

import (
	"testing"
)

func TestGetAllContributorsAndIsContributor(t *testing.T) {
	posts := []Post{
		{
			Title:        "Post 1",
			AuthorGitHub: "dafagareth",
			Contributors: []string{"budisantoso", "alexr"},
			Draft:        false,
		},
		{
			Title:        "Post 2",
			AuthorGitHub: "dafagareth",
			Contributors: []string{"budisantoso"},
			Draft:        false,
		},
	}

	contribs := GetAllContributors(posts)
	if len(contribs) != 3 {
		t.Fatalf("expected 3 unique contributors, got %d", len(contribs))
	}

	if !IsContributor(posts, "budisantoso") {
		t.Errorf("expected budisantoso to be contributor")
	}
	if !IsContributor(posts, "DAFAGARETH") {
		t.Errorf("expected DAFAGARETH case-insensitive to be contributor")
	}
	if IsContributor(posts, "unknown_user") {
		t.Errorf("expected unknown_user to not be contributor")
	}
}
