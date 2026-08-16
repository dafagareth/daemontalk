package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func (h *Handler) AdminPostExportMD(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		h.NotFound(w, r)
		return
	}

	slug := cleanSlug(r.URL.Query().Get("slug"))
	if slug == "" || !safeSlugRegex.MatchString(slug) {
		http.Error(w, "Invalid post slug", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join("content/posts", slug+".md")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		filePath = filepath.Join("content/posts", slug+".md.archive")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.md\"", slug))
	_, _ = w.Write(data)
}

func (h *Handler) AdminPostFileArchive(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		h.NotFound(w, r)
		return
	}

	slug := cleanSlug(r.FormValue("slug"))
	if slug == "" || !safeSlugRegex.MatchString(slug) {
		http.Error(w, "Invalid post slug", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join("content/posts", slug+".md")
	if _, err := os.Stat(filePath); err == nil {
		archPath := filepath.Join("content/posts", slug+".md.archive")
		_ = os.Rename(filePath, archPath)
	}

	h.ReloadFilePosts()
	h.RefreshPosts()
	http.Redirect(w, r, "/admin#content", http.StatusSeeOther)
}

func (h *Handler) AdminPostFileRestore(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		h.NotFound(w, r)
		return
	}

	slug := cleanSlug(r.FormValue("slug"))
	if slug == "" || !safeSlugRegex.MatchString(slug) {
		http.Error(w, "Invalid post slug", http.StatusBadRequest)
		return
	}

	archPath := filepath.Join("content/posts", slug+".md.archive")
	if _, err := os.Stat(archPath); err == nil {
		filePath := filepath.Join("content/posts", slug+".md")
		_ = os.Rename(archPath, filePath)
	}

	h.ReloadFilePosts()
	h.RefreshPosts()
	http.Redirect(w, r, "/admin#content", http.StatusSeeOther)
}

func (h *Handler) AdminPostFileDelete(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		h.NotFound(w, r)
		return
	}

	slug := cleanSlug(r.FormValue("slug"))
	if slug == "" || !safeSlugRegex.MatchString(slug) {
		http.Error(w, "Invalid post slug", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join("content/posts", slug+".md")
	archPath := filepath.Join("content/posts", slug+".md.archive")
	_ = os.Remove(filePath)
	_ = os.Remove(archPath)

	imgDir := filepath.Join("web/static/images/posts", slug)
	_ = os.RemoveAll(imgDir)

	h.ReloadFilePosts()
	h.RefreshPosts()
	http.Redirect(w, r, "/admin#content", http.StatusSeeOther)
}
