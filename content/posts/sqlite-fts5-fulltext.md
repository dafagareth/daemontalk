---
title: "Full-Text Search in SQLite with FTS5"
slug: c2a08b5e
aliases: sqlite-fts5-fulltext
date: 2026-04-14
tags: [storage, backend, tools]
lang: en
draft: false
---

SQLite includes a full-text search extension called FTS5 that supports ranked results, prefix queries, and phrase matching. It is compiled into most SQLite distributions by default and requires no additional dependencies.

## Fun Facts

**Fact 1.** FTS5 is the fifth iteration of SQLite's full-text search module. Its predecessor, FTS4, is still supported but no longer receives new features. FTS5 uses a significantly redesigned index format that enables pluggable tokenizers and the BM25 ranking function natively.

**Fact 2.** SQLite's FTS5 stores its inverted index across five internal shadow tables (e.g., `<name>_data`, `<name>_idx`, `<name>_docsize`). These are regular SQLite tables and are visible in `.tables` output, but modifying them directly corrupts the index.

**Fact 3.** BM25 (Best Match 25) was published by Robertson and Walker in 1994 and remains one of the most effective ranking algorithms for keyword search. Elasticsearch uses it as its default similarity model.

---

## Tips and Tricks

### 1. FTS5 vs LIKE: When to Use Which

`LIKE` scans every row and applies a pattern match. It is acceptable for small tables or when the pattern has a non-wildcard prefix (allowing index use). FTS5 maintains an inverted index and scales to millions of rows without full-table scans.

Use `LIKE` when:
- The table has fewer than ~10,000 rows and no ranking is needed.
- The query is a simple prefix match on an indexed column.

Use FTS5 when:
- You need relevance ranking.
- Users can input arbitrary multi-word queries.
- The text corpus is large.

```sql
-- LIKE: no ranking, sequential scan on large tables
SELECT title FROM articles WHERE body LIKE '%kernel namespace%';

-- FTS5: ranked, uses inverted index
SELECT title FROM articles_fts WHERE articles_fts MATCH 'kernel namespace';
```

### 2. Creating a Virtual FTS5 Table

An FTS5 virtual table indexes one or more text columns. The `content` option links it to an existing base table, avoiding data duplication.

```sql
-- Create a base table
CREATE TABLE articles (
  id    INTEGER PRIMARY KEY,
  title TEXT NOT NULL,
  body  TEXT NOT NULL
);

-- Create a content-linked FTS5 index
CREATE VIRTUAL TABLE articles_fts USING fts5(
  title,
  body,
  content='articles',
  content_rowid='id'
);

-- Populate the FTS index (required after linking to content table)
INSERT INTO articles_fts(articles_fts) VALUES ('rebuild');
```

After the initial rebuild, keep the index current with triggers:

```sql
CREATE TRIGGER articles_ai AFTER INSERT ON articles BEGIN
  INSERT INTO articles_fts(rowid, title, body)
  VALUES (new.id, new.title, new.body);
END;

CREATE TRIGGER articles_ad AFTER DELETE ON articles BEGIN
  INSERT INTO articles_fts(articles_fts, rowid, title, body)
  VALUES ('delete', old.id, old.title, old.body);
END;

CREATE TRIGGER articles_au AFTER UPDATE ON articles BEGIN
  INSERT INTO articles_fts(articles_fts, rowid, title, body)
  VALUES ('delete', old.id, old.title, old.body);
  INSERT INTO articles_fts(rowid, title, body)
  VALUES (new.id, new.title, new.body);
END;
```

### 3. Ranking with BM25

FTS5 exposes BM25 as the `rank` column. Lower (more negative) values indicate higher relevance. Always sort ascending.

```sql
-- Search and rank by relevance
SELECT
  a.id,
  a.title,
  articles_fts.rank
FROM articles_fts
JOIN articles a ON a.id = articles_fts.rowid
WHERE articles_fts MATCH 'linux container'
ORDER BY articles_fts.rank
LIMIT 10;
```

To weight columns differently (title matches count more than body matches):

```sql
-- Column weights: title = 10.0, body = 1.0
SELECT title, rank
FROM articles_fts
WHERE articles_fts MATCH 'linux container'
ORDER BY bm25(articles_fts, 10.0, 1.0)
LIMIT 10;
```

### 4. Highlighting and Snippets

The `highlight()` and `snippet()` auxiliary functions help build search UIs without post-processing in application code.

```sql
-- Wrap matched tokens in <b> tags
SELECT
  highlight(articles_fts, 0, '<b>', '</b>') AS title_hl,
  highlight(articles_fts, 1, '<b>', '</b>') AS body_hl
FROM articles_fts
WHERE articles_fts MATCH 'kernel namespace'
ORDER BY rank
LIMIT 5;

-- Return a short excerpt around the first match
SELECT
  a.title,
  snippet(articles_fts, 1, '<mark>', '</mark>', '...', 16) AS excerpt
FROM articles_fts
JOIN articles a ON a.id = articles_fts.rowid
WHERE articles_fts MATCH 'cgroup isolation'
ORDER BY rank
LIMIT 5;
```

The fourth argument to `snippet()` is the maximum number of tokens in the returned fragment. Keep it between 10 and 32 for readable excerpts.

### 5. Performance Limits to Keep in Mind

FTS5 is well-suited to read-heavy workloads with infrequent bulk writes. There are several limits worth knowing before committing to it:

- **Write amplification.** Each INSERT or DELETE to a content table with triggers causes multiple writes to the shadow tables. On write-heavy workloads, this adds up.
- **No partial updates.** Updating a single field requires deleting and re-inserting the full document into the FTS index.
- **Index size.** FTS5 indexes are typically 1.5x to 3x the size of the raw text, depending on tokenization and corpus vocabulary.
- **No fuzzy matching.** FTS5 supports prefix queries (`linux*`) but not edit-distance fuzzy matching. For that, you need a separate approach such as trigram indexing or a dedicated search engine.
- **WAL mode recommended.** Enable WAL for better read/write concurrency when FTS queries and inserts happen concurrently.

```bash
# Check FTS5 index size on disk
sqlite3 mydb.sqlite "SELECT SUM(pgsize) FROM dbstat WHERE name LIKE 'articles_fts%';"
```

For datasets under a few million short documents, FTS5 is a practical choice that avoids the operational overhead of running a separate search service.
