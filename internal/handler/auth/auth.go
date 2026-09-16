package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	authstore "daemontalk/internal/auth"
	"daemontalk/internal/handler/common"
	"daemontalk/internal/post"
	"daemontalk/web/templates/layouts"
)

const oauthStateCookie = "daemontalk_oauth_state"

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.Auth == nil {
			next.ServeHTTP(w, r)
			return
		}

		rawToken := authstore.GetSessionTokenFromRequest(r)
		if rawToken != "" {
			tokenHash := authstore.HashToken(rawToken)
			user, err := h.Auth.GetSessionUser(tokenHash)
			if err != nil {
				slog.Warn("failed to query session user", "error", err)
			} else if user != nil {
				if user.Role == "member" && post.IsContributor(h.GetAllPosts(), user.Username) {
					user.Role = "contributor"
				}
				r = r.WithContext(authstore.WithUser(r.Context(), user))
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) AuthGitHub(w http.ResponseWriter, r *http.Request) {
	if h.GitHubOAuth == nil {
		http.Error(w, "GitHub OAuth is not configured on this server", http.StatusServiceUnavailable)
		return
	}

	state, err := authstore.GenerateRandomToken()
	if err != nil {
		http.Error(w, "Failed to generate security state", http.StatusInternalServerError)
		return
	}

	returnTo := sanitizeReturnTo(r.URL.Query().Get("return_to"))
	if returnTo == "/socket" {
		if ref := r.Header.Get("Referer"); ref != "" {
			returnTo = sanitizeReturnTo(ref)
		}
	}

	statePayload := state + "|" + returnTo
	http.SetCookie(w, &http.Cookie{
		Name:     oauthStateCookie,
		Value:    statePayload,
		Path:     "/",
		Expires:  time.Now().Add(10 * time.Minute),
		MaxAge:   600,
		HttpOnly: true,
		Secure:   h.IsProduction,
		SameSite: http.SameSiteLaxMode,
	})

	url := h.GitHubOAuth.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) AuthGitHubCallback(w http.ResponseWriter, r *http.Request) {
	if h.GitHubOAuth == nil || h.Auth == nil {
		http.Error(w, "OAuth is not configured", http.StatusServiceUnavailable)
		return
	}

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if code == "" || state == "" {
		http.Error(w, "Invalid OAuth callback request", http.StatusBadRequest)
		return
	}

	stateCookie, err := r.Cookie(oauthStateCookie)
	if err != nil || stateCookie.Value == "" {
		http.Error(w, "Security state cookie missing or expired. Please try again.", http.StatusBadRequest)
		return
	}

	parts := splitStateCookie(stateCookie.Value)
	if parts[0] != state {
		http.Error(w, "Invalid security state mismatch", http.StatusBadRequest)
		return
	}
	returnTo := sanitizeReturnTo(parts[1])

	http.SetCookie(w, &http.Cookie{
		Name:     oauthStateCookie,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	})

	ghUser, err := h.GitHubOAuth.ExchangeToken(r.Context(), code)
	if err != nil {
		slog.Error("failed to exchange github oauth token", "error", err)
		http.Error(w, "Failed to authenticate with GitHub", http.StatusInternalServerError)
		return
	}

	dbUser, err := h.Auth.UpsertUser(*ghUser)
	if err != nil {
		slog.Error("failed to upsert user in auth db", "error", err)
		http.Error(w, "Failed to save user session", http.StatusInternalServerError)
		return
	}

	rawToken, err := authstore.GenerateRandomToken()
	if err != nil {
		http.Error(w, "Failed to generate session", http.StatusInternalServerError)
		return
	}

	tokenHash := authstore.HashToken(rawToken)
	if _, err := h.Auth.CreateSession(dbUser.ID, tokenHash, authstore.SessionDuration); err != nil {
		slog.Error("failed to create session in auth db", "error", err)
		http.Error(w, "Failed to persist session", http.StatusInternalServerError)
		return
	}

	authstore.SetSessionCookie(w, rawToken, h.IsProduction)
	http.Redirect(w, r, returnTo, http.StatusSeeOther)
}

func (h *Handler) AuthLogout(w http.ResponseWriter, r *http.Request) {
	if h.Auth != nil {
		rawToken := authstore.GetSessionTokenFromRequest(r)
		if rawToken != "" {
			tokenHash := authstore.HashToken(rawToken)
			_ = h.Auth.DeleteSession(tokenHash)
		}
	}

	authstore.ClearSessionCookie(w, h.IsProduction)

	returnTo := sanitizeReturnTo(r.URL.Query().Get("return_to"))
	if returnTo == "/socket" {
		if ref := r.Header.Get("Referer"); ref != "" {
			returnTo = sanitizeReturnTo(ref)
		}
	}
	if returnTo == "/socket" {
		returnTo = "/"
	}
	http.Redirect(w, r, returnTo, http.StatusSeeOther)
}

func (h *Handler) AuthMe(w http.ResponseWriter, r *http.Request) {
	user := authstore.GetUser(r.Context())
	w.Header().Set("Content-Type", "application/json")
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]any{"authenticated": false})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{
		"authenticated": true,
		"user":          user,
	})
}

func (h *Handler) AuthBadge(w http.ResponseWriter, r *http.Request) {
	user := authstore.GetUser(r.Context())
	returnTo := sanitizeReturnTo(r.Header.Get("HX-Current-URL"))
	if returnTo == "/socket" {
		if ref := r.Header.Get("Referer"); ref != "" {
			returnTo = sanitizeReturnTo(ref)
		}
	}
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = common.LangFromRequest(r)
	}
	currentPath := r.URL.Query().Get("path")
	if currentPath == "" {
		currentPath = "/"
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	mode := r.URL.Query().Get("mode")
	if mode == "mobile" {
		_ = layouts.AuthNavBadgeMobile(user, returnTo, lang, currentPath).Render(r.Context(), w)
		return
	}
	_ = layouts.AuthNavBadge(user, returnTo, lang, currentPath).Render(r.Context(), w)
}

func splitStateCookie(val string) [2]string {
	for i := 0; i < len(val); i++ {
		if val[i] == '|' {
			return [2]string{val[:i], sanitizeReturnTo(val[i+1:])}
		}
	}
	return [2]string{val, "/socket"}
}

func sanitizeReturnTo(returnTo string) string {
	returnTo = strings.TrimSpace(returnTo)
	if returnTo == "" || returnTo == "/auth/github" {
		return "/socket"
	}
	if strings.HasPrefix(returnTo, "http://") || strings.HasPrefix(returnTo, "https://") {
		if u, err := url.Parse(returnTo); err == nil {
			path := u.EscapedPath()
			if path == "" {
				path = "/"
			}
			if u.RawQuery != "" {
				path += "?" + u.RawQuery
			}
			if path == "/auth/github" {
				return "/socket"
			}
			return path
		}
		return "/socket"
	}
	if !strings.HasPrefix(returnTo, "/") || strings.HasPrefix(returnTo, "//") {
		return "/socket"
	}
	return returnTo
}
