---
title: "Mengelola Systemd di Server Linux"
slug: mengelola-systemd-linux
aliases: []
date: 2026-08-28
author: "Daemontalk Editorial"
tags: ["Linux", "Sysadmin"]
lang: id
draft: false
type: post
cover: "https://images.unsplash.com/photo-1629654297299-c8506221ca97?auto=format&fit=crop&q=80&w=1600"
coverCaption: "Linux command line interface"
coverSource: "https://unsplash.com"
readTime: 6
description: "Pelajari cara mengendalikan daemon, menulis unit service kustom untuk aplikasi Anda, serta menganalisis log secara efektif dengan journalctl."
---

**Systemd** telah mendominasi sebagai sistem inisialisasi (*init system*) default pada hampir semua distribusi Linux modern, termasuk Ubuntu, Debian, CentOS, RHEL, dan Arch Linux.

Tugas utamanya adalah menjadi "induk dari semua proses" (PID 1) yang memulai, mengawasi, dan menghentikan layanan (service/daemon) yang berjalan di latar belakang server Anda.

## Manajemen Layanan dengan systemctl

Perintah utama untuk berinteraksi dengan layanan adalah `systemctl`.

```bash
# Mengecek status terkini (menampilkan log singkat dan status aktif)
sudo systemctl status nginx

# Menjalankan, menghentikan, dan memulai ulang layanan
sudo systemctl start nginx // [!code ++]
sudo systemctl stop nginx // [!code --]
sudo systemctl restart nginx

# Memastikan layanan langsung otomatis menyala saat server restart
sudo systemctl enable nginx
```

> [!TIP]
> Jika Anda hanya mengubah konfigurasi aplikasi (seperti konfigurasi Nginx) tanpa perlu mematikan total prosesnya, gunakan perintah `sudo systemctl reload nginx`. Ini menghindari *downtime* singkat.

## Menulis Custom Service File

Jika Anda membangun aplikasi sendiri (misal *binary* Golang atau server Node.js), Anda perlu membuat file definisi `.service` agar aplikasi Anda dikelola oleh systemd.

File konfigurasi ini umumnya diletakkan di `/etc/systemd/system/myapp.service`.

```ini
[Unit]
Description=Backend API Go App
After=network.target

[Service]
Type=simple
User=appuser
WorkingDirectory=/opt/myapp
ExecStart=/opt/myapp/api-server
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Beberapa konfigurasi kunci:
- `After=network.target`: Memastikan layanan ini tidak dijalankan sebelum jaringan internet server aktif.
- `Restart=on-failure`: Sangat penting. Jika aplikasi Anda *crash* karena kesalahan *runtime* atau kehabisan memori, systemd akan otomatis menghidupkannya kembali.

> [!IMPORTANT]
> Setiap kali Anda membuat atau mengubah file `.service` secara manual, Anda wajib memberitahu systemd untuk membaca ulang konfigurasinya dengan menjalankan perintah: `sudo systemctl daemon-reload`.

## Investigasi Log dengan journalctl

Berbeda dengan sistem log lama berupa file teks (*syslog*), systemd menggunakan log biner terpusat yang disebut **journald**. Alat baca log ini adalah `journalctl`.

- **Melihat log real-time** dari layanan aplikasi kita:
  `journalctl -u myapp.service -f`
- **Melihat log berdasarkan waktu**:
  `journalctl -u myapp.service --since "1 hour ago"`
- **Melihat log dari sesi booting sebelumnya**:
  `journalctl -b -1`

> [!WARNING]
> Karena formatnya *binary*, ukuran log journald bisa membengkak drastis. Pastikan Anda mengatur batas ukuran di `/etc/systemd/journald.conf` dengan menambahkan parameter seperti `SystemMaxUse=500M`.

## Referensi

```references
- title: "systemd System and Service Manager"
  author: "Lennart Poettering"
  year: 2024
  publisher: "freedesktop.org"
  url: "https://systemd.io/"
```
