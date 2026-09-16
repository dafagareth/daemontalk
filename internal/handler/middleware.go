package handler

import (
	"net/http"
	"time"

	"daemontalk/internal/handler/common"
)

func SecurityHeaders(next http.Handler) http.Handler {
	return common.SecurityHeaders(next)
}

func StaticCacheControl(next http.Handler) http.Handler {
	return common.StaticCacheControl(next)
}

func NewRateLimiter(limit int, window time.Duration) func(http.Handler) http.Handler {
	return common.NewRateLimiter(limit, window)
}

func countablePath(p string) bool {
	return common.CountablePath(p)
}

func (h *Handler) Analytics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.Comments != nil && r.Method == http.MethodGet && common.CountablePath(r.URL.Path) {
			path := r.URL.Path
			go func() { _ = h.Comments.IncrementPageView(path) }()
		}
		next.ServeHTTP(w, r)
	})
}
