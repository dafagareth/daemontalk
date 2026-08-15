---
title: "Monitoring Resource Starvation with Linux PSI and cgroups v2"
slug: e9a1b3c8
aliases: [linux-cgroups-v2-psi]
date: 2026-05-08
tags: [linux, devops, performance]
lang: en
draft: false
type: post
---

Traditional Linux load average metrics combine active CPU execution with uninterruptible I/O waits, obscuring the precise root cause of system slowdowns. Pressure Stall Information (PSI) introduced in kernel 4.20 quantifies exact hardware resource starvation across CPU, memory, and I/O subsystems. In cgroups v2, PSI metrics enable fine-grained monitoring and proactive memory reclaiming before kernel Out-Of-Memory (OOM) panics occur.

## Fun Facts

**Fact 1.** PSI distinguishes between `some` pressure (at least one task is stalled waiting for resources) and `full` pressure (all non-idle tasks are stalled simultaneously, representing complete throughput loss).

**Fact 2.** The Linux kernel uses a windowed moving average calculation across 10-second, 60-second, and 300-second intervals to expose PSI metrics via pseudofiles in `/proc/pressure/`.

**Fact 3.** Android was an early adopter of PSI, using vendor daemon triggers on memory stall thresholds to terminate background activities before memory pressure freezes the display compositor.

---

## Tips and Tricks

### 1. Inspect System-Wide Pressure Stall Metrics
Read `/proc/pressure/cpu`, `/proc/pressure/memory`, and `/proc/pressure/io` to view real-time starvation percentages.

```bash
cat /proc/pressure/memory
# Output:
# some avg10=0.00 avg60=0.00 avg300=0.00 total=0
# full avg10=0.00 avg60=0.00 avg300=0.00 total=0
```

### 2. Read Cgroup-Specific Pressure in cgroups v2
Navigate to a target cgroup hierarchy to inspect localized resource starvation on specific container workloads.

```bash
cat /sys/fs/cgroup/system.slice/memory.pressure
cat /sys/fs/cgroup/system.slice/io.pressure
```

### 3. Register Epoll Events on Memory Stall Thresholds
Write a lightweight monitoring loop that uses `epoll` or `poll` to trigger immediate alerts when 10-second memory stalls exceed a 20 percent threshold.

```bash
# Register a custom trigger on memory pressure pseudofile
echo "some 200000 1000000" > /proc/pressure/memory
```

### 4. Enable systemd-oomd with PSI Thresholds
Configure `systemd-oomd` to kill memory-hogging cgroups when memory stall percentages exceed user-defined thresholds.

```ini
# /etc/systemd/oomd.conf
[OOM]
SwapUsedLimit=80%
DefaultMemoryPressureLimit=60%
DefaultMemoryPressureDurationSec=30s
```

### 5. Enforce PSI Memory Relief on Specific Systemd Services
Apply `ManagedOOMMemoryPressure` settings directly inside unit file configurations to isolate production microservices from uncontrolled swapping.

```ini
[Unit]
Description=Backend API Microservice

[Service]
ExecStart=/usr/local/bin/api-server
MemoryAccounting=true
CPUAccounting=true
IOAccounting=true
ManagedOOMMemoryPressure=kill
ManagedOOMMemoryPressureLimit=50%
```

### 6. Monitor PSI Metrics Prometheus Exporter
Scrape cgroup v2 PSI metrics into Prometheus using Node Exporter to generate proactive alerting graphs on cluster resource degradation.

```bash
node_exporter --collector.cgroups --collector.powersupply-class.disabled
```
