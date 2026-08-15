package comment

import (
	"path/filepath"
	"testing"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.db")
	s, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestAddAndList(t *testing.T) {
	s := newTestStore(t)

	if _, err := s.Add("post-a", "Alice", "first"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if _, err := s.Add("post-a", "Bob", "second"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if _, err := s.Add("post-b", "Carol", "other post"); err != nil {
		t.Fatalf("Add: %v", err)
	}

	got, err := s.ListBySlug("post-a")
	if err != nil {
		t.Fatalf("ListBySlug: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 comments for post-a, got %d", len(got))
	}
	// Oldest first
	if got[0].Name != "Alice" || got[1].Name != "Bob" {
		t.Errorf("wrong order: %q, %q", got[0].Name, got[1].Name)
	}
}

func TestAddValidation(t *testing.T) {
	s := newTestStore(t)

	cases := []struct{ name, body string }{
		{"", "body"},
		{"name", ""},
		{"   ", "body"},
		{"name", "   "},
	}
	for _, c := range cases {
		if _, err := s.Add("slug", c.name, c.body); err != ErrInvalid {
			t.Errorf("Add(%q,%q): expected ErrInvalid, got %v", c.name, c.body, err)
		}
	}
}

func TestAddTrimsAndTruncates(t *testing.T) {
	s := newTestStore(t)

	longName := ""
	for i := 0; i < MaxNameLen+20; i++ {
		longName += "x"
	}
	c, err := s.Add("slug", "  "+longName+"  ", "  hi  ")
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if len([]rune(c.Name)) != MaxNameLen {
		t.Errorf("name not truncated: len=%d", len([]rune(c.Name)))
	}
	if c.Body != "hi" {
		t.Errorf("body not trimmed: %q", c.Body)
	}
}

func TestDelete(t *testing.T) {
	s := newTestStore(t)

	c, _ := s.Add("slug", "Alice", "to delete")
	if err := s.Delete(c.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	got, _ := s.ListBySlug("slug")
	if len(got) != 0 {
		t.Errorf("expected 0 after delete, got %d", len(got))
	}
}

func TestViewCount(t *testing.T) {
	s := newTestStore(t)

	if n, _ := s.ViewCount("never-viewed"); n != 0 {
		t.Errorf("unviewed post should be 0, got %d", n)
	}

	for want := 1; want <= 3; want++ {
		n, err := s.IncrementView("slug")
		if err != nil {
			t.Fatalf("IncrementView: %v", err)
		}
		if n != want {
			t.Errorf("after %d increments, got %d", want, n)
		}
	}
}

func TestListAll(t *testing.T) {
	s := newTestStore(t)

	s.Add("post-a", "Alice", "a")
	s.Add("post-b", "Bob", "b")
	s.Add("post-a", "Carol", "c")

	all, err := s.ListAll()
	if err != nil {
		t.Fatalf("ListAll: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("expected 3 comments, got %d", len(all))
	}
	// Newest first — Carol was added last.
	if all[0].Name != "Carol" {
		t.Errorf("expected newest first (Carol), got %q", all[0].Name)
	}
}

func TestAllViewCounts(t *testing.T) {
	s := newTestStore(t)

	s.IncrementView("post-a")
	s.IncrementView("post-a")
	s.IncrementView("post-b")

	counts, err := s.AllViewCounts()
	if err != nil {
		t.Fatalf("AllViewCounts: %v", err)
	}
	if counts["post-a"] != 2 {
		t.Errorf("post-a: expected 2, got %d", counts["post-a"])
	}
	if counts["post-b"] != 1 {
		t.Errorf("post-b: expected 1, got %d", counts["post-b"])
	}
}
