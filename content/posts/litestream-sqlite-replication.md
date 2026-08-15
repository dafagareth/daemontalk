---
title: "Replikasi Kontinu SQLite ke S3 dengan Litestream"
slug: e4b8a1c9
aliases: [litestream-sqlite-replication]
date: 2026-04-28
tags: [sqlite, storage, devops]
lang: id
draft: false
type: post
---

SQLite merupakan mesin database embedded yang efisien untuk aplikasi web monolitik. Namun, strategi backup tradisional sering kali mengganggu operasi I/O atau menghasilkan salinan file yang tidak konsisten. Litestream menyelesaikan masalah ini dengan melakukan streaming Write-Ahead Log (WAL) secara langsung ke cloud storage.

## Fun Fact

**Fact 1.** SQLite menggunakan satu file tunggal di sistem berkas, namun melakukan penyalinan biasa saat proses penulisan berlangsung dapat mengakibatkan kerusakan struktur database.

**Fact 2.** Litestream berjalan sebagai proses daemon terpisah yang memantau file WAL tanpa memblokir pembacaan atau penulisan dari aplikasi.

**Fact 3.** Fitur Point-In-Time Recovery pada Litestream memungkinkan pemulihan status database hingga tingkat milidetik sebelum insiden kerusakan data terjadi.

---

## Tips dan Trik

### 1. Keterbatasan Backup Tradisional dan Arsitektur WAL Streaming

Pendekatan backup SQLite konvensional biasanya mengandalkan perintah `sqlite3 db.sqlite ".backup copy.db"` atau snapshot tingkat sistem berkas. Metode ini membutuhkan koordinasi lock yang dapat memblokir transaksi penulisan.

Litestream memanfaatkan mode Write-Ahead Log (WAL) bawaan SQLite. Saat aplikasi menulis data ke halaman WAL, Litestream membaca halaman baru tersebut secara asinkron dan mengirimkannya ke objek penyimpanan seperti AWS S3 atau MinIO.

### 2. Konfigurasi Litestream untuk Replikasi ke S3

Buat file konfigurasi `/etc/litestream.yml` untuk menentukan lokasi database lokal dan target replikasi:

```yaml
access-key-id: AKIAIOSFODNN7EXAMPLE
secret-access-key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY

dbs:
  - path: /var/lib/app/production.db
    replicas:
      - type: s3
        bucket: my-app-backups
        path: db
        endpoint: https://s3.us-west-004.backblazeb2.com
```

### 3. Mengintegrasikan Litestream Daemon dengan Systemd

Agar proses replikasi berjalan secara kontinu di latar belakang, buat unit file systemd di `/etc/systemd/system/litestream.service`:

```ini
[Unit]
Description=Litestream SQLite Replication Service
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/litestream replicate -config /etc/litestream.yml
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Jalankan perintah berikut untuk mengaktifkan layanan:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now litestream
```

### 4. Prosedur Point-In-Time Recovery (PITR)

Apabila file database utama mengalami kerusakan atau terhapus secara tidak sengaja, data dapat dipulihkan ke titik waktu tertentu dengan argumen `-timestamp`:

```bash
litestream restore -o /var/lib/app/production.db \
  -timestamp 2026-04-28T14:30:00Z \
  https://s3.us-west-004.backblazeb2.com/my-app-backups/db
```

Perintah di atas mengunduh snapshot basis beserta rantai WAL terkait, lalu menyusun kembali file database hingga stempel waktu yang ditentukan.

### 5. Analisis Efisiensi Biaya: SQLite + Litestream vs Managed PostgreSQL

Mengoperasikan instance managed PostgreSQL pada penyedia cloud umumnya membutuhkan biaya minimal 15 hingga 50 USD per bulan untuk spesifikasi dasar.

Kombinasi SQLite dan Litestream pada VPS kecil seharga 5 USD per bulan yang terhubung ke cloud storage S3 hanya menghabiskan biaya penyimpanan beberapa sen per gigabyte. Model arsitektur ini menekan pengeluaran operasional infrastruktur secara signifikan tanpa mengorbankan durabilitas data.
