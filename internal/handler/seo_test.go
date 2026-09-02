package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSitemapAfterPivot(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
	rec := httptest.NewRecorder()
	h.Sitemap(rec, req)

	body := rec.Body.String()
	if !strings.Contains(body, seoBaseURL+"/colophon") {
		t.Error("sitemap harus memuat /colophon")
	}
	for _, dead := range []string{"/about", "/uses", "/now", "/projects</loc>", "/blog</loc>"} {
		if strings.Contains(body, dead) {
			t.Errorf("sitemap masih memuat route mati %q", dead)
		}
	}
}
