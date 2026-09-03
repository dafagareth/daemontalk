package forum

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

// Open opens or creates the SQLite database for Discussions / Forum.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open forum db: %w", err)
	}

	if _, err := db.Exec(`
		PRAGMA journal_mode=WAL;
		PRAGMA synchronous=NORMAL;
		PRAGMA busy_timeout=5000;
		PRAGMA foreign_keys=ON;

		CREATE TABLE IF NOT EXISTS users (
			id           INTEGER PRIMARY KEY AUTOINCREMENT,
			provider     TEXT NOT NULL DEFAULT 'github',
			provider_id  TEXT NOT NULL UNIQUE,
			username     TEXT NOT NULL,
			display_name TEXT NOT NULL,
			email        TEXT,
			avatar_url   TEXT NOT NULL,
			github_url   TEXT NOT NULL,
			role         TEXT NOT NULL DEFAULT 'member',
			created_at   DATETIME NOT NULL,
			updated_at   DATETIME NOT NULL
		);

		CREATE TABLE IF NOT EXISTS forum_topics (
			id              INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id         INTEGER NOT NULL,
			title           TEXT NOT NULL,
			slug            TEXT NOT NULL UNIQUE,
			category        TEXT NOT NULL DEFAULT 'qna',
			tags            TEXT NOT NULL DEFAULT '[]',
			body_md         TEXT NOT NULL,
			body_html       TEXT NOT NULL,
			solved_reply_id INTEGER DEFAULT NULL,
			views_count     INTEGER NOT NULL DEFAULT 0,
			votes_count     INTEGER NOT NULL DEFAULT 0,
			replies_count   INTEGER NOT NULL DEFAULT 0,
			pinned          BOOLEAN NOT NULL DEFAULT 0,
			created_at      DATETIME NOT NULL,
			updated_at      DATETIME NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_topics_slug ON forum_topics(slug);
		CREATE INDEX IF NOT EXISTS idx_topics_cat ON forum_topics(category);
		CREATE INDEX IF NOT EXISTS idx_topics_created ON forum_topics(created_at DESC);

		CREATE TABLE IF NOT EXISTS forum_replies (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			topic_id    INTEGER NOT NULL REFERENCES forum_topics(id) ON DELETE CASCADE,
			parent_id   INTEGER DEFAULT NULL REFERENCES forum_replies(id) ON DELETE CASCADE,
			user_id     INTEGER NOT NULL,
			body_md     TEXT NOT NULL,
			body_html   TEXT NOT NULL,
			is_solution BOOLEAN NOT NULL DEFAULT 0,
			votes_count INTEGER NOT NULL DEFAULT 0,
			created_at  DATETIME NOT NULL,
			updated_at  DATETIME NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_replies_topic ON forum_replies(topic_id);
		CREATE INDEX IF NOT EXISTS idx_replies_parent ON forum_replies(parent_id);

		CREATE TABLE IF NOT EXISTS forum_votes (
			user_id     INTEGER NOT NULL,
			target_type TEXT NOT NULL, -- 'topic' or 'reply'
			target_id   INTEGER NOT NULL,
			created_at  DATETIME NOT NULL,
			PRIMARY KEY (user_id, target_type, target_id)
		);

		CREATE TABLE IF NOT EXISTS forum_topic_views (
			topic_id   INTEGER NOT NULL REFERENCES forum_topics(id) ON DELETE CASCADE,
			viewer_key TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			PRIMARY KEY (topic_id, viewer_key)
		);

		CREATE INDEX IF NOT EXISTS idx_topic_views_topic ON forum_topic_views(topic_id);
	`); err != nil {
		return nil, fmt.Errorf("migrate forum db: %w", err)
	}

	return &Store{db: db}, nil
}

// Close closes the database connection.
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// AnonymizeUser anonymizes user references in forum topics, replies, and deletes their votes.
func (s *Store) AnonymizeUser(userID int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Set user_id = 0 for topics and replies to preserve discussion integrity without author PII
	if _, err := tx.Exec(`UPDATE forum_topics SET user_id = 0 WHERE user_id = ?`, userID); err != nil {
		return err
	}
	if _, err := tx.Exec(`UPDATE forum_replies SET user_id = 0 WHERE user_id = ?`, userID); err != nil {
		return err
	}
	// Delete user votes
	if _, err := tx.Exec(`DELETE FROM forum_votes WHERE user_id = ?`, userID); err != nil {
		return err
	}

	return tx.Commit()
}
