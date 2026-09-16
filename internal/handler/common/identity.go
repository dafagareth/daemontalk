package common

import (
	"crypto/sha256"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"daemontalk/internal/auth"
)

func ClientIP(r *http.Request) string {
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

func GenerateAnonymousName(id string) string {
	hash := sha256.Sum256([]byte(id))
	return fmt.Sprintf("anonym_%x", hash[:3])
}

func GetVisitorIdentity(w http.ResponseWriter, r *http.Request) string {
	var visitorID string

	cookie, err := r.Cookie(CookieVisitorID)
	if err != nil || cookie.Value == "" {
		visitorID = fmt.Sprintf("%d-%s", time.Now().UnixNano(), ClientIP(r))

		http.SetCookie(w, &http.Cookie{
			Name:     CookieVisitorID,
			Value:    visitorID,
			Path:     "/",
			Expires:  time.Now().AddDate(CookieVisitorExpiryYears, 0, 0),
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
	} else {
		visitorID = cookie.Value
	}

	return GenerateAnonymousName(visitorID)
}

func GetViewerKey(w http.ResponseWriter, r *http.Request, user *auth.User) string {
	if user != nil && user.ID > 0 {
		return fmt.Sprintf("u:%d", user.ID)
	}
	cookie, err := r.Cookie(CookieVisitorID)
	if err == nil && cookie.Value != "" {
		return fmt.Sprintf("v:%s", cookie.Value)
	}
	_ = GetVisitorIdentity(w, r)
	if c, err := r.Cookie(CookieVisitorID); err == nil && c.Value != "" {
		return fmt.Sprintf("v:%s", c.Value)
	}
	return fmt.Sprintf("ip:%s", ClientIP(r))
}

func IsBot(r *http.Request) bool {
	ua := strings.ToLower(r.UserAgent())
	if ua == "" {
		return true
	}
	botKeywords := []string{
		"bot", "crawler", "spider", "googlebot", "bingbot",
		"yandex", "baidu", "duckduck", "slurp", "headless",
		"curl", "wget", "python", "httpclient", "ahrefs", "semrush",
	}
	for _, kw := range botKeywords {
		if strings.Contains(ua, kw) {
			return true
		}
	}
	return false
}
