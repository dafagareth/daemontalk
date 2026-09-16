---
title: "Berhenti Mengandalkan Git Stash dengan Git Worktree"
slug: "git-worktree-tanpa-stash"
aliases: []
author: "daemontalk team"
contributors: []
tags: ["tools", "git", "workflow"]
lang: "id"
status: "published"
type: "dispatch"
readTime: 4
cover: "https://images.unsplash.com/photo-1618401471353-b98aedd04e11?auto=format&fit=crop&w=1200&q=80"
coverCaption: "Bekerja di beberapa branch paralel tanpa konflik stash"
coverSource: "https://unsplash.com"
description: "Bolak-balik git stash saat ada hotfix mendadak sering memicu konflik dan merusak konteks kerja. Git worktree memungkinkan Anda membuka branch terpisah di direktori berbeda secara paralel."
series: ""
series_part: 0
---

<span class="drop-cap">S</span>kenario klasik di tengah hari kerja: Anda sedang sibuk merefaktor arsitektur database di sebuah branch fitur dengan belasan berkas belum di-commit. Tiba-tiba, sebuah bug kritis di lingkungan produksi memaksa Anda memperbaikinya detik itu juga.

Kebiasaan lama sebagian besar developer adalah mengetik:

```bash
git stash
git checkout main
# perbaiki bug, commit, lalu push
git checkout fitur-database
git stash pop
```

Lalu bencana kecil pun terjadi: Anda disambut rentetan konflik merge di layar terminal. Kode setengah matang yang tadinya Anda ingat mendadak harus diselesaikan bersamaan dengan perubahan branch main.

Stash Anda akhirnya menumpuk di laci tersembunyi (`stash@{0}`, `stash@{1}`) dan sering kali terbengkalai berbulan-bulan.

Ada cara yang jauh lebih bersih: `git worktree`.

## Satu Repositori, Banyak Direktori Kerja

Fitur bawaan Git ini memungkinkan satu repositori lokal memiliki beberapa direktori kerja (working tree) yang terhubung ke satu folder `.git` yang sama.

Alih-alih mengacak-acak branch Anda yang sekarang, Anda cukup membuat direktori terpisah untuk menangani hotfix:

```bash
git worktree add ../proyek-hotfix -b hotfix-login main
```

Perintah di atas membuat folder baru bernama `proyek-hotfix` di luar folder kerja utama Anda, sekaligus membuat branch baru bernama `hotfix-login` dari branch `main`.

Buka terminal atau editor baru di direktori `../proyek-hotfix`. Perbaiki bug di sana, jalankan tes, lalu commit dan push seperti biasa.

Direktori kerja utama Anda tidak tersentuh sama sekali. Seluruh perubahan yang belum di-commit tetap aman di tempat asalnya tanpa perlu di-stash.

## Membersihkan Worktree yang Selesai

Setelah perbaikan di-merge dan tugas selesai, Anda bisa melihat daftar worktree yang aktif dengan perintah:

```bash
git worktree list
```

Untuk menghapus direktori kerja yang sudah tidak digunakan, cukup jalankan:

```bash
git worktree remove ../proyek-hotfix
```

Folder kerja akan dihapus secara otomatis dan repositori lokal Anda kembali rapi.

## Kapan Harus Menggunakan Worktree?

Teknik ini sangat berguna pada situasi:

- Menjalankan test suite atau proses build yang memakan waktu lama di satu branch sambil tetap menulis kode di branch lain.
- Membandingkan antarmuka atau perilaku dua branch berbeda secara berdampingan di peramban tanpa harus restart server lokal berulang kali.
- Melakukan tinjauan kode (code review) PR rekan tim secara lokal tanpa mengorbankan status direktori kerja Anda saat ini.

Git stash dirancang untuk menyimpan jeda sesaat, bukan untuk menopang alur kerja paralel. Gunakan `git worktree` dan biarkan setiap tugas berjalan di jalurnya masing-masing.

```references
- title: "git-worktree Documentation"
  author: "Git Core Team"
  year: 2024
  publisher: "Git SCM"
  url: "https://git-scm.com/docs/git-worktree"
```
