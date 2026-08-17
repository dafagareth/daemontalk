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
		`SELECT id, post_slug, name, body, parent_id, created_at FROM comments WHERE post_slug = ? ORDER BY created_at ASC`,
		slug,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Comment
	for rows.Next() {
		var c Comment
		var parentID sql.NullInt64
		if err := rows.Scan(&c.ID, &c.PostSlug, &c.Name, &c.Body, &parentID, &c.CreatedAt); err != nil {
			return nil, err
		}
		if parentID.Valid {
			val := parentID.Int64
			c.ParentID = &val
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// Add validates and inserts a new root comment, returning the stored row.
func (s *Store) Add(slug, name, body string) (Comment, error) {
	return s.AddWithParent(slug, name, body, nil)
}

// AddWithParent validates and inserts a new comment or reply.
func (s *Store) AddWithParent(slug, name, body string, parentID *int64) (Comment, error) {
	name = strings.TrimSpace(name)
	body = strings.TrimSpace(body)

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
	if parentID != nil && *parentID > 0 {
		var parentSlug string
		err := s.db.QueryRow(`SELECT post_slug FROM comments WHERE id = ?`, *parentID).Scan(&parentSlug)
		if err != nil || parentSlug != slug {
			// If parent not found or slug mismatch, reject or fallback to root
			parentID = nil
		}
	}

	c := Comment{
		PostSlug:  slug,
		Name:      name,
		Body:      body,
		ParentID:  parentID,
		CreatedAt: time.Now().UTC(),
	}

	var res sql.Result
	var err error
	if parentID != nil && *parentID > 0 {
		res, err = s.db.Exec(
			`INSERT INTO comments (post_slug, name, body, parent_id, created_at) VALUES (?, ?, ?, ?, ?)`,
			c.PostSlug, c.Name, c.Body, *c.ParentID, c.CreatedAt,
		)
	} else {
		res, err = s.db.Exec(
			`INSERT INTO comments (post_slug, name, body, parent_id, created_at) VALUES (?, ?, ?, NULL, ?)`,
			c.PostSlug, c.Name, c.Body, c.CreatedAt,
		)
	}

	if err != nil {
		return Comment{}, err
	}
	c.ID, _ = res.LastInsertId()
	return c, nil
}

// Delete removes a comment by ID and its recursive descendants.
func (s *Store) Delete(id int64) error {
	_, err := s.db.Exec(`DELETE FROM comments WHERE id = ? OR parent_id = ?`, id, id)
	return err
}

// ListAll returns all comments across every post, newest first.
func (s *Store) ListAll() ([]Comment, error) {
	rows, err := s.db.Query(
		`SELECT id, post_slug, name, body, parent_id, created_at FROM comments ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Comment
	for rows.Next() {
		var c Comment
		var parentID sql.NullInt64
		if err := rows.Scan(&c.ID, &c.PostSlug, &c.Name, &c.Body, &parentID, &c.CreatedAt); err != nil {
			return nil, err
		}
		if parentID.Valid {
			val := parentID.Int64
			c.ParentID = &val
		}
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
