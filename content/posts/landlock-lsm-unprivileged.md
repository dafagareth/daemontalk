---
title: "Isolasi Proses Unprivileged menggunakan Landlock LSM secara Mandiri"
slug: 3f4e5d6c
aliases: []
date: 2026-08-11
tags: [security]
lang: id
draft: false
type: post
cover: ""
---

Pendahuluan
Keamanan sistem operasi modern sangat bergantung pada mekanisme isolasi proses yang presisi. Secara tradisional, pembatasan hak akses komprehensif pada sistem berbasis Linux memerlukan hak istimewa administratif (root). Hal ini menimbulkan masalah dalam desain perangkat lunak yang berprinsip pada hak akses minimal (least privilege). Untuk mengatasi hambatan tersebut, kernel Linux versi 5.13 memperkenalkan Landlock. Landlock adalah modul keamanan Linux (Linux Security Module atau LSM) yang memungkinkan aplikasi tanpa hak akses administrator (unprivileged) untuk membatasi dirinya sendiri dan proses turunannya dalam sebuah lingkungan sandbox.

Arsitektur Dasar dan Komponen Internal
Landlock berbeda dari mekanisme sandboxing lain seperti Seccomp-BPF. Sementara Seccomp-BPF fokus pada pemfilteran panggilan sistem (system calls) secara prosedural, Landlock menggunakan pendekatan semantik berbasis objek terhadap sistem file (filesystem). Komponen utama dari Landlock adalah ruleset. Ruleset merupakan struktur data yang menampung koleksi aturan akses. Aturan-aturan ini mendefinisikan interaksi yang diizinkan terhadap file dan direktori tertentu, seperti izin membaca, menulis, atau mengeksekusi file.

Proses inisialisasi Landlock melibatkan tiga tahapan sistem panggilan yang spesifik. Pertama, proses harus memanggil landlock_create_ruleset untuk membuat wadah aturan yang kosong. Kedua, proses menggunakan landlock_add_rule untuk mendefinisikan batasan akses untuk struktur direktori tertentu, yang diidentifikasi melalui file descriptor. Terakhir, proses memanggil landlock_restrict_self untuk menerapkan kebijakan keamanan yang telah dikonfigurasi.

Keistimewaan dan Pewarisan Hak Akses
Salah satu karakteristik krusial dari Landlock adalah sifatnya yang tidak dapat dibatalkan (irreversible). Setelah landlock_restrict_self dieksekusi dengan sukses, proses yang memanggilnya akan terikat secara permanen pada batasan yang didefinisikan dalam ruleset. Lebih jauh lagi, kebijakan ini diturunkan kepada setiap thread anak yang diciptakan melalui pemanggilan fork atau clone.

Sebelum Landlock dapat diaktifkan oleh proses tanpa hak istimewa, ada prasyarat keamanan yang harus dipenuhi. Proses tersebut wajib memiliki properti no_new_privs yang diset melalui prctl. Properti ini menjamin bahwa proses dan anak-anaknya tidak akan pernah memperoleh hak istimewa baru, bahkan jika mereka mengeksekusi biner yang memiliki bit set-user-ID (SUID). Integrasi erat antara Landlock dan no_new_privs mencegah eskalasi hak istimewa secara tidak disengaja maupun eksploitatif.

Komparasi Kinerja dan Efektivitas
Dibandingkan dengan mekanisme isolasi konvensional seperti chroot atau namespaces, Landlock menawarkan tingkat granularitas yang jauh lebih tinggi. Chroot secara prinsip rentan terhadap pelepasan paksa (breakouts) jika tidak dikonfigurasi bersama kapabilitas kernel tambahan. Network dan mount namespaces membutuhkan kompleksitas administrasi jaringan dan disk yang tidak sepele. Landlock menghindari kerumitan ini dengan membatasi ruang pandang sistem file langsung dari dalam proses, tanpa memodifikasi lingkungan globa sistem.

Dalam skenario aplikasi peladen web, Landlock memungkinkan peladen untuk membatasi dirinya sendiri hanya untuk membaca dari direktori berkas statis, meskipun pada dasarnya aplikasi tersebut dijalankan oleh pengguna yang mungkin memiliki akses baca dan tulis penuh terhadap direktori konfigurasi dan direktori penguna lainnya. Dengan demikian, celah kerentanan seperti Path Traversal dapat diredam secara efektif. Jika seorang penyerang berhasil mengeksploitasi cacat dalam penguraian URL, penyerang tersebut tetap tidak dapat membaca file sensitif seperti /etc/shadow karena direktori tersebut tidak secara eksplisit diizinkan dalam ruleset Landlock.

Tantangan Adopsi dalam Pengembangan Perangkat Lunak
Meskipun Landlock menjanjikan peningkatan keamanan yang signifikan, implementasinya dalam basis kode aplikasi yang ada membutuhkan analisis mendalam terhadap operasi sistem file. Aplikasi yang secara konstan memerlukan pembuatan file sementara di berbagai lokasi disk mungkin menghadapi kendala operasional jika direktori target tidak terantisipasi selama pembentukan ruleset. Oleh karena itu, pengembang harus merancang profil keamanan yang ketat namun memadai untuk siklus hidup penuh aplikasi.

Evolusi modul keamanan Landlock terus berlanjut seiring dengan rilis versi kernel yang lebih baru. Iterasi pengembangan baru memperkenalkan fitur pemfilteran yang lebih canggih, melingkupi operasi penggantian nama file (rename), pembuatan tautan (link creation), dan pembatasan pada modifikasi metadata objek sistem file. 

Kesimpulan
Landlock LSM merepresentasikan lompatan signifikan dalam paradigma keamanan sistem Linux. Dengan mendelegasikan kemampuan konstruksi sandbox langsung ke level aplikasi pengguna, pengembang kini dapat menerapkan pertahanan berlapis (defense-in-depth) secara terprogram. Pemanfaatan Landlock membatasi luasan area serangan (attack surface) secara drastis dalam peristiwa kompromi aplikasi. Seiring dengan kematangan ekosistem dokumentasi dan perbaikan API secara terus menerus, adopsi Landlock diproyeksikan menjadi standar industri bagi perangkat lunak yang berfokus pada ketahanan dan keamanan tinggi.
[^1][^2]

## Referensi

[^1]: Salaün, M. "Landlock: Unprivileged Access Control." Linux Kernel Documentation, 2021.
[^2]: Edge, J. "Sandboxing with Landlock." LWN.net, 2020.