---
title: "Caddy Server: HTTPS Otomatis dan Reverse Proxy Tanpa Ribet"
slug: 5c6d7e8f
aliases: [caddy-vs-nginx-go]
date: 2026-06-15
tags: [networking, devops, go]
lang: id
draft: false
---

Nginx adalah web server luar biasa, namun mengatur sertifikat SSL Let's Encrypt dengan certbot, konfigurasi cron renewal, dan blok SSL boilerplate sering kali membosankan.

**Caddy**, web server open-source yang ditulis murni dalam bahasa Go, menyelesaikan masalah ini dengan *Automatic HTTPS by default*.

## Fun Fact

**Caddy dibuat oleh Matt Holt pada tahun 2015.**
Matt Holt merancang Caddy dengan satu gol utama: membuat internet terenkripsi secara otomatis tanpa beban setup manual bagi sysadmin.

**Caddy adalah web server pertama di dunia yang mendukung ACME (Let's Encrypt) secara terintegrasi.**
Begitu kamu mengetik nama domain di `Caddyfile`, Caddy otomatis meminta, memverifikasi, menginstal, dan memperbarui sertifikat TLS di latar belakang.

**Mendukung HTTP/3 (QUIC) secara out-of-the-box.**
Tanpa perlu compile manual library eksternal, Caddy langsung mendukung protokol transport UDP HTTP/3 terbaru.

---

## Tips dan Trik

### 1. Reverse Proxy ke Aplikasi Go Hanya dalam 3 Baris `Caddyfile`

Setup reverse proxy dengan SSL otomatis menjadi sesederhana ini:

```caddy
api.example.com {
    reverse_proxy localhost:8080
}
```

### 2. Sajikan File Statis dengan Kompresi Zstandard & Gzip

Optimalkan pengiriman asset CSS, JS, dan gambar:

```caddy
static.example.com {
    root * /var/www/public
    file_server
    encode zstd gzip
}
```

### 3. Load Balancing Multi-Backend dengan Healthcheck

Caddy dapat membagi beban request ke beberapa instance backend Go secara otomatis:

```caddy
app.example.com {
    reverse_proxy localhost:8001 localhost:8002 localhost:8003 {
        lb_policy round_robin
        health_uri /health
        health_interval 5s
    }
}
```

### 4. Format dan Validasi Caddyfile Sebelum Reload

Pastikan tidak ada syntax error dengan utilitas bawaan:

```bash
caddy fmt --overwrite Caddyfile
caddy validate --config Caddyfile
```

### 5. Zero-Downtime Config Reload

Perbarui konfigurasi server tanpa memutus koneksi client yang sedang aktif:

```bash
caddy reload
```
