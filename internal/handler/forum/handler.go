package forum

import (
	"net/http"

	"daemontalk/internal/forum"
	"daemontalk/internal/handler/common"
)

type Handler struct {
	Forum      *forum.Store
	AdminToken string
}

func (h *Handler) IsAdmin(r *http.Request) bool {
	return common.IsAdmin(h.AdminToken, r)
}
