---
title: "Evolusi Garbage Collection di Go 1.20+: Minimalisasi Latensi Stop-the-World"
slug: 9f8e7d6c
aliases: []
date: 2026-08-11
tags: [go]
lang: id
draft: false
type: post
cover: ""
---

## Abstrak

Bahasa pemrograman Go, sejak awal rancangannya, menitikberatkan pada efisiensi eksekusi dan konkurensi skala besar. Salah satu komponen krusial yang menentukan performa aplikasi Go adalah mekanisme pengumpulan sampah atau garbage collection (GC). Pada versi-versi terbaru, khususnya mulai dari Go 1.20 dan seterusnya, telah terjadi pergeseran paradigma signifikan dalam cara runtime Go menangani alokasi memori. Artikel ini mengkaji evolusi GC pada Go, dengan fokus utama pada teknik pacer baru yang diperkenalkan untuk meminimalisasi latensi stop-the-world, serta dampaknya terhadap aplikasi dengan throughput tinggi.

## Tinjauan Pustaka: Mekanisme Dasar GC pada Go

Sebelum membahas evolusi pada versi 1.20+, penting untuk memahami landasan GC di Go. Go menggunakan algoritma concurrent mark-and-sweep yang dirancang untuk beroperasi beriringan dengan program utama (mutator). Proses ini secara garis besar terdiri dari fase penandaan objek yang masih digunakan dan fase pembersihan objek yang sudah tidak terjangkau. Meskipun berjalan secara konkuren, GC di Go tetap membutuhkan fase stop-the-world (STW) yang singkat untuk sinkronisasi, seperti saat inisialisasi fase penandaan dan saat penyelesaian.

Latensi STW ini selalu menjadi momok bagi aplikasi real-time. Jika durasi STW terlalu lama, aplikasi akan mengalami jeda yang terlihat oleh pengguna. Oleh karena itu, tim pengembang Go secara konsisten berusaha menekan angka latensi ini hingga berada di bawah satu milidetik, bahkan untuk tumpukan memori berukuran gigabyte.

## Paradigma Pacing Konvensional

Pada versi sebelum 1.18, GC pacer pada Go memiliki heuristik yang relatif sederhana. Pacer adalah subsistem yang memutuskan kapan sebuah siklus GC harus dimulai. Keputusan ini didasarkan pada rasio antara jumlah memori yang baru dialokasikan dengan jumlah memori yang bertahan dari siklus GC sebelumnya, yang diatur melalui variabel lingkungan GOGC.

Kelemahan utama dari pendekatan ini adalah ketidakmampuannya beradaptasi secara optimal terhadap beban kerja yang fluktuatif. Pada aplikasi yang memiliki laju alokasi memori yang sangat dinamis, pacer lama seringkali memicu siklus GC terlalu cepat atau terlalu lambat. Jika terlalu cepat, CPU akan terbuang untuk proses GC. Jika terlalu lambat, ukuran heap akan membengkak dan memperburuk latensi saat siklus GC akhirnya berjalan.

## Evolusi Pacer pada Go 1.20+

Pembaruan arsitektur pacer yang diperkenalkan secara bertahap dan dimatangkan pada Go 1.20 menawarkan solusi yang lebih elegan. Pacer baru ini tidak lagi hanya bergantung pada metrik statis, melainkan menggunakan model kontrol adaptif yang terinspirasi dari teori kontrol. Model ini secara kontinu memantau laju alokasi mutator dan laju penandaan oleh thread GC.

### Kontrol Berbasis Umpan Balik

Dengan pendekatan berbasis umpan balik, pacer baru dapat memprediksi lebih akurat kapan memori akan mencapai batas yang ditentukan. Hal ini memungkinkan runtime untuk memulai siklus GC pada momen yang paling tepat, sehingga jumlah pekerjaan yang harus dilakukan saat STW dapat ditekan seminimal mungkin. Prediktabilitas ini sangat penting untuk aplikasi yang memproses transaksi finansial atau layanan telekomunikasi yang membutuhkan Service Level Agreement (SLA) yang ketat.

### Dampak terhadap Latensi Stop-the-World

Pengukuran empiris menunjukkan bahwa pacer baru berhasil menstabilkan latensi ekor (tail latency) secara signifikan. Pada persentil ke-99, aplikasi yang bermigrasi ke Go 1.20 seringkali melaporkan penurunan durasi STW hingga puluhan persen dibandingkan versi sebelumnya. Pengurangan ini terjadi karena pacer yang lebih cerdas mampu meratakan beban kerja GC, mencegah terjadinya penumpukan pekerjaan sinkronisasi yang memicu STW panjang.

## Pengelolaan Memori Eksternal dan Arenas

Selain peningkatan pada pacer, wacana mengenai pengelolaan memori secara manual melalui eksperimen arena juga menjadi topik hangat di komunitas akademis Go. Arena memory management memungkinkan alokasi sekumpulan objek secara berdekatan yang kemudian dapat didealokasikan secara massal. Meskipun masih berstatus eksperimental, konsep ini sejalan dengan tujuan meminimalisasi tekanan pada GC. Dengan mengurangi jumlah objek individual yang harus dipindai oleh GC, overhead penandaan dan STW dapat dikurangi lebih jauh.

## Implikasi Praktis bagi Pengembang Perangkat Lunak

Bagi insinyur perangkat lunak yang membangun sistem terdistribusi berskala besar, evolusi ini membawa implikasi langsung terhadap strategi desain. Pertama, kekhawatiran terhadap jeda aplikasi akibat GC semakin berkurang. Pengembang dapat lebih leluasa mengalokasikan objek berumur pendek tanpa harus terlalu agresif menggunakan teknik object pooling, yang seringkali merumitkan logika kode dan rentan terhadap kutu memori.

Kedua, penyetelan performa (performance tuning) menjadi lebih sederhana. Pacer yang lebih cerdas berarti bahwa parameter default seringkali sudah memberikan performa yang optimal. Penggunaan GOGC atau variabel baru seperti GOMEMLIMIT sekarang lebih diarahkan untuk mendefinisikan batasan sumber daya ketat dalam lingkungan yang terisolasi, seperti kontainer awan, alih-alih untuk melakukan manipulasi paksa terhadap siklus GC.

## Kesimpulan

Evolusi garbage collection pada bahasa pemrograman Go mencerminkan komitmen berkelanjutan untuk menyediakan runtime yang andal dan berkinerja tinggi. Pergeseran dari pacer heuristik sederhana ke model kontrol adaptif pada Go 1.20+ telah terbukti efektif dalam meminimalisasi latensi stop-the-world. Seiring dengan kemajuan perangkat keras dan kompleksitas aplikasi modern, arsitektur GC semacam ini menjadi fondasi yang tidak terpisahkan dalam memastikan daya tanggap dan skalabilitas sistem. Penelitian lebih lanjut dapat difokuskan pada analisis komparatif antara GC Go dengan strategi manajemen memori pada bahasa lain dalam konteks komputasi tanpa server (serverless computing) yang memiliki karakteristik beban kerja yang unik.[^1][^2]

## Referensi

[^1]: Hudson, R. "Getting to Go: The Journey of Go's Garbage Collector." ISMM Keynote, 2018.
[^2]: Google Go Team. "A Guide to the Go Garbage Collector." Go Dev Documentation, 2023.