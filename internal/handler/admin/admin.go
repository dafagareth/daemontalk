package admin

import (
	"log/slog"
	"net/http"
	"strconv"

	"daemontalk/internal/comment"
	"daemontalk/internal/handler/common"
	"daemontalk/internal/post"
	"daemontalk/web/templates/admin"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) Admin(w http.ResponseWriter, r *http.Request) {
	if h.AdminToken != "" {
		if tok := r.URL.Query().Get("admin"); tok != "" {
			if tok == h.AdminToken {
				common.SetAdminCookie(w, r, tok)
			}
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
			return
		}
	}

	if !h.IsAdmin(r) {
		h.HandleNotFound(w, r)
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

	var archivedPosts []post.Post
	if arc, err := post.LoadArchived(h.GetContentPath("posts")); err == nil {
		archivedPosts = arc
	}

	stats := admin.Stats{
		Posts:         h.GetAllPosts(),
		FilePosts:     h.GetFilePosts(),
		ArchivedPosts: archivedPosts,
		Views:         views,
		Comments:      allComments,
		TopPages:      topPages,
		TotalHits:     totalHits,
	}

	common.Render(w, r, admin.Layout("admin", r.URL.Path, admin.Page(stats)))
}

func (h *Handler) AdminDeleteComment(w http.ResponseWriter, r *http.Request) {
	if !h.IsAdmin(r) {
		h.HandleNotFound(w, r)
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
