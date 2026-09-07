package handler

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

const contentSecurityPolicy = "default-src 'self'; " +
	"script-src 'self' 'unsafe-inline' https://unpkg.com https://cdn.jsdelivr.net https://static.cloudflareinsights.com; " +
	"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com https://cdn.jsdelivr.net; " +
	"font-src 'self' https://fonts.gstatic.com https://cdn.jsdelivr.net; " +
	"img-src 'self' data: https:; " +
	"connect-src 'self' https://*.wikipedia.org https://cloudflareinsights.com; " +
	"frame-ancestors 'none'; " +
	"base-uri 'self'; " +
	"form-action 'self'"

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		h.Set("Content-Security-Policy", contentSecurityPolicy)
		next.ServeHTTP(w, r)
	})
}

func StaticCacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if v := r.URL.Query().Get("v"); v != "" && v != "dev" {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "no-cache, must-revalidate")
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) Analytics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.Comments != nil && r.Method == http.MethodGet && countablePath(r.URL.Path) {
			path := r.URL.Path
			go func() { _ = h.Comments.IncrementPageView(path) }()
		}
		next.ServeHTTP(w, r)
	})
}

func countablePath(p string) bool {
	skipPrefixes := []string{
		"/static/", "/id/static/", "/admin", "/api/",
		"/wp-", "/wordpress", "/php", "/.env", "/.git",
	}
	for _, pre := range skipPrefixes {
		if strings.HasPrefix(p, pre) {
			return false
		}
	}
	skipSuffixes := []string{
		"/og.png", "/comments/stream", "/comments",
		".php", ".env", ".git", ".bak", ".sql", ".asp", ".aspx", ".jsp",
	}
	for _, suf := range skipSuffixes {
		if strings.HasSuffix(p, suf) {
			return false
		}
	}
	skipExact := []string{
		"/og.png", "/rss.xml", "/feed.xml", "/feed.json",
		"/sitemap.xml", "/sitemap-index.xml", "/sitemap-en.xml", "/sitemap-id.xml",
		"/robots.txt", "/healthz", "/favicon.ico", "/manifest.json", "/blog/posts",
	}
	for _, e := range skipExact {
		if p == e {
			return false
		}
	}
	return true
}

type rateLimiter struct {
	mu     sync.Mutex
	hits   map[string][]time.Time
	limit  int
	window time.Duration
	lastGC time.Time
}

func NewRateLimiter(limit int, window time.Duration) func(http.Handler) http.Handler {
	rl := &rateLimiter{
		hits:   make(map[string][]time.Time),
		limit:  limit,
		window: window,
		lastGC: time.Now(),
	}
	return rl.middleware
}

func (rl *rateLimiter) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.allow(clientIP(r)) {
			w.Header().Set("Retry-After", "60")
			http.Error(w, "Too many requests. Please slow down.", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rl *rateLimiter) allow(ip string) bool {
	now := time.Now()
	cutoff := now.Add(-rl.window)

	rl.mu.Lock()
	defer rl.mu.Unlock()

	if now.Sub(rl.lastGC) > rl.window {
		for k, ts := range rl.hits {
			if len(ts) == 0 || ts[len(ts)-1].Before(cutoff) {
				delete(rl.hits, k)
			}
		}
		rl.lastGC = now
	}

	var recent []time.Time
	for _, t := range rl.hits[ip] {
		if t.After(cutoff) {
			recent = append(recent, t)
		}
	}
	if len(recent) >= rl.limit {
		rl.hits[ip] = recent
		return false
	}
	rl.hits[ip] = append(recent, now)
	return true
}

func clientIP(r *http.Request) string {

	if cfip := r.Header.Get("CF-Connecting-IP"); cfip != "" {
		return strings.TrimSpace(cfip)
	}

	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if i := strings.IndexByte(xff, ','); i >= 0 {
			return strings.TrimSpace(xff[:i])
		}
		return strings.TrimSpace(xff)
	}
	if rip := r.Header.Get("X-Real-IP"); rip != "" {
		return strings.TrimSpace(rip)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
