package handler

import (
	"net/http"

	"daemontalk/internal/handler/blog"
)

func (h *Handler) BlogHandler() *blog.Handler {
	return &blog.Handler{
		AllPosts:     h.AllPosts,
		VisiblePosts: h.VisiblePosts,
		Comments:     h.Comments,
		AdminToken:   h.AdminToken,
		SMTPHost:     h.SMTPHost,
		SMTPPort:     h.SMTPPort,
		SMTPUser:     h.SMTPUser,
		SMTPPass:     h.SMTPPass,
		SMTPTo:       h.SMTPTo,
	}
}

func (h *Handler) BlogIndex(w http.ResponseWriter, r *http.Request) {
	h.BlogHandler().BlogIndex(w, r)
}

func (h *Handler) BlogPost(w http.ResponseWriter, r *http.Request) {
	h.BlogHandler().BlogPost(w, r)
}

func (h *Handler) BlogPostsPartial(w http.ResponseWriter, r *http.Request) {
	h.BlogHandler().BlogPostsPartial(w, r)
}

func (h *Handler) TagIndex(w http.ResponseWriter, r *http.Request) {
	h.BlogHandler().TagIndex(w, r)
}

func (h *Handler) TagPostsPartial(w http.ResponseWriter, r *http.Request) {
	h.BlogHandler().TagPostsPartial(w, r)
}

func (h *Handler) RedirectTag(w http.ResponseWriter, r *http.Request) {
	h.BlogHandler().RedirectTag(w, r)
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	h.BlogHandler().DeleteComment(w, r)
}

func (h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {
	h.BlogHandler().PostComment(w, r)
}

func (h *Handler) CommentsPartial(w http.ResponseWriter, r *http.Request) {
	h.BlogHandler().CommentsPartial(w, r)
}

func (h *Handler) EditCommentForm(w http.ResponseWriter, r *http.Request) {
	h.BlogHandler().EditCommentForm(w, r)
}

func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	h.BlogHandler().UpdateComment(w, r)
}

func (h *Handler) ReportComment(w http.ResponseWriter, r *http.Request) {
	h.BlogHandler().ReportComment(w, r)
}

func (h *Handler) StreamComments(w http.ResponseWriter, r *http.Request) {
	h.BlogHandler().StreamComments(w, r)
}

func (h *Handler) PostReaction(w http.ResponseWriter, r *http.Request) {
	h.BlogHandler().PostReaction(w, r)
}
