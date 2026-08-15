1. 

---
title: "Sudah Waktunya Tinggalkan XAMPP dan Laragon"
slug: 61e38488
aliases: [docker-pengganti-xampp-laragon]
date: 2026-06-10
tags: [docker, php, laravel, development]
lang: id
draft: false
---
Hampir setiap mahasiswa teknik informatika di Indonesia memulai perjalanan web development dengan XAMPP atau Laragon. Keduanya mudah dipasang, langsung jalan, dan cukup untuk mengerjakan tugas kuliah. Tidak ada yang salah dengan itu.

Tapi ada satu masalah yang sering muncul dan jarang dibahas secara serius: lingkungan pengembangan yang tidak konsisten. Kode yang berjalan di laptop sendiri, tiba-tiba bermasalah di server kampus atau di komputer teman satu tim. Versi PHP berbeda, ekstensi tidak terpasang, konfigurasi Apache saling tumpang tindih antara satu proyek dengan proyek lain. Ini bukan masalah teknis yang sepele, ini masalah yang menghabiskan waktu.

Docker menyelesaikan masalah itu dari akarnya.

---

## Bukan Soal Keren-kerenan

Satu kesalahpahaman yang sering muncul adalah Docker dianggap sebagai alat untuk developer "level atas" atau sesuatu yang dipelajari nanti setelah lulus. Padahal justru sebaliknya.

Docker membuat kamu bisa mendefinisikan environment secara eksplisit: versi PHP berapa, ekstensi apa saja, database apa, versi berapa. Semua ditulis dalam file teks biasa dan bisa dimasukkan ke repositori Git. Ketika teman kamu clone repo yang sama dan jalankan satu perintah, mereka mendapat environment yang identik dengan milikmu. Tidak ada lagi sesi debugging "tapi di laptopku jalan".

Ini relevan sekali untuk pengerjaan tugas kelompok di kampus.

---

## Setup Dasar untuk PHP

Berikut struktur direktori yang cukup untuk sebagian besar proyek PHP mahasiswa, termasuk yang menggunakan Laravel atau CodeIgniter.

```
project/
├── docker-compose.yml
├── nginx/
│   └── default.conf
├── php/
│   └── Dockerfile
└── src/
    └── index.php
```

### `docker-compose.yml`

```yaml
services:
  app:
    build: ./php
    volumes:
      - ./src:/var/www/html
    networks:
      - dev

  webserver:
    image: nginx:alpine
    ports:
      - "8080:80"
    volumes:
      - ./src:/var/www/html
      - ./nginx/default.conf:/etc/nginx/conf.d/default.conf
    depends_on:
      - app
    networks:
      - dev

  database:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: secret
      MYSQL_DATABASE: appdb
      MYSQL_USER: appuser
      MYSQL_PASSWORD: secret
    ports:
      - "3306:3306"
    volumes:
      - dbdata:/var/lib/mysql
    networks:
      - dev

volumes:
  dbdata:

networks:
  dev:
```

### `php/Dockerfile`

```dockerfile
FROM php:8.2-fpm

RUN apt-get update && apt-get install -y \
    libpng-dev \
    libjpeg-dev \
    libfreetype6-dev \
    zip \
    unzip \
    git \
    curl

RUN docker-php-ext-configure gd --with-freetype --with-jpeg \
    && docker-php-ext-install gd pdo pdo_mysql mbstring exif pcntl bcmath

COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

WORKDIR /var/www/html
```

### `nginx/default.conf`

```nginx
server {
    listen 80;
    index index.php index.html;
    root /var/www/html/public;

    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    location ~ \.php$ {
        fastcgi_pass app:9000;
        fastcgi_index index.php;
        fastcgi_param SCRIPT_FILENAME $realpath_root$fastcgi_script_name;
        include fastcgi_params;
    }
}
```

Dengan tiga file ini, jalankan:

```bash
docker compose up -d
```

Buka `http://localhost:8080` dan environment PHP kamu sudah berjalan.

---

## Untuk Proyek Laravel

Jika kamu mengerjakan tugas berbasis Laravel, ada satu langkah tambahan setelah container berjalan.

```bash
# Masuk ke dalam container
docker compose exec app bash

# Di dalam container
composer install
cp .env.example .env
php artisan key:generate
php artisan migrate
```

Isi `.env` untuk koneksi database:

```
DB_CONNECTION=mysql
DB_HOST=database
DB_PORT=3306
DB_DATABASE=appdb
DB_USERNAME=appuser
DB_PASSWORD=secret
```

Perhatikan `DB_HOST` menggunakan nama service `database`, bukan `localhost`. Di dalam jaringan Docker, container berkomunikasi lewat nama service, bukan IP lokal.

---

## Ganti PHP Tanpa Pusing

Ini salah satu kelebihan paling nyata dibanding XAMPP. Di XAMPP, mengganti versi PHP berarti download ulang atau memodifikasi konfigurasi yang rawan rusak. Di Docker, cukup ubah satu baris di Dockerfile.

```dockerfile
# Dari ini
FROM php:8.2-fpm

# Menjadi ini
FROM php:8.1-fpm
```

Lalu rebuild:

```bash
docker compose up -d --build
```

Kalau proyek lain butuh PHP versi berbeda, tidak ada yang perlu diubah di sistem. Masing-masing proyek punya containernya sendiri, terisolasi sempurna.

---

## Banyak Proyek, Tidak Saling Ganggu

Di XAMPP, menjalankan dua proyek sekaligus sering bermasalah karena berbagi satu port 80 dan satu instalasi PHP. Di Docker, setiap proyek punya port sendiri.

Proyek A jalan di `localhost:8080`, proyek B di `localhost:8081`. Tidak ada konfigurasi virtual host yang rumit, tidak ada restart Apache setiap kali ganti proyek.

---

## Perintah yang Sering Dipakai

```bash
# Nyalakan semua service
docker compose up -d

# Matikan semua service
docker compose down

# Lihat log real-time
docker compose logs -f

# Masuk ke container PHP
docker compose exec app bash

# Jalankan artisan dari luar container
docker compose exec app php artisan migrate

# Hapus container dan volume (reset database)
docker compose down -v
```

---

## Soal Kurva Belajar

Docker memang butuh waktu untuk dipahami di awal. Konsep image, container, volume, dan network terasa asing kalau belum pernah menyentuhnya. Tapi waktu yang diinvestasikan di sini jauh lebih berguna dibanding waktu yang terbuang karena environment bermasalah setiap kali berganti mesin atau berganti anggota tim.

Lebih dari itu, Docker adalah keterampilan yang langsung relevan di dunia kerja. Hampir semua perusahaan teknologi saat ini menggunakan container di alur deployment mereka. Belajar Docker dari sekarang berarti kamu tidak perlu memulai dari nol setelah lulus.

XAMPP dan Laragon cukup untuk belajar sintaks PHP. Tapi kalau kamu sudah mengerjakan proyek yang melibatkan tim, versi PHP yang spesifik, atau deployment ke server, Docker bukan pilihan tambahan, melainkan kebutuhan dasar.
