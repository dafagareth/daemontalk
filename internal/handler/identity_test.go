package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetVisitorIdentitySetsCookieAndGeneratesAnonymName(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/guestbook", nil)
	rec := httptest.NewRecorder()

	name1 := GetVisitorIdentity(rec, req)
	if !strings.HasPrefix(name1, "anonym_") {
		t.Errorf("expected anonym_ prefix, got %q", name1)
	}

	// Verify cookie was set
	cookies := rec.Result().Cookies()
	var visitorCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "visitor_id" {
			visitorCookie = c
			break
		}
	}
	if visitorCookie == nil {
		t.Fatal("visitor_id cookie was not set")
	}

	// Subsequent request with the same cookie must produce the EXACT same anonymous name
	req2 := httptest.NewRequest(http.MethodGet, "/guestbook", nil)
	req2.AddCookie(visitorCookie)
	rec2 := httptest.NewRecorder()

	name2 := GetVisitorIdentity(rec2, req2)
	if name1 != name2 {
		t.Errorf("expected identical name %q across requests with same cookie, got %q", name1, name2)
	}
}
