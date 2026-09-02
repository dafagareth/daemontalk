package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"daemontalk/internal/post"
	"daemontalk/internal/postdb"
	"github.com/go-chi/chi/v5"
)

func TestAdminPostEditForbiddenWithoutAuth(t *testing.T) {
	h := &Handler{AdminToken: "secret"}

	req := httptest.NewRequest(http.MethodGet, "/admin/posts/1/edit", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()
	h.AdminPostEdit(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("unauthenticated GET /admin/posts/1/edit: got %d, want 404", rec.Code)
	}
}

func newAdminAuthedHandler(t *testing.T) *Handler {
	t.Helper()
	pdb, err := postdb.Open(filepath.Join(t.TempDir(), "posts.db"))
	if err != nil {
		t.Fatalf("open postdb: %v", err)
	}
	return &Handler{
		AdminToken: "secret",
		PostDB:     pdb,
		FilePosts: []post.Post{
			{Title: "Post File", Slug: "post-file"},
		},
	}
}

type autosaveResp struct {
	ID      int64  `json:"id"`
	Slug    string `json:"slug"`
	SavedAt string `json:"savedAt"`
}

func autosave(t *testing.T, h *Handler, id int64, title, body string) autosaveResp {
	return autosaveWithSlug(t, h, id, title, body, "")
}

func autosaveWithSlug(t *testing.T, h *Handler, id int64, title, body, slug string) autosaveResp {
	t.Helper()
	payload := fmt.Sprintf(`{"id":%d,"title":%q,"body":%q,"slug":%q}`, id, title, body, slug)
	req := httptest.NewRequest(http.MethodPost, "/admin/posts/autosave", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "admin_token", Value: h.AdminToken})
	rec := httptest.NewRecorder()
	h.AdminPostAutosave(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("autosave: got %d, want 200 (body: %s)", rec.Code, rec.Body.String())
	}
	var resp autosaveResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("autosave response bukan JSON: %v (%s)", err, rec.Body.String())
	}
	return resp
}

func publishForm(h *Handler, id int64, form url.Values) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/posts/%d/publish", id), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "admin_token", Value: h.AdminToken})
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", fmt.Sprintf("%d", id))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()
	h.AdminPostPublish(rec, req)
	return rec
}

func TestAdminAutosaveForbiddenWithoutAuth(t *testing.T) {
	h := &Handler{AdminToken: "secret"}

	req := httptest.NewRequest(http.MethodPost, "/admin/posts/autosave", strings.NewReader(`{"id":0,"title":"x","body":"y"}`))
	rec := httptest.NewRecorder()
	h.AdminPostAutosave(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("unauthenticated autosave: got %d, want 404", rec.Code)
	}
}

func TestAdminAutosaveCreatesDraft(t *testing.T) {
	h := newAdminAuthedHandler(t)

	resp := autosave(t, h, 0, "Tulisan Baru", "Isi tulisan dari editor.")
	if resp.ID == 0 {
		t.Fatal("autosave create: id 0, mau id baru")
	}
	if len(resp.Slug) != 8 {
		t.Errorf("slug: dapat %q, mau short ID 8 karakter", resp.Slug)
	}

	p, err := h.PostDB.Get(resp.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if !p.Draft {
		t.Error("post autosave harus Draft")
	}
	if p.Title != "Tulisan Baru" || p.BodyMD != "Isi tulisan dari editor." {
		t.Errorf("isi tidak sesuai: %+v", p)
	}
}

func TestAdminAutosaveSlugConflictGetsSuffix(t *testing.T) {
	h := newAdminAuthedHandler(t)

	// "Post File" slug collides with post file → suffix -2.
	resp := autosaveWithSlug(t, h, 0, "Post File", "Isi apapun.", "post-file")
	if resp.Slug != "post-file-2" {
		t.Errorf("slug bentrok: dapat %q, mau post-file-2", resp.Slug)
	}
}

func TestAdminAutosaveUpdatesExisting(t *testing.T) {
	h := newAdminAuthedHandler(t)

	created := autosave(t, h, 0, "Judul Awal", "Isi awal.")
	updated := autosave(t, h, created.ID, "Judul Baru", "Isi baru.")

	if updated.ID != created.ID {
		t.Errorf("id berubah: %d → %d", created.ID, updated.ID)
	}
	if updated.Slug != created.Slug {
		t.Errorf("slug berubah saat autosave: %q → %q", created.Slug, updated.Slug)
	}
	p, err := h.PostDB.Get(created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if p.Title != "Judul Baru" || p.BodyMD != "Isi baru." {
		t.Errorf("update tidak tersimpan: %+v", p)
	}
}

func TestAdminPostPublishDraft(t *testing.T) {
	h := newAdminAuthedHandler(t)
	created := autosave(t, h, 0, "Tulisan Baru", "Isi tulisan.")

	rec := publishForm(h, created.ID, url.Values{
		"slug":   {"tulisan-baru"},
		"date":   {"2026-07-05"},
		"lang":   {"id"},
		"tags":   {"uji, editor"},
		"action": {"publish"},
	})
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("publish: got %d, want 303 (body: %s)", rec.Code, rec.Body.String())
	}
	if loc := rec.Header().Get("Location"); loc != "/blog/tulisan-baru" {
		t.Errorf("redirect: dapat %q, mau /blog/tulisan-baru", loc)
	}

	p, err := h.PostDB.Get(created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if p.Draft {
		t.Error("setelah publish harus Draft=false")
	}
	found := false
	for _, ap := range h.AllPosts() {
		if ap.Slug == "tulisan-baru" {
			found = true
		}
	}
	if !found {
		t.Error("post publish tidak muncul di AllPosts()")
	}
}

func TestAdminPostPublishSlugConflict(t *testing.T) {
	h := newAdminAuthedHandler(t)
	created := autosave(t, h, 0, "Tulisan Baru", "Isi tulisan.")

	rec := publishForm(h, created.ID, url.Values{
		"slug":   {"post-file"}, // collides with post file
		"date":   {"2026-07-05"},
		"lang":   {"id"},
		"action": {"publish"},
	})
	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("slug bentrok: got %d, want 422", rec.Code)
	}
}

func TestAdminPostPublishSlugLockedOncePublished(t *testing.T) {
	h := newAdminAuthedHandler(t)
	created := autosave(t, h, 0, "Tulisan Baru", "Isi tulisan.")

	rec := publishForm(h, created.ID, url.Values{
		"slug": {"slug-pertama"}, "date": {"2026-07-05"}, "lang": {"id"}, "action": {"publish"},
	})
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("publish pertama: got %d", rec.Code)
	}

	// Already published → slug from form ignored.
	rec = publishForm(h, created.ID, url.Values{
		"slug": {"slug-kedua"}, "date": {"2026-07-05"}, "lang": {"id"}, "action": {"publish"},
	})
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("publish kedua: got %d (body: %s)", rec.Code, rec.Body.String())
	}
	p, err := h.PostDB.Get(created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if p.Slug != "slug-pertama" {
		t.Errorf("slug post published harus terkunci: dapat %q, mau slug-pertama", p.Slug)
	}
}

func TestAdminPostEditorDoesNotIncludePublicNav(t *testing.T) {
	h := newAdminAuthedHandler(t)
	id, _ := h.PostDB.Create(postdb.WebPost{
		Slug: "draft-post", Title: "Draft Post", BodyMD: "draft body",
		Lang: "en", Draft: true, Date: "2026-07-01",
	})

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/posts/%d/edit", id), nil)
	req.AddCookie(&http.Cookie{Name: "admin_token", Value: h.AdminToken})
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", fmt.Sprintf("%d", id))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()
	h.AdminPostEdit(rec, req)

	body := rec.Body.String()
	for _, notWant := range []string{`href="/colophon"`, `href="/search"`, "buymeacoffee.com"} {
		if strings.Contains(body, notWant) {
			t.Errorf("editor shell should not include public nav/footer, found %q", notWant)
		}
	}
	if !strings.Contains(body, "View Site") {
		t.Error("editor shell should include 'View Site' link back to public site")
	}
}

func TestAdminPostEditorDraftLabelsAreEnglish(t *testing.T) {
	h := newAdminAuthedHandler(t)
	id, _ := h.PostDB.Create(postdb.WebPost{
		Slug: "draft-post-2", Title: "Draft Post 2", BodyMD: "draft body",
		Lang: "en", Draft: true, Date: "2026-07-01",
	})

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/posts/%d/edit", id), nil)
	req.AddCookie(&http.Cookie{Name: "admin_token", Value: h.AdminToken})
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", fmt.Sprintf("%d", id))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()
	h.AdminPostEdit(rec, req)

	body := rec.Body.String()
	if !strings.Contains(body, "Publish") {
		t.Error("draft post button should be labeled 'Publish'")
	}
	if !strings.Contains(body, "Draft") {
		t.Error("draft post chip should be 'Draft'")
	}
	if strings.Contains(body, "Terbitkan") || strings.Contains(body, "Simpan sebagai draft") {
		t.Error("editor should not contain Indonesian labels")
	}
}

func TestAdminPostEditorPublishedLabelsAreEnglish(t *testing.T) {
	h := newAdminAuthedHandler(t)
	id, err := h.PostDB.Create(postdb.WebPost{
		Slug: "already-published", Title: "Already Published", BodyMD: "body",
		Lang: "en", Draft: false, Date: "2026-07-01",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/posts/%d/edit", id), nil)
	req.AddCookie(&http.Cookie{Name: "admin_token", Value: h.AdminToken})
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", fmt.Sprintf("%d", id))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()
	h.AdminPostEdit(rec, req)

	body := rec.Body.String()
	if !strings.Contains(body, "Update") {
		t.Error("published post button should be labeled 'Update'")
	}
	if !strings.Contains(body, "Published") {
		t.Error("published post chip should be 'Published'")
	}
	if strings.Contains(body, "Perbarui") || strings.Contains(body, "Terbit") {
		t.Error("editor should not contain Indonesian labels")
	}
}
