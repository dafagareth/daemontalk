package comment

import (
	"database/sql"
	"errors"
	"fmt"

	_ "modernc.org/sqlite"
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
		PRAGMA foreign_keys=ON;
		PRAGMA temp_store=MEMORY;
		PRAGMA mmap_size=268435456;
		PRAGMA cache_size=-2000;
		
		CREATE TABLE IF NOT EXISTS comments (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			post_slug  TEXT NOT NULL,
			name       TEXT NOT NULL,
			body       TEXT NOT NULL,
			parent_id  INTEGER DEFAULT NULL REFERENCES comments(id) ON DELETE CASCADE,
			created_at DATETIME NOT NULL
		);

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

	// Migrate existing database instances that might not have parent_id column yet
	_, _ = db.Exec(`ALTER TABLE comments ADD COLUMN parent_id INTEGER DEFAULT NULL REFERENCES comments(id) ON DELETE CASCADE;`)
	_, _ = db.Exec(`ALTER TABLE comments ADD COLUMN user_id INTEGER DEFAULT NULL;`)
	_, _ = db.Exec(`ALTER TABLE comments ADD COLUMN avatar_url TEXT DEFAULT '';`)
	_, _ = db.Exec(`ALTER TABLE comments ADD COLUMN is_verified BOOLEAN DEFAULT 0;`)
	_, _ = db.Exec(`ALTER TABLE comments ADD COLUMN github_url TEXT DEFAULT '';`)
	_, _ = db.Exec(`ALTER TABLE comments ADD COLUMN is_reported BOOLEAN DEFAULT 0;`)

	// Create indexes safely after ensuring all columns exist
	if _, err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_comments_slug ON comments(post_slug, created_at);
		CREATE INDEX IF NOT EXISTS idx_comments_parent ON comments(parent_id);
	`); err != nil {
		return nil, fmt.Errorf("migrate indexes: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }
