## Filosofi & Standar Editorial

Daemontalk adalah publikasi rekayasa sistem independen yang mengutamakan artikel teknis mendalam, berorientasi eksperimen nyata, dan dapat direproduksi (*reproducible*).

- **Lugas dan Bebas Basa-Basi**: Awali tulisan langsung ke inti persoalan teknis, diagram arsitektur, atau cuplikan kode. Hindari pengantar yang bertele-tele.
- **Verifikasi & Reproduksibilitas**: Setiap klaim performa atau analisis wajib didukung cuplikan kode uji, perintah shell, log diagnostik, atau diagram arsitektur.
- **Referensi Otoritatif**: Lengkapi setiap tulisan dengan blok referensi ke RFC resmi, kode sumber kernel Linux, dokumen arsitektur CPU, atau *paper* riset ilmiah.

---

## Domain & Topik yang Diminati

1. **Sistem Operasi & Linux Kernel**: eBPF/XDP, cgroups v2, penjadwalan CPU (EEVDF), Linux OOM Killer internals, system call profiling, dan zero-copy I/O (`sendfile`, `io_uring`).
2. **Konkurensi & Runtime Bahasa**: Go GC (Tri-color mark-and-sweep), goroutine scheduler, Rust borrow checker, dan struktur data lock-free CAS.
3. **Storage & Database**: LSM-Tree compaction (RocksDB), PostgreSQL MVCC, algoritma konsensus Raft/Paxos, dan Write-Ahead Logging (WAL).
4. **Protokol Jaringan & Kriptografi**: QUIC dan HTTP/3, gRPC multiplexing, TLS 1.3 Perfect Forward Secrecy, dan arsitektur mitigasi DDoS terabit.
5. **Analisis Insiden (Post-Mortem & RCA)**: Kronologi Root Cause Analysis, pembedahan log forensik, dan rekonstruksi insiden produksi.

---

## Format Frontmatter YAML & Markdown

Setiap artikel disimpan di `content/posts/nama-topik.md` dengan frontmatter YAML di awal berkas:

```yaml
---
title: "Arsitektur Zero-Copy: Melipatgandakan Throughput dengan Syscall sendfile"
slug: "performance-zero-copy-sendfile"
aliases: []
date: 2026-08-30
author: "Nama Anda"
contributors: ["username-github"]
tags: ["performance", "low-level", "systems", "linux"]
lang: "id"
draft: false
description: "Membongkar cara kerja mekanisme zero-copy kernel Linux yang membuat Nginx dan Kafka mampu mentransfer data berkecepatan tinggi."
cover: "https://images.unsplash.com/photo-1558494949-ef010cbdcc31?auto=format&fit=crop&w=1200&q=80"
coverCaption: "Data Transfer Optimization"
coverSource: "https://unsplash.com"
readTime: 6
---
```

### Elemen Markdown Khusus

- **Callouts**: Gunakan `> [!NOTE]`, `> [!TIP]`, `> [!IMPORTANT]`, atau `> [!WARNING]`.
- **Referensi Terstruktur**:
```references
- title: "Linux man-pages: sendfile(2)"
  url: "https://man7.org/linux/man-pages/man2/sendfile.2.html"
- title: "RFC 9000: QUIC Protocol"
  url: "https://datatracker.ietf.org/doc/html/rfc9000"
```
