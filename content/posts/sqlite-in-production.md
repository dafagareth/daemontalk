---
title: "SQLite in Production: More Capable Than You Think"
slug: 832ffc9a
aliases: [sqlite-in-production]
date: 2025-09-22
tags: [sqlite, database, backend]
lang: en
draft: false
---

SQLite has a reputation as "the embedded database for mobile apps and local tools." That reputation undersells it significantly. SQLite handles more production workloads than most developers realize, and the reasons to reach for Postgres or MySQL are more specific than "SQLite doesn't scale."

This is what you need to know before dismissing it.

## What SQLite Actually Is

SQLite is not a client-server database. There is no daemon, no connection pool, no network round-trip. The database is a single file your application reads and writes directly through a C library. That architecture has real implications, both good and bad.

Good: reads are extremely fast because there is no network overhead. A local read is a disk read.

Bad: write concurrency is limited. By default, only one writer at a time.

## WAL Mode Changes the Write Story

The default journal mode (DELETE) blocks all reads during a write. Enabling WAL (Write-Ahead Log) changes this:

```sql
PRAGMA journal_mode = WAL;
```

With WAL, readers never block writers, and writers never block readers. Multiple readers can run concurrently with one writer. This covers the majority of web application access patterns where reads vastly outnumber writes.

Set this once when you open the database. It persists.

Other pragmas worth setting at startup:

```sql
PRAGMA busy_timeout = 5000;    -- wait up to 5s instead of failing immediately on lock
PRAGMA synchronous = NORMAL;   -- safe with WAL, faster than FULL
PRAGMA foreign_keys = ON;
PRAGMA cache_size = -64000;    -- 64MB page cache
```

## What It Handles Well

SQLite performs well for applications with:

- Read-heavy workloads (analytics dashboards, blog backends, documentation sites)
- Low to moderate write rates (hundreds per second is achievable)
- A single application server
- Datasets up to tens of gigabytes

The SQLite documentation lists benchmarks showing it outperforms client-server databases on the same machine for workloads that fit this profile. The network round-trip that Postgres requires is not free.

## Real Numbers

A SQLite database with WAL enabled can typically handle:

- 50,000+ simple reads per second
- 1,000+ writes per second (with batching via transactions)

For a personal project, a startup with one server, or any application that has not yet proven it needs horizontal scaling, these numbers are not a bottleneck.

## When to Move to Postgres

The actual reasons to choose Postgres over SQLite:

- **Multiple application servers writing to the same database.** SQLite is a file. You cannot have two servers on different machines write to it simultaneously.
- **Write throughput above SQLite's ceiling.** If you are processing thousands of writes per second with no batching opportunity, you will hit limits.
- **Full-text search at scale.** SQLite has FTS5, but Postgres's full-text search and extensions like pg_trgm are more capable for large datasets.
- **Team expertise.** If your team knows Postgres deeply, the operational cost of managing SQLite is not worth the simplicity gain.

## The Deployment Advantage

A SQLite-backed application is one binary and one file. No separate database server to install, configure, monitor, or back up with special tooling. Backup is `cp database.db database.db.bak` or `sqlite3 database.db ".backup backup.db"`.

For personal projects and small production services, this matters. You can run a complete production stack on a $5 VPS with nothing but the application binary, the database file, and a systemd unit.

## Go and SQLite

If you are writing Go, the `modernc.org/sqlite` driver compiles SQLite into the binary with no CGO dependency:

```go
import (
    "database/sql"
    _ "modernc.org/sqlite"
)

db, err := sql.Open("sqlite", "./data/app.db")
```

Single binary, no external dependencies, production-ready. That is a compelling stack.
