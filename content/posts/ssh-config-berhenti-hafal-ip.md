---
title: "Berhenti Menghafal IP Server dengan SSH Config"
slug: "ssh-config-berhenti-hafal-ip"
aliases: []
author: "daemontalk team"
contributors: []
tags: ["wire", "linux", "devops"]
lang: "id"
status: "published"
type: "dispatch"
readTime: 4
cover: "https://images.unsplash.com/photo-1629654297299-c8506221ca97?auto=format&fit=crop&w=1200&q=80"
coverCaption: "Akses terminal rapi tanpa deretan angka IP"
coverSource: "https://unsplash.com"
description: "Menghafal alamat IP, port kustom, dan flag SSH adalah kebiasaan buruk. Cukup satu berkas ~/.ssh/config untuk merapikan seluruh alur kerja remote server Anda."
series: ""
series_part: 0
---

<span class="drop-cap">B</span>erapa banyak alamat IP server yang tercecer di notepad, riwayat shell terminal, atau pesan chat tim Anda?

Bagi sebagian developer, rutinitas masuk ke server produksi terlihat seperti ini:

```bash
ssh -i ~/.ssh/keys/client_prod_ed25519 -p 2202 deployer@103.145.22.84
```

Perintah di atas panjang, lambat diketik, dan sangat rawan salah sasaran saat Anda terburu-buru menangani insiden tengah malam.

Solusinya sudah ada di laptop Anda sejak puluhan tahun lalu: berkas `~/.ssh/config`.

## Sederhanakan Perintah dalam Satu Alias

Alih-alih mengetik semua flag di terminal, Anda bisa memetakan kredensial server ke dalam alias pendek di `~/.ssh/config`:

```ssh-config
Host prod-db
    HostName 103.145.22.84
    User deployer
    Port 2202
    IdentityFile ~/.ssh/keys/client_prod_ed25519
```

Setelah berkas ini disimpan, Anda cukup mengetik:

```bash
ssh prod-db
```

Klien SSH akan membaca nama host, user, port, dan kunci privat secara otomatis. Tab-completion di bash maupun zsh bahkan langsung mengenali alias tersebut.

## Melompati Bastion Host Tanpa Ribet

Skenario umum di infrastruktur cloud: server basis data berada di subnet privat dan hanya bisa diakses melalui bastion host (jump host).

Cara konvensional memaksa Anda membuat SSH tunnel manual yang rumit. Dengan berkas config, gunakan opsi `ProxyJump`:

```ssh-config
Host bastion
    HostName 103.145.22.10
    User jumpuser
    IdentityFile ~/.ssh/bastion_key

Host internal-db
    HostName 10.0.2.45
    User postgres
    IdentityFile ~/.ssh/internal_key
    ProxyJump bastion
```

Ketik `ssh internal-db`, dan OpenSSH akan otomatis menyambungkan koneksi melalui bastion secara transparan di latar belakang.

## Menjaga Koneksi Tetap Hidup

Pernahkah sesi terminal Anda membeku (hang) saat ditinggal mengambil kopi selama sepuluh menit?

Masalah ini terjadi karena router atau firewall perantara memutus koneksi idle TCP. Anda bisa mencegahnya untuk semua koneksi dengan menambahkan blok default di bagian paling bawah berkas:

```ssh-config
Host *
    ServerAliveInterval 60
    ServerAliveCountMax 3
    AddKeysToAgent yes
```

Pengaturan ini memerintahkan klien mengirim paket *keep-alive* setiap 60 detik. Flag `AddKeysToAgent yes` juga otomatis mendaftarkan passphrase kunci ke ssh-agent saat pertama kali dipakai.

## Pisahkan Berkas Berdasarkan Proyek

Jika Anda mengelola puluhan server untuk berbagai klien, berkas config bisa membengkak hingga ratusan baris.

Gunakan direktif `Include` di baris pertama `~/.ssh/config` untuk memecahnya per proyek:

```ssh-config
Include ~/.ssh/conf.d/*
```

Buat direktori `~/.ssh/conf.d/`, lalu simpan konfigurasi per lingkungan, misalnya `kantor.conf`, `pribadi.conf`, atau `klien-a.conf`.

Berhenti menyiksa memori kepala Anda dengan mengingat deretan angka IP dan port acak. Luangkan lima menit untuk merapikan `~/.ssh/config`, dan kembalikan fokus Anda ke pekerjaan rekayasa sistem yang sesungguhnya.

```references
- title: "ssh_config(5) OpenSSH Client Configuration File"
  author: "OpenBSD Project"
  year: 2024
  publisher: "OpenBSD Manual Pages"
  url: "https://man.openbsd.org/ssh_config.5"
```
