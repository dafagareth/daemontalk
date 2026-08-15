package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"daemontalk/internal/post"
)

func TestTerminalPage(t *testing.T) {
	h := &Handler{}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/terminal", nil)
	h.Terminal(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "terminalContent") && !strings.Contains(body, "term-window") {
		t.Error("expected terminal frame or window elements in page response")
	}
}

func TestTerminalDataEndpoint(t *testing.T) {
	posts := []post.Post{
		{
			Title:       "Terminal Test Post",
			Slug:        "terminal-test",
			Description: "A description of terminal behavior",
			Tags:        []string{"unix", "test"},
			Date:        time.Now(),
			Body:        "<div>Hello terminal world</div>",
			Draft:       false,
		},
	}
	h := &Handler{
		FilePosts: posts,
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/terminal/data", nil)
	h.TerminalData(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("expected application/json content type, got %q", contentType)
	}

	var data terminalDataResponse
	if err := json.NewDecoder(rec.Body).Decode(&data); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	if data.Host != "daemontalk.local" {
		t.Errorf("expected host daemontalk.local, got %q", data.Host)
	}

	if len(data.Posts) != 1 || data.Posts[0].Slug != "terminal-test" {
		t.Errorf("expected post 'terminal-test' in response, got %+v", data.Posts)
	}

	if data.Tags["unix"] != 1 || data.Tags["test"] != 1 {
		t.Errorf("expected tag counts, got %+v", data.Tags)
	}
}
