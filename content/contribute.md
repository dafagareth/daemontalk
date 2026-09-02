## Philosophy & Editorial Standards

Daemontalk is an independent systems engineering publication dedicated to practical, verifiable, and reproducible technical dispatches. We welcome contributions from software engineers, systems researchers, infrastructure architects, and open-source developers passionate about low-level computing and platform engineering.

**Straight to the Point**: Dive directly into the technical core, architectural blueprints, or code examples in the opening section. Avoid conversational fluff or generic introductory filler that adds no engineering value.

**Verification & Reproducibility**: Every technical claim or performance benchmark must be supported by actionable test code, diagnostic logs, shell commands, or architecture diagrams that enable readers to independently reproduce and verify your findings.

**Cited & Authoritative References**: Conclude every in-depth dispatch with verified citations to authoritative specifications, RFC standards, Linux kernel source code repositories, processor manuals, or academic research papers.

---

## Accepted Contribution Pillars

Daemontalk is an open-source platform open to multiple forms of technical contribution:

**1. Writing Technical Dispatches (`content/posts/`)**:
Authoring in-depth analyses, kernel investigations, network protocol benchmarks, database internals, and systems post-mortems.

**2. Codebase Improvements & Bug Fixes (*Core Engine*)**:
- Go Backend (HTTP handlers, SQLite storage, CLI tools, Goldmark markdown extensions).
- Terminal UI (*Bubble Tea / Lip Gloss*).
- Templ components & Tailwind CSS styling.
- Performance profiling, security hardening, and database query optimizations.

**3. Content Corrections, Typos & Code Snippet Updates**:
Found a bug in a code example, an outdated shell flag, broken references, or a typo in a published article? Submit a Pull Request directly against the corresponding `.md` file.

**4. Translations & Localization (*i18n*)**:
Help localize articles and user interface strings into English (`.md`), Indonesian (`.id.md`), or Spanish (`.es.md`).

---

## Topics & Domains of Interest

We actively seek dispatches across the following core systems domains:

**Operating Systems & Linux Kernel**: eBPF and XDP packet processing, cgroups v2 resource isolation, CPU scheduling algorithms (EEVDF), memory management internals (Linux OOM Killer mechanics), system call profiling, and zero-copy I/O paths (`sendfile`, `io_uring`).

**Concurrency & Language Runtimes**: Go runtime internals (Tri-color Garbage Collector, goroutine scheduler), Rust memory management (*borrow checker*, *lifetimes*), and lock-free atomic data structures powered by Compare-And-Swap (CAS).

**Storage Engines & Databases**: LSM-Tree compaction mechanics (RocksDB), PostgreSQL MVCC internals, distributed consensus algorithms (Raft, Paxos), Write-Ahead Logging (WAL), and storage hardware alignment.

**Network Protocols & Cryptography**: Deep dives into the QUIC and HTTP/3 protocol stacks, gRPC HTTP/2 multiplexing, TLS 1.3 Perfect Forward Secrecy, cryptography implementations, and terabit-scale network defense architectures.

**Production Incident Analysis (Post-Mortem & RCA)**: Chronological root-cause investigations, log forensics, failure reconstructions, and post-incident architectural mitigations.

---

## Frontmatter Specification & Markdown Formatting

Technical dispatches are stored as Markdown (`.md`) files in `content/posts/` with standard YAML frontmatter at the top:

```yaml
---
title: "Zero-Copy Architecture: Multiplying Throughput with sendfile Syscall"
slug: "performance-zero-copy-sendfile"
aliases: []
date: 2026-08-30
author: "Your Name or Handle"
contributors: ["github-handle"]
tags: ["performance", "low-level", "systems", "linux"]
lang: "en"
draft: false
description: "Exploring the Linux kernel zero-copy data path and sendfile system call that powers high-throughput data transfer in Nginx and Kafka."
cover: "https://images.unsplash.com/photo-1558494949-ef010cbdcc31?auto=format&fit=crop&w=1200&q=80"
coverCaption: "Data Transfer Optimization"
coverSource: "https://unsplash.com"
readTime: 6
---
```

### Supported Custom Markdown Elements

**Callout & Alert Blocks**: Use standard GitHub Markdown syntax: `> [!NOTE]`, `> [!TIP]`, `> [!IMPORTANT]`, and `> [!WARNING]`.

**Structured References Block**: Include an indexed YAML references block at the end of the post:

```references
- title: "Linux man-pages: sendfile(2)"
  url: "https://man7.org/linux/man-pages/man2/sendfile.2.html"
- title: "RFC 9000: QUIC - A UDP-Based Multiplexed and Secure Transport"
  url: "https://datatracker.ietf.org/doc/html/rfc9000"
```

**Quick Metrics Block**: Use ` ```stat ` blocks to highlight key benchmarks or numbers visually.

---

## Git Workflow & Submission

**Step 1: Fork & Clone the Repository**
```bash
git clone https://github.com/USERNAME/daemontalk.git
cd daemontalk
```

**Step 2: Create a Dedicated Branch**
- **New Dispatch**: `git checkout -b post/your-topic-slug`
- **Bug Fix / Code**: `git checkout -b fix/bug-description` or `feat/feature-name`
- **Content Correction**: `git checkout -b docs/fix-slug`

**Step 3: Test & Verify Locally**
Ensure templates compile and the full test suite passes:
```bash
make build
go test ./...
```

**Step 4: Open a Pull Request**
Submit a **Pull Request** to `main` at `https://github.com/dafagareth/daemontalk` with a concise summary of your changes.

---

## Alternative Submission Pathways

**Community Discussions & Feature Ideas**:
For interactive questions or RFCs, sign in with your **GitHub OAuth** account and start a new topic on our **Discussions Forum (`/socket`)**.

**Via Direct Email**:
If you prefer submitting article drafts via email for preliminary editorial review, send your Markdown file directly to **realdaemontalk@gmail.com** with the subject line `[Contribution Draft] Your Article Title`.

---

## Copyright & Contributor Attribution

**Author Ownership**: Authors retain full copyright and ownership of their original work.

**Content Licensing**: Written text is published under the Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International (CC BY-NC-SA 4.0) license, and software source code under permissive open-source terms.

**Profile Attribution**: Every merged contributor receives official attribution in article bylines, repository git history, and platform contributor badges.
