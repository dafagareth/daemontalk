package handler

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestHandler() *Handler {
	return &Handler{
		AdminToken: "test-admin-token",
	}
}

func TestAdminPostUploadMD_Unauthorized(t *testing.T) {
	h := setupTestHandler()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.md")
	_, _ = part.Write([]byte("# Test Title\n\nContent here"))
	_ = writer.Close()

	req := httptest.NewRequest("POST", "/admin/posts/upload-md", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	h.AdminPostUploadMD(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for unauthorized, got %d", w.Code)
	}
}

func TestAdminPostUploadMD_Success(t *testing.T) {
	h := setupTestHandler()

	rawMarkdown := `---
title: "Testing Direct Upload"
slug: "test-direct-upload"
date: 2026-08-16
tags: ["go", "systems"]
lang: "en"
draft: false
description: "A test markdown post"
---

Hello world from automated test!
`

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "upload.md")
	_, _ = part.Write([]byte(rawMarkdown))
	_ = writer.Close()

	req := httptest.NewRequest("POST", "/admin/posts/upload-md", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(&http.Cookie{Name: "admin_token", Value: "test-admin-token"})
	w := httptest.NewRecorder()

	h.AdminPostUploadMD(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 redirect, got %d, body: %s", w.Code, w.Body.String())
	}

	savedFile := filepath.Join("content/posts", "test-direct-upload.md")
	defer os.Remove(savedFile)

	data, err := os.ReadFile(savedFile)
	if err != nil {
		t.Fatalf("expected file to be written, got err: %v", err)
	}
	if !strings.Contains(string(data), "Testing Direct Upload") {
		t.Fatalf("file content does not match uploaded content: %s", string(data))
	}
}

func TestAdminPostFileEditAndSave(t *testing.T) {
	h := setupTestHandler()

	slug := "test-edit-save-cycle"
	filePath := filepath.Join("content/posts", slug+".md")
	initContent := `---
title: "Initial Title"
slug: "test-edit-save-cycle"
date: 2026-08-16
tags: ["systems"]
lang: "en"
draft: false
---

Initial body content.
`
	_ = os.WriteFile(filePath, []byte(initContent), 0644)
	defer os.Remove(filePath)

	// 1. Test Edit GET
	reqEdit := httptest.NewRequest("GET", "/admin/posts/file-edit?slug="+slug, nil)
	reqEdit.AddCookie(&http.Cookie{Name: "admin_token", Value: "test-admin-token"})
	wEdit := httptest.NewRecorder()
	h.AdminPostFileEdit(wEdit, reqEdit)

	if wEdit.Code != http.StatusOK {
		t.Fatalf("expected 200 on file edit, got %d", wEdit.Code)
	}
	if !strings.Contains(wEdit.Body.String(), "content/posts/test-edit-save-cycle.md") {
		t.Fatalf("editor view missing filename: %s", wEdit.Body.String())
	}

	// 2. Test Save POST
	updatedContent := `---
title: "Updated Title"
slug: "test-edit-save-cycle"
date: 2026-08-16
tags: ["systems", "linux"]
lang: "en"
draft: false
---

Updated body content via web editor.
`
	body := strings.NewReader("slug=" + slug + "&content=" + strings.ReplaceAll(updatedContent, "\n", "%0A") + "&action=save")
	reqSave := httptest.NewRequest("POST", "/admin/posts/file-save", body)
	reqSave.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqSave.AddCookie(&http.Cookie{Name: "admin_token", Value: "test-admin-token"})
	wSave := httptest.NewRecorder()
	h.AdminPostFileSave(wSave, reqSave)

	if wSave.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 on save, got %d", wSave.Code)
	}

	saved, _ := os.ReadFile(filePath)
	if !strings.Contains(string(saved), "Updated Title") {
		t.Fatalf("file not updated on disk: %s", string(saved))
	}
}

func TestAdminUploadImage(t *testing.T) {
	h := setupTestHandler()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("slug", "test-img-post")
	part, _ := writer.CreateFormFile("image", "architecture.png")
	_, _ = io.WriteString(part, "fake PNG binary data")
	_ = writer.Close()

	req := httptest.NewRequest("POST", "/admin/upload-image", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(&http.Cookie{Name: "admin_token", Value: "test-admin-token"})
	w := httptest.NewRecorder()

	h.AdminUploadImage(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on image upload, got %d: %s", w.Code, w.Body.String())
	}

	expectedDir := filepath.Join("web/static/images/posts", "test-img-post")
	defer os.RemoveAll(expectedDir)

	if !strings.Contains(w.Body.String(), "/static/images/posts/test-img-post/architecture.png") {
		t.Fatalf("unexpected JSON response: %s", w.Body.String())
	}
}

func TestAdminPostFileArchiveRestoreDelete(t *testing.T) {
	h := setupTestHandler()

	slug := "test-lifecycle-post"
	filePath := filepath.Join("content/posts", slug+".md")
	archPath := filepath.Join("content/posts", slug+".md.archive")
	initContent := "---\ntitle: \"Lifecycle Test\"\nslug: \"test-lifecycle-post\"\ndate: 2026-08-16\n---\nBody content."
	_ = os.WriteFile(filePath, []byte(initContent), 0644)
	defer os.Remove(filePath)
	defer os.Remove(archPath)

	// 1. Archive
	reqArch := httptest.NewRequest("POST", "/admin/posts/file-archive", strings.NewReader("slug="+slug))
	reqArch.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqArch.AddCookie(&http.Cookie{Name: "admin_token", Value: "test-admin-token"})
	wArch := httptest.NewRecorder()
	h.AdminPostFileArchive(wArch, reqArch)

	if wArch.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 on archive, got %d", wArch.Code)
	}
	if _, err := os.Stat(archPath); os.IsNotExist(err) {
		t.Fatalf("expected .md.archive to exist after archive")
	}

	// 2. Restore
	reqRestore := httptest.NewRequest("POST", "/admin/posts/file-restore", strings.NewReader("slug="+slug))
	reqRestore.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqRestore.AddCookie(&http.Cookie{Name: "admin_token", Value: "test-admin-token"})
	wRestore := httptest.NewRecorder()
	h.AdminPostFileRestore(wRestore, reqRestore)

	if wRestore.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 on restore, got %d", wRestore.Code)
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("expected .md to exist after restore")
	}

	// 3. Delete
	reqDelete := httptest.NewRequest("POST", "/admin/posts/file-delete", strings.NewReader("slug="+slug))
	reqDelete.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqDelete.AddCookie(&http.Cookie{Name: "admin_token", Value: "test-admin-token"})
	wDelete := httptest.NewRecorder()
	h.AdminPostFileDelete(wDelete, reqDelete)

	if wDelete.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 on delete, got %d", wDelete.Code)
	}
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Fatalf("expected file to be deleted")
	}
}
