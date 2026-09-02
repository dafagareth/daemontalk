package forum

import (
	"path/filepath"
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
