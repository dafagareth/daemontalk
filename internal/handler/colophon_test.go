package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestColophonPage(t *testing.T) {
	h := &Handler{}

	req := httptest.NewRequest(http.MethodGet, "/colophon", nil)
	rec := httptest.NewRecorder()
	h.Colophon(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /colophon: got %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{"Colophon", "Platform &amp; Host"} {
		if !strings.Contains(body, want) {
			t.Errorf("body missing %q", want)
		}
	}
}
