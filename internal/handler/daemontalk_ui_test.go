package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"daemontalk/internal/post"
)

func TestDaemontalkUILayoutRendering(t *testing.T) {
	posts := []post.Post{
		{
			Title:       "First Major Lead Story",
			Slug:        "first-lead-story",
			Description: "An in-depth breakdown of high-performance architecture.",
			Tags:        []string{"go", "architecture"},
			Date:        time.Now(),
			ReadTime:    6,
			Cover:       "/static/images/custom-lead.png",
		},
		{
			Title:       "Second Featured Top Story",
			Slug:        "second-story",
			Description: "Deep dive into Linux kernel internals.",
			Tags:        []string{"linux"},
			Date:        time.Now(),
			ReadTime:    4,
		},
		{
			Title:       "Third Story Without Custom Cover",
			Slug:        "third-story",
			Description: "Docker container tuning tips and tricks.",
			Tags:        []string{"docker"},
			Date:        time.Now(),
			ReadTime:    5,
		},
		{
			Title:       "Fourth Story In Secondary Grid",
			Slug:        "fourth-story",
			Description: "Git rebase workflow guide.",
			Tags:        []string{"git"},
			Date:        time.Now(),
			ReadTime:    3,
		},
		{
			Title:       "Fifth Story In Stream Grid",
			Slug:        "fifth-story",
			Description: "SQLite performance hacks.",
			Tags:        []string{"sqlite"},
			Date:        time.Now(),
			ReadTime:    2,
		},
	}

	h := &Handler{
		FilePosts: posts,
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	h.BlogIndex(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()

	// Check Editorial Lead & Stories
	if !strings.Contains(body, "First Major Lead Story") {
		t.Error("expected lead story title")
	}
	if !strings.Contains(body, "/static/images/custom-lead.png") {
		t.Error("expected custom cover image in grid")
	}
	if !strings.Contains(body, "Second Featured Top Story") {
		t.Error("expected second story in stream")
	}
	if !strings.Contains(body, "Fifth Story In Stream Grid") {
		t.Error("expected fifth story in stream")
	}
	if !strings.Contains(body, "Moth Icon") {
		t.Error("expected Moth Icon fallback for posts without explicit cover")
	}
}

func TestTagFiltering(t *testing.T) {
	posts := []post.Post{
		{Title: "Go Post One", Slug: "go-1", Tags: []string{"go"}, Date: time.Now()},
		{Title: "Docker Post One", Slug: "docker-1", Tags: []string{"docker"}, Date: time.Now()},
		{Title: "Go Post Two", Slug: "go-2", Tags: []string{"go"}, Date: time.Now()},
	}
	h := &Handler{FilePosts: posts}

	// Filter by tag=go
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?tag=go", nil)
	h.BlogIndex(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Go Post One") || !strings.Contains(body, "Go Post Two") {
		t.Error("expected go posts to be rendered")
	}
	if strings.Contains(body, "Docker Post One") {
		t.Error("did not expect docker post when filtered by tag=go")
	}
}

func TestDaemontalkUIHTMXPartial(t *testing.T) {
	posts := make([]post.Post, 50)
	now := time.Now()
	for i := 0; i < 50; i++ {
		posts[i] = post.Post{
			Title:    fmt.Sprintf("Post Number %d", i+1),
			Slug:     fmt.Sprintf("post-%d", i+1),
			Tags:     []string{"go"},
			Date:     now.Add(-time.Duration(i) * time.Minute),
			ReadTime: 3,
		}
	}

	h := &Handler{
		FilePosts: posts,
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/blog/posts?page=2&lang=en", nil)
	h.BlogPostsPartial(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	expectedPost := fmt.Sprintf("Post Number %d", DefaultPostsPerPage+1)
	if !strings.Contains(body, expectedPost) {
		t.Errorf("expected %q in HTMX partial response", expectedPost)
	}
	if !strings.Contains(body, "hx-target=\"#blog-grid\"") {
		t.Error("expected HTMX load more target to be #blog-grid")
	}
	if !strings.Contains(body, `id="load-more-wrap" class="mt-12 mb-8 text-center"`) {
		t.Error("expected load-more-wrap to retain classes in OOB swap")
	}
}

func TestTagPagination(t *testing.T) {
	posts := make([]post.Post, 25)
	now := time.Now()
	for i := 0; i < 25; i++ {
		posts[i] = post.Post{
			Title: fmt.Sprintf("Tools Post %d", i+1),
			Slug:  fmt.Sprintf("tools-%d", i+1),
			Tags:  []string{"tools"},
			Date:  now.Add(-time.Duration(i) * time.Minute),
		}
	}

	h := &Handler{FilePosts: posts}

	// 1. Test ?tag=tools route on BlogIndex
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?tag=tools", nil)
	h.BlogIndex(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "page=2") {
		t.Errorf("expected pagination link with page=2 in ?tag=tools response")
	}

	// 2. Test /blog/tag/tools route on TagIndex
	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/blog/tag/tools", nil)
	// Add chi route context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("tag", "tools")
	req2 = req2.WithContext(context.WithValue(req2.Context(), chi.RouteCtxKey, rctx))
	h.TagIndex(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec2.Code)
	}
	body2 := rec2.Body.String()
	if !strings.Contains(body2, "page=2") {
		t.Errorf("expected pagination link with page=2 in /blog/tag/tools response")
	}
}
