---
title: "Zsh vs Fish: Memilih Shell Harian yang Nyaman dan Cepat"
slug: c9d0e1f2
aliases: [zsh-fish-shell-tips]
date: 2026-07-22
tags: [linux, shell, terminal]
lang: id
draft: false
---

Terminal adalah rumah kedua bagi developer backend dan sysadmin. Dua shell modern yang paling banyak digunakan di Linux dan macOS adalah **Zsh** dan **Fish**.

Keduanya menawarkan fitur yang jauh melampaui Bash standar, tapi dengan filosofi desain yang berbeda.

## Fun Fact

**Zsh (Z Shell) dinamai dari nama profesor di Universitas Princeton.**
Paul Falstad membuat Zsh pada tahun 1990 saat masih kuliah di Princeton University. Ia menamai shell tersebut berdasarkan nama asisten profesor Zhong Shao.

**Fish (Friendly Interactive Shell) sengaja tidak 100% kompatibel dengan POSIX.**
Fish dirilis pada 2005 oleh Axel Liljencrantz dengan moto "konfigurasi nol (*works out of the box*)". Fish sengaja membuang sintaks Bash kuno agar sintaks skripnya lebih bersih dan intuitif.

**macOS beralih ke Zsh sebagai default shell sejak macOS Catalina (2019).**
Alasan utamanya adalah lisensi GNU GPLv3 pada Bash versi 4+, sedangkan Zsh memakai lisensi bergaya MIT yang lebih ramah bagi ekosistem komersial Apple.

---

## Tips dan Trik

### 1. Gunakan Prompt Cepat Seperti Starship di Kedua Shell

Alih-alih memakai framework berat seperti Oh-My-Zsh yang memperlambat startup terminal, pasang prompt berbasis Rust seperti **Starship**:

```bash
# Pasang starship
curl -sS https://starship.rs/install.sh | sh

# Di ~/.zshrc:
eval "$(starship init zsh)"

# Di ~/.config/fish/config.fish:
starship init fish | source
```

### 2. Auto-Suggestions Instan di Zsh

Fitur auto-suggestion abu-abu yang terkenal di Fish bisa dihadirkan di Zsh dengan plugin `zsh-autosuggestions`:

```bash
# Clone plugin
git clone https://github.com/zsh-users/zsh-autosuggestions ~/.zsh/zsh-autosuggestions

# Tambahkan ke ~/.zshrc
source ~/.zsh/zsh-autosuggestions/zsh-autosuggestions.zsh
```

### 3. Syntax Highlighting Real-Time

Beri warna hijau untuk perintah valid dan merah untuk typo sebelum kamu menekan Enter:

```bash
# Zsh
git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ~/.zsh/zsh-syntax-highlighting
source ~/.zsh/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh
```

### 4. Mengatur Alias Cepat dan Universal

Di Fish, alias otomatis disimpan sebagai fungsi modular:

```fish
# Di terminal fish langsung:
alias g='git'
alias l='ls -lah'
funcsave g
funcsave l
```

### 5. Benchmark Startup Time Shell Kamu

Pastikan waktu booting shell di bawah 50 milidetik agar membuka tab baru tidak terasa lag:

```bash
# Benchmark Zsh
for i in $(seq 1 10); do /usr/bin/time zsh -i -c exit; done

# Benchmark Fish
fish --profile-startup /tmp/fish.profile -c exit
```
