---
title: "Simpan Gambar di Database atau di Folder Server?"
slug: "simpan-gambar-database-atau-folder-server"
aliases: []
date: 2026-09-24
author: "daemontalk team"
contributors: []
tags: ["craft", "backend", "architecture", "database"]
lang: "id"
status: "published"
type: "dispatch"
readTime: 4
cover: "https://images.unsplash.com/photo-1555066931-4365d14bab8c?auto=format&fit=crop&w=1200&q=80"
coverCaption: "Penyimpanan data relasional dan berkas biner memiliki kebutuhan infrastruktur yang berbeda"
coverSource: "https://unsplash.com"
description: "Dilema klasik saat membuat fitur unggah foto profil: menyimpan berkas biner di kolom database atau menaruhnya di folder server? Ini perbandingan teknis dan alasan pemisahan keduanya."
series: ""
series_part: 0
---

<span class="drop-cap">F</span>ormulir pendaftaran pengguna baru saja selesai dirangkai. Input teks seperti nama, email, dan kata sandi sudah aman tersimpan ke tabel. Sekarang giliran fitur unggah foto profil. Di titik ini, banyak pengembang backend pemula berhenti sejenak dan bertanya: berkas gambarnya harus disimpan di mana?

Pilihan pertama yang sering terlintas adalah menyimpannya langsung ke database. Bukankah database mendukung kolom tipe `BLOB` (*Binary Large Object*) atau string teks panjang yang bisa menampung representasi Base64?

Pilihan kedua adalah menyimpan berkas fisik gambar di folder server (misalnya direktori `/var/www/uploads/`), lalu database hanya mencatat alamat teks lokasinya.

Di atas kertas, pilihan pertama tampak rapi. Dalam praktik nyata di server produksi, pilihan kedua hampir selalu menjadi standar industri.

## Daya Tarik Menyimpan Gambar di Kolom Database

Alasan utama yang membuat pemula menyukai opsi menyimpan file di database adalah kenyamanan pengelolaan:

1. **Semua data berada di satu tempat.** Saat membuat salinan cadangan (*backup* database) dengan `pg_dump` atau `mysqldump`, seluruh foto pengguna otomatis ikut tercadangkan dalam satu berkas arsip.
2. **Transaksi atomik terlindungi.** Jika proses pendaftaran akun gagal di tengah jalan, seluruh baris (termasuk gambarnya) otomatis batal lewat mekanisme *rollback*. Tidak ada risiko file gambar yatim piatu tertinggal di hard disk.
3. **Tidak ada urusan izin folder.** Pengembang tidak perlu repot mengatur hak akses direktori (*file permissions*) seperti `chmod` atau `chown` di sistem operasi Linux.

Namun, kenyamanan awal ini datang dengan harga operasional yang sangat mahal begitu aplikasi mulai digunakan banyak orang.

## Kenapa Ukuran Database Bisa Cepat Membengkak

Database relasional seperti PostgreSQL, MySQL, dan SQLite dioptimalkan untuk mengolah data terstruktur yang berukuran kecil dan seragam: angka, tanggal, dan teks pendek.

Ketika Anda memasukkan berkas biner berukuran 2 MB ke dalam kolom tabel, beberapa masalah langsung muncul:

- **Ukuran backup membengkak tak terkendali.** Basis data dengan 50.000 pengguna yang hanya menyimpan teks biasanya berukuran di bawah 100 MB. Proses backup selesai dalam hitungan detik. Jika setiap pengguna memiliki avatar 2 MB di tabel yang sama, ukuran database melonjak menjadi 100 GB. Membuat backup harian mendadak menjadi beban berat bagi CPU dan disk I/O server.
- **Pemborosan memori cache database.** Database mengalokasikan RAM khusus (*buffer pool*) untuk menyimpan baris data yang sering diakses agar query berjalan cepat. Ketika query memuat baris data yang berisi byte gambar berukuran besar, buffer pool tersebut akan cepat sesak, sehingga baris data relasional lain terdorong keluar dari RAM.
- **Beban koneksi jaringan backend.** Setiap kali aplikasi memanggil `SELECT * FROM users`, jaringan antara server web dan server database dipaksa mentransfer megabyte data biner yang belum tentu dibutuhkan saat itu.

## Pola yang Benar: Simpan Lokasi, Bukan Fisiknya

Solusi yang terbukti tahan lama adalah memisahkan tanggung jawab:

```text
[Browser / Klien]
       |
       | 1. Kirim multipart form (file foto)
       v
[Backend API]
       |
       +---> 2. Tulis berkas fisik ke disk:  /uploads/avatars/u-984.jpg
       |
       +---> 3. Simpan string path ke SQL:   INSERT INTO users (avatar_url) 
                                             VALUES ('/uploads/avatars/u-984.jpg');
```

Database hanya menyimpan sebaris string sederhana berisi path atau URL publik:

```sql
SELECT id, username, avatar_path 
FROM users 
WHERE id = 42;
```

Ukuran kolom `avatar_path` hanya sekitar 30 sampai 50 byte. Query tetap ringan, buffer pool tetap hemat, dan tabel tidak terbebani oleh ukuran biner gambar.

Untuk menyajikan gambar tersebut ke browser pengguna, serahkan tugasnya ke web server (seperti Nginx atau Caddy) atau handler file statis bawaan aplikasi. Web server modern menggunakan pemanggilan sistem operasi seperti `sendfile` di Linux, yang menyalurkan berkas langsung dari disk ke kartu jaringan (*zero-copy*) tanpa perlu memproses data di dalam memori aplikasi Anda.

## Kapan Membutuhkan Object Storage?

Menyimpan berkas di folder lokal server adalah titik awal terbaik untuk aplikasi dengan satu server. 

Namun ketika aplikasi Anda berkembang dan mulai berjalan di beberapa server di balik load balancer, atau menggunakan container Docker yang sifat penyimpanannya sementara (*ephemeral*), folder lokal tidak lagi cukup. Gambar yang diunggah ke Server A tidak akan ada di Server B.

Di tahap itulah Anda mulai memindahkan berkas fisik ke layanan **Object Storage** (seperti Amazon S3, Cloudflare R2, atau MinIO). Pola databasenya tetap sama persis: database tetap hanya menyimpan URL string gambar, sementara berkas fisiknya ditangani oleh infrastruktur penyimpanan berkas yang terpisah.

Gunakan database untuk apa yang menjadi keahliannya: mencari relasi, menyaring data, dan menjaga integritas tabel. Untuk berkas gambar, biarkan disk atau object storage yang menyimpannya.
