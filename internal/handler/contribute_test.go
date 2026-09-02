package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestContributePage(t *testing.T) {
	h := &Handler{
		ContentDir: "../../content",
	}

	req := httptest.NewRequest(http.MethodGet, "/contribute", nil)
	rec := httptest.NewRecorder()
	h.Contribute(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /contribute: got status %d, want 200", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Download daemontalk-template.md") {
		t.Error("body missing 'Download daemontalk-template.md' button")
	}
}

func TestDownloadTemplate(t *testing.T) {
	h := &Handler{
		ContentDir: "../../content",
	}

	req := httptest.NewRequest(http.MethodGet, "/download/daemontalk-template.md", nil)
	rec := httptest.NewRecorder()
	h.DownloadTemplate(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /download/daemontalk-template.md: got status %d, want 200", rec.Code)
	}

	disp := rec.Header().Get("Content-Disposition")
	if !strings.Contains(disp, "attachment") || !strings.Contains(disp, "daemontalk-template.md") {
		t.Errorf("expected attachment header for daemontalk-template.md, got %q", disp)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "title:") || !strings.Contains(body, "slug:") {
		t.Error("downloaded file body does not contain expected template frontmatter")
	}
}
