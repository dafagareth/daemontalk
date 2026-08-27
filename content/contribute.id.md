## Pendekatan Tulisan

Fokus utama di sini adalah tulisan yang lugas, berbasis eksperimen nyata, dan dapat direproduksi. Jika Anda ingin berbagi catatan eksplorasi lab, investigasi masalah produksi, atau analisis sistem perkuliahan, tulisan Anda sangat disambut baik.

Sebaiknya langsung masuk ke inti masalah, arsitektur, atau contoh kode sejak awal. Hindari kalimat pembuka yang bertele-tele agar pembaca langsung memahami substansi teknis yang dibahas.

## Format Tulisan

Ada dua jenis format tulisan yang biasa dimuat di Daemontalk:

- **Catatan Teknis Mendalam (Deep-Dives)**: Membahas konsep sistem secara mendalam, seperti eksplorasi kernel Linux, konkurensi Go, protokol jaringan, atau arsitektur penyimpanan.
- **Catatan Insiden & Debugging (RCA)**: Menceritakan proses pencarian akar masalah dari kendala teknis nyata, langkah pelacakan log diagnostik, serta solusi mitigasinya. Gunakan tag seperti `["rca", "incident"]` atau tag sistem terkait.

## Format Frontmatter

Setiap artikel disimpan di folder `content/posts/` dalam format file Markdown dengan konfigurasi YAML di awal berkas:

```yaml
---
title: "Zero-Copy I/O dengan io_uring di Go"
slug: "7f8a9b1c"
date: "2026-08-08"
author: "Nama atau Callsign Anda"
tags: ["linux", "go", "performance", "storage"]
lang: "id"
draft: false
description: "Menjelajahi batching I/O asinkron menggunakan antarmuka syscall io_uring Linux di Go."
---
```

Gunakan `lang: "id"` untuk bahasa Indonesia dan `lang: "en"` untuk bahasa Inggris. Pastikan seluruh cuplikan kode menyertakan penanda bahasa (seperti `go`, `bash`, atau `c`) agar penyorotan sintaks berjalan optimal.

Anda dapat mengunduh file template contoh lengkap di sini: [Download template.md](/download/template.md).

## Cara Berkontribusi

**Melalui GitHub Pull Request**
1. Fork dan clone repositori ini: `git clone https://github.com/dafagareth/daemontalk`.
2. Buat branch baru untuk tulisan Anda: `git checkout -b post/topik-anda`.
3. Tambahkan file Markdown baru di `content/posts/topik-anda.md` (atau jalankan `./scripts/post.sh new --uid "Judul Anda"`).
4. Uji tampilan dan jalankan build lokal dengan perintah `make build`.
5. Buka Pull Request di GitHub dengan ringkasan singkat isi tulisan.

**Melalui Email**
Jika tidak menggunakan GitHub, Anda juga dapat mengirimkan draf Markdown atau file patch git langsung ke **realdaemontalk@gmail.com** dengan subjek judul tulisan.

## Hak Cipta dan Lisensi

Seluruh tulisan di Daemontalk diterbitkan di bawah lisensi Creative Commons Attribution-ShareAlike 4.0 International (CC BY-SA 4.0), sedangkan contoh kode berada di bawah Lisensi MIT. Hak cipta orisinal tetap menjadi milik penulis sepenuhnya.
