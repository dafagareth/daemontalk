package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"daemontalk/internal/auth"
	"daemontalk/internal/forum"
	"daemontalk/internal/i18n"
	"daemontalk/internal/post"
	"daemontalk/web/templates"
)

const oauthStateCookie = "daemontalk_oauth_state"

// AuthMiddleware inspects session cookies and injects the authenticated user into the request context.
func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.Auth == nil {
			next.ServeHTTP(w, r)
			return
		}

		rawToken := auth.GetSessionTokenFromRequest(r)
		if rawToken != "" {
			tokenHash := auth.HashToken(rawToken)
			user, err := h.Auth.GetSessionUser(tokenHash)
			if err != nil {
				slog.Warn("failed to query session user", "error", err)
			} else if user != nil {
				if user.Role == "member" && post.IsContributor(h.AllPosts(), user.Username) {
					user.Role = "contributor"
				}
				r = r.WithContext(auth.WithUser(r.Context(), user))
			}
		}

		next.ServeHTTP(w, r)
	})
}

// AuthGitHub initiates the GitHub OAuth redirect.
func (h *Handler) AuthGitHub(w http.ResponseWriter, r *http.Request) {
	if h.GitHubOAuth == nil {
		http.Error(w, "GitHub OAuth is not configured on this server", http.StatusServiceUnavailable)
		return
	}

	state, err := auth.GenerateRandomToken()
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

	// Store state + returnTo in cookie for verification
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

// AuthGitHubCallback processes the redirect back from GitHub.
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

	// Verify state cookie
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

	// Clear state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     oauthStateCookie,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Exchange authorization code for User profile
	ghUser, err := h.GitHubOAuth.ExchangeToken(r.Context(), code)
	if err != nil {
		slog.Error("failed to exchange github oauth token", "error", err)
		http.Error(w, "Failed to authenticate with GitHub", http.StatusInternalServerError)
		return
	}

	// Upsert user in database
	dbUser, err := h.Auth.UpsertUser(*ghUser)
	if err != nil {
		slog.Error("failed to upsert user in auth db", "error", err)
		http.Error(w, "Failed to save user session", http.StatusInternalServerError)
		return
	}

	// Create Session Token
	rawToken, err := auth.GenerateRandomToken()
	if err != nil {
		http.Error(w, "Failed to generate session", http.StatusInternalServerError)
		return
	}

	tokenHash := auth.HashToken(rawToken)
	if _, err := h.Auth.CreateSession(dbUser.ID, tokenHash, auth.SessionDuration); err != nil {
		slog.Error("failed to create session in auth db", "error", err)
		http.Error(w, "Failed to persist session", http.StatusInternalServerError)
		return
	}

	// Set HttpOnly cookie
	auth.SetSessionCookie(w, rawToken, h.IsProduction)
	http.Redirect(w, r, returnTo, http.StatusSeeOther)
}

// AuthLogout destroys the active session and clears the cookie.
func (h *Handler) AuthLogout(w http.ResponseWriter, r *http.Request) {
	if h.Auth != nil {
		rawToken := auth.GetSessionTokenFromRequest(r)
		if rawToken != "" {
			tokenHash := auth.HashToken(rawToken)
			_ = h.Auth.DeleteSession(tokenHash)
		}
	}

	auth.ClearSessionCookie(w, h.IsProduction)

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

// AuthMe returns the current authenticated user profile as JSON.
func (h *Handler) AuthMe(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
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

// AuthBadge renders the HTMX user profile or login badge in the top navigation bar.
func (h *Handler) AuthBadge(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	returnTo := sanitizeReturnTo(r.Header.Get("HX-Current-URL"))
	if returnTo == "/socket" {
		if ref := r.Header.Get("Referer"); ref != "" {
			returnTo = sanitizeReturnTo(ref)
		}
	}
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = langFromRequest(r)
	}
	currentPath := r.URL.Query().Get("path")
	if currentPath == "" {
		currentPath = "/"
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	mode := r.URL.Query().Get("mode")
	if mode == "mobile" {
		_ = templates.AuthNavBadgeMobile(user, returnTo, lang, currentPath).Render(r.Context(), w)
		return
	}
	_ = templates.AuthNavBadge(user, returnTo, lang, currentPath).Render(r.Context(), w)
}

// AuthExport downloads all user data (profile, topics, replies, comments) as a JSON file.
func (h *Handler) AuthExport(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		http.Redirect(w, r, "/auth/github", http.StatusSeeOther)
		return
	}

	data := map[string]any{
		"exported_at": time.Now().UTC().Format(time.RFC3339),
		"platform":    "daemontalk",
		"user":        user,
	}

	if h.Forum != nil {
		if contrib, err := h.Forum.GetUserContributions(user.ID); err == nil {
			data["forum_topics"] = contrib.Topics
			data["forum_replies"] = contrib.Replies
		}
	}

	if h.Comments != nil {
		if userComments, err := h.Comments.ListByUserID(user.ID); err == nil {
			data["article_comments"] = userComments
		}
	}

	filename := fmt.Sprintf("daemontalk-data-%s.json", user.Username)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	_ = json.NewEncoder(w).Encode(data)
}

// AuthDeleteAccount permanently purges the user's account and anonymizes their public forum contributions.
func (h *Handler) AuthDeleteAccount(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// 1. Anonymize forum topics and replies
	if h.Forum != nil {
		_ = h.Forum.AnonymizeUser(user.ID)
	}

	// 2. Anonymize article comments
	if h.Comments != nil {
		_ = h.Comments.AnonymizeUserComments(user.ID)
	}

	// 3. Delete user & cascade delete sessions
	if h.Auth != nil {
		_ = h.Auth.DeleteUser(user.ID)
	}

	// 4. Clear session cookie
	auth.ClearSessionCookie(w, h.IsProduction)

	lang := langFromRequest(r)
	redirectURL := "/"
	if lang != "en" && lang != "" {
		redirectURL = "/" + lang
	}
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func splitStateCookie(val string) [2]string {
	for i := 0; i < len(val); i++ {
		if val[i] == '|' {
			return [2]string{val[:i], sanitizeReturnTo(val[i+1:])}
		}
	}
	return [2]string{val, "/socket"}
}

// sanitizeReturnTo ensures all redirect destinations stay on relative paths and do not redirect to localhost or external sites.
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

// AuthSettings renders the user settings page.
func (h *Handler) AuthSettings(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	lang := langFromRequest(r)
	currentPath := r.URL.Path
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = templates.SettingsPage(i18n.Get(lang), user, lang, currentPath).Render(r.Context(), w)
}

// AuthUserProfile renders the public user profile page.
func (h *Handler) AuthUserProfile(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var profileUser *auth.User
	if h.Auth != nil {
		user, err := h.Auth.GetUserByUsername(username)
		if err == nil {
			profileUser = user
			if profileUser != nil && profileUser.Role == "member" && post.IsContributor(h.AllPosts(), profileUser.Username) {
				profileUser.Role = "contributor"
			}
		}
	}

	var stats forum.UserStats
	var recentTopics []*forum.Topic
	if h.Forum != nil && profileUser != nil {
		if s, err := h.Forum.GetUserStats(profileUser.Username); err == nil {
			stats = s
		}
		var currentUserID int64
		if u := auth.GetUser(r.Context()); u != nil {
			currentUserID = u.ID
		}
		if ts, _, err := h.Forum.ListTopics("", "", "", profileUser.Username, "latest", 5, 0, currentUserID); err == nil {
			recentTopics = ts
		}
	}

	lang := langFromRequest(r)
	currentPath := r.URL.Path
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = templates.UserProfilePage(i18n.Get(lang), profileUser, stats, recentTopics, lang, currentPath).Render(r.Context(), w)
}
