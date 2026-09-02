package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"daemontalk/internal/post"
)

func TestStatsPage(t *testing.T) {
	posts := []post.Post{
		{
			Title:    "Stats Post 1",
			Slug:     "stats-1",
			Tags:     []string{"go"},
			Lang:     "en",
			Date:     time.Now(),
			ReadTime: 5,
		},
	}
	h := &Handler{
		FilePosts: posts,
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	h.Stats(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	// The stats page displays numbers and summaries. We can assert strings related to stats are present.
	if !strings.Contains(body, "Total Articles Published") {
		t.Error("expected statistics elements in Stats response page")
	}
}
