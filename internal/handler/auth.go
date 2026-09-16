package handler

import (
	"net/http"

	authhandler "daemontalk/internal/handler/auth"
)

func (h *Handler) AuthHandler() *authhandler.Handler {
	return &authhandler.Handler{
		Auth:         h.Auth,
		GitHubOAuth:  h.GitHubOAuth,
		Comments:     h.Comments,
		Forum:        h.Forum,
		AllPosts:     h.AllPosts,
		IsProduction: h.IsProduction,
	}
}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return h.AuthHandler().AuthMiddleware(next)
}

func (h *Handler) AuthGitHub(w http.ResponseWriter, r *http.Request) {
	h.AuthHandler().AuthGitHub(w, r)
}

func (h *Handler) AuthGitHubCallback(w http.ResponseWriter, r *http.Request) {
	h.AuthHandler().AuthGitHubCallback(w, r)
}

func (h *Handler) AuthLogout(w http.ResponseWriter, r *http.Request) {
	h.AuthHandler().AuthLogout(w, r)
}

func (h *Handler) AuthMe(w http.ResponseWriter, r *http.Request) {
	h.AuthHandler().AuthMe(w, r)
}

func (h *Handler) AuthBadge(w http.ResponseWriter, r *http.Request) {
	h.AuthHandler().AuthBadge(w, r)
}

func (h *Handler) AuthSettings(w http.ResponseWriter, r *http.Request) {
	h.AuthHandler().AuthSettings(w, r)
}

func (h *Handler) AuthMyProfile(w http.ResponseWriter, r *http.Request) {
	h.AuthHandler().AuthMyProfile(w, r)
}

func (h *Handler) AuthUserProfile(w http.ResponseWriter, r *http.Request) {
	h.AuthHandler().AuthUserProfile(w, r)
}

func (h *Handler) AuthExport(w http.ResponseWriter, r *http.Request) {
	h.AuthHandler().AuthExport(w, r)
}

func (h *Handler) AuthDeleteAccount(w http.ResponseWriter, r *http.Request) {
	h.AuthHandler().AuthDeleteAccount(w, r)
}
