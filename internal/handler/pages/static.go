package pages

import (
	"net/http"

	"daemontalk/internal/handler/common"
	"daemontalk/internal/i18n"
	"daemontalk/web/templates/layouts"
	pagestmpl "daemontalk/web/templates/pages"
	"daemontalk/web/templates/shared"
)

func (h *Handler) About(w http.ResponseWriter, r *http.Request) {
	common.RenderMarkdownPage(w, r, h.ContentDir, "about", "about", shared.PageMeta{
		Description: "About Daemontalk philosophy, editorial values, and open tech publishing.",
	}, pagestmpl.AboutPage)
}

func (h *Handler) Accessibility(w http.ResponseWriter, r *http.Request) {
	common.RenderMarkdownPage(w, r, h.ContentDir, "accessibility", "accessibility", shared.PageMeta{
		Description: "Accessibility statement and commitment to digital inclusion on daemontalk.com.",
	}, pagestmpl.AccessibilityPage)
}

func (h *Handler) License(w http.ResponseWriter, r *http.Request) {
	common.RenderMarkdownPage(w, r, h.ContentDir, "license", "license", shared.PageMeta{
		Description: "License and copyright information for content, code snippets, and projects.",
	}, pagestmpl.LicensePage)
}

func (h *Handler) Privacy(w http.ResponseWriter, r *http.Request) {
	common.RenderMarkdownPage(w, r, h.ContentDir, "privacy", "privacy", shared.PageMeta{
		Description: "Privacy policy for daemontalk.com — zero tracking, minimal data retention, and your rights.",
	}, pagestmpl.PrivacyPage)
}

func (h *Handler) Terms(w http.ResponseWriter, r *http.Request) {
	common.RenderMarkdownPage(w, r, h.ContentDir, "terms", "terms", shared.PageMeta{
		Description: "Terms of service and usage conditions for daemontalk.com and associated tools.",
	}, pagestmpl.TermsPage)
}

func (h *Handler) Saved(w http.ResponseWriter, r *http.Request) {
	lang := common.LangFromRequest(r)
	ui := i18n.Get(lang)
	common.Render(w, r, layouts.Layout(ui, lang, "saved", r.URL.Path, shared.PageMeta{
		Description: "Your saved posts reading list.",
	}, pagestmpl.SavedPage(ui, lang)))
}

func (h *Handler) ContactPage(w http.ResponseWriter, r *http.Request) {
	lang := common.LangFromRequest(r)
	ui := i18n.Get(lang)
	common.Render(w, r, layouts.Layout(ui, lang, "contact", r.URL.Path, shared.PageMeta{
		Description: "Send a message to daemontalk.",
	}, pagestmpl.ContactPage(ui, lang)))
}

func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	lang := common.LangFromRequest(r)
	ui := i18n.Get(lang)
	w.WriteHeader(http.StatusNotFound)
	_ = layouts.Layout(ui, lang, "404", r.URL.Path, shared.PageMeta{
		Description: ui.NotFound_Body,
	}, pagestmpl.NotFound(ui, lang)).Render(r.Context(), w)
}
