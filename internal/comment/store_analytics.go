package comment

import (
	"database/sql"
)

// PageView is one row of the lightweight per-path analytics counter.
type PageView struct {
	Path  string
	Count int
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

// IncrementPageView bumps the hit counter for a request path.
func (s *Store) IncrementPageView(path string) error {
	_, err := s.db.Exec(`
		INSERT INTO pageviews (path, count) VALUES (?, 1)
		ON CONFLICT(path) DO UPDATE SET count = count + 1
	`, path)
	return err
}

// TopPageViews returns the most-visited legitimate paths, highest first.
func (s *Store) TopPageViews(limit int) ([]PageView, error) {
	rows, err := s.db.Query(`
		SELECT path, count FROM pageviews 
		WHERE path NOT LIKE '%.php%' 
		  AND path NOT LIKE '%/comments%' 
		  AND path NOT LIKE '%wp-%'
		  AND path NOT LIKE '%.env%'
		  AND path NOT LIKE '%/api/%'
		ORDER BY count DESC LIMIT ?
	`, limit)
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
