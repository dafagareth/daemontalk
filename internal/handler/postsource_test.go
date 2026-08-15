package handler

import (
	"path/filepath"
	"testing"
	"time"

	"daemontalk/internal/post"
	"daemontalk/internal/postdb"
)

func TestAllPostsMergesDB(t *testing.T) {
	pdb, err := postdb.Open(filepath.Join(t.TempDir(), "posts.db"))
	if err != nil {
		t.Fatalf("open postdb: %v", err)
	}
	id, err := pdb.Create(postdb.WebPost{
		Slug: "dari-web", Title: "Post dari Web", BodyMD: "Isi post dari editor.",
		Lang: "id", Date: "2026-07-04",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	h := &Handler{
		FilePosts: []post.Post{{Title: "Dari File", Slug: "dari-file", Date: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}},
		PostDB:    pdb,
	}

	all := h.AllPosts()
	if len(all) != 2 {
		t.Fatalf("AllPosts: dapat %d, mau 2 (%+v)", len(all), all)
	}
	// Post DB bertanggal lebih baru → harus di urutan pertama.
	if all[0].Slug != "dari-web" || all[1].Slug != "dari-file" {
		t.Errorf("urutan salah: %s, %s", all[0].Slug, all[1].Slug)
	}

	if err := pdb.Delete(id); err != nil {
		t.Fatal(err)
	}
	h.RefreshPosts()
	if got := len(h.AllPosts()); got != 1 {
		t.Errorf("setelah delete+refresh: dapat %d, mau 1", got)
	}
}
