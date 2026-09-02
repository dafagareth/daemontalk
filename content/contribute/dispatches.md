## Philosophy & Editorial Standards

Daemontalk is an independent systems engineering publication dedicated to practical, verifiable, and reproducible technical dispatches.

- **Straight to the Point**: Dive directly into the technical core, architectural blueprints, or code examples in the opening section. Avoid conversational fluff.
- **Verification & Reproducibility**: Every technical claim or benchmark must be supported by actionable test code, diagnostic logs, shell commands, or architecture diagrams.
- **Authoritative Citations**: Conclude every in-depth dispatch with verified citations to authoritative specifications, RFC standards, Linux kernel source code, or academic research papers.

---

## Topics & Domains of Interest

1. **Operating Systems & Linux Kernel**: eBPF/XDP packet processing, cgroups v2, CPU scheduling (EEVDF), Linux OOM Killer internals, system call profiling, and zero-copy I/O (`sendfile`, `io_uring`).
2. **Concurrency & Language Runtimes**: Go runtime internals (Tri-color GC, goroutine scheduler), Rust memory management, and lock-free CAS data structures.
3. **Storage Engines & Databases**: LSM-Tree compaction (RocksDB), PostgreSQL MVCC internals, distributed consensus (Raft/Paxos), and Write-Ahead Logging (WAL).
4. **Network Protocols & Cryptography**: QUIC & HTTP/3 protocol stacks, gRPC multiplexing, TLS 1.3 Perfect Forward Secrecy, and terabit-scale network defense architectures.
5. **Production Incident Analysis (Post-Mortem & RCA)**: Chronological root-cause investigations, log forensics, and post-incident architectural mitigations.

---

## Frontmatter Specification & Formatting

Technical dispatches are stored as Markdown files in `content/posts/your-topic.md` with standard YAML frontmatter:

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

### Special Markdown Syntax

- **Callout Blocks**: Use standard GitHub syntax `> [!NOTE]`, `> [!TIP]`, `> [!IMPORTANT]`, or `> [!WARNING]`.
- **Structured References**:
```references
- title: "Linux man-pages: sendfile(2)"
  url: "https://man7.org/linux/man-pages/man2/sendfile.2.html"
- title: "RFC 9000: QUIC Protocol"
  url: "https://datatracker.ietf.org/doc/html/rfc9000"
```
