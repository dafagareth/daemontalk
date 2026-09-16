package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	authstore "daemontalk/internal/auth"
	"daemontalk/internal/forum"
	"daemontalk/internal/handler/common"
	"daemontalk/internal/i18n"
	"daemontalk/internal/post"
	"daemontalk/web/templates/layouts"
	pagestmpl "daemontalk/web/templates/pages"
	"daemontalk/web/templates/shared"
)

func (h *Handler) AuthSettings(w http.ResponseWriter, r *http.Request) {
	user := authstore.GetUser(r.Context())
	lang := common.LangFromRequest(r)
	ui := i18n.Get(lang)
	currentPath := r.URL.Path
	meta := shared.PageMeta{
		Description: "User Settings and Preferences",
		Type:        "website",
	}
	common.Render(w, r, layouts.Layout(ui, lang, "settings", currentPath, meta,
		pagestmpl.SettingsPage(ui, user, lang)))
}

func (h *Handler) AuthMyProfile(w http.ResponseWriter, r *http.Request) {
	user := authstore.GetUser(r.Context())
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

	var profileUser *authstore.User
	if h.Auth != nil {
		user, err := h.Auth.GetUserByUsername(username)
		if err == nil {
			profileUser = user
			if profileUser != nil && profileUser.Role == "member" && post.IsContributor(h.GetAllPosts(), profileUser.Username) {
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
		if u := authstore.GetUser(r.Context()); u != nil {
			currentUserID = u.ID
		}
		if ts, _, err := h.Forum.ListTopics("", "", "", profileUser.Username, "latest", 5, 0, currentUserID); err == nil {
			recentTopics = ts
		}
	}

	lang := common.LangFromRequest(r)
	ui := i18n.Get(lang)
	currentPath := r.URL.Path
	desc := "User not found"
	if profileUser != nil {
		desc = "Public profile for " + profileUser.Username
	}
	meta := shared.PageMeta{
		Description: desc,
		Type:        "website",
	}
	common.Render(w, r, layouts.Layout(ui, lang, "profile", currentPath, meta,
		pagestmpl.UserProfilePage(ui, profileUser, stats, recentTopics, lang, currentPath)))
}

func (h *Handler) AuthExport(w http.ResponseWriter, r *http.Request) {
	user := authstore.GetUser(r.Context())
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
	user := authstore.GetUser(r.Context())
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

	authstore.ClearSessionCookie(w, h.IsProduction)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
