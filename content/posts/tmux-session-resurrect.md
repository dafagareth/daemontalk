---
title: "Tmux Persistensi: Menyelamatkan Workspace Saat Server Reboot"
slug: 3b4c5d6e
aliases: [tmux-session-resurrect]
date: 2026-04-28
tags: [tmux, terminal, workflow]
lang: id
draft: false
---

Tmux sangat hebat dalam menjaga sesi terminal tetap berjalan saat koneksi SSH terputus. Namun, ketika laptop atau VPS di-reboot, seluruh layout jendela (*panes*), direktori kerja, dan sesi Tmux yang kamu susun dengan rapi biasanya akan lenyap.

Dengan plugin persistensi `tmux-resurrect` dan `tmux-continuum`, kamu bisa mengembalikan workspace kerja 100% seperti semula setelah restart.

## Fun Fact

**Tmux ditulis dalam C oleh Nicholas Marriott pada tahun 2007.**
Nicholas merancangnya sebagai alternatif modern yang lebih bersih dan modular dibanding GNU Screen yang sudah berusia puluhan tahun.

**Tmux menggunakan arsitektur Client-Server terpisah.**
Server Tmux mengelola semua state jendela dan pseudo-terminal (pty), sementara klien Tmux hanyalah perantara input/output tampilan.

**TPM (Tmux Plugin Manager) dibuat oleh Bruno Sutic pada 2014.**
TPM memungkinkan instalasi plugin Tmux semudah menekan shortcut keyboard `Prefix + I`.

---

## Tips dan Trik

### 1. Pasang Tmux Plugin Manager (TPM)

Langkah awal untuk memasang ekosistem plugin Tmux:

```bash
git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm
```

### 2. Konfigurasi `tmux-resurrect` dan `tmux-continuum` di `~/.tmux.conf`

Tambahkan baris berikut agar sesi disimpan secara berkala setiap 15 menit dan dipulihkan otomatis:

```tmux
# ~/.tmux.conf
set -g @plugin 'tmux-plugins/tpm'
set -g @plugin 'tmux-plugins/tmux-sensible'
set -g @plugin 'tmux-plugins/tmux-resurrect'
set -g @plugin 'tmux-plugins/tmux-continuum'

# Auto-restore sesi terakhir saat tmux server dijalankan
set -g @continuum-restore 'on'

# Pulihkan juga program seperti Neovim
set -g @resurrect-strategy-nvim 'session'

# Inisialisasi TPM (harus di baris paling bawah)
run '~/.tmux/plugins/tpm/tpm'
```

### 3. Simpan dan Restore Manual Kapan Saja

Di dalam Tmux:
- `Prefix + Ctrl+s` : Simpan state semua sesi dan pane saat ini
- `Prefix + Ctrl+r` : Restore kembali seluruh state sesi yang tersimpan

### 4. Aktifkan Mouse Scrolling & Resizing

Jangan biarkan scrolling terhambat:

```tmux
set -g mouse on
```

### 5. Atur Base Index Jendela Mulai dari Angka 1

Tombol angka `1` di keyboard jauh lebih dekat ke tangan kiri dibanding angka `0`:

```tmux
set -g base-index 1
setw -g pane-base-index 1
```
