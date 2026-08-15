---
title: "Menjalankan Aplikasi sebagai systemd Service"
slug: 539915c7
aliases: [systemd-service-pemula]
date: 2025-06-15
tags: [linux, systemd, sysadmin]
lang: id
draft: false
---

Kamu sudah membuat aplikasi, mengunggahnya ke server, dan menjalankannya. Semua berfungsi. Lalu kamu menutup sesi SSH, dan aplikasinya ikut mati. Atau server reboot karena pembaruan, dan aplikasinya tidak menyala kembali. Kamu harus masuk lagi, menjalankannya lagi, secara manual.

Cara amatir mengatasi ini adalah dengan `nohup` atau menjalankan proses di latar belakang dengan `&`. Cara yang benar di server Linux modern adalah membuatnya sebagai **systemd service**. Setelah itu, aplikasi menyala saat boot, hidup kembali jika crash, dan dikelola dengan perintah yang sama seperti service sistem lainnya.

## Apa Itu systemd

systemd adalah sistem init yang dipakai hampir semua distribusi Linux modern: Ubuntu, Debian, Fedora, Arch, dan lainnya. Ia bertugas menjalankan dan mengelola proses yang berjalan di latar belakang, yang disebut service atau unit. Database, web server, dan SSH yang berjalan di server kamu semuanya dikelola systemd.

Aplikasimu bisa bergabung dalam pengelolaan yang sama. Yang kamu perlukan hanya satu file teks.

## Unit File Pertama

Unit file diletakkan di `/etc/systemd/system/`. Buat file bernama `myapp.service`:

```ini
# /etc/systemd/system/myapp.service

[Unit]
Description=Aplikasi Web Saya
After=network.target

[Service]
Type=simple
User=deploy
WorkingDirectory=/var/www/myapp
ExecStart=/var/www/myapp/app
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Mari pahami setiap bagiannya.

Bagian `[Unit]` berisi metadata. `Description` adalah penjelasan singkat. `After=network.target` memberi tahu systemd untuk menjalankan service ini setelah jaringan siap, penting untuk aplikasi yang butuh koneksi.

Bagian `[Service]` adalah inti pengaturan. `Type=simple` cocok untuk program yang berjalan di latar depan dan tidak melakukan fork sendiri, yang merupakan mayoritas aplikasi modern. `User=deploy` menjalankan aplikasi sebagai user biasa, bukan root, sebuah praktik keamanan yang penting. `WorkingDirectory` dan `ExecStart` menentukan dari mana dan apa yang dijalankan.

Dua baris yang membuat systemd istimewa: `Restart=on-failure` membuat systemd otomatis menjalankan ulang aplikasi jika ia crash, dan `RestartSec=5` memberi jeda 5 detik sebelum mencoba lagi agar tidak terjadi restart beruntun yang membebani sistem.

Bagian `[Install]` menentukan kapan service diaktifkan. `WantedBy=multi-user.target` berarti service ini dijalankan saat sistem mencapai mode operasi normal, yaitu saat boot biasa.

## Mengaktifkan dan Menjalankan

Setelah membuat file, beri tahu systemd untuk membaca ulang konfigurasinya, lalu aktifkan dan jalankan service.

```bash
# Muat ulang konfigurasi setelah membuat atau mengubah unit file
sudo systemctl daemon-reload

# Aktifkan agar otomatis jalan saat boot
sudo systemctl enable myapp

# Jalankan sekarang
sudo systemctl start myapp
```

Perhatikan perbedaan `enable` dan `start`. `enable` mendaftarkan service agar menyala otomatis saat boot, tapi tidak langsung menjalankannya. `start` menjalankannya sekarang, tapi tidak mengatur perilaku saat boot. Biasanya kamu memakai keduanya. Ada juga jalan pintas: `systemctl enable --now myapp` melakukan keduanya sekaligus.

## Memeriksa Status dan Log

Setelah berjalan, kamu bisa memeriksa keadaannya kapan saja.

```bash
sudo systemctl status myapp
```

Perintah ini menampilkan apakah service aktif, sejak kapan berjalan, process ID-nya, dan beberapa baris log terakhir. Ini hal pertama yang dilihat saat sesuatu bermasalah.

Untuk log lengkap, systemd menyediakan `journalctl`:

```bash
# Lihat seluruh log service ini
sudo journalctl -u myapp

# Ikuti log secara real-time, seperti tail -f
sudo journalctl -u myapp -f

# Lihat log sejak boot terakhir saja
sudo journalctl -u myapp -b
```

Salah satu keunggulan menjalankan aplikasi sebagai service adalah seluruh outputnya otomatis tercatat. Apa pun yang aplikasimu cetak ke standard output dan standard error masuk ke journal, lengkap dengan cap waktu, tanpa kamu perlu mengatur file log sendiri.

## Mengubah Konfigurasi

Saat kamu mengubah unit file, systemd tidak langsung tahu. Ia perlu diberi tahu untuk membaca ulang.

```bash
sudo systemctl daemon-reload
sudo systemctl restart myapp
```

Lupa menjalankan `daemon-reload` setelah mengedit unit file adalah kesalahan yang sangat umum. Perubahan seolah tidak berpengaruh karena systemd masih memakai versi lama di memorinya.

## Variabel Environment

Aplikasi sering butuh konfigurasi lewat environment variable. Ada dua cara memberikannya. Untuk nilai sederhana, langsung di unit file:

```ini
[Service]
Environment=PORT=8080
Environment=NODE_ENV=production
```

Untuk konfigurasi yang lebih banyak atau berisi rahasia, lebih rapi memakai file terpisah:

```ini
[Service]
EnvironmentFile=/etc/myapp/env
```

Lalu isi `/etc/myapp/env` dengan pasangan kunci-nilai, dan atur permissionnya agar hanya bisa dibaca yang berhak:

```bash
sudo chmod 600 /etc/myapp/env
```

Memisahkan konfigurasi rahasia ke file dengan permission ketat lebih aman daripada menuliskannya langsung di unit file yang permissionnya lebih longgar.

---

Membuat systemd service membutuhkan satu file teks dan beberapa perintah, tapi mengubah cara aplikasimu berjalan secara fundamental. Tidak ada lagi proses yang mati saat SSH terputus, tidak ada lagi aplikasi yang tidak menyala setelah reboot, tidak ada lagi aplikasi yang tumbang permanen setelah crash. Ditambah pencatatan log otomatis dan perintah pengelolaan yang seragam, ini adalah cara standar dan benar untuk menjalankan apa pun di server Linux. Begitu kamu terbiasa, `nohup &` mulai terasa seperti menambal sesuatu yang sebenarnya sudah ada solusi rapinya.
