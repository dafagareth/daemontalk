---
title: "Btop & Htop: Monitor Resource Linux dengan Elegan"
slug: a1b2c3d4
aliases: [btop-htop-monitor-linux]
date: 2026-07-28
tags: [linux, cli, tips]
lang: id
draft: false
---

Memantau performa server atau laptop Linux tidak harus membosankan dengan tampilan `top` bawaan yang kaku. Distro Linux modern punya deretan TUI (Terminal User Interface) monitor yang interaktif dan kaya visual.

Dua alat yang paling populer adalah `htop` dan `btop`. Berikut hal menarik dan tips penggunaannya di terminal harian.

## Fun Fact

**`top` asli sudah ada sejak tahun 1984.**
Perintah `top` pertama kali ditulis oleh William LeFebvre untuk BSD UNIX. Desainnya yang sederhana dibuat agar muat pada layar terminal beresolusi 80x24 karakter.

**`htop` ditulis ulang dalam C oleh Hisham Muhammad.**
Hisham menciptakan `htop` pada 2004 karena frustrasi `top` tidak mendukung scroll horizontal, pemilihan proses dengan kursor mouse, dan visualisasi bar CPU multi-core.

**`btop` adalah evolusi dari `bpytop` dan `bashtop` (ya, Bash!).**
Developer asal Swedia, Aristocratos, awalnya membuat tool ini murni menggunakan script Bash (`bashtop`). Setelah lambat, ia porting ke Python (`bpytop`), dan akhirnya ke C++ murni (`btop`) untuk efisiensi CPU mendekati 0%.

**Btop mendukung visualisasi GPU NVIDIA dan AMD secara native.**
Tanpa perlu tool terpisah seperti `nvidia-smi`, btop langsung bisa membaca utilisasi VRAM, clock speed, dan watt GPU di layar monitor yang sama.

---

## Tips dan Trik

### 1. Navigasi Tree Process di Htop dengan Tombol `F5` atau `t`

Melihat proses dalam struktur hirarki (parent-child) memudahkan identifikasi sub-proses atau worker pool yang memakan banyak memori.

```bash
htop
# Tekan F5 atau huruf 't' untuk toggle Tree View
```

### 2. Custom Layout dan Tema Warna di Btop

Btop menyediakan puluhan tema bawaan (Nord, Dracula, Gruvbox, Tokyo Night) langsung dari menu options tanpa perlu edit file config manual.

```bash
# Buka btop, lalu tekan ESC atau tombol 'm'
# Masuk ke menu Options -> Color theme
```

### 3. Filter Proses Spesifik dengan Tombol `F4` atau `/`

Alih-alih mencari manual dengan scroll ratusan proses, tekan `/` di `htop` atau `f` di `btop` untuk mengetik keyword nama aplikasi.

```bash
# Di htop:
# Tekan F4, ketik "nginx" atau "postgres"
```

### 4. Kirim Signal SIGKILL (`kill -9`) Tanpa Keluar dari Monitor

Pilih proses yang hang atau zombie, lalu tekan tombol `F9` (htop) atau `k` (btop). Pilih signal yang diinginkan (SIGTERM 15 untuk gracefully stop, atau SIGKILL 9 untuk paksa terminate).

```bash
# Shortcut htop:
# Sorot baris proses -> F9 -> pilih 9 (SIGKILL) -> Enter
```

### 5. Jalankan Htop Hanya untuk User Tertentu

Jika server dipakai bersama tim, pantau hanya proses milik akunmu agar layar tidak penuh:

```bash
htop -u $USER
```

### 6. Simpan Konfigurasi Custom Btop di `~/.config/btop/btop.conf`

Kamu bisa mengubah update rate (default 2000ms menjadi 1000ms agar lebih responsif) atau mematikan box network/disk jika hanya butuh monitor CPU/RAM.

```ini
# ~/.config/btop/btop.conf
update_ms = 1000
shown_boxes = "cpu mem proc"
```
