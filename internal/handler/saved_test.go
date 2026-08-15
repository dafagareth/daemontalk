package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSavedPage(t *testing.T) {
	h := &Handler{}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/saved", nil)
	h.Saved(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "saved-posts") && !strings.Contains(body, "saved") {
		t.Error("expected saved page elements to be in response")
	}
}
