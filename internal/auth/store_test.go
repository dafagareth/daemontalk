package auth

import (
	"context"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"
)

func TestAuthStore_UserAndSession(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_auth.db")

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open auth store: %v", err)
	}
	defer store.Close()

	// 1. Upsert User
	u1 := User{
		Provider:    "github",
		ProviderID:  "123456",
		Username:    "octocat",
		DisplayName: "Mona Lisa Octocat",
		Email:       "octocat@github.com",
		AvatarURL:   "https://github.com/images/error/octocat_happy.gif",
		GitHubURL:   "https://github.com/octocat",
		Role:        "member",
	}

	savedUser, err := store.UpsertUser(u1)
	if err != nil {
		t.Fatalf("failed to upsert user: %v", err)
	}
	if savedUser.ID <= 0 {
		t.Fatalf("expected positive user ID, got %d", savedUser.ID)
	}
	if savedUser.Username != "octocat" {
		t.Fatalf("expected username octocat, got %s", savedUser.Username)
	}

	// 2. Count Users
	if count := store.CountUsers(); count != 1 {
		t.Errorf("expected 1 user, got %d", count)
	}

	// 3. Get User By Username
	byName, err := store.GetUserByUsername("octocat")
	if err != nil || byName == nil {
		t.Fatalf("expected to find user by username octocat, got %v (err: %v)", byName, err)
	}
	if byName.ID != savedUser.ID {
		t.Errorf("expected ID %d, got %d", savedUser.ID, byName.ID)
	}

	// Non-existent username
	_, err = store.GetUserByUsername("nonexistent")
	if err == nil {
		t.Errorf("expected error for nonexistent user")
	}

	// 4. Create Session
	token := "raw-secure-token-12345"
	tokenHash := HashToken(token)
	session, err := store.CreateSession(savedUser.ID, tokenHash, 24*time.Hour)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	if session.UserID != savedUser.ID {
		t.Errorf("expected session user id %d, got %d", savedUser.ID, session.UserID)
	}

	// 5. Get Session User
	sessionUser, err := store.GetSessionUser(tokenHash)
	if err != nil || sessionUser == nil {
		t.Fatalf("failed to get session user: %v", err)
	}
	if sessionUser.Username != "octocat" {
		t.Errorf("expected session user octocat, got %s", sessionUser.Username)
	}

	// Expired session test
	expiredTokenHash := HashToken("expired-token")
	_, err = store.CreateSession(savedUser.ID, expiredTokenHash, -1*time.Hour)
	if err != nil {
		t.Fatalf("failed to create expired session: %v", err)
	}
	expiredUser, err := store.GetSessionUser(expiredTokenHash)
	if err != nil || expiredUser != nil {
		t.Errorf("expected nil for expired session user, got %v (err: %v)", expiredUser, err)
	}

	// 6. Delete Session
	if err := store.DeleteSession(tokenHash); err != nil {
		t.Fatalf("failed to delete session: %v", err)
	}
	sessionUser, err = store.GetSessionUser(tokenHash)
	if err != nil || sessionUser != nil {
		t.Fatalf("expected nil user after session deletion, got %v (err: %v)", sessionUser, err)
	}

	// 7. Delete User
	if err := store.DeleteUser(savedUser.ID); err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}
	if count := store.CountUsers(); count != 0 {
		t.Errorf("expected 0 users after delete, got %d", count)
	}
}

func TestAuthContext(t *testing.T) {
	// Nil context
	if u := GetUser(nil); u != nil {
		t.Errorf("expected nil user from nil context")
	}

	// Empty context
	ctx := context.Background()
	if u := GetUser(ctx); u != nil {
		t.Errorf("expected nil user from empty context")
	}

	// With user
	user := &User{ID: 42, Username: "gopher", Role: "admin"}
	ctx = WithUser(ctx, user)
	retrieved := GetUser(ctx)
	if retrieved == nil || retrieved.ID != 42 || retrieved.Username != "gopher" {
		t.Errorf("expected user with ID 42, got %v", retrieved)
	}
}

func TestSessionHelpers(t *testing.T) {
	tok1, err := GenerateRandomToken()
	if err != nil || len(tok1) == 0 {
		t.Fatalf("failed to generate random token: %v", err)
	}
	tok2, err := GenerateRandomToken()
	if err != nil || tok1 == tok2 {
		t.Fatalf("expected unique random tokens")
	}

	h1 := HashToken(tok1)
	h2 := HashToken(tok1)
	if h1 != h2 {
		t.Errorf("hash must be deterministic")
	}

	// Cookie test
	rec := httptest.NewRecorder()
	SetSessionCookie(rec, tok1, false)

	resp := rec.Result()
	cookies := resp.Cookies()
	if len(cookies) == 0 || cookies[0].Name != SessionCookieName || cookies[0].Value != tok1 {
		t.Errorf("expected session cookie with token, got cookies: %v", cookies)
	}

	// Extract from request
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(cookies[0])
	extracted := GetSessionTokenFromRequest(req)
	if extracted != tok1 {
		t.Errorf("expected extracted token %s, got %s", tok1, extracted)
	}

	// Clear cookie test
	rec2 := httptest.NewRecorder()
	ClearSessionCookie(rec2, false)
	clearCookies := rec2.Result().Cookies()
	if len(clearCookies) == 0 || clearCookies[0].MaxAge != -1 {
		t.Errorf("expected cleared session cookie with MaxAge -1")
	}
}
