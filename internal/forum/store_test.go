package forum

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestForumStore_CRUDAndInteractions(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_forum.db")

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open forum db: %v", err)
	}
	defer store.Close()

	// 1. Create Topic
	t1, err := store.CreateTopic(Topic{
		UserID:   42,
		Title:    "Mengapa Linux io_uring Lebih Cepat dari Epoll?",
		Category: "kernel",
		Tags:     []string{"linux", "kernel", "io_uring"},
		BodyMD:   "Ini adalah pertanyaan tentang arsitektur `io_uring` pada kernel Linux modern.",
	})
	if err != nil {
		t.Fatalf("failed to create topic: %v", err)
	}
	if t1.ID <= 0 {
		t.Fatalf("expected positive topic ID, got %d", t1.ID)
	}
	if t1.Slug == "" {
		t.Fatalf("expected non-empty slug")
	}

	// 2. Fetch Topic
	fetched, err := store.GetTopicBySlug(t1.Slug, 42)
	if err != nil || fetched == nil {
		t.Fatalf("failed to get topic by slug: %v", err)
	}
	if fetched.Title != t1.Title {
		t.Errorf("expected title %s, got %s", t1.Title, fetched.Title)
	}
	if len(fetched.Tags) != 3 {
		t.Errorf("expected 3 tags, got %d", len(fetched.Tags))
	}
	if !fetched.IsOwner {
		t.Errorf("expected is_owner to be true for user 42")
	}

	// 3. Add Replies
	r1, err := store.CreateReply(Reply{
		TopicID: t1.ID,
		UserID:  99,
		BodyMD:  "Karena io_uring menggunakan dua lock-free ring buffer (SQ & CQ) yang di-share dengan kernel.",
	})
	if err != nil {
		t.Fatalf("failed to create reply 1: %v", err)
	}

	// Nested reply to r1
	r2, err := store.CreateReply(Reply{
		TopicID:  t1.ID,
		ParentID: r1.ID,
		UserID:   42,
		BodyMD:   "Terima kasih atas penjelasannya! Sangat jelas.",
	})
	if err != nil {
		t.Fatalf("failed to create reply 2: %v", err)
	}

	replies, err := store.GetTopicReplies(t1.ID, 42)
	if err != nil {
		t.Fatalf("failed to get replies: %v", err)
	}
	if len(replies) != 1 {
		t.Fatalf("expected 1 root reply, got %d", len(replies))
	}
	if len(replies[0].Children) != 1 || replies[0].Children[0].ID != r2.ID {
		t.Fatalf("expected 1 child reply under r1")
	}

	// 4. Mark Accepted Solution
	if err := store.MarkSolution(t1.ID, r1.ID, 42); err != nil {
		t.Fatalf("failed to mark solution: %v", err)
	}
	updatedTopic, _ := store.GetTopicBySlug(t1.Slug, 42)
	if !updatedTopic.IsSolved || updatedTopic.SolvedReplyID != r1.ID {
		t.Errorf("expected topic to be solved with reply id %d", r1.ID)
	}

	// 5. Upvote Voting
	votes, hasVoted, err := store.Vote(42, "topic", t1.ID)
	if err != nil {
		t.Fatalf("failed to vote topic: %v", err)
	}
	if votes != 1 || !hasVoted {
		t.Errorf("expected 1 vote and hasVoted=true, got votes=%d, hasVoted=%v", votes, hasVoted)
	}

	// Toggle vote off
	votes, hasVoted, err = store.Vote(42, "topic", t1.ID)
	if err != nil {
		t.Fatalf("failed to toggle vote: %v", err)
	}
	if votes != 0 || hasVoted {
		t.Errorf("expected 0 votes and hasVoted=false, got votes=%d, hasVoted=%v", votes, hasVoted)
	}

	// 6. List Topics Filter
	list, total, err := store.ListTopics("kernel", "", "", "", "latest", 10, 0, 42)
	if err != nil {
		t.Fatalf("failed to list topics: %v", err)
	}
	if total != 1 || len(list) != 1 {
		t.Errorf("expected 1 topic in list, got %d (total: %d)", len(list), total)
	}

	// 7. Delete Reply & Delete Topic
	if err := store.DeleteReply(r2.ID, 42, false); err != nil {
		t.Fatalf("failed to delete reply 2: %v", err)
	}

	if err := store.DeleteTopic(t1.ID, 42, false); err != nil {
		t.Fatalf("failed to delete topic: %v", err)
	}

	deletedTopic, _ := store.GetTopicBySlug(t1.Slug, 42)
	if deletedTopic != nil {
		t.Errorf("expected nil topic after deletion")
	}
}

func TestRecordTopicViewDeduplication(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_views.db")

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open forum db: %v", err)
	}
	defer store.Close()

	topic, err := store.CreateTopic(Topic{
		UserID:   1,
		Title:    "Test Views Topic",
		Category: "general",
		Tags:     []string{"test"},
		BodyMD:   "Testing view deduplication.",
	})
	if err != nil {
		t.Fatalf("failed to create topic: %v", err)
	}

	// Initial views should be 0
	fetched, _ := store.GetTopicBySlug(topic.Slug, 0)
	if fetched.ViewsCount != 0 {
		t.Fatalf("expected 0 initial views, got %d", fetched.ViewsCount)
	}

	// First view from user 100
	rec, err := store.RecordTopicView(topic.ID, "u:100")
	if err != nil || !rec {
		t.Fatalf("expected recorded=true, got %v, err: %v", rec, err)
	}
	fetched, _ = store.GetTopicBySlug(topic.Slug, 0)
	if fetched.ViewsCount != 1 {
		t.Fatalf("expected 1 view, got %d", fetched.ViewsCount)
	}

	// Refresh 10 times by user 100 -> should NOT increment!
	for i := 0; i < 10; i++ {
		rec, err := store.RecordTopicView(topic.ID, "u:100")
		if err != nil || rec {
			t.Fatalf("expected recorded=false on refresh, got %v", rec)
		}
	}
	fetched, _ = store.GetTopicBySlug(topic.Slug, 0)
	if fetched.ViewsCount != 1 {
		t.Fatalf("expected views count to stay 1 after refreshes, got %d", fetched.ViewsCount)
	}

	// Different viewer: visitor anon_xyz
	rec, err = store.RecordTopicView(topic.ID, "v:anon_xyz")
	if err != nil || !rec {
		t.Fatalf("expected recorded=true for new visitor, got %v", rec)
	}
	fetched, _ = store.GetTopicBySlug(topic.Slug, 0)
	if fetched.ViewsCount != 2 {
		t.Fatalf("expected views count 2, got %d", fetched.ViewsCount)
	}

	// Refresh by visitor anon_xyz
	rec, _ = store.RecordTopicView(topic.ID, "v:anon_xyz")
	if rec {
		t.Fatalf("expected recorded=false on visitor refresh")
	}
	fetched, _ = store.GetTopicBySlug(topic.Slug, 0)
	if fetched.ViewsCount != 2 {
		t.Fatalf("expected views count to stay 2, got %d", fetched.ViewsCount)
	}
}

func TestRenderMarkdown_XSSPrevention(t *testing.T) {
	cleanMD := `# Valid Title

**Bold Text**`
	htmlClean := string(RenderMarkdown(cleanMD))
	if !strings.Contains(htmlClean, "Valid Title") {
		t.Errorf("expected Valid Title to be present in %s", htmlClean)
	}
	if !strings.Contains(htmlClean, "<strong>Bold Text</strong>") {
		t.Errorf("expected <strong>Bold Text</strong> in %s", htmlClean)
	}

	dangerousScript := `<script>alert('xss')</script>`
	htmlScript := string(RenderMarkdown(dangerousScript))
	if strings.Contains(htmlScript, "<script") || strings.Contains(htmlScript, "alert('xss')") {
		t.Errorf("script tag or payload was not stripped: %s", htmlScript)
	}

	dangerousImg := `<img src="x" onerror="alert(1)"/>`
	htmlImg := string(RenderMarkdown(dangerousImg))
	if strings.Contains(htmlImg, "onerror") {
		t.Errorf("onerror event handler was not stripped: %s", htmlImg)
	}

	dangerousLink := `[Click here](javascript:alert(1))`
	htmlLink := string(RenderMarkdown(dangerousLink))
	if strings.Contains(htmlLink, "javascript:") {
		t.Errorf("javascript: URI was not stripped: %s", htmlLink)
	}
}
