---
title: "Linux Pipes and Redirection: The Parts Most People Skip"
slug: 861b26db
aliases: [linux-pipes-redirection]
date: 2025-08-14
tags: [linux, bash, cli]
lang: en
draft: false
---

Most people learn `|` and `>` early and stop there. That covers 80% of daily use. But the remaining 20% shows up often enough that not knowing it means reaching for awkward workarounds.

This covers the parts that come up in real work.

## stderr Is Not stdout

Every process has two output streams: stdout (fd 1) and stderr (fd 2). The pipe `|` only connects stdout. If a command prints errors, they go to your terminal regardless of what you pipe.

Redirect stderr to stdout to capture both:

```bash
command 2>&1 | grep "error"
```

Redirect each to different files:

```bash
command > output.log 2> error.log
```

Discard errors entirely:

```bash
command 2>/dev/null
```

## Appending vs. Overwriting

```bash
command > file.txt    # overwrites
command >> file.txt   # appends
```

A common mistake is using `>` in a loop and wondering why only the last iteration is in the file.

## tee: Write and Pass Through

`tee` writes to a file and passes stdin to stdout simultaneously. Useful when you want to log output while still seeing it in the terminal:

```bash
./build.sh | tee build.log
```

Append instead of overwrite:

```bash
./deploy.sh | tee -a deploy.log
```

## Process Substitution

Process substitution lets you use a command's output as if it were a file. The syntax `<(command)` creates a temporary file descriptor.

Compare two command outputs directly:

```bash
diff <(sort file1.txt) <(sort file2.txt)
```

Without process substitution, you would need to create two temporary files, run diff, then delete them.

Feed multiple outputs into a command that expects files:

```bash
paste <(cut -d, -f1 data.csv) <(cut -d, -f3 data.csv)
```

## Here Documents

A here-doc passes multiline text to a command without creating a temporary file:

```bash
cat <<EOF
line one
line two
line three
EOF
```

Common use: writing config files in scripts:

```bash
cat > /etc/app/config.toml <<EOF
[server]
port = 8080
host = "0.0.0.0"
EOF
```

Use `<<'EOF'` (quoted) to prevent variable expansion inside the block:

```bash
cat <<'EOF'
This $variable will not be expanded
EOF
```

## Named Pipes (FIFOs)

A named pipe is a file that connects two processes. One process writes to it, another reads from it, and no data hits disk.

```bash
mkfifo /tmp/pipe
command_a > /tmp/pipe &
command_b < /tmp/pipe
```

Less common than anonymous pipes, but useful when you cannot chain commands directly (for example, two commands in different terminals or scripts).

## Subshell Output with $()

Command substitution captures stdout into a variable:

```bash
files=$(find . -name "*.go" -newer go.mod)
echo "Modified since last go.mod change: $files"
```

This is cleaner than backticks and can be nested:

```bash
echo "Go version: $(go version | cut -d' ' -f3)"
```

These primitives compose well. Understanding them lets you build pipelines that would otherwise require a temporary script or an external tool.
