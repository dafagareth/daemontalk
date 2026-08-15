// Package postdb menyimpan post yang ditulis lewat editor web di SQLite.
// Post file markdown di content/posts tetap menjadi sumber terpisah; keduanya
// digabung di layer handler.
package postdb

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type WebPost struct {
	ID          int64
	Slug        string
	Title       string
	Description string
	BodyMD      string
	Tags        string // dipisah koma: "go, web"
	Lang        string
	Cover       string
	Draft       bool
	Date        string // YYYY-MM-DD (tanggal publish)
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

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
		
		CREATE TABLE IF NOT EXISTS posts (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			slug        TEXT NOT NULL UNIQUE,
			title       TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			body_md     TEXT NOT NULL,
			tags        TEXT NOT NULL DEFAULT '',
			lang        TEXT NOT NULL DEFAULT 'id',
			cover       TEXT NOT NULL DEFAULT '',
			draft       INTEGER NOT NULL DEFAULT 1,
			date        TEXT NOT NULL,
			created_at  DATETIME NOT NULL,
			updated_at  DATETIME NOT NULL
		);
	`); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return &Store{db: db}, nil
}

const cols = "id, slug, title, description, body_md, tags, lang, cover, draft, date, created_at, updated_at"

func scan(row interface{ Scan(...any) error }) (WebPost, error) {
	var p WebPost
	err := row.Scan(&p.ID, &p.Slug, &p.Title, &p.Description, &p.BodyMD,
		&p.Tags, &p.Lang, &p.Cover, &p.Draft, &p.Date, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

// List returns all web posts, newest publish date first.
func (s *Store) List() ([]WebPost, error) {
	rows, err := s.db.Query("SELECT " + cols + " FROM posts ORDER BY date DESC, id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []WebPost
	for rows.Next() {
		p, err := scan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *Store) Get(id int64) (WebPost, error) {
	return scan(s.db.QueryRow("SELECT "+cols+" FROM posts WHERE id = ?", id))
}

func (s *Store) GetBySlug(slug string) (WebPost, error) {
	return scan(s.db.QueryRow("SELECT "+cols+" FROM posts WHERE slug = ?", slug))
}

func (s *Store) Create(p WebPost) (int64, error) {
	now := time.Now().UTC()
	res, err := s.db.Exec(`
		INSERT INTO posts (slug, title, description, body_md, tags, lang, cover, draft, date, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.Slug, p.Title, p.Description, p.BodyMD, p.Tags, p.Lang, p.Cover, p.Draft, p.Date, now, now)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) Update(p WebPost) error {
	_, err := s.db.Exec(`
		UPDATE posts SET slug = ?, title = ?, description = ?, body_md = ?, tags = ?,
			lang = ?, cover = ?, draft = ?, date = ?, updated_at = ?
		WHERE id = ?`,
		p.Slug, p.Title, p.Description, p.BodyMD, p.Tags, p.Lang, p.Cover, p.Draft, p.Date,
		time.Now().UTC(), p.ID)
	return err
}

func (s *Store) Delete(id int64) error {
	_, err := s.db.Exec("DELETE FROM posts WHERE id = ?", id)
	return err
}

// ToMarkdown menyintesis dokumen frontmatter+body sehingga post DB bisa
// dirender lewat post.Parse — pipeline yang sama dengan post file.
func (p WebPost) ToMarkdown() []byte {
	var b strings.Builder
	b.WriteString("---\n")
	fmt.Fprintf(&b, "title: %q\n", p.Title)
	fmt.Fprintf(&b, "slug: %s\n", p.Slug)
	if p.Date != "" {
		fmt.Fprintf(&b, "date: %s\n", p.Date)
	}
	if p.Lang != "" {
		fmt.Fprintf(&b, "lang: %s\n", p.Lang)
	}
	fmt.Fprintf(&b, "draft: %t\n", p.Draft)
	if p.Cover != "" {
		fmt.Fprintf(&b, "cover: %q\n", p.Cover)
	}
	if tags := p.TagList(); len(tags) > 0 {
		quoted := make([]string, len(tags))
		for i, t := range tags {
			quoted[i] = fmt.Sprintf("%q", t)
		}
		fmt.Fprintf(&b, "tags: [%s]\n", strings.Join(quoted, ", "))
	}
	b.WriteString("---\n\n")
	b.WriteString(p.BodyMD)
	return []byte(b.String())
}

// TagList memecah field Tags ("a, b") menjadi slice bersih.
func (p WebPost) TagList() []string {
	var out []string
	for _, t := range strings.Split(p.Tags, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}

// Close closes the underlying SQLite database connection.
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
