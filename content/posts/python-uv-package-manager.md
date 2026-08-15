---
title: "uv: The Fast Python Package Manager"
slug: e9b4f712
aliases: python-uv-package-manager
date: 2026-06-05
tags: [python, tools, rust]
lang: en
draft: false
---

`uv` is a Python package manager and project tool written in Rust by Astral, the team behind the `ruff` linter. It replaces `pip`, `pip-tools`, `virtualenv`, `pyenv`, and parts of `poetry` with a single binary. The speed difference is significant: operations that take seconds with `pip` complete in milliseconds with `uv`, because it parallelizes downloads and uses a shared global cache.

## Fun Facts

**Fact 1.** Astral published benchmarks showing `uv` resolving and installing Django in 11 milliseconds compared to 3.1 seconds for `pip`. The difference comes from Rust's concurrency model and an efficient dependency resolver adapted from the `pubgrub` algorithm.

**Fact 2.** `uv` maintains a global package cache at `~/.cache/uv` by default. When you create a new virtual environment and install packages already in the cache, it hard-links files instead of copying them, making disk usage nearly zero for repeated installs.

**Fact 3.** `uv` can download and manage Python interpreter versions directly, similar to `pyenv`, without requiring a separate tool or shell shim. It stores them under `~/.local/share/uv/python/`.

---

## Tips and Tricks

### 1. Install uv and Understand What It Replaces

`uv` is distributed as a standalone binary with no Python dependency of its own.

```bash
# Install via the official installer script
curl -LsSf https://astral.sh/uv/install.sh | sh

# Or via pip if you already have Python
pip install uv

# Verify installation
uv --version

# What uv replaces:
# pip            -> uv pip install / uv add
# virtualenv     -> uv venv
# pyenv          -> uv python install / uv python pin
# pip-tools      -> uv pip compile / uv lock
# poetry init    -> uv init
```

After installation, `uv` is available as a single binary at `~/.cargo/bin/uv` or `~/.local/bin/uv` depending on the install method.

### 2. Replace pyenv for Python Version Management

Instead of installing `pyenv`, configuring shell shims, and running `pyenv install 3.12.3`, `uv` handles Python versions directly.

```bash
# List available Python versions
uv python list

# Install a specific Python version
uv python install 3.12.3
uv python install 3.11

# Pin a version for the current project directory
uv python pin 3.12.3
cat .python-version

# Run a one-off command with a specific Python version
uv run --python 3.11 python --version

# Show where uv stores Python interpreters
ls ~/.local/share/uv/python/
```

The `.python-version` file is compatible with `pyenv`, so existing tooling that reads it continues to work.

### 3. Initialize and Manage a Project

`uv init` creates a project with a `pyproject.toml` and a lockfile, similar to `cargo new` for Rust projects.

```bash
# Create a new project
uv init myproject
cd myproject
ls
# pyproject.toml  README.md  .python-version  src/myproject/__init__.py

# Add dependencies
uv add requests httpx pydantic

# Add a development-only dependency
uv add --dev pytest ruff

# View the generated lockfile
cat uv.lock | head -40

# Remove a dependency
uv remove httpx

# Install all dependencies from the lockfile (like npm ci)
uv sync
```

The `uv.lock` file records exact versions of every package in the dependency tree, including transitive dependencies. Committing it to version control gives reproducible installs across machines.

### 4. Replace pip in an Existing Project

For existing projects that use `requirements.txt` or `setup.py`, `uv` provides a `pip`-compatible interface.

```bash
# Create a virtual environment
uv venv .venv

# Activate it (same as before)
source .venv/bin/activate

# Install from requirements.txt
uv pip install -r requirements.txt

# Install in editable mode (replaces pip install -e .)
uv pip install -e .

# Compile a locked requirements file from requirements.in
uv pip compile requirements.in -o requirements.txt

# Sync environment to exactly match requirements.txt
uv pip sync requirements.txt

# Benchmark: time a fresh install of Django
time uv pip install django
# real    0m0.089s  (approximate, cache warm)
```

`uv pip sync` removes packages from the environment that are not in the requirements file, making the environment a precise match for the specification.

### 5. Run Scripts and Tools Without Installing

`uv run` executes a script or tool in an isolated environment without permanently installing anything. This is equivalent to `npx` in the Node.js ecosystem.

```bash
# Run a Python script with its dependencies declared inline (PEP 723)
cat > script.py << 'EOF'
# /// script
# dependencies = ["httpx", "rich"]
# ///
import httpx
from rich import print
r = httpx.get("https://httpbin.org/json")
print(r.json())
EOF

uv run script.py

# Run a tool (like ruff or black) without installing it permanently
uv tool run ruff check .
uv tool run black --check src/

# Install a tool globally (makes it available as a shell command)
uv tool install ruff
ruff --version

# Equivalent of pipx run
uvx black --check src/
```

`uvx` is a shorthand alias for `uv tool run`. It creates a temporary isolated environment, installs the requested tool, runs it, and then discards the environment. The package cache means subsequent runs of the same tool version are fast.
