---
title: "Caddy: Web Server Modern dengan HTTPS Otomatis"
slug: d5b39c7a
aliases: caddy-automtls-praktis
date: 2026-07-15
tags: [devops, networking, tools]
lang: id
draft: false
---

Caddy adalah web server yang ditulis dalam Go dengan kemampuan manajemen sertifikat TLS otomatis. Ia dirancang agar konfigurasi minimal sudah menghasilkan HTTPS yang benar, bukan sebagai fitur tambahan yang perlu dikonfigurasi secara terpisah.

## Fakta Menarik

**Fakta 1.** Caddy menyimpan sertifikat dan kunci privat yang diperoleh dari ACME ke dalam direktori data-nya secara terenkripsi di disk. Pada restart, ia tidak perlu meminta sertifikat baru selama sertifikat yang ada belum kedaluwarsa.

**Fakta 2.** Sejak versi 2.4, Caddy mendukung ZeroSSL sebagai ACME CA default selain Let's Encrypt, dengan fallback otomatis. Ini memberikan redundansi jika salah satu CA mengalami gangguan.

**Fakta 3.** Format konfigurasi asli Caddy secara internal adalah JSON (Caddy API). Caddyfile adalah bahasa konfigurasi tingkat tinggi yang dikompilasi menjadi JSON tersebut. Anda dapat menggunakan keduanya, bahkan secara bersamaan melalui endpoint `/config/` API-nya.

---

## Tips dan Trik

### 1. Perbedaan Caddy vs Nginx dari Sisi Konfigurasi

Nginx memerlukan beberapa langkah untuk mengaktifkan HTTPS: instalasi Certbot atau acme.sh, konfigurasi blok `ssl_certificate`, dan cronjob pembaruan. Caddy menangani seluruh siklus hidup ini secara internal.

Perbandingan konfigurasi untuk reverse proxy sederhana:

**Nginx (dengan sertifikat manual):**
```nginx
server {
    listen 443 ssl;
    server_name app.example.com;

    ssl_certificate     /etc/letsencrypt/live/app.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/app.example.com/privkey.pem;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

**Caddy (Caddyfile):**
```caddy
app.example.com {
    reverse_proxy localhost:8080
}
```

Caddy secara otomatis memperoleh sertifikat, mengonfigurasi HTTP-to-HTTPS redirect, dan memperbarui sertifikat sebelum kedaluwarsa. Tidak ada langkah tambahan yang diperlukan.

### 2. Cara Kerja Automatic TLS via ACME

Saat Caddy pertama kali menerima request untuk domain yang belum memiliki sertifikat, ia menjalankan tantangan ACME (HTTP-01 atau TLS-ALPN-01) secara otomatis:

```
Browser -> Caddy (port 443) -> tantangan TLS-ALPN-01
                             -> Let's Encrypt / ZeroSSL
                             <- sertifikat diterbitkan
Caddy menyimpan sertifikat -> $CADDY_DATA_DIR/certificates/
Pembaruan otomatis ~30 hari sebelum kedaluwarsa
```

Untuk domain yang tidak dapat diakses dari internet (misalnya intranet), gunakan tantangan DNS-01 dengan plugin DNS provider yang sesuai:

```caddy
{
    acme_dns cloudflare {env.CF_API_TOKEN}
}

internal.example.com {
    reverse_proxy localhost:8080
}
```

### 3. Caddyfile Minimal untuk Aplikasi Go dan Node.js

Contoh konfigurasi untuk menjalankan dua aplikasi pada subdomain berbeda:

```caddy
# /etc/caddy/Caddyfile

# Opsi global
{
    email admin@example.com
    # Aktifkan staging Let's Encrypt untuk testing (hapus untuk produksi)
    # acme_ca https://acme-staging-v02.api.letsencrypt.org/directory
}

# Aplikasi Go di port 9000
api.example.com {
    reverse_proxy localhost:9000 {
        header_up X-Forwarded-Proto {scheme}
        header_up X-Real-IP {remote_host}
    }
    log {
        output file /var/log/caddy/api-access.log
    }
}

# Aplikasi Node.js di port 3000
app.example.com {
    reverse_proxy localhost:3000 {
        header_up X-Forwarded-Proto {scheme}
        header_up X-Real-IP {remote_host}
    }
    # Sajikan file statis langsung dari Caddy
    handle_path /static/* {
        root * /var/www/app/static
        file_server
    }
}

# Redirect www ke apex
www.example.com {
    redir https://example.com{uri} permanent
}
```

Uji validitas Caddyfile sebelum menerapkan:

```bash
caddy validate --config /etc/caddy/Caddyfile
```

### 4. Mengelola Caddy dengan systemd

Paket Caddy dari repository resmi sudah menyertakan unit file systemd. Jika menginstal secara manual, buat unit file berikut:

```ini
# /etc/systemd/system/caddy.service
[Unit]
Description=Caddy Web Server
Documentation=https://caddyserver.com/docs/
After=network.target network-online.target
Requires=network-online.target

[Service]
Type=notify
User=caddy
Group=caddy
ExecStart=/usr/bin/caddy run --environ --config /etc/caddy/Caddyfile
ExecReload=/usr/bin/caddy reload --config /etc/caddy/Caddyfile --force
TimeoutStopSec=5s
LimitNOFILE=1048576
PrivateTmp=true
ProtectSystem=full
AmbientCapabilities=CAP_NET_BIND_SERVICE

[Install]
WantedBy=multi-user.target
```

Perintah manajemen yang umum digunakan:

```bash
# Aktifkan dan jalankan Caddy
sudo systemctl enable --now caddy

# Muat ulang konfigurasi tanpa downtime (zero-downtime reload)
sudo systemctl reload caddy

# Atau gunakan perintah Caddy langsung untuk reload
sudo caddy reload --config /etc/caddy/Caddyfile

# Lihat log Caddy melalui journald
sudo journalctl -u caddy -f

# Periksa status sertifikat yang dikelola Caddy
sudo caddy trust    # tambahkan CA lokal ke trust store sistem (untuk dev)
```

### 5. Menggunakan Caddy API untuk Konfigurasi Dinamis

Selain Caddyfile, Caddy menyediakan REST API di `localhost:2019` untuk mengubah konfigurasi tanpa restart:

```bash
# Lihat konfigurasi aktif saat ini
curl -s http://localhost:2019/config/ | python3 -m json.tool

# Tambahkan route baru secara dinamis
curl -X POST http://localhost:2019/config/apps/http/servers/srv0/routes \
  -H "Content-Type: application/json" \
  -d '{
    "match": [{"host": ["newsite.example.com"]}],
    "handle": [{
      "handler": "reverse_proxy",
      "upstreams": [{"dial": "localhost:4000"}]
    }]
  }'

# Simpan konfigurasi aktif ke file
curl -s http://localhost:2019/config/ > caddy-backup.json
```

API ini berguna untuk skenario di mana konfigurasi dihasilkan secara programatik, misalnya dalam platform multi-tenant yang perlu mendaftarkan domain baru tanpa me-restart proses web server.
