---
title: "SQLite WAL Mode: Mengatasi Locking Saat Traffic Tinggi"
slug: 3c4d5e6f
aliases: [sqlite-wal-concurrency]
date: 2026-06-28
tags: [sqlite, database, backend]
lang: id
draft: false
---

Banyak developer menganggap SQLite tidak sanggup melayani traffic concurrent karena masalah `database is locked`. Sebenarnya, 90% masalah tersebut terjadi karena database masih berjalan di mode rollback journal default, bukan **WAL (Write-Ahead Logging)**.

Dengan konfigurasi pragmas yang tepat, SQLite mampu melayani ribuan request per detik tanpa gangguan.

## Fun Fact

**SQLite adalah software database yang paling banyak di-deploy di seluruh planet.**
Setiap smartphone Android, iPhone, browser Chrome, Firefox, macOS, Windows 10/11, dan mobil Tesla memiliki puluhan database SQLite aktif.

**SQLite ditulis dan dirancang oleh Dr. D. Richard Hipp pada tahun 2000.**
Ia mendesainnya saat bekerja untuk Angkatan Laut AS di atas kapal perang perusak rudal USS Coronado.

**WAL mode memisahkan operasi pembacaan (Readers) dan penulisan (Writers).**
Dalam mode WAL, pembaca tidak pernah memblokir penulis, dan penulis tidak pernah memblokir pembaca.

---

## Tips dan Trik

### 1. Aktifkan WAL Mode Segera Setelah Membuka Database

Jalankan perintah PRAGMA ini saat inisialisasi koneksi aplikasi:

```sql
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA foreign_keys = ON;
```

### 2. Atur `busy_timeout` di Koneksi Driver Go

Beri toleransi beberapa detik bagi SQLite untuk menunggu giliran tulis alih-alih langsung melempar error `busy`:

```go
db, err := sql.Open("sqlite3", "file:data.db?_journal_mode=WAL&_busy_timeout=5000")
```

### 3. Batasi `MaxOpenConns` untuk Penulis di Go

Jika aplikasi melakukan banyak operasi mutasi data, gunakan pool koneksi tunggal untuk proses write agar tidak terjadi perebutan lock internal:

```go
db.SetMaxOpenConns(1) // Atau pisahkan pool Read dan Write
```

### 4. Simpan Database Temporary di Memory

Tingkatkan kecepatan query sorting dan indeks sementara dengan mengarahkan temp store ke RAM:

```sql
PRAGMA temp_store = MEMORY;
PRAGMA cache_size = -64000; -- Alokasikan cache sekitar 64MB
```

### 5. Lakukan Backup Aman Saat Database Sedang Berjalan

Gunakan SQLite online backup API atau command line tool bawaan tanpa perlu mematikan aplikasi:

```bash
sqlite3 data.db ".backup 'data_backup.db'"
```
