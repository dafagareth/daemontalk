package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strconv"

	"daemontalk/internal/comment"
	"daemontalk/internal/i18n"
	"daemontalk/internal/post"
	"daemontalk/web/templates"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)
	slug := chi.URLParam(r, "slug")

	if !h.isAdmin(r) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if h.Comments == nil {
		http.Error(w, "comments unavailable", http.StatusServiceUnavailable)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if err := h.Comments.Delete(id); err != nil {
		log.Printf("delete comment %d: %v", id, err)
	}
	h.renderCommentList(w, r, ui, slug, true)
}
func (h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)
	slug := chi.URLParam(r, "slug")

	if _, ok := post.FindBySlug(h.AllPosts(), slug); !ok {
		http.NotFound(w, r)
		return
	}
	if h.Comments == nil {
		http.Error(w, "comments unavailable", http.StatusServiceUnavailable)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Honeypot: a filled "website" field means a bot — silently drop it but
	// return the unchanged list so the bot gets no signal.
	if r.PostFormValue("website") != "" {
		h.renderCommentList(w, r, ui, slug, h.isAdmin(r))
		return
	}

	name := GetVisitorIdentity(w, r)
	body := r.PostFormValue("body")

	// Spam check: score high-risk submissions silently.
	if spamScore(name, body) > spamThreshold {
		h.renderCommentList(w, r, ui, slug, h.isAdmin(r))
		return
	}

	if _, err := h.Comments.Add(slug, name, body); err != nil {
		if err == comment.ErrInvalid {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			log.Printf("add comment for %s: %v", slug, err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		// Still return the current list so the UI stays consistent.
	} else if h.SMTPHost != "" && h.SMTPTo != "" {
		go h.sendCommentNotification(slug, name, body)
	}
	h.renderCommentList(w, r, ui, slug, h.isAdmin(r))
}
func (h *Handler) renderCommentList(w http.ResponseWriter, r *http.Request, ui i18n.UI, slug string, isAdmin bool) {
	visitorName := GetVisitorIdentity(w, r)
	comments, err := h.Comments.ListBySlug(slug)
	if err != nil {
		log.Printf("load comments for %s: %v", slug, err)
	}
	if err := templates.CommentList(ui, comments, isAdmin, slug, langFromRequest(r), visitorName).Render(r.Context(), w); err != nil {
		log.Printf("render error: %v", err)
	}
}
func (h *Handler) sendCommentNotification(slug, name, body string) {
	port := h.SMTPPort
	if port == "" {
		port = "587"
	}
	subject := stripCRLF(fmt.Sprintf("New comment on: %s", slug))
	msgBody := fmt.Sprintf("Post: %s\nFrom: %s\n\n%s", stripCRLF(slug), stripCRLF(name), body)
	msg := []byte("To: " + h.SMTPTo + "\r\n" +
		"From: " + h.SMTPUser + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		msgBody)
	auth := smtp.PlainAuth("", h.SMTPUser, h.SMTPPass, h.SMTPHost)
	if err := smtp.SendMail(h.SMTPHost+":"+port, auth, h.SMTPUser, []string{h.SMTPTo}, msg); err != nil {
		log.Printf("comment notification: %v", err)
	}
}
