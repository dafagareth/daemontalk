---
title: "Cgroups & OOM Killer: Mengapa Proses Kamu Tiba-tiba Mati?"
slug: 7c8d9e0f
aliases: [linux-cgroups-memory-limits]
date: 2026-05-24
tags: [linux, sysadmin, architecture]
lang: id
draft: false
---

Pernahkah kamu mendapati proses web server, build Go, atau database di VPS tiba-tiba menghilang tanpa pesan error di aplikasi? Kemungkinan besar proses tersebut dieksekusi mati oleh **Linux Kernel OOM (Out-of-Memory) Killer** melalui batasan **cgroups**.

Memahami cara kerja cgroups dan OOM score adalah kunci menjaga uptime server production.

## Fun Fact

**Cgroups (Control Groups) awalnya dirancang oleh engineer Google.**
Rohit Seth dan Paul Menage memulainya pada 2006 dengan nama "Process Containers", sebelum akhirnya diubah namanya menjadi cgroups dan dimerge ke Linux Kernel 2.6.24. Cgroups adalah fondasi utama kelahiran Docker dan LXC!

**Linux Kernel secara default melakukan memory overcommit.**
Kernel Linux berjanji memberikan alokasi memori ke proses melebihi kapasitas RAM fisik yang sebenarnya tersedia (mirip maskapai penerbangan yang overbooking tiket).

**OOM Killer memiliki sistem penilaian skor (*badness score*).**
Kernel menghitung proses mana yang paling layak dibunuh berdasarkan jumlah RAM yang digunakan dan nilai `oom_score_adj`.

---

## Tips dan Trik

### 1. Periksa Log Kernel untuk Bukti Pembunuhan OOM

Jika suatu proses mendadak exit code 137 (`128 + 9 (SIGKILL)`), cek `dmesg`:

```bash
sudo dmesg -T | grep -i -E "oom|out of memory|killed process"
```

### 2. Lindungi Proses Kritis dari OOM Killer dengan `oom_score_adj`

Beri skor proteksi negatif (antara -1000 hingga 1000) pada proses seperti SSH daemon atau Database agar tidak menjadi target utama:

```bash
# Lindungi proses PID 1234 agar tidak disentuh OOM Killer
echo -1000 | sudo tee /proc/1234/oom_score_adj
```

### 3. Batasi Memory Service via Systemd Unit

Manfaatkan cgroups v2 langsung dari file `.service` tanpa perlu tool tambahan:

```ini
[Service]
MemoryMax=1G
MemoryHigh=800M
CPUQuota=50%
```

### 4. Konfigurasi `vm.overcommit_memory` di Sysctl

Untuk database performa tinggi seperti Redis atau PostgreSQL, atur overcommit handling sesuai anjuran resmi:

```ini
# /etc/sysctl.d/99-memory.conf
vm.overcommit_memory = 1
vm.swappiness = 10
```

### 5. Pantau cgroups Resource Usage dengan `systemd-cgtop`

Lihat penggunaan CPU, Memory, dan Disk IO per control group secara real-time:

```bash
systemd-cgtop -m
```
