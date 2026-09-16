package pages

import (
	"net/http"

	"daemontalk/internal/github"
	"daemontalk/internal/handler/common"
	"daemontalk/internal/i18n"
	"daemontalk/web/templates/layouts"
	pagestmpl "daemontalk/web/templates/pages"
	"daemontalk/web/templates/shared"
)

func (h *Handler) Colophon(w http.ResponseWriter, r *http.Request) {
	lang := common.LangFromRequest(r)
	ui := i18n.Get(lang)

	ghStats := github.Fetch("dafagareth", h.GitHubToken)

	meta := shared.PageMeta{
		Description: "DaemonTalk specifications, typography, licenses, and architecture.",
	}
	common.Render(w, r, layouts.Layout(ui, lang, "colophon", r.URL.Path, meta, pagestmpl.Colophon(ui, lang, ghStats)))
}
