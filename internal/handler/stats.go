package handler

import (
	"log/slog"
	"net/http"
	"time"

	"daemontalk/internal/i18n"
	"daemontalk/web/templates"
)

func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)
	isAdmin := h.isAdmin(r)

	data := templates.StatsData{}
	tagSet := make(map[string]bool)

	for _, p := range h.AllPosts() {
		if p.Draft && !isAdmin {
			continue
		}
		if !p.PublishAt.IsZero() && p.PublishAt.After(time.Now()) && !isAdmin {
			continue
		}
		if p.Type == "til" {
			data.TotalTIL++
		} else {
			data.TotalPosts++
		}
		data.TotalWords += p.ReadTime * 200
		for _, t := range p.Tags {
			tagSet[t] = true
		}
		switch p.Lang {
		case "en":
			data.LangEN++
		case "id":
			data.LangID++
		}
	}
	data.TotalTags = len(tagSet)

	if h.Comments != nil {
		if pv, err := h.Comments.TopPageViews(5); err == nil {
			for _, p := range pv {
				data.TopPages = append(data.TopPages, templates.StatTopPage{Path: p.Path, Count: p.Count})
			}
		} else {
			slog.Error("stats top pageviews query failed", "error", err)
		}
		if tv, err := h.Comments.TotalPageViews(); err == nil {
			data.TotalViews = tv
		}
		if cs, err := h.Comments.ListAll(); err == nil {
			for _, c := range cs {
				if c.PostSlug != GuestbookSlug {
					data.TotalComments++
				}
			}
		}
	}

	meta := templates.PageMeta{Description: "Site statistics for daemontalk.com — posts, views, and more."}
	h.Render(w, r, templates.Layout(ui, lang, "stats", r.URL.Path, meta, templates.StatsPage(ui, data)))
}

