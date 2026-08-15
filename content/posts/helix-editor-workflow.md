---
title: "Helix: A Selection-First Terminal Text Editor"
slug: c5e9b1d3
aliases: [helix-editor-workflow]
date: 2026-06-30
tags: [tools, rust, terminal]
lang: en
draft: false
---

Helix is a terminal text editor written in Rust that ships with built-in LSP support, tree-sitter parsing, and multiple cursors. It requires no plugin system because its core feature set covers what most developers need from a daily editor.

## Fun Facts

**Fact 1.** Helix uses a selection-first editing model derived from Kakoune, not Vim. In Vim you type an operator (`d`, `c`, `y`) and then a motion. In Helix you make a selection first and then act on it, which means you always see exactly what will be affected before committing a change.

**Fact 2.** Tree-sitter integration is built in at the parser level, not bolted on as a plugin. This means syntax highlighting is accurate for nested languages (HTML inside JavaScript inside a template literal, for example) and structural navigation works without configuration.

**Fact 3.** Helix ships as a single static binary. The entire editor, including the grammar library and default themes, weighs under 15 MB on Linux. There is no package manager, no init.lua, and no startup plugin loading sequence.

---

## Tips and Tricks

### 1. Understand the Selection-First Model

Every editing action in Helix begins with a selection. The cursor is always the tail of a selection, even if the selection covers only one character. This changes how you think about editing:

| Task | Vim | Helix |
|---|---|---|
| Delete a word | `dw` | `w` then `d` |
| Change inside parens | `ci(` | `mi(` then `c` |
| Delete to end of line | `D` | `gl` then `d` |
| Yank a line | `yy` | `x` then `y` |

In Helix, `x` selects the entire current line. Running `x` multiple times extends the selection one line at a time. This is more predictable than Vim's line-oriented commands.

### 2. Install Helix

Pre-built binaries for Linux are available on the GitHub releases page. Alternatively, build from source:

```bash
git clone https://github.com/helix-editor/helix
cd helix
cargo build --release --locked
# Binary is at ./target/release/hx
# Copy runtime directory alongside the binary
cp -r runtime ~/.config/helix/
```

On Arch Linux, install from the official repositories:

```bash
sudo pacman -S helix
```

Check the installation:

```bash
hx --version
# helix 24.03 (ab1d33f3)
```

### 3. Configure LSP for a Language

Helix picks up language servers that are already installed on `$PATH`. For Rust:

```bash
rustup component add rust-analyzer
```

For Go:

```bash
go install golang.org/x/tools/gopls@latest
```

Open a source file in Helix and press `Space l` to open the language server log. If the server is running, diagnostics appear inline. No configuration file is needed for standard language servers.

### 4. Essential Keybindings for Daily Editing

```
Normal mode:
  w / b / e     move by word forward / backward / end
  f<char>        find next occurrence of char on line
  t<char>        move to just before char on line
  mi(            select inside parentheses (m = match)
  ma{            select around curly braces
  %              select entire file
  C              add cursor on next matching selection
  ,              collapse to primary cursor
  &              align cursors
  Space f        fuzzy-find file (requires file picker)
  Space s        fuzzy-find symbol in current file
  ]d / [d        jump to next / previous diagnostic
  Space a        code action (LSP)
  K              show hover documentation (LSP)
```

The `C` key for adding cursors is the fastest path to multi-cursor editing: select a word with `w`, then press `C` repeatedly to clone the cursor on each subsequent occurrence.

### 5. Startup Time vs. Neovim

Helix starts faster than a fully configured Neovim instance because there is no Lua runtime initialization or plugin loading:

```bash
# Measure cold-start time to first rendered frame (averaged over 20 runs)
hyperfine --warmup 3 \
  'hx --quit' \
  'nvim --headless +qa'
```

On a typical developer laptop, Helix opens in 8-15 ms. A lean Neovim setup with lazy-loading takes 40-80 ms; a heavier configuration with many plugins regularly exceeds 200 ms. For quick file edits this difference is noticeable.

### 6. Minimal Configuration

Helix configuration lives at `~/.config/helix/config.toml`. A practical starting point:

```toml
[editor]
line-number = "relative"
cursorline = true
color-modes = true
auto-save = true

[editor.cursor-shape]
insert = "bar"
normal = "block"
select = "underline"

[editor.indent-guides]
render = true

[keys.normal]
# Use space+q to quit without saving (mirrors :q!)
"space" = { q = ":q!" }
```

Run `hx --health` to see which language servers and tree-sitter grammars are detected on the current system.
