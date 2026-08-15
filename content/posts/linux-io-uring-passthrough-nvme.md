---
title: "Akses NVMe Langsung dengan io_uring Passthrough di Linux"
slug: d5a1f4e9
aliases: [linux-io-uring-passthrough-nvme]
date: 2026-07-25
tags: [linux, storage, performance]
lang: id
draft: false
type: post
---

Operasi I/O pada SSD NVMe modern sering terhambat oleh overhead perangkat lunak di dalam kernel Linux. Dengan fitur `io_uring` passthrough (`io_uring_cmd`), aplikasi userspace dapat mengirimkan perintah NVMe secara langsung ke pengontrol drive tanpa melewati antarmuka block layer konvensional.

## Fun Fact

**Fact 1.** Fitur io_uring passthrough pertama kali diperkenalkan pada Linux Kernel 5.19 untuk memungkinkan eksekusi perintah NVMe generik dari pengguna.

**Fact 2.** Dalam pengarsitekturan I/O tradisional, setiap request harus melintasi lapisan Virtual File System (VFS), I/O Scheduler, dan Block Layer sebelum sampai ke driver NVMe.

**Fact 3.** Pengujian throughput menggunakan `fio` engine `io_uring_cmd` mencatatkan penurunan latensi p99 hingga 20% dibandingkan metode pembacaan blok `/dev/nvmeXnY` standar.

---

## Tips dan Trik

### 1. Arsitektur io_uring Passthrough (io_uring_cmd)

Arsitektur I/O standar kernel Linux mengonversi setiap request menjadi struktur `struct bio` dan `struct request`. Proses pengubahan ini mengonsumsi siklus CPU yang signifikan ketika drive NVMe mampu memproses jutaan IOPS.

`io_uring_cmd` menyediakan jalur pintas di mana aplikasi dapat membungkus perintah NVMe 64-byte langsung ke dalam Submission Queue Entry (SQE) `io_uring`, yang kemudian ditangani secara langsung oleh driver `/dev/ngXnY`.

### 2. Memeriksa Perangkat Character Device NVMe

Gunakan pustaka `nvme-cli` untuk mengidentifikasi antarmuka character device `/dev/ngXnY` pada sistem Linux:

```bash
# Menampilkan daftar perangkat NVMe generic character device
ls -l /dev/ng*

# Mengatur hak akses perangkat tanpa hak akses root
sudo chmod 666 /dev/ng0n1
```

### 3. Mengirim Perintah NVMe Read dari Userspace Menggunakan C

Berikut adalah pustaka C sederhana menggunakan `liburing` untuk menyiapkan instruksi NVMe Read langsung ke character device:

```c
#include <liburing.h>
#include <linux/nvme_ioctl.h>

void prepare_nvme_read(struct io_uring_sqe *sqe, int fd, void *buf, uint64_t slba, uint16_t nlb) {
    struct nvme_uring_cmd *cmd;
    
    io_uring_prep_cmd(sqe, IORING_OP_URIS_CMD, fd);
    cmd = (struct nvme_uring_cmd *)sqe->cmd;
    cmd->opcode = 0x02; /* NVMe Read Opcode */
    cmd->addr = (__u64)buf;
    cmd->data_len = nlb * 512;
    cmd->slba = slba;
    cmd->nlb = nlb - 1;
}
```

### 4. Benchmark Throughput IOPS Menggunakan fio Engine io_uring_cmd

Buat file konfigurasi pengujian `/etc/fio/nvme-passthrough.fio` untuk mengukur performa IOPS maksimum:

```ini
[global]
filename=/dev/ng0n1
ioengine=io_uring_cmd
cmd_type=nvme
iodepth=64
numjobs=4
thread=1
group_reporting=1

[nvme-read-test]
rw=randread
bs=4k
time_based=1
runtime=30s
```

### 5. Menjalankan Pengujian dan Mengevaluasi Latensi

Eksekusi pengujian fio untuk membandingkan throughput dan latensi:

```bash
# Menjalankan pengujian performa passthrough
fio /etc/fio/nvme-passthrough.fio
```

Hasil benchmark menunjukkan pengurangan overhead alokasi memori kernel dan peningkatan throughput IOPS hingga 15-20% pada beban kerja acak 4KB.
