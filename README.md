# daemontalk

[![CI](https://github.com/dafagareth/daemontalk/actions/workflows/ci.yml/badge.svg)](https://github.com/dafagareth/daemontalk/actions/workflows/ci.yml)
[![Deploy](https://github.com/dafagareth/daemontalk/actions/workflows/deploy.yml/badge.svg)](https://github.com/dafagareth/daemontalk/actions/workflows/deploy.yml)
[![License](https://img.shields.io/badge/License-Non--Commercial-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8.svg)](go.mod)

`daemontalk` is an independent tech publication, digital notebook, and open computing platform for developers and enthusiasts. Built as a single standalone Go binary, it provides deep-dive software articles, community discussions, an in-browser web terminal, and an interactive terminal reader accessible over SSH.

<p align="center">
  <img src="docs/architecture.svg" alt="Daemontalk Architecture" width="100%">
</p>

## Primary Capabilities

- Interactive Bubble Tea terminal client accessible directly via SSH (`ssh ssh.daemontalk.com -p 2222`).
- In-browser virtual UNIX shell (`/terminal`) with a simulated POSIX file system and 40+ commands.
- High-performance SSR blogging engine with in-memory full-text search and zero client JS frameworks.
- Community discussions forum (`/discussions`) with GitHub OAuth 2.0, voting, and self-service data privacy.
- Anonymous and verified SQLite-backed comment engine with persistent visitor handles.
- Automated OpenGraph image generator, RSS 2.0 feed, JSON Feed, and XML Sitemaps.

## Quick Start

Launch the interactive TUI over SSH:
```bash
$ ssh ssh.daemontalk.com -p 2222
```

Stream daily tech dispatches via curl:
```bash
$ curl -sL daemontalk.com/daily
```

Install the standalone CLI binary locally:
```bash
$ curl -sL https://daemontalk.com/install.sh | bash
```

## Development & Makefile Workflows

All development workflows, testing, code generation, container orchestration, and content utilities are driven via the [Makefile](Makefile). Check the `Makefile` for the complete set of targets:

- `make dev` — Start live-reloading dev environment (Go `air` + `templ` watch + Tailwind watch).
- `make build` — Full build pipeline: compile templates, minify Tailwind CSS, and build binary.
- `make test` / `make test-race` — Run unit tests and race detection suite.
- `make new-post` / `make validate-posts` — CLI utilities for managing and validating markdown articles.
- `make docker-up` / `make docker-down` — Production Docker Compose lifecycle.

## Documentation and Specifications

- Refer to [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for internal codebase architecture, package breakdowns, and database schemas.
- Refer to [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) for local development conventions.
- Refer to [CONTRIBUTING.md](CONTRIBUTING.md) for article writing guidelines and Pull Request workflows.

## License

Source code is released under the Non-Commercial Source-Available License. The "daemontalk" name and original written content are reserved. See [LICENSE](LICENSE) for terms.
