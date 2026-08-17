package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"daemontalk/internal/comment"
	"daemontalk/internal/post"

	"github.com/go-chi/chi/v5"
)

func newTestCommentHandler(t *testing.T) (*Handler, *comment.Store) {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "comments.db")
	cs, err := comment.Open(dbPath)
	if err != nil {
		t.Fatalf("comment.Open: %v", err)
	}
	t.Cleanup(func() { cs.Close() })

	testPosts := []post.Post{
		{
			Slug:  "test-post",
			Title: "Test Article",
			Date:  time.Now(),
		},
	}

	h := &Handler{
		Comments:   cs,
		FilePosts:  testPosts,
		AdminToken: "admin-secret",
	}
	return h, cs
}

func TestPostCommentAndReply(t *testing.T) {
	h, cs := newTestCommentHandler(t)

	// 1. Post root comment
	form := url.Values{
		"body": {"Root discussion topic"},
	}
	req := httptest.NewRequest(http.MethodPost, "/blog/test-post/comments", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("slug", "test-post")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	h.PostComment(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Root discussion topic") {
		t.Error("expected response to contain root comment")
	}

	comments, err := cs.ListBySlug("test-post")
	if err != nil || len(comments) != 1 {
		t.Fatalf("expected 1 comment in store, got %d (err: %v)", len(comments), err)
	}
	rootID := comments[0].ID

	// 2. Post reply to root comment
	replyForm := url.Values{
		"body":      {"This is a thoughtful reply"},
		"parent_id": {fmt.Sprintf("%d", rootID)},
	}
	reqReply := httptest.NewRequest(http.MethodPost, "/blog/test-post/comments", strings.NewReader(replyForm.Encode()))
	reqReply.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rctxReply := chi.NewRouteContext()
	rctxReply.URLParams.Add("slug", "test-post")
	reqReply = reqReply.WithContext(context.WithValue(reqReply.Context(), chi.RouteCtxKey, rctxReply))

	recReply := httptest.NewRecorder()
	h.PostComment(recReply, reqReply)

	if recReply.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for reply, got %d", recReply.Code)
	}

	replyBody := recReply.Body.String()
	if !strings.Contains(replyBody, "This is a thoughtful reply") {
		t.Error("expected response to contain reply comment")
	}
	if !strings.Contains(replyBody, "Root discussion topic") {
		t.Error("expected response to contain parent comment")
	}

	commentsAfter, _ := cs.ListBySlug("test-post")
	if len(commentsAfter) != 2 {
		t.Fatalf("expected 2 comments in store, got %d", len(commentsAfter))
	}
	if commentsAfter[1].ParentID == nil || *commentsAfter[1].ParentID != rootID {
		t.Errorf("expected parent_id to be %d, got %v", rootID, commentsAfter[1].ParentID)
	}
}

func TestDeleteCommentWithReplies(t *testing.T) {
	h, cs := newTestCommentHandler(t)

	c1, _ := cs.Add("test-post", "Alice", "Root comment")
	_, _ = cs.AddWithParent("test-post", "Bob", "Reply to Alice", &c1.ID)

	// Delete as Admin
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/blog/test-post/comments/%d/delete", c1.ID), nil)
	req.AddCookie(&http.Cookie{Name: CookieAdminToken, Value: "admin-secret"})
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("slug", "test-post")
	rctx.URLParams.Add("id", fmt.Sprintf("%d", c1.ID))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	h.DeleteComment(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	remaining, _ := cs.ListBySlug("test-post")
	if len(remaining) != 0 {
		t.Errorf("expected 0 comments after deleting parent, got %d", len(remaining))
	}
}
