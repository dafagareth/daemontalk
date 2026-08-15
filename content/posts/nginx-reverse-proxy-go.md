---
title: "Setup Nginx sebagai Reverse Proxy untuk Aplikasi Go"
slug: bcb16548
aliases: [nginx-reverse-proxy-go]
date: 2025-10-05
tags: [nginx, go, devops, linux]
lang: id
draft: false
---

Aplikasi Go bisa langsung mendengarkan di port 80 atau 443. Secara teknis tidak ada yang menghalangi. Tapi menempatkan Nginx di depannya adalah praktik yang benar karena beberapa alasan: Nginx menangani terminasi SSL, kompresi, caching aset statis, dan rate limiting jauh lebih efisien, sementara proses Go fokus pada logika aplikasi.

## Konfigurasi Minimal

Buat file di `/etc/nginx/sites-available/daemontalk.com`:

```nginx
server {
    listen 80;
    server_name daemontalk.com www.daemontalk.com;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Aktifkan dan test:

```bash
ln -s /etc/nginx/sites-available/daemontalk.com /etc/nginx/sites-enabled/
nginx -t
systemctl reload nginx
```

## Header yang Penting

Tiga header di atas bukan hiasan. Tanpa `X-Real-IP`, aplikasi Go kamu akan melihat semua request datang dari `127.0.0.1`, yaitu Nginx itu sendiri. Log akses jadi tidak berguna, rate limiting berdasarkan IP tidak bisa bekerja.

Di sisi Go, baca IP asli dari header tersebut:

```go
func realIP(r *http.Request) string {
    if ip := r.Header.Get("X-Real-IP"); ip != "" {
        return ip
    }
    return r.RemoteAddr
}
```

`X-Forwarded-Proto` memberitahu aplikasi apakah request datang lewat HTTP atau HTTPS, berguna ketika kamu perlu membangun URL absolut atau mengatur cookie dengan flag `Secure`.

## SSL dengan Let's Encrypt

Pasang Certbot:

```bash
apt install certbot python3-certbot-nginx
certbot --nginx -d daemontalk.com -d www.daemontalk.com
```

Certbot otomatis mengubah konfigurasi Nginx untuk menangani HTTPS dan menambahkan redirect dari HTTP. Setelah selesai, konfigurasinya akan terlihat seperti ini:

```nginx
server {
    listen 443 ssl;
    server_name daemontalk.com www.daemontalk.com;

    ssl_certificate /etc/letsencrypt/live/daemontalk.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/daemontalk.com/privkey.pem;
    include /etc/letsencrypt/options-ssl-nginx.conf;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

server {
    listen 80;
    server_name daemontalk.com www.daemontalk.com;
    return 301 https://$host$request_uri;
}
```

Sertifikat diperbarui otomatis lewat systemd timer yang dipasang Certbot.

## Aset Statis Langsung dari Nginx

Kalau aplikasi Go kamu melayani file statis (CSS, JS, gambar), Nginx bisa menanganinya langsung tanpa meneruskan request ke Go. Ini lebih cepat:

```nginx
location /static/ {
    alias /home/deploy/portfolio/web/static/;
    expires 30d;
    add_header Cache-Control "public, immutable";
    gzip_static on;
}

location / {
    proxy_pass http://127.0.0.1:8080;
    # ... header lainnya
}
```

Pastikan direktori `/home/deploy/portfolio/web/static/` bisa dibaca oleh user `www-data` (user yang menjalankan Nginx).

## Timeout dan Buffer

Untuk endpoint yang butuh waktu lebih lama (misalnya operasi database berat atau generasi gambar):

```nginx
location / {
    proxy_pass http://127.0.0.1:8080;
    proxy_read_timeout 30s;
    proxy_connect_timeout 5s;
    proxy_send_timeout 30s;

    proxy_buffering on;
    proxy_buffer_size 4k;
    proxy_buffers 8 4k;
}
```

Default `proxy_read_timeout` Nginx adalah 60 detik. Sesuaikan dengan kebutuhan aplikasimu.

## Verifikasi

Setelah konfigurasi aktif, cek bahwa semuanya berjalan:

```bash
nginx -t                        # validasi konfigurasi
curl -I https://daemontalk.com  # cek header respons
journalctl -u nginx --since "5 minutes ago"  # lihat log
```

Dengan setup ini, Nginx menangani koneksi SSL dan aset statis, sementara proses Go hanya menerima request yang relevan dari localhost.
