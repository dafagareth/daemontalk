package common

import (
	"net/http"
)

func IsAdmin(adminToken string, r *http.Request) bool {
	if adminToken == "" {
		return false
	}
	c, err := r.Cookie(CookieAdminToken)
	if err != nil {
		return false
	}
	return c.Value == adminToken
}

func SetAdminCookie(w http.ResponseWriter, r *http.Request, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieAdminToken,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   CookieAdminMaxAge,
	})
}
