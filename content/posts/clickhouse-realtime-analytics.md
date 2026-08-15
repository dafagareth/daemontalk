---
title: "Arsitektur ClickHouse untuk Analitik Log Skala Terabita"
slug: 7f8a3d1e
aliases: [clickhouse-realtime-analytics]
date: 2026-04-30
tags: [database, storage, performance]
lang: id
draft: false
type: post
---

ClickHouse adalah sistem manajemen basis data kolom terdistribusi yang dirancang untuk analisis OLAP (Online Analytical Processing) real-time pada data berukuran terabita hingga petabita. Arsitektur penyimpanannya mengandalkan mesin tabel MergeTree yang mengurutkan data berdasarkan kunci primer dan membaginya ke dalam partisi fisik terkompresi. Tulisan ini membahas mekanika internal ClickHouse, mencakup kompresi data, eksekusi query ter-vektorisasi, perbandingan efisiensi agregasi dengan PostgreSQL, dan perancangan skema tabel log produksi.

## Fakta Menarik

**Fakta 1.** ClickHouse awalnya dikembangkan pada tahun 2009 di Yandex untuk menganalisis data log lalu lintas web Yandex.Metrica sebelum dirilis sebagai proyek open-source berlisensi Apache 2.0 pada tahun 2016.

**Fakta 2.** Mesin tabel MergeTree tidak melakukan pembaruan data secara langsung pada disk (in-place update), melainkan menulis partisi baru secara append-only dan menggabungkannya secara asynchronous di latar belakang.

**Fakta 3.** Eksekusi query ter-vektorisasi di ClickHouse memproses data dalam bentuk blok array berukuran 8192 baris, memanfaatkan instruksi SIMD (Single Instruction, Multiple Data) pada CPU modern untuk mempercepat operasi agregasi.

---

## Tips dan Trik

### 1. Strukturkan Tabel Menggunakan Engine MergeTree dan Primary Key Efektif

Mesin tabel MergeTree mengelompokkan data berdasarkan rentang indeks. Pilih urutan kunci pada klausa `ORDER BY` berdasarkan kolom yang paling sering digunakan dalam klausa `WHERE`.

```sql
CREATE TABLE sys_logs
(
    timestamp DateTime64(3, 'UTC'),
    service_name LowCardinality(String),
    hostname LowCardinality(String),
    level Enum8('DEBUG' = 1, 'INFO' = 2, 'WARN' = 3, 'ERROR' = 4),
    message String,
    duration_us UInt32
)
ENGINE = MergeTree()
PRIMARY KEY (service_name, timestamp)
ORDER BY (service_name, timestamp, level);
```

### 2. Terapkan Kombinasi Algoritma Kompresi Per Kolom

Gunakan codec kompresi spesifik seperti `DoubleDelta` atau `Delta` digabungkan dengan `ZSTD` atau `LZ4` untuk menghemat ruang simpan hingga 80 persen pada kolom numerik dan stempel waktu.

```sql
ALTER TABLE sys_logs MODIFY COLUMN
    timestamp DateTime64(3, 'UTC') CODEC(DoubleDelta, ZSTD(1)),
    duration_us UInt32 CODEC(T64, ZSTD(1)),
    message String CODEC(ZSTD(3));
```

### 3. Manfaatkan Vectorized Query Execution Melalui Agregasi Array

Struktur data berorientasi kolom memungkinkan ClickHouse memproses miliaran baris per detik. Buat query agregasi yang memaksimalkan pembacaan kolom berurutan di memori.

```sql
SELECT
    service_name,
    level,
    count() AS total_events,
    quantileExact(0.99)(duration_us) AS p99_latency_us
FROM sys_logs
WHERE timestamp >= now() - INTERVAL 1 HOUR
GROUP BY service_name, level
ORDER BY total_events DESC;
```

### 4. Perhatikan Perbedaan Performa Agregasi PostgreSQL vs ClickHouse

PostgreSQL menyimpan baris data secara utuh (row-oriented database), sehingga analisis agregasi memerlukan pembacaan seluruh blok baris dari disk. ClickHouse hanya membaca kolom yang dipanggil.

```sql
/* Di PostgreSQL, query ini harus membaca baris utuh dari tabel 100GB: */
/* SELECT service_name, COUNT(*) FROM logs GROUP BY service_name; */

/* Di ClickHouse, query ini hanya membaca file kolom service_name (misal 500MB): */
SELECT service_name, count()
FROM sys_logs
GROUP BY service_name;
```

### 5. Atur Retensi Data Otomatis Menggunakan Kebijakan TTL

Kelola ukuran disk pada kluster produksi dengan menambahkan klausa `TTL` pada definisi tabel untuk menghapus data lama secara otomatis tanpa mengganggu proses penulisan.

```sql
ALTER TABLE sys_logs
MODIFY TTL timestamp + INTERVAL 30 DAY;
```
