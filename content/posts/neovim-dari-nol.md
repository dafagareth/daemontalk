---
title: "Neovim dari Nol: Membangun Editor yang Benar-Benar Milikmu"
slug: f4ef8090
aliases: [neovim-dari-nol]
date: 2025-08-30
tags: [neovim, editor, productivity]
lang: id
draft: false
---

Setiap kali seseorang menyebut Vim atau Neovim, percakapan biasanya berhenti di dua hal: lelucon soal "cara keluar dari Vim", dan ketakutan bahwa mempelajarinya butuh waktu berbulan-bulan. Keduanya membuat banyak orang tidak pernah mencoba.

Kenyataannya, Neovim tidak harus dipelajari sekaligus. Kamu bisa mulai dengan editor yang nyaris kosong, lalu menambahkan satu hal pada satu waktu sesuai kebutuhan. Yang kamu dapatkan di akhir bukan sekadar editor, melainkan editor yang setiap bagiannya kamu pahami karena kamu sendiri yang menyusunnya.

## Kenapa Repot

Pertanyaan yang wajar: kenapa tidak pakai VS Code saja yang langsung jalan?

Jawabannya bukan kecepatan editing atau gengsi. Jawabannya adalah pemahaman dan kendali. Saat kamu menyusun konfigurasi sendiri, kamu tahu persis apa yang berjalan dan kenapa. Tidak ada plugin misterius yang memakan memori, tidak ada fitur yang tidak pernah kamu pakai. Editor itu seringan dan secepat yang kamu buat.

Selain itu, Vim motion, cara berpindah dan mengedit teks tanpa mouse, adalah keterampilan yang terbawa ke mana-mana. Banyak editor lain, termasuk VS Code, punya mode emulasi Vim. Mempelajari motionnya satu kali bermanfaat seumur hidup.

## Mulai dari Gerakan, Bukan Konfigurasi

Sebelum menyentuh konfigurasi apa pun, luangkan waktu memahami konsep inti Vim: **mode**. Tidak seperti editor biasa, Vim punya mode terpisah untuk mengetik dan untuk bernavigasi.

- **Normal mode**: untuk bergerak dan menjalankan perintah (mode default)
- **Insert mode**: untuk mengetik teks, masuk dengan `i`
- **Visual mode**: untuk memilih teks, masuk dengan `v`
- **Command mode**: untuk perintah seperti simpan dan keluar, masuk dengan `:`

Kembali ke Normal mode selalu dengan `Esc`. Untuk keluar dari editor, dari Normal mode ketik `:q` lalu Enter. Untuk menyimpan, `:w`. Untuk keduanya, `:wq`. Lelucon soal "tidak bisa keluar dari Vim" sepenuhnya hilang begitu kamu tahu satu baris ini.

Gerakan dasar di Normal mode:

```
h j k l    kiri, bawah, atas, kanan
w          loncat ke kata berikutnya
b          loncat ke kata sebelumnya
0          awal baris
$          akhir baris
gg         awal file
G          akhir file
dd         hapus satu baris
yy         salin satu baris
p          tempel
```

Kuasai ini dulu sebelum mengonfigurasi apa pun. Editor secanggih apa pun tidak berguna kalau kamu belum bisa bergerak di dalamnya.

## Struktur Konfigurasi

Neovim membaca konfigurasi dari `~/.config/nvim/`. Berbeda dengan Vim lama yang memakai bahasa scriptnya sendiri, Neovim memakai **Lua**, bahasa pemrograman yang sungguhan dan jauh lebih mudah dibaca.

File utamanya adalah `~/.config/nvim/init.lua`. Mulai dengan pengaturan dasar yang membuat editor langsung lebih nyaman:

```lua
-- ~/.config/nvim/init.lua

-- Tampilkan nomor baris
vim.opt.number = true
vim.opt.relativenumber = true

-- Indentasi
vim.opt.tabstop = 4
vim.opt.shiftwidth = 4
vim.opt.expandtab = true

-- Pencarian
vim.opt.ignorecase = true
vim.opt.smartcase = true

-- Tampilan
vim.opt.termguicolors = true
vim.opt.scrolloff = 8
```

`relativenumber` menampilkan nomor baris relatif terhadap posisi kursor, yang memudahkan perpindahan cepat. `smartcase` membuat pencarian mengabaikan huruf besar-kecil kecuali kamu mengetik huruf kapital. `scrolloff = 8` menjaga selalu ada delapan baris di atas dan bawah kursor agar konteks tetap terlihat. Pengaturan kecil ini sudah mengubah pengalaman tanpa satu plugin pun.

## Menambahkan Plugin

Kemampuan modern seperti autocomplete, penyorotan sintaks yang lebih baik, dan pencarian file datang dari plugin. Neovim modern memakai plugin manager untuk mengelolanya. Salah satu yang populer adalah `lazy.nvim`.

Filosofi yang penting di sini: **tambahkan plugin hanya saat kamu merasakan kebutuhannya**, bukan menyalin konfigurasi orang lain yang berisi lima puluh plugin sekaligus. Setiap plugin yang kamu pasang sebaiknya menyelesaikan masalah nyata yang kamu alami.

Beberapa yang hampir selalu berguna:

- **telescope.nvim**: pencarian file dan teks yang sangat cepat
- **nvim-treesitter**: penyorotan sintaks yang memahami struktur kode
- **nvim-lspconfig**: dukungan Language Server untuk autocomplete dan deteksi error

Mulai dengan satu, pahami cara kerjanya, baru tambahkan berikutnya. Konfigurasi yang tumbuh perlahan jauh lebih mudah dirawat daripada yang disalin utuh dari internet lalu tidak pernah dipahami.

## Soal LSP

Language Server Protocol adalah yang memberi editor kemampuan "pintar": autocomplete yang sadar konteks, melompat ke definisi fungsi, menampilkan error saat mengetik. Inilah yang membuat VS Code terasa canggih, dan Neovim bisa melakukan hal yang sama lewat `nvim-lspconfig`.

Setiap bahasa punya language servernya sendiri yang perlu dipasang terpisah, misalnya `gopls` untuk Go atau `pyright` untuk Python. Setelah terpasang dan dikonfigurasi, Neovim memberi pengalaman pengembangan yang setara dengan editor besar, tapi dengan kendali penuh atas setiap bagiannya.

## Bersabar di Minggu Pertama

Jujur saja: minggu pertama akan terasa lebih lambat dibanding editor lama kamu. Tanganmu akan otomatis meraih mouse, lalu teringat tidak perlu. Kamu akan lupa menekan `i` sebelum mengetik. Ini normal dan dialami semua orang.

Tapi otot ingatan terbentuk lebih cepat dari dugaan. Dalam dua sampai tiga minggu, gerakan mulai terasa alami. Dan begitu alami, mengedit teks tanpa pernah mengangkat tangan dari keyboard menjadi sesuatu yang sulit ditinggalkan.

---

Neovim bukan untuk semua orang, dan itu tidak masalah. Tapi reputasinya yang menakutkan sebagian besar tidak beralasan. Kamu tidak perlu menguasai semuanya sekaligus, tidak perlu menyalin konfigurasi rumit milik orang lain, dan tidak akan terjebak tanpa bisa keluar. Mulai dari gerakan dasar, tambahkan satu pengaturan dan satu plugin pada satu waktu, dan biarkan editor itu tumbuh mengikuti kebutuhanmu. Hasilnya adalah alat yang benar-benar kamu pahami sampai ke akarnya, sesuatu yang jarang bisa dikatakan tentang perkakas yang kita pakai setiap hari.
