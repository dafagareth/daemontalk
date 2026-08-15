### v1.0.0 (August 2026) · Initial Release

*   **Core Architecture**: Single standalone Go binary powered by `chi` routing, `a-h/templ` server-side rendering, and Tailwind CSS v4. Zero heavy JS frameworks, zero trackers.
*   **Bilingual Routing**: Native English and Indonesian (`/id`) localized routes with language-aware UI and metadata.
*   **Interactive Terminal & CLI**: In-browser virtual UNIX shell (`/terminal`) with command history, plus curl-friendly endpoints (`/daily`, `/recipes`, `/p/:slug`).
*   **Anonymous Comments & Guestbook**: Lightweight SQLite storage with persistent deterministic visitor handles (`anonym_<hex>`) and grouped conversational messaging.
*   **Editorial Markdown Engine**: Server-rendered syntax highlighting, responsive carousels (` ```carousel `), galleries (` ```gallery `), accordions (` ```faq `), author cards (` ```author `), footnotes, and table of contents.
*   **Reading Experience**: Font scaling (A+/A-), Serif toggle, multi-palette theme switcher (Light, Sepia, Dark), reading list bookmarks, and in-memory fulltext search.
*   **Feeds & SEO**: Automated OpenGraph card generator, RSS 2.0 feed, JSON Feed, Sitemap XML, and hardened Content Security Policy (CSP).
