package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAccessibilityPage(t *testing.T) {
	h := &Handler{
		ContentDir: "../../content",
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/accessibility", nil)
	h.Accessibility(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Accessibility Statement") {
		t.Error("expected 'Accessibility Statement' in page")
	}
	if !strings.Contains(body, "WCAG 2.1 Level AA") {
		t.Error("expected WCAG standard mention")
	}
}
