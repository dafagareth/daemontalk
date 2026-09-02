package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"
)

const (
	SessionCookieName = "daemontalk_session"
	SessionDuration   = 30 * 24 * time.Hour // 30 days
)

// GenerateRandomToken generates a cryptographically secure random string.
func GenerateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashToken hashes a raw token with SHA-256 for secure database storage.
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// SetSessionCookie sets a secure, HttpOnly session cookie on the response.
func SetSessionCookie(w http.ResponseWriter, token string, isProduction bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(SessionDuration),
		MaxAge:   int(SessionDuration.Seconds()),
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearSessionCookie removes the session cookie from the client.
func ClearSessionCookie(w http.ResponseWriter, isProduction bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteLaxMode,
	})
}

// GetSessionTokenFromRequest extracts the session cookie value from an incoming HTTP request.
func GetSessionTokenFromRequest(r *http.Request) string {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil || cookie.Value == "" {
		return ""
	}
	return cookie.Value
}
