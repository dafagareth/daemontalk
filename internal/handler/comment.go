package handler

import (
	"fmt"
	"html"
	"log/slog"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"sync"
	"time"

	"daemontalk/internal/auth"
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

	visitorName := GetVisitorIdentity(w, r)
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

	r.Body = http.MaxBytesReader(w, r.Body, 64*1024)
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
	authUser := auth.GetUser(r.Context())

	name := GetVisitorIdentity(w, r)
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
		// Logged in but posting anonymously. We keep 'name' as visitorName (anonym_xyz),
		// but attach the userID so they still own the comment and can delete it.
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

	// Spam check: score high-risk submissions silently.
	if spamScore(name, body) > spamThreshold {
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
	authUser := auth.GetUser(r.Context())
	
	comments, err := h.Comments.ListBySlug(slug)
	if err != nil {
		slog.Error("load comments for slug failed", "slug", slug, "error", err)
	}
	h.Render(w, r, templates.CommentList(ui, comments, isAdmin, slug, langFromRequest(r), visitorName, authUser))
}
func (h *Handler) sendCommentNotification(slug, name, body string) {
	port := h.SMTPPort
	if port == "" {
		port = "587"
	}
	subject := stripCRLF(fmt.Sprintf("New comment on: %s", slug))
	cleanedBody := cleanEmailBody(body)
	msgBody := fmt.Sprintf("Post: %s\r\nFrom: %s\r\n\r\n%s", stripCRLF(slug), stripCRLF(name), cleanedBody)
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

// cleanEmailBody strips unprintable control characters to prevent header/boundary manipulation.
func cleanEmailBody(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r == '\n' || r == '\r' || r == '\t' || (r >= 32 && r != 127) {
			b.WriteRune(r)
		}
	}
	return b.String()
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

func blogPrefix(lang string) string {
	if lang == "id" {
		return "/id"
	}
	return ""
}

func (h *Handler) EditCommentForm(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	prefix := blogPrefix(lang)
	slug := chi.URLParam(r, "slug")
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	
	c, err := h.Comments.GetByID(id)
	if err != nil || c == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	visitorName := GetVisitorIdentity(w, r)
	user := auth.GetUser(r.Context())
	isOwner := (user != nil && c.UserID != nil && *c.UserID == user.ID) || (c.Name == visitorName)
	isAdmin := h.isAdmin(r)

	if !isAdmin && !isOwner {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	
	if !isAdmin && time.Since(c.CreatedAt) > 10*time.Minute {
		http.Error(w, "edit window expired", http.StatusForbidden)
		return
	}

	cancelLabel := "Cancel"
	saveLabel := "Save"
	if lang == "id" {
		cancelLabel = "Batal"
		saveLabel = "Simpan"
	}

	formHTML := fmt.Sprintf(`
		<form class="mt-2 space-y-2" hx-post="%s/blog/%s/comments/%d/update" hx-target="#comment-list" hx-swap="outerHTML">
			<textarea name="body" required maxlength="2000" rows="3" class="w-full px-3 py-2 text-sm rounded-none border border-[var(--c-link)] bg-surface text-text focus:outline-none focus:ring-1 focus:ring-[var(--c-link)] resize-y">%s</textarea>
			<div class="flex items-center justify-end gap-2">
				<button type="button" class="text-xs font-mono text-muted hover:text-text cursor-pointer px-3 py-1.5 border border-border bg-surface transition-colors" hx-get="%s/blog/%s/comments" hx-target="#comment-list" hx-swap="outerHTML">%s</button>
				<button type="submit" class="px-4 py-1.5 text-xs font-mono bg-[var(--c-link)] text-white hover:brightness-110 transition-colors cursor-pointer rounded-none font-medium">%s</button>
			</div>
		</form>
	`, prefix, slug, id, html.EscapeString(c.Body), prefix, slug, cancelLabel, saveLabel)
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(formHTML))
}

func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	c, err := h.Comments.GetByID(id)
	if err != nil || c == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	visitorName := GetVisitorIdentity(w, r)
	user := auth.GetUser(r.Context())
	isOwner := (user != nil && c.UserID != nil && *c.UserID == user.ID) || (c.Name == visitorName)
	isAdmin := h.isAdmin(r)

	if !isAdmin && !isOwner {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	
	if !isAdmin && time.Since(c.CreatedAt) > 10*time.Minute {
		http.Error(w, "edit window expired", http.StatusForbidden)
		return
	}

	body := strings.TrimSpace(r.FormValue("body"))
	if body != "" && len(body) <= comment.MaxBodyLen {
		_ = h.Comments.UpdateBody(id, body)
	}

	lang := langFromRequest(r)
	ui := i18n.Get(lang)
	h.renderCommentList(w, r, ui, slug, isAdmin)
}

func (h *Handler) ReportComment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err == nil {
		_ = h.Comments.Report(id)
	}
	
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<span class="text-[var(--c-link)] font-bold px-2 py-1 bg-surface border border-border mt-1">Reported!</span>`))
}
