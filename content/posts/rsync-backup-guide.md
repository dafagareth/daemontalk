---
title: "Rsync: Transfer File Cepat dan Aman Tanpa Drama"
slug: 5a6b7c8d
aliases: [rsync-backup-guide]
date: 2026-07-08
tags: [linux, sysadmin, tips]
lang: id
draft: false
---

Mengirim file puluhan gigabyte antar-server sering gagal di tengah jalan jika hanya mengandalkan `scp` atau FTP biasa. `rsync` (Remote Sync) adalah standar emas transfer data di dunia UNIX.

Fitur *delta-transfer algorithm*-nya hanya mengirim bagian file yang berubah, menghemat bandwidth dan waktu secara drastis.

## Fun Fact

**Rsync dibuat oleh Andrew Tridgell dan Paul Mackerras pada tahun 1996.**
Andrew Tridgell juga merupakan pencipta Samba (protokol sharing file Windows di Linux).

**Algoritma Rsync menggunakan Rolling Checksum (Adler-32) + MD5/MD4.**
Rsync membagi file menjadi blok-blok kecil, menghitung checksum di kedua sisi, dan hanya mentransfer *byte offset* yang tidak cocok.

**Satu garis miring (`/`) di ujung path bisa mengubah nasib transfer.**
Di rsync, `src/` (dengan trailing slash) menyalin *isi* folder, sedangkan `src` (tanpa slash) menyalin *folder itu sendiri*.

---

## Tips dan Trik

### 1. Flag Wajib untuk Mirror Lengkap: `-avzP`

Kombinasi flag paling serbaguna:
- `-a` (archive: pertahankan permission, symlink, timestamp, user ownership)
- `-v` (verbose output)
- `-z` (kompresi data saat transmisi jaringan)
- `-P` (tampilkan progress bar & izinkan resume jika koneksi terputus)

```bash
rsync -avzP /data/backup/ user@remote-server:/var/backups/
```

### 2. Gunakan Opsi `--delete` untuk Sinkronisasi Presisi

Jika file di direktori sumber dihapus, opsi ini memastikan file serupa di sisi tujuan juga dibersihkan:

```bash
rsync -avzP --delete /var/www/site/ user@backup-server:/var/www/site/
```

### 3. Simulasi Eksekusi dengan `--dry-run` atau `-n`

Sebelum menjalankan sinkronisasi besar yang berisiko menimpa data, uji coba perintah tanpa mengubah file nyata:

```bash
rsync -avzP --dry-run /source/ /destination/
```

### 4. Batasi Penggunaan Bandwidth dengan `--bwlimit`

Jangan biarkan proses backup menyedot seluruh pipa koneksi internet kantor atau server:

```bash
# Batasi kecepatan maksimal 5000 KB/s (5 MB/s)
rsync -avzP --bwlimit=5000 /big-data/ user@remote:/storage/
```

### 5. Abaikan Direktori Tertentu dengan `--exclude`

Lewati folder build atau cache yang tidak perlu ikut disinkronkan:

```bash
rsync -avzP --exclude='.git' --exclude='node_modules' --exclude='tmp/' ./ user@vps:/app/
```
