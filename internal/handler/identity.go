package handler

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"
)

// GetVisitorIdentity checks for a visitor tracking cookie. If it doesn't exist,
// it generates one and sets it. It returns a consistent, pseudo-anonymous handle
// (e.g., "anonym_a3f89c") deterministically hashed from that cookie.
func GetVisitorIdentity(w http.ResponseWriter, r *http.Request) string {
	cookieName := "visitor_id"
	var visitorID string

	cookie, err := r.Cookie(cookieName)
	if err != nil || cookie.Value == "" {
		// Generate a new random visitor ID based on time and client IP
		visitorID = fmt.Sprintf("%d-%s", time.Now().UnixNano(), clientIP(r))

		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Value:    visitorID,
			Path:     "/",
			Expires:  time.Now().AddDate(10, 0, 0), // 10 years
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
	} else {
		visitorID = cookie.Value
	}

	return generateAnonymousName(visitorID)
}

// generateAnonymousName deterministically hashes the visitor ID into an anonym_<hex> handle.
func generateAnonymousName(id string) string {
	hash := sha256.Sum256([]byte(id))
	return fmt.Sprintf("anonym_%x", hash[:3])
}
