---
title: "Strace & Lsof: Mengintip Apa yang Sebenarnya Dilakukan Program"
slug: 9c0d1e2f
aliases: [strace-lsof-debugging]
date: 2026-05-05
tags: [linux, debugging, tools]
lang: id
draft: false
---

Ketika sebuah program binary melempar pesan error ambigu seperti "Permission denied", "File not found", atau tiba-tiba hang tanpa output, apa yang harus kamu lakukan jika kamu tidak memiliki akses ke source code-nya?

Jawabannya adalah **strace** (System Call Tracer) dan **lsof** (List Open Files). Kedua tool ini memungkinkan kamu melihat langsung interaksi aplikasi dengan Linux Kernel.

## Fun Fact

**Di sistem UNIX/Linux, segala sesuatu adalah file (*Everything is a file*).**
Bukan hanya file dokumen di disk, tetapi socket jaringan, device hardware, pipe terminal, dan memori proses diakses melalui file descriptor. Itulah mengapa `lsof` bisa menampilkan koneksi internet aktif!

**`strace` memanfaatkan syscall `ptrace` yang kuat.**
`ptrace` adalah antarmuka kernel yang sama yang digunakan debugger seperti GDB untuk menghentikan, memeriksa, dan memodifikasi memori proses lain.

**Strace bisa melacak file konfigurasi apa yang gagal dimuat.**
Banyak aplikasi memeriksa beberapa lokasi file config secara berurutan (misal: `~/.config/`, `/etc/`, `/usr/local/etc/`). Strace dapat memperlihatkan path persis yang dicari oleh binary.

---

## Tips dan Trik

### 1. Lacak File Apa Saja yang Dibuka oleh Binary

Filter hanya syscall `open` dan `openat` untuk menemukan file konfigurasi yang hilang:

```bash
strace -e trace=open,openat myapp
```

### 2. Pasang Strace ke Proses yang Sedang Berjalan Tanpa Restart

Attach ke PID proses yang sedang hang atau freeze:

```bash
sudo strace -p 1234 -s 256
```

### 3. Hitung Waktu dan Statistik Syscall dengan `-c`

Ketahui di mana program menghabiskan waktu paling lama (apakah di I/O disk, network read, atau futex locking):

```bash
strace -c myapp
```

### 4. Temukan Proses Mana yang Sedang Mengunci Port Jaringan

Jika kamu mendapat error `bind: address already in use` di port 8080:

```bash
sudo lsof -i :8080
```

### 5. Lihat Semua File dan Socket yang Dibuka oleh Suatu PID

```bash
lsof -p 1234
```
