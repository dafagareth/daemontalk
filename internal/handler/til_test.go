package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"daemontalk/internal/post"
)

func TestTILPage(t *testing.T) {
	posts := []post.Post{
		{
			Title: "Today I Learned Rust Lifetimes",
			Slug:  "til-rust-lifetimes",
			Tags:  []string{"rust"},
			Lang:  "en",
			Type:  "til",
			Date:  time.Now(),
			Body:  "<div>lifetimes explanation</div>",
		},
		{
			Title: "Regular Blog Post",
			Slug:  "regular-blog-post",
			Tags:  []string{"go"},
			Lang:  "en",
			Date:  time.Now(),
			Body:  "<div>blog post</div>",
		},
	}
	h := &Handler{
		FilePosts: posts,
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/til", nil)
	h.TIL(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Today I Learned Rust Lifetimes") {
		t.Error("expected TIL post title to be in rendered page")
	}
	if strings.Contains(body, "Regular Blog Post") {
		t.Error("did not expect regular blog post in TIL page")
	}
}
