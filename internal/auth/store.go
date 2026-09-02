package auth

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

// Open opens or creates the SQLite database for authentication.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open auth db: %w", err)
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

		CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
		CREATE INDEX IF NOT EXISTS idx_users_provider ON users(provider, provider_id);

		CREATE TABLE IF NOT EXISTS sessions (
			token_hash TEXT PRIMARY KEY,
			user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			expires_at DATETIME NOT NULL,
			created_at DATETIME NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
		CREATE INDEX IF NOT EXISTS idx_sessions_expires ON sessions(expires_at);
	`); err != nil {
		return nil, fmt.Errorf("migrate auth db: %w", err)
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

// UpsertUser inserts or updates a user from an OAuth provider.
func (s *Store) UpsertUser(u User) (*User, error) {
	now := time.Now().UTC()
	query := `
		INSERT INTO users (provider, provider_id, username, display_name, email, avatar_url, github_url, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(provider_id) DO UPDATE SET
			username = excluded.username,
			display_name = excluded.display_name,
			email = excluded.email,
			avatar_url = excluded.avatar_url,
			github_url = excluded.github_url,
			updated_at = excluded.updated_at
		RETURNING id, provider, provider_id, username, display_name, email, avatar_url, github_url, role, created_at, updated_at
	`
	var user User
	var email sql.NullString
	err := s.db.QueryRow(query,
		u.Provider, u.ProviderID, u.Username, u.DisplayName, u.Email, u.AvatarURL, u.GitHubURL, u.Role, now, now,
	).Scan(
		&user.ID, &user.Provider, &user.ProviderID, &user.Username, &user.DisplayName,
		&email, &user.AvatarURL, &user.GitHubURL, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("upsert user: %w", err)
	}
	if email.Valid {
		user.Email = email.String
	}
	return &user, nil
}

// CreateSession creates a persistent session token for a user.
func (s *Store) CreateSession(userID int64, tokenHash string, duration time.Duration) (*Session, error) {
	now := time.Now().UTC()
	expires := now.Add(duration)
	query := `INSERT INTO sessions (token_hash, user_id, expires_at, created_at) VALUES (?, ?, ?, ?)`
	_, err := s.db.Exec(query, tokenHash, userID, expires, now)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}
	return &Session{
		TokenHash: tokenHash,
		UserID:    userID,
		ExpiresAt: expires,
		CreatedAt: now,
	}, nil
}

// GetSessionUser returns the user associated with an active, unexpired session token.
func (s *Store) GetSessionUser(tokenHash string) (*User, error) {
	query := `
		SELECT u.id, u.provider, u.provider_id, u.username, u.display_name, u.email, u.avatar_url, u.github_url, u.role, u.created_at, u.updated_at
		FROM sessions s
		JOIN users u ON u.id = s.user_id
		WHERE s.token_hash = ? AND s.expires_at > ?
	`
	var user User
	var email sql.NullString
	err := s.db.QueryRow(query, tokenHash, time.Now().UTC()).Scan(
		&user.ID, &user.Provider, &user.ProviderID, &user.Username, &user.DisplayName,
		&email, &user.AvatarURL, &user.GitHubURL, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get session user: %w", err)
	}
	if email.Valid {
		user.Email = email.String
	}
	return &user, nil
}

// DeleteSession invalidates a session token.
func (s *Store) DeleteSession(tokenHash string) error {
	_, err := s.db.Exec(`DELETE FROM sessions WHERE token_hash = ?`, tokenHash)
	return err
}

// CountUsers returns the total count of registered community members.
func (s *Store) CountUsers() int {
	var count int
	_ = s.db.QueryRow(`SELECT COUNT(1) FROM users`).Scan(&count)
	return count
}

// DeleteUser permanently purges a user record and their active sessions.
func (s *Store) DeleteUser(userID int64) error {
	_, err := s.db.Exec(`DELETE FROM users WHERE id = ?`, userID)
	return err
}

// GetUserByUsername retrieves a user by their username.
func (s *Store) GetUserByUsername(username string) (*User, error) {
	var user User
	err := s.db.QueryRow(
		"SELECT id, provider, provider_id, username, display_name, avatar_url, github_url, email, role, created_at, updated_at FROM users WHERE username = ?",
		username,
	).Scan(
		&user.ID, &user.Provider, &user.ProviderID, &user.Username, &user.DisplayName,
		&user.AvatarURL, &user.GitHubURL, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
