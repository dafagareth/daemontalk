package comment

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

const (
	MaxNameLen = 50
	MaxBodyLen = 2000
)

// ErrInvalid is returned when a comment fails validation.
var ErrInvalid = errors.New("invalid comment")

type Store struct {
	db *sql.DB
}

// Open opens (or creates) the SQLite database at path and ensures the schema.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if _, err := db.Exec(`
		PRAGMA journal_mode=WAL;
		PRAGMA synchronous=NORMAL;
		PRAGMA busy_timeout=5000;
		
		CREATE TABLE IF NOT EXISTS comments (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			post_slug  TEXT NOT NULL,
			name       TEXT NOT NULL,
			body       TEXT NOT NULL,
			created_at DATETIME NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_comments_slug ON comments(post_slug, created_at);

		CREATE TABLE IF NOT EXISTS views (
			post_slug TEXT PRIMARY KEY,
			count     INTEGER NOT NULL DEFAULT 0
		);

		CREATE TABLE IF NOT EXISTS pageviews (
			path  TEXT PRIMARY KEY,
			count INTEGER NOT NULL DEFAULT 0
		);

		CREATE TABLE IF NOT EXISTS reactions (
			post_slug TEXT NOT NULL,
			emoji     TEXT NOT NULL,
			count     INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY (post_slug, emoji)
		);
	`); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

// ListBySlug returns all comments for a post, oldest first.
func (s *Store) ListBySlug(slug string) ([]Comment, error) {
	rows, err := s.db.Query(
		`SELECT id, post_slug, name, body, created_at FROM comments WHERE post_slug = ? ORDER BY created_at ASC`,
		slug,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.PostSlug, &c.Name, &c.Body, &c.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// Add validates and inserts a new comment, returning the stored row.
func (s *Store) Add(slug, name, body string) (Comment, error) {
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

	c := Comment{
		PostSlug:  slug,
		Name:      name,
		Body:      body,
		CreatedAt: time.Now().UTC(),
	}
	res, err := s.db.Exec(
		`INSERT INTO comments (post_slug, name, body, created_at) VALUES (?, ?, ?, ?)`,
		c.PostSlug, c.Name, c.Body, c.CreatedAt,
	)
	if err != nil {
		return Comment{}, err
	}
	c.ID, _ = res.LastInsertId()
	return c, nil
}

// Delete removes a comment by ID (used for moderation).
func (s *Store) Delete(id int64) error {
	_, err := s.db.Exec(`DELETE FROM comments WHERE id = ?`, id)
	return err
}

// IncrementView bumps the view counter for a post and returns the new total.
func (s *Store) IncrementView(slug string) (int, error) {
	if _, err := s.db.Exec(`
		INSERT INTO views (post_slug, count) VALUES (?, 1)
		ON CONFLICT(post_slug) DO UPDATE SET count = count + 1
	`, slug); err != nil {
		return 0, err
	}
	return s.ViewCount(slug)
}

// ViewCount returns the current view total for a post.
func (s *Store) ViewCount(slug string) (int, error) {
	var n int
	err := s.db.QueryRow(`SELECT count FROM views WHERE post_slug = ?`, slug).Scan(&n)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return n, err
}

// ListAll returns all comments across every post, newest first.
func (s *Store) ListAll() ([]Comment, error) {
	rows, err := s.db.Query(
		`SELECT id, post_slug, name, body, created_at FROM comments ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.PostSlug, &c.Name, &c.Body, &c.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// PageView is one row of the lightweight per-path analytics counter.
type PageView struct {
	Path  string
	Count int
}

// IncrementPageView bumps the hit counter for a request path.
func (s *Store) IncrementPageView(path string) error {
	_, err := s.db.Exec(`
		INSERT INTO pageviews (path, count) VALUES (?, 1)
		ON CONFLICT(path) DO UPDATE SET count = count + 1
	`, path)
	return err
}

// TopPageViews returns the most-visited paths, highest first.
func (s *Store) TopPageViews(limit int) ([]PageView, error) {
	rows, err := s.db.Query(
		`SELECT path, count FROM pageviews ORDER BY count DESC LIMIT ?`, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []PageView
	for rows.Next() {
		var pv PageView
		if err := rows.Scan(&pv.Path, &pv.Count); err != nil {
			return nil, err
		}
		out = append(out, pv)
	}
	return out, rows.Err()
}

// TotalPageViews returns the sum of all path hit counts.
func (s *Store) TotalPageViews() (int, error) {
	var n int
	err := s.db.QueryRow(`SELECT COALESCE(SUM(count), 0) FROM pageviews`).Scan(&n)
	return n, err
}

// GetReactions returns emoji→count map for a post.
func (s *Store) GetReactions(slug string) (map[string]int, error) {
	rows, err := s.db.Query(`SELECT emoji, count FROM reactions WHERE post_slug = ?`, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string]int)
	for rows.Next() {
		var emoji string
		var count int
		if err := rows.Scan(&emoji, &count); err != nil {
			return nil, err
		}
		out[emoji] = count
	}
	return out, rows.Err()
}

// IncrementReaction bumps the reaction count for a post+emoji and returns all reactions for that post.
func (s *Store) IncrementReaction(slug, emoji string) (map[string]int, error) {
	_, err := s.db.Exec(`
		INSERT INTO reactions (post_slug, emoji, count) VALUES (?, ?, 1)
		ON CONFLICT(post_slug, emoji) DO UPDATE SET count = count + 1
	`, slug, emoji)
	if err != nil {
		return nil, err
	}
	return s.GetReactions(slug)
}

// DecrementReaction decrements the reaction count for a post+emoji (capping at 0) and returns all reactions for that post.
func (s *Store) DecrementReaction(slug, emoji string) (map[string]int, error) {
	_, err := s.db.Exec(`
		UPDATE reactions SET count = CASE WHEN count > 0 THEN count - 1 ELSE 0 END 
		WHERE post_slug = ? AND emoji = ?
	`, slug, emoji)
	if err != nil {
		return nil, err
	}
	return s.GetReactions(slug)
}


// AllViewCounts returns a slug→count map for every post that has been viewed.
func (s *Store) AllViewCounts() (map[string]int, error) {
	rows, err := s.db.Query(`SELECT post_slug, count FROM views`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string]int)
	for rows.Next() {
		var slug string
		var count int
		if err := rows.Scan(&slug, &count); err != nil {
			return nil, err
		}
		out[slug] = count
	}
	return out, rows.Err()
}
