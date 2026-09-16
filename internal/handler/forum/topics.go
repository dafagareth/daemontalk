package forum

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"daemontalk/internal/auth"
	"daemontalk/internal/forum"
	"daemontalk/internal/handler/common"
	"daemontalk/internal/i18n"
	forumtmpl "daemontalk/web/templates/forum"
	"daemontalk/web/templates/layouts"
	"daemontalk/web/templates/shared"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) Discussions(w http.ResponseWriter, r *http.Request) {
	lang := common.LangFromRequest(r)
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

	var activeTopic *forum.Topic
	var activeReplies []*forum.Reply

	if len(topics) > 0 && h.Forum != nil {
		activeSlug := r.URL.Query().Get("slug")
		if activeSlug == "" {
			activeSlug = topics[0].Slug
		}
		activeTopic, _ = h.Forum.GetTopicBySlug(activeSlug, currentUserID)
		if activeTopic != nil {
			activeReplies, _ = h.Forum.GetTopicReplies(activeTopic.ID, currentUserID)
		}
	}

	title := "Socket · Discussions & Q&A"

	meta := shared.PageMeta{
		Description: "Daemontalk open tech and systems discussions, debugging Q&A, and incident post-mortems.",
	}

	common.Render(w, r, layouts.Layout(ui, lang, title, r.URL.Path, meta,
		forumtmpl.IndexPage(ui, lang, user, topics, tag, sortOrder, search, author, total, page, totalPages, activeTopic, activeReplies, false),
	))
}

func (h *Handler) DiscussionsNew(w http.ResponseWriter, r *http.Request) {
	lang := common.LangFromRequest(r)
	ui := i18n.Get(lang)
	user := auth.GetUser(r.Context())

	title := "New Discussion"

	common.Render(w, r, layouts.Layout(ui, lang, title, r.URL.Path, shared.PageMeta{
		Description: "Start a new tech discussion or ask a question on Daemontalk.",
	}, forumtmpl.NewPage(ui, lang, user)))
}

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

	lang := common.LangFromRequest(r)
	prefix := common.URLPrefix(lang)
	http.Redirect(w, r, prefix+"/socket/"+topic.Slug, http.StatusSeeOther)
}

func (h *Handler) DiscussionsDetail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.NotFound(w, r)
		return
	}

	lang := common.LangFromRequest(r)
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

	if !common.IsBot(r) && !topic.IsOwner && !h.IsAdmin(r) {
		viewerKey := common.GetViewerKey(w, r, user)
		if recorded, _ := h.Forum.RecordTopicView(topic.ID, viewerKey); recorded {
			topic.ViewsCount++
		}
	}

	replies, err := h.Forum.GetTopicReplies(topic.ID, currentUserID)
	if err != nil {
		slog.Warn("failed to load topic replies", "topic_id", topic.ID, "error", err)
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = forumtmpl.TopicPane(ui, lang, user, topic, replies).Render(r.Context(), w)
		return
	}

	var topics []*forum.Topic
	var total int
	if h.Forum != nil {
		topics, total, _ = h.Forum.ListTopics("", "", "", "", "latest", 25, 0, currentUserID)
	}
	totalPages := (total + 25 - 1) / 25
	if totalPages < 1 {
		totalPages = 1
	}

	meta := shared.PageMeta{
		Description: topic.Title,
	}

	title := topic.Title + " · Socket"

	common.Render(w, r, layouts.Layout(ui, lang, title, r.URL.Path, meta,
		forumtmpl.IndexPage(ui, lang, user, topics, "", "latest", "", "", total, 1, totalPages, topic, replies, true),
	))
}

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

	isAdmin := h.IsAdmin(r) || user.Role == "admin"
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
