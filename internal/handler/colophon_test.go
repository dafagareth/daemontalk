package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"daemontalk/internal/project"
)

func TestColophonPage(t *testing.T) {
	h := &Handler{
		AllProjects: []project.Project{
			{Name: "svault", Slug: "svault", Description: "Local encrypted secret vault.", Status: "active", Featured: true},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/colophon", nil)
	rec := httptest.NewRecorder()
	h.Colophon(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /colophon: got %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{"Colophon", "svault", "id=\"projects\"", "Platform &amp; Host"} {
		if !strings.Contains(body, want) {
			t.Errorf("body missing %q", want)
		}
	}
}

func TestColophonPageID(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodGet, "/id/colophon", nil)
	rec := httptest.NewRecorder()
	h.Colophon(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /id/colophon: got %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Spesifikasi & Kolofon") {
		t.Error("halaman /id/colophon harus merender dengan benar")
	}
}
