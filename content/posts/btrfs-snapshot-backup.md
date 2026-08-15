---
title: "Snapshot Btrfs sebagai Strategi Backup Inkremental yang Efisien"
slug: a7f3c2e8
aliases: btrfs-snapshot-backup
date: 2026-03-12
tags: [linux, storage, tools]
lang: id
draft: false
---

Btrfs menyediakan mekanisme snapshot berbasis Copy-on-Write (CoW) yang memungkinkan pengambilan gambar sistem berkas pada titik waktu tertentu tanpa menyalin data secara fisik. Dikombinasikan dengan fitur `btrfs send` dan `btrfs receive`, ini menjadi fondasi strategi backup inkremental yang ringan dan efisien.

## Fakta Menarik

**Fakta 1.** Snapshot Btrfs pada awalnya tidak mengonsumsi ruang disk tambahan karena hanya menyimpan referensi ke blok data yang sama. Ruang tambahan hanya dikonsumsi saat data asli atau snapshot dimodifikasi.

**Fakta 2.** `btrfs send` menghasilkan aliran data yang merepresentasikan perbedaan antara dua snapshot. Aliran ini dapat dikirim melalui jaringan, ditulis ke file, atau langsung diteruskan ke `btrfs receive` di perangkat lain.

**Fakta 3.** Snapper, alat manajemen snapshot yang populer, mendukung Btrfs dan dapat dikonfigurasi untuk membuat snapshot otomatis sebelum setiap operasi `pacman` atau `zypper`, memungkinkan rollback sistem yang cepat.

---

## Tips dan Trik

### 1. Perbedaan Snapshot dan Salinan Biasa di Btrfs

Pada sistem berkas konvensional, menyalin direktori berarti menduplikasi seluruh data secara fisik. Di Btrfs, snapshot adalah subvolume baru yang berbagi blok data dengan subvolume sumber melalui mekanisme CoW.

```bash
# Buat subvolume sebagai titik awal
sudo btrfs subvolume create /mnt/data/dokumen

# Tambahkan beberapa file
echo "Konten penting" | sudo tee /mnt/data/dokumen/file.txt

# Lihat penggunaan disk subvolume
sudo btrfs subvolume show /mnt/data/dokumen

# Bandingkan: cp -r (menyalin semua data secara fisik)
time sudo cp -r /mnt/data/dokumen /mnt/data/dokumen-cp

# Snapshot (hampir instan, tidak menyalin data)
time sudo btrfs subvolume snapshot /mnt/data/dokumen /mnt/data/dokumen-snap
```

### 2. Membuat dan Mengelola Snapshot

```bash
# Buat snapshot read-only (direkomendasikan untuk backup)
sudo btrfs subvolume snapshot -r /mnt/data/dokumen \
  /mnt/snapshots/dokumen-$(date +%Y%m%d-%H%M%S)

# Daftar semua subvolume dan snapshot
sudo btrfs subvolume list /mnt

# Lihat detail snapshot tertentu
sudo btrfs subvolume show /mnt/snapshots/dokumen-20260312-100000

# Hapus snapshot yang sudah tidak diperlukan
sudo btrfs subvolume delete /mnt/snapshots/dokumen-20260101-120000

# Periksa penggunaan ruang yang sebenarnya
sudo btrfs filesystem du -s /mnt/snapshots/
```

### 3. Mengirim Snapshot ke Disk Eksternal dengan btrfs send/receive

Ini adalah inti dari backup inkremental Btrfs. Snapshot pertama dikirim secara penuh, snapshot berikutnya hanya mengirim selisih (delta).

```bash
# Pastikan disk eksternal juga diformat Btrfs
# sudo mkfs.btrfs /dev/sdb1
sudo mount /dev/sdb1 /mnt/backup

# Kirim snapshot pertama (penuh)
SNAP1="dokumen-20260312-100000"
sudo btrfs send /mnt/snapshots/$SNAP1 | \
  sudo btrfs receive /mnt/backup/

# Buat snapshot kedua setelah beberapa waktu
sudo btrfs subvolume snapshot -r /mnt/data/dokumen \
  /mnt/snapshots/dokumen-$(date +%Y%m%d-%H%M%S)

SNAP2="dokumen-20260312-180000"

# Kirim hanya delta antara SNAP1 dan SNAP2
sudo btrfs send -p /mnt/snapshots/$SNAP1 \
  /mnt/snapshots/$SNAP2 | \
  sudo btrfs receive /mnt/backup/

# Verifikasi snapshot di disk eksternal
sudo btrfs subvolume list /mnt/backup
```

```bash
# Kirim melalui SSH ke server remote
sudo btrfs send /mnt/snapshots/$SNAP1 | \
  ssh user@backup-server "sudo btrfs receive /backup/data/"

# Kirim delta secara inkremental melalui SSH
sudo btrfs send -p /mnt/snapshots/$SNAP1 \
  /mnt/snapshots/$SNAP2 | \
  ssh user@backup-server "sudo btrfs receive /backup/data/"
```

### 4. Integrasi dengan Snapper

Snapper mengotomasi pembuatan dan penghapusan snapshot berdasarkan jadwal atau kejadian sistem.

```bash
# Instal snapper
sudo pacman -S snapper    # Arch Linux
# sudo apt install snapper  # Debian/Ubuntu

# Buat konfigurasi snapper untuk subvolume home
sudo snapper -c home create-config /home

# Lihat konfigurasi yang ada
sudo snapper list-configs

# Edit kebijakan penyimpanan snapshot
sudo nano /etc/snapper/configs/home
```

```ini
# Bagian penting dalam /etc/snapper/configs/home

TIMELINE_CREATE="yes"
TIMELINE_CLEANUP="yes"

# Jumlah snapshot yang dipertahankan
TIMELINE_LIMIT_HOURLY="5"
TIMELINE_LIMIT_DAILY="7"
TIMELINE_LIMIT_WEEKLY="4"
TIMELINE_LIMIT_MONTHLY="6"
TIMELINE_LIMIT_YEARLY="2"
```

```bash
# Buat snapshot manual dengan deskripsi
sudo snapper -c home create --description "sebelum upgrade kernel"

# Daftar semua snapshot
sudo snapper -c home list

# Lihat perbedaan antara dua snapshot
sudo snapper -c home diff 3..5

# Kembalikan file tertentu dari snapshot
sudo snapper -c home undochange 3..0 /home/dd/.bashrc

# Aktifkan timer snapper
sudo systemctl enable --now snapper-timeline.timer
sudo systemctl enable --now snapper-cleanup.timer
```
