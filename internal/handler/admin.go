package handler

import (
	"log"
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
					Name:     "admin_token",
					Value:    tok,
					Path:     "/",
					HttpOnly: true,
					SameSite: http.SameSiteLaxMode,
					MaxAge:   60 * 60 * 24 * 30,
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
			log.Printf("admin: view counts: %v", err)
		} else {
			views = v
		}
		if cs, err := h.Comments.ListAll(); err != nil {
			log.Printf("admin: list comments: %v", err)
		} else {
			allComments = cs
		}
		if tp, err := h.Comments.TopPageViews(10); err != nil {
			log.Printf("admin: top pages: %v", err)
		} else {
			topPages = tp
		}
		if n, err := h.Comments.TotalPageViews(); err != nil {
			log.Printf("admin: total hits: %v", err)
		} else {
			totalHits = n
		}
	}

	var webPosts []postdb.WebPost
	if h.PostDB != nil {
		if wp, err := h.PostDB.List(); err != nil {
			log.Printf("admin: list web posts: %v", err)
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

	err := templates.AdminLayout("admin", r.URL.Path, templates.AdminPage(stats)).Render(r.Context(), w)
	if err != nil {
		log.Printf("render error: %v", err)
	}
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
			log.Printf("admin delete comment %d: %v", id, err)
		}
	}

	w.WriteHeader(http.StatusOK)
}
