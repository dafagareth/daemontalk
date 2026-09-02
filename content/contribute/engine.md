## Daemontalk Core Engine Development

Daemontalk is built from scratch as a high-performance modular monolith in Go with minimal dependencies.

### Architecture Overview

- **Backend & Routing**: Go 1.23+, Chi HTTP Router, pure-Go SQLite (`modernc.org/sqlite`), and Goldmark markdown parser.
- **Frontend & Rendering**: A-h Templ (type-safe templating), Tailwind CSS v4 CLI, and HTMX.
- **Terminal UI (TUI)**: Charmbracelet Bubble Tea, Lip Gloss styling, and SSH Daemon server (`/tuisrv`).
- **Auth & Discussions**: GitHub OAuth 2.0, HMAC-SHA256 session cookies, and discussions platform (`/socket`).

---

## Local Development Workflow

```bash
# 1. Clone repository and install dependencies
git clone https://github.com/dafagareth/daemontalk.git
cd daemontalk
go mod download
npm install

# 2. Compile templates & Tailwind CSS
make build

# 3. Run complete test suite
go test -count=1 ./...

# 4. Start local server (open http://localhost:8080)
./daemontalk
```

---

## Code Guidelines & PR Standards

1. **Formatting**: Ensure Go code is formatted with `go fmt ./...`.
2. **Testing**: Include unit tests for every new HTTP handler, storage query, or parser feature in `*_test.go`.
3. **Branching**: Use clear branch names such as `fix/bug-description` or `feat/feature-name`.
