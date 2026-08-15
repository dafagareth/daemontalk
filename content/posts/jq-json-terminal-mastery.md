---
title: "Jq Power User: Memfilter JSON Raksasa Langsung di CLI"
slug: 3e4f5a6b
aliases: [jq-json-terminal-mastery]
date: 2026-05-30
tags: [cli, tools, terminal]
lang: id
draft: false
---

Ketika berhadapan dengan respons API berukuran megabyte atau dump log database berformat JSON, membuka file di GUI editor sering membuat komputer hang. `jq` adalah prosesor JSON baris perintah yang luar biasa fleksibel, setara dengan `sed` dan `awk` khusus data terstruktur.

Berikut trik filtering yang akan menghemat waktu inspeksi data.

## Fun Fact

**Jq diciptakan oleh Stephen Dolan pada tahun 2012.**
Stephen Dolan adalah peneliti ilmu komputer dan kontributor inti bahasa OCaml. Sintaks `jq` dirancang seperti bahasa pemrograman fungsional murni.

**Jq mendukung streaming parser untuk file gigabyte.**
Dengan flag `--stream`, `jq` mampu memproses file JSON multi-gigabyte tanpa harus memuat seluruh struktur dokumen ke dalam RAM sekaligus.

**Output Jq bisa diubah ke format CSV, TSV, atau raw text.**
Kamu bisa mengubah response REST API kompleks menjadi tabel CSV siap olah di spreadsheet hanya dengan satu baris command.

---

## Tips dan Trik

### 1. Ekstrak Array of Objects Menjadi Plain Table / TSV

Ambil hanya kolom tertentu (misal: ID dan Email) tanpa tanda kutip ganda:

```bash
curl -s https://api.example.com/users | jq -r '.[] | [.id, .email, .role] | @tsv'
```

### 2. Filter Berdasarkan Kondisi (*Conditional Query*)

Cari data item yang memenuhi kriteria spesifik (misal: status aktif dan skor > 80):

```bash
jq '.data[] | select(.active == true and .score > 80) | {name: .name, score: .score}' records.json
```

### 3. Modifikasi Nilai Key Tanpa Merusak Struktur Lain

Gunakan update operator `|=` untuk mengubah field tertentu:

```bash
jq '.config.timeout |= 60 | .config.retries |= 5' settings.json
```

### 4. Ekstrak Semua Keys atau Tipe Data dari Dokumen Asing

Membantu memahami schema data JSON yang tidak memiliki dokumentasi:

```bash
jq 'keys' payload.json
jq 'map(type) | unique' payload.json
```

### 5. Format JSON Mentah Langsung di Clipboard

Kombinasikan dengan clipboard utility (`xclip` atau `pbcopy`):

```bash
# Pretty-print JSON di clipboard
xclip -o | jq . | xclip -selection clipboard
```
