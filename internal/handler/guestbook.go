package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"daemontalk/internal/comment"
	"daemontalk/internal/i18n"
	"daemontalk/web/templates"
)

func (h *Handler) Guestbook(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)

	var entries []comment.Comment
	if h.Comments != nil {
		if cs, err := h.Comments.ListBySlug(GuestbookSlug); err != nil {
			slog.Error("guestbook list query failed", "error", err)
		} else {
			entries = cs
		}
	}

	visitorName := GetVisitorIdentity(w, r)
	if h.isAdmin(r) {
		visitorName = "daemontalk"
	}

	meta := templates.PageMeta{Description: "Sign the guestbook and leave a message at daemontalk.com"}
	h.Render(w, r, templates.Layout(ui, lang, "guestbook", r.URL.Path, meta,
		templates.GuestbookPage(ui, entries, lang, visitorName)))
}

func (h *Handler) PostGuestbook(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)

	if h.Comments == nil {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 64*1024)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Honeypot: silently drop bot submissions.
	if r.PostFormValue("website") != "" {
		h.renderGuestbookList(w, r, ui)
		return
	}

	name := GetVisitorIdentity(w, r)
	customName := strings.TrimSpace(r.PostFormValue("name"))
	if customName != "" && len(customName) <= 30 {
		lower := strings.ToLower(customName)
		if !h.isAdmin(r) && (lower == "daemontalk" || lower == "dafa gareth" || lower == "admin" || lower == "author") {
			customName = name
		}
		name = customName
	}
	if h.isAdmin(r) {
		name = "daemontalk"
	}

	category := strings.TrimSpace(r.PostFormValue("category"))
	if len([]rune(category)) > 20 {
		category = string([]rune(category)[:20])
	}
	body := strings.TrimSpace(r.PostFormValue("body"))
	if len([]rune(body)) > comment.MaxBodyLen {
		body = string([]rune(body)[:comment.MaxBodyLen])
	}
	if category != "" && !strings.HasPrefix(body, "[") {
		body = fmt.Sprintf("[%s] %s", strings.ToUpper(category), body)
	}

	// Spam check: silently drop high-risk submissions.
	if spamScore(name, body) > spamThreshold {
		h.renderGuestbookList(w, r, ui)
		return
	}

	if _, err := h.Comments.Add(GuestbookSlug, name, body); err != nil {
		if err == comment.ErrInvalid {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			slog.Error("add guestbook entry failed", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	h.renderGuestbookList(w, r, ui)
}

func (h *Handler) renderGuestbookList(w http.ResponseWriter, r *http.Request, ui i18n.UI) {
	var entries []comment.Comment
	if h.Comments != nil {
		if cs, err := h.Comments.ListBySlug(GuestbookSlug); err != nil {
			slog.Error("guestbook render list query failed", "error", err)
		} else {
			entries = cs
		}
	}
	h.Render(w, r, templates.GuestbookList(ui, entries))
}
