package handler

import (
	"fmt"
	"log/slog"
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

	filePath := filepath.Join(h.getContentPath("posts"), slug+".md")
	data, err := os.ReadFile(filePath)
	if err != nil {
		archPath := filepath.Join(h.getContentPath("posts"), slug+".md.archive")
		data, err = os.ReadFile(archPath)
		if err != nil {
			http.Error(w, "Post file not found", http.StatusNotFound)
			return
		}
	}

	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.md\"", slug))
	w.WriteHeader(http.StatusOK)
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

	filePath := filepath.Join(h.getContentPath("posts"), slug+".md")
	if _, err := os.Stat(filePath); err == nil {
		archPath := filepath.Join(h.getContentPath("posts"), slug+".md.archive")
		if err := os.Rename(filePath, archPath); err != nil {
			slog.Error("archive post file failed", "slug", slug, "error", err)
		}
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

	archPath := filepath.Join(h.getContentPath("posts"), slug+".md.archive")
	if _, err := os.Stat(archPath); err == nil {
		filePath := filepath.Join(h.getContentPath("posts"), slug+".md")
		if err := os.Rename(archPath, filePath); err != nil {
			slog.Error("restore post file failed", "slug", slug, "error", err)
		}
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

	filePath := filepath.Join(h.getContentPath("posts"), slug+".md")
	archPath := filepath.Join(h.getContentPath("posts"), slug+".md.archive")
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		slog.Warn("delete post file failed", "path", filePath, "error", err)
	}
	if err := os.Remove(archPath); err != nil && !os.IsNotExist(err) {
		slog.Warn("delete archived post file failed", "path", archPath, "error", err)
	}

	imgDir := filepath.Join("web/static/images/posts", slug)
	_ = os.RemoveAll(imgDir)

	h.ReloadFilePosts()
	h.RefreshPosts()
	http.Redirect(w, r, "/admin#content", http.StatusSeeOther)
}
