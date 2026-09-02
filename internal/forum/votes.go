package forum

import (
	"encoding/json"
	"fmt"
	"time"
)

// Vote handles upvoting and toggling votes on topics and replies.
func (s *Store) Vote(userID int64, targetType string, targetID int64) (int, bool, error) {
	if userID <= 0 {
		return 0, false, fmt.Errorf("authentication required to vote")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return 0, false, err
	}
	defer tx.Rollback()

	var exists int
	_ = tx.QueryRow(`SELECT COUNT(1) FROM forum_votes WHERE user_id = ? AND target_type = ? AND target_id = ?`, userID, targetType, targetID).Scan(&exists)

	hasVoted := false
	if exists > 0 {
		// Remove vote (toggle off)
		if _, err := tx.Exec(`DELETE FROM forum_votes WHERE user_id = ? AND target_type = ? AND target_id = ?`, userID, targetType, targetID); err != nil {
			return 0, false, err
		}
		if targetType == "topic" {
			_, _ = tx.Exec(`UPDATE forum_topics SET votes_count = MAX(0, votes_count - 1) WHERE id = ?`, targetID)
		} else {
			_, _ = tx.Exec(`UPDATE forum_replies SET votes_count = MAX(0, votes_count - 1) WHERE id = ?`, targetID)
		}
		hasVoted = false
	} else {
		// Add vote
		if _, err := tx.Exec(`INSERT INTO forum_votes (user_id, target_type, target_id, created_at) VALUES (?, ?, ?, ?)`, userID, targetType, targetID, time.Now().UTC()); err != nil {
			return 0, false, err
		}
		if targetType == "topic" {
			_, _ = tx.Exec(`UPDATE forum_topics SET votes_count = votes_count + 1 WHERE id = ?`, targetID)
		} else {
			_, _ = tx.Exec(`UPDATE forum_replies SET votes_count = votes_count + 1 WHERE id = ?`, targetID)
		}
		hasVoted = true
	}

	var newCount int
	if targetType == "topic" {
		_ = tx.QueryRow(`SELECT votes_count FROM forum_topics WHERE id = ?`, targetID).Scan(&newCount)
	} else {
		_ = tx.QueryRow(`SELECT votes_count FROM forum_replies WHERE id = ?`, targetID).Scan(&newCount)
	}

	if err := tx.Commit(); err != nil {
		return 0, false, err
	}

	return newCount, hasVoted, nil
}

// GetUserStats retrieves aggregated forum statistics for a given username.
func (s *Store) GetUserStats(username string) (UserStats, error) {
	var stats UserStats
	query := `
		SELECT 
			COALESCE((SELECT COUNT(1) FROM forum_topics t JOIN users u ON u.id = t.user_id WHERE u.username = ? COLLATE NOCASE), 0),
			COALESCE((SELECT COUNT(1) FROM forum_replies r JOIN users u ON u.id = r.user_id WHERE u.username = ? COLLATE NOCASE), 0),
			COALESCE((SELECT COUNT(1) FROM forum_replies r JOIN users u ON u.id = r.user_id WHERE u.username = ? COLLATE NOCASE AND r.is_solution = 1), 0),
			COALESCE((SELECT SUM(t.votes_count) FROM forum_topics t JOIN users u ON u.id = t.user_id WHERE u.username = ? COLLATE NOCASE), 0) +
			COALESCE((SELECT SUM(r.votes_count) FROM forum_replies r JOIN users u ON u.id = r.user_id WHERE u.username = ? COLLATE NOCASE), 0)
	`
	err := s.db.QueryRow(query, username, username, username, username, username).Scan(
		&stats.TopicsCount,
		&stats.RepliesCount,
		&stats.SolutionsCount,
		&stats.TotalVotes,
	)
	return stats, err
}

// Stats contains aggregate numbers for forum activities.
type Stats struct {
	TotalTopics  int
	TotalReplies int
	SolvedTopics int
	TotalVotes   int
}

// GetStats returns aggregated metrics for the forum.
func (s *Store) GetStats() Stats {
	var st Stats
	_ = s.db.QueryRow(`SELECT COUNT(1) FROM forum_topics`).Scan(&st.TotalTopics)
	_ = s.db.QueryRow(`SELECT COUNT(1) FROM forum_replies`).Scan(&st.TotalReplies)
	_ = s.db.QueryRow(`SELECT COUNT(1) FROM forum_topics WHERE solved_reply_id IS NOT NULL`).Scan(&st.SolvedTopics)
	_ = s.db.QueryRow(`SELECT COUNT(1) FROM forum_votes`).Scan(&st.TotalVotes)
	return st
}

// UserContributions holds all topics and replies authored by a user.
type UserContributions struct {
	Topics  []*Topic `json:"topics"`
	Replies []*Reply `json:"replies"`
}

// GetUserContributions retrieves all topics and replies authored by a user.
func (s *Store) GetUserContributions(userID int64) (UserContributions, error) {
	var contrib UserContributions
	// Fetch topics
	tRows, err := s.db.Query(`SELECT id, title, slug, category, tags, body_md, created_at FROM forum_topics WHERE user_id = ? ORDER BY created_at DESC`, userID)
	if err == nil {
		defer tRows.Close()
		for tRows.Next() {
			var t Topic
			var tagsJSON string
			if err := tRows.Scan(&t.ID, &t.Title, &t.Slug, &t.Category, &tagsJSON, &t.BodyMD, &t.CreatedAt); err == nil {
				_ = json.Unmarshal([]byte(tagsJSON), &t.Tags)
				contrib.Topics = append(contrib.Topics, &t)
			}
		}
	}

	// Fetch replies
	rRows, err := s.db.Query(`SELECT id, topic_id, body_md, is_solution, created_at FROM forum_replies WHERE user_id = ? ORDER BY created_at DESC`, userID)
	if err == nil {
		defer rRows.Close()
		for rRows.Next() {
			var r Reply
			if err := rRows.Scan(&r.ID, &r.TopicID, &r.BodyMD, &r.IsSolution, &r.CreatedAt); err == nil {
				contrib.Replies = append(contrib.Replies, &r)
			}
		}
	}

	return contrib, nil
}
