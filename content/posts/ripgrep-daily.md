---
title: "ripgrep: A Faster grep That Respects Your Time"
slug: 5ad9c07a
aliases: [ripgrep-daily]
date: 2025-07-08
tags: [linux, tools, cli]
lang: en
draft: false
---

`grep` ships on every Unix system and gets the job done. But searching a large codebase with it is slow, it walks into directories you do not care about (node_modules, .git, build output), and the output needs extra flags just to be readable.

`ripgrep` (`rg`) addresses all three. It uses Rust's regex engine, searches in parallel threads, respects `.gitignore` by default, and colors output without any configuration.

## Installation

```bash
# Arch Linux
pacman -S ripgrep

# Debian / Ubuntu
apt install ripgrep

# macOS
brew install ripgrep
```

## Basic Usage

The syntax is close enough to grep that the muscle memory transfers:

```bash
rg "pattern"           # search current directory
rg "pattern" src/      # search specific path
rg -i "config"         # case-insensitive
rg -w "err"            # whole word only
```

## Filtering by File Type

This is where rg earns its place. Instead of `--include="*.go"`:

```bash
rg "http.HandleFunc" -t go
```

List all supported types:

```bash
rg --type-list
```

Multiple types at once:

```bash
rg "import" -t go -t ts
```

## Fixed Strings

When your search pattern contains characters with special regex meaning (dots, parentheses, brackets), use `-F` to treat it literally:

```bash
rg -F "fmt.Println("
rg -F "err != nil"
```

## Useful Flags

| Flag | Effect |
|---|---|
| `-l` | Print file names only |
| `-c` | Count matches per file |
| `-A 3` | Show 3 lines after each match |
| `-B 3` | Show 3 lines before each match |
| `-C 3` | Show 3 lines around each match |
| `-g "*.go"` | Glob filter |
| `--no-ignore` | Search gitignored files too |
| `--hidden` | Include hidden files and directories |
| `-0` | NUL-separated output (for xargs -0) |

## Searching Ignored Paths

By default, rg skips anything in `.gitignore` and hidden directories. When you actually need to search there:

```bash
rg "pattern" --no-ignore --hidden
```

## Feeding Into Other Tools

```bash
# List files containing a pattern, then open them in your editor
rg -l "TODO" | xargs nvim

# Rename a symbol across all Go files
rg -l "OldFuncName" -t go | xargs sed -i 's/OldFuncName/NewFuncName/g'

# Count total matches across the project
rg "console.log" -c | awk -F: '{sum += $2} END {print sum}'
```

## One Alias Worth Setting

```bash
alias grep='rg'
```

It is not a perfect drop-in replacement (some grep flags differ), but for interactive use in a terminal, this alias means you always get the fast path without thinking about it.
