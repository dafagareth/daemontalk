# Architecture

Daemontalk is a single Go binary combining server-side rendering, embedded SQLite, and static assets with zero client JS frameworks.

## Subsystems

- **Web Server (`internal/router`)**: Chi router with security headers, rate limiting, gzip compression, and session middleware.
- **Rendering (`web/templates/`)**: Go `templ` components compiled to Go code, styled with standalone Tailwind CSS v4.
- **Content Engine (`internal/post`)**: In-memory index of markdown articles in `content/posts/` with YAML frontmatter, Chroma syntax highlighting, reading time calculation, and tag filtering.
- **Forum / Socket (`internal/forum`)**: Discussions and Q&A engine (`/socket`) with upvoting, solution marking, and view tracking.
- **Comments (`internal/comment`)**: Hierarchical article comments with visitor attribution and SSE live streaming.
- **Auth (`internal/auth`)**: GitHub OAuth 2.0 with cryptographically hashed session tokens in SQLite.

## Data Storage

SQLite in WAL mode (`PRAGMA journal_mode=WAL`) with foreign key enforcement:

| Table | Purpose | Key Fields |
|-------|---------|------------|
| `users` | Member accounts from GitHub OAuth | `id`, `provider_id`, `username`, `role` |
| `sessions` | Active login sessions | `token_hash` (SHA-256), `user_id`, `expires_at` |
| `topics` | Discussion threads (`/socket`) | `id`, `user_id`, `slug`, `solved_reply_id` |
| `replies` | Threaded discussion replies | `id`, `topic_id`, `user_id`, `parent_id`, `is_solution` |
| `topic_votes` | Prevent duplicate upvotes | `(topic_id, user_id)` PK |
| `forum_topic_views` | Deduplicated topic view counts | `(topic_id, viewer_key)` PK |
| `post_views` | Deduplicated article view counts | `(post_slug, viewer_key)` PK |
| `comments` | Article comments & threaded replies | `id`, `post_slug`, `user_id`, `parent_id` |

## Request Flow

1. **Routing**: Chi routes requests through middleware (`RealIP`, `Recoverer`, `SecurityHeaders`, `RateLimiter`, `SessionMiddleware`).
2. **Visitor Tracking**: `GetViewerKey` generates a persistent identifier (`u:<id>` for members, `v:<uuid>` via cookie) to deduplicate views and prevent bot skew.
3. **Rendering**: Handlers query SQLite or the in-memory post index, passing data to `templ` components directly into the `http.ResponseWriter`.

## Auth & Privacy

- **Login**: GitHub OAuth 2.0 (`/auth/github` -> `/auth/github/callback`) -> upsert user -> generate 256-bit token -> store SHA-256 hash -> set `HttpOnly`, `SameSite=Lax` cookie.
- **Export (`GET /auth/export`)**: Serializes user profile, topics, and replies into JSON.
- **Delete Account (`POST /auth/delete-account`)**: Anonymizes user contributions (author to `[Deleted User]`, username to `ghost`), removes sessions, deletes user row, and clears cookie.
