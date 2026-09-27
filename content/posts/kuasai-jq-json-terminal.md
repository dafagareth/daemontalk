---
title: "Kuasai jq: Membedah JSON di Terminal Tanpa Menulis Skrip"
slug: "kuasai-jq-json-terminal"
aliases: []
author: "daemontalk team"
contributors: []
tags: ["tools", "cli", "json", "api"]
lang: "id"
status: "published"
type: "dispatch"
readTime: 4
cover: "https://images.unsplash.com/photo-1555066931-4365d14bab8c?auto=format&fit=crop&w=1200&q=80"
coverCaption: "Streaming dan transformasi data terstruktur langsung di pipeline terminal"
coverSource: "https://unsplash.com"
description: "Mengolah respons API dan berkas log JSON kompleks langsung dari command-line dengan filter, ekstraksi field, dan transformasi format tanpa skrip Python dadakan."
series: ""
series_part: 0
---

<span class="drop-cap">S</span>etiap backend engineer pernah berada di situasi ini: Anda melakukan panggilan HTTP dengan `curl` ke sebuah endpoint internal, dan terminal Anda seketika dihantam dinding teks JSON tanpa spasi sepanjang puluhan ribu karakter.

Reaksi pertama yang lazim terlihat adalah menyalin output raksasa tersebut, lalu menempelkannya ke situs web formatter JSON daring. Selain berisiko membocorkan token otentikasi atau data privat pengguna, langkah ini jelas memperlambat alur kerja harian Anda.

Alternatif lainnya adalah menulis skrip Python atau Node.js satu baris (`python -m json.tool` atau `python -c "import sys, json..."`) hanya untuk mengambil satu nilai string di dalam array bersarang.

Padahal, ekosistem Unix sudah memiliki solusi definitif untuk masalah ini: `jq`.

## Lebih dari Sekadar Pretty-Printer

`jq` sering kali hanya dipandang sebagai alat untuk memberi warna dan indentasi rapi pada JSON melalui perintah `curl -s https://api.example.com | jq .`. 

Namun, kekuatan sejati `jq` terletak pada kemampuannya bertindak sebagai bahasa pemrograman fungsional mini yang dirancang khusus untuk memotong, memfilter, dan merekonstruksi struktur data secara streaming.

> "Write programs that do one thing and do it well. Write programs to work together. Write programs to handle text streams, because that is a universal interface." — Doug McIlroy

Kekuatan ini membuat `jq` menyatu mulus dengan filosofi pipeline Unix. Anda menerima aliran byte, menyaring field yang relevan, dan meneruskannya ke utilitas lain seperti `grep`, `sort`, atau `xargs`.

## Pola Filter Harian yang Paling Sering Dibutuhkan

Berikut beberapa pola penyaringan `jq` yang paling sering menyelamatkan waktu dalam pekerjaan backend sehari-hari:

### 1. Mengambil Nilai Spesifik dari Array

Jika Anda ingin mengekstrak daftar email dari koleksi user tanpa membawa atribut lainnya:

```bash
curl -s https://api.internal/users | jq -r '.[].email'
```

Flag `-r` (raw output) menginstruksikan `jq` agar mencetak teks murni tanpa tanda kutip ganda, sehingga siap dialirkan langsung ke baris perintah berikutnya.

### 2. Memfilter Elemen Berdasarkan Kondisi Logika

Sering kali kita hanya membutuhkan data yang memenuhi status tertentu, misalnya daftar transaksi yang berstatus gagal:

```bash
cat transactions.json | jq '.[] | select(.status == "FAILED" and .amount > 500000)'
```

Fungsi `select()` menyaring setiap elemen secara deklaratif tanpa memerlukan iterasi manual atau percabangan if-else di skrip bash.

### 3. Membentuk Ulang Struktur Objek Baru

Anda bisa mengambil segelintir field penting dari objek kompleks lalu menggabungkannya ke bentuk ringkas:

```bash
cat server-metrics.json | jq '.nodes[] | {id: .node_id, ip: .network.private_ip, load: .stats.cpu}'
```

Output yang dihasilkan adalah objek JSON baru yang ramping dan siap dikonsumsi oleh modul analitik Anda.

### 4. Mengubah Objek Menjadi Format Tab/CSV

Ketika Anda ingin memasukkan data API ke spreadsheet atau mengolahnya dengan awk:

```bash
cat inventory.json | jq -r '.items[] | [.sku, .name, .price] | @tsv'
```

Operator `@tsv` atau `@csv` secara otomatis menangani escaping karakter pembatas, menghilangkan bug pemformatan yang sering muncul saat menggunakan skrip regex manual.

## Berhenti Membuat Skrip Sekali Pakai

Menulis skrip ad-hoc untuk membedah response JSON adalah bentuk over-engineering kecil yang menumpuk jadi pemborosan waktu. 

Dengan menguasai beberapa operator dasar `jq`, manipulasi data terstruktur di terminal menjadi seringkas mengetik sebaris perintah.

```references
- title: "jq Manual and Reference Documentation"
  author: "Stephen Dolan"
  year: 2024
  publisher: "GitHub Pages"
  url: "https://jqlang.github.io/jq/manual/"
```
