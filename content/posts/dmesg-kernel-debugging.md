---
title: "Reading Kernel Messages with dmesg"
slug: e9d62f4b
aliases: dmesg-kernel-debugging
date: 2026-05-19
tags: [linux, debugging, tools]
lang: en
draft: false
---

The kernel's ring buffer records hardware events, driver messages, and subsystem errors from the moment the system starts. `dmesg` reads this buffer and is often the first place to look when hardware misbehaves or a driver fails silently.

## Fun Facts

**Fact 1.** The kernel ring buffer has a fixed size, configurable at compile time via `CONFIG_LOG_BUF_SHIFT`. On most distributions this is 18, meaning the buffer holds 2^18 = 256 KB of log data. On systems with many CPUs it is larger. Old messages are overwritten when the buffer fills.

**Fact 2.** Since kernel 3.5, `dmesg` can display timestamps in human-readable format with `-T`. Before that, all timestamps were seconds since boot (the `[12345.678901]` format), which required manual arithmetic to correlate with wall-clock time.

**Fact 3.** A kernel Oops is not a panic. An Oops records a fault and may allow the kernel to continue running, though the faulting process is usually killed and the system may become unstable afterward. A panic is a deliberate halt when the kernel determines it cannot recover.

---

## Tips and Tricks

### 1. Log Levels and Facility Codes

Each kernel message has a priority encoded as `<facility><level>`. The level follows syslog conventions:

| Level | Value | Meaning |
|---|---|---|
| KERN_EMERG | 0 | System is unusable |
| KERN_ALERT | 1 | Action must be taken immediately |
| KERN_CRIT | 2 | Critical conditions |
| KERN_ERR | 3 | Error conditions |
| KERN_WARNING | 4 | Warning conditions |
| KERN_NOTICE | 5 | Normal but significant |
| KERN_INFO | 6 | Informational |
| KERN_DEBUG | 7 | Debug-level messages |

```bash
# Show only errors and above (levels 0-3)
dmesg --level=emerg,alert,crit,err

# Show level alongside each message
dmesg -x

# Human-readable timestamps, show levels, color output
dmesg -Tx --color=always | less -R
```

### 2. Filtering by Time and Severity

`dmesg` supports time-based filtering since util-linux 2.30.

```bash
# Messages from the last 10 minutes
dmesg --since "10 minutes ago"

# Messages between two timestamps
dmesg --since "2026-05-19 08:00" --until "2026-05-19 09:00"

# Combine time filter with level filter
dmesg --since "1 hour ago" --level=err,warn

# Follow new messages in real time (like tail -f)
dmesg -w
```

For scripting, the raw timestamp (seconds since boot) is more reliable than the human-readable form:

```bash
# Get boot time in epoch seconds
BOOT_TIME=$(date -d "$(uptime -s)" +%s)

# Convert a raw dmesg timestamp to wall time
RAW_TS=12345.678901
echo "$(date -d @$(echo "$BOOT_TIME + $RAW_TS" | bc) +'%Y-%m-%d %H:%M:%S')"
```

### 3. Common Hardware Errors and What They Mean

**ACPI errors** appear during boot when the firmware provides malformed tables:

```
ACPI Error: AE_NOT_FOUND, While evaluating Sleep State [\_S1_]
ACPI BIOS Error (bug): Could not resolve symbol [\_SB.PCI0.XHC.RHUB]
```

These are usually harmless unless accompanied by a device that fails to initialize. Set `acpi=off` in kernel parameters only as a last resort.

**EDAC memory errors** report correctable (CE) and uncorrectable (UE) ECC RAM faults:

```
EDAC MC0: 1 CE memory read error on CPU_SrcID#0_Ha#0_Chan#0_DIMM#0
EDAC MC0: 1 UE memory read error on CPU_SrcID#0_Ha#0_Chan#1_DIMM#1
```

Correctable errors are recoverable. Uncorrectable errors often trigger a machine check exception and system crash. Track CE counts over time with `edac-util -s 0`.

**NVMe timeouts** indicate a drive that is not responding within the expected window:

```bash
# Filter NVMe errors
dmesg | grep -i nvme | grep -iE "error|timeout|reset|failed"

# Example output:
# nvme nvme0: I/O 23 QID 1 timeout, reset controller
# nvme nvme0: controller is down; will reset: CSTS=0x3, PCI_STATUS=0x10
```

A single timeout followed by a successful reset is often benign. Repeated resets indicate failing hardware or a firmware bug.

### 4. Decoding a Kernel Oops Call Trace

An Oops contains the faulting instruction, register state, and a call trace. The trace uses kernel symbol names and offsets.

```
BUG: kernel NULL pointer dereference, address: 0000000000000010
RIP: 0010:some_driver_function+0x42/0x180 [mymodule]
Call Trace:
 <TASK>
 other_function+0x1a/0x60 [mymodule]
 __do_something+0x88/0x100
 process_one_work+0x1f5/0x390
 worker_thread+0x4d/0x3d0
```

To decode the `+0x42/0x180` offset back to a source line, use `addr2line` against the kernel module:

```bash
# Find the module file
find /lib/modules/$(uname -r) -name 'mymodule.ko*'

# Decode the offset (strip KASLR by using the relative offset only)
addr2line -e /lib/modules/$(uname -r)/kernel/drivers/misc/mymodule.ko \
  -f -i 0x42

# Alternatively, use the kernel's own script
scripts/decode_stacktrace.sh /usr/lib/debug/lib/modules/$(uname -r)/vmlinux \
  < oops.txt
```

### 5. journalctl -k as a More Flexible Alternative

`journalctl -k` reads the same kernel messages but through the systemd journal, which persists across reboots and supports richer filtering.

```bash
# Current boot kernel messages, errors and above only
journalctl -k -p err

# Kernel messages from the previous boot
journalctl -k -b -1

# List available boots
journalctl --list-boots

# Kernel messages matching a pattern, with context
journalctl -k --grep="nvme" -n 50

# Export to plain text for sharing
journalctl -k -b 0 --no-pager > kernel-boot.log
```

One practical difference: `dmesg` shows the live ring buffer and is always available, even without systemd. `journalctl -k` gives you historical data across reboots and integrates with journal filtering, which makes it more useful for post-incident analysis on systems that use systemd.
