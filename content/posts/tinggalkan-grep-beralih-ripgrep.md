---
title: "Tinggalkan Grep Lawas dan Beralih ke Ripgrep"
slug: "tinggalkan-grep-beralih-ripgrep"
aliases: []
author: "daemontalk team"
contributors: []
tags: ["tools", "cli", "terminal"]
lang: "id"
status: "published"
type: "dispatch"
readTime: 4
cover: "https://images.unsplash.com/photo-1526374965328-7f61d4dc18c5?auto=format&fit=crop&w=1200&q=80"
coverCaption: "Pencarian kode instan tanpa memindai node_modules"
coverSource: "https://unsplash.com"
description: "Mencari string di repositori besar dengan grep bawaan sering membekukan terminal. Ripgrep (rg) memindai kode secara instan dengan otomatis menghormati .gitignore."
series: ""
series_part: 0
---

<span class="drop-cap">M</span>engetik `grep -rn "AuthToken" .` di dalam proyek web modern adalah resep pasti untuk membekukan jendela terminal Anda selama beberapa puluh detik.

Penyebab utamanya sederhana: `grep` lawas tidak tahu apa yang sedang terjadi di ekosistem software modern. Alat ini akan dengan patuh memeriksa jutaan baris kode di dalam folder `node_modules`, file biner di `bin/`, atau riwayat kompresi di `.git/`.

Anda tentu bisa menambahkan deretan flag `--exclude-dir`, tetapi mengetik perintah panjang setiap kali ingin mencari sebaris fungsi jelas membuang waktu.

Solusi modern untuk masalah ini adalah `ripgrep` (perintah `rg`).

## Kenapa Ripgrep Jauh Lebih Kepatuhan?

Ditulis dalam bahasa Rust oleh Andrew Gallant, ripgrep dirancang dari nol untuk kecepatan pencarian kode sumber.

Alat ini secara otomatis membaca berkas `.gitignore`, `.ignore`, dan berkas tersembunyi. Folder seperti `node_modules`, `vendor/`, atau `target/` langsung dilewati tanpa perlu Anda instruksikan secara manual.

Selain itu, ripgrep memanfaatkan akselerasi instruksi SIMD di tingkat CPU dan pembagian thread paralel cerdas. Hasilnya, pencarian di repositori jutaan baris berkas sering kali selesai dalam pecahan detik.

## Perintah Harian yang Sering Dipakai

Sintaks ripgrep sangat ringkas. Untuk mencari teks, Anda cukup mengetik kata kuncinya:

```bash
rg "AuthToken"
```

Hasilnya langsung rapi: nama file, nomor baris, dan sorotan warna kata kunci tercetak jelas di terminal.

Beberapa opsi harian yang sangat berguna:

```bash
# Hanya mencari di berkas Go
rg "HandleUser" -t go

# Hanya mencari di berkas Rust atau Python
rg "connect_db" -t rust -t py

# Tampilkan 2 baris konteks sebelum dan sesudah kecocokan
rg -C 2 "database connection failed"

# Hanya cetak nama berkas yang memuat string
rg -l "deprecated_api"
```

## Menemukan Berkas dan Kolaborasi dengan Tool Lain

Jika Anda ingin mencari daftar nama berkas tanpa memindai isinya, gunakan opsi `--files`:

```bash
rg --files | rg "middleware"
```

Perintah ini jauh lebih cepat daripada menggunakan utilitas `find . -name "*middleware*"`.

Untuk operasi pencarian dan penggantian massal, gabungkan ripgrep dengan `sed`:

```bash
rg -l "OldConfigName" | xargs sed -i 's/OldConfigName/NewConfigName/g'
```

Menunggu terminal menyelesaikan pencarian adalah salah satu pengikis fokus terbesar seorang developer. Pasang ripgrep sekarang, dan jangan pernah lagi membuang waktu menunggu terminal memindai folder pihak ketiga.

```references
- title: "ripgrep User Guide and Implementation Architecture"
  author: "Andrew Gallant"
  year: 2024
  publisher: "GitHub"
  url: "https://github.com/BurntSushi/ripgrep"
```
