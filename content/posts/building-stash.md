---
title: "Building stash: an encrypted secret vault in Go"
slug: 59cf9f8b
aliases: [building-stash]
date: 2026-06-05
tags: [go, cli, security]
lang: en
draft: false
series: building-stash
series_part: 1
---

Every developer has been there: you clone a repo, create a `.env` file, fill it with API keys and database passwords, and then spend the next five minutes hoping you haven't accidentally staged it.

I got tired of that workflow. So I built **stash**, a local encrypted vault for secrets.

## The core idea

One binary. No server. No cloud. Your secrets live in `~/.stash/vault.enc`, encrypted with AES-256-GCM and a key derived from your master password via Argon2id.

```bash
$ stash init
$ stash unlock
$ stash set DB_PASSWORD supersecret
$ stash exec -- npm start
```

That last command injects all secrets as environment variables into the child process, with no `.env` file needed at all.

## What I learned

**Argon2id** is the right choice for password-based key derivation. The memory-hardness makes brute force attacks expensive even with GPUs. I set memory to 64MB, time to 3 iterations, slow enough to be secure and fast enough that unlocking feels instant.

**AES-256-GCM** provides authenticated encryption. If anyone tampers with `vault.enc`, decryption fails and the vault refuses to open. This is the right default: you never want silent corruption.

**Session design** is where the trade-offs get interesting. The derived key is stored in `/tmp/.stash_session` with `0600` permissions. Root can read it. That's a known trade-off: convenience vs. perfect security. I documented it explicitly in the README rather than pretending it doesn't exist.

## What's next

The vault works well for local development. I'm thinking about namespaces for multi-project workflows and a `check` command to validate which `.env.example` keys are missing from the vault.

Source: [github.com/dafagareth/svault](https://github.com/dafagareth/svault)
