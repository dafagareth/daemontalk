# Contributing to DaemonTalk

Guidelines for submitting technical articles, writing incident dispatches, and contributing code to DaemonTalk.

## Local Development Setup

Ensure you have Go (>= 1.25), Node.js (>= 20), `templ`, and `air` installed.

```bash
# Clone the repository
git clone https://github.com/dafagareth/daemontalk.git
cd daemontalk

# Install tools and dependencies
go install github.com/a-h/templ/cmd/templ@v0.3.1020
go install github.com/air-verse/air@latest
npm install

# Start the live-reloading development server (HTTP :8080 & SSH :2222)
make dev
```

## Writing Technical Articles

Articles live in `content/posts/` as standard Markdown files with YAML frontmatter.

**1. Generating a New Post**
Use the helper script to create a post with a random 8-character hex UID:

```bash
./scripts/post.sh uid "Your Article Title"
```

**2. Required Frontmatter Format**

```yaml
---
title: "Zero-Copy I/O with io_uring in Go"
slug: "7f8a9b1c"
date: "2026-08-19"
tags: ["linux", "go", "systems", "performance"]
lang: "en"
draft: false
type: "post"
summary: "Exploring asynchronous I/O batching and ring-buffer submissions using io_uring in Go."
---
```

**3. Writing Rules**
- Technical depth over length: prioritize reproducible shell outputs, benchmarks, architecture diagrams, and runnable code snippets.
- Use explicit language identifiers on fenced code blocks (e.g. ```` ```go ````, ```` ```bash ````, ```` ```c ````) for syntax highlighting.
- Keep illustrations in `web/static/images/posts/<slug>/` using modern `.webp` or `.png` format.

## Code Contributions

DaemonTalk follows a minimalist, zero-heavy-client-JS architecture:

- **Routing & HTTP:** `go-chi/chi/v5` in `internal/router/` and `internal/handler/`.
- **HTML Templates:** Type-safe AOT-compiled templates in `web/templates/*.templ`.
- **SSH Terminal Reader (TUI):** `charmbracelet/wish` and `bubbletea` in `internal/tui/` and `internal/tuisrv/`.
- **Storage:** Embedded SQLite with WAL mode in `internal/comment/` and `internal/postdb/`.

## Development & Makefile Reference

The project includes a comprehensive [Makefile](Makefile) covering all building, testing, code generation, and content workflows. Run or check the `Makefile` for all available targets:

- `make dev` — Start live-reloading dev environment (Go `air` + `templ` watch + Tailwind watch).
- `make build` — Full build pipeline (templ generate + minified CSS + JS bundle + binary).
- `make run` — Compile and execute the binary locally on port 8080.
- `make test` — Run all unit tests.
- `make test-race` — Run unit tests with the Go race detector enabled.
- `make test-cover` — Check test coverage percentages across all packages.
- `make fmt` & `make lint` — Format Go source files, templ components, and run `go vet`.
- `make new-post` / `make new-uid` — Scaffold new post files with metadata.
- `make validate-posts` — Validate Markdown frontmatter and tags across all posts.
- `make docker-build` / `make docker-up` — Build and run local Docker containers.
- `make clean` — Clean up compiled binaries and temporary test artifacts.

## Quality Checks & Testing

Always ensure tests pass and post frontmatter is valid before opening a pull request:

```bash
# Run unit tests via Makefile
make test

# Format code and templates
make fmt

# Validate Markdown posts
make validate-posts
```

## Pull Request Workflow

1. Fork the repository and create a branch from `main`.
2. Make your changes and write commit messages following the Conventional Commits specification (`feat:`, `fix:`, `docs:`).
3. Ensure all tests pass (`go test ./...`).
4. Push to your fork and submit a Pull Request against `main`.
