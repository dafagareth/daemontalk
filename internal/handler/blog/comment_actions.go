package blog

import (
	"fmt"
	"html"
	"net/http"
	"strconv"
	"strings"
	"time"

	"daemontalk/internal/auth"
	"daemontalk/internal/comment"
	"daemontalk/internal/handler/common"
	"daemontalk/internal/i18n"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) CommentsPartial(w http.ResponseWriter, r *http.Request) {
	h.renderCommentList(w, r, i18n.Get(common.LangFromRequest(r)), chi.URLParam(r, "slug"), h.isAdmin(r))
}

func (h *Handler) EditCommentForm(w http.ResponseWriter, r *http.Request) {
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

	visitorName := common.GetVisitorIdentity(w, r)
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

	formHTML := fmt.Sprintf(`
		<form class="mt-2 space-y-2" hx-post="/blog/%s/comments/%d/update" hx-target="#comment-list" hx-swap="outerHTML">
			<textarea name="body" required maxlength="2000" rows="3" class="w-full px-3 py-2 text-sm rounded-none border border-border bg-surface text-text focus:outline-none focus:border-text resize-y">%s</textarea>
			<div class="flex items-center justify-end gap-2">
				<button type="button" class="text-xs font-mono text-muted hover:text-text cursor-pointer px-3 py-1.5 border border-border bg-surface transition-colors" hx-get="/blog/%s/comments" hx-target="#comment-list" hx-swap="outerHTML">Cancel</button>
				<button type="submit" class="px-4 py-1.5 text-xs font-sans bg-slate-600 hover:bg-slate-500 text-white font-bold transition-colors cursor-pointer rounded-none">Save</button>
			</div>
		</form>
	`, slug, id, html.EscapeString(c.Body), slug)

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

	visitorName := common.GetVisitorIdentity(w, r)
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

	lang := common.LangFromRequest(r)
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
