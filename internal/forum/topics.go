package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"regexp"
	"strings"
	"time"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-z0-9]+`)

// Slugify generates a clean URL slug from a title.
func Slugify(title string) string {
	s := strings.ToLower(strings.TrimSpace(title))
	s = nonAlphanumericRegex.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if len(s) > 80 {
		s = s[:80]
	}
	if s == "" {
		s = fmt.Sprintf("topic-%d", time.Now().Unix())
	}
	return s
}

// CreateTopic adds a new discussion thread or question.
func (s *Store) CreateTopic(t Topic) (*Topic, error) {
	now := time.Now().UTC()
	baseSlug := Slugify(t.Title)
	slug := baseSlug

	// Ensure unique slug
	for i := 1; ; i++ {
		var exists int
		_ = s.db.QueryRow(`SELECT COUNT(1) FROM forum_topics WHERE slug = ?`, slug).Scan(&exists)
		if exists == 0 {
			break
		}
		slug = fmt.Sprintf("%s-%d", baseSlug, i)
	}

	tagsJSON, err := json.Marshal(t.Tags)
	if err != nil {
		tagsJSON = []byte("[]")
	}

	bodyHTML := string(RenderMarkdown(t.BodyMD))

	query := `
		INSERT INTO forum_topics (user_id, title, slug, category, tags, body_md, body_html, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, created_at, updated_at
	`
	var id int64
	var createdAt, updatedAt time.Time
	err = s.db.QueryRow(query, t.UserID, t.Title, slug, t.Category, string(tagsJSON), t.BodyMD, bodyHTML, now, now).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("create topic: %w", err)
	}

	t.ID = id
	t.Slug = slug
	t.BodyHTML = template.HTML(bodyHTML)
	t.CreatedAt = createdAt
	t.UpdatedAt = updatedAt
	return &t, nil
}

// GetTopicBySlug retrieves a topic with author details and vote status.
func (s *Store) GetTopicBySlug(slug string, currentUserID int64) (*Topic, error) {
	query := `
		SELECT t.id, t.user_id, u.display_name, u.username, u.avatar_url, u.github_url,
		       t.title, t.slug, t.category, t.tags, t.body_md, t.body_html,
		       t.solved_reply_id, t.views_count, t.votes_count, t.replies_count, t.pinned,
		       t.created_at, t.updated_at,
		       EXISTS(SELECT 1 FROM forum_votes WHERE user_id = ? AND target_type = 'topic' AND target_id = t.id) AS user_voted
		FROM forum_topics t
		LEFT JOIN users u ON u.id = t.user_id
		WHERE t.slug = ?
	`
	var t Topic
	var tagsJSON string
	var solvedID sql.NullInt64
	var userVoted int
	var authorName, authorUser, authorAvatar, authorGH sql.NullString

	err := s.db.QueryRow(query, currentUserID, slug).Scan(
		&t.ID, &t.UserID, &authorName, &authorUser, &authorAvatar, &authorGH,
		&t.Title, &t.Slug, &t.Category, &tagsJSON, &t.BodyMD, &t.BodyHTML,
		&solvedID, &t.ViewsCount, &t.VotesCount, &t.RepliesCount, &t.Pinned,
		&t.CreatedAt, &t.UpdatedAt, &userVoted,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get topic by slug: %w", err)
	}

	t.AuthorName = authorName.String
	t.AuthorUsername = authorUser.String
	t.AuthorAvatar = authorAvatar.String
	t.AuthorGitHub = authorGH.String
	if t.AuthorName == "" {
		t.AuthorName = "Deleted User"
	}
	if t.AuthorUsername == "" {
		t.AuthorUsername = "ghost"
	}
	if t.AuthorAvatar == "" {
		t.AuthorAvatar = "/static/images/deleted-user.png"
	}
	if solvedID.Valid {
		t.SolvedReplyID = solvedID.Int64
		t.IsSolved = true
	}
	t.UserVoted = (userVoted == 1)
	t.IsOwner = (currentUserID > 0 && t.UserID == currentUserID)

	_ = json.Unmarshal([]byte(tagsJSON), &t.Tags)
	return &t, nil
}

// ListTopics retrieves topics with filtering, pagination, and sorting.
func (s *Store) ListTopics(category, tag, search, author, sortOrder string, limit, offset int, currentUserID int64) ([]*Topic, int, error) {
	whereClauses := []string{"1=1"}
	var whereArgs []any

	if category != "" && category != "all" {
		whereClauses = append(whereClauses, "t.category = ?")
		whereArgs = append(whereArgs, category)
	}
	if tag != "" {
		whereClauses = append(whereClauses, "t.tags LIKE ?")
		whereArgs = append(whereArgs, "%\""+tag+"\"%")
	}
	if search != "" {
		whereClauses = append(whereClauses, "(t.title LIKE ? OR t.body_md LIKE ?)")
		whereArgs = append(whereArgs, "%"+search+"%", "%"+search+"%")
	}
	if author != "" {
		whereClauses = append(whereClauses, "u.username = ? COLLATE NOCASE")
		whereArgs = append(whereArgs, author)
	}

	whereSQL := strings.Join(whereClauses, " AND ")

	// Total count
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(1) FROM forum_topics t LEFT JOIN users u ON u.id = t.user_id WHERE %s", whereSQL)
	if err := s.db.QueryRow(countQuery, whereArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	orderBy := "t.pinned DESC, t.created_at DESC"
	switch sortOrder {
	case "votes":
		orderBy = "t.pinned DESC, t.votes_count DESC, t.created_at DESC"
	case "replies":
		orderBy = "t.pinned DESC, t.replies_count DESC, t.created_at DESC"
	case "unsolved":
		whereSQL += " AND t.solved_reply_id IS NULL"
		orderBy = "t.pinned DESC, t.created_at DESC"
	}

	query := fmt.Sprintf(`
		SELECT t.id, t.user_id, u.display_name, u.username, u.avatar_url, u.github_url,
		       t.title, t.slug, t.category, t.tags, t.body_md, t.body_html,
		       t.solved_reply_id, t.views_count, t.votes_count, t.replies_count, t.pinned,
		       t.created_at, t.updated_at,
		       EXISTS(SELECT 1 FROM forum_votes WHERE user_id = ? AND target_type = 'topic' AND target_id = t.id) AS user_voted
		FROM forum_topics t
		LEFT JOIN users u ON u.id = t.user_id
		WHERE %s
		ORDER BY %s
		LIMIT ? OFFSET ?
	`, whereSQL, orderBy)

	queryArgs := make([]any, 0, len(whereArgs)+3)
	queryArgs = append(queryArgs, currentUserID)
	queryArgs = append(queryArgs, whereArgs...)
	queryArgs = append(queryArgs, limit, offset)

	rows, err := s.db.Query(query, queryArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("list topics: %w", err)
	}
	defer rows.Close()

	var topics []*Topic
	for rows.Next() {
		var t Topic
		var tagsJSON string
		var solvedID sql.NullInt64
		var userVoted int
		var authorName, authorUser, authorAvatar, authorGH sql.NullString

		if err := rows.Scan(
			&t.ID, &t.UserID, &authorName, &authorUser, &authorAvatar, &authorGH,
			&t.Title, &t.Slug, &t.Category, &tagsJSON, &t.BodyMD, &t.BodyHTML,
			&solvedID, &t.ViewsCount, &t.VotesCount, &t.RepliesCount, &t.Pinned,
			&t.CreatedAt, &t.UpdatedAt, &userVoted,
		); err != nil {
			continue
		}

		t.AuthorName = authorName.String
		t.AuthorUsername = authorUser.String
		t.AuthorAvatar = authorAvatar.String
		t.AuthorGitHub = authorGH.String
		if t.AuthorName == "" {
			t.AuthorName = "Deleted User"
		}
		if t.AuthorUsername == "" {
			t.AuthorUsername = "ghost"
		}
		if t.AuthorAvatar == "" {
			t.AuthorAvatar = "/static/images/deleted-user.png"
		}
		if solvedID.Valid {
			t.SolvedReplyID = solvedID.Int64
			t.IsSolved = true
		}
		t.UserVoted = (userVoted == 1)
		_ = json.Unmarshal([]byte(tagsJSON), &t.Tags)
		topics = append(topics, &t)
	}

	return topics, total, nil
}

// IncrementTopicViews increments view counter for a topic.
func (s *Store) IncrementTopicViews(topicID int64) {
	_, _ = s.db.Exec(`UPDATE forum_topics SET views_count = views_count + 1 WHERE id = ?`, topicID)
}

// DeleteTopic deletes a topic and all its replies and votes (owner or admin).
func (s *Store) DeleteTopic(topicID int64, currentUserID int64, isAdmin bool) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var ownerID int64
	if err := tx.QueryRow(`SELECT user_id FROM forum_topics WHERE id = ?`, topicID).Scan(&ownerID); err != nil {
		return fmt.Errorf("topic not found: %w", err)
	}

	if !isAdmin && ownerID != currentUserID {
		return fmt.Errorf("unauthorized to delete topic")
	}

	if _, err := tx.Exec(`DELETE FROM forum_votes WHERE target_type = 'topic' AND target_id = ?`, topicID); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM forum_replies WHERE topic_id = ?`, topicID); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM forum_topics WHERE id = ?`, topicID); err != nil {
		return err
	}

	return tx.Commit()
}
