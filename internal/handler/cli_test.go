package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"daemontalk/internal/post"
)

func TestIsCLIRequest(t *testing.T) {
	reqCurl, _ := http.NewRequest("GET", "/", nil)
	reqCurl.Header.Set("User-Agent", "curl/8.5.0")
	if !IsCLIRequest(reqCurl) {
		t.Errorf("expected curl to be detected as CLI request")
	}

	reqWget, _ := http.NewRequest("GET", "/", nil)
	reqWget.Header.Set("User-Agent", "Wget/1.21.4")
	if !IsCLIRequest(reqWget) {
		t.Errorf("expected wget to be detected as CLI request")
	}

	reqBrowser, _ := http.NewRequest("GET", "/", nil)
	reqBrowser.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36")
	if IsCLIRequest(reqBrowser) {
		t.Errorf("expected browser to not be detected as CLI request")
	}

	reqPlain, _ := http.NewRequest("GET", "/?plain=1", nil)
	if !IsCLIRequest(reqPlain) {
		t.Errorf("expected ?plain=1 to be detected as CLI request")
	}
}

func TestCLIMainAndDaily(t *testing.T) {
	h := &Handler{
		FilePosts: []post.Post{
			{
				Title:   "eBPF Linux Observability",
				Slug:    "ebpf-linux-observability",
				Date:    time.Now(),
				Lang:    "en",
				Tags:        []string{"linux", "ebpf"},
				Description: "Deep dive into eBPF tracing",
				Body:        "## eBPF Internals\n\nDetailed content here.",
			},
		},
	}
	h.RefreshPosts()

	// Test CLIMain
	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = "localhost:8080"
	rec := httptest.NewRecorder()
	h.CLIMain(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "daemontalk") || !strings.Contains(body, "eBPF Linux Observability") {
		t.Errorf("unexpected CLIMain response: %s", body)
	}

	// Test CLIDaily (CLI mode)
	reqDailyCLI, _ := http.NewRequest("GET", "/daily", nil)
	reqDailyCLI.Host = "localhost:8080"
	reqDailyCLI.Header.Set("User-Agent", "curl/8.5.0")
	recDailyCLI := httptest.NewRecorder()
	h.CLIDaily(recDailyCLI, reqDailyCLI)

	if recDailyCLI.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recDailyCLI.Code)
	}
	if !strings.Contains(recDailyCLI.Body.String(), "DAILY TECH BRIEFING") {
		t.Errorf("unexpected CLIDaily CLI response: %s", recDailyCLI.Body.String())
	}

	// Test CLIDaily (Browser mode)
	reqDailyBrowser, _ := http.NewRequest("GET", "/daily", nil)
	reqDailyBrowser.Host = "localhost:8080"
	reqDailyBrowser.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64)")
	reqDailyBrowser.Header.Set("Accept", "text/html")
	recDailyBrowser := httptest.NewRecorder()
	h.CLIDaily(recDailyBrowser, reqDailyBrowser)

	if recDailyBrowser.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recDailyBrowser.Code)
	}
	if !strings.Contains(recDailyBrowser.Body.String(), "Daily Technical Briefing") {
		t.Errorf("unexpected CLIDaily Browser response: %s", recDailyBrowser.Body.String())
	}

	// Test CLIRecipes
	reqRecipes, _ := http.NewRequest("GET", "/recipes", nil)
	recRecipes := httptest.NewRecorder()
	h.CLIRecipes(recRecipes, reqRecipes)

	if recRecipes.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recRecipes.Code)
	}
	if !strings.Contains(recRecipes.Body.String(), "eBPF RECIPES") {
		t.Errorf("unexpected CLIRecipes response: %s", recRecipes.Body.String())
	}

	// Test CLIPost
	r := chi.NewRouter()
	r.Get("/p/{slug}", h.CLIPost)

	reqPost, _ := http.NewRequest("GET", "/p/ebpf-linux-observability", nil)
	recPost := httptest.NewRecorder()
	r.ServeHTTP(recPost, reqPost)

	if recPost.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recPost.Code)
	}
	if !strings.Contains(recPost.Body.String(), "eBPF Internals") {
		t.Errorf("unexpected CLIPost response: %s", recPost.Body.String())
	}
}
