package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"daemontalk/internal/post"
	"github.com/go-chi/chi/v5"
)

// BlogIndex is used as "/" handler after pivot; must render the list of
// posts with page title "daemontalk" (not "Blog · daemontalk").
func TestBlogIndexAsHome(t *testing.T) {
	h := &Handler{
		FilePosts: []post.Post{{Title: "Tulisan Pertama", Slug: "tulisan-pertama"}},
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.BlogIndex(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /: got %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Tulisan Pertama") {
		t.Error("home harus berisi daftar tulisan")
	}
	if !strings.Contains(body, "<title>daemontalk</title>") {
		t.Error("home harus berjudul 'daemontalk', bukan 'Blog · daemontalk'")
	}
}

func TestBlogPostAliasRedirectsToCanonicalShortID(t *testing.T) {
	h := &Handler{
		FilePosts: []post.Post{
			{
				Title:   "Arch Linux Tips",
				Slug:    "803461ff",
				Aliases: []string{"archlinux-funfact-tips"},
			},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/blog/archlinux-funfact-tips", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("slug", "archlinux-funfact-tips")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	h.BlogPost(rec, req)

	if rec.Code != http.StatusMovedPermanently {
		t.Fatalf("GET /blog/archlinux-funfact-tips: got %d, want 301", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/blog/803461ff" {
		t.Errorf("Location header: got %q, want /blog/803461ff", loc)
	}
}
