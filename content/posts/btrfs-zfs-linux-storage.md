---
title: "Btrfs & ZFS: Fitur Snapshot dan CoW untuk Developer"
slug: 9e0f1a2b
aliases: [btrfs-zfs-linux-storage]
date: 2026-07-02
tags: [linux, storage, tips]
lang: id
draft: false
---

Filesystem tradisional seperti `ext4` sangat stabil, namun filesystem generasi baru berbasis **Copy-on-Write (CoW)** seperti **Btrfs** dan **ZFS** menawarkan fitur canggih: snapshot instan, kompresi transparan, dan proteksi dari silent data corruption.

Bagi developer, ini berarti kamu bisa menduplikasi database 100GB dalam waktu 0,1 detik tanpa memakan ruang disk tambahan.

## Fun Fact

**ZFS awalnya dikembangkan oleh Sun Microsystems untuk Solaris OS.**
ZFS dirancang oleh Jeff Bonwick dan Matt Ahrens pada 2001 dengan arsitektur 128-bit, kapasitas teoritis yang mampu menampung seluruh data di muka bumi tanpa batas kapasitas.

**Btrfs (B-tree FS) diusulkan oleh insinyur Oracle Chris Mason pada 2007.**
Fedora dan OpenSUSE telah menjadikan Btrfs sebagai filesystem default untuk instalasi desktop dan workstation mereka.

**Snapshot CoW tidak menyalin data fisik saat dibuat.**
Snapshot hanya membuat referensi metadata baru ke blok data yang sudah ada. Disk hanya terpakai ketika ada blok data lama yang dimodifikasi.

---

## Tips dan Trik

### 1. Duplikasi File Instan dengan `cp --reflink=always`

Di filesystem Btrfs atau XFS dengan reflink enabled, kamu bisa menduplikasi file besar secara instan (*zero-copy*):

```bash
cp --reflink=always database_production.db database_test.db
```

### 2. Buat Subvolume dan Snapshot Instan di Btrfs

Ambil snapshot root sistem sebelum melakukan upgrade paket besar:

```bash
# Buat snapshot read-only
sudo btrfs subvolume snapshot -r / /snapshots/root_before_upgrade
```

### 3. Aktifkan Kompresi Transparan `zstd` di `/etc/fstab`

Hemat ruang penyimpanan SSD hingga 30-50% tanpa penurunan performa:

```text
# /etc/fstab
UUID=xxxx-xxxx / btrfs rw,noatime,compress=zstd:1,subvol=@ 0 0
```

### 4. Periksa Integritas Data dengan Btrfs Scrub

Jalankan background scrub untuk memverifikasi checksum blok dan memperbaiki kerusakan data secara otomatis jika memakai mode RAID1:

```bash
sudo btrfs scrub start /
sudo btrfs scrub status /
```

### 5. Pantau Ruang Bebas Btrfs yang Akurat

Perintah `df -h` standar sering keliru membaca sisa disk pada filesystem CoW. Gunakan perintah bawaan:

```bash
sudo btrfs filesystem usage /
```
