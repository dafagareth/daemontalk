package handler

import (
	"net/http"

	"daemontalk/internal/handler/admin"
)

func (h *Handler) AdminHandler() *admin.Handler {
	return &admin.Handler{
		AdminToken:      h.AdminToken,
		Comments:        h.Comments,
		ContentDir:      h.ContentDir,
		FilePosts:       h.GetFilePosts,
		ReloadFilePosts: h.ReloadFilePosts,
		VisiblePosts:    h.VisiblePosts,
		AllPosts:        h.AllPosts,
		RefreshPosts:    h.RefreshPosts,
		NotFound:        h.NotFound,
	}
}

func (h *Handler) Admin(w http.ResponseWriter, r *http.Request) {
	h.AdminHandler().Admin(w, r)
}

func (h *Handler) AdminDeleteComment(w http.ResponseWriter, r *http.Request) {
	h.AdminHandler().AdminDeleteComment(w, r)
}

func (h *Handler) AdminPostUploadMD(w http.ResponseWriter, r *http.Request) {
	h.AdminHandler().AdminPostUploadMD(w, r)
}

func (h *Handler) AdminUploadImage(w http.ResponseWriter, r *http.Request) {
	h.AdminHandler().AdminUploadImage(w, r)
}

func (h *Handler) AdminPostFileEdit(w http.ResponseWriter, r *http.Request) {
	h.AdminHandler().AdminPostFileEdit(w, r)
}

func (h *Handler) AdminPostFileSave(w http.ResponseWriter, r *http.Request) {
	h.AdminHandler().AdminPostFileSave(w, r)
}

func (h *Handler) AdminPostExportMD(w http.ResponseWriter, r *http.Request) {
	h.AdminHandler().AdminPostExportMD(w, r)
}

func (h *Handler) AdminPostFileArchive(w http.ResponseWriter, r *http.Request) {
	h.AdminHandler().AdminPostFileArchive(w, r)
}

func (h *Handler) AdminPostFileRestore(w http.ResponseWriter, r *http.Request) {
	h.AdminHandler().AdminPostFileRestore(w, r)
}

func (h *Handler) AdminPostFileDelete(w http.ResponseWriter, r *http.Request) {
	h.AdminHandler().AdminPostFileDelete(w, r)
}
