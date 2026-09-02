package router

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"daemontalk/internal/auth"
	"daemontalk/internal/comment"
	"daemontalk/internal/forum"
	"daemontalk/internal/handler"
	"daemontalk/internal/post"
	"daemontalk/internal/postdb"
	"daemontalk/internal/project"
)

func newTestHandler(t *testing.T) *handler.Handler {
	tmpDir := t.TempDir()

	commentStore, err := comment.Open(filepath.Join(tmpDir, "test_comments.db"))
	if err != nil {
		t.Fatalf("failed to open test comments db: %v", err)
	}

	authStore, err := auth.Open(filepath.Join(tmpDir, "test_auth.db"))
	if err != nil {
		t.Fatalf("failed to open test auth db: %v", err)
	}

	forumStore, err := forum.Open(filepath.Join(tmpDir, "test_forum.db"))
	if err != nil {
		t.Fatalf("failed to open test forum db: %v", err)
	}

	postDBStore, err := postdb.Open(filepath.Join(tmpDir, "test_posts.db"))
	if err != nil {
		t.Fatalf("failed to open test postdb: %v", err)
	}

	t.Cleanup(func() {
		_ = commentStore.Close()
		_ = authStore.Close()
		_ = forumStore.Close()
		_ = postDBStore.Close()
	})

	h := &handler.Handler{
		ContentDir:  "content",
		AllProjects: project.All,
		FilePosts: []post.Post{
			{Title: "Test Dispatch", Slug: "test-dispatch", Description: "Test summary"},
		},
		PostDB:   postDBStore,
		Comments: commentStore,
		Auth:     authStore,
		Forum:    forumStore,
	}
	h.RefreshPosts()
	return h
}

func TestRouterEndpoints(t *testing.T) {
	h := newTestHandler(t)
	r := New(h)

	tests := []struct {
		method         string
		target         string
		expectedStatus int
		expectedHeader string
		expectedLoc    string
	}{
		{"GET", "/healthz", http.StatusOK, "", ""},
		{"GET", "/robots.txt", http.StatusOK, "", ""},
		{"GET", "/manifest.json", http.StatusOK, "", ""},
		{"GET", "/rss.xml", http.StatusOK, "", ""},
		{"GET", "/sitemap.xml", http.StatusOK, "", ""},
		{"GET", "/", http.StatusOK, "", ""},
		{"GET", "/colophon", http.StatusOK, "", ""},
		{"GET", "/blog", http.StatusMovedPermanently, "Location", "/"},
		{"GET", "/projects", http.StatusMovedPermanently, "Location", "/colophon#projects"},
		{"GET", "/uses", http.StatusMovedPermanently, "Location", "/colophon"},
		{"GET", "/now", http.StatusMovedPermanently, "Location", "/colophon"},
		{"GET", "/terminal", http.StatusMovedPermanently, "Location", "/"},
		{"GET", "/links", http.StatusNotFound, "", ""},
		{"GET", "/nonexistent-route-12345", http.StatusNotFound, "", ""},
	}

	for _, tc := range tests {
		t.Run(tc.method+" "+tc.target, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.target, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, rec.Code)
			}
			if tc.expectedLoc != "" {
				loc := rec.Header().Get(tc.expectedHeader)
				if loc != tc.expectedLoc {
					t.Errorf("expected Location %q, got %q", tc.expectedLoc, loc)
				}
			}
		})
	}
}

func TestSecurityHeadersMiddleware(t *testing.T) {
	h := newTestHandler(t)
	r := New(h)

	req := httptest.NewRequest("GET", "/healthz", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	headers := []struct {
		header   string
		expected string
	}{
		{"X-Content-Type-Options", "nosniff"},
		{"X-Frame-Options", "DENY"},
		{"Referrer-Policy", "strict-origin-when-cross-origin"},
	}

	for _, h := range headers {
		if val := rec.Header().Get(h.header); val != h.expected {
			t.Errorf("expected header %s: %q, got %q", h.header, h.expected, val)
		}
	}
}
