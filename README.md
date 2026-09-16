# daemontalk

[![CI](https://github.com/dafagareth/daemontalk/actions/workflows/ci.yml/badge.svg)](https://github.com/dafagareth/daemontalk/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-Non--Commercial-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8.svg)](go.mod)

Independent systems engineering journal and discussion platform. Single standalone Go binary with SQLite, pure SSR via `templ`, standalone Tailwind CSS, and zero client JS frameworks.

## Quick Start

```bash
# Start local dev server (Air + templ + Tailwind)
make dev

# Or build standalone binary
make build
./daemontalk
```

## Useful Commands

- `make dev` Run dev server with hot reload
- `make build` Compile templates, Tailwind CSS, and binary
- `make test` Run test suite
- `make lint` Run format and vet checks (`fmt` + `vet`)
- `make validate-posts` Validate article frontmatter and links

## Documentation

- [Architecture](docs/ARCHITECTURE.md) — System design and schemas
- [Contributing](CONTRIBUTING.md) — Editorial guidelines and code contributions
- [License](LICENSE) — Non-Commercial Source-Available License
