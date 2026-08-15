# 🚀 Panduan Deployment DaemonTalk ke VPS Debian (Rumahweb)

Panduan langkah demi langkah untuk melakukan deploy **daemontalk.com** pada VPS Debian menggunakan **Docker** dan **Caddy Web Server** (Automatic SSL/HTTPS).

---

## 1. Setup DNS Domain di Rumahweb / Cloudflare

Masuk ke panel DNS Domain Anda (Rumahweb / Cloudflare) dan tambahkan 2 record berikut:

| Type | Name / Host | Target / IP Address | TTL |
| :--- | :--- | :--- | :--- |
| **A** | `@` (atau kosong) | `IP_PUBLIC_VPS_ANDA` | Auto / 300 |
| **CNAME** | `www` | `daemontalk.com` | Auto / 300 |

*Tunggu 5–15 menit hingga DNS terpropagasi ke seluruh dunia.*

---

## 2. Setup Awal VPS Debian

Login ke VPS Anda via terminal:
```bash
ssh root@IP_PUBLIC_VPS_ANDA
```

Jalankan script setup otomatis yang sudah disiapkan:
```bash
curl -sL https://raw.githubusercontent.com/dafagareth/daemontalk/main/scripts/setup-vps.sh | bash
```

*Script di atas akan otomatis menginstal Docker, Docker Compose, Caddy Server, dan mengaktifkan Firewall UFW.*

---

## 3. Clone Repository & Konfigurasi Lingkungan

```bash
# Clone repository ke VPS
git clone https://github.com/dafagareth/daemontalk.git /opt/daemontalk
cd /opt/daemontalk

# Buat file .env dari template
cp .env.example .env
nano .env
```

Isi variabel penting di `.env`:
```env
PORT=8080
ENV=production
ADMIN_TOKEN=ganti_dengan_token_rahasia_anda_disini
```

---

## 4. Konfigurasi Caddy (Automatic SSL/HTTPS)

Salin `Caddyfile` DaemonTalk ke direktori Caddy sistem:
```bash
sudo cp /opt/daemontalk/Caddyfile /etc/caddy/Caddyfile
sudo systemctl reload caddy
```

*Caddy akan otomatis meminta sertifikat SSL dari Let's Encrypt / ZeroSSL untuk `daemontalk.com` dan `www.daemontalk.com`.*

---

## 5. Jalankan Aplikasi dengan Docker Compose

```bash
cd /opt/daemontalk
docker compose up -d --build
```

Periksa status container:
```bash
docker compose ps
docker compose logs -f
```

---

## 6. Update / Deployment Rutin (CI/CD atau Manual)

Jika di kemudian hari ada pembaruan kode atau artikel baru di GitHub, cukup jalankan perintah ini di VPS:
```bash
cd /opt/daemontalk
git pull origin main
docker compose up -d --build
```

---

## 7. Verifikasi Akses

- **Browser Web**: Buka `https://daemontalk.com` (cek gembok hijau SSL).
- **Terminal CLI**: Buka terminal dan jalankan:
  ```bash
  curl -sL https://daemontalk.com/daily
  ```
- **TUI Client**:
  ```bash
  go run ./cmd/tui
  ```
