# daemontalk

`daemontalk.com` is a personal tech portal, systems engineering publication, and interactive UNIX playground built with Go, templ, HTMX, Tailwind CSS, and Charm Bubble Tea. It features an editorial tech-journalism layout, runs without a heavy client-side JavaScript framework, renders syntax-highlighted code on the server, and includes both an in-browser virtual UNIX terminal and a zero-install interactive TUI over SSH.

---

## ⚡ Quick Access & Terminal First

Access DaemonTalk directly from your favorite terminal without installing any dependencies:

```bash
# 1. Launch the interactive TUI over SSH (with 7 developer themes & keyboard navigation)
ssh daemontalk.com -p 2222

# 2. Stream the daily engineering dispatch briefing via curl
curl -sL daemontalk.com/daily

# 3. Read specific tags or systems cheat sheets
curl -sL daemontalk.com/t/linux
curl -sL daemontalk.com/recipes

# 4. Install the standalone CLI binary on your machine
curl -sL https://daemontalk.com/install.sh | bash
```

---

## ✨ Features

### 1. Terminal User Interface (TUI) & SSH Server
* **Interactive Bubble Tea Client**: Standalone terminal application (`cmd/tui`) with full keyboard navigation, split view, full-screen article reader mode (`Enter`), and instant browser action triggers (`o` to open cover image, `w` to open web article).
* **7 Developer Themes**: Instant theme cycling via `t` key:
  * ❄️ **Nord** *(Default)*
  * 🌃 **Tokyo Night**
  * 🍵 **Catppuccin Mocha**
  * 🪵 **Gruvbox Dark**
  * 🧛 **Dracula**
  * 🌲 **Rose Pine**
  * 🎨 **Monokai Pro**
* **Instant SSH Server**: Connect with `ssh daemontalk.com -p 2222` to launch the interactive TUI anywhere in the world.

### 2. Editorial Portal & Reading Experience
* **Bloomberg-Style Portal**: Clean typography, newsroom-style lead dispatches, trending reads driven by SQLite analytics, and category streams for Linux, Systems, Go, Rust, DevOps, Storage, and Security.
* **Interactive Systems Radar (`/graph`)**: Interactive canvas-based Knowledge Graph connecting Linux kernel architectures, language runtimes, memory safety, and storage engines.
* **Server-Side Syntax Highlighting**: Chroma syntax highlighting executed server-side to ensure zero client-side JavaScript execution overhead.

### 3. In-Browser Virtual UNIX Terminal (`/terminal`)
* Complete virtual file system (`/home/visitor`, `/etc`, `/bin`, `/var/log`, `/dev`).
* 40+ POSIX utilities: `cat`, `grep`, `awk`, `sed`, `find`, `curl`, `jq`, `top`, `ps`, `nc`, `dig`, `tree`, `env`.
* Full support for shell operators: pipes (`|`), redirection (`>`, `>>`), and chaining (`&&`, `||`, `;`).

### 4. Content Engine & Admin Studio (`/admin`)
* Author posts directly as Markdown files with YAML frontmatter in `content/posts/` or through the Web Studio editor at `/admin/posts/new`.
* **Weekly Digest Generator (`/admin/digest`)**: Auto-generates ready-to-publish Markdown newsletter summaries of recent dispatches.
* Pure Go SQLite backend in WAL mode using `modernc.org/sqlite`.

---

## 🏗️ Architecture

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

## 📁 Directory Structure

```
├── cmd/
│   └── tui/                 # Standalone TUI client binary entrypoint
├── content/
│   └── posts/               # Markdown technical dispatches with metadata
├── docs/
│   └── VPS_SETUP.md         # Production deployment guide (Debian + Docker + Caddy)
├── internal/
│   ├── comment/             # SQLite comments and pageview tracking
│   ├── handler/             # HTTP route handlers, SEO, RSS, sitemaps, graph, install
│   ├── highlight/           # Server-side Chroma syntax theme generator
│   ├── i18n/                # Bilingual string translations (EN / ID)
│   ├── og/                  # Dynamic OpenGraph PNG social card renderer
│   ├── post/                # Markdown file parser and post repository
│   ├── postdb/              # Web Studio SQLite persistence layer
│   ├── tui/                 # Shared Bubble Tea TUI components, themes, and views
│   └── tuisrv/              # Wish SSH server wrapper for remote TUI sessions
├── scripts/
│   └── setup-vps.sh         # 1-click VPS provisioner (Docker, Caddy, UFW)
├── web/
│   ├── static/
│   │   ├── css/             # Tailwind input.css, compiled main.css, chroma.css
│   │   ├── js/              # Terminal VFS engine and quill editor
│   │   └── images/          # Logos, icons, and asset diagrams
│   └── templates/           # Modular a-h/templ templates
├── Caddyfile                # Production reverse proxy and automatic SSL config
├── Dockerfile               # Multi-stage production container build
├── docker-compose.yml       # Production container orchestration
├── Makefile                 # Build, watch, test, and development tasks
└── main.go                  # Server entry point (HTTP + SSH TUI listener)
```

---

## 🚀 Getting Started

### Option 1: Run with Docker Compose (Recommended)

```bash
# 1. Clone repository
git clone https://github.com/dafagareth/daemontalk.git
cd daemontalk

# 2. Configure environment
cp .env.example .env

# 3. Start development stack (with live reload)
docker compose -f docker-compose.dev.yml up -d --build
```

Open [http://localhost:8080](http://localhost:8080) in your browser.

---

### Option 2: Local Host Development

**Prerequisites:**
- Go >= 1.23
- Node.js >= 20 & npm (for Tailwind CSS)
- templ: `go install github.com/a-h/templ/cmd/templ@latest`

```bash
# Install dependencies
npm install

# Run development server (runs Templ, Tailwind, and Go server in parallel)
make dev

# Run local TUI client
go run ./cmd/tui
```

---

## 🛠️ Development & Build Commands

```bash
# Generate Go code from .templ templates
make generate

# Compile Tailwind CSS bundle
make css

# Run unit tests
make test

# Build production binaries (server + TUI)
make build
go build -o bin/daemontalk-tui ./cmd/tui
```

---

## ⚙️ Environment Variables

| Variable | Description | Default |
| :--- | :--- | :--- |
| `PORT` | HTTP server listening port | `8080` |
| `SSH_PORT` | Wish SSH server port for remote TUI sessions | `2222` |
| `ADMIN_TOKEN` | Secret query token for `/admin?admin=<TOKEN>` access | `dev` |
| `BASE_URL` | Canonical base URL for feeds and OpenGraph cards | `https://daemontalk.com` |
| `ENV` | Environment mode (`development` / `production`) | `development` |

---

## 📄 License

Source code is available for personal and educational use under a Non-Commercial Source-Available License. The "daemontalk" brand, trademark, and original written articles are strictly reserved. See [LICENSE](LICENSE) for details.