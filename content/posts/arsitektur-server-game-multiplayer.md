---
title: "Arsitektur Server pada Game Multiplayer Modern"
slug: arsitektur-server-game-multiplayer
aliases: []
date: 2026-08-28
author: "Daemontalk Editorial"
tags: ["Gaming", "GameDev"]
lang: id
draft: false
type: post
cover: "https://images.unsplash.com/photo-1542751371-adc38448a05e?auto=format&fit=crop&q=80&w=1600"
coverCaption: "Pengontrol konsol game modern"
coverSource: "https://unsplash.com"
readTime: 6
description: "Mengenal model Server-Authoritative, prediksi klien (Client Prediction), dan trik manipulasi latensi agar game kompetitif terasa responsif."
---

Membangun aplikasi web itu menantang, tetapi membangun *server* untuk *game multiplayer* kompetitif (seperti Valorant, CS2, atau Dota 2) jauh lebih sulit. 

*Game server* harus menyinkronkan status fisik dunia virtual (*state*) untuk puluhan hingga ratusan pemain dalam waktu puluhan milidetik secara *real-time*, sambil mencegah kecurangan (*cheating*).

## Jangan Pernah Percaya Klien (Server-Authoritative)

Aturan pertama dalam *multiplayer game development* adalah sama dengan prinsip keamanan siber: **Jangan pernah percaya data dari klien**.

Jika komputer pemain memiliki otoritas untuk menentukan apakah tembakannya kena atau tidak, *hacker* bisa dengan mudah mengirim paket data palsu ke server ("Saya baru saja menembak tepat di kepala lawan").

```text
Model Peer-to-Peer (Rentan Cheat):
[Player A] <--- (Menentukan Tembakan) ---> [Player B]

Model Server-Authoritative (Aman):
[Player A] ---> (Kirim Input Tombol) ---> [Game Server]
                                              |
                                              v
                                      (Server Menghitung Fisika)
                                              |
                                              v
[Player B] <--- (Menerima Hasil Akhir) <------+
```

Dalam model *Server-Authoritative*, klien murni berfungsi sebagai "layar dan *joystick* bodoh". Server-lah yang melakukan kalkulasi fisika sejati.

## Dilema Latensi (Lag)

Masalah besar dari *Server-Authoritative* adalah latensi jaringan. Jika Anda menekan tombol maju, dan baru 100ms kemudian karakter di layar Anda bergerak (menunggu izin server), *game* akan terasa sangat berat dan tidak responsif.

Untuk mengatasi ini, *game engine* modern (seperti Unreal atau Unity) menggunakan teknik manipulasi waktu.

### 1. Client Prediction (Prediksi Klien)
Begitu pemain menekan tombol maju, komputer klien langsung memajukan karakter di layarnya **tanpa menunggu server**. Secara bersamaan, input dikirim ke server. Jika server setuju, tidak ada masalah. Jika server menolak (misal ternyata karakter menabrak dinding yang tak terlihat di klien), server akan memaksa klien untuk mengoreksi posisi (*rubber-banding*).

### 2. Lag Compensation (Kompensasi Lag)
Saat Anda menembak musuh yang bergerak, tembakan Anda memiliki jeda perjalanan ke server. Untuk memastikan tembakan terasa adil, server menyimpan riwayat pergerakan (rekaman masa lalu) setiap pemain selama sekian milidetik terakhir.

Ketika server menerima sinyal "menembak" dari Anda, server akan **memutar balik waktu** (secara matematis) untuk mencocokkan di mana posisi musuh tepat di saat Anda menarik pelatuk di layar Anda.

> [!NOTE]
> Inilah alasan mengapa kadang Anda merasa sudah bersembunyi di balik tembok, tapi masih tertembak mati. Di layar musuh (di masa lalu yang dikompensasi server), Anda belum masuk ke balik tembok.

## Protokol Komunikasi: UDP vs TCP

Sebagian besar data *gameplay* *real-time* (posisi kordinat X, Y, Z) dikirim melalui **UDP**, bukan TCP. UDP tidak meminta konfirmasi penerimaan, sehingga jauh lebih cepat. Jika ada paket koordinat pergerakan yang hilang di jalan (*packet loss*), *game* tidak akan peduli dan hanya menunggu koordinat berikutnya di detik selanjutnya.

> [!IMPORTANT]
> Menggunakan TCP untuk koordinat posisi akan memicu fenomena *head-of-line blocking* di mana jika satu paket hilang, seluruh *game* akan *freeze* menunggu paket itu dikirim ulang.

## Referensi Pembelajaran

```references
- title: "Source Multiplayer Networking"
  author: "Valve Developer Community"
  year: 2024
  publisher: "Valve Corporation"
  url: "https://developer.valvesoftware.com/wiki/Source_Multiplayer_Networking"
  
- title: "Fast-Paced Multiplayer"
  author: "Gabriel Gambetta"
  year: 2016
  publisher: "gabrielgambetta.com"
```
