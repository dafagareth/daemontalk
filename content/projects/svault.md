## Overview

**svault** is a local, encrypted secret vault for the command line. It keeps API
keys, passwords, database URLs, and any other sensitive string inside a single
encrypted file on your own machine. No cloud sync, no telemetry, no account, no
running server required.

The entire tool ships as one static binary that runs on Linux, macOS, and
Windows. Unlock it once with a master password, work for a 30-minute session,
and it locks itself automatically when the timer expires.

```
$ svault set DB_PASSWORD supersecret123
OK  DB_PASSWORD saved

$ svault get DB_PASSWORD
supersecret123
```

> **Naming note:** svault was originally called `stash`. It was renamed in
> v2.0.0 to avoid clashing with `git stash` and an unrelated AUR package.
> The command, vault directory (`~/.svault`), and environment variables all
> use the `svault` prefix from v2.0.0 onward.

## Why svault?

Most developers manage dozens of secrets across multiple projects. Plain `.env`
files risk leaking into git history, full-featured password managers are
excessive for CLI scripts, and cloud-based vaults introduce network latency and
external account dependencies.

| Problem | How svault solves it |
|---|---|
| Plain-text `.env` files can be committed to Git | Every value is encrypted with AES-256-GCM on disk |
| Enterprise secret managers are too heavy | One binary, no server, no daemon, no setup |
| Secrets scattered across many projects | One isolated namespace per project inside a single vault |
| No audit trail of secret access | A local, append-only audit log records every operation |
| Retyping master passwords constantly | Unlock once, work for up to 30 minutes |

## Install

### Arch Linux (AUR)

```bash
yay -S svault
# or
paru -S svault
```

### Quick install: Linux / macOS

```bash
curl -fsSL https://raw.githubusercontent.com/dafagareth/svault/main/install.sh | sh
```

The script detects your OS and architecture, downloads the matching release
binary, verifies its checksum, and places it in `/usr/local/bin` (or
`~/.local/bin` if that path is not writable). Override the version with
`SVAULT_VERSION=v2.0.0`, or the install location with `SVAULT_BIN_DIR`.

### Quick install: Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/dafagareth/svault/main/install.ps1 | iex
```

Installs to `%LOCALAPPDATA%\svault` and adds that path to your user `PATH`.

### Download a release binary

Pick the correct file for your platform from the
[Releases](https://github.com/dafagareth/svault/releases) page:

| Platform | File |
|---|---|
| Linux x86_64 | `svault-linux-amd64` |
| Linux ARM64 | `svault-linux-arm64` |
| macOS Intel | `svault-darwin-amd64` |
| macOS Apple Silicon | `svault-darwin-arm64` |
| Windows x86_64 | `svault-windows-amd64.exe` |

```bash
chmod +x svault-linux-amd64
sudo mv svault-linux-amd64 /usr/local/bin/svault
```

### Build from source

```bash
git clone git@github.com:dafagareth/svault.git
cd svault
sudo make install               # builds and installs to /usr/local/bin
sudo make install PREFIX=/usr   # alternative install prefix
sudo make uninstall             # removes the binary
```

**Requirements:** Go 1.25 or later. Homebrew (macOS) and Scoop (Windows)
manifests exist but are not yet published. Use the install scripts above in
the meantime.

## Quick Start

```bash
# Step 1: initialize the vault once. You will be prompted for a master password.
$ svault init
Master password: ********
Confirm password: ********
Vault initialized at ~/.svault/vault.enc

# Step 2: unlock a 30-minute session.
$ svault unlock
Vault unlocked. Session valid for 30 minutes.

# Step 3: store secrets.
$ svault set DB_PASSWORD supersecret123
$ svault set JWT_SECRET=myjwtsecret         # KEY=VALUE syntax also works

# Step 4: use them.
$ svault get DB_PASSWORD                    # print a single value
$ svault list                               # list all keys (values stay hidden)
$ svault export > .env                      # write a standard .env file
$ svault exec -- npm start                  # inject secrets directly into a command

# Step 5: lock when done (or let the session expire naturally).
$ svault lock
```

## Command Reference

### Initialization and authentication

```bash
svault init      # create a new vault, prompts for a master password
svault unlock    # unlock the vault, starts a 30-minute session
svault lock      # lock the vault immediately
svault status    # show whether the vault is locked or unlocked
```

`svault status --short` prints a compact representation (lock icon or remaining
minutes) designed for use inside shell prompts.

### Secret management

```bash
svault set KEY VALUE      # store a secret (KEY=VALUE syntax also works)
svault set KEY --stdin    # read the value from stdin, keeping it out of shell history
svault get KEY            # retrieve a stored value
svault edit KEY           # open a value in $EDITOR, ideal for multiline secrets
svault delete KEY         # permanently remove a secret
svault list               # list all keys in the active namespace
svault search PATTERN     # search keys by name; use --all to search every namespace
svault rename OLD NEW     # rename a key while preserving its value
svault move KEY --to NS   # move a key to a different namespace
```

Reading from stdin protects sensitive values from shell history:

```bash
echo "supersecret" | svault set DB_PASSWORD --stdin
pbpaste | svault set API_KEY --stdin       # pipe from clipboard on macOS
```

### Clipboard and password generation

```bash
svault copy KEY                        # copy a value, auto-clears clipboard after 30s
svault generate                        # generate a random password, copy to clipboard
svault generate --length 32 --save DB_PASSWORD
svault generate --no-symbols           # alphanumeric only
svault open GITHUB                     # open GITHUB_URL in browser and copy GITHUB_PASS
```

`svault open KEY` resolves `KEY_URL` (or uses `KEY` directly if it holds an HTTP
URL), opens it in the default browser, then copies `KEY_PASS`, `KEY_PASSWORD`,
or `KEY_TOKEN` to the clipboard if any of those exist.

### Shell integration

```bash
svault exec -- npm start               # run a command with secrets injected as env vars
svault exec --ns production -- ./deploy.sh
eval $(svault env)                     # load all secrets into the current shell session
eval $(svault env --ns production)
```

### Namespace management

```bash
svault use NAMESPACE           # set the active namespace
svault ns list                 # list all namespaces with key counts
svault ns rename OLD NEW       # rename a namespace
svault ns delete NAMESPACE     # delete a namespace and all its secrets
svault diff staging production # compare two namespaces side by side
```

Example output of `svault diff`:

```
= DB_PASSWORD                   same
< JWT_SECRET                    only in [staging]
> STRIPE_KEY                    only in [production]
~ REDIS_URL                     value differs

[staging] vs [production]: 1 same, 3 differ
```

### Import, export, and verification

```bash
svault export > .env                   # export active namespace as KEY=VALUE lines
svault export --ns production > .env.production
svault import .env.example             # import secrets from an existing .env file
svault check                           # verify vault against .env.example
svault check .env.production           # verify against a different file
```

Example output of `svault check`:

```
OK    DB_PASSWORD
OK    JWT_SECRET
MISS  STRIPE_KEY

2/3 keys present in namespace [default], 1 missing
```

### Maintenance and diagnostics

```bash
svault rotate                  # change the master password and re-encrypt the vault
svault backup ~/safe/vault.bak
svault backup                  # auto-named, timestamped backup saved to ~/.svault/
svault restore ~/safe/vault.bak
svault info                    # vault version, path, active namespace, key count, session state
svault log                     # display the audit log
svault doctor                  # health-check the installation and vault
svault version                 # print the version number (--short for number only)
```

After `svault rotate`, the previous password no longer works. The active session
is refreshed automatically, so no re-unlock is needed.

Running `svault doctor` is the recommended first step when something seems wrong.
It checks the vault directory and file, file permissions, config and audit log,
session state, git availability for auto-namespace detection, and clipboard
tooling:

```
[OK  ] vault directory              /home/you/.svault
[OK  ] vault file                   412 bytes, mode 0600
[OK  ] config file                  /home/you/.svault/config.json
[OK  ] audit log                    /home/you/.svault/vault.log
[OK  ] session                      unlocked, 24m remaining
[OK  ] git (auto-namespace)         current namespace: grocyvo
[OK  ] clipboard                    wl-copy

All checks passed.
```

The commands `doctor`, `init`, `version`, `completion`, and `help` run without
an existing vault. All other commands exit early with a descriptive message if
no vault is found.

## Namespaces

Each project gets its own isolated namespace, so two projects can each define a
`DB_PASSWORD` without conflict.

### Automatic detection

Inside a git repository, the namespace is derived from the repo name
automatically. No manual `svault use` is required.

```bash
~/grocyvo$ svault set DB_URL=postgres://localhost/grocyvo
~/grocyvo$ svault info
Namespace  : grocyvo (from git, 2 total)

~/portfolio$ svault get DB_URL     # reads from the 'portfolio' namespace instead
```

Detection priority, from highest to lowest:

1. The `--ns` flag on the current command
2. The `SVAULT_NS` environment variable
3. The current git repository name
4. The active namespace set with `svault use`
5. `default`

### Manual control

```bash
svault use production               # switch the active namespace
svault set --ns staging DB_URL=...  # one-off override without switching namespaces
svault ns list                      # view all namespaces and their key counts
```

## Sessions

When you unlock the vault, svault derives the encryption key from your master
password and writes it to a session file. This allows subsequent commands to run
without retyping the password.

- Default session length: **30 minutes**, configurable via `SVAULT_SESSION_TTL`.
- The session ends when the TTL expires or when you run `svault lock`.
- The session file is stored at `/tmp/.svault_session` with permission `0600`
  and is deleted automatically on expiry.

## Storage Layout

```
~/.svault/
  vault.enc       # your secrets, encrypted with AES-256-GCM (JSON structure inside)
  vault.log       # append-only audit log
  config.json     # active namespace setting

/tmp/
  .svault_session # derived session key, mode 0600, auto-deleted after TTL
```

The raw layout of `vault.enc` on disk:

```
[ 16 bytes salt ][ 12 bytes nonce ][ ciphertext ]
```

## Configuration

| Environment variable | Default | Purpose |
|---|---|---|
| `SVAULT_SESSION_TTL` | `30` | Session duration in minutes |
| `SVAULT_NS` | (unset) | Force a specific namespace, overriding git detection |

## Shell Completion

svault generates completion scripts for Bash, Zsh, and Fish:

```bash
# Bash
svault completion bash | sudo tee /etc/bash_completion.d/svault > /dev/null

# Zsh
svault completion zsh > "${fpath[1]}/_svault"

# Fish
svault completion fish > ~/.config/fish/completions/svault.fish
```

Restart your shell or source the generated file directly to activate completion.

## Security Model

svault is built on a small set of deliberate cryptographic choices:

- **Encryption: AES-256-GCM.** An authenticated cipher. Any tampering with the
  vault file is detected at decryption time before any data is exposed.
- **Key derivation: Argon2id.** The current recommended algorithm for
  password-based key derivation. Each vault receives a fresh, randomly generated
  16-byte salt.
- **Randomness: `crypto/rand`.** All salts, nonces, and generated passwords use
  the cryptographically secure RNG. `math/rand` is never used for anything
  security-related.

The salt is stored in plaintext inside `vault.enc`. This is standard practice:
the salt is not a secret. Its role is to make precomputed (rainbow table) attacks
against the master password infeasible.

### Security guarantees

- **The master password cannot be recovered.** Losing it means losing access to
  the vault permanently, by design. Keep a copy of the password in a safe
  location.
- Secret **values** are never written to logs, temp files, or standard output
  except through `get`, `copy`, `env`, and `export`.
- Every write operation first creates a `vault.enc.bak` and several timestamped
  rollback copies. A failed write cannot destroy an earlier clean state.
- Concurrent svault processes are serialized using an **exclusive file lock**.
  Two simultaneous writes cannot corrupt each other.
- The session file is created with permission `0600`, readable only by the
  owning user.
- Secret keys must conform to valid shell variable naming rules, ensuring that
  `export` and `env` output is always safe to `eval`.

### Known trade-off: session key in /tmp

While the vault is unlocked, the derived encryption key is stored as plaintext
in `/tmp/.svault_session` (permission `0600`). This is an intentional
convenience trade-off to avoid retyping the master password constantly.

- Regular users on the same machine **cannot** read this file.
- **Root, or any process with root-equivalent privileges, can.**
- The file is deleted automatically after the TTL or immediately upon `svault lock`.

On shared servers or systems where other users have root access, run
`svault lock` as soon as you finish working rather than waiting for the TTL.

## Common Scenarios

### Replace a project's .env file

```bash
cd ~/myproject             # namespace becomes "myproject" automatically
svault import .env         # import existing secrets from the file
rm .env                    # remove the plaintext file
svault exec -- npm run dev # run with secrets injected, no file on disk
```

### Onboard to a project with .env.example

```bash
svault check               # which keys from .env.example are not yet set?
# OK    DATABASE_URL
# MISS  STRIPE_KEY         <- set this one next
svault set STRIPE_KEY=sk_test_...
```

### Keep secrets out of shell history

```bash
echo "$GENERATED" | svault set API_KEY --stdin   # value never appears as an argument
svault edit MULTILINE_CERT                        # opens $EDITOR for long values
```

### Promote secrets from staging to production

```bash
svault diff staging production       # review differences before touching anything
svault move STRIPE_KEY --to production
```

### CI and automated deployments

```bash
export SVAULT_NS=production
svault exec -- ./deploy.sh           # secrets are scoped to that process only
```

### Rotate the master password

```bash
svault rotate   # prompts for old and new passwords, re-encrypts the entire vault
                # the old password is invalidated immediately; session is refreshed
```

### Migrate from the old stash (pre-v2.0)

```bash
mv ~/.stash ~/.svault    # move the vault directory to the new location
svault unlock            # unlock as usual with the same master password
```

## FAQ

```faq
Q: I forgot my master password. Can I recover my secrets?
A: No. There is no recovery path, backdoor, or reset mechanism. This is intentional: a recovery option would weaken the security guarantee that only someone with the password can read the vault. Your only options are to restore from a backup made while you still knew the password, or to re-initialize the vault and re-enter your secrets. Keep your master password in a separate safe location.

Q: Does svault send any data to the internet?
A: Never. There is no network code in the primary execution path. All data stays in `~/.svault` on your local machine. The only network activity that could occur is when the install scripts download the binary from GitHub Releases during initial setup, or when you open a URL stored in the vault with `svault open`.

Q: Can two different projects use the same secret name?
A: Yes. Namespaces keep them completely separate. `myproject/DB_PASSWORD` and `other-project/DB_PASSWORD` are independent values. When you run svault inside a git repository, the namespace is set automatically from the repo name, so there is nothing special to configure.

Q: Is it safe to use svault on a shared server?
A: With caution. The vault file and session file are both created with permission `0600`, which prevents other regular users from reading them. However, anyone with root access on the machine can read the session key from `/tmp/.svault_session` while the vault is unlocked. On a shared or multi-user server, always run `svault lock` immediately after you finish working instead of relying on the session TTL.

Q: Does svault require a daemon or background service?
A: No. svault is a plain command-line tool. It starts, performs its operation, and exits. The only persistent state is the encrypted vault file, the audit log, and the short-lived session file in `/tmp`. Nothing runs in the background between commands.

Q: What happens if the session file in /tmp is deleted?
A: svault treats a missing session file the same as an expired session. The next command that requires an unlocked vault will refuse to run and prompt you to run `svault unlock` again. Your vault and all stored secrets remain intact; only the in-memory session state is lost.

Q: Can I use svault in a CI/CD pipeline or a Docker container?
A: Yes, with some planning. The recommended approach is to unlock the vault at the start of the pipeline script and use `svault exec` to inject secrets into individual commands. Set `SVAULT_NS` to the correct namespace and use `SVAULT_SESSION_TTL` to extend the session length if your pipeline takes longer than 30 minutes. The session file will be scoped to the container or runner, so it will not persist between runs.

Q: How do I share secrets with a teammate?
A: svault is a local-only tool and does not have a built-in sharing mechanism. To share secrets, export them to a `.env` file with `svault export` and transfer the file through a secure channel (an encrypted message, a secrets manager, or a secure file share). The recipient can then import the file with `svault import`. Never send the `.env` file over email or unencrypted chat.

Q: What happens if I run svault init when a vault already exists?
A: svault detects the existing vault file and refuses to overwrite it without confirmation. Your existing secrets are not affected unless you explicitly confirm the re-initialization. If you want to start fresh, remove `~/.svault/vault.enc` manually first, then run `svault init`.

Q: Can I have multiple vaults for different purposes?
A: Not directly. svault uses a single `vault.enc` file and separates secrets by namespace within that file. For most use cases, namespaces are sufficient: each project gets its own isolated namespace, and you can create as many as you need. If you genuinely need a completely separate vault with a different master password, you would need to run a second instance with a different home directory or manage the vault files manually.

Q: What clipboard tools does svault support?
A: svault detects the available clipboard tool based on your environment: **Linux (Wayland):** `wl-copy` (from `wl-clipboard`) **Linux (X11):** `xsel` or `xclip` **macOS:** `pbcopy` (built in) Run `svault doctor` to confirm which tool was detected. If none is found, clipboard commands (`copy`, `generate`, `open`) will fail with a descriptive error.

Q: What happens if two svault processes run at the same time?
A: svault uses an exclusive file lock on every write operation. If two processes attempt to write simultaneously, the second one waits until the first completes and then proceeds. Reads are not locked and can run concurrently. This prevents the race condition where two concurrent writes corrupt each other and silently lose data.

Q: How do I back up my vault?
A: Use the built-in backup command: `svault backup ~/safe/vault.bak` copies the vault to a specific path. `svault backup` (no path) creates a timestamped copy in `~/.svault/` automatically. svault also creates an automatic rollback copy on every write, so a bad write never destroys the previous good state. For long-term disaster recovery, copy `~/.svault/vault.enc` to external storage or a separate encrypted backup location.

Q: Can I use svault without git installed?
A: Yes. Git is only used for automatic namespace detection. If git is not installed or you are not inside a git repository, svault falls back to the namespace set by `svault use`, the `SVAULT_NS` environment variable, or the `default` namespace. All other features work normally. `svault doctor` will report a warning about missing git, but no commands will fail because of it.

Q: Can I automate svault in shell scripts without interactive password prompts?
A: The intended workflow is to call `svault unlock` once interactively at the start of your session, then use all other commands non-interactively within the 30-minute window. Inside a CI pipeline where interactive input is not possible, unlock using a heredoc or piped input, then use `SVAULT_SESSION_TTL` to keep the session alive for the duration of the pipeline. For fully automated environments, consider whether a dedicated secrets manager with machine-identity authentication is more appropriate.

Q: What are the Argon2id settings svault uses?
A: svault uses Argon2id tuned for roughly 200ms of derivation time on a modern laptop. This provides meaningful resistance to brute-force and dictionary attacks while keeping the unlock command fast enough for daily interactive use. The exact parameters (memory, iterations, parallelism) are stored alongside the salt so that future versions can change defaults without breaking existing vaults.
```

## Tech Stack

- **Go 1.25+**, compiled to a single static binary with no CGO and no runtime dependencies
- **`golang.org/x/crypto`** for Argon2id key derivation and AES-256-GCM encryption
- **Cobra** for the command tree and shell completion generation
- Cross-compiled for Linux, macOS, and Windows (amd64 and arm64) in CI via GitHub Actions

## Status

Shipped at v2.0.0. Packaged and available on the AUR. Homebrew and Scoop
manifests are written but not yet submitted to their respective registries.
Source code, release binaries, the full security policy, and the changelog
are available in the [GitHub repository](https://github.com/dafagareth/svault).
