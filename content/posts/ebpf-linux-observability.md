---
title: "eBPF for Linux Observability and Tracing"
slug: a3f7e1c9
aliases: [ebpf-linux-observability]
date: 2026-05-28
tags: [linux, performance, debugging]
lang: en
draft: false
---

eBPF (extended Berkeley Packet Filter) lets you run sandboxed programs inside the Linux kernel without writing a kernel module or rebooting. It has become the standard substrate for production-grade observability, networking, and security tooling on modern Linux systems.

## Fun Facts

**Fact 1.** The original BPF was introduced in 1992 for packet filtering in `tcpdump`. The "extended" variant, merged in Linux 3.18 (2014), broadened the instruction set and added a general-purpose map abstraction, making it useful far beyond networking.

**Fact 2.** Unlike kernel modules, eBPF programs cannot crash the kernel. The in-kernel verifier statically analyzes every program before it runs, rejecting any code with unbounded loops, out-of-bounds memory access, or uninitialized reads.

**Fact 3.** Cilium, the CNCF-graduated CNI plugin used by Google Kubernetes Engine and many others, replaces iptables entirely with eBPF programs, achieving lower latency and better scalability on large clusters.

---

## Tips and Tricks

### 1. Understand eBPF vs. Kernel Modules

A kernel module runs at ring 0 with no safety net: a null pointer dereference panics the system. An eBPF program is loaded via `bpf(2)`, verified, JIT-compiled, and attached to a hook point such as a kprobe, tracepoint, or network interface. If verification fails, the load syscall returns an error and nothing runs.

The trade-off is expressiveness: eBPF programs cannot call arbitrary kernel functions (only a curated set of helpers), and loops must terminate in a provably finite number of iterations. For most observability tasks this is not a real constraint.

### 2. Install bpftrace for One-Liner Tracing

`bpftrace` provides an awk-like language for writing eBPF probes. On Debian/Ubuntu:

```bash
sudo apt install bpftrace
```

On Fedora/RHEL:

```bash
sudo dnf install bpftrace
```

Kernel headers and `debugfs` must be mounted. Verify:

```bash
mount | grep debugfs
# /sys/kernel/debug type debugfs (rw,nosuid,nodev,noexec,relatime)
```

### 3. Trace Syscall Latency with a bpftrace One-Liner

The following program measures the time each `read(2)` call spends in the kernel and prints a histogram on exit:

```bash
sudo bpftrace -e '
tracepoint:syscalls:sys_enter_read { @start[tid] = nsecs; }
tracepoint:syscalls:sys_exit_read  /@start[tid]/ {
  @latency_ns = hist(nsecs - @start[tid]);
  delete(@start[tid]);
}
END { print(@latency_ns); }
'
```

Send `SIGINT` (Ctrl-C) to stop. The histogram buckets show where latency is concentrated, distinguishing fast cache hits from slow disk reads.

### 4. Use bcc Tools for Deeper Analysis

The BCC (BPF Compiler Collection) project ships over 100 ready-made tools. A few practical ones:

```bash
# Show I/O latency distribution per disk, updated every second
sudo /usr/share/bcc/tools/biolatency -d sda 1

# Trace TCP connections being established
sudo /usr/share/bcc/tools/tcpconnect

# Profile CPU usage by kernel and user stack at 99 Hz, capture for 30 s
sudo /usr/share/bcc/tools/profile -F 99 -f 30 > flamegraph.folded
```

Install on Debian/Ubuntu with `sudo apt install bpfcc-tools`. The tools are Python scripts that compile embedded C eBPF programs at runtime using LLVM.

### 5. Attach to User-Space Probes with USDT

eBPF is not limited to the kernel. User Statically Defined Tracing (USDT) probes allow attaching to instrumented points in user-space binaries. The Go runtime, CPython, PostgreSQL, and many others ship with USDT probes.

Trace PostgreSQL query execution:

```bash
sudo bpftrace -e '
usdt:/usr/lib/postgresql/16/bin/postgres:postgresql:query__start {
  printf("%s\n", str(arg0));
}
'
```

This prints every SQL statement at the moment it begins executing, with zero modifications to the application binary.

### 6. Read the Verifier Log When a Load Fails

When `bpftool prog load` or a library fails with `EACCES` or `EINVAL`, the verifier has rejected the program. Read the rejection reason:

```bash
sudo bpftool prog load my_prog.o /sys/fs/bpf/my_prog 2>&1 | head -40
```

Common rejection reasons include accessing a map value without checking for a NULL pointer first, and calling a helper that is not permitted for the chosen program type (e.g., `bpf_probe_write_user` is restricted to `BPF_PROG_TYPE_KPROBE`).
