package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"daemontalk/internal/postdb"
	"daemontalk/web/templates"

	"github.com/go-chi/chi/v5"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var slugCleanRe = regexp.MustCompile(`[^a-z0-9-]+`)

// generateShortID generates a random 8-character hex string (e.g. "8f2a1b4c").

// uniqueShortID generates a collision-free short ID.

// slugify generates kebab-case slug from title (server-side fallback for
// generator JS di form editor).

// editorMD renders markdown into plain HTML to load into Quill —
// intentionally without chroma/heading-id/lazy-img to allow round-trip HTML↔markdown
// di editor tetap bersih.
var editorMD = goldmark.New(goldmark.WithExtensions(extension.Strikethrough))

// slugTaken reports whether slug is already used by a post file or another DB post.

// uniqueSlug finds a free slug with suffix -2, -3, … — used by autosave
// so draft creation never fails due to slug collision.

func (h *Handler) renderEditor(w http.ResponseWriter, r *http.Request, p postdb.WebPost, errMsg string, status int) {
	w.WriteHeader(status)
	h.Render(w, r, templates.AdminLayout("admin", r.URL.Path,
		templates.AdminPostEditor(p, mdToEditorHTML(p.BodyMD), errMsg)))
}

// AdminPostNew displays an empty write page (draft created upon first autosave).
func (h *Handler) AdminPostNew(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		h.NotFound(w, r)
		return
	}
	p := postdb.WebPost{
		Slug:  h.uniqueShortID(0),
		Lang:  "en",
		Draft: true,
		Date:  time.Now().Format("2006-01-02"),
	}
	h.renderEditor(w, r, p, "", http.StatusOK)
}

// AdminPostEdit displays the editor prefilled with an existing post.
func (h *Handler) AdminPostEdit(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		h.NotFound(w, r)
		return
	}
	if h.PostDB == nil {
		http.Error(w, "post db disabled", http.StatusServiceUnavailable)
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	p, err := h.PostDB.Get(id)
	if err != nil {
		h.NotFound(w, r)
		return
	}
	h.renderEditor(w, r, p, "", http.StatusOK)
}

// AdminPostAutosave accepts JSON payload from the editor, silently creates a draft
// if ID == 0 or updates an existing draft, and returns the assigned ID.
func (h *Handler) AdminPostAutosave(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		h.NotFound(w, r)
		return
	}
	if h.PostDB == nil {
		http.Error(w, "post db disabled", http.StatusServiceUnavailable)
		return
	}

	var req struct {
		ID    int64  `json:"id"`
		Title string `json:"title"`
		Body  string `json:"body"`
		Slug  string `json:"slug"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	var p postdb.WebPost
	if req.ID == 0 {
		title := strings.TrimSpace(req.Title)
		if title == "" {
			title = "Draft " + time.Now().Format("02 Jan 15:04")
		}
		slug := strings.TrimSpace(req.Slug)
		if slug != "" {
			slug = h.uniqueSlug(slugify(slug), 0)
		} else {
			slug = h.uniqueShortID(0)
		}
		p = postdb.WebPost{
			Slug:   slug,
			Title:  title,
			BodyMD: req.Body,
			Lang:   "en",
			Draft:  true,
			Date:   time.Now().Format("2006-01-02"),
		}
		id, err := h.PostDB.Create(p)
		if err != nil {
			slog.Error("autosave create post failed", "error", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		p.ID = id
	} else {
		existing, err := h.PostDB.Get(req.ID)
		if err != nil {
			h.NotFound(w, r)
			return
		}
		if reqTitle := strings.TrimSpace(req.Title); reqTitle != "" {
			existing.Title = reqTitle
		}
		existing.BodyMD = req.Body
		if err := h.PostDB.Update(existing); err != nil {
			slog.Error("autosave update post failed", "id", req.ID, "error", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		p = existing
	}

	h.RefreshPosts()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"id":      p.ID,
		"slug":    p.Slug,
		"savedAt": time.Now().Format("15:04"),
	})
}

// validateWebPost returns an error message (empty = valid). excludeID
// ignores the post itself when updating.

// AdminPostPublish handles publish modal form: metadata + action
// (publish/draft). Title and body are retrieved from the last autosave.
// Slug can only be changed while draft; once published, slug is locked to keep URL stable.
func (h *Handler) AdminPostPublish(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		h.NotFound(w, r)
		return
	}
	if h.PostDB == nil {
		http.Error(w, "post db disabled", http.StatusServiceUnavailable)
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	p, err := h.PostDB.Get(id)
	if err != nil {
		h.NotFound(w, r)
		return
	}

	if t := strings.TrimSpace(r.FormValue("title")); t != "" {
		p.Title = t
	}
	if b := r.FormValue("body"); strings.TrimSpace(b) != "" {
		p.BodyMD = b
	}

	if p.Draft {
		slug := slugify(r.FormValue("slug"))
		if slug == "" {
			slug = p.Slug
		}
		if slug == "" {
			slug = h.uniqueShortID(id)
		}
		p.Slug = slug
	}
	if date := r.FormValue("date"); date != "" {
		if _, err := time.Parse("2006-01-02", date); err == nil {
			p.Date = date
		}
	}
	if lang := r.FormValue("lang"); lang == "id" {
		p.Lang = "id"
	} else {
		p.Lang = "en"
	}
	p.Tags = strings.TrimSpace(r.FormValue("tags"))
	p.Cover = strings.TrimSpace(r.FormValue("cover"))
	p.Description = strings.TrimSpace(r.FormValue("description"))
	p.Draft = r.FormValue("action") != "publish"

	if msg := h.validateWebPost(p, id); msg != "" {
		h.renderEditor(w, r, p, msg, http.StatusUnprocessableEntity)
		return
	}
	if err := h.PostDB.Update(p); err != nil {
		slog.Error("publish post update failed", "id", id, "error", err)
		h.renderEditor(w, r, p, "Failed to save to database.", http.StatusInternalServerError)
		return
	}
	h.RefreshPosts()

	if p.Draft {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/blog/"+p.Slug, http.StatusSeeOther)
}

// AdminPostDelete menghapus post DB.
func (h *Handler) AdminPostDelete(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		h.NotFound(w, r)
		return
	}
	if h.PostDB == nil {
		http.Error(w, "post db disabled", http.StatusServiceUnavailable)
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if err := h.PostDB.Delete(id); err != nil {
		slog.Error("admin delete post failed", "id", id, "error", err)
	}
	h.RefreshPosts()
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
