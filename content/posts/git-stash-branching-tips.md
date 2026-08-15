---
title: "Git Stash & Bisect: Senjata Rahasia Menyelamatkan Kode"
slug: 1c2d3e4f
aliases: [git-stash-branching-tips]
date: 2026-07-12
tags: [git, workflow, tips]
lang: id
draft: false
---

Ketika berada di tengah-tengah pengerjaan fitur dan tiba-tiba ada bug kritis yang harus segera diperbaiki di branch `main`, apa yang biasa kamu lakukan?

Git memiliki deretan sub-command canggih yang sering diabaikan tapi mampu menghemat waktu berjam-jam saat troubleshooting.

## Fun Fact

**Git diciptakan Linus Torvalds hanya dalam waktu sekitar 10 hari.**
Pada April 2005, setelah komunitas Linux dilarang memakai BitKeeper secara gratis, Linus menyendiri dan merancang core sistem Git dari awal dalam hitungan hari.

**`git bisect` menggunakan algoritma Binary Search (O(log N)).**
Jika ada bug di antara 1.000 commit terakhir, `git bisect` hanya membutuhkan sekitar 10 langkah pengujian untuk menemukan commit persis yang menjadi penyebabnya.

**Stash di Git sebenarnya adalah commit biasa yang tidak di-referensikan oleh branch.**
Saat kamu menjalankan `git stash`, Git membuat dua atau tiga commit rahasia di bawah `.git/refs/stash`.

---

## Tips dan Trik

### 1. Simpan Stash dengan Pesan Deskriptif

Jangan biarkan stash list penuh dengan tulisan `WIP on main:` yang membingungkan:

```bash
git stash push -m "refactor auth middleware sebelum hotfix"
```

### 2. Sertakan File Baru (*Untracked Files*) Saat Stash

Secara default, file baru yang belum pernah di-`git add` tidak akan ikut tersimpan dalam stash. Tambahkan flag `-u`:

```bash
git stash -u
```

### 3. Otomatisasi Pencarian Bug dengan `git bisect run`

Jika kamu punya script test otomatis (misalnya unit test yang gagal jika bug ada), biarkan Git mencari commit yang rusak secara otomatis:

```bash
git bisect start
git bisect bad                 # commit sekarang rusak
git bisect good v1.4.0         # commit rilis lalu masih aman
git bisect run go test ./...   # biarkan Git mencari otomatis!
```

### 4. Buka Branch Baru Langsung dari Stash

Jika kode yang kamu simpan di stash ternyata terlalu besar dan butuh branch mandiri:

```bash
git stash branch feature-auth-v2 stash@{0}
```

### 5. Intip Isi Perubahan Stash Tanpa Me-restore

Gunakan `git stash show -p` untuk melihat patch diff:

```bash
git stash show -p stash@{0}
```
