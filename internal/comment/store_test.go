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

func TestThreadedRepliesAndBuildTree(t *testing.T) {
	s := newTestStore(t)

	// Root comment 1
	c1, err := s.Add("post-tree", "Alice", "Root comment 1")
	if err != nil {
		t.Fatalf("Add root 1: %v", err)
	}

	// Reply to Root 1
	c1Reply1, err := s.AddWithParent("post-tree", "Bob", "Reply 1 to Alice", &c1.ID)
	if err != nil {
		t.Fatalf("Add reply 1: %v", err)
	}

	// Nested reply to Reply 1 (grandchild)
	c1Grandchild, err := s.AddWithParent("post-tree", "Carol", "Reply to Bob", &c1Reply1.ID)
	if err != nil {
		t.Fatalf("Add grandchild: %v", err)
	}

	// Root comment 2
	c2, err := s.Add("post-tree", "Dave", "Root comment 2")
	if err != nil {
		t.Fatalf("Add root 2: %v", err)
	}

	list, err := s.ListBySlug("post-tree")
	if err != nil {
		t.Fatalf("ListBySlug: %v", err)
	}
	if len(list) != 4 {
		t.Fatalf("expected 4 flat comments, got %d", len(list))
	}

	tree := BuildTree(list)
	if len(tree) != 2 {
		t.Fatalf("expected 2 root comments in tree, got %d", len(tree))
	}

	// Check root 1 (Alice)
	if tree[0].ID != c1.ID || tree[0].Name != "Alice" {
		t.Errorf("expected root 1 to be Alice, got %q", tree[0].Name)
	}
	if len(tree[0].Replies) != 1 {
		t.Fatalf("expected 1 reply under Alice, got %d", len(tree[0].Replies))
	}

	// Check reply 1 (Bob)
	bobReply := tree[0].Replies[0]
	if bobReply.ID != c1Reply1.ID || bobReply.Name != "Bob" {
		t.Errorf("expected reply to be Bob, got %q", bobReply.Name)
	}
	if len(bobReply.Replies) != 1 {
		t.Fatalf("expected 1 grandchild reply under Bob, got %d", len(bobReply.Replies))
	}

	// Check grandchild (Carol)
	carolReply := bobReply.Replies[0]
	if carolReply.ID != c1Grandchild.ID || carolReply.Name != "Carol" {
		t.Errorf("expected grandchild to be Carol, got %q", carolReply.Name)
	}

	// Check root 2 (Dave)
	if tree[1].ID != c2.ID || tree[1].Name != "Dave" {
		t.Errorf("expected root 2 to be Dave, got %q", tree[1].Name)
	}
	if len(tree[1].Replies) != 0 {
		t.Errorf("expected 0 replies under Dave, got %d", len(tree[1].Replies))
	}
}
