## Overview

This site (**daemontalk.com**) is a personal portfolio and tech blog built with Go. No complex JavaScript frameworks or heavy build pipelines. Just server-rendered HTML, a sprinkle of HTMX for interactive features, and a SQLite database for comments and analytics.

## Architecture

The site is a single Go binary that serves:

- **Static pages**: home, about, uses, now, resume, changelog, links
- **Blog**: markdown posts with frontmatter, syntax highlighting, table of contents, and series navigation
- **Projects**: project showcase with detail pages
- **Guestbook**: community message board
- **TIL**: bite-sized engineering discoveries
- **Admin**: lightweight dashboard for moderation and basic analytics

## Tech Stack

- **Go** with the standard library + Chi router
- **templ** for type-safe server-side HTML templating
- **HTMX** for dynamic interactions without heavy JS bundles
- **Tailwind CSS v4** for clean utility styling
- **goldmark** for Markdown rendering with Chroma syntax highlighting
- **modernc.org/sqlite** (pure-Go SQLite driver) for persistence
- Deployed on a single Linux VPS with systemd

## Design Goals

- **Zero JavaScript framework dependencies**: every page works reliably without client JS.
- **Fast cold start**: all posts are indexed in memory on startup for instant response times.
- **Simple ops**: one binary, one SQLite database, one systemd service.
- **Bilingual**: full support for English (`/`) and Indonesian (`/id/`).

## Features Shipped

- Dark/light/system theme toggle
- Per-post OG image generation
- RSS 2.0 + JSON Feed
- Sitemap + robots.txt
- PWA manifest
- Rate limiting + honeypot spam protection
- Scheduled posts via `publish_at` frontmatter
- Emoji reactions (👍 ❤️ 🔥)
- Reading time estimate and view counter
- Print styles and anchor copy links
- Keyboard shortcuts (press `?`)
