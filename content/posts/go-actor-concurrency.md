---
title: "Pemodelan Konkurensi Aktor menggunakan Goroutine pada Sistem Terdistribusi"
slug: 5b4a3c2d
aliases: []
date: 2026-08-11
tags: [go]
lang: id
draft: false
type: post
cover: ""
---

## Pengantar

Sistem terdistribusi modern dihadapkan pada tantangan untuk menangani ribuan bahkan jutaan proses konkuren secara efisien. Secara tradisional, manajemen status bersama dalam lingkungan multithreading diatasi dengan primitif sinkronisasi seperti mutex dan semafor. Namun, pendekatan ini sering menimbulkan masalah berupa kebuntuan (deadlock) dan pelambatan akibat pertentangan sumber daya (resource contention). Model aktor (actor model) muncul sebagai paradigma alternatif yang kuat, di mana unit komputasi independen berkomunikasi semata-mata melalui pertukaran pesan. Artikel ini mengeksplorasi penerapan model aktor menggunakan fitur primitif konkurensi di bahasa Go, yaitu goroutine dan channel.

## Prinsip Dasar Model Aktor

Dalam teori komputasi konkuren, model aktor mendefinisikan aktor sebagai entitas dasar komputasi primitif. Ketika menerima sebuah pesan, sebuah aktor dapat melakukan tiga aksi utama, yaitu mengubah status internalnya, mengirim pesan kepada aktor lain, atau membuat aktor baru. Karakteristik utama dari model ini adalah isolasi status, yang berarti tidak ada status yang dibagikan secara langsung antar aktor. Semua interaksi harus dimediasi melalui mekanisme penyampaian pesan asinkron.

Isolasi ini secara inheren menghilangkan masalah pertentangan memori. Oleh karena itu, model aktor sangat cocok diimplementasikan pada arsitektur perangkat keras modern yang didominasi oleh prosesor multicore dan sistem terdistribusi di mana memori tidak lagi berbagi secara fisik.

## Pemetaan Aktor ke dalam Idiom Go

Bahasa Go dirancang dengan filosofi konkurensi yang berpusat pada proses sekuensial yang saling berkomunikasi (Communicating Sequential Processes atau CSP). Meskipun CSP memiliki dasar teoritis yang sedikit berbeda dengan model aktor, goroutine dan channel menyediakan fondasi yang sangat baik untuk membangun kerangka kerja berbasis aktor.

### Goroutine sebagai Unit Aktor

Goroutine adalah benang eksekusi ringan yang dikelola oleh runtime Go. Biaya pembuatan dan penjadwalan goroutine sangat rendah dibandingkan dengan thread sistem operasi tradisional. Hal ini memungkinkan pembuatan aktor dalam jumlah masif. Dalam implementasi praktis, sebuah aktor direpresentasikan sebagai goroutine yang berjalan terus-menerus dalam sebuah kalang (loop) tak berujung, menunggu pesan masuk untuk diproses.

### Channel sebagai Kotak Surat (Mailbox)

Komponen esensial kedua dalam pemodelan aktor di Go adalah channel. Channel berfungsi ganda sebagai antrean pesan dan mekanisme sinkronisasi. Sebuah channel Go dapat bertindak sebagai kotak surat untuk goroutine aktor. Pengirim dapat meletakkan pesan ke dalam channel, dan goroutine aktor akan mengambil serta memproses pesan tersebut satu per satu secara berurutan. Karakteristik penguncian internal pada channel memastikan bahwa operasi pengiriman dan penerimaan pesan bersifat thread-safe.

## Arsitektur Sistem Terdistribusi Berbasis Aktor

Menerapkan model aktor secara lokal dalam sebuah proses Go adalah langkah pertama. Tantangan sesungguhnya terletak pada perluasan model ini ke lingkungan terdistribusi, di mana aktor dapat berada pada mesin fisik yang berbeda dalam sebuah klaster.

### Resolusi Alamat dan Perutean Pesan

Dalam sistem terdistribusi, pengirim pesan tidak perlu mengetahui lokasi fisik dari aktor penerima. Sistem memerlukan mekanisme penamaan atau resolusi alamat (address resolution). Pendekatan yang umum digunakan adalah dengan mengimplementasikan lapisan perantara atau layanan registrasi di mana aktor dapat mendaftarkan ID unik mereka. Ketika sebuah pesan dikirim, lapisan perutean akan mengidentifikasi lokasi aktor, baik lokal maupun jarak jauh, dan meneruskan pesan tersebut menggunakan protokol jaringan yang sesuai seperti gRPC atau TCP murni.

### Toleransi Kesalahan (Fault Tolerance)

Kelebihan signifikan dari model aktor adalah kemampuannya dalam menangani kegagalan sistem. Terinspirasi dari mesin virtual Erlang, sistem aktor sering mengadopsi hierarki pengawasan (supervision trees). Dalam desain ini, sebuah aktor induk bertanggung jawab untuk memantau siklus hidup aktor anaknya. Jika sebuah aktor anak mengalami kegagalan, aktor induk dapat mengambil keputusan untuk memulai ulang (restart) anak tersebut, meneruskan kesalahan ke tingkat yang lebih tinggi, atau menghentikan eksekusi. Di Go, mekanisme ini dapat disimulasikan menggunakan goroutine pengawas yang memonitor channel sinyal dari goroutine pekerja.

## Evaluasi Kinerja dan Keterbatasan

Implementasi aktor murni menggunakan goroutine menawarkan kinerja komputasi yang mengesankan berkat model pemetaan benang hibrida dalam runtime Go. Proses perpindahan konteks antar goroutine terjadi pada ruang pengguna (user space), yang jauh lebih cepat dibandingkan perpindahan konteks sistem operasi.

Namun, terdapat beberapa tantangan yang perlu diantisipasi. Penggunaan channel secara berlebihan tanpa desain aliran data yang cermat dapat berujung pada kebuntuan, terutama jika ukuran antrean pesan (buffer) tidak dikonfigurasi dengan tepat. Selain itu, garbage collector pada Go harus memindai setiap goroutine aktif beserta tumpukannya. Sistem dengan puluhan juta aktor yang sebagian besar tidak aktif (idle) dapat memberikan beban tambahan pada fase pemindaian GC.

## Kesimpulan

Model aktor menawarkan kerangka konseptual yang kokoh untuk mengelola kompleksitas aplikasi konkuren skala besar. Bahasa Go, dengan goroutine dan channel, menyediakan elemen-elemen fundamental yang memfasilitasi konstruksi sistem berbasis aktor secara alami dan efisien. Penyatuan paradigma aktor dengan kekuatan jaringan bawaan Go membuka jalan bagi pengembangan layanan terdistribusi yang sangat tangguh, skalabel, dan tahan uji. Penelitian mendatang perlu difokuskan pada perbandingan antara pustaka aktor Go yang sudah ada dengan pendekatan kustom dari nol, serta optimalisasi latensi perutean pesan antar jaringan (inter-network message routing latency).[^1][^2]

## Referensi

[^1]: Hoare, C. A. R. "Communicating Sequential Processes." Communications of the ACM, 1978.
[^2]: Hewitt, C., et al. "A Universal Modular Actor Formalism for Artificial Intelligence." IJCAI, 1973.