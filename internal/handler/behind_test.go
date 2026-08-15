package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"daemontalk/internal/project"
)

func TestBehindPage(t *testing.T) {
	h := &Handler{
		AllProjects: []project.Project{
			{Name: "svault", Slug: "svault", Description: "Local encrypted secret vault.", Status: "active", Featured: true},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/behind", nil)
	rec := httptest.NewRecorder()
	h.Behind(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /behind: got %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{"Behind this website", "svault", "id=\"projects\"", "Languages"} {
		if !strings.Contains(body, want) {
			t.Errorf("body missing %q", want)
		}
	}
}

func TestBehindPageID(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodGet, "/id/behind", nil)
	rec := httptest.NewRecorder()
	h.Behind(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /id/behind: got %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Languages") {
		t.Error("halaman /id/behind harus merender dengan benar")
	}
}
