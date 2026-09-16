package handler

import (
	"net/http"

	"daemontalk/internal/auth"
	"daemontalk/internal/handler/common"
)

func GetVisitorIdentity(w http.ResponseWriter, r *http.Request) string {
	return common.GetVisitorIdentity(w, r)
}

func GetViewerKey(w http.ResponseWriter, r *http.Request, user *auth.User) string {
	return common.GetViewerKey(w, r, user)
}

func isBot(r *http.Request) bool {
	return common.IsBot(r)
}

func clientIP(r *http.Request) string {
	return common.ClientIP(r)
}

func generateAnonymousName(id string) string {
	return common.GenerateAnonymousName(id)
}
