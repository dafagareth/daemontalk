package handler

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// contentSecurityPolicy whitelists exactly the external origins the site uses:
// htmx + fonts from their CDNs, inline styles/scripts (theme + page scripts),
// and self for everything else. Adjust if a new CDN is introduced.
const contentSecurityPolicy = "default-src 'self'; " +
	"script-src 'self' 'unsafe-inline' https://unpkg.com; " +
	"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; " +
	"font-src 'self' https://fonts.gstatic.com; " +
	"img-src 'self' data: https:; " +
	"connect-src 'self'; " +
	"frame-ancestors 'none'; " +
	"base-uri 'self'; " +
	"form-action 'self'"

// SecurityHeaders sets a baseline of hardening headers on every response.
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

// Analytics records a page hit for GET requests to HTML pages. Static assets,
// generated images, feeds and the admin area are skipped. Recording happens in
// a goroutine so it never adds latency to the response.
func (h *Handler) Analytics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.Comments != nil && r.Method == http.MethodGet && countablePath(r.URL.Path) {
			path := r.URL.Path
			go func() { _ = h.Comments.IncrementPageView(path) }()
		}
		next.ServeHTTP(w, r)
	})
}

// countablePath reports whether a path represents a real page view worth
// tracking (as opposed to an asset, image, feed, or admin route).
func countablePath(p string) bool {
	skipPrefixes := []string{"/static/", "/id/static/", "/admin"}
	for _, pre := range skipPrefixes {
		if strings.HasPrefix(p, pre) {
			return false
		}
	}
	skipExact := []string{"/og.png", "/rss.xml", "/feed.xml", "/feed.json", "/sitemap.xml", "/robots.txt", "/healthz", "/favicon.ico", "/manifest.json", "/blog/posts"}
	for _, e := range skipExact {
		if p == e {
			return false
		}
	}
	// Skip generated per-post OG images (/blog/<slug>/og.png).
	if strings.HasSuffix(p, "/og.png") {
		return false
	}
	return true
}

// rateLimiter is a simple per-IP sliding-window limiter kept in memory.
type rateLimiter struct {
	mu     sync.Mutex
	hits   map[string][]time.Time
	limit  int
	window time.Duration
	lastGC time.Time
}

// NewRateLimiter allows `limit` requests per `window` per client IP.
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

	// Occasional garbage collection of stale IP buckets.
	if now.Sub(rl.lastGC) > rl.window {
		for k, ts := range rl.hits {
			if len(ts) == 0 || ts[len(ts)-1].Before(cutoff) {
				delete(rl.hits, k)
			}
		}
		rl.lastGC = now
	}

	recent := rl.hits[ip][:0]
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
	// 1. Cloudflare native header
	if cfip := r.Header.Get("CF-Connecting-IP"); cfip != "" {
		return trimSpace(cfip)
	}
	// 2. Standard proxies (NGINX, HAProxy, etc)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if i := indexByte(xff, ','); i >= 0 {
			return trimSpace(xff[:i])
		}
		return trimSpace(xff)
	}
	if rip := r.Header.Get("X-Real-IP"); rip != "" {
		return trimSpace(rip)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func indexByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}

func trimSpace(s string) string {
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
