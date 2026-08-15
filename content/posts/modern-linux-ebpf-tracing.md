---
title: "eBPF Tracing for Modern Linux Infrastructure"
slug: e7b91024
aliases: [modern-linux-ebpf-tracing]
date: 2025-11-14
tags: [linux, devops, tools]
lang: en
draft: false
type: post
---

Extended Berkeley Packet Filter (eBPF) technology provides programmable runtime introspection within the Linux kernel without requiring external kernel modules. Modern Linux kernels execute verified eBPF bytecode through a Just-In-Time (JIT) compiler, allowing engineers to trace kernel execution paths and monitor system calls with minimal overhead.

## Fun Facts

**Fact 1.** The Linux kernel verifier evaluates program control flow up to 1,000,000 instructions to ensure safety, preventing invalid pointer dereferences and infinite loops before JIT execution.

**Fact 2.** BPF Type Format (BTF) provides compact type information embedded directly in the kernel image, enabling Compile Once: Run Everywhere (CO-RE) portability across different kernel builds.

**Fact 3.** Modern Linux distributions automatically mount the dedicated eBPF file system at `/sys/fs/bpf`, allowing userspace utilities to pin BPF maps and retain persistent tracing data across process restarts.

---

## Tips and Tricks

### 1. Query Kernel Tracepoints and Filter Functions

Before attaching eBPF probes, inspect available tracepoints and kernel functions directly via the tracefs interface:

```bash
cat /sys/kernel/tracing/available_filter_functions | grep -E '^vfs_' | head -n 10
```

You can also query tracepoints using `bpftool` to verify kernel BTF support:

```bash
sudo bpftool feature probe kernel
```

### 2. Trace VFS Read Latency with bpftrace One-Liners

Measure Virtual File System (VFS) read duration by calculating elapsed nanoseconds between enter and exit kprobes:

```bash
sudo bpftrace -e '
kprobe:vfs_read { @start[tid] = nsecs; }
kretprobe:vfs_read /@start[tid]/ {
  @dur_ns = hist(nsecs - @start[tid]);
  delete(@start[tid]);
}
'
```

This outputs a power-of-two histogram of read call latencies when interrupted with `Ctrl-C`.

### 3. Pin eBPF Maps to bpffs for Persistent Metrics

Create a hash map directly from the terminal using `bpftool` and pin it to `/sys/fs/bpf/` for out-of-process access:

```bash
sudo bpftool map create /sys/fs/bpf/my_map type hash key 4 value 8 entries 1024 name my_map
sudo bpftool map show pinned /sys/fs/bpf/my_map
```

Inspect map contents at any time using:

```bash
sudo bpftool map dump pinned /sys/fs/bpf/my_map
```

### 4. Export Kernel Data Types into C Headers via BTF

Extract exact kernel data structure definitions directly from running kernel BTF metadata into a header file for custom eBPF development:

```bash
sudo bpftool btf dump file /sys/kernel/btf/vmlinux format c > vmlinux.h
head -n 25 vmlinux.h
```

Including `vmlinux.h` in your C source code eliminates external kernel header dependencies.

### 5. Profile Process Off-CPU Time to Find Blocking I/O

Track time spent in sleep or blocked state per process using `offcputime` from the BPF Compiler Collection:

```bash
sudo /usr/share/bcc/tools/offcputime -p $(pgrep -n nginx) 5
```

This captures kernel stack traces whenever the target process relinquishes CPU control, pinpointing lock contention or disk I/O bottlenecks.

### 6. Inspect Verifier Logs for Failed Bytecode Loading

When loading compiled object files fails, view detailed verifier execution state traces using `bpftool`:

```bash
sudo bpftool prog load trace_prog.o /sys/fs/bpf/trace_prog type kprobe log_level 1
```

Check kernel diagnostic messages with `dmesg`:

```bash
sudo dmesg -T | grep -i bpf | tail -n 20
```
