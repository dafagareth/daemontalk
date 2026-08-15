---
title: "Berhenti Mengetik Perintah SSH yang Panjang"
slug: 396e35d3
aliases: [ssh-config-tips]
date: 2024-12-18
tags: [ssh, linux, tips]
lang: id
draft: false
---

Kalau kamu sering terhubung ke server lewat SSH, mungkin kamu hafal perintah seperti ini di luar kepala:

```bash
ssh -i ~/.ssh/keys/produksi.pem -p 2222 deploy@203.0.113.45
```

Mengetiknya berulang kali setiap hari adalah pemborosan, dan rawan salah ketik. Yang lebih buruk, ketika kamu mengelola lima atau enam server, mengingat IP dan port masing-masing menjadi beban tersendiri.

Ada satu file yang menyelesaikan semua ini, dan banyak orang tidak tahu keberadaannya: `~/.ssh/config`.

## File Config SSH

SSH membaca file `~/.ssh/config` setiap kali dijalankan. Di sana kamu bisa mendefinisikan koneksi dengan nama pendek beserta seluruh detailnya.

```
# ~/.ssh/config

Host produksi
    HostName 203.0.113.45
    User deploy
    Port 2222
    IdentityFile ~/.ssh/keys/produksi.pem
```

Dengan ini, perintah panjang tadi menjadi:

```bash
ssh produksi
```

SSH membaca alias `produksi`, mencari konfigurasinya, dan mengisi sendiri hostname, user, port, dan key yang sesuai. Tidak ada lagi yang perlu diingat selain nama yang kamu tentukan sendiri.

Manfaat ini menyebar ke alat lain. `scp`, `rsync`, dan `git` yang memakai SSH juga membaca config yang sama:

```bash
scp laporan.pdf produksi:/var/www/uploads/
rsync -avz ./dist/ produksi:/var/www/app/
```

## Beberapa Server Sekaligus

File config bisa memuat banyak host. Susun seluruh server kamu di satu tempat.

```
Host staging
    HostName 203.0.113.46
    User deploy
    Port 22
    IdentityFile ~/.ssh/keys/staging.pem

Host db-internal
    HostName 10.0.1.20
    User admin
    IdentityFile ~/.ssh/keys/internal.pem

Host vps-pribadi
    HostName vps.contoh.com
    User dafa
```

Setelah ini, `ssh staging`, `ssh db-internal`, dan `ssh vps-pribadi` langsung bekerja. Konfigurasi server menjadi dokumentasi yang hidup, tersimpan rapi alih-alih tersebar di catatan atau ingatan.

## Pengaturan Default untuk Semua Host

Blok `Host *` berlaku untuk semua koneksi. Ini tempat yang tepat untuk pengaturan umum.

```
Host *
    ServerAliveInterval 60
    ServerAliveCountMax 3
    AddKeysToAgent yes
```

`ServerAliveInterval 60` mengirim sinyal kecil setiap 60 detik agar koneksi tidak terputus saat idle, masalah yang sering muncul ketika kamu meninggalkan sesi SSH sebentar lalu kembali dan mendapati koneksi sudah mati. `AddKeysToAgent yes` menambahkan key ke ssh-agent secara otomatis sehingga kamu tidak perlu memasukkan passphrase berulang kali.

## Melompat Lewat Server Perantara

Sering kali server internal tidak bisa diakses langsung dari internet. Kamu harus masuk dulu ke server perantara (bastion atau jump host), baru dari sana terhubung ke server tujuan. Tanpa config, ini berarti dua langkah SSH manual.

`ProxyJump` menyatukannya menjadi satu perintah.

```
Host bastion
    HostName 203.0.113.10
    User jump

Host server-internal
    HostName 10.0.2.30
    User admin
    ProxyJump bastion
```

Dengan ini, `ssh server-internal` otomatis melompat lewat `bastion` terlebih dahulu. SSH mengurus seluruh rantai koneksi tanpa kamu perlu masuk dua kali.

## Mempercepat Koneksi Berulang

Jika kamu sering membuka beberapa sesi ke server yang sama, connection multiplexing membuat koneksi kedua dan seterusnya hampir instan dengan memakai ulang koneksi pertama.

```
Host *
    ControlMaster auto
    ControlPath ~/.ssh/sockets/%r@%h-%p
    ControlPersist 600
```

Jangan lupa membuat direktori socketnya terlebih dahulu:

```bash
mkdir -p ~/.ssh/sockets
```

Setelah koneksi pertama terbentuk, sesi berikutnya ke host yang sama tidak perlu melakukan handshake dan autentikasi ulang. `ControlPersist 600` menjaga koneksi master tetap hidup selama 600 detik setelah sesi terakhir ditutup.

## Soal Permission File

SSH cukup ketat soal keamanan. Jika file config atau key punya permission terlalu longgar, SSH akan menolak memakainya.

```bash
chmod 600 ~/.ssh/config
chmod 600 ~/.ssh/keys/*
chmod 700 ~/.ssh
```

Ini bukan rewel tanpa alasan. Private key yang bisa dibaca user lain di sistem yang sama adalah risiko keamanan nyata. SSH sengaja menolak bekerja sampai izinnya diperketat.

---

File `~/.ssh/config` adalah salah satu peningkatan kualitas hidup terbesar yang bisa kamu dapatkan dengan usaha paling kecil. Lima belas menit menyusun seluruh server ke dalam satu file menghemat ribuan ketikan ke depan, mengurangi kesalahan, dan mengubah koneksi yang rumit menjadi satu kata pendek. Kalau kamu sering bekerja dengan server jarak jauh, ini bukan optimasi, melainkan kebersihan kerja yang mendasar.
