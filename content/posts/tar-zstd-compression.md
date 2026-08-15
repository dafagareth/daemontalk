---
title: "Zstd vs Gzip: Kompresi Kilat untuk Backup Data Raksasa"
slug: 7f8a9b0c
aliases: [tar-zstd-compression]
date: 2026-04-20
tags: [linux, cli, storage]
lang: id
draft: false
---

Selama lebih dari dua dekade, `tar -czvf` dengan kompresi Gzip adalah kombinasi standar untuk membuat arsip file di Linux. Namun, di era database puluhan gigabyte dan multi-core CPU, algoritma modern **Zstandard (Zstd)** dari Meta telah mengubah standar kompresi data.

Zstd mampu mengompresi data 3x hingga 5x lebih cepat dari Gzip dengan rasio kompresi yang setara atau bahkan lebih padat.

## Fun Fact

**Zstandard dirancang oleh Yann Collet di Facebook (Meta) pada 2015.**
Yann Collet juga adalah pencipta LZ4, algoritma kompresi real-time tercepat yang dipakai di kernel Linux dan game engine.

**Linux Kernel, Arch Linux, dan Fedora beralih ke Zstd.**
Arch Linux mengganti kompresi paket `.pkg.tar.xz` menjadi `.pkg.tar.zst` pada 2020, memangkas waktu dekompresi package install hingga 14 kali lebih cepat!

**Zstd mendukung custom dictionary training.**
Jika kamu mengompresi ribuan file kecil dengan struktur mirip (seperti log JSON atau data REST API), Zstd bisa dilatih dengan 'kamus' khusus untuk meningkatkan rasio kompresi hingga ratusan persen.

---

## Tips dan Trik

### 1. Buat Arsip Tar dengan Zstd Multi-Thread

Manfaatkan seluruh core CPU yang ada di komputermu dengan flag `-T0`:

```bash
tar --zstd -cf backup.tar.zst /path/to/data/
```

Atau menggunakan command `zstd` langsung:

```bash
tar -cf - /path/to/data | zstd -T0 -3 -o backup.tar.zst
```

### 2. Ekstrak Arsip `.tar.zst` dengan Cepat

Tar modern di Linux otomatis mengenali format Zstd tanpa perlu flag khusus:

```bash
tar -xf backup.tar.zst
```

### 3. Pilih Level Kompresi Sesuai Kebutuhan

Zstd mendukung level 1 (super cepat) hingga 19 (rasio maksimal, atau 22 dengan `--ultra`):

```bash
# Mode cepat untuk transfer network real-time (level 1)
zstd -1 -T0 bigfile.iso

# Mode arsip dingin / cold backup (level 19)
zstd -19 -T0 database_dump.sql
```

### 4. Bandingkan Rasio dan Kecepatan dengan Benchmark Bawaan

Uji coba kompresi langsung pada file sampel di RAM tanpa menulis file ke disk:

```bash
zstd -b3 -e19 large_log.log
```

### 5. Dekompresi Paralel Super Kilat

Zstd dirancang agar waktu dekompresi hampir selalu mencapai batas kecepatan throughput I/O SSD/NVMe:

```bash
unzstd -T0 backup.tar.zst
```
