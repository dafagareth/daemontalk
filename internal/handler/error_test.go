package handler

import (
	"errors"
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

func TestForbiddenHandler(t *testing.T) {
	h := &Handler{}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/secret", nil)
	h.Forbidden(rec, req, "invalid_admin_signature")

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "403") {
		t.Error("expected body to contain 403")
	}
	if !strings.Contains(body, "Access forbidden") {
		t.Error("expected body to contain Access forbidden")
	}
}

func TestForbiddenIndonesian(t *testing.T) {
	h := &Handler{}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/id/admin/secret", nil)
	h.Forbidden(rec, req, "no_permission")

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Akses ditolak") {
		t.Error("expected Indonesian forbidden title")
	}
}

func TestInternalErrorHandler(t *testing.T) {
	h := &Handler{}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/some-broken-endpoint", nil)
	h.InternalError(rec, req, errors.New("db_connection_timed_out"))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "500") {
		t.Error("expected body to contain 500")
	}
	if !strings.Contains(body, "Internal server error") {
		t.Error("expected body to contain Internal server error")
	}
}

func TestCustomErrorHandler(t *testing.T) {
	h := &Handler{}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/rate-limited", nil)
	h.CustomError(rec, req, http.StatusTooManyRequests, "Rate Limit Exceeded", "Please slow down your requests.")

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected status 429, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "429") {
		t.Error("expected body to contain 429")
	}
	if !strings.Contains(body, "Rate Limit Exceeded") {
		t.Error("expected body to contain Rate Limit Exceeded")
	}
}
