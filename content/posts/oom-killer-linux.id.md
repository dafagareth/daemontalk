---
title: "Server RAM 64GB Tapi Selalu Crash? Anatomi Kelam OOM Killer di Linux"
slug: anatomi-oom-killer-linux
aliases: []
date: 2026-08-19
author: "Daemontalk Editorial"
tags: ["linux", "systems", "kernel", "devops"]
lang: id
draft: false
type: post
cover: ""
coverCaption: ""
coverSource: ""
readTime: 7
description: "Pernahkah server dengan RAM besar tiba-tiba membunuh proses database atau aplikasi Anda? Mari membedah arsitektur virtual memory dan heuristik Out-Of-Memory (OOM) Killer di Linux secara akademis."
series: "Linux Kernel Deep Dive"
series_part: 1
---

Anda menyewa *Virtual Private Server* (VPS) dengan RAM 64GB yang cukup mahal. Aplikasi Golang dan database PostgreSQL Anda berjalan lancar selama berminggu-minggu. Namun tiba-tiba, di tengah malam saat *traffic* melonjak, database mati seketika tanpa pesan *error* yang jelas di log aplikasi.

Selamat, Anda baru saja menjadi korban dari algojo sistem operasi paling ditakuti: **Linux Out-Of-Memory (OOM) Killer**.

Bagaimana mungkin sistem dengan RAM sebesar itu bisa kehabisan memori secara mendadak? Artikel ini akan membedah anatomi *virtual memory* di Linux, bagaimana kernel "berbohong" tentang kapasitas RAM kepada aplikasi Anda, dan matematika di balik penentuan proses mana yang akan "dibunuh" oleh sistem.

## 1. Ilusi Memori: Linux Overcommit

Sebagian besar *software engineer* menganggap RAM sebagai sumber daya statis yang terukur pasti. Faktanya, arsitektur *Virtual Memory Manager* (VMM) di Linux dibangun di atas prinsip **optimisme yang buta** atau yang dikenal sebagai *memory overcommit*.

Saat sebuah proses memanggil fungsi `malloc()` di bahasa C (atau saat *runtime* Go meminta memori ke OS), kernel Linux akan selalu menjawab "Ya", terlepas dari apakah masih ada sisa RAM fisik (RAM + Swap) atau tidak. Linux hanya mengalokasikan ruang di *Virtual Address Space* (VMA), bukan RAM fisik yang sebenarnya.

Hal ini dapat dipahami melalui persamaan alokasi memori aktual:

$$
M_{\text{actual}} = \sum_{i=1}^{n} \text{RSS}_i + \text{Page Cache} + \text{Slab}
$$

Di mana:
- **RSS (Resident Set Size)** adalah memori fisik yang benar-benar diakses/disentuh (*page fault*) oleh aplikasi.
- **Page Cache** adalah sisa memori yang digunakan kernel untuk *caching disk I/O*.

> [!WARNING]
> Masalah terjadi ketika semua aplikasi mulai benar-benar menulis (*write*) ke alamat virtual yang dijanjikan tersebut secara bersamaan. Jika total RSS melonjak melewati kapasitas memori fisik yang ada, sistem akan mengalami panik.

## 2. Sang Algojo: OOM Killer Heuristics

Ketika cadangan memori bebas (*free memory pages*) jatuh di bawah ambang batas kritis (ditentukan oleh `vm.min_free_kbytes`), Linux tidak punya pilihan selain membebaskan memori secepat mungkin. Di sinilah **OOM Killer** terbangun.

OOM Killer tidak membunuh proses secara acak. Ia menggunakan fungsi akademis dan deterministik yang disebut `oom_badness()` di dalam *source code* kernel (`mm/oom_kill.c`).

Fungsi ini menghitung skor "kejahatan" (*badness score*) dari 0 hingga 1000 untuk setiap proses yang berjalan. Skor tertinggi akan dieksekusi.

Rumus sederhananya berbanding lurus dengan memori yang dikonsumsi:

$$
\text{Badness Score} = \left( \frac{\text{RSS} + \text{Swap} + \text{Pagetables}}{\text{Total RAM} + \text{Total Swap}} \right) \times 1000 + \text{oom\_score\_adj}
$$

```stat
- value: "+1000"
  label: "Maximum Badness"
  description: "Proses ini 100% dipastikan akan dibunuh oleh OOM Killer."

- value: "-1000"
  label: "OOM Disabled"
  description: "Proses dijamin kebal dari eksekusi (biasanya daemon sistem seperti sshd)."

- value: "3%"
  label: "Sisa Memori"
  description: "Batas kritis rata-rata di mana kernel akan memicu algoritma OOM."
```

## 3. Melindungi Sistem Anda (Mitigasi)

Bagaimana cara agar database Anda (`postgres` atau `mysqld`) tidak menjadi korban salah sasaran dari proses *memory-leak* di aplikasi web Anda?

Anda bisa memanipulasi variabel `oom_score_adj` untuk memanipulasi perhitungan skor di atas. Nilai negatif melindungi proses, sementara nilai positif menjadikannya tumbal.

Berikut adalah tiga cara profesional untuk melindungi servis Anda di tataran infrastruktur:

```tabs
=== Systemd (Disarankan)
# /etc/systemd/system/postgresql.service.d/override.conf
[Service]
# Berikan kekebalan dari OOM Killer (nilai antara -1000 hingga 1000)
# -1000 berarti kebal total, -500 sangat dilindungi.
OOMScoreAdjust=-900

=== Docker / Docker Compose
# docker-compose.yml
services:
  db:
    image: postgres:15
    # Melindungi kontainer database dari pembunuhan kernel host
    oom_score_adj: -500
    deploy:
      resources:
        limits:
          memory: 4G

=== Sysctl (Global Policy)
# Mencegah overcommit berlebihan di tingkat OS
# /etc/sysctl.d/99-custom.conf
vm.overcommit_memory = 2
vm.overcommit_ratio = 80
```

> [!TIP]
> Jangan pernah mengubah `vm.overcommit_memory = 2` kecuali Anda menggunakan bahasa pemrograman yang tidak melakukan *greedy allocation* (seperti C/C++). Bahasa modern berdasar *Garbage Collector* seperti Java, Python, atau Go sering kali gagal secara sporadis (panic) jika `overcommit` dimatikan.

## Kesimpulan

Menyalahkan aplikasi saat terjadi *crash* di server dengan RAM besar seringkali adalah analisis yang dangkal. Memahami cara Linux mengelola *virtual memory* dan matematika heuristik di balik OOM Killer adalah perbedaan antara administrator sistem amatir dan profesional sejati (*Systems Engineer*). 

Dengan mengelola limit memori (melalui `cgroups` atau Docker) dan memanipulasi `oom_score_adj`, Anda bisa memastikan bahwa jika krisis memori terjadi, proses aplikasi tak penting-lah yang dikorbankan, bukan *core database* penyimpan data pengguna Anda.

## Referensi Akademis & Literatur

```references
- title: "Linux Kernel Documentation: Out Of Memory Management"
  author: "The Linux Kernel Authors"
  year: 2024
  url: "https://www.kernel.org/doc/gorman/html/understand/understand016.html"

- title: "Systems Performance: Enterprise and the Cloud, 2nd Edition"
  author: "Brendan Gregg"
  year: 2020
  publisher: "Addison-Wesley Professional"
  url: "https://www.brendangregg.com/systems-performance-2nd-edition-book.html"
```

```author
name: Daemontalk Editorial
role: Systems Architecture Desk
avatar: /static/logo/favicon-32x32.png
bio: Menyelami lapisan terdalam sistem Linux, backend Go, dan optimasi arsitektur infrastruktur terdistribusi tingkat produksi.
github: https://github.com/dafagareth/daemontalk
website: https://daemontalk.com
```
