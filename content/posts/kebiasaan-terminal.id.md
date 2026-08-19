---
title: "Kebiasaan Kecil yang Membuat Sesi Terminal Saya Lebih Tenang"
slug: kebiasaan-terminal
aliases: []
date: 2026-08-19
author: "Daemontalk Editorial"
tags: ["linux", "tools", "productivity"]
lang: id
draft: false
type: post
cover: ""
coverCaption: ""
coverSource: ""
readTime: 4
description: "Bukan tentang plugin atau konfigurasi rumit. Hanya beberapa kebiasaan kecil yang diam-diam mengubah cara saya bekerja di terminal setiap hari."
series: ""
series_part: 0
---

Ada fase yang hampir semua orang lewati ketika mulai serius menggunakan terminal: fase konfigurasi obsesif.

Kita menghabiskan satu akhir pekan penuh memasang plugin, menulis ratusan baris alias, menyesuaikan prompt hingga terlihat persis seperti milik seseorang di Reddit. Lalu kita membuka terminal keesokan harinya dan menyadari bahwa kita masih mengetik perintah yang sama seperti sebelumnya.

Yang benar-benar mengubah cara saya bekerja bukan konfigurasi. Melainkan beberapa kebiasaan sederhana yang saya bangun perlahan.

## Membaca output dengan sabar

Dulu, ketika sebuah perintah gagal, reaksi pertama saya adalah langsung mencari solusi di Google. Salin pesan error, tempel di browser, buka link pertama, salin solusinya, tempel di terminal.

Sekarang saya memaksa diri untuk membaca output error sampai habis terlebih dahulu. Seringkali jawabannya sudah ada di sana. Linux tidak segan memberi tahu apa yang salah dengan bahasa yang cukup jelas, misalnya `permission denied`, `no such file or directory`, atau `address already in use`. Tiga detik membaca menghemat tiga menit browsing.

## Menggunakan `man` sebelum menebak-nebak

`man` adalah salah satu perintah yang paling sering diabaikan oleh pengguna baru. Padahal hampir setiap pertanyaan tentang opsi sebuah perintah bisa dijawab dengan membukanya.

Ketika saya lupa apakah `tar` menggunakan `-x` atau `-e` untuk ekstrak, jawabannya ada di `man tar`. Ketika saya tidak yakin apa yang dilakukan `grep -P`, jawabannya ada di sana. Tidak perlu membuka tab baru.

Memang butuh waktu untuk terbiasa membaca dokumentasi dalam format itu. Tapi investasinya sepadan.

## Tidak menyimpan terlalu banyak alias

Saya pernah punya file `.bashrc` dengan lebih dari delapan puluh alias. Masalahnya, saya hanya ingat sekitar dua belas di antaranya. Selebihnya hanya menambah noise ketika saya membuka file itu untuk mencari sesuatu.

Sekarang saya hanya membuat alias untuk perintah yang benar-benar saya ketik lebih dari lima kali sehari. Dan sebelum membuat alias baru, saya bertanya pada diri sendiri: apakah perintah aslinya memang susah diingat, atau saya hanya malas mengetiknya?

Kebanyakan jawabannya adalah yang kedua.

## Bekerja dengan direktori yang bersih

Ini terdengar sepele, tapi direktori home yang penuh file acak membuat saya merasa tidak tenang setiap kali membuka terminal. File-file hasil eksperimen yang tidak pernah dibersihkan, skrip lama yang sudah tidak relevan, folder `test2` dan `test2-final` yang entah berisi apa.

Sekarang saya punya aturan sederhana: setiap selesai mengerjakan sesuatu, saya bersihkan. Bukan karena ada orang yang akan melihatnya, tapi karena lingkungan yang rapi membantu saya berpikir lebih jernih.

## Menutup terminal ketika selesai

Kedengarannya konyol. Tapi dulu saya punya kebiasaan membiarkan terminal terbuka seharian, dengan puluhan tab yang masing-masing berisi sesi yang sudah tidak aktif. Berpindah-pindah antara tab itu memakan energi mental yang tidak terasa tapi nyata.

Sekarang ketika saya selesai dengan sebuah tugas, saya tutup tab yang tidak lagi diperlukan. Terminal saya hampir selalu hanya berisi satu atau dua tab aktif.

Kebiasaan-kebiasaan ini tidak mengubah produktivitas secara dramatis dalam semalam. Tapi setelah beberapa bulan, perbedaannya terasa. Sesi terminal menjadi lebih tenang, lebih fokus, dan lebih jarang frustrasi.

Mungkin itulah yang sebenarnya dicari ketika kita menghabiskan waktu mengonfigurasi lingkungan kerja kita — bukan tampilan yang keren, melainkan ketenangan.
