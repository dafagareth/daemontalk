---
title: "Query Lambat Biasanya Bukan Salah Database, Tapi Lupa Bikin Index"
slug: "query-lambat-lupa-bikin-index"
aliases: []
date: 2026-09-24
author: "daemontalk team"
contributors: []
tags: ["craft", "database", "backend", "performance"]
lang: "id"
status: "published"
type: "dispatch"
readTime: 4
cover: "https://images.unsplash.com/photo-1558494949-ef010cbdcc31?auto=format&fit=crop&w=1200&q=80"
coverCaption: "Tabel besar tanpa index memaksa mesin membaca setiap blok memori dari nol"
coverSource: "https://unsplash.com"
description: "Saat aplikasi mulai lambat di server, reaksi pertama developer sering menyalahkan database atau berencana memasang Redis. Padahal masalahnya hampir selalu sama: tabel dipaksa membaca setiap baris dari awal karena lupa index."
series: ""
series_part: 0
---

<span class="drop-cap">S</span>ebuah endpoint API yang diuji di laptop terasa sangat kencang. Waktu responnya hanya 3 milidetik. Namun dua bulan kemudian, ketika aplikasi diunggah ke server dan mulai diisi data transaksi oleh pengguna nyata, endpoint yang sama mendadak butuh waktu 2,5 detik untuk memuat satu halaman profil.

Reaksi pertama yang lazim muncul di kepala programmer pemula biasanya seragam: database relasional dianggap lambat, server kurang RAM, atau arsitekturnya harus segera dipasangi Redis untuk *caching*.

Padahal jika query tersebut diperiksa langsung di konsol SQL, biang keladinya sangat sederhana: mesin database dipaksa membaca ratusan ribu baris data satu per satu hanya untuk mencari satu baris milik Anda.

## Membaca Buku Tanpa Indeks Halaman

Bayangkan Anda memegang buku kamus setebal 1.000 halaman, tetapi seluruh katanya disusun acak tanpa urutan abjad. Jika seseorang meminta Anda mencari arti kata "kernel", apa yang harus Anda lakukan?

Anda terpaksa membuka lembar pertama, memeriksa setiap kata, lalu pindah ke lembar kedua, ketiga, sampai lembar terakhir. Proses melelahkan ini di dunia database disebut **Full Table Scan** (atau *Sequential Scan*).

Saat data di tabel `orders` baru berisi 20 baris selama tahap pengembangan, mesin database menyelesaikan pembacaan seluruh tabel dalam sekejap mata. Di laptop, perbedaan antara membaca 20 baris dengan membaca 1 baris sama-sama terasa instan.

Masalah baru meledak saat tabel tersebut tumbuh menjadi 200.000 baris. Tanpa petunjuk apa pun, query sederhana seperti ini:

```sql
SELECT id, total_amount, created_at
FROM orders
WHERE user_id = 4821;
```

akan memaksa mesin database mengangkat 200.000 baris data dari disk ke memori, mencocokkan nilai `user_id` satu per satu, lalu membuang 199.990 baris sisanya.

## Melihat Bukti dengan EXPLAIN

Sebelum menebak-nebak kenapa query terasa lambat, biasakan menggunakan perintah `EXPLAIN` atau `EXPLAIN ANALYZE` di PostgreSQL, MySQL, atau SQLite.

Perintah ini memperlihatkan rencana eksekusi (*query execution plan*) yang dibuat oleh query planner database sebelum data diambil.

Tanpa index pada kolom `user_id`, output Postgres akan terlihat seperti ini:

```text
Seq Scan on orders  (cost=0.00..4250.00 rows=12 width=32)
  Filter: (user_id = 4821)
  Rows Removed by Filter: 199988
Execution Time: 38.452 ms
```

Perhatikan baris `Seq Scan` dan angka `Rows Removed by Filter: 199988`. Database membaca seluruh baris hanya untuk menemukan 12 transaksi milik pengguna tersebut.

Sekarang tambahkan satu baris perintah index:

```sql
CREATE INDEX idx_orders_user_id ON orders(user_id);
```

Jalankan kembali query yang sama dengan `EXPLAIN ANALYZE`:

```text
Index Scan using idx_orders_user_id on orders  (cost=0.29..14.30 rows=12 width=32)
  Index Cond: (user_id = 4821)
Execution Time: 0.084 ms
```

Dari 38 milidetik turun menjadi 0,08 milidetik. Lebih dari 400 kali lipat lebih cepat, hanya dengan satu baris DDL.

Index bekerja seperti indeks di halaman belakang buku teks. Database menyusun nilai `user_id` ke dalam struktur pohon seimbang (kebanyakan menggunakan B-Tree). Mesin cukup melakukan beberapa lompatan pencarian biner untuk langsung menemukan lokasi fisik baris data di storage.

## Kolom Mana Saja yang Perlu Diberi Index?

Memberi index bukan berarti Anda harus membuat index di semua kolom tabel. Jadikan tiga kriteria ini sebagai patokan utama:

1. **Kolom pada klausa WHERE yang sering dipanggil.** Misalnya kolom `user_id` pada tabel transaksi, atau kolom `status` jika query sering memfilter data aktif.
2. **Kolom penghubung antar-tabel (Foreign Key untuk JOIN).** Kolom yang digunakan pada klausa `ON a.order_id = b.id`.
3. **Kolom untuk pengurutan berulang (ORDER BY).** Jika halaman utama selalu menampilkan data berdasarkan `created_at DESC`, index pada kolom tersebut mencegah database melakukan pengurutan ulang di memori (*in-memory sort*) setiap kali request masuk.

## Biaya yang Harus Dibayar

Index bukan fitur gratis. Setiap index membutuhkan ruang penyimpanan tambahan di disk. 

Selain itu, setiap kali aplikasi menjalankan operasi tulis (`INSERT`, `UPDATE`, atau `DELETE`), mesin database tidak hanya mengubah data tabel utama, tetapi juga harus memperbarui struktur pohon index yang terpasang pada tabel tersebut.

Jika sebuah tabel memiliki sepuluh index berbeda, satu operasi `INSERT` berarti satu kali tulis ke tabel dan sepuluh kali tulis ke index. Karena itu, buatlah index berdasarkan pola query yang benar-benar dipanggil aplikasi, bukan berdasarkan spekulasi di masa depan.

Sebelum memikirkan arsitektur caching yang rumit atau memperbesar kapasitas CPU server, periksa dulu log query aplikasi Anda. Sering kali, masalah performa backend selesai hanya dengan satu baris `CREATE INDEX`.
