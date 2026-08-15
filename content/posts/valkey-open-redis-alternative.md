---
title: "Valkey: Alternatif Redis Terbuka Berlisensi BSD dari Linux Foundation"
slug: c2d8e4f1
aliases: [valkey-open-redis-alternative]
date: 2026-06-28
tags: [database, devops, tools]
lang: id
draft: false
type: post
---

Valkey adalah proyek penyimpanan struktur data in-memory berlisensi BSD 3-Clause yang didirikan di bawah naungan Linux Foundation sebagai respon terhadap perubahan lisensi Redis 7.2+ menjadi dual-license non-open-source (RSALv2 dan SSPLv1). Valkey mempertahankan kompatibilitas penuh dengan protokol Redis RESP2 dan RESP3, sembari memperkenalkan arsitektur I/O multithreading baru pada versi 8.0. Tulisan ini menguraikan latar belakang peralihan proyek, arsitektur internal Valkey 8.0, peningkatan throughput, serta tahapan migrasi kluster tanpa waktu henti (zero-downtime migration).

## Fakta Menarik

**Fakta 1.** Proyek Valkey didukung oleh perusahaan skala besar seperti Amazon Web Services, Google Cloud, Oracle, Ericsson, dan Snap untuk menjamin kepemilikan netral dan keberlanjutan lisensi open-source.

**Fakta 2.** Valkey 8.0 merombak penanganan jaringan dengan mengisolasi pemrosesan I/O socket dan parsing protokol ke worker thread terpisah, sehingga benang utama (main thread) fokus pada eksekusi perintah di memori.

**Fakta 3.** Valkey tetap kompatibel secara drop-in dengan CLI `redis-cli` dan library klien Redis di berbagai bahasa pemrograman tanpa perlu mengubah nama perintah.

---

## Tips dan Trik

### 1. Pahami Alasan Perubahan Lisensi dan Migrasi ke Valkey

Lisensi RSALv2 dan SSPLv1 melarang penggunaan Redis sebagai layanan terkelola komersial tanpa lisensi berbayar dari Redis Inc. Lisensi BSD 3-Clause pada Valkey menjamin kebebasan komersial dan integrasi ulang tanpa pembatasan hukum.

```bash
# Periksa versi dan informasi rilis server Valkey
valkey-server -v
valkey-cli info server | grep valkey_version
```

### 2. Konfigurasikan Multithreaded I/O pada Valkey 8.0

Aktifkan peningkatan I/O terulir pada file konfigurasi `valkey.conf` untuk mengeksploitasi inti CPU tambahan pada instance server besar.

```ini
# Tentukan jumlah thread I/O berdasarkan inti CPU yang tersedia
io-threads 8

# Izinkan thread I/O membaca dan menulis ke socket client
io-threads-do-reads yes

# Alokasikan buffer memori per-thread untuk mengurangi kontensi lock
io-threads-mode async-per-client
```

### 3. Ukur Peningkatan Throughput Menggunakan Valkey Benchmark

Valkey 8.0 mampu mencapai throughput hingga 1.2 juta operasi per detik pada single instance berkat arsitektur threading yang diperbarui.

```bash
# Jalankan pengujian performa benchmark perintah SET dan GET
valkey-benchmark -h 127.0.0.1 -p 6379 -c 100 -n 1000000 -t set,get -P 16
```

### 4. Lakukan Migrasi Kluster Redis ke Valkey Tanpa Downtime

Gunakan teknik replikasi master-replica di mana node Valkey bergabung ke kluster Redis yang ada sebagai replica sebelum dipromosikan menjadi failover master.

```bash
# Hubungkan node Valkey baru sebagai replica dari master Redis
valkey-cli -h 10.0.0.2 -p 6379 REPLICATOF 10.0.0.1 6379

# Verifikasi status sinkronisasi data dari master
valkey-cli -h 10.0.0.2 -p 6379 INFO replication

# Promosikan node Valkey menjadi master independen setelah sinkronisasi selesai
valkey-cli -h 10.0.0.2 -p 6379 REPLICATOF NO ONE
```

### 5. Sesuaikan Konfigurasi Memori dan Eviction Policy

Gunakan algoritma pengosongan memori `allkeys-lru` atau `volatile-lfu` untuk menjaga stabilitas memori RAM server saat menangani beban memori tinggi.

```ini
# Batasi penggunaan memori maksimum hingga 16 Gigabyte
maxmemory 16gb

# Tentukan kebijakan pengosongan data saat memori penuh
maxmemory-policy allkeys-lru
```
