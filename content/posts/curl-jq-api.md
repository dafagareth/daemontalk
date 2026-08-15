---
title: "curl dan jq: Menguji API Tanpa Membuka Postman"
slug: cdf722f3
aliases: [curl-jq-api]
date: 2025-03-09
tags: [curl, api, tools]
lang: id
cover: /static/images/curl.png
draft: false
---

Saat mengembangkan atau menguji API, banyak orang langsung membuka Postman atau Insomnia. Keduanya alat yang baik, tapi untuk pengujian cepat, membuka aplikasi GUI, membuat request, dan mengatur tab terasa berlebihan. Sering kali yang kamu butuhkan hanya satu baris di terminal.

`curl` mengirim request, `jq` mengolah jawaban JSON-nya. Kombinasi keduanya cukup untuk sebagian besar pengujian API sehari-hari, dan punya keunggulan yang tidak dimiliki GUI: mudah disalin, dibagikan, dan dimasukkan ke dalam script.

## Request Paling Dasar

```bash
curl https://api.contoh.com/users
```

Tanpa opsi tambahan, `curl` melakukan GET dan menampilkan body respons apa adanya. Untuk API yang mengembalikan JSON, outputnya akan berupa satu baris panjang yang sulit dibaca. Di sinilah `jq` masuk.

```bash
curl -s https://api.contoh.com/users | jq
```

`jq` tanpa argumen merapikan dan mewarnai JSON. Opsi `-s` pada curl ("silent") menyembunyikan progress bar agar tidak mengganggu pipeline. Dua tambahan kecil ini sudah membuat output jauh lebih nyaman dibaca.

## Mengambil Bagian Tertentu dengan jq

Kekuatan sebenarnya `jq` adalah menyaring data. Anggap responsnya seperti ini:

```json
{
  "users": [
    { "id": 1, "name": "Dafa", "active": true },
    { "id": 2, "name": "Sari", "active": false }
  ]
}
```

Ambil hanya array users:

```bash
curl -s https://api.contoh.com/users | jq '.users'
```

Ambil nama setiap user:

```bash
curl -s https://api.contoh.com/users | jq '.users[].name'
```

Saring hanya user yang aktif:

```bash
curl -s https://api.contoh.com/users | jq '.users[] | select(.active == true)'
```

Bentuk ulang data menjadi struktur baru:

```bash
curl -s https://api.contoh.com/users | jq '.users[] | {nama: .name, aktif: .active}'
```

Sintaks `jq` perlu sedikit pembiasaan, tapi pola dasarnya konsisten: titik untuk mengakses field, `[]` untuk menelusuri array, `select()` untuk menyaring, dan `|` untuk merangkai operasi.

## POST dengan Body JSON

Mengirim data sama mudahnya. Gunakan `-X` untuk metode dan `-d` untuk body.

```bash
curl -s -X POST https://api.contoh.com/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Budi", "email": "budi@contoh.com"}' | jq
```

Header `Content-Type` penting karena memberi tahu server bahwa body yang dikirim berformat JSON. Tanpa itu, banyak server salah menafsirkan data dan mengembalikan error.

Untuk body yang panjang, lebih rapi membacanya dari file:

```bash
curl -s -X POST https://api.contoh.com/users \
  -H "Content-Type: application/json" \
  -d @user-baru.json | jq
```

Tanda `@` memberi tahu curl untuk membaca isi dari file `user-baru.json`.

## Autentikasi

Sebagian besar API butuh token. Kirim lewat header Authorization.

```bash
curl -s https://api.contoh.com/profil \
  -H "Authorization: Bearer eyJhbGc..." | jq
```

Agar token tidak berserakan di riwayat shell, simpan di variabel environment:

```bash
export TOKEN="eyJhbGc..."
curl -s https://api.contoh.com/profil \
  -H "Authorization: Bearer $TOKEN" | jq
```

Cara ini juga membuat perintah lebih mudah disalin dan dibagikan tanpa membocorkan token asli.

## Melihat yang Tidak Terlihat

Kadang masalah bukan di body, melainkan di status code atau header. Opsi `-i` menyertakan header respons, sementara `-v` menampilkan seluruh percakapan termasuk request yang dikirim.

```bash
# Tampilkan header respons bersama body
curl -i https://api.contoh.com/users

# Mode verbose: lihat request dan respons lengkap
curl -v https://api.contoh.com/users
```

Untuk sekadar mengecek status code tanpa body, ini trik yang berguna:

```bash
curl -s -o /dev/null -w "%{http_code}\n" https://api.contoh.com/users
```

`-o /dev/null` membuang body, dan `-w "%{http_code}"` mencetak hanya status code. Berguna saat kamu hanya ingin memastikan endpoint mengembalikan 200, bukan 404 atau 500.

## Menyusun Menjadi Script

Karena semuanya berupa perintah teks, mudah merangkainya menjadi pemeriksaan otomatis.

```bash
#!/bin/bash
set -euo pipefail

base="https://api.contoh.com"

# Ambil token sekali
token=$(curl -s -X POST "$base/login" \
  -H "Content-Type: application/json" \
  -d '{"email": "tes@contoh.com", "password": "rahasia"}' \
  | jq -r '.token')

# Pakai token untuk request berikutnya
curl -s "$base/profil" \
  -H "Authorization: Bearer $token" \
  | jq '.name'
```

Perhatikan `jq -r`. Opsi `-r` ("raw") mengeluarkan string tanpa tanda kutip, sehingga nilai token bisa langsung dipakai di perintah berikutnya. Tanpa `-r`, hasilnya akan menyertakan tanda kutip JSON yang merepotkan.

---

Postman tetap berguna untuk koleksi request yang kompleks dan kerja tim. Tapi untuk memeriksa apakah sebuah endpoint berfungsi, melihat bentuk respons, atau menyaring satu nilai dari JSON yang besar, `curl` dan `jq` lebih cepat dan tidak meminta kamu meninggalkan terminal. Keduanya juga keterampilan yang terbawa ke mana-mana: ada di hampir setiap server Linux, bisa dimasukkan ke script, dan tidak pernah berubah antarversi. Sekali terbiasa, membuka aplikasi GUI untuk pengujian sederhana mulai terasa seperti langkah yang tidak perlu.
