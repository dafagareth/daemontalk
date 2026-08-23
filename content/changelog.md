### v1.0.6 (August 24, 2026) · Complete Feature Eradication

- **TIL Complete Eradication**: Scoured the remaining codebase and completely eradicated leftover "Today I Learned" (TIL) dependencies, including obsolete struct models, localization dictionaries, and mock data in backend tests.
- **Author Parser Deprecation**: Nuked the redundant custom markdown parser for the ` ```author ` code block. The system no longer renders large author cards to keep the article reading experience as clean and focused as possible.

### v1.0.5 (August 23, 2026) · TUI LaTeX Graceful Degradation

- **Terminal UI Math Support**: Fixed an issue where LaTeX math blocks (`$$`) and inline math (`$`) would break terminal rendering by bleeding into markdown styling. The TUI now utilizes a preprocessing engine that elegantly degrades LaTeX math into syntax-highlighted code blocks, ensuring formulas remain perfectly readable and structurally intact over SSH connections.

### v1.0.4 (August 23, 2026) · Feature Pruning & Optimization

- **TIL Feature Deprecation**: Completely removed the "Today I Learned" (TIL) micro-blogging feature and its associated routes. This decision aligns with the publication's focus on high-quality, long-form technical articles, eliminating redundant UI clutter and streamlining the core architecture.
- **Guestbook UI Refinement**: Relocated the Guestbook submit button to the bottom-right corner of the input form to improve layout consistency and visual hierarchy.

### v1.0.3 (August 23, 2026) · Mobile Sidebar Redesign & UI Cleanup

- **Mobile Sidebar Overhaul**: Completely redesigned the mobile navigation drawer. Replaced the generic slide-in links with a dedicated top search bar, prominent topic streams, and native sans-serif fonts for better touch readability.
- **Clock Feature Deprecation**: Removed the redundant live clock indicator from both desktop and mobile navigation to maintain a cleaner, distraction-free header.
- **Social Ecosystem**: Added a dedicated, left-aligned social media icon row at the bottom of the mobile drawer. Introduced a Facebook icon and switched the Threads icon from a bulky solid fill to a clean, outlined SVG to match the platform's design language. Updated external handles to point directly to `daemontalk`.
- **Docs Standardization**: Eradicated redundant markup and "AI slop" across all `.md` documentation files. Converted messy HTML `details` blocks to the project's native ` ```faq ` syntax and scrubbed excessive bullet numbering from legal and contribution guides.

### v1.0.2 (August 22, 2026) · CI/CD Pipeline Patch

- **Deployment Synchronization**: Fixed a race condition in GitHub Actions where the VPS would deploy a stale image before the new GHCR Docker image finished building. Deployments now strictly wait for the build pipeline to complete.

### v1.0.1 (August 21, 2026) · UX & Infrastructure Refinements

- **Live Search**: Upgraded the search engine to use HTMX for real-time, instant search suggestions as you type. Fixed aggressive highlight bug for single-character queries.
- **Editor Auto-Save**: The Admin Web Studio now features robust `localStorage` auto-saving to prevent draft loss on accidental closure.
- **Typography Overrides**: Enforced rigid sans-serif fonts for markdown headings globally to maintain design hierarchy when the article body is toggled to Serif reading mode.
- **CSS Fix**: Resolved an inheritance bug where a global `[id]` selector inadvertently forced sans-serif typography onto the entire article body, overriding the reading mode toggle.
- **Canonical WWW Routing**: The apex domain (`daemontalk.com`) now strictly performs a 301 permanent redirect to `www.daemontalk.com` via Caddy.
- **Docs Update**: Simplified the project `README.md` and `VPS_SETUP.md` guides and swapped the ASCII art for a clean, scalable vector (SVG) architecture diagram.

### v1.0.0 (August 2026) · Initial Release

- **Core Architecture**: Single standalone Go binary powered by `chi` routing, `a-h/templ` server-side rendering, and Tailwind CSS v4. Zero heavy JS frameworks, zero trackers.
- **Bilingual Routing**: Native English and Indonesian (`/id`) localized routes with language-aware UI and metadata.
- **Interactive Terminal & CLI**: In-browser virtual UNIX shell (`/terminal`) with command history, plus curl-friendly endpoints (`/daily`, `/recipes`, `/p/:slug`).
- **Anonymous Comments & Guestbook**: Lightweight SQLite storage with persistent deterministic visitor handles (`anonym_<hex>`) and grouped conversational messaging.
- **Editorial Markdown Engine**: Server-rendered syntax highlighting, responsive carousels (` ```carousel `), galleries (` ```gallery `), accordions (` ```faq `), author cards (` ```author `), footnotes, and table of contents.
- **Reading Experience**: Font scaling (A+/A-), Serif toggle, multi-palette theme switcher (Light, Sepia, Dark), reading list bookmarks, and in-memory fulltext search.
- **Feeds & SEO**: Automated OpenGraph card generator, RSS 2.0 feed, JSON Feed, Sitemap XML, and hardened Content Security Policy (CSP).
