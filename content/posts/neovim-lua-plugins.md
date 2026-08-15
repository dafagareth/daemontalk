---
title: "Neovim Lua: Trik Menulis Config Modular Tanpa Bloatware"
slug: 9a0b1c2d
aliases: [neovim-lua-plugins]
date: 2026-06-08
tags: [neovim, editor, tips]
lang: id
draft: false
---

Migrasi dari Vimscript tradisional ke Lua adalah lompatan terbesar dalam ekosistem Neovim. Lua adalah bahasa scripting yang sangat cepat, ringan, dan membuat manajemen plugin menjadi modular serta menyenangkan.

Berikut tips praktis menyusun `init.lua` yang bersih, cepat (startup di bawah 25ms), dan minim bloatware.

## Fun Fact

**LuaJIT adalah salah satu dynamic language JIT compiler tercepat di dunia.**
Neovim memanfaatkan LuaJIT yang dibuat oleh Mike Pall. Kecepatan eksekusi LuaJIT bisa puluhan kali lebih cepat dibanding interpreter Vimscript lama.

**Neovim bermula dari proyek fork Vim oleh Thiago de Arruda pada 2014.**
Thiago memulai Neovim setelah pull request-nya untuk menambahkan fitur async job control ditolak oleh pencipta Vim, Bram Moolenaar.

**Struktur `lua/` di runtimepath otomatis dikenali Neovim.**
Setiap file Lua di bawah folder `~/.config/nvim/lua/myconfig/` bisa langsung di-*import* dengan perintah `require('myconfig')`.

---

## Tips dan Trik

### 1. Struktur Folder Modular yang Rapi

Pisahkan konfigurasi berdasarkan tanggung jawab:

```text
~/.config/nvim/
├── init.lua
└── lua/
    ├── core/
    │   ├── options.lua
    │   └── keymaps.lua
    └── plugins/
        ├── lsp.lua
        └── treesitter.lua
```

### 2. Gunakan `lazy.nvim` untuk Lazy-Loading Otomatis

Manajer plugin modern seperti `lazy.nvim` hanya memuat plugin saat dibutuhkan (misal saat membuka filetype tertentu atau menekan keybinding):

```lua
-- lua/plugins/lsp.lua
return {
    "neovim/nvim-lspconfig",
    event = { "BufReadPre", "BufNewFile" },
    dependencies = {
        "hrsh7th/nvim-cmp",
    },
    config = function()
        -- setup gopls, rust-analyzer, etc.
    end
}
```

### 3. Profil Waktu Startup dengan `:Lazy profile`

Lacak plugin mana yang memperlambat waktu booting editor:

```vim
:Lazy profile
```

### 4. Setting Opsi Penting di `options.lua`

Optimalkan pengalaman scrolling dan editing:

```lua
local opt = vim.opt

opt.number = true
opt.relativenumber = true
opt.tabstop = 4
opt.shiftwidth = 4
opt.expandtab = true
opt.smartindent = true
opt.wrap = false
opt.termguicolors = true
opt.scrolloff = 8
opt.signcolumn = "yes"
```

### 5. Atur Leader Key Menjadi Spasi

Tombol spasi jauh lebih mudah dijangkau ibu jari dibanding tombol `\` bawaan Vim:

```lua
vim.g.mapleader = " "
vim.g.maplocalleader = " "

-- Contoh shortcut simpan file:
vim.keymap.set("n", "<leader>w", "<cmd>write<cr>", { desc = "Save file" })
```
