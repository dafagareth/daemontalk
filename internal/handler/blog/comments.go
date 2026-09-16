package blog

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"

	"daemontalk/internal/auth"
	"daemontalk/internal/comment"
	"daemontalk/internal/handler/common"
	"daemontalk/internal/i18n"
	"daemontalk/internal/post"
	blogtmpl "daemontalk/web/templates/blog"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	lang := common.LangFromRequest(r)
	ui := i18n.Get(lang)
	slug := chi.URLParam(r, "slug")

	isAdmin := h.isAdmin(r)
	if h.Comments == nil {
		http.Error(w, "comments unavailable", http.StatusServiceUnavailable)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	visitorName := common.GetVisitorIdentity(w, r)
	user := auth.GetUser(r.Context())

	if !isAdmin {
		c, err := h.Comments.GetByID(id)
		if err != nil || c == nil {
			http.Error(w, "comment not found", http.StatusNotFound)
			return
		}
		isOwner := (user != nil && c.UserID != nil && *c.UserID == user.ID) || (c.Name == visitorName)
		if !isOwner {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
	}

	if err := h.Comments.Delete(id); err != nil {
		slog.Error("delete comment failed", "id", id, "error", err)
	}
	h.renderCommentList(w, r, ui, slug, isAdmin)
}

func (h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {
	lang := common.LangFromRequest(r)
	ui := i18n.Get(lang)
	slug := chi.URLParam(r, "slug")

	if _, ok := post.FindBySlug(h.getAllPosts(), slug); !ok {
		http.NotFound(w, r)
		return
	}
	if h.Comments == nil {
		http.Error(w, "comments unavailable", http.StatusServiceUnavailable)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 64*1024)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if r.PostFormValue("website") != "" {
		h.renderCommentList(w, r, ui, slug, h.isAdmin(r))
		return
	}

	isAdmin := h.isAdmin(r)
	authUser := auth.GetUser(r.Context())

	name := common.GetVisitorIdentity(w, r)
	var avatarURL, ghURL string
	var isVerified bool
	var userID *int64

	isAnon := r.PostFormValue("is_anonymous") == "true"

	if authUser != nil && !isAnon {
		name = authUser.Username
		avatarURL = authUser.AvatarURL
		ghURL = authUser.GitHubURL
		isVerified = true
		userID = &authUser.ID
	} else if authUser != nil && isAnon {
		userID = &authUser.ID
	} else if isAdmin {
		name = "daemontalk"
		isVerified = true
	}
	body := r.PostFormValue("body")

	var parentID *int64
	if pIDStr := strings.TrimSpace(r.PostFormValue("parent_id")); pIDStr != "" {
		if pid, err := strconv.ParseInt(pIDStr, 10, 64); err == nil && pid > 0 {
			parentID = &pid
		}
	}

	if common.SpamScore(name, body) > common.SpamThreshold {
		h.renderCommentList(w, r, ui, slug, isAdmin)
		return
	}

	if _, err := h.Comments.AddAdvanced(comment.Comment{
		PostSlug:   slug,
		Name:       name,
		Body:       body,
		ParentID:   parentID,
		UserID:     userID,
		AvatarURL:  avatarURL,
		IsVerified: isVerified,
		GitHubURL:  ghURL,
	}); err != nil {
		if err == comment.ErrInvalid {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			slog.Error("add comment failed", "slug", slug, "parent_id", parentID, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		broadcastNewComment(slug)
		if h.SMTPHost != "" && h.SMTPTo != "" {
			go h.sendCommentNotification(slug, name, body)
		}
	}
	h.renderCommentList(w, r, ui, slug, isAdmin)
}

func (h *Handler) renderCommentList(w http.ResponseWriter, r *http.Request, ui i18n.UI, slug string, isAdmin bool) {
	visitorName := common.GetVisitorIdentity(w, r)
	if isAdmin {
		visitorName = "daemontalk"
	}
	authUser := auth.GetUser(r.Context())

	comments, err := h.Comments.ListBySlug(slug)
	if err != nil {
		slog.Error("load comments for slug failed", "slug", slug, "error", err)
	}
	h.render(w, r, blogtmpl.CommentList(ui, comments, isAdmin, slug, common.LangFromRequest(r), visitorName, authUser))
}

func (h *Handler) sendCommentNotification(slug, name, body string) {
	port := h.SMTPPort
	if port == "" {
		port = "587"
	}
	subject := common.StripCRLF(fmt.Sprintf("New comment on: %s", slug))
	cleanedBody := cleanEmailBody(body)
	msgBody := fmt.Sprintf("Post: %s\r\nFrom: %s\r\n\r\n%s", common.StripCRLF(slug), common.StripCRLF(name), cleanedBody)
	msg := []byte("To: " + h.SMTPTo + "\r\n" +
		"From: " + h.SMTPUser + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		msgBody)
	auth := smtp.PlainAuth("", h.SMTPUser, h.SMTPPass, h.SMTPHost)
	if err := smtp.SendMail(h.SMTPHost+":"+port, auth, h.SMTPUser, []string{h.SMTPTo}, msg); err != nil {
		slog.Error("send comment notification failed", "slug", slug, "error", err)
	}
}

func cleanEmailBody(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r == '\n' || r == '\r' || r == '\t' || (r >= 32 && r != 127) {
			b.WriteRune(r)
		}
	}
	return b.String()
}
