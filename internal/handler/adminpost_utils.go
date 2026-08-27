package handler

import (
	"bytes"
	"crypto/rand"
	"daemontalk/internal/postdb"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

func generateShortID() string {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%x", time.Now().UnixNano())[:8]
	}
	return hex.EncodeToString(b)
}

func (h *Handler) uniqueShortID(excludeID int64) string {
	for {
		id := generateShortID()
		if !h.slugTaken(id, excludeID) {
			return id
		}
	}
}

func slugify(title string) string {
	s := strings.ToLower(strings.TrimSpace(title))
	s = strings.ReplaceAll(s, " ", "-")
	s = slugCleanRe.ReplaceAllString(s, "")
	s = strings.Trim(strings.ReplaceAll(s, "--", "-"), "-")
	return s
}

func mdToEditorHTML(md string) string {
	var buf bytes.Buffer
	if err := editorMD.Convert([]byte(md), &buf); err != nil {
		slog.Error("editor markdown convert failed", "error", err)
		return ""
	}
	return buf.String()
}

func (h *Handler) slugTaken(slug string, excludeID int64) bool {
	h.filePostsMu.RLock()
	for _, fp := range h.FilePosts {
		if fp.Slug == slug {
			h.filePostsMu.RUnlock()
			return true
		}
	}
	h.filePostsMu.RUnlock()

	if h.PostDB != nil {
		if existing, err := h.PostDB.GetBySlug(slug); err == nil && existing.ID != excludeID {
			return true
		}
	}
	return false
}

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
