package handler

import (
	"log"
	"net/http"

	"daemontalk/internal/comment"
	"daemontalk/internal/i18n"
	"daemontalk/web/templates"
)

func (h *Handler) Guestbook(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)

	var entries []comment.Comment
	if h.Comments != nil {
		if cs, err := h.Comments.ListBySlug("__guestbook__"); err != nil {
			log.Printf("guestbook list: %v", err)
		} else {
			entries = cs
		}
	}

	visitorName := GetVisitorIdentity(w, r)

	meta := templates.PageMeta{Description: "Sign the guestbook and leave a message at daemontalk.com"}
	err := templates.Layout(ui, lang, "guestbook", r.URL.Path, meta,
		templates.GuestbookPage(ui, entries, lang, visitorName)).Render(r.Context(), w)
	if err != nil {
		log.Printf("render error: %v", err)
	}
}

func (h *Handler) PostGuestbook(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)

	if h.Comments == nil {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
		return
	}
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
	body := r.PostFormValue("body")

	// Spam check: silently drop high-risk submissions.
	if spamScore(name, body) > spamThreshold {
		h.renderGuestbookList(w, r, ui)
		return
	}

	if _, err := h.Comments.Add("__guestbook__", name, body); err != nil {
		if err == comment.ErrInvalid {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			log.Printf("add guestbook entry: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	h.renderGuestbookList(w, r, ui)
}

func (h *Handler) renderGuestbookList(w http.ResponseWriter, r *http.Request, ui i18n.UI) {
	var entries []comment.Comment
	if h.Comments != nil {
		if cs, err := h.Comments.ListBySlug("__guestbook__"); err != nil {
			log.Printf("guestbook render list: %v", err)
		} else {
			entries = cs
		}
	}
	if err := templates.GuestbookList(ui, entries).Render(r.Context(), w); err != nil {
		log.Printf("render error: %v", err)
	}
}
