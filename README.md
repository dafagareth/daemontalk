# daemontalk

[![CI](https://github.com/dafagareth/daemontalk/actions/workflows/ci.yml/badge.svg)](https://github.com/dafagareth/daemontalk/actions/workflows/ci.yml)
[![Deploy](https://github.com/dafagareth/daemontalk/actions/workflows/deploy.yml/badge.svg)](https://github.com/dafagareth/daemontalk/actions/workflows/deploy.yml)
[![License](https://img.shields.io/badge/License-Non--Commercial-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8.svg)](go.mod)

`daemontalk` is an independent systems engineering publication, research notebook, and source-available computing platform. Built as a single standalone Go binary, it provides deep-dive software engineering articles, community discussions, an in-browser web terminal, and an interactive terminal reader accessible over SSH.

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

Stream daily engineering dispatches via curl:
```bash
$ curl -sL daemontalk.com/daily
```

Install the standalone CLI binary locally:
```bash
$ curl -sL https://daemontalk.com/install.sh | bash
```

## Documentation and Specifications

Refer to [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for internal codebase architecture, package breakdowns, and database schemas.

Refer to [docs/VPS_SETUP.md](docs/VPS_SETUP.md) for production deployment instructions using Docker, Caddy, and UFW.

Refer to [CONTRIBUTING.md](CONTRIBUTING.md) for development setup (Docker vs local) and contribution standards.

## License

Source code is released under the Non-Commercial Source-Available License. The "daemontalk" name and original written content are reserved. See [LICENSE](LICENSE) for terms.
