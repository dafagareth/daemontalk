package handler

import (
	"net/http"

	"daemontalk/internal/handler/forum"
)

func (h *Handler) ForumHandler() *forum.Handler {
	return &forum.Handler{
		Forum:      h.Forum,
		AdminToken: h.AdminToken,
	}
}

func (h *Handler) Discussions(w http.ResponseWriter, r *http.Request) {
	h.ForumHandler().Discussions(w, r)
}

func (h *Handler) DiscussionsNew(w http.ResponseWriter, r *http.Request) {
	h.ForumHandler().DiscussionsNew(w, r)
}

func (h *Handler) DiscussionsCreate(w http.ResponseWriter, r *http.Request) {
	h.ForumHandler().DiscussionsCreate(w, r)
}

func (h *Handler) DiscussionsDetail(w http.ResponseWriter, r *http.Request) {
	h.ForumHandler().DiscussionsDetail(w, r)
}

func (h *Handler) DiscussionsDeleteTopic(w http.ResponseWriter, r *http.Request) {
	h.ForumHandler().DiscussionsDeleteTopic(w, r)
}

func (h *Handler) DiscussionsReply(w http.ResponseWriter, r *http.Request) {
	h.ForumHandler().DiscussionsReply(w, r)
}

func (h *Handler) DiscussionsSolve(w http.ResponseWriter, r *http.Request) {
	h.ForumHandler().DiscussionsSolve(w, r)
}

func (h *Handler) DiscussionsVote(w http.ResponseWriter, r *http.Request) {
	h.ForumHandler().DiscussionsVote(w, r)
}

func (h *Handler) DiscussionsDeleteReply(w http.ResponseWriter, r *http.Request) {
	h.ForumHandler().DiscussionsDeleteReply(w, r)
}
