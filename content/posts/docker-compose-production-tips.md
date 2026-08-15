---
title: "Docker Compose: Trik Production yang Sering Terlewatkan"
slug: 3a4b5c6d
aliases: [docker-compose-production-tips]
date: 2026-07-19
tags: [docker, devops, tips]
lang: id
draft: false
---

Banyak developer mengira Docker Compose hanya cocok untuk *development* lokal. Faktanya, untuk skala single-node VPS (aplikasi mikro, web server, database internal), Compose sering kali jauh lebih stabil dan mudah di-maintain dibanding setup Kubernetes yang berlebihan.

Berikut beberapa fakta dan teknik konfigurasi esensial untuk Docker Compose.

## Fun Fact

**Docker Compose awalnya bernama Fig.**
Fig dibuat oleh startup asal London bernama Orchard Laboratories. Docker kemudian mengakuisisi Orchard pada 2014 dan mengubah nama Fig menjadi Docker Compose.

**Docker Compose V2 ditulis ulang dari Python ke Go.**
Compose V1 lawas adalah program Python (`docker-compose` dengan tanda hubung). Mulai versi 2, Compose ditulis ulang dalam Go dan diintegrasikan langsung sebagai sub-command Docker CLI (`docker compose`).

**File `.env` otomatis dibaca tanpa flag `--env-file`.**
Docker Compose secara native akan mencari file `.env` di direktori yang sama dengan `docker-compose.yml` untuk melakukan interpolasi variabel lingkungan.

---

## Tips dan Trik

### 1. Batasi Log Driver Agar Harddisk Tidak Penuh

Docker secara default menyimpan log JSON tanpa batas ukuran. Di server production, batasi ukuran dan rotasi log di level service:

```yaml
services:
  web:
    image: myapp:latest
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

### 2. Atur Healthcheck untuk Zero-Downtime Dependency

Gunakan `depends_on` dengan kondisi `service_healthy` agar web app tidak booting sebelum database benar-benar siap menerima koneksi:

```yaml
services:
  db:
    image: postgres:16-alpine
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U $$POSTGRES_USER -d $$POSTGRES_DB"]
      interval: 5s
      timeout: 5s
      retries: 5

  api:
    image: myapi:latest
    depends_on:
      db:
        condition: service_healthy
```

### 3. Batasi Memory dan CPU per Service

Cegah satu container memakan seluruh RAM host yang memicu kernel OOM Killer:

```yaml
services:
  worker:
    image: queue-worker:latest
    deploy:
      resources:
        limits:
          cpus: '0.50'
          memory: 512M
        reservations:
          memory: 128M
```

### 4. Gunakan `restart: unless-stopped`

Opsi ini memastikan service otomatis bangkit kembali jika host VPS di-reboot atau crash mendadak, kecuali jika kamu secara eksplisit mematikannya dengan `docker compose stop`.

```yaml
services:
  nginx:
    image: nginx:alpine
    restart: unless-stopped
```

### 5. Pisahkan File Dev dan Production dengan Override

Jalankan environment yang berbeda tanpa duplikasi kode YAML:

```bash
# Production deployment
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```
