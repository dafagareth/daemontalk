package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNotFoundHandler(t *testing.T) {
	h := &Handler{}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/random-unknown-page", nil)
	h.NotFound(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "404") {
		t.Error("expected body to contain 404")
	}
	if !strings.Contains(body, "Page not found") {
		t.Error("expected body to contain Page not found")
	}
	if !strings.Contains(body, "Go home") {
		t.Error("expected body to contain Go home button")
	}
}

func TestNotFoundIndonesian(t *testing.T) {
	h := &Handler{}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/id/halaman-ngaco", nil)
	h.NotFound(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "404") {
		t.Error("expected body to contain 404")
	}
	if !strings.Contains(body, "Halaman tidak ditemukan") {
		t.Error("expected body to contain Indonesian not found title")
	}
	if !strings.Contains(body, "Ke beranda") {
		t.Error("expected body to contain Indonesian home button")
	}
}
