---
title: "Kodinganmu Terlalu Rapi Sampai Lupa Rilis"
slug: "kodingan-terlalu-rapi-lupa-rilis"
aliases: []
author: "daemontalk team"
contributors: []
tags: ["essays", "opinion", "culture"]
lang: "id"
status: "published"
type: "dispatch"
readTime: 5
cover: "https://images.unsplash.com/photo-1507238691740-187a5b1d37b8?auto=format&fit=crop&w=1200&q=80"
coverCaption: "Arsitektur berlapis tanpa pengguna yang memakainya"
coverSource: "https://unsplash.com"
description: "Obsesi pada abstraksi dini dan aturan Clean Code dogmatis sering kali menjadi pelarian programmer dari kenyataan bisnis: kodenya sempurna, tapi produknya kehabisan waktu."
series: ""
series_part: 0
---

<span class="drop-cap">P</span>ernahkah Anda memeriksa sebuah pull request yang tugasnya cuma menambahkan satu kolom nomor telepon di formulir registrasi, tetapi menyentuh delapan belas berkas berbeda?

Di dalamnya ada berkas Controller, UseCase, RepositoryInterface, RepositoryImpl, Dto, RequestValidator, dan dua berkas Mapper terpisah.

Kodenya luar biasa rapi. Semua interface terpasang presisi, mock test lulus seratus persen, dan pohon direktori proyek terlihat megah seperti buku teks arsitektur software.

Masalahnya cuma satu: tiket itu butuh waktu pengerjaan tiga pekan. Saat fiturnya akhirnya meluncur ke produksi, pengguna ternyata tidak pernah mengisi kolom tersebut.

## Jebakan Abstraksi Dini

Banyak developer merasa berdosa jika menulis kode yang sederhana dan langsung pada intinya.

Setiap kali membaca buku arsitektur, timbul dorongan bawah sadar untuk mengantisipasi masa depan yang belum tentu terjadi. Kita sibuk menyiapkan skenario: bagaimana kalau nanti database Postgres diganti MongoDB, atau bagaimana kalau protokol HTTP diganti gRPC?

Kenyataan di lapangan jauh lebih membosankan: dalam sembilan puluh sembilan persen kasus produk digital, database utama tidak pernah diganti sampai perusahaannya pivot atau kehabisan modal.

Abstraksi yang dibuat terlalu dini bukan tanda kematangan rekayasa, melainkan beban kognitif yang dipinjam dari masa depan tanpa bunga yang jelas.

## Modul Dangkal dan Beban Pikiran

Dalam bukunya mengenai filosofi desain software, John Ousterhout menyinggung konsep yang ia sebut sebagai modul dangkal (shallow modules).

Modul dangkal adalah kelas atau fungsi yang antarmukanya lebar, tetapi pekerjaan di dalamnya nyaris tidak ada. Kelas tersebut cuma berfungsi sebagai perantara yang meneruskan argumen ke kelas lain di bawahnya.

Modul jenis ini tidak menyelesaikan kompleksitas. Ia hanya menyembunyikan logika dan memindahkan beban ingatan ke kepala developer yang membacanya.

Ketika terjadi kesalahan di produksi, Anda terpaksa membuka enam tab editor sekaligus hanya untuk menelusuri satu pemanggilan query SQL yang sebenarnya bisa ditulis dalam sepuluh baris.

## Kode yang Bagus adalah Kode yang Gampang Dihapus

Pada fase awal pengembangan produk, hal paling berharga bukanlah kesempurnaan struktur kode, melainkan kecepatan memvalidasi kebutuhan pengguna.

Fitur yang Anda bangun minggu ini punya peluang lima puluh persen untuk dibuang bulan depan karena pasar tidak menginginkannya.

Kent Beck, pelopor metodologi Extreme Programming dan Test-Driven Development, memiliki adagium terkenal yang sering kita lupakan:

> "Make it work, make it right, make it fast."
> (Kent Beck)

Ironisnya, banyak developer hari ini membalik urutan tersebut: sibuk menjadikannya *right* dan *fast* dengan arsitektur berlapis, padahal pembuktian bahwa solusinya *work* bagi pengguna saja belum selesai.

Dan Abramov pernah menulis pengakuan serupa dalam esainya: dorongan untuk membersihkan kode secara berlebihan sering kali didorong oleh rasa bangga semu seorang programmer, bukan kebutuhan nyata dari proyek yang sedang dikerjakan.

Jangan jadikan kerapian folder sebagai tempat persembunyian dari ketakutan merilis produk ke dunia nyata. Tulis kode yang jelas, selesaikan masalah pengguna hari ini, dan rapikan kembali arsitekturnya saat produk Anda sudah membuktikan nilainya.

```references
- title: "A Philosophy of Software Design"
  author: "John Ousterhout"
  year: 2018
  publisher: "Yaknyam Press"
  url: "https://web.stanford.edu/~ouster/cgi-bin/book.php"

- title: "Goodbye, Clean Code"
  author: "Dan Abramov"
  year: 2020
  publisher: "Overreacted"
  url: "https://overreacted.io/goodbye-clean-code/"
```
