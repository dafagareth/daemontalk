package pages

import (
	"log/slog"
	"net/http"
	"time"

	"daemontalk/internal/handler/common"
	"daemontalk/internal/i18n"
	"daemontalk/web/templates/layouts"
	pagestmpl "daemontalk/web/templates/pages"
	"daemontalk/web/templates/shared"
)

func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	lang := common.LangFromRequest(r)
	ui := i18n.Get(lang)
	isAdmin := h.IsAdmin(r)

	data := pagestmpl.StatsData{}
	tagSet := make(map[string]bool)

	for _, p := range h.GetAllPosts() {
		if p.Draft && !isAdmin {
			continue
		}
		if !p.PublishAt.IsZero() && p.PublishAt.After(time.Now()) && !isAdmin {
			continue
		}
		data.TotalPosts++
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

	if h.Forum != nil {
		fst := h.Forum.GetStats()
		data.TotalTopics = fst.TotalTopics
		data.TotalReplies = fst.TotalReplies
		data.SolvedTopics = fst.SolvedTopics
		data.TotalVotes = fst.TotalVotes
	}

	if h.Auth != nil {
		data.TotalUsers = h.Auth.CountUsers()
	}

	if h.Comments != nil {
		if pv, err := h.Comments.TopPageViews(5); err == nil {
			for _, p := range pv {
				data.TopPages = append(data.TopPages, pagestmpl.StatTopPage{Path: p.Path, Count: p.Count})
			}
		} else {
			slog.Error("stats top pageviews query failed", "error", err)
		}
		if tv, err := h.Comments.TotalPageViews(); err == nil {
			data.TotalViews = tv
		}
		if cs, err := h.Comments.ListAll(); err == nil {
			for _, c := range cs {
				if c.PostSlug != common.GuestbookSlug {
					data.TotalComments++
				}
			}
		}
	}

	meta := shared.PageMeta{Description: "Site statistics for daemontalk.com — posts, views, and more."}
	common.Render(w, r, layouts.Layout(ui, lang, "stats", r.URL.Path, meta, pagestmpl.StatsPage(ui, lang, data)))
}
