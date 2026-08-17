# daemontalk

daemontalk is an engineering publication, systems knowledge graph, and interactive UNIX playground built with Go, templ, HTMX, Tailwind CSS, and Charm Bubble Tea. It delivers a fast, server-rendered editorial platform with zero client-side JavaScript framework overhead, syntax-highlighted code executed on the server, an in-browser virtual UNIX terminal, and a terminal user interface (TUI) accessible over SSH.

---

## Quick Access

Access daemontalk directly from the terminal without installing dependencies:

```bash
# Launch interactive TUI over SSH
ssh ssh.daemontalk.com -p 2222

# Stream daily engineering dispatches
curl -sL daemontalk.com/daily

# Read specific tag streams or cheat sheets
curl -sL daemontalk.com/t/linux
curl -sL daemontalk.com/recipes

# Install standalone CLI binary
curl -sL https://daemontalk.com/install.sh | bash
```

---

## Features

### Terminal User Interface (TUI) & SSH Server
- **Interactive Bubble Tea Client**: Standalone terminal application (`cmd/tui`) featuring split view, full-screen article reader mode (`Enter`), and instant browser triggers (`o` for image, `w` for web article).
- **7 Developer Color Themes**: Switchable at runtime with `t` key (Nord, Tokyo Night, Catppuccin Mocha, Gruvbox Dark, Dracula, Rose Pine, Monokai Pro).
- **Wish SSH Listener**: Direct remote access via SSH port `2222` with PTY negotiation and ANSI/Glamour markdown rendering.

### Editorial Web Platform
- **Content Engine**: Markdown posts with YAML frontmatter in `content/posts/` alongside a built-in SQLite Web Studio editor at `/admin/posts/new`.
- **Knowledge Graph (`/graph`)**: Canvas-based interactive systems radar connecting Linux kernel, eBPF, language runtimes, memory safety, and storage architectures.
- **Server-Side Rendering**: High-performance HTML compilation via `a-h/templ` and server-side Chroma code highlighting.
- **Weekly Digest Generator (`/admin/digest`)**: Automated Markdown newsletter compiler for editorial distribution.

### Virtual UNIX Terminal (`/terminal`)
- In-browser virtual file system (`/home/visitor`, `/etc`, `/bin`, `/var/log`, `/dev`).
- Over 40 POSIX-compliant utilities including `cat`, `grep`, `awk`, `sed`, `find`, `curl`, `jq`, `top`, `ps`, `nc`, `dig`, and `tree`.
- Support for pipelines (`|`), redirection (`>`, `>>`), and chaining (`&&`, `||`, `;`).

---

## Architecture

```
                                  +-------------------------+
                                  |  Terminal / SSH Client  |
                                  +------------+------------+
                                               | (port 2222)
                                               v
                                  +-------------------------+
                                  |   Wish SSH TUI Server   |
                                  | (Bubble Tea, Glamour)   |
                                  +-------------------------+

                                  +-------------------------+
                                  |   Web Browser / cURL    |
                                  +------------+------------+
                                               | (port 80/443 -> 8080)
                                               v
                                  +-------------------------+
                                  |    Go Chi Web Server    |
                                  |   (Middlewares, SEO)    |
                                  +-----+-------------+-----+
                                        |             |
                     +------------------+             +------------------+
                     v                                                   v
         +-----------------------+                           +-----------------------+
         | Content Engine (Disk) |                           | SQLite Engine (WAL)   |
         | - content/posts/*.md  |                           | - Page Views / Hits   |
         | - Goldmark Parser     |                           | - Comments & Reactions|
         | - Chroma Highlighting |                           | - Web Studio Posts    |
         +-----------+-----------+                           +-----------+-----------+
                     |                                                   |
                     +------------------+             +------------------+
                                        |             |
                                        v             v
                                  +-------------------------+
                                  |  a-h/templ Components   |
                                  |  (SSR HTML + HTMX + CSS)|
                                  +-------------------------+
```

---

## Project Structure

```
├── cmd/
│   └── tui/                 # Standalone TUI client binary
├── content/
│   └── posts/               # Markdown technical dispatches with metadata
├── docs/
│   └── VPS_SETUP.md         # Production deployment guide (Debian, Docker, Caddy)
├── internal/
│   ├── comment/             # SQLite comment and view count repository
│   ├── handler/             # Route handlers, SEO, RSS, sitemaps, graph, install
│   ├── highlight/           # Server-side Chroma syntax highlighter
│   ├── i18n/                # Bilingual translations (EN / ID)
│   ├── og/                  # Dynamic OpenGraph PNG generator
│   ├── post/                # Markdown file parser and post loader
│   ├── postdb/              # Web Studio SQLite persistence layer
│   ├── tui/                 # Shared Bubble Tea TUI components, themes, and views
│   └── tuisrv/              # Wish SSH server wrapper for remote TUI sessions
├── scripts/
│   └── setup-vps.sh         # Automated VPS provisioner (Docker, Caddy, UFW)
├── web/
│   ├── static/
│   │   ├── css/             # Tailwind input.css, compiled main.css, chroma.css
│   │   ├── js/              # Terminal VFS engine and quill editor
│   │   └── images/          # Assets and diagrams
│   └── templates/           # Modular a-h/templ component templates
├── Caddyfile                # Reverse proxy and TLS certificate automation
├── Dockerfile               # Multi-stage production container build
├── docker-compose.yml       # Production container orchestration
├── Makefile                 # Build, watch, test, and development automation
└── main.go                  # Application entry point (HTTP & SSH servers)
```

---

## Getting Started

### Using Docker Compose (Recommended)

```bash
# Clone repository
git clone https://github.com/dafagareth/daemontalk.git
cd daemontalk

# Configure environment
cp .env.example .env

# Start development stack
docker compose -f docker-compose.dev.yml up -d --build
```

Access the web interface at `http://localhost:8080`.

---

### Local Bare-Metal Setup

**Prerequisites:**
- Go >= 1.23
- Node.js >= 20 & npm
- templ CLI: `go install github.com/a-h/templ/cmd/templ@latest`

```bash
# Install frontend assets
npm install

# Run development server (Templ, Tailwind, Air live-reload)
make dev

# Run local TUI client
go run ./cmd/tui
```

---

## Build Commands

```bash
# Generate Go code from .templ files
make generate

# Compile Tailwind CSS bundle
make css

# Run unit tests
make test

# Build production binaries
make build
go build -o bin/daemontalk-tui ./cmd/tui
```

---

## Configuration

| Variable | Description | Default |
| :--- | :--- | :--- |
| `PORT` | HTTP server listening port | `8080` |
| `SSH_PORT` | Wish SSH server port for remote TUI sessions | `2222` |
| `ADMIN_TOKEN` | Secret query token for `/admin?admin=<TOKEN>` access | `dev` |
| `BASE_URL` | Canonical base URL for feeds and OpenGraph metadata | `https://daemontalk.com` |
| `ENV` | Environment mode (`development` / `production`) | `development` |

---

## License

Source code is released under the terms of the Non-Commercial Source-Available License. The "daemontalk" name, trademark, and original written content are reserved. See [LICENSE](LICENSE) for details.
