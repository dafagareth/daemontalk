---
title: "Building stash, part 2: namespaces, testing, and knowing when to stop"
slug: 33ad268c
aliases: [building-stash-part-2]
date: 2026-06-07
tags: [go, cli, security, testing]
lang: en
draft: false
series: building-stash
series_part: 2
---

In [part 1](/blog/building-stash) I built the core of stash: an encrypted vault with AES-256-GCM, Argon2id key derivation, and a session model. It worked, but it was a tool for one project at a time. This is the story of the second iteration, where it became something I actually use every day.

## Namespaces that detect themselves

The first real friction was switching between projects. Each project has its own `DB_URL`, its own `API_KEY`. Storing them all flat meant constant collisions.

Namespaces solved the storage side: each project gets its own isolated map of secrets. But I still had to run `stash use myproject` every time I changed directories. That command got old fast.

The fix was to detect the namespace automatically from the git repository:

```go
func gitRepoName() string {
    out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
    if err != nil {
        return ""
    }
    return filepath.Base(strings.TrimSpace(string(out)))
}
```

Now when I am inside `~/grocyvo`, stash uses the `grocyvo` namespace. Inside `~/portfolio`, it uses `portfolio`. No command needed.

The detection follows a clear priority order: an explicit `--ns` flag wins, then a `STASH_NS` environment variable, then the git repo name, then a manually set active namespace, and finally `default`. Each layer exists for a reason. The flag is for one-off overrides, the env var is for scripting and CI, and git detection is the everyday path.

## Testing the part that was scary to test

The commands all touch the filesystem: they read `~/.stash/vault.enc`, write a session token to `/tmp`, and shell out to git. That makes them awkward to test without trashing the real vault.

The trick was isolating every external dependency per test:

```go
func setupVault(t *testing.T) {
    home := t.TempDir()
    t.Setenv("HOME", home)
    t.Setenv("STASH_NS", "default")
    storage.SetSessionFile(filepath.Join(home, ".session"))
    // init a fresh vault and unlock a session...
}
```

Setting `HOME` to a temp directory redirects the whole vault path. Setting `STASH_NS` makes namespace detection deterministic, so a test never accidentally picks up the git repo it runs inside. A small exported `SetSessionFile` hook isolates the session token.

With that helper in place, testing a command becomes ordinary:

```go
func TestRunSet_EqualsForm(t *testing.T) {
    setupVault(t)
    runSet(newCmd(), []string{"API_KEY=sk-123"})
    got, _ := readSecret(t, "default", "API_KEY")
    if got != "sk-123" {
        t.Errorf("got %q, want sk-123", got)
    }
}
```

The command package went from zero tests to covering the real behaviour: both `KEY VALUE` and `KEY=VALUE` forms, namespace moves, the protected `default` namespace, password generation length and charset. Coverage is not a goal in itself, but writing these tests caught two genuine edge cases I had not considered, including a value containing an `=` sign being split in the wrong place.

## Keeping secrets out of shell history

One small feature I am happy with: reading a value from stdin.

```bash
echo "supersecret" | stash set DB_PASSWORD --stdin
```

Typing `stash set DB_PASSWORD supersecret` puts the secret straight into your shell history file in plain text. The `--stdin` flag avoids that entirely. It is a five-line change that removes a real footgun.

## Knowing when to stop

The harder lesson was restraint. Once the foundation was solid, it was tempting to keep adding commands: a password generator, clipboard support, a browser opener, a guestbook of features. I added the ones that earned their place and stopped.

A tool that does eight things well is more trustworthy than one that does twenty things vaguely. For a security tool especially, every feature is surface area. The discipline of saying "this is enough" is part of the engineering, not separate from it.

The full source, with tests and documentation split across README, CONTRIBUTING, and SECURITY, is on [GitHub](https://github.com/dafagareth/svault).
