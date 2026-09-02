package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"daemontalk/internal/auth"
	"daemontalk/internal/comment"
	"daemontalk/internal/forum"
	"github.com/go-chi/chi/v5"
)

func TestAuthEndpoints(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "auth.db")
	authDB, _ := auth.Open(dbPath)
	forumDB, _ := forum.Open(dbPath)
	commDB, _ := comment.Open(filepath.Join(tmpDir, "comments.db"))
	t.Cleanup(func() {
		commDB.Close()
		forumDB.Close()
		authDB.Close()
	})

	u, _ := authDB.UpsertUser(auth.User{
		Provider:    "github",
		ProviderID:  "12345",
		Username:    "octocat",
		DisplayName: "Octocat",
		AvatarURL:   "https://github.com/octocat.png",
		GitHubURL:   "https://github.com/octocat",
		Role:        "member",
	})

	// Add forum topic & reply
	_, _ = forumDB.CreateTopic(forum.Topic{
		UserID:   u.ID,
		Title:    "Understanding Linux Schedulers",
		Category: "kernel",
		BodyMD:   "Deep dive into EEVDF.",
	})

	// Add comment
	_, _ = commDB.AddAdvanced(comment.Comment{
		PostSlug: "os-oom-killer",
		Name:     u.DisplayName,
		Body:     "Great explanation of oom_score!",
		UserID:   &u.ID,
	})

	token := "session-token-for-test-999"
	tokenHash := auth.HashToken(token)
	_, _ = authDB.CreateSession(u.ID, tokenHash, 24*time.Hour)

	h := &Handler{
		Auth:     authDB,
		Forum:    forumDB,
		Comments: commDB,
	}
	r := chi.NewRouter()
	r.Use(h.AuthMiddleware)
	r.Get("/auth/me", h.AuthMe)
	r.Get("/auth/logout", h.AuthLogout)
	r.Get("/auth/export", h.AuthExport)
	r.Post("/auth/delete-account", h.AuthDeleteAccount)

	// 1. Test /auth/me unauthenticated
	reqUnauth := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	recUnauth := httptest.NewRecorder()
	r.ServeHTTP(recUnauth, reqUnauth)

	if recUnauth.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for unauthenticated /auth/me, got %d", recUnauth.Code)
	}

	// 2. Test /auth/me authenticated with session cookie
	reqAuth := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	reqAuth.AddCookie(&http.Cookie{
		Name:  auth.SessionCookieName,
		Value: token,
	})
	recAuth := httptest.NewRecorder()
	r.ServeHTTP(recAuth, reqAuth)

	if recAuth.Code != http.StatusOK {
		t.Fatalf("expected 200 for authenticated /auth/me, got %d", recAuth.Code)
	}

	var meResp map[string]any
	if err := json.NewDecoder(recAuth.Body).Decode(&meResp); err != nil {
		t.Fatalf("failed to parse /auth/me json response: %v", err)
	}
	if meResp["authenticated"] != true {
		t.Errorf("expected authenticated true")
	}

	// 3. Test /auth/export
	reqExport := httptest.NewRequest(http.MethodGet, "/auth/export", nil)
	reqExport.AddCookie(&http.Cookie{
		Name:  auth.SessionCookieName,
		Value: token,
	})
	recExport := httptest.NewRecorder()
	r.ServeHTTP(recExport, reqExport)

	if recExport.Code != http.StatusOK {
		t.Fatalf("expected 200 for /auth/export, got %d", recExport.Code)
	}
	if !strings.Contains(recExport.Header().Get("Content-Disposition"), "daemontalk-data-octocat.json") {
		t.Errorf("expected attachment filename in export header, got %s", recExport.Header().Get("Content-Disposition"))
	}
	var exportData map[string]any
	if err := json.NewDecoder(recExport.Body).Decode(&exportData); err != nil {
		t.Fatalf("failed to parse export json: %v", err)
	}
	if exportData["platform"] != "daemontalk" {
		t.Errorf("expected platform daemontalk in export")
	}

	// 4. Test /auth/delete-account
	reqDelete := httptest.NewRequest(http.MethodPost, "/auth/delete-account", nil)
	reqDelete.AddCookie(&http.Cookie{
		Name:  auth.SessionCookieName,
		Value: token,
	})
	recDelete := httptest.NewRecorder()
	r.ServeHTTP(recDelete, reqDelete)

	if recDelete.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 for /auth/delete-account, got %d", recDelete.Code)
	}

	// Verify session and user deleted in DB
	deletedSession, _ := authDB.GetSessionUser(tokenHash)
	if deletedSession != nil {
		t.Errorf("expected session to be deleted from auth db")
	}
	if count := authDB.CountUsers(); count != 0 {
		t.Errorf("expected 0 users after deletion, got %d", count)
	}

	// Verify forum topics still exist but author is anonymized
	topics, _, err := forumDB.ListTopics("", "", "", "", "", 10, 0, 0)
	if err != nil || len(topics) == 0 {
		t.Fatalf("expected forum topic to remain preserved")
	}
	if topics[0].UserID != 0 {
		t.Errorf("expected topic user_id to be anonymized to 0, got %d", topics[0].UserID)
	}

	// 5. Test /auth/logout
	reqLogout := httptest.NewRequest(http.MethodGet, "/auth/logout", nil)
	recLogout := httptest.NewRecorder()
	r.ServeHTTP(recLogout, reqLogout)

	if recLogout.Code != http.StatusSeeOther {
		t.Errorf("expected 303 for logout, got %d", recLogout.Code)
	}
}
