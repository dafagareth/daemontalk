package handler

import (
	"fmt"
	"net/http"

	"daemontalk/internal/i18n"
	"daemontalk/web/templates"
)

// NotFound renders the custom 404 Not Found error page.
func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)
	w.WriteHeader(http.StatusNotFound)
	_ = templates.Layout(ui, lang, "404", r.URL.Path, templates.PageMeta{
		Description: ui.NotFound_Body,
	}, templates.NotFound(ui, lang)).Render(r.Context(), w)
}

// Forbidden renders the custom 403 Forbidden error page.
func (h *Handler) Forbidden(w http.ResponseWriter, r *http.Request, reason string) {
	if reason == "" {
		reason = "unauthorized_token_or_ip"
	}
	lang := langFromRequest(r)
	ui := i18n.Get(lang)
	w.WriteHeader(http.StatusForbidden)
	_ = templates.Layout(ui, lang, "403", r.URL.Path, templates.PageMeta{
		Description: ui.Forbidden_Body,
	}, templates.Forbidden(ui, lang, reason)).Render(r.Context(), w)
}

// InternalError renders the custom 500 Internal Server Error page.
func (h *Handler) InternalError(w http.ResponseWriter, r *http.Request, err error) {
	detail := "system_fault"
	if err != nil {
		detail = err.Error()
	}
	lang := langFromRequest(r)
	ui := i18n.Get(lang)
	w.WriteHeader(http.StatusInternalServerError)
	_ = templates.Layout(ui, lang, "500", r.URL.Path, templates.PageMeta{
		Description: ui.ServerError_Body,
	}, templates.ServerError(ui, lang, detail)).Render(r.Context(), w)
}

// CustomError renders a custom HTTP status error page with a specified title and message.
func (h *Handler) CustomError(w http.ResponseWriter, r *http.Request, status int, title, message string) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)
	w.WriteHeader(status)
	_ = templates.Layout(ui, lang, fmt.Sprintf("%d · %s", status, title), r.URL.Path, templates.PageMeta{
		Description: message,
	}, templates.CustomError(status, title, message, ui, lang)).Render(r.Context(), w)
}
