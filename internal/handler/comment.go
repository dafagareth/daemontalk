package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"sync"
	"time"

	"daemontalk/internal/comment"
	"daemontalk/internal/i18n"
	"daemontalk/internal/post"
	"daemontalk/web/templates"

	"github.com/go-chi/chi/v5"
)

var (
	commentSubscribers   = make(map[chan string]string)
	commentSubscribersMu sync.RWMutex
)

func broadcastNewComment(slug string) {
	commentSubscribersMu.RLock()
	defer commentSubscribersMu.RUnlock()
	for ch, s := range commentSubscribers {
		if s == slug {
			select {
			case ch <- slug:
			default:
			}
		}
	}
}

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
		slog.Error("delete comment failed", "id", id, "error", err)
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

	isAdmin := h.isAdmin(r)
	name := GetVisitorIdentity(w, r)
	if isAdmin {
		name = "daemontalk"
	}
	body := r.PostFormValue("body")

	var parentID *int64
	if pIDStr := strings.TrimSpace(r.PostFormValue("parent_id")); pIDStr != "" {
		if pid, err := strconv.ParseInt(pIDStr, 10, 64); err == nil && pid > 0 {
			parentID = &pid
		}
	}

	// Spam check: score high-risk submissions silently.
	if spamScore(name, body) > spamThreshold {
		h.renderCommentList(w, r, ui, slug, isAdmin)
		return
	}

	if _, err := h.Comments.AddWithParent(slug, name, body, parentID); err != nil {
		if err == comment.ErrInvalid {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			slog.Error("add comment failed", "slug", slug, "parent_id", parentID, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		// Still return the current list so the UI stays consistent.
	} else {
		broadcastNewComment(slug)
		if h.SMTPHost != "" && h.SMTPTo != "" {
			go h.sendCommentNotification(slug, name, body)
		}
	}
	h.renderCommentList(w, r, ui, slug, isAdmin)
}
func (h *Handler) renderCommentList(w http.ResponseWriter, r *http.Request, ui i18n.UI, slug string, isAdmin bool) {
	visitorName := GetVisitorIdentity(w, r)
	if isAdmin {
		visitorName = "daemontalk"
	}
	comments, err := h.Comments.ListBySlug(slug)
	if err != nil {
		slog.Error("load comments for slug failed", "slug", slug, "error", err)
	}
	h.Render(w, r, templates.CommentList(ui, comments, isAdmin, slug, langFromRequest(r), visitorName))
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
		slog.Error("send comment notification failed", "slug", slug, "error", err)
	}
}

func (h *Handler) StreamComments(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	// If accessed directly by human in browser address bar (HTML request), redirect to article
	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "text/html") && !strings.Contains(accept, "text/event-stream") {
		http.Redirect(w, r, "/blog/"+slug+"#comments", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache, no-transform")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Send initial comment byte to establish stream immediately
	fmt.Fprintf(w, ": connected\n\n")
	flusher.Flush()

	ch := make(chan string, 1)
	commentSubscribersMu.Lock()
	commentSubscribers[ch] = slug
	commentSubscribersMu.Unlock()

	defer func() {
		commentSubscribersMu.Lock()
		delete(commentSubscribers, ch)
		commentSubscribersMu.Unlock()
	}()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	notify := r.Context().Done()
	for {
		select {
		case <-notify:
			return
		case <-ticker.C:
			// Send SSE keepalive heartbeat so Cloudflare proxy never times out (Error 524)
			fmt.Fprintf(w, ": ping\n\n")
			flusher.Flush()
		case s := <-ch:
			fmt.Fprintf(w, "event: new_comment\ndata: %s\n\n", s)
			flusher.Flush()
		}
	}
}

func (h *Handler) CommentsPartial(w http.ResponseWriter, r *http.Request) {
	h.renderCommentList(w, r, i18n.Get(langFromRequest(r)), chi.URLParam(r, "slug"), h.isAdmin(r))
}
