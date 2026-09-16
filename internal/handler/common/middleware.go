package common

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

const contentSecurityPolicy = "default-src 'self'; " +
	"script-src 'self' 'unsafe-inline' https://unpkg.com https://static.cloudflareinsights.com; " +
	"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; " +
	"font-src 'self' https://fonts.gstatic.com; " +
	"img-src 'self' data: https:; " +
	"connect-src 'self' https://cloudflareinsights.com; " +
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

func CountablePath(p string) bool {
	skipPrefixes := []string{
		"/static/", "/admin", "/api/",
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
		"/sitemap.xml", "/sitemap-index.xml", "/sitemap-en.xml",
		"/robots.txt", "/healthz", "/favicon.ico", "/manifest.json", "/blog/posts",
	}
	for _, e := range skipExact {
		if p == e {
			return false
		}
	}
	return true
}

type RateLimiter struct {
	mu     sync.Mutex
	hits   map[string][]time.Time
	limit  int
	window time.Duration
	lastGC time.Time
}

func NewRateLimiter(limit int, window time.Duration) func(http.Handler) http.Handler {
	rl := &RateLimiter{
		hits:   make(map[string][]time.Time),
		limit:  limit,
		window: window,
		lastGC: time.Now(),
	}
	return rl.Middleware
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.Allow(ClientIP(r)) {
			w.Header().Set("Retry-After", "60")
			http.Error(w, "Too many requests. Please slow down.", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) Allow(ip string) bool {
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
