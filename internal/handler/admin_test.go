package handler

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"daemontalk/internal/post"
	"daemontalk/internal/postdb"
	"daemontalk/web/templates"
)

func TestIsAdmin(t *testing.T) {
	h := &Handler{AdminToken: "secret"}

	// No cookie → not admin
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	if h.isAdmin(req) {
		t.Error("no cookie should not be admin")
	}

	// Wrong token → not admin
	req = httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.AddCookie(&http.Cookie{Name: "admin_token", Value: "wrong"})
	if h.isAdmin(req) {
		t.Error("wrong token should not be admin")
	}

	// Correct token → admin
	req = httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.AddCookie(&http.Cookie{Name: "admin_token", Value: "secret"})
	if !h.isAdmin(req) {
		t.Error("correct token should be admin")
	}
}

func TestIsAdminDisabledWhenNoToken(t *testing.T) {
	h := &Handler{AdminToken: ""} // moderation disabled

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.AddCookie(&http.Cookie{Name: "admin_token", Value: "anything"})
	if h.isAdmin(req) {
		t.Error("when AdminToken is empty, nobody should be admin")
	}
}

func TestAdminForbiddenWithoutAuth(t *testing.T) {
	h := &Handler{AdminToken: "secret"}

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	h.Admin(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("unauthenticated /admin: got %d, want 404", rec.Code)
	}
}

func TestAdminLoginSetsCookie(t *testing.T) {
	h := &Handler{AdminToken: "secret"}

	req := httptest.NewRequest(http.MethodGet, "/admin?admin=secret", nil)
	rec := httptest.NewRecorder()
	h.Admin(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("login should redirect (303), got %d", rec.Code)
	}
	var found bool
	for _, c := range rec.Result().Cookies() {
		if c.Name == "admin_token" && c.Value == "secret" {
			found = true
		}
	}
	if !found {
		t.Error("login with correct token should set admin_token cookie")
	}
}

func TestAdminDeleteForbiddenWithoutAuth(t *testing.T) {
	h := &Handler{AdminToken: "secret"}

	req := httptest.NewRequest(http.MethodPost, "/admin/comments/1/delete", nil)
	rec := httptest.NewRecorder()
	h.AdminDeleteComment(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("unauthenticated delete: got %d, want 404", rec.Code)
	}
}

func TestConfirmModalMarkupPresent(t *testing.T) {
	h := &Handler{AdminToken: "secret"}

	// 403 Page is also wrapped in the same Layout() — just to make sure
	// modal konfirmasi global (dipakai admin & blog) benar-benar terpasang.
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	h.Admin(rec, req)

	body := rec.Body.String()
	for _, want := range []string{
		`id="confirm-overlay"`,
		`id="confirm-message"`,
		"data-confirm-ok",
		"data-confirm-cancel",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("halaman harus memuat modal konfirmasi global, tidak ditemukan %q", want)
		}
	}
}

func TestAdminPageIsHumanized(t *testing.T) {
	pdb, err := postdb.Open(filepath.Join(t.TempDir(), "posts.db"))
	if err != nil {
		t.Fatalf("open postdb: %v", err)
	}
	if _, err := pdb.Create(postdb.WebPost{
		Slug: "tulisan-web", Title: "Tulisan dari Web", BodyMD: "isi",
		Lang: "id", Date: "2026-07-01",
	}); err != nil {
		t.Fatalf("create webpost: %v", err)
	}

	h := &Handler{
		AdminToken: "secret",
		PostDB:     pdb,
		FilePosts: []post.Post{
			{Title: "Belajar Go dari Nol", Slug: "belajar-go-dari-nol", Date: time.Date(2026, 7, 2, 0, 0, 0, 0, time.UTC)},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.AddCookie(&http.Cookie{Name: "admin_token", Value: "secret"})
	rec := httptest.NewRecorder()
	h.Admin(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /admin authenticated: got %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{
		"Channel Dashboard",
		"Channel Content",
		"Belajar Go dari Nol",
		"No comments yet.",
		"content/posts/*.md",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("dashboard body missing %q", want)
		}
	}
	for _, notWant := range []string{"Halo, Dafa", "Semua tulisan", "Belum ada komentar"} {
		if strings.Contains(body, notWant) {
			t.Errorf("dashboard body should not contain %q", notWant)
		}
	}
}

func TestAdminPageDoesNotIncludePublicNav(t *testing.T) {
	h := &Handler{AdminToken: "secret"}

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.AddCookie(&http.Cookie{Name: "admin_token", Value: "secret"})
	rec := httptest.NewRecorder()
	h.Admin(rec, req)

	body := rec.Body.String()
	for _, notWant := range []string{`href="/colophon"`, `href="/search"`, "buymeacoffee.com"} {
		if strings.Contains(body, notWant) {
			t.Errorf("admin shell should not include public nav/footer, found %q", notWant)
		}
	}
}

func TestAdminPageIncludesMinimalShell(t *testing.T) {
	h := &Handler{AdminToken: "secret"}

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.AddCookie(&http.Cookie{Name: "admin_token", Value: "secret"})
	rec := httptest.NewRecorder()
	h.Admin(rec, req)

	body := rec.Body.String()
	for _, want := range []string{"View Site", "toggleTheme()", `name="robots" content="noindex,nofollow"`, `data-tab="dashboard"`, `data-tab="content"`} {
		if !strings.Contains(body, want) {
			t.Errorf("admin shell body missing %q", want)
		}
	}
}

func TestPublicPageStillHasPublicNavAndFooter(t *testing.T) {
	h := &Handler{}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.BlogIndex(rec, req)

	body := rec.Body.String()
	for _, want := range []string{`href="/about"`, "github.com", `hreflang="en"`} {
		if !strings.Contains(body, want) {
			t.Errorf("public page should still include %q after headTags() extraction", want)
		}
	}
	if strings.Contains(body, `name="robots" content="noindex,nofollow"`) {
		t.Error("public page should not be noindex")
	}
}

func TestAdminUnauthorizedReturnsNotFoundPage(t *testing.T) {
	h := &Handler{AdminToken: "secret"}

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	h.Admin(rec, req)

	body := rec.Body.String()
	if !strings.Contains(body, "404") {
		t.Error("Unauthorized admin access should return 404 error page")
	}
}

func TestAdminToggleRadar(t *testing.T) {
	h := &Handler{AdminToken: "secret"}

	// Initially disabled
	templates.SetRadarEnabled(false)

	// 1. Check Graph returns 404 when disabled
	reqGraph := httptest.NewRequest(http.MethodGet, "/graph", nil)
	recGraph := httptest.NewRecorder()
	h.Graph(recGraph, reqGraph)
	if recGraph.Code != http.StatusNotFound {
		t.Fatalf("expected 404 when radar is disabled, got %d", recGraph.Code)
	}

	// 2. Toggle to enabled
	reqToggle := httptest.NewRequest(http.MethodPost, "/admin/settings/toggle-radar", nil)
	reqToggle.AddCookie(&http.Cookie{Name: "admin_token", Value: "secret"})
	recToggle := httptest.NewRecorder()
	h.AdminToggleRadar(recToggle, reqToggle)

	if recToggle.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 on toggle, got %d", recToggle.Code)
	}
	if !templates.IsRadarEnabled() {
		t.Fatal("expected radar to be enabled after toggle")
	}

	// 3. Check Graph returns 200 when enabled
	reqGraph2 := httptest.NewRequest(http.MethodGet, "/graph", nil)
	recGraph2 := httptest.NewRecorder()
	h.Graph(recGraph2, reqGraph2)
	if recGraph2.Code != http.StatusOK {
		t.Fatalf("expected 200 when radar is enabled, got %d", recGraph2.Code)
	}

	// 4. Toggle back to disabled
	recToggle2 := httptest.NewRecorder()
	h.AdminToggleRadar(recToggle2, reqToggle)
	if templates.IsRadarEnabled() {
		t.Fatal("expected radar to be disabled after second toggle")
	}
}
