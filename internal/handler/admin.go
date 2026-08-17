package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"daemontalk/internal/comment"
	"daemontalk/internal/postdb"
	"daemontalk/web/templates"
)

func (h *Handler) Admin(w http.ResponseWriter, r *http.Request) {
	// Allow login via ?admin=TOKEN (same mechanism as blog posts)
	if h.AdminToken != "" {
		if tok := r.URL.Query().Get("admin"); tok != "" {
			if tok == h.AdminToken {
				http.SetCookie(w, &http.Cookie{
					Name:     CookieAdminToken,
					Value:    tok,
					Path:     "/",
					HttpOnly: true,
					SameSite: http.SameSiteLaxMode,
					MaxAge:   CookieAdminMaxAge,
				})
			}
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
			return
		}
	}

	if !h.isAdmin(r) {
		h.NotFound(w, r)
		return
	}

	views := make(map[string]int)
	var allComments []comment.Comment
	var topPages []comment.PageView
	totalHits := 0
	if h.Comments != nil {
		if v, err := h.Comments.AllViewCounts(); err != nil {
			slog.Error("admin view counts query failed", "error", err)
		} else {
			views = v
		}
		if cs, err := h.Comments.ListAll(); err != nil {
			slog.Error("admin list comments query failed", "error", err)
		} else {
			allComments = cs
		}
		if tp, err := h.Comments.TopPageViews(10); err != nil {
			slog.Error("admin top pages query failed", "error", err)
		} else {
			topPages = tp
		}
		if n, err := h.Comments.TotalPageViews(); err != nil {
			slog.Error("admin total hits query failed", "error", err)
		} else {
			totalHits = n
		}
	}

	var webPosts []postdb.WebPost
	if h.PostDB != nil {
		if wp, err := h.PostDB.List(); err != nil {
			slog.Error("admin list web posts failed", "error", err)
		} else {
			webPosts = wp
		}
	}

	stats := templates.AdminStats{
		Posts:        h.AllPosts(),
		WebPosts:     webPosts,
		Views:        views,
		Comments:     allComments,
		TopPages:     topPages,
		TotalHits:    totalHits,
		RadarEnabled: templates.IsRadarEnabled(),
	}

	h.Render(w, r, templates.AdminLayout("admin", r.URL.Path, templates.AdminPage(stats)))
}

// AdminToggleRadar toggles the systems radar feature flag.
func (h *Handler) AdminToggleRadar(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		h.NotFound(w, r)
		return
	}

	newState := !templates.IsRadarEnabled()
	templates.SetRadarEnabled(newState)

	http.Redirect(w, r, "/admin#dashboard", http.StatusSeeOther)
}

// AdminDeleteComment deletes a comment and returns an empty response so HTMX
// removes the card from the admin dashboard via outerHTML swap.
func (h *Handler) AdminDeleteComment(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		h.NotFound(w, r)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if h.Comments != nil {
		if err := h.Comments.Delete(id); err != nil {
			slog.Error("admin delete comment failed", "id", id, "error", err)
		}
	}

	w.WriteHeader(http.StatusOK)
}

