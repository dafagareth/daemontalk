## Overview

This site (**daemontalk.com**) is an independent technology publication and systems knowledge graph built entirely in Go. It explicitly avoids complex JavaScript frameworks and heavy build pipelines in favor of server-rendered HTML, precise HTMX interactions, and a lightweight SQLite database.

## Architecture

The system compiles down to a single standalone Go binary that independently serves:

- **Static Pages**: Core routing for home, about, uses, now, resume, changelog, and links.
- **Blog Engine**: Markdown parser with frontmatter metadata, server-side Chroma syntax highlighting, and dynamic table of contents.
- **Project Showcase**: Detailed markdown-driven project portfolios.
- **Community Guestbook**: Lightweight, persistent message board using cookie-bound anonymous identities.
- **Admin Studio**: Secure dashboard for moderation, content drafting, and telemetry analytics.

## Tech Stack

- **Go 1.26+** using the standard library and `go-chi/chi` for routing.
- **a-h/templ** for type-safe, compiled server-side HTML generation.
- **HTMX** for dynamic, real-time interactions (like live search) without JavaScript bloat.
- **Tailwind CSS v4** for strict, utility-first styling.
- **goldmark** for robust Markdown rendering and extensions.
- **modernc.org/sqlite** (pure-Go SQLite driver) for persistent data storage.
- Deployed on a lightweight Debian VPS managed by `systemd` and `Caddy`.

## Design Principles

- **Zero JS Dependency**: Core reading and navigation must function perfectly with JavaScript disabled.
- **Instant Cold Starts**: All posts and metadata are indexed in-memory on startup for sub-millisecond response times.
- **Operational Simplicity**: One binary, one SQLite database, one systemd service. No external caching layers or databases required.
- **Native Localization**: Full bilingual routing support for English (`/`) and Indonesian (`/id/`).

## Key Capabilities

- Granular typography control (Dark/Light themes, Serif/Sans toggle, dynamic font scaling).
- Automated OpenGraph preview card generation for social sharing.
- Comprehensive syndication via RSS 2.0 and JSON Feed.
- Integrated SEO essentials (Sitemap XML, `robots.txt`, structured schema).
- Rate limiting and honeypot mechanics for spam prevention.
- Scheduled publishing using `publish_at` frontmatter.
