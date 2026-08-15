package handler

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"daemontalk/internal/postdb"
	"daemontalk/web/templates"
)

var slugCleanRe = regexp.MustCompile(`[^a-z0-9-]+`)

// generateShortID generates a random 8-character hex string (e.g. "8f2a1b4c").
func generateShortID() string {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%x", time.Now().UnixNano())[:8]
	}
	return hex.EncodeToString(b)
}

// uniqueShortID generates a collision-free short ID.
func (h *Handler) uniqueShortID(excludeID int64) string {
	for {
		id := generateShortID()
		if !h.slugTaken(id, excludeID) {
			return id
		}
	}
}

// slugify menghasilkan slug kebab-case dari judul (fallback server-side untuk
// generator JS di form editor).
func slugify(title string) string {
	s := strings.ToLower(strings.TrimSpace(title))
	s = strings.ReplaceAll(s, " ", "-")
	s = slugCleanRe.ReplaceAllString(s, "")
	s = strings.Trim(strings.ReplaceAll(s, "--", "-"), "-")
	return s
}

// editorMD merender markdown menjadi HTML polos untuk dimuat ke Quill —
// sengaja tanpa chroma/heading-id/lazy-img supaya round-trip HTML↔markdown
// di editor tetap bersih.
var editorMD = goldmark.New(goldmark.WithExtensions(extension.Strikethrough))

func mdToEditorHTML(md string) string {
	var buf bytes.Buffer
	if err := editorMD.Convert([]byte(md), &buf); err != nil {
		log.Printf("editor html: %v", err)
		return ""
	}
	return buf.String()
}

// slugTaken melapor apakah slug sudah dipakai post file atau post DB lain.
func (h *Handler) slugTaken(slug string, excludeID int64) bool {
	for _, fp := range h.FilePosts {
		if fp.Slug == slug {
			return true
		}
	}
	if h.PostDB != nil {
		if existing, err := h.PostDB.GetBySlug(slug); err == nil && existing.ID != excludeID {
			return true
		}
	}
	return false
}

// uniqueSlug mencari slug bebas dengan sufiks -2, -3, … — dipakai autosave
// supaya pembuatan draft tidak pernah gagal karena slug bentrok.
func (h *Handler) uniqueSlug(base string, excludeID int64) string {
	if base == "" {
		return h.uniqueShortID(excludeID)
	}
	candidate := base
	for i := 2; h.slugTaken(candidate, excludeID); i++ {
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
	return candidate
}

func (h *Handler) renderEditor(w http.ResponseWriter, r *http.Request, p postdb.WebPost, errMsg string, status int) {
	w.WriteHeader(status)
	if err := templates.AdminLayout("admin", r.URL.Path,
		templates.AdminPostEditor(p, mdToEditorHTML(p.BodyMD), errMsg)).Render(r.Context(), w); err != nil {
		log.Printf("render error: %v", err)
	}
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

// AdminPostEdit displays the editor for an existing DB post.
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
		http.NotFound(w, r)
		return
	}
	h.renderEditor(w, r, p, "", http.StatusOK)
}

// AdminPostAutosave handles JSON autosave from Quill/title inputs.
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
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 2<<20)).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	var p postdb.WebPost
	if req.ID == 0 {
		title := strings.TrimSpace(req.Title)
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
			log.Printf("autosave create: %v", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		p.ID = id
	} else {
		existing, err := h.PostDB.Get(req.ID)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if reqTitle := strings.TrimSpace(req.Title); reqTitle != "" {
			existing.Title = reqTitle
		}
		existing.BodyMD = req.Body
		if err := h.PostDB.Update(existing); err != nil {
			log.Printf("autosave update %d: %v", req.ID, err)
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
func (h *Handler) validateWebPost(p postdb.WebPost, excludeID int64) string {
	if strings.TrimSpace(p.Title) == "" {
		return "Title is required."
	}
	if !p.Draft && strings.TrimSpace(p.BodyMD) == "" {
		return "Body content is required before publishing."
	}
	if p.Slug == "" {
		return "Slug cannot be empty."
	}
	if h.slugTaken(p.Slug, excludeID) {
		return fmt.Sprintf("Slug %q is already in use by another post.", p.Slug)
	}
	return ""
}

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
		http.NotFound(w, r)
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
		log.Printf("publish post %d: %v", id, err)
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
		log.Printf("admin delete post %d: %v", id, err)
	}
	h.RefreshPosts()
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
