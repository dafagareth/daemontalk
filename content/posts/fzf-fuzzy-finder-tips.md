---
title: "Fzf: Fuzzy Finder yang Mengubah Cara Navigasi Terminal"
slug: e5f6a7b8
aliases: [fzf-fuzzy-finder-tips]
date: 2026-07-25
tags: [terminal, cli, tips]
lang: id
draft: false
---

`fzf` adalah interactive Unix filter untuk command line yang bekerja dengan input apapun: file, command history, git commit, proses, atau daftar hostname SSH.

Sekali terbiasa dengan `fzf`, kamu tidak akan pernah lagi mengetik path direktori panjang secara manual.

## Fun Fact

**Dibuat oleh Junegunn Choi dalam bahasa Go.**
Junegunn Choi (developer di balik vim-plug yang legendaris) menulis `fzf` pada 2013 dengan fokus kecepatan render dan responsivitas instan di puluhan ribu baris teks.

**Fzf bisa menerima input dari STDIN apapun.**
`fzf` bukan sekadar file searcher. Setiap output yang bisa di-pipe (`cat`, `git log`, `ps`, `docker ps`) bisa langsung disalurkan ke `fzf` untuk disaring secara interaktif.

**Mendukung algoritma fuzzy pencocokan karakter non-sekuensial.**
Mengetik `fzp` akan langsung mencocokkan `fuzzy_finder_project.go` karena fzf menghitung skor kedekatan dan batas kata (*word boundary*).

---

## Tips dan Trik

### 1. Aktifkan Keybinding Resmi: `Ctrl+R` dan `Alt+C`

Setelah menginstal fzf, integrasikan ke `.bashrc` atau `.zshrc`. `Ctrl+R` akan menggantikan pencarian history bawaan terminal dengan filter fuzzy yang sangat cepat, sementara `Alt+C` membuka fuzzy cd ke subdirektori.

```bash
# Tambahkan ke ~/.zshrc atau ~/.bashrc
source <(fzf --zsh) # atau source <(fzf --bash)
```

### 2. Gabungkan `fzf` dengan `bat` untuk Preview File Instan

Dengan opsi `--preview`, fzf akan merender preview syntax highlighted dari file yang sedang disorot di panel samping:

```bash
fzf --preview 'bat --style=numbers --color=always --line-range :500 {}'
```

### 3. Ganti Default Searcher dengan `ripgrep` / `fd`

Secara default fzf menggunakan `find`. Menggantinya dengan `fd` akan mengabaikan file `.git` dan `node_modules` secara otomatis sehingga pencarian jauh lebih cepat:

```bash
# ~/.bashrc atau ~/.zshrc
export FZF_DEFAULT_COMMAND='fd --type f --strip-cwd-prefix --hidden --exclude .git'
export FZF_CTRL_T_COMMAND="$FZF_DEFAULT_COMMAND"
```

### 4. Fuzzy Git Checkout Branch

Buat alias Git untuk berpindah branch tanpa harus mengingat ejaan nama branch:

```bash
git-checkout-fuzzy() {
  local branch=$(git branch --all | grep -v HEAD | string trim | fzf | sed "s/.* //" | sed "s#remotes/[^/]*/##")
  if [ -n "$branch" ]; then
    git checkout "$branch"
  fi
}
```

### 5. Interactive Docker Container Kill

Pilih container yang berjalan dan matikan langsung lewat antarmuka fuzzy:

```bash
docker ps | sed 1d | fzf -m | awk '{print $1}' | xargs -r docker stop
```
