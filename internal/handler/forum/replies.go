package forum

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"daemontalk/internal/auth"
	"daemontalk/internal/forum"
	"daemontalk/internal/handler/common"
	"daemontalk/internal/i18n"
	forumtmpl "daemontalk/web/templates/forum"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) DiscussionsReply(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	topicIDStr := chi.URLParam(r, "id")
	topicID, err := strconv.ParseInt(topicIDStr, 10, 64)
	if err != nil || topicID <= 0 {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	bodyMD := strings.TrimSpace(r.FormValue("body"))
	if bodyMD == "" {
		http.Error(w, "Reply content cannot be empty", http.StatusBadRequest)
		return
	}

	var parentID int64
	if pStr := r.FormValue("parent_id"); pStr != "" {
		if p, err := strconv.ParseInt(pStr, 10, 64); err == nil {
			parentID = p
		}
	}

	reply, err := h.Forum.CreateReply(forum.Reply{
		TopicID:  topicID,
		ParentID: parentID,
		UserID:   user.ID,
		BodyMD:   bodyMD,
	})
	if err != nil {
		slog.Error("failed to create reply", "error", err)
		http.Error(w, "Failed to submit reply", http.StatusInternalServerError)
		return
	}

	reply.AuthorName = user.DisplayName
	reply.AuthorUsername = user.Username
	reply.AuthorAvatar = user.AvatarURL
	reply.AuthorGitHub = user.GitHubURL
	reply.IsOwner = true

	if r.Header.Get("HX-Request") == "true" {
		lang := common.LangFromRequest(r)
		ui := i18n.Get(lang)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = forumtmpl.ReplyItem(ui, lang, user, reply, topicID, user.ID).Render(r.Context(), w)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func (h *Handler) DiscussionsSolve(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	topicIDStr := chi.URLParam(r, "id")
	topicID, err := strconv.ParseInt(topicIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	replyIDStr := r.URL.Query().Get("reply_id")
	replyID, _ := strconv.ParseInt(replyIDStr, 10, 64)

	if err := h.Forum.MarkSolution(topicID, replyID, user.ID); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Refresh", "true")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func (h *Handler) DiscussionsVote(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Please login with GitHub to vote", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	targetType := r.FormValue("type")
	targetID, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err != nil || (targetType != "topic" && targetType != "reply") {
		http.Error(w, "Invalid vote target", http.StatusBadRequest)
		return
	}

	newCount, hasVoted, err := h.Forum.Vote(user.ID, targetType, targetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = forumtmpl.VoteButton(targetType, targetID, newCount, hasVoted).Render(r.Context(), w)
}

func (h *Handler) DiscussionsDeleteReply(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid reply ID", http.StatusBadRequest)
		return
	}

	isAdmin := h.IsAdmin(r) || user.Role == "admin"
	if err := h.Forum.DeleteReply(id, user.ID, isAdmin); err != nil {
		slog.Error("delete forum reply failed", "id", id, "error", err)
		http.Error(w, "Failed to delete reply", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Refresh", "true")
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}
