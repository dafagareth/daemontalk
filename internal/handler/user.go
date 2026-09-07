package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"daemontalk/internal/auth"
	"daemontalk/internal/forum"
	"daemontalk/internal/i18n"
	"daemontalk/internal/post"
	"daemontalk/web/templates"
)

func (h *Handler) AuthSettings(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	lang := langFromRequest(r)
	currentPath := r.URL.Path
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = templates.SettingsPage(i18n.Get(lang), user, lang, currentPath).Render(r.Context(), w)
}

func (h *Handler) AuthMyProfile(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user != nil {
		http.Redirect(w, r, "/u/"+user.Username, http.StatusSeeOther)
		return
	}
	returnTo := r.URL.Query().Get("return_to")
	if returnTo == "" {
		returnTo = "/"
	}
	http.Redirect(w, r, "/auth/github?return_to="+returnTo, http.StatusSeeOther)
}

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

func (h *Handler) AuthDeleteAccount(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if h.Forum != nil {
		_ = h.Forum.AnonymizeUser(user.ID)
	}

	if h.Comments != nil {
		_ = h.Comments.AnonymizeUserComments(user.ID)
	}

	if h.Auth != nil {
		_ = h.Auth.DeleteUser(user.ID)
	}

	auth.ClearSessionCookie(w, h.IsProduction)

	lang := langFromRequest(r)
	redirectURL := "/"
	if lang != "en" && lang != "" {
		redirectURL = "/" + lang
	}
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
