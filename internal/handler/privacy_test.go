package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPrivacyPage(t *testing.T) {
	h := &Handler{
		ContentDir: "../../content",
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/privacy", nil)
	h.Privacy(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Privacy Policy") {
		t.Error("expected 'Privacy Policy' title in english privacy page")
	}
	if !strings.Contains(body, "zero-tracker") {
		t.Error("expected zero-tracker declaration")
	}

	// Test Indonesian locale
	recID := httptest.NewRecorder()
	reqID := httptest.NewRequest(http.MethodGet, "/id/privacy", nil)
	h.Privacy(recID, reqID)

	bodyID := recID.Body.String()
	if !strings.Contains(bodyID, "Kebijakan Privasi") {
		t.Error("expected 'Kebijakan Privasi' title in indonesian privacy page")
	}
	if !strings.Contains(bodyID, "zero-tracker") {
		t.Error("expected indonesian zero-tracker description")
	}
}
