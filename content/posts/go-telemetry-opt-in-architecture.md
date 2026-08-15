---
title: "Arsitektur Go Telemetry: Pengumpulan Data Performa Toolchain Berbasis Opt-In"
slug: 9a3f2b1c
aliases: [go-telemetry-opt-in-architecture]
date: 2026-07-22
tags: [go, security, architecture]
lang: id
cover: /static/images/golang.png
draft: false
type: post
---
Sistem Go Telemetry yang diperkenalkan secara stabil pada Go 1.23 merancang pengumpulan data metrik penggunaan toolchain compiler Go secara transparansi tinggi dan transparan. Pendekatan ini memprioritaskan privasi pengembang dengan menerapkan model eksplisit opt-in, pencatatan counter lokal terenkapsulasi, dan teknik differential privacy sebelum data agregat diunggah. Tulisan ini membedah arsitektur dasar Go Telemetry, mekanisme p-bit counter, penerapan privasi diferensial, serta cara verifikasi data menggunakan CLI `gotelemetry`.
![GO Architecture](../web/static/images/golang.png)

## Fakta Menarik

**Fakta 1.** Telemetri di Go secara bawaan berada pada status nonaktif (off) untuk seluruh pemasangan standar, sehingga tidak ada data yang dikirim ke server Go tanpa persetujuan manual pengembang.

**Fakta 2.** Laporan telemetri yang diunggah hanya berisi counter agregat mingguan berbasis tanggal dan nama fungsi internal compiler, tanpa menyertakan alamat IP, nama proyek, isi kode sumber, atau identitas unik komputer.

**Fakta 3.** Go Telemetry mengimplementasikan algoritma Differential Privacy berbasis sampel p-bit untuk memastikan bahwa keberadaan data dari satu mesin pengembang tidak dapat dideduksi secara statistik dari dataset umum.

---

## Tips dan Trik

### 1. Pahami Status Operasional Go Telemetry

Go Telemetry menyediakan tiga mode operasi: `off` (bawaan), `local` (hanya mencatat di disk lokal), dan `on` (mencatat dan mengunggah laporan mingguan).

```bash
# Periksa mode aktif Go Telemetry saat ini
gotelemetry env

# Aktifkan perekaman data lokal tanpa mengunggah ke server
gotelemetry local

# Aktifkan telemetri penuh (opt-in)
gotelemetry on
```

### 2. Inspeksi File Counter Lokal pada Mesin Pengembang

Data telemetri disimpan dalam bentuk file biner terpetakan memori (mmap) pada direktori cache pengguna sebelum dirangkum menjadi laporan.

```bash
# Lihat daftar file counter telemetri lokal di Linux
ls -la ~/.cache/go-telemetry/

# Cetak isi counter biner dalam format teks transparan
gotelemetry view
```

### 3. Pahami Mekanisme p-bit Counter untuk Sampling Efisien

Proses penghitungan metrik memanfaatkan counter atomic berukuran 64-bit yang diperbarui langsung oleh runtime compiler Go saat mengeksekusi fungsi tertentu seperti `compile/flag:v` atau `link/latency`.

```go
package main

import (
    "fmt"
    "golang.org/x/telemetry/counter"
)

func main() {
    /* Menambah counter lokal untuk memantau eksekusi fitur internal */
    counter.Inc("mytool/execution/count")
    fmt.Println("Counter telemetri lokal berhasil diperbarui.")
}
```

### 4. Pelajari Teknik Differential Privacy dan Penambahan Noise

Sebelum laporan mingguan diunggah, Go Telemetry menghitung nilai ekspresi p-bit dengan menambahkan noise acak Gaussian. Algoritma ini menjamin bahwa laporan dari seribu pengembang tidak dapat diuraikan balik untuk mengidentifikasi satu individu tertentu.

```bash
# Cetak laporan mingguan akhir yang siap diunggah ke telemetri server
gotelemetry view -json
```

### 5. Nonaktifkan Telemetri Secara Permanen Jika Diperlukan

Bila lingkungan CI/CD atau kebijakan keamanan organisasi melarang perekaman metrik lokal, pastikan mode telemetri dimatikan secara tegas.

```bash
# Matikan pengumpulan data telemetri secara total
gotelemetry off

# Verifikasi variabel lingkungan GOTELEMETRY
export GOTELEMETRY=off
```
