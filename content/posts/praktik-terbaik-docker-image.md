---
title: "Praktik Terbaik Membangun Docker Image"
slug: praktik-terbaik-docker-image
aliases: []
date: 2026-08-28
author: "Daemontalk Editorial"
tags: ["DevOps", "Docker"]
lang: id
draft: false
type: post
cover: "https://images.unsplash.com/photo-1605745341112-85968b19335b?auto=format&fit=crop&q=80&w=1600"
coverCaption: "Shipping containers at a port"
coverSource: "https://unsplash.com"
readTime: 6
description: "Panduan mendalam untuk mengoptimalkan ukuran image, memanfaatkan sistem caching layer, dan meningkatkan keamanan lingkungan Docker."
---

Membangun image Docker yang efisien adalah kunci untuk proses *Continuous Integration* (CI) yang cepat, *deployment* yang stabil, dan keamanan yang terjamin di lingkungan produksi.

Image yang berukuran besar tidak hanya memakan waktu jaringan untuk di-*pull*, tetapi juga meningkatkan *attack surface* (area permukaan serangan) karena memuat *library* OS yang tidak diperlukan.

## 1. Gunakan file `.dockerignore`

Langkah pertama sebelum menulis Dockerfile adalah membuat `.dockerignore`. Fungsinya mirip `.gitignore`. Ini mencegah Docker menyalin file-file seperti `.git/`, `node_modules/`, atau *log* lokal ke dalam *build context*.

```text
# Contoh .dockerignore
.git
.env
node_modules/
vendor/
*.log
```

Mengabaikan direktori besar membuat proses *build* awal jauh lebih cepat karena Docker tidak perlu mentransfer data tersebut ke *daemon*.

## 2. Multi-Stage Builds

Gunakan multi-stage build untuk memisahkan *environment* build dengan *environment* runtime.

```dockerfile
# Stage 1: Build (menggunakan image lengkap dengan compiler)
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go // [!code hl]

# Stage 2: Runtime (image sangat kecil dan polos)
FROM scratch // [!code ++]
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

> [!TIP]
> Pada contoh di atas, kita menggunakan `scratch` (image kosong sama sekali) atau `alpine` pada stage kedua. Hasilnya, ukuran image bisa turun drastis dari ~800MB menjadi ~15MB.

## 3. Optimalkan Layer Caching

Docker mengeksekusi instruksi dari atas ke bawah dan melakukan *cache* pada setiap layernya. Jika sebuah layer berubah, semua layer di bawahnya harus di-*build* ulang.

Oleh karena itu, salin file *dependencies* terlebih dahulu sebelum menyalin seluruh kode sumber (seperti contoh Golang di atas, atau `package.json` untuk Node.js). 

```dockerfile
# URUTAN YANG BENAR:
COPY package.json package-lock.json ./
RUN npm install
# Jika source code berubah, npm install TIDAK dijalankan ulang
COPY . .
```

## 4. Hindari Menjalankan sebagai Root

Secara default, aplikasi di dalam container Docker berjalan dengan hak akses *root*. Hal ini berisiko tinggi jika ada kerentanan *Remote Code Execution* (RCE) pada aplikasi Anda. 

Gunakan pengguna khusus (non-root):

```dockerfile
# Alpine: Membuat grup dan user non-root
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser
```

> [!IMPORTANT]
> Selalu terapkan prinsip *least privilege*. Pastikan *user* tersebut hanya memiliki izin *read/write* pada direktori yang mutlak diperlukan saja.

## 5. Pemindaian Keamanan (Vulnerability Scanning)

Selalu pindai image Anda sebelum di-*push* ke *registry* (seperti Docker Hub atau ECR). Alat seperti **Trivy** dapat disisipkan ke dalam *pipeline* CI/CD untuk menggagalkan *build* jika ditemukan celah kritis (CVE).

## Referensi Terverifikasi

```references
- title: "Best practices for writing Dockerfiles"
  author: "Docker Docs"
  year: 2024
  publisher: "Docker Inc."
  url: "https://docs.docker.com/develop/develop-images/dockerfile_best-practices/"
```
