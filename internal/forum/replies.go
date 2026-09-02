package forum

import (
	"database/sql"
	"fmt"
	"html/template"
	"time"
)

// CreateReply inserts a new reply to a topic and increments replies_count.
func (s *Store) CreateReply(r Reply) (*Reply, error) {
	now := time.Now().UTC()
	bodyHTML := string(RenderMarkdown(r.BodyMD))

	var parentID sql.NullInt64
	if r.ParentID > 0 {
		parentID.Int64 = r.ParentID
		parentID.Valid = true
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO forum_replies (topic_id, parent_id, user_id, body_md, body_html, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		RETURNING id, created_at, updated_at
	`
	var id int64
	var createdAt, updatedAt time.Time
	if err := tx.QueryRow(query, r.TopicID, parentID, r.UserID, r.BodyMD, bodyHTML, now, now).Scan(&id, &createdAt, &updatedAt); err != nil {
		return nil, fmt.Errorf("insert reply: %w", err)
	}

	if _, err := tx.Exec(`UPDATE forum_topics SET replies_count = replies_count + 1, updated_at = ? WHERE id = ?`, now, r.TopicID); err != nil {
		return nil, fmt.Errorf("update topic replies count: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	r.ID = id
	r.BodyHTML = template.HTML(bodyHTML)
	r.CreatedAt = createdAt
	r.UpdatedAt = updatedAt
	return &r, nil
}

// GetTopicReplies returns all replies for a topic organized into a tree.
func (s *Store) GetTopicReplies(topicID int64, currentUserID int64) ([]*Reply, error) {
	query := `
		SELECT r.id, r.topic_id, r.parent_id, r.user_id, u.display_name, u.username, u.avatar_url, u.github_url,
		       r.body_md, r.body_html, r.is_solution, r.votes_count, r.created_at, r.updated_at,
		       EXISTS(SELECT 1 FROM forum_votes WHERE user_id = ? AND target_type = 'reply' AND target_id = r.id) AS user_voted
		FROM forum_replies r
		LEFT JOIN users u ON u.id = r.user_id
		WHERE r.topic_id = ?
		ORDER BY r.is_solution DESC, r.votes_count DESC, r.created_at ASC
	`
	rows, err := s.db.Query(query, currentUserID, topicID)
	if err != nil {
		return nil, fmt.Errorf("get replies: %w", err)
	}
	defer rows.Close()

	var allReplies []*Reply
	replyMap := make(map[int64]*Reply)

	for rows.Next() {
		var r Reply
		var parentID sql.NullInt64
		var userVoted int
		var authorName, authorUser, authorAvatar, authorGH sql.NullString

		if err := rows.Scan(
			&r.ID, &r.TopicID, &parentID, &r.UserID, &authorName, &authorUser, &authorAvatar, &authorGH,
			&r.BodyMD, &r.BodyHTML, &r.IsSolution, &r.VotesCount, &r.CreatedAt, &r.UpdatedAt, &userVoted,
		); err != nil {
			continue
		}

		r.AuthorName = authorName.String
		r.AuthorUsername = authorUser.String
		r.AuthorAvatar = authorAvatar.String
		r.AuthorGitHub = authorGH.String
		if r.AuthorName == "" {
			r.AuthorName = "Deleted User"
		}
		if r.AuthorUsername == "" {
			r.AuthorUsername = "ghost"
		}
		if r.AuthorAvatar == "" {
			r.AuthorAvatar = "/static/images/deleted-user.png"
		}
		if parentID.Valid {
			r.ParentID = parentID.Int64
		}
		r.UserVoted = (userVoted == 1)
		r.IsOwner = (currentUserID > 0 && r.UserID == currentUserID)

		replyMap[r.ID] = &r
		allReplies = append(allReplies, &r)
	}

	// Organize into top-level and children
	var rootReplies []*Reply
	for _, r := range allReplies {
		if r.ParentID > 0 {
			if parent, ok := replyMap[r.ParentID]; ok {
				parent.Children = append(parent.Children, r)
				continue
			}
		}
		rootReplies = append(rootReplies, r)
	}

	return rootReplies, nil
}

// MarkSolution marks a reply as the accepted solution for a topic.
func (s *Store) MarkSolution(topicID int64, replyID int64, currentUserID int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Verify topic owner or admin
	var topicOwnerID int64
	if err := tx.QueryRow(`SELECT user_id FROM forum_topics WHERE id = ?`, topicID).Scan(&topicOwnerID); err != nil {
		return fmt.Errorf("topic not found: %w", err)
	}
	if topicOwnerID != currentUserID {
		return fmt.Errorf("only topic creator can accept solution")
	}

	// Reset all solutions for this topic
	if _, err := tx.Exec(`UPDATE forum_replies SET is_solution = 0 WHERE topic_id = ?`, topicID); err != nil {
		return err
	}

	// Set this reply as solution
	if replyID > 0 {
		if _, err := tx.Exec(`UPDATE forum_replies SET is_solution = 1 WHERE id = ? AND topic_id = ?`, replyID, topicID); err != nil {
			return err
		}
		if _, err := tx.Exec(`UPDATE forum_topics SET solved_reply_id = ? WHERE id = ?`, replyID, topicID); err != nil {
			return err
		}
	} else {
		// Unmark solution
		if _, err := tx.Exec(`UPDATE forum_topics SET solved_reply_id = NULL WHERE id = ?`, topicID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// DeleteReply removes or redacts a reply (author or admin).
func (s *Store) DeleteReply(replyID int64, currentUserID int64, isAdmin bool) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var ownerID int64
	var topicID int64
	if err := tx.QueryRow(`SELECT user_id, topic_id FROM forum_replies WHERE id = ?`, replyID).Scan(&ownerID, &topicID); err != nil {
		return fmt.Errorf("reply not found: %w", err)
	}

	if !isAdmin && ownerID != currentUserID {
		return fmt.Errorf("unauthorized to delete reply")
	}

	// Check if this reply has child replies
	var childCount int
	if err := tx.QueryRow(`SELECT COUNT(1) FROM forum_replies WHERE parent_id = ?`, replyID).Scan(&childCount); err != nil {
		return err
	}

	if childCount > 0 {
		// Redact content if children exist to maintain conversation tree
		if _, err := tx.Exec(`UPDATE forum_replies SET body_md = '[Komentar dihapus oleh penulis]', body_html = '<p class="text-muted italic">[Komentar dihapus oleh penulis]</p>' WHERE id = ?`, replyID); err != nil {
			return err
		}
	} else {
		// Hard delete if no child replies
		if _, err := tx.Exec(`DELETE FROM forum_votes WHERE target_type = 'reply' AND target_id = ?`, replyID); err != nil {
			return err
		}
		if _, err := tx.Exec(`DELETE FROM forum_replies WHERE id = ?`, replyID); err != nil {
			return err
		}
		if _, err := tx.Exec(`UPDATE forum_topics SET replies_count = MAX(0, replies_count - 1) WHERE id = ?`, topicID); err != nil {
			return err
		}
	}

	return tx.Commit()
}
