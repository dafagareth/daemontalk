## Filosofi & Standar Kualitas

Daemontalk adalah publikasi rekayasa sistem komputer independen yang mengutamakan artikel teknis mendalam, berorientasi eksperimen nyata, dan dapat direproduksi (*reproducible*). Kami menyambut kontribusi dari siapa pun—mahasiswa, rekayasawan perangkat lunak (*software engineers*), arsitek sistem, maupun peneliti yang ingin membagikan catatan teknisnya atau mengembangkan ekosistem Daemontalk.

**Lugas dan Bebas Basa-Basi**: Awali tulisan langsung ke inti persoalan teknis, diagram arsitektur, atau cuplikan kode. Hindari pengantar yang bertele-tele dan klise pembuka yang tidak bernilai teknis.

**Verifikasi & Reproduksibilitas**: Setiap hipotesis atau klaim performa wajib dilengkapi dengan cuplikan kode uji, perintah terminal, log diagnostik, atau diagram arsitektur yang memungkinkan pembaca menguji dan memverifikasi sendiri temuannya.

**Referensi Terverifikasi**: Lengkapi setiap artikel dengan blok referensi ke dokumen resmi otoritatif, seperti RFC standar internet, repositori kode sumber kernel Linux, manual arsitektur prosesor, atau *paper* riset ilmiah.

---

## Pilar Kontribusi yang Diterima

Daemontalk terbuka untuk berbagai bentuk kontribusi teknis:

**1. Penulisan Artikel & Riset Sistem (`content/posts/`)**:
Menulis artikel analisis mendalam, eksplorasi kernel, protokol jaringan, basis data, atau studi kasus sistem terdistribusi.

**2. Pengembangan Kode Sumber & Perbaikan Bug (*Core Engine*)**:
- Backend Go (HTTP Handlers, SQLite storage, CLI tools, parser Goldmark).
- Antarmuka Terminal UI (*Bubble Tea / Lip Gloss*).
- Templ Components & Styling Tailwind CSS.
- Peningkatan performa, *security hardening*, dan optimasi query database.

**3. Koreksi Konten & Pembaruan Cuplikan Kode**:
Menemukan kesalahan logika kode, perintah terminal yang usang, tautan rujukan mati, atau *typo* pada artikel yang sudah terbit? Anda dapat langsung mengajukan Pull Request perbaikan pada berkas `.md` terkait.

**4. Terjemahan & Lokalisasi (*i18n*)**:
Membantu menerjemahkan artikel atau teks antarmuka UI ke dalam Bahasa Indonesia (`.id.md`), Inggris (`.md`), atau Spanyol (`.es.md`).

---

## Domain & Topik yang Diminati

Berikut adalah bidang utama yang menjadi fokus publikasi di Daemontalk:

**Sistem Operasi & Linux Kernel**: Eksplorasi mekanisme eBPF/XDP, isolasi resource dengan cgroups v2, penjadwalan CPU (EEVDF), manajemen memori (*OOM Killer internals*), *system calls*, dan I/O *zero-copy* (`sendfile`, `io_uring`).

**Konkurensi & Runtime Bahasa Pemrograman**: Arsitektur runtime Go (Tri-color Garbage Collector, scheduler goroutine), manajemen memori Rust (*borrow checker*, *lifetimes*), dan struktur data konkuren *lock-free* berbasis *Compare-And-Swap* (CAS).

**Penyimpanan & Basis Data**: Mekanisme LSM-Tree compaction (RocksDB), arsitektur PostgreSQL MVCC, algoritma konsensus terdistribusi (Raft, Paxos), Write-Ahead Logging (WAL), dan optimasi I/O penyimpanan.

**Protokol Jaringan & Kriptografi**: Bedah protokol QUIC dan HTTP/3, gRPC multiplexing di atas HTTP/2, enkripsi TLS 1.3 dengan *Perfect Forward Secrecy*, dan arsitektur mitigasi serangan jaringan berskala terabit.

**Analisis Insiden Produksi (Post-Mortem & RCA)**: Kronologi investigasi akar masalah (*Root Cause Analysis*), pembedahan rekaman log diagnostik, rekonstruksi kegagalan sistem, dan mitigasi arsitektur pasca-insiden.

---

## Spesifikasi Frontmatter & Format Markdown

Setiap artikel teknis disimpan dalam format berkas Markdown (`.md`) di dalam direktori `content/posts/` dengan konfigurasi frontmatter YAML di awal berkas:

```yaml
---
title: "Arsitektur Zero-Copy: Melipatgandakan Throughput dengan Syscall sendfile"
slug: "performance-zero-copy-sendfile"
aliases: []
date: 2026-08-30
author: "Nama Anda atau Callsign"
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

### Elemen Parser Khusus Daemontalk

**Blok Callout & Peringatan**: Gunakan sintaks standar GitHub Markdown seperti `> [!NOTE]`, `> [!TIP]`, `> [!IMPORTANT]`, dan `> [!WARNING]`.

**Blok Referensi Terstruktur**: Letakkan blok referensi YAML di akhir artikel untuk menyajikan rujukan resmi yang terindeks rapi:

```references
- title: "Dokumentasi Linux man page sendfile(2)"
  url: "https://man7.org/linux/man-pages/man2/sendfile.2.html"
- title: "RFC 9000: QUIC - A UDP-Based Multiplexed and Secure Transport"
  url: "https://datatracker.ietf.org/doc/html/rfc9000"
```

**Blok Statistik Cepat**: Gunakan blok ` ```stat ` untuk menyajikan angka metrik kunci atau data performa secara visual.

---

## Alur Kerja Git (*Git Workflow*) & Pengiriman

**Langkah 1: Fork & Clone Repositori**
```bash
git clone https://github.com/USERNAME/daemontalk.git
cd daemontalk
```

**Langkah 2: Buat Cabang (*Branch*) Khusus Sesuai Jenis Kontribusi**
- **Artikel Baru**: `git checkout -b post/nama-topik`
- **Perbaikan Bug / Kode**: `git checkout -b fix/deskripsi-bug` atau `feat/nama-fitur`
- **Koreksi Artikel / Typo**: `git checkout -b docs/perbaikan-slug`

**Langkah 3: Uji Coba & Verifikasi Lokal**
Pastikan generator template dan seluruh unit test berjalan sukses:
```bash
make build
go test ./...
```

**Langkah 4: Buka Pull Request**
Buka **Pull Request** di repositori utama `https://github.com/dafagareth/daemontalk`. Sertakan ringkasan singkat perubahan yang Anda lakukan.

---

## Jalur Alternatif

**Partisipasi Forum Diskusi & Tanya Jawab**:
Untuk diskusi teknis, tanya jawab konfigurasi, atau diskusi ide fitur, Anda dapat langsung masuk menggunakan akun **GitHub OAuth** dan membuat topik baru di halaman **Forum Diskusi (`/socket`)**.

**Melalui Email**:
Jika Anda ingin mengirimkan draf naskah awal untuk ditinjau oleh tim redaksi sebelum membuka PR, kirimkan berkas Markdown Anda ke **realdaemontalk@gmail.com** dengan subjek `[Draft Kontribusi] Judul Tulisan Anda`.

---

## Hak Cipta & Atribusi Kontributor

**Kepemilikan Hak Cipta**: Penulis memegang hak cipta penuh dan kepemilikan atas karya orisinal yang ditulis.

**Lisensi Publikasi**: Artikel teks di Daemontalk diterbitkan di bawah lisensi Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International (CC BY-NC-SA 4.0), dan kode sumber platform berada di bawah lisensi terbuka permisif.

**Atribusi Profil & Badge**: Setiap kontributor yang karyanya digabungkan (*merged*) akan tercantum di byline artikel, riwayat git repositori, serta mendapatkan badge kontributor resmi pada ekosistem platform.
