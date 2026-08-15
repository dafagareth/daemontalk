---
title: "tmux: Satu Terminal, Banyak Sesi yang Tidak Pernah Hilang"
slug: 879fd99d
aliases: [tmux-terminal-multiplexer]
date: 2025-01-27
tags: [tmux, terminal, productivity]
lang: id
draft: false
---

Pernah menjalankan proses panjang di server lewat SSH, lalu koneksi terputus, dan seluruh proses ikut mati? Atau membuka lima tab terminal untuk menjalankan server, melihat log, dan mengedit file secara bersamaan, lalu kehilangan jejak mana yang mana?

`tmux` menyelesaikan kedua masalah ini. Ia adalah terminal multiplexer: satu terminal yang bisa membagi diri menjadi banyak panel, banyak window, dan yang terpenting, tetap hidup meski koneksi kamu putus.

## Konsep Dasar

tmux punya tiga tingkatan:

- **session**: sebuah ruang kerja yang menampung semuanya
- **window**: seperti tab di dalam session
- **pane**: pembagian layar di dalam satu window

Satu session bisa punya banyak window, dan satu window bisa dibagi menjadi banyak pane. Semuanya berjalan di server tmux yang terpisah dari terminal kamu. Inilah kenapa proses tetap berjalan meski terminalnya ditutup.

Mulai session baru:

```bash
tmux
```

Atau dengan nama agar mudah dikenali nanti:

```bash
tmux new -s kerja
```

## Prefix Key

Semua perintah tmux diawali dengan kombinasi **prefix**, yang secara default adalah `Ctrl-b`. Kamu menekan `Ctrl-b`, melepaskannya, lalu menekan tombol perintah.

Beberapa yang paling sering dipakai:

```
Ctrl-b c     buat window baru
Ctrl-b n     pindah ke window berikutnya
Ctrl-b p     pindah ke window sebelumnya
Ctrl-b 0-9   loncat ke window nomor tertentu
Ctrl-b ,     ganti nama window
Ctrl-b %     bagi pane secara vertikal
Ctrl-b "     bagi pane secara horizontal
Ctrl-b o     pindah ke pane berikutnya
Ctrl-b x     tutup pane saat ini
```

Daftar ini terlihat banyak, tapi dalam praktik kamu hanya akan sering memakai lima atau enam. Sisanya menyusul dengan sendirinya.

## Fitur Pembunuh: Detach dan Attach

Inilah alasan utama orang memakai tmux. Kamu bisa melepaskan diri dari session tanpa menghentikannya, lalu menyambung kembali kapan saja, bahkan dari mesin lain.

```
Ctrl-b d
```

Perintah ini melakukan **detach**. Session beserta semua proses di dalamnya tetap berjalan di latar belakang. Terminal kamu kembali normal. Kamu bisa menutup terminal, mematikan laptop, atau kehilangan koneksi SSH, dan session itu tetap utuh.

Untuk menyambung kembali:

```bash
# Lihat session yang sedang berjalan
tmux ls

# Sambung ke session bernama "kerja"
tmux attach -t kerja
```

Bayangkan skenario ini: kamu menjalankan proses build atau migrasi database yang makan waktu satu jam di server. Tanpa tmux, kamu harus menjaga koneksi SSH tetap hidup selama itu. Dengan tmux, kamu detach, tutup laptop, pulang, lalu attach lagi dari rumah untuk melihat hasilnya. Prosesnya tidak pernah terganggu.

## Membagi Layar untuk Alur Kerja Nyata

Pane membuat satu window menjadi ruang kerja lengkap. Pola yang umum saat mengembangkan aplikasi:

```bash
# Mulai dengan satu pane, lalu bagi
# Ctrl-b %  → bagi jadi kiri-kanan
# Ctrl-b "  → bagi pane kanan jadi atas-bawah
```

Hasilnya bisa berupa: pane kiri untuk editor, pane kanan atas menjalankan server, pane kanan bawah menampilkan log. Semua terlihat sekaligus dalam satu layar, tanpa berpindah-pindah tab.

Untuk berpindah antar pane dengan lebih nyaman, banyak orang menambahkan navigasi ala vim ke konfigurasi.

## Konfigurasi yang Membuat tmux Nyaman

Default tmux fungsional tapi kaku. File `~/.tmux.conf` memperbaikinya. Berikut konfigurasi minimal yang sudah meningkatkan pengalaman secara signifikan:

```bash
# ~/.tmux.conf

# Ganti prefix ke Ctrl-a, lebih mudah dijangkau
unbind C-b
set -g prefix C-a
bind C-a send-prefix

# Mulai penomoran window dari 1, bukan 0
set -g base-index 1

# Navigasi pane ala vim
bind h select-pane -L
bind j select-pane -D
bind k select-pane -U
bind l select-pane -R

# Split dengan tombol yang lebih intuitif
bind | split-window -h
bind - split-window -v

# Aktifkan mouse untuk klik dan resize pane
set -g mouse on

# Perbesar history scrollback
set -g history-limit 10000
```

Setelah mengubah file, muat ulang tanpa keluar dari tmux:

```
Ctrl-b :
source-file ~/.tmux.conf
```

Banyak orang mengganti prefix ke `Ctrl-a` karena lebih mudah dijangkau oleh ibu jari dan jari telunjuk dibanding `Ctrl-b`. Mengaktifkan mouse juga menghilangkan friksi bagi yang belum hafal semua shortcut.

## Mode Scroll dan Copy

Secara default, scroll dengan mouse di tmux berperilaku aneh. Masuk ke copy mode untuk menggulir riwayat dan menyalin teks:

```
Ctrl-b [     masuk copy mode, lalu gulir dengan panah atau PgUp
q            keluar dari copy mode
```

Di copy mode kamu bisa menggulir ke atas untuk melihat output lama yang sudah lewat dari layar, berguna saat mencari pesan error yang muncul beberapa menit lalu.

---

tmux butuh beberapa hari untuk terbiasa, terutama menghafal prefix dan beberapa shortcut inti. Tapi imbalannya sepadan: proses yang tidak pernah mati karena koneksi putus, ruang kerja yang tertata dalam satu layar, dan kemampuan berpindah mesin tanpa kehilangan konteks. Bagi siapa pun yang banyak bekerja di terminal, terutama lewat SSH ke server, tmux dengan cepat berubah dari hal yang dipelajari menjadi hal yang tidak bisa ditinggalkan.
