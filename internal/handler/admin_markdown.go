package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"daemontalk/internal/post"
	"daemontalk/web/templates"
)

var safeSlugRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// cleanSlug ensures a slug string is strictly safe for file naming.
func cleanSlug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	var out []rune
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			out = append(out, r)
		}
	}
	return string(out)
}

// AdminPostUploadMD handles bulk / single .md file uploads directly to content/posts.
func (h *Handler) AdminPostUploadMD(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		h.NotFound(w, r)
		return
	}

	// Max upload size
	if err := r.ParseMultipartForm(MaxMarkdownUploadSize); err != nil {
		http.Error(w, "File upload size exceeded", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		files = r.MultipartForm.File["file"]
	}
	if len(files) == 0 {
		http.Error(w, "No markdown files selected", http.StatusBadRequest)
		return
	}

	uploadedCount := 0
	var lastSlug string

	for _, fileHeader := range files {
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if ext != ".md" && ext != ".markdown" {
			continue
		}

		file, err := fileHeader.Open()
		if err != nil {
			slog.Error("open uploaded markdown file failed", "filename", fileHeader.Filename, "error", err)
			continue
		}

		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, io.LimitReader(file, 5<<20)); err != nil {
			_ = file.Close()
			continue
		}
		_ = file.Close()

		rawBytes := buf.Bytes()
		p, err := post.Parse(rawBytes)
		if err != nil {
			slog.Error("parse uploaded markdown failed", "filename", fileHeader.Filename, "error", err)
			continue
		}

		slug := cleanSlug(p.Slug)
		if slug == "" {
			if p.Title != "" {
				slug = cleanSlug(slugify(p.Title))
			} else {
				slug = cleanSlug(strings.TrimSuffix(fileHeader.Filename, ext))
			}
		}
		if slug == "" {
			slug = generateShortID()
		}

		// Ensure content/posts directory exists
		if err := os.MkdirAll("content/posts", 0755); err != nil {
			slog.Error("mkdir content/posts failed", "error", err)
			http.Error(w, "Failed to create directory", http.StatusInternalServerError)
			return
		}

		targetPath := filepath.Join("content/posts", slug+".md")
		if err := os.WriteFile(targetPath, rawBytes, 0644); err != nil {
			slog.Error("write uploaded markdown file failed", "path", targetPath, "error", err)
			http.Error(w, "Failed to save file to disk", http.StatusInternalServerError)
			return
		}

		// Create associated post images directory if missing
		imgDir := filepath.Join("web/static/images/posts", slug)
		_ = os.MkdirAll(imgDir, 0755)

		uploadedCount++
		lastSlug = slug
	}

	if uploadedCount == 0 {
		http.Error(w, "No valid .md files could be parsed and uploaded", http.StatusBadRequest)
		return
	}

	// Reload all markdown posts and refresh the DB snapshot
	h.ReloadFilePosts()
	h.RefreshPosts()

	if uploadedCount == 1 && lastSlug != "" {
		http.Redirect(w, r, "/admin/posts/file-edit?slug="+lastSlug, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin#content", http.StatusSeeOther)
}

// AdminUploadImage handles image uploads directly from the Markdown Editor UI.
func (h *Handler) AdminUploadImage(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		h.NotFound(w, r)
		return
	}

	// Max image size
	if err := r.ParseMultipartForm(MaxImageUploadSize); err != nil {
		http.Error(w, "Image size exceeds limit", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "No image file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" && ext != ".svg" {
		http.Error(w, "Unsupported image type. Allowed: jpg, png, gif, webp, svg", http.StatusBadRequest)
		return
	}

	slug := cleanSlug(r.FormValue("slug"))
	if slug == "" {
		slug = "general"
	}

	destDir := filepath.Join("web/static/images/posts", slug)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		slog.Error("mkdir image dest failed", "dir", destDir, "error", err)
		http.Error(w, "Failed to create directory", http.StatusInternalServerError)
		return
	}

	// Clean image filename
	baseName := strings.TrimSuffix(filepath.Base(header.Filename), ext)
	cleanName := cleanSlug(baseName)
	if cleanName == "" {
		cleanName = fmt.Sprintf("img-%x", time.Now().UnixNano())
	}
	finalFilename := cleanName + ext
	destPath := filepath.Join(destDir, finalFilename)

	out, err := os.Create(destPath)
	if err != nil {
		slog.Error("create image file failed", "path", destPath, "error", err)
		http.Error(w, "Failed to write image file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		slog.Error("copy image data failed", "error", err)
		http.Error(w, "Failed to stream image content", http.StatusInternalServerError)
		return
	}

	relURL := fmt.Sprintf("/static/images/posts/%s/%s", slug, finalFilename)
	markdownCode := fmt.Sprintf("![%s](%s)", cleanName, relURL)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success":  true,
		"url":      relURL,
		"markdown": markdownCode,
		"filename": finalFilename,
	})
}

// AdminPostFileEdit displays the raw Markdown file editor for an existing repository post.
func (h *Handler) AdminPostFileEdit(w http.ResponseWriter, r *http.Request) {
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
	isArchived := false
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		filePath = filepath.Join("content/posts", slug+".md.archive")
		isArchived = true
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	p, _ := post.Parse(data)

	h.Render(w, r, templates.AdminLayout("admin", r.URL.Path,
		templates.AdminMarkdownEditor(slug, string(data), p, isArchived, "")))
}

// AdminPostFileSave saves edited Markdown content directly back to the disk file.
func (h *Handler) AdminPostFileSave(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		h.NotFound(w, r)
		return
	}

	slug := cleanSlug(r.FormValue("slug"))
	if slug == "" || !safeSlugRegex.MatchString(slug) {
		http.Error(w, "Invalid slug parameter", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if strings.TrimSpace(content) == "" {
		http.Error(w, "Content cannot be empty", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join("content/posts", slug+".md")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		archPath := filepath.Join("content/posts", slug+".md.archive")
		if _, errArch := os.Stat(archPath); errArch == nil {
			filePath = archPath
		}
	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		slog.Error("save markdown file failed", "path", filePath, "error", err)
		http.Error(w, "Failed to save markdown file to disk", http.StatusInternalServerError)
		return
	}

	h.ReloadFilePosts()
	h.RefreshPosts()

	action := r.FormValue("action")
	if action == "view" {
		http.Redirect(w, r, "/blog/"+slug, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin#content", http.StatusSeeOther)
}

// AdminPostExportMD downloads the raw .md file directly from the browser.

// AdminPostFileArchive renames a markdown file to .md.archive to hide it from public feed.

// AdminPostFileRestore renames a .md.archive file back to .md.

// AdminPostFileDelete permanently deletes a markdown file and its associated images directory.
