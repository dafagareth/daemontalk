---
title: "Kelola Dotfiles dengan Git Tanpa Tooling Tambahan"
slug: 7f9d3771
aliases: [dotfiles-dengan-git]
date: 2026-01-20
tags: [linux, git, tools, dotfiles]
lang: id
draft: false
---

Dotfiles adalah file konfigurasi yang tersebar di seluruh home directory: `.bashrc`, `.zshrc`, `.gitconfig`, `~/.config/nvim/`, dan seterusnya. Ketika kamu ganti mesin atau install ulang sistem, setup konfigurasi dari nol membutuhkan waktu berjam-jam.

Ada banyak tool yang dibuat untuk mengelola dotfiles, tapi semuanya menambah lapisan abstraksi yang tidak perlu. Git sudah cukup, dengan satu trik kecil.

## Metode Bare Repository

Idenya sederhana: buat repository Git biasa, tapi simpan direktori `.git`-nya di tempat lain (misalnya `~/.dotfiles`), bukan di dalam direktori yang sama dengan file yang dikelola.

Setup awal:

```bash
git init --bare $HOME/.dotfiles
```

Buat alias untuk menjalankan perintah Git terhadap repo ini tanpa harus masuk ke direktorinya:

```bash
alias dots='git --git-dir=$HOME/.dotfiles --work-tree=$HOME'
```

Tambahkan alias ini ke `.zshrc` atau `.bashrc` supaya tersedia di sesi berikutnya.

Satu konfigurasi penting agar Git tidak menampilkan ribuan file yang tidak ditrack:

```bash
dots config status.showUntrackedFiles no
```

## Menambahkan File

Sekarang kamu bisa mengelola file konfigurasi dari mana saja:

```bash
dots add ~/.zshrc
dots add ~/.gitconfig
dots add ~/.config/nvim/init.lua
dots commit -m "tambah konfigurasi awal"
```

Push ke remote:

```bash
dots remote add origin git@github.com:dafagareth/dotfiles.git
dots push -u origin main
```

Itu saja. Tidak perlu symlink, tidak perlu tool tambahan.

## Penggunaan Sehari-hari

Setelah setup awal, alur kerjanya sama persis dengan Git biasa, hanya menggunakan alias `dots`:

```bash
dots status                  # lihat perubahan
dots diff ~/.zshrc           # diff sebelum commit
dots add ~/.config/nvim/     # tambah direktori
dots commit -m "update nvim config"
dots push
```

## Setup di Mesin Baru

Clone ke bare repository:

```bash
git clone --bare git@github.com:dafagareth/dotfiles.git $HOME/.dotfiles
```

Definisikan alias sementara:

```bash
alias dots='git --git-dir=$HOME/.dotfiles --work-tree=$HOME'
```

Checkout file-file konfigurasi ke home directory:

```bash
dots checkout
```

Jika ada file yang sudah ada (misalnya `.bashrc` default dari sistem), Git akan menolak. Backup dulu:

```bash
mkdir -p ~/.dotfiles-backup
dots checkout 2>&1 | grep "^\s" | awk '{print $1}' | xargs -I{} mv {} ~/.dotfiles-backup/{}
dots checkout
```

Terakhir, tambahkan konfigurasi:

```bash
dots config status.showUntrackedFiles no
```

## Apa yang Perlu Ditrack

Track konfigurasi yang benar-benar personal dan butuh waktu untuk dikonfigurasi ulang:

```
~/.zshrc
~/.bashrc
~/.gitconfig
~/.gitignore_global
~/.ssh/config          # bukan private key, hanya config
~/.config/nvim/
~/.config/tmux/
~/.config/alacritty/
~/.config/starship.toml
```

Jangan track hal-hal yang berisi kredensial, data spesifik mesin, atau yang di-generate otomatis:

```
~/.ssh/id_*            # private keys
~/.aws/credentials
~/.config/gh/hosts.yml
~/.local/share/        # data aplikasi
```

Buat `~/.gitignore` untuk global ignores, lalu track file itu sendiri.

## Keunggulan Metode Ini

Tidak ada tool tambahan yang bisa usang atau tidak kompatibel. Siapapun yang tahu Git sudah bisa memahami cara kerjanya. Repository bisa di-clone langsung dari GitHub dan langsung dipakai. History perubahan konfigurasi tersimpan dengan baik, lengkap dengan pesan commit.

Satu-satunya kekurangan: tidak ada UI yang bagus, hanya Git murni. Tapi untuk sesuatu yang sifatnya infrastruktur personal, itu justru keunggulan.
