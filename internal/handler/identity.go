package handler

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"
	"time"

	"daemontalk/internal/auth"
)

func GetVisitorIdentity(w http.ResponseWriter, r *http.Request) string {
	var visitorID string

	cookie, err := r.Cookie(CookieVisitorID)
	if err != nil || cookie.Value == "" {

		visitorID = fmt.Sprintf("%d-%s", time.Now().UnixNano(), clientIP(r))

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

	return generateAnonymousName(visitorID)
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
	return fmt.Sprintf("ip:%s", clientIP(r))
}

func isBot(r *http.Request) bool {
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

func generateAnonymousName(id string) string {
	hash := sha256.Sum256([]byte(id))
	return fmt.Sprintf("anonym_%x", hash[:3])
}
