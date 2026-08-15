---
title: "UFW & Iptables: Mengamankan VPS Linux dari Serangan Brute Force"
slug: 5e6f7a8b
aliases: [ufw-iptables-firewall-tips]
date: 2026-05-12
tags: [linux, security, sysadmin]
lang: id
draft: false
---

Setiap IP publik VPS baru yang aktif di internet akan diserang ribuan percobaan brute-force bot otomatis dalam hitungan menit. Membiarkan semua port terbuka adalah undangan terbuka bagi attacker.

**UFW (Uncomplicated Firewall)** dan backend packet filter Linux (`iptables` / `nftables`) menyediakan dinding pertahanan pertama yang tangguh namun mudah diatur.

## Fun Fact

**UFW dikembangkan oleh Canonical untuk sistem operasi Ubuntu.**
UFW dibuat pada 2008 sebagai antarmuka ramah pemula yang membungkus sintaks `iptables` yang terkenal rumit dan membingungkan.

**Iptables kini telah digantikan oleh `nftables` di kernel modern.**
Mulai Linux kernel 3.13, `nftables` menyederhanakan arsitektur pemrosesan paket jaringan dengan bytecode engine yang jauh lebih cepat. Perintah `ufw` modern otomatis mengelola rules di atas `nftables`.

**Rule default firewall yang aman adalah "Deny All Incoming, Allow All Outgoing".**
Hanya port yang secara eksplisit kamu izinkan (seperti SSH, HTTP, HTTPS) yang boleh menerima paket data dari luar.

---

## Tips dan Trik

### 1. Jangan Kunci Dirimu Sendiri! Izinkan SSH Terlebih Dahulu

Urutan eksekusi sangat krusial saat menyalakan UFW untuk pertama kali:

```bash
# 1. Atur default policy
sudo ufw default deny incoming
sudo ufw default allow outgoing

# 2. Izinkan port SSH (jangan sampai lupa!)
sudo ufw allow 22/tcp

# 3. Baru aktifkan firewall
sudo ufw enable
```

### 2. Batasi Akses SSH Hanya dari Alamat IP Tertentu

Jika kamu memiliki static IP atau koneksi VPN:

```bash
sudo ufw allow from 203.0.113.50 to any port 22 proto tcp
```

### 3. Aktifkan Rate Limiting untuk Mencegah Brute Force

UFW punya fitur `limit` bawaan yang otomatis memblokir IP jika melakukan lebih dari 6 koneksi dalam 30 detik:

```bash
sudo ufw limit ssh
```

### 4. Izinkan Port Web Server Standar Sekaligus

Gunakan profil aplikasi bawaan:

```bash
sudo ufw allow "Nginx Full" # Membuka port 80 (HTTP) dan 443 (HTTPS)
```

### 5. Hapus Rule Berdasarkan Nomor Indeks

Hindari salah hapus dengan melihat urutan nomor rule:

```bash
sudo ufw status numbered
sudo ufw delete 3 # Menghapus rule nomor 3
```
