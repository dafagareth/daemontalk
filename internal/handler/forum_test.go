package handler

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"daemontalk/internal/auth"
	"daemontalk/internal/forum"
	"github.com/go-chi/chi/v5"
)

func TestDiscussionsHandler(t *testing.T) {
	tmpDir := t.TempDir()
	authDB, err := auth.Open(filepath.Join(tmpDir, "auth.db"))
	if err != nil {
		t.Fatalf("open auth db: %v", err)
	}
	defer authDB.Close()

	forumDB, err := forum.Open(filepath.Join(tmpDir, "forum.db"))
	if err != nil {
		t.Fatalf("open forum db: %v", err)
	}
	defer forumDB.Close()

	u, _ := authDB.UpsertUser(auth.User{
		Provider:    "github",
		ProviderID:  "999",
		Username:    "daemontalk_dev",
		DisplayName: "Daemontalk Dev",
		AvatarURL:   "https://github.com/daemontalk.png",
		GitHubURL:   "https://github.com/daemontalk",
		Role:        "member",
	})

	_, _ = forumDB.CreateTopic(forum.Topic{
		UserID:   u.ID,
		Title:    "Membongkar Concurrency Memory Model di Go",
		Category: "go",
		Tags:     []string{"go", "concurrency"},
		BodyMD:   "Bagaimana cara kerja atomic memory ordering di Go runtime?",
	})

	h := &Handler{
		Auth:  authDB,
		Forum: forumDB,
	}

	r := chi.NewRouter()
	r.Use(h.AuthMiddleware)
	r.Get("/socket", h.Discussions)
	r.Get("/socket/new", h.DiscussionsNew)
	r.Get("/socket/{slug}", h.DiscussionsDetail)
	r.Get("/discussions", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/socket", http.StatusMovedPermanently)
	})
	r.Get("/guestbook", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/socket", http.StatusMovedPermanently)
	})

	// 1. Test /socket
	req := httptest.NewRequest(http.MethodGet, "/socket", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 for /socket, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Socket") {
		t.Errorf("expected Socket page title in HTML")
	}
	if !strings.Contains(body, "Membongkar Concurrency Memory Model di Go") {
		t.Errorf("expected seeded topic title to appear in socket list")
	}

	// 2. Test /discussions 301 Redirect to /socket
	discReq := httptest.NewRequest(http.MethodGet, "/discussions", nil)
	discRec := httptest.NewRecorder()
	r.ServeHTTP(discRec, discReq)
	if discRec.Code != http.StatusMovedPermanently || discRec.Header().Get("Location") != "/socket" {
		t.Errorf("expected 301 to /socket, got %d, loc: %s", discRec.Code, discRec.Header().Get("Location"))
	}

	// 3. Test /guestbook 301 Redirect to /socket
	gbReq := httptest.NewRequest(http.MethodGet, "/guestbook", nil)
	gbRec := httptest.NewRecorder()
	r.ServeHTTP(gbRec, gbReq)

	if gbRec.Code != http.StatusMovedPermanently {
		t.Errorf("expected status 301 for /guestbook, got %d", gbRec.Code)
	}
	if loc := gbRec.Header().Get("Location"); loc != "/socket" {
		t.Errorf("expected Location /socket, got %s", loc)
	}

	// 4. Test /socket/new
	newReq := httptest.NewRequest(http.MethodGet, "/socket/new", nil)
	newRec := httptest.NewRecorder()
	r.ServeHTTP(newRec, newReq)

	if newRec.Code != http.StatusOK {
		t.Errorf("expected status 200 for /socket/new, got %d", newRec.Code)
	}

	// 5. Test /socket/{slug} detail
	detailReq := httptest.NewRequest(http.MethodGet, "/socket/membongkar-concurrency-memory-model-di-go", nil)
	detailRec := httptest.NewRecorder()
	r.ServeHTTP(detailRec, detailReq)

	if detailRec.Code != http.StatusOK {
		t.Errorf("expected status 200 for topic detail, got %d", detailRec.Code)
	}
	detailBody := detailRec.Body.String()
	if !strings.Contains(detailBody, "atomic memory ordering") {
		t.Errorf("expected topic body markdown rendered in HTML")
	}
}

func TestDiscussionsDetail_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	forumDB, _ := forum.Open(filepath.Join(tmpDir, "forum.db"))
	defer forumDB.Close()

	h := &Handler{Forum: forumDB}
	r := chi.NewRouter()
	r.Get("/socket/{slug}", h.DiscussionsDetail)

	req := httptest.NewRequest(http.MethodGet, "/socket/non-existent-slug-xyz", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}
}
