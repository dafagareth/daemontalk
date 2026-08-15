package postdb

import (
	"path/filepath"
	"testing"

	"daemontalk/internal/post"
)

func openTest(t *testing.T) *Store {
	t.Helper()
	s, err := Open(filepath.Join(t.TempDir(), "posts.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	return s
}

func TestCRUD(t *testing.T) {
	s := openTest(t)

	id, err := s.Create(WebPost{
		Slug: "halo-dunia", Title: "Halo Dunia", BodyMD: "Isi **tebal**.",
		Tags: "go, web", Lang: "id", Draft: true, Date: "2026-07-04",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	p, err := s.Get(id)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if p.Title != "Halo Dunia" || !p.Draft || p.Lang != "id" {
		t.Errorf("Get salah: %+v", p)
	}

	p.Title = "Halo Lagi"
	p.Draft = false
	if err := s.Update(p); err != nil {
		t.Fatalf("Update: %v", err)
	}
	p2, _ := s.Get(id)
	if p2.Title != "Halo Lagi" || p2.Draft {
		t.Errorf("Update tidak tersimpan: %+v", p2)
	}

	if err := s.Delete(id); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := s.Get(id); err == nil {
		t.Error("Get setelah Delete harus error")
	}
}

func TestListOrderAndUniqueSlug(t *testing.T) {
	s := openTest(t)
	if _, err := s.Create(WebPost{Slug: "lama", Title: "Lama", BodyMD: "x", Date: "2026-01-01"}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Create(WebPost{Slug: "baru", Title: "Baru", BodyMD: "x", Date: "2026-07-01"}); err != nil {
		t.Fatal(err)
	}

	list, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 2 || list[0].Slug != "baru" {
		t.Errorf("urutan List salah: %+v", list)
	}

	if _, err := s.Create(WebPost{Slug: "baru", Title: "Dobel", BodyMD: "x", Date: "2026-07-02"}); err == nil {
		t.Error("slug dobel harus error")
	}
}

func TestToMarkdownRoundTrip(t *testing.T) {
	wp := WebPost{
		Slug: "post-web", Title: "Post \"Web\"", BodyMD: "Pembuka.\n\n## Judul Bagian\n\nIsi.",
		Tags: "go, cli", Lang: "id", Draft: true, Date: "2026-07-04",
	}
	p, err := post.Parse(wp.ToMarkdown())
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if p.Title != `Post "Web"` || p.Slug != "post-web" || !p.Draft {
		t.Errorf("frontmatter salah: %+v", p)
	}
	if len(p.Tags) != 2 || p.Tags[0] != "go" {
		t.Errorf("tags: %v", p.Tags)
	}
	if len(p.TOC) != 1 {
		t.Errorf("TOC: %+v", p.TOC)
	}
	if p.Date.Format("2006-01-02") != "2026-07-04" {
		t.Errorf("date: %v", p.Date)
	}
}
