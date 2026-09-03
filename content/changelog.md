### v1.4.0 (September 3, 2026) · Fast-Track Sync, Cloudflare CSP & UI Polish

- **Fast-Track Content Sync**: Implemented a bypass workflow (`sync-content.yml`) that triggers via webhook for markdown-only changes, skipping the 3-minute Docker build and deploying instantly (~10s).
- **Git Log Date Fallback & UTC Fix**: Fixed git history resolution in Docker (`safe.directory` bypass) and synchronized container timezone to `Asia/Jakarta` to ensure accurate article publication dates.
- **Cloudflare Analytics Integration**: Hardened the Content-Security-Policy (CSP) by explicitly whitelisting Cloudflare's Web Analytics beacon in the router middleware.
- **UI Enhancements**: Implemented an auto-dimming Focus Mode (`opacity-25`) for desktop sidebars and replaced the static Table of Contents with a collapsible `<details>` element for cleaner reading.

### v1.3.2 (September 2, 2026) · CI/CD Deploy Script Execution Fix

- **VPS Deployment Permissions**: Added automatic `chmod +x` and explicit `bash` invocation in GitHub Actions CD deployment workflow to prevent execution permission errors on remote hosts.

### v1.3.1 (September 2, 2026) · Webhook Method Handling & Tag Push Support

- **Webhook GET & Method Resilience**: Handled `GET` requests gracefully on `/api/webhook/github` and configured Caddy with HTTP 308 permanent redirect to preserve HTTP methods across reverse proxies.
- **Git Tag Push Support**: Extended the GitHub webhook listener to trigger automatic article reloads on git release tag pushes (`refs/tags/*`) in addition to branch pushes.

### v1.3.0 (September 2, 2026) · GitHub OAuth2, Socket Discussions & Modular Licensing

- **GitHub OAuth2 & User Profiles**: Integrated seamless GitHub authentication (`/auth/github`), user session persistence via SQLite, dynamic navbar badge, and clean user profile pages (`/u/:username`).
- **Socket Discussions Forum**: Built a native, server-rendered community discussion platform (`/socket`) featuring topic categorization, threaded replies, upvoting, and solution validation.
- **PolyForm Noncommercial & CC BY-NC-SA 4.0 Licensing**: Standardized codebase under the PolyForm Noncommercial License 1.0.0, editorial dispatches under CC BY-NC-SA 4.0, and strict All Rights Reserved brand protection.
- **Git-Driven Instant Publishing**: Implemented HMAC-SHA256 verified GitHub webhook (`/api/webhook/github`) for hot-reloading markdown dispatches automatically on push with zero downtime.
- **Tag Portal & UI Redesign**: Overhauled tag archive portals (`/blog/tag/:tag`) with a modern 4-column framed grid, streamlined mobile navigation drawer, and modernized the Contributor Directory.
- **Legacy Architecture Pruning**: Completely eradicated legacy modules including `/links`, web `/terminal`, `/guestbook`, and obsolete runner engines.
- **Colophon Specification**: Added a dedicated system architecture and infrastructure blueprint page (`/colophon`).

### v1.2.0 (August 28, 2026) · Security Hardening, Full i18n & Architecture Refinements

- **Security & Authorization Hardening**: Strengthened access boundaries, enforced strict URL validation to prevent XSS, added payload size limits, and patched potential draft content leaks across syndication and CLI endpoints.
- **Universal Bilingual Experience**: Completed full localized UI dictionary coverage across all modals, search dropdowns, footer links, and guestbook components for English and Indonesian (`/id`).
- **Pagination & Navigation Fixes**: Resolved a pagination gap where articles were skipped on "Load More", and fixed alias redirects to preserve the active language route.
- **Streamlined Archive River**: Simplified the Chronological Archive River list to focus strictly on headlines, thumbnails, and metadata without description clutter.
- **Palette Simplification**: Deprecated legacy Sepia mode to deliver refined Light and Dark themes with dedicated code syntax styling.
- **SSH TUI Concurrency & Resilience**: Isolated themes and states per SSH connection session, eliminated server-crashing error exits, and optimized file lookups to O(1).

### v1.1.0 (August 27, 2026) · Mobile UI & Codebase Optimization

- **Mobile View Optimizations**: Enforced strict edge-to-edge (padding-free) constraints and consistent 16:9 aspect ratios for thumbnail images on mobile viewports across the portal.
- **Reading List Redesign**: Optimized the 'Saved Dispatches' view by converting it into a clean, minimalist single-column ledger list on mobile devices.
- **Navigation Cleanup**: Streamlined the mobile navigation sidebar, added missing essential links, updated labels, and removed unnecessary social icons to prioritize readability.
- **Codebase Refactoring**: Executed a comprehensive template cleanup, stripping away excessive styling annotations, legacy comments, and "AI slop" from the HTML structure, resulting in a significantly cleaner developer experience.
- **Search Bug Fix**: Resolved a fallback rendering bug where submitting a search query on the HTMX form resulted in a partial, unstyled HTML fragment being served instead of the complete page layout.

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

- **Deployment Synchronization**: Fixed a race condition in automated CI/CD deployment pipelines where the VPS would deploy a stale container image before the new build finished. Deployments now strictly wait for the build pipeline to complete.

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
- **Reading Experience**: Font scaling (A+/A-), Serif toggle, theme switcher (Light, Dark), reading list bookmarks, and in-memory fulltext search.
- **Feeds & SEO**: Automated OpenGraph card generator, RSS 2.0 feed, JSON Feed, Sitemap XML, and hardened Content Security Policy (CSP).

### v1.3.3 (2026-09-02)
- **Fix**: Implemented Git Log history fallback for article dates (posts no longer require frontmatter `date`).
- **Fix**: Removed `date:` field from `new-post` CLI template.
- **CI/CD**: Added Fast-Track content sync via `.github/workflows/sync-content.yml`.
- **CI/CD**: Ignored `content/**` changes in Docker builds to speed up publishing from 5 minutes to 5 seconds.
- **Chore**: Updated OG default meta title and JSON Feed description to "Technology publication & community portal".
