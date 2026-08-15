---
title: "Menjalankan Container Rootless dengan Podman"
slug: b8d2f4a6
aliases: [podman-rootless-container]
date: 2026-04-05
tags: [devops, linux, security]
lang: id
draft: false
---

Podman adalah container runtime yang dapat menjalankan container tanpa daemon dan tanpa hak akses root. Arsitekturnya yang berbeda dari Docker menjadikannya pilihan yang lebih aman untuk lingkungan produksi maupun mesin pengembang.

## Fakta Menarik

**Fakta 1.** Docker menggunakan daemon terpusat (`dockerd`) yang berjalan sebagai root. Podman tidak memiliki daemon sama sekali; setiap perintah `podman` adalah proses biasa yang berakhir setelah selesai, sehingga tidak ada layanan latar yang memegang soket root secara permanen.

**Fakta 2.** Podman menggunakan `newuidmap` dan `newgidmap` untuk memetakan UID pengguna biasa ke rentang UID di dalam namespace container, memungkinkan proses container berjalan sebagai "root" di dalam container tetapi sebagai pengguna tidak istimewa di sistem host.

**Fakta 3.** Red Hat Enterprise Linux 8 ke atas tidak menyertakan Docker sama sekali; Podman adalah container runtime resmi yang didukung penuh dan tersedia langsung dari repositori paket distribusi.

---

## Tips dan Trik

### 1. Instalasi Podman

Pada sistem berbasis Debian/Ubuntu:

```bash
sudo apt install podman
```

Pada Fedora/RHEL/CentOS Stream:

```bash
sudo dnf install podman
```

Setelah instalasi, pastikan subUID dan subGID sudah dikonfigurasi untuk pengguna Anda:

```bash
grep "$(whoami)" /etc/subuid /etc/subgid
# contoh output:
# /etc/subuid:dd:100000:65536
# /etc/subgid:dd:100000:65536
```

Jika baris tersebut belum ada, tambahkan secara manual:

```bash
sudo usermod --add-subuids 100000-165535 --add-subgids 100000-165535 "$(whoami)"
```

### 2. Menjalankan Container Rootless

Tidak diperlukan konfigurasi khusus untuk menjalankan container pertama. Podman secara otomatis menggunakan mode rootless jika dijalankan sebagai pengguna biasa:

```bash
podman run --rm -it docker.io/library/alpine:3.20 sh
```

Di dalam container, perintah `id` akan menampilkan `uid=0(root)`, tetapi di luar container proses tersebut berjalan sebagai UID pengguna biasa. Untuk memverifikasi:

```bash
# Di terminal lain, saat container sedang berjalan
ps -o pid,user,comm -p "$(pgrep -n conmon)"
```

### 3. Menjalankan Pod Multi-Container

Podman mendukung konsep pod yang kompatibel dengan Kubernetes. Berikut contoh pod dengan Nginx dan sidecar logger:

```bash
# Buat pod dengan port forwarding
podman pod create --name webpod -p 8080:80

# Jalankan container Nginx di dalam pod
podman run -d --pod webpod --name nginx docker.io/library/nginx:1.27-alpine

# Jalankan container sidecar yang membaca log Nginx
podman run -d --pod webpod --name logger \
  -v /var/log/nginx:/var/log/nginx:ro \
  docker.io/library/busybox:1.36 \
  sh -c 'tail -F /var/log/nginx/access.log 2>/dev/null || sleep infinity'

# Cek status pod
podman pod ps
```

Seluruh container dalam satu pod berbagi network namespace, sehingga mereka dapat berkomunikasi melalui `localhost`.

### 4. Integrasi dengan systemd via podman generate systemd

Podman dapat membuat unit file systemd secara otomatis untuk mengelola siklus hidup container sebagai layanan pengguna:

```bash
# Buat unit file untuk container yang sudah ada
mkdir -p ~/.config/systemd/user
podman generate systemd --new --name nginx > \
  ~/.config/systemd/user/container-nginx.service

# Aktifkan dan jalankan sebagai layanan pengguna (tanpa sudo)
systemctl --user daemon-reload
systemctl --user enable --now container-nginx.service

# Periksa status
systemctl --user status container-nginx.service
```

Opsi `--new` membuat unit file yang akan membuat ulang container dari image setiap kali layanan dimulai, bukan sekadar memulai ulang container yang ada. Ini memudahkan pembaruan image.

### 5. Menggunakan Volume dan Mengelola Image

Podman menyimpan semua data di direktori pengguna, bukan di direktori sistem:

```bash
# Lokasi penyimpanan image dan container
podman info --format '{{.Store.GraphRoot}}'
# ~/.local/share/containers/storage

# Buat volume bernama
podman volume create mydata

# Pasang volume ke container
podman run --rm -v mydata:/data alpine:3.20 sh -c 'echo "hello" > /data/test.txt'

# Baca kembali dari volume
podman run --rm -v mydata:/data alpine:3.20 cat /data/test.txt

# Hapus image yang tidak terpakai
podman image prune -f
```

Karena semua penyimpanan berada di direktori pengguna, penghapusan container rootless tidak memerlukan hak akses root sama sekali.
