---
title: "Nix and Reproducible Builds"
slug: e4b7d9c8
aliases: [nix-reproducible-builds]
date: 2026-05-02
tags: [devops, tools, linux]
lang: en
draft: false
---

Nix is a package manager and build system built around a single principle: the output of a build is determined entirely by its inputs. Given the same inputs, every build produces bit-for-bit identical output, on any machine, at any time.

## Fun Facts

**Fact 1.** Every package in the Nix store lives at a path like `/nix/store/7dq3h...glibc-2.39/`. The hash prefix is derived from the complete closure of inputs: source, compiler, flags, dependencies, and their dependencies. Changing any input changes the hash, creating a new path rather than overwriting the old one.

**Fact 2.** Nix evaluates package descriptions written in the Nix expression language, which is a pure, lazy, dynamically typed functional language. Side effects (network access, file writes) are forbidden during evaluation; they are only permitted during the sandboxed build phase.

**Fact 3.** NixOS is a Linux distribution where the entire operating system configuration, from kernel parameters to installed services and their configs, is described in a single set of Nix files and is fully reproducible. Rolling back to any previous system generation takes one command.

---

## Tips and Tricks

### 1. Understand the Core Problem Nix Solves

The "works on my machine" problem has two common forms. First, implicit dependencies: a build works because a library happens to be installed globally, but the build description never declares it. Second, version drift: the same `apt install` command installs different versions on different machines over time.

Nix solves both by making all dependencies explicit in derivations and by storing every version in a content-addressed store. Nothing is implicit; nothing is mutable.

### 2. Install Nix

The official installer works on any Linux distribution and on macOS:

```bash
sh <(curl -L https://nixos.org/nix/install) --daemon
```

The `--daemon` flag sets up a multi-user installation where the Nix daemon runs as root but individual users can install packages without elevated privileges. After installation, source the profile:

```bash
. /etc/profile.d/nix.sh
```

Verify:

```bash
nix --version
# nix (Nix) 2.24.1
```

### 3. Write a Minimal flake.nix for a Go Project

Nix Flakes are the modern way to declare reproducible environments. A minimal `flake.nix` for a Go project:

```nix
{
  description = "My Go service";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in {
        # Development shell: nix develop
        devShells.default = pkgs.mkShell {
          buildInputs = [
            pkgs.go_1_22
            pkgs.gopls
            pkgs.golangci-lint
            pkgs.delve
          ];
          shellHook = ''
            export GOPATH="$HOME/.cache/go"
            echo "Go $(go version | awk '{print $3}')"
          '';
        };

        # Reproducible build: nix build
        packages.default = pkgs.buildGoModule {
          pname = "myservice";
          version = "0.1.0";
          src = ./.;
          vendorHash = null; # set to actual hash after first build
        };
      });
}
```

Pin the exact revision by recording the lock file:

```bash
nix flake lock          # creates flake.lock
git add flake.nix flake.lock
```

### 4. Use nix develop as a Portable Dev Shell

`nix develop` drops you into a shell with exactly the tools declared in `devShells.default`, regardless of what is installed globally on the system:

```bash
# Enter the dev shell
nix develop

# Run the project
go build ./...
go test ./...

# Exit returns you to the normal shell
exit
```

To activate the shell automatically when entering a directory, use `direnv` with the Nix integration:

```bash
# Install direnv
nix profile install nixpkgs#direnv

# Create .envrc in the project root
echo 'use flake' > .envrc
direnv allow
```

After this, entering the directory automatically activates the Nix shell and leaving it deactivates it, with no manual intervention.

### 5. How NixOS Takes This to the OS Level

On NixOS, the entire system is a derivation. The configuration typically lives in `/etc/nixos/configuration.nix`:

```nix
{ config, pkgs, ... }: {
  # Networking
  networking.hostName = "myserver";
  networking.firewall.allowedTCPPorts = [ 80 443 ];

  # Services
  services.nginx.enable = true;
  services.postgresql = {
    enable = true;
    package = pkgs.postgresql_16;
  };

  # System packages
  environment.systemPackages = with pkgs; [
    htop
    git
    ripgrep
  ];

  system.stateVersion = "24.05";
}
```

Apply the configuration and switch atomically:

```bash
sudo nixos-rebuild switch
```

If the new configuration has a problem, roll back to the previous generation:

```bash
sudo nixos-rebuild switch --rollback
# or select a generation from the boot menu
```

Every generation is kept in the bootloader menu by default, making system recovery straightforward without needing a live USB.
