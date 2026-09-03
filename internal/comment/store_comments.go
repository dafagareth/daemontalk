package comment

import (
	"database/sql"
	"strings"
	"time"
)

const (
	MaxNameLen = 50
	MaxBodyLen = 2000
)

// ListBySlug returns all comments for a post, oldest first.
func (s *Store) ListBySlug(slug string) ([]Comment, error) {
	rows, err := s.db.Query(
		`SELECT id, post_slug, name, body, parent_id, created_at, user_id, avatar_url, is_verified, github_url, is_reported 
		 FROM comments WHERE post_slug = ? ORDER BY created_at ASC`,
		slug,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Comment
	for rows.Next() {
		var c Comment
		var parentID, userID sql.NullInt64
		var avatarURL, ghURL sql.NullString
		var isVerified sql.NullBool
		var isReported sql.NullBool
		if err := rows.Scan(&c.ID, &c.PostSlug, &c.Name, &c.Body, &parentID, &c.CreatedAt, &userID, &avatarURL, &isVerified, &ghURL, &isReported); err != nil {
			return nil, err
		}
		if parentID.Valid {
			val := parentID.Int64
			c.ParentID = &val
		}
		if userID.Valid {
			val := userID.Int64
			c.UserID = &val
		}
		c.AvatarURL = avatarURL.String
		c.IsVerified = isVerified.Bool
		c.GitHubURL = ghURL.String
		c.IsReported = isReported.Bool
		out = append(out, c)
	}
	return out, rows.Err()
}

// AddAdvanced inserts a comment with optional user authentication info.
func (s *Store) AddAdvanced(c Comment) (Comment, error) {
	name := strings.TrimSpace(c.Name)
	body := strings.TrimSpace(c.Body)

	if name == "" || body == "" {
		return Comment{}, ErrInvalid
	}
	if len([]rune(name)) > MaxNameLen {
		name = string([]rune(name)[:MaxNameLen])
	}
	if len([]rune(body)) > MaxBodyLen {
		body = string([]rune(body)[:MaxBodyLen])
	}

	// If replying to a parent comment, verify the parent exists and belongs to the same post
	if c.ParentID != nil && *c.ParentID > 0 {
		var parentSlug string
		err := s.db.QueryRow(`SELECT post_slug FROM comments WHERE id = ?`, *c.ParentID).Scan(&parentSlug)
		if err != nil || parentSlug != c.PostSlug {
			c.ParentID = nil
		}
	}

	c.Name = name
	c.Body = body
	c.CreatedAt = time.Now().UTC()

	var res sql.Result
	var err error
	var parentArg any
	if c.ParentID != nil && *c.ParentID > 0 {
		parentArg = *c.ParentID
	}
	var userArg any
	if c.UserID != nil && *c.UserID > 0 {
		userArg = *c.UserID
	}

	query := `
		INSERT INTO comments (post_slug, name, body, parent_id, created_at, user_id, avatar_url, is_verified, github_url)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	res, err = s.db.Exec(query, c.PostSlug, c.Name, c.Body, parentArg, c.CreatedAt, userArg, c.AvatarURL, c.IsVerified, c.GitHubURL)
	if err != nil {
		return Comment{}, err
	}
	c.ID, _ = res.LastInsertId()
	return c, nil
}

// GetByID retrieves a single comment by ID.
func (s *Store) GetByID(id int64) (*Comment, error) {
	var c Comment
	var parentID, userID sql.NullInt64
	var avatarURL, ghURL sql.NullString
	var isVerified sql.NullBool
		var isReported sql.NullBool
	err := s.db.QueryRow(
		`SELECT id, post_slug, name, body, parent_id, created_at, user_id, avatar_url, is_verified, github_url, is_reported 
		 FROM comments WHERE id = ?`,
		id,
	).Scan(&c.ID, &c.PostSlug, &c.Name, &c.Body, &parentID, &c.CreatedAt, &userID, &avatarURL, &isVerified, &ghURL, &isReported)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if parentID.Valid {
		val := parentID.Int64
		c.ParentID = &val
	}
	if userID.Valid {
		val := userID.Int64
		c.UserID = &val
	}
	c.AvatarURL = avatarURL.String
	c.IsVerified = isVerified.Bool
	c.GitHubURL = ghURL.String
		c.IsReported = isReported.Bool
	return &c, nil
}

// Delete removes a comment by ID and its recursive descendants via CASCADE.
func (s *Store) Delete(id int64) error {
	_, err := s.db.Exec(`DELETE FROM comments WHERE id = ?`, id)
	return err
}

// ListAll returns all comments across every post, newest first.
func (s *Store) ListAll() ([]Comment, error) {
	rows, err := s.db.Query(
		`SELECT id, post_slug, name, body, parent_id, created_at, user_id, avatar_url, is_verified, github_url, is_reported 
		 FROM comments ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Comment
	for rows.Next() {
		var c Comment
		var parentID, userID sql.NullInt64
		var avatarURL, ghURL sql.NullString
		var isVerified sql.NullBool
		var isReported sql.NullBool
		if err := rows.Scan(&c.ID, &c.PostSlug, &c.Name, &c.Body, &parentID, &c.CreatedAt, &userID, &avatarURL, &isVerified, &ghURL, &isReported); err != nil {
			return nil, err
		}
		if parentID.Valid {
			val := parentID.Int64
			c.ParentID = &val
		}
		if userID.Valid {
			val := userID.Int64
			c.UserID = &val
		}
		c.AvatarURL = avatarURL.String
		c.IsVerified = isVerified.Bool
		c.GitHubURL = ghURL.String
		c.IsReported = isReported.Bool
		out = append(out, c)
	}
	return out, rows.Err()
}

// BuildTree converts a flat list of comments into a hierarchical tree structure.
func BuildTree(comments []Comment) []Comment {
	if len(comments) == 0 {
		return nil
	}

	type node struct {
		c       Comment
		replies []*node
	}

	nodes := make(map[int64]*node, len(comments))
	var rootNodes []*node

	for _, c := range comments {
		nodes[c.ID] = &node{c: c}
	}

	for _, c := range comments {
		n := nodes[c.ID]
		if c.ParentID != nil && *c.ParentID > 0 {
			if parent, ok := nodes[*c.ParentID]; ok {
				parent.replies = append(parent.replies, n)
				continue
			}
		}
		rootNodes = append(rootNodes, n)
	}

	var convert func(n *node) Comment
	convert = func(n *node) Comment {
		res := n.c
		if len(n.replies) > 0 {
			res.Replies = make([]Comment, len(n.replies))
			for i, child := range n.replies {
				res.Replies[i] = convert(child)
			}
		}
		return res
	}

	result := make([]Comment, len(rootNodes))
	for i, r := range rootNodes {
		result[i] = convert(r)
	}
	return result
}

// AnonymizeUserComments anonymizes comments authored by a user.
func (s *Store) AnonymizeUserComments(userID int64) error {
	_, err := s.db.Exec(`
		UPDATE comments 
		SET user_id = NULL, name = 'Deleted User', avatar_url = '/static/images/deleted-user.png', is_verified = 0, github_url = ''
		WHERE user_id = ?
	`, userID)
	return err
}

// ListByUserID returns all comments authored by a user.
func (s *Store) ListByUserID(userID int64) ([]Comment, error) {
	rows, err := s.db.Query(
		`SELECT id, post_slug, name, body, parent_id, created_at, user_id, avatar_url, is_verified, github_url, is_reported 
		 FROM comments WHERE user_id = ? ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Comment
	for rows.Next() {
		var c Comment
		var parentID, uid sql.NullInt64
		var avatarURL, ghURL sql.NullString
		var isVerified sql.NullBool
		var isReported sql.NullBool
		if err := rows.Scan(&c.ID, &c.PostSlug, &c.Name, &c.Body, &parentID, &c.CreatedAt, &uid, &avatarURL, &isVerified, &ghURL, &isReported); err != nil {
			return nil, err
		}
		if parentID.Valid {
			val := parentID.Int64
			c.ParentID = &val
		}
		if uid.Valid {
			val := uid.Int64
			c.UserID = &val
		}
		c.AvatarURL = avatarURL.String
		c.IsVerified = isVerified.Bool
		c.GitHubURL = ghURL.String
		c.IsReported = isReported.Bool
		out = append(out, c)
	}
	return out, rows.Err()
}

// UpdateBody updates the body of a comment by ID.
func (s *Store) UpdateBody(id int64, body string) error {
	_, err := s.db.Exec(`UPDATE comments SET body = ? WHERE id = ?`, body, id)
	return err
}

// Report marks a comment as reported
func (s *Store) Report(id int64) error {
	_, err := s.db.Exec(`UPDATE comments SET is_reported = 1 WHERE id = ?`, id)
	return err
}
