---
title: "Pengembangan Custom CPU Scheduler Linux Menggunakan sched_ext dan eBPF"
slug: 3b9c8d2a
aliases: [kernel-sched-ext-custom-scheduler]
date: 2026-07-05
tags: [linux, performance, ebpf]
lang: id
draft: false
type: post
---

Fitur `sched_ext` (Extensible Scheduler Class) pada Linux 6.12 memungkinkan pengembang mengimplementasikan penjadwalan CPU kustom di tingkat kernel menggunakan eBPF. Mekanisme ini menggantikan atau melengkapi scheduler standar kernel tanpa perlu mengompilasi ulang kernel secara keseluruhan. Tulisan ini menjelaskan keterbatasan scheduler EEVDF dalam beban kerja latensi rendah, arsitektur kfunc `sched_ext`, struktur dispatch queue, dan contoh program eBPF scheduler.

## Fakta Menarik

**Fakta 1.** `sched_ext` resmi didukung dan masuk ke dalam mainline Linux Kernel versi 6.12 setelah dipelopori oleh insinyur Meta dan Canonical untuk beban kerja khusus.

**Fakta 2.** Jika program eBPF scheduler mengalami kecelakaan (crash) atau deadlock, kernel Linux secara otomatis mengembalikan penjadwalan CPU ke scheduler standar EEVDF tanpa menyebabkan kernel panic.

**Fakta 3.** Proyek sched-ext menyediakan berbagai scheduler siap pakai di user-space seperti `scx_bpfland` dan `scx_lavd` yang dirancang khusus untuk meningkatkan frame rate pada aplikasi game.

---

## Tips dan Trik

### 1. Pahami Keterbatasan Scheduler EEVDF untuk Latensi Rendah

Scheduler standar EEVDF (Earliest Eligible Virtual Deadline First) dirancang untuk pembagian wajar (fairness) sumber daya CPU pada server serbaguna. Namun, pada aplikasi gaming atau sistem audio real-time, penundaan akibat context switching dapat menyebabkan jitter yang tidak diinginkan.

```bash
# Periksa ketersediaan sched_ext pada sistem Linux 6.12+
zgrep CONFIG_SCHED_CLASS_EXT /proc/config.gz

# Periksa status sysfs scheduler yang sedang aktif
cat /sys/kernel/sched_ext/state
```

### 2. Pahami Fungsi Callback dan kfunc Utama pada sched_ext

Program eBPF `sched_ext` berinteraksi dengan kernel melalui callback seperti `select_cpu`, `enqueue`, dan `dispatch`, serta fungsi penolong (kfunc) seperti `scx_bpf_dispatch`.

```c
/* Struktur callback dasar sched_ext dalam eBPF */
SEC("struct_ops/scx_simple_select_cpu")
s32 BPF_PROG(simple_select_cpu, struct task_struct *p, s32 prev_cpu, u64 wake_flags)
{
    /* Kembalikan CPU sebelumnya untuk menjaga kestabilan cache L1/L2 */
    return prev_cpu;
}
```

### 3. Masukkan Task ke Dispatch Queue Menggunakan scx_bpf_dispatch

Proses dispatch bertugas memindahkan tugas yang siap jalan ke antrean eksekusi CPU (Local Dispatch Queue atau Shared Dispatch Queue).

```c
SEC("struct_ops/scx_simple_enqueue")
void BPF_PROG(simple_enqueue, struct task_struct *p, u64 enq_flags)
{
    /* Masukkan tugas langsung ke Shared Dispatch Queue dengan slice waktu 5ms */
    scx_bpf_dispatch(p, SCX_DSQ_GLOBAL, 5000000, enq_flags);
}
```

### 4. Buat Implementasi Kernel Struct Ops Sederhana

Gabungkan callback ke dalam struktur `sched_ext_ops` agar dapat dimuat oleh verifier eBPF ke dalam kernel space.

```c
SEC(".struct_ops")
struct sched_ext_ops simple_ops = {
    .select_cpu = (void *)simple_select_cpu,
    .enqueue    = (void *)simple_enqueue,
    .name       = "scx_simple",
};
```

### 5. Jalankan dan Pantau Scheduler dari User Space

Gunakan loader berbasis libbpf di user space untuk memuat program eBPF dan memantau statistik penjadwalan CPU.

```bash
# Kompilasi kode sumber eBPF scheduler
clang -O2 -target bpf -c scx_simple.bpf.c -o scx_simple.bpf.o

# Registrasikan objek struct_ops ke kernel menggunakan bpftool
sudo bpftool struct_ops register scx_simple.bpf.o

# Verifikasi status registrasi scheduler di kernel
sudo bpftool struct_ops list
```
