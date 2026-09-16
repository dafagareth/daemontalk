# Development

## Prerequisites

- Go 1.26+
- `templ` CLI: `go install github.com/a-h/templ/cmd/templ@latest`
- `air` (optional, for live reload): `go install github.com/air-verse/air@latest`

Tailwind CSS v4 standalone binary is automatically fetched by `Makefile` if not present. No Node.js required.

## Local Workflow

```bash
# Live reloading (Air + templ watch + Tailwind watch)
make dev

# Full build
make build
./daemontalk
```

## Common Commands

- `make dev` — Run live-reload development server on `:8080`
- `make build` — Compile templates, Tailwind CSS, and binary
- `make test` — Run all unit tests (`go test ./...`)
- `make test-race` — Run tests with race detector (`go test -race ./...`)
- `make lint` — Run `go fmt` and `go vet` checks
- `make validate-posts` — Verify markdown frontmatter, dates, and images

## Conventions

- **Rendering**: Go `templ` components under `web/templates/`. Zero client-side JS frameworks.
- **Data Access**: Parameterized queries only in `internal/<module>/store*.go`.
- **Auth**: Always resolve active user from `auth.GetUser(r.Context())`.
- **Styling**: Tailwind CSS v4 utility classes. Keep custom CSS in `web/static/css/input.css` to a minimum.
