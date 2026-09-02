package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGitHubWebhook(t *testing.T) {
	h := &Handler{
		ContentDir: "../../content",
	}

	payload := `{"ref": "refs/heads/main"}`
	req := httptest.NewRequest(http.MethodPost, "/api/webhook/github", bytes.NewBufferString(payload))
	req.Header.Set("X-GitHub-Event", "push")
	rec := httptest.NewRecorder()

	h.GitHubWebhook(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d (body: %s)", rec.Code, rec.Body.String())
	}
}
