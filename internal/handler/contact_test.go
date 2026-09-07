package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestValidEmail(t *testing.T) {
	valid := []string{
		"a@b.co",
		"dafagareth@gmail.com",
		"first.last@sub.domain.io",
	}
	for _, e := range valid {
		if !validEmail(e) {
			t.Errorf("validEmail(%q) = false, want true", e)
		}
	}

	invalid := []string{
		"",
		"notanemail",
		"@nolocal.com",
		"two@@at.com",
		"name <a@b.co>",
		"a@b.co\r\nBcc: x@y.z",
	}
	for _, e := range invalid {
		if validEmail(e) {
			t.Errorf("validEmail(%q) = true, want false", e)
		}
	}
}

func TestStripCRLF(t *testing.T) {
	got := stripCRLF("hello\r\nBcc: evil@x.com")
	if strings.ContainsAny(got, "\r\n") {
		t.Errorf("stripCRLF left newline: %q", got)
	}
	if got != "helloBcc: evil@x.com" {
		t.Errorf("unexpected: %q", got)
	}
}

func TestContactRejectsInvalid(t *testing.T) {
	h := &Handler{}

	form := url.Values{
		"name":    {"Tester"},
		"email":   {"not-an-email"},
		"message": {"hi"},
	}
	req := httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.Contact(rec, req)

	body := rec.Body.String()
	if !strings.Contains(body, "c-warn") {
		t.Errorf("expected error styling for invalid email, got: %s", body)
	}
}

func TestContactHoneypot(t *testing.T) {
	h := &Handler{}

	form := url.Values{
		"name":    {"Bot"},
		"email":   {"bot@spam.com"},
		"message": {"spam"},
		"website": {"filled-by-bot"},
	}
	req := httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.Contact(rec, req)

	if !strings.Contains(rec.Body.String(), "c-ok") {
		t.Errorf("honeypot should return success-looking response")
	}
}

func TestContactAcceptsValid(t *testing.T) {
	h := &Handler{}

	form := url.Values{
		"name":    {"Tester"},
		"email":   {"tester@example.com"},
		"message": {"Hello there"},
	}
	req := httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.Contact(rec, req)

	if !strings.Contains(rec.Body.String(), "c-ok") {
		t.Errorf("valid submission should succeed, got: %s", rec.Body.String())
	}
}
