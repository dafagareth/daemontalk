# Daemontalk System Architecture and Internal Codebase Reference

This document provides a comprehensive technical overview of the Daemontalk internal architecture, subsystem design, data flows, and storage mechanisms.

---

## 1. Subsystem Architecture Overview

Daemontalk is engineered as a single, self-contained Go binary that combines a server-rendered web application, a headless SSH TUI server, an embedded SQLite persistence layer, and static asset delivery.

```
                         +-----------------------------------+
                         |      External Client Requests     |
                         +-----------------------------------+
                                   /               \
                          HTTP/HTTPS               SSH :2222
                                 /                   \
           +-------------------------+     +-------------------------+
           |     Chi HTTP Router     |     |   Wish SSH Daemon / TUI |
           |   (internal/router)     |     |     (internal/tuisrv)   |
           +-------------------------+     +-------------------------+
                        |                               |
       +----------------+----------------+              |
       |                |                |              |
+--------------+ +--------------+ +--------------+ +---------------+
|  Auth/OAuth  | |  Forum Engine| | PostDB/Index | | Bubble Tea TUI|
|internal/auth | |internal/forum| |internal/post | | internal/tui  |
+--------------+ +--------------+ +--------------+ +---------------+
       |                |                |
       +----------------+----------------+
                        |
            +-----------------------+
            |  Embedded SQLite DB   |
            |      (WAL Mode)       |
            +-----------------------+
```

---

## 2. Internal Package Breakdown

### `internal/auth`
Manages GitHub OAuth 2.0 authentication, session lifecycle, and privacy operations.
* `Store`: SQLite repository managing `users` and `sessions` tables.
* `UpsertUser`: Atomically creates or updates user records using `provider_id` conflict resolution.
* `CreateSession`: Generates a high-entropy 256-bit cryptographic token via `crypto/rand`, stores its SHA-256 hash in SQLite, and sets an `HttpOnly`, `SameSite=Lax` session cookie.
* `GetUserByToken`: Resolves authenticated users from incoming session hashes.
* `DeleteUser`: Purges user record and active sessions upon self-service account termination.

### `internal/forum`
Implements the community discussions engine (`/discussions`).
* `Store`: Manages `topics`, `replies`, `topic_votes`, and `forum_topic_views` tables.
* `CreateTopic` & `CreateReply`: Handles Markdown parsing, `bluemonday.UGCPolicy()` HTML sanitization, tag normalization, and parent-child reply relationships.
* `VoteTopic` & `VoteReply`: Atomically manages user upvotes and prevents duplicate voting via composite primary keys `(target_id, user_id)`.
* `SetSolved`: Marks a specific reply as the accepted solution by the topic owner or admin.
* `RecordTopicView`: Deduplicates topic views in `forum_topic_views` based on unique visitor keys.
* `AnonymizeUser`: Converts all contributions of a deleted user to author `[Deleted User]` and username `ghost`, preserving discussion continuity.

### `internal/comment`
Provides the hierarchical commenting engine and unique view tracking for blog posts.
* `Store`: Embedded SQLite repository storing article comments, threaded replies, and `post_views`.
* `RecordPostView`: Deduplicates post views using `(post_slug, viewer_key)` composite index in `post_views`.
* `AddAdvanced`: Supports verified user attribution from GitHub sessions as well as deterministic guest handles (`anonym_<hex>`).
* `AnonymizeUserComments`: Replaces personal metadata with `Deleted User` and default avatar upon account purge.

### `internal/post` and `internal/postdb`
Handles Markdown article ingestion, frontmatter extraction, and in-memory indexing.
* `Post`: Struct representation containing metadata, reading time calculations, tags, multi-language mappings, and rendered HTML bodies.
* `Load`: Scans `content/posts/*.md`, extracts YAML frontmatter, parses custom extension blocks (callouts, stats, structured references), and caches processed entries in memory.
* `Database`: In-memory thread-safe repository enabling full-text search, tag aggregation, and multi-language routing (`en`, `id`, `es`).

### `internal/handler`
Contains HTTP presentation controllers, view rendering functions, and visitor identification.
* `Handler`: Central struct coordinating routers, stores, post databases, and layout rendering.
* `GetViewerKey`: Resolves unique visitor identity (`u:<id>` for members, `v:<uuid>` for anonymous readers via 10-year `dt_vid` cookie) and suppresses bot/crawler traffic.
* `Discussions*`: Forum listing with dynamic tag filters, topic view with view deduplication, reply submission, and upvoting endpoints.
* `Auth*`: OAuth initiation, callback verification, JSON data export (`/auth/export`), and account deletion (`/auth/delete-account`).
* `Blog*`: Post rendering, dynamic JSON-LD metadata attribution (`p.Author`), tag filtering, RSS 2.0, JSON Feed, and sitemap generation.
* `renderMarkdownPage`: Generic renderer for static markdown pages (`about`, `privacy`, `terms`, `accessibility`, `contribute`).

### `internal/router`
Constructs the Chi routing table, security middlewares, and localized path prefixes.
* Localized sub-routers for Indonesian (`/id/*`) and Spanish (`/es/*`) with fallback to root English paths.
* Security pipeline: `Recoverer`, `RealIP`, `Compress`, and `SessionMiddleware`.

### `internal/tui` and `internal/tuisrv`
Implements the terminal user interface and SSH daemon.
* Built using `charmbracelet/bubbletea` and `charmbracelet/wish`.
* Listens on `:2222` to serve an interactive keyboard-driven article reader directly to standard terminal clients without a browser.

### `internal/highlight`
Chroma-based syntax highlighting engine applied to server-rendered code snippets.

---

## 3. Database Schema and Storage Model

All persistent data is stored in a single SQLite database file configured with Write-Ahead Logging (`PRAGMA journal_mode=WAL`) and foreign key enforcement (`PRAGMA foreign_keys=ON`).

```sql
-- Users and authentication
CREATE TABLE users (
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

-- Active login sessions
CREATE TABLE sessions (
    token_hash TEXT PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL
);

-- Discussion topics
CREATE TABLE topics (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id       INTEGER NOT NULL,
    category      TEXT NOT NULL,
    title         TEXT NOT NULL,
    slug          TEXT NOT NULL UNIQUE,
    body_md       TEXT NOT NULL,
    body_html     TEXT NOT NULL,
    tags          TEXT NOT NULL DEFAULT '[]',
    views_count   INTEGER NOT NULL DEFAULT 0,
    replies_count INTEGER NOT NULL DEFAULT 0,
    votes_count   INTEGER NOT NULL DEFAULT 0,
    solved_reply_id INTEGER,
    pinned        INTEGER NOT NULL DEFAULT 0,
    created_at    DATETIME NOT NULL,
    updated_at    DATETIME NOT NULL
);

-- Discussion replies
CREATE TABLE replies (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    topic_id    INTEGER NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    user_id     INTEGER NOT NULL,
    parent_id   INTEGER REFERENCES replies(id) ON DELETE CASCADE,
    body_md     TEXT NOT NULL,
    body_html   TEXT NOT NULL,
    votes_count INTEGER NOT NULL DEFAULT 0,
    is_solution INTEGER NOT NULL DEFAULT 0,
    created_at  DATETIME NOT NULL,
    updated_at  DATETIME NOT NULL
);

-- Upvotes
CREATE TABLE topic_votes (
    topic_id   INTEGER NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    user_id    INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    PRIMARY KEY (topic_id, user_id)
);

-- Deduplicated unique topic views
CREATE TABLE forum_topic_views (
    topic_id   INTEGER NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    viewer_key TEXT NOT NULL,
    viewed_at  DATETIME NOT NULL,
    PRIMARY KEY (topic_id, viewer_key)
);

-- Deduplicated unique blog post views
CREATE TABLE post_views (
    post_slug  TEXT NOT NULL,
    viewer_key TEXT NOT NULL,
    viewed_at  DATETIME NOT NULL,
    PRIMARY KEY (post_slug, viewer_key)
);
```

---

## 4. Authentication and Request Flow

1. User clicks **Login with GitHub** (`/auth/github`).
2. Server generates a cryptographically secure `state` parameter and redirects the client to GitHub OAuth endpoint.
3. GitHub redirects back to `/auth/github/callback` with authorization `code` and `state`.
4. Server exchanges `code` for an access token, fetches the user profile from the GitHub API, and upserts the record in the `users` table using `provider_id`.
5. A random 256-bit token is generated, hashed with SHA-256, and stored in the `sessions` table.
6. The unhashed token is set in an `HttpOnly`, `Secure` (in production), `SameSite=Lax` cookie named `daemontalk_session`.
7. Subsequent requests pass through `SessionMiddleware`, which extracts the token, computes the SHA-256 hash, and attaches the active `auth.User` struct to the request context.

---

## 5. Privacy and Account Deletion Lifecycle

* **Data Export (`GET /auth/export`)**: Serializes all user profile fields, created topics, and submitted replies into a downloadable JSON payload.
* **Account Deletion (`POST /auth/delete-account`)**:
  1. Purges user session tokens from the `sessions` table.
  2. Updates forum topics and replies authored by the user: sets author display name to `Deleted User`, username to `ghost`, and avatar to `/static/images/deleted-user.png`.
  3. Anonymizes blog comments by clearing `user_id`, setting author to `Deleted User`, and removing verified badges.
  4. Deletes the record from the `users` table.
  5. Clears the session cookie and redirects the client.
