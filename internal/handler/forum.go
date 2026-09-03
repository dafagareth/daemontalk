package handler

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"daemontalk/internal/auth"
	"daemontalk/internal/forum"
	"daemontalk/internal/i18n"
	"daemontalk/web/templates"
	"github.com/go-chi/chi/v5"
)

// Discussions lists forum topics with categories, search, and sorting.
func (h *Handler) Discussions(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)
	user := auth.GetUser(r.Context())

	var currentUserID int64
	if user != nil {
		currentUserID = user.ID
	}

	category := r.URL.Query().Get("category")
	tag := r.URL.Query().Get("tag")
	author := strings.TrimSpace(r.URL.Query().Get("author"))
	search := strings.TrimSpace(r.URL.Query().Get("q"))
	sortOrder := r.URL.Query().Get("sort")
	if sortOrder == "" {
		sortOrder = "latest"
	}

	page := 1
	if p, err := strconv.Atoi(r.URL.Query().Get("p")); err == nil && p > 0 {
		page = p
	}
	limit := 20
	offset := (page - 1) * limit

	var topics []*forum.Topic
	var total int
	var err error

	if h.Forum != nil {
		topics, total, err = h.Forum.ListTopics(category, tag, search, author, sortOrder, limit, offset, currentUserID)
		if err != nil {
			slog.Error("failed to list forum topics", "error", err)
		}
	}

	totalPages := (total + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}

	title := "Discussions & Q&A"
	if lang == "id" {
		title = "Forum Diskusi & Tanya Jawab"
	}

	meta := templates.PageMeta{
		Description: "Daemontalk open tech and systems discussions, debugging Q&A, and incident post-mortems.",
	}

	h.Render(w, r, templates.Layout(ui, lang, title, r.URL.Path, meta,
		templates.DiscussionsIndexPage(ui, lang, user, topics, tag, sortOrder, search, author, total, page, totalPages),
	))
}

// DiscussionsNew shows the form to start a new discussion or ask a question.
func (h *Handler) DiscussionsNew(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)
	user := auth.GetUser(r.Context())

	title := "New Discussion"
	if lang == "id" {
		title = "Buat Topik Baru"
	}

	h.Render(w, r, templates.Layout(ui, lang, title, r.URL.Path, templates.PageMeta{
		Description: "Start a new tech discussion or ask a question on Daemontalk.",
	}, templates.DiscussionsNewPage(ui, lang, user)))
}

// DiscussionsCreate handles the submission of a new topic.
func (h *Handler) DiscussionsCreate(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		http.Redirect(w, r, "/auth/github?return_to=/socket/new", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	category := strings.TrimSpace(r.FormValue("category"))
	bodyMD := strings.TrimSpace(r.FormValue("body"))
	tagsRaw := strings.TrimSpace(r.FormValue("tags"))

	if title == "" || bodyMD == "" {
		http.Error(w, "Title and content cannot be empty", http.StatusBadRequest)
		return
	}

	if category == "" {
		category = "general"
	}

	var tags []string
	if tagsRaw != "" {
		for _, t := range strings.Split(tagsRaw, ",") {
			cleaned := strings.ToLower(strings.TrimSpace(t))
			if cleaned != "" {
				tags = append(tags, cleaned)
			}
		}
	}

	topic, err := h.Forum.CreateTopic(forum.Topic{
		UserID:   user.ID,
		Title:    title,
		Category: category,
		Tags:     tags,
		BodyMD:   bodyMD,
	})
	if err != nil {
		slog.Error("failed to create forum topic", "error", err)
		http.Error(w, "Failed to create topic", http.StatusInternalServerError)
		return
	}

	lang := langFromRequest(r)
	prefix := urlPrefix(lang)
	http.Redirect(w, r, prefix+"/socket/"+topic.Slug, http.StatusSeeOther)
}

// DiscussionsDetail views a single topic and its replies.
func (h *Handler) DiscussionsDetail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.NotFound(w, r)
		return
	}

	lang := langFromRequest(r)
	ui := i18n.Get(lang)
	user := auth.GetUser(r.Context())

	var currentUserID int64
	if user != nil {
		currentUserID = user.ID
	}

	topic, err := h.Forum.GetTopicBySlug(slug, currentUserID)
	if err != nil || topic == nil {
		http.NotFound(w, r)
		return
	}

	// Record view only for unique human viewers (exclude bots, CLI, author, and admins)
	if !IsCLIRequest(r) && !isBot(r) && !topic.IsOwner && !h.isAdmin(r) {
		viewerKey := GetViewerKey(w, r, user)
		if recorded, _ := h.Forum.RecordTopicView(topic.ID, viewerKey); recorded {
			topic.ViewsCount++
		}
	}

	replies, err := h.Forum.GetTopicReplies(topic.ID, currentUserID)
	if err != nil {
		slog.Warn("failed to load topic replies", "topic_id", topic.ID, "error", err)
	}

	meta := templates.PageMeta{
		Description: topic.Title,
	}

	h.Render(w, r, templates.Layout(ui, lang, topic.Title+" · Discussions", r.URL.Path, meta,
		templates.DiscussionsDetailPage(ui, lang, user, topic, replies),
	))
}

// DiscussionsReply submits a reply to a discussion.
func (h *Handler) DiscussionsReply(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	topicIDStr := chi.URLParam(r, "id")
	topicID, err := strconv.ParseInt(topicIDStr, 10, 64)
	if err != nil || topicID <= 0 {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	bodyMD := strings.TrimSpace(r.FormValue("body"))
	if bodyMD == "" {
		http.Error(w, "Reply content cannot be empty", http.StatusBadRequest)
		return
	}

	var parentID int64
	if pStr := r.FormValue("parent_id"); pStr != "" {
		if p, err := strconv.ParseInt(pStr, 10, 64); err == nil {
			parentID = p
		}
	}

	reply, err := h.Forum.CreateReply(forum.Reply{
		TopicID:  topicID,
		ParentID: parentID,
		UserID:   user.ID,
		BodyMD:   bodyMD,
	})
	if err != nil {
		slog.Error("failed to create reply", "error", err)
		http.Error(w, "Failed to submit reply", http.StatusInternalServerError)
		return
	}

	// Populate author fields for immediate UI render
	reply.AuthorName = user.DisplayName
	reply.AuthorUsername = user.Username
	reply.AuthorAvatar = user.AvatarURL
	reply.AuthorGitHub = user.GitHubURL
	reply.IsOwner = true

	// If request is from HTMX, render just the reply card
	if r.Header.Get("HX-Request") == "true" {
		lang := langFromRequest(r)
		ui := i18n.Get(lang)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = templates.DiscussionsReplyItem(ui, lang, user, reply, topicID, user.ID).Render(r.Context(), w)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

// DiscussionsSolve marks or unmarks a reply as the accepted solution.
func (h *Handler) DiscussionsSolve(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	topicIDStr := chi.URLParam(r, "id")
	topicID, err := strconv.ParseInt(topicIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	replyIDStr := r.URL.Query().Get("reply_id")
	replyID, _ := strconv.ParseInt(replyIDStr, 10, 64)

	if err := h.Forum.MarkSolution(topicID, replyID, user.ID); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Refresh", "true")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

// DiscussionsVote handles upvoting topics or replies via HTMX.
func (h *Handler) DiscussionsVote(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Please login with GitHub to vote", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	targetType := r.FormValue("type") // "topic" or "reply"
	targetID, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err != nil || (targetType != "topic" && targetType != "reply") {
		http.Error(w, "Invalid vote target", http.StatusBadRequest)
		return
	}

	newCount, hasVoted, err := h.Forum.Vote(user.ID, targetType, targetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render HTMX vote button response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = templates.DiscussionsVoteButton(targetType, targetID, newCount, hasVoted).Render(r.Context(), w)
}

// DiscussionsDeleteTopic deletes a topic (owner or admin).
func (h *Handler) DiscussionsDeleteTopic(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	isAdmin := h.isAdmin(r) || user.Role == "admin"
	if err := h.Forum.DeleteTopic(id, user.ID, isAdmin); err != nil {
		slog.Error("delete forum topic failed", "id", id, "error", err)
		http.Error(w, "Failed to delete topic", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/socket")
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, "/socket", http.StatusSeeOther)
}

// DiscussionsDeleteReply deletes or redacts a reply (author or admin).
func (h *Handler) DiscussionsDeleteReply(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid reply ID", http.StatusBadRequest)
		return
	}

	isAdmin := h.isAdmin(r) || user.Role == "admin"
	if err := h.Forum.DeleteReply(id, user.ID, isAdmin); err != nil {
		slog.Error("delete forum reply failed", "id", id, "error", err)
		http.Error(w, "Failed to delete reply", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Refresh", "true")
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}
