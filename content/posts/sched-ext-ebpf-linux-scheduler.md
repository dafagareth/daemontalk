---
title: "sched_ext Mendarat di Linux: Menulis Scheduler CPU Sendiri Pakai eBPF"
slug: "sched-ext-ebpf-linux-scheduler"
aliases: []
date: 2026-09-28
author: "daemontalk team"
contributors: []
tags: ["radar", "linux", "kernel", "ebpf"]
lang: "id"
status: "published"
type: "dispatch"
readTime: 5
cover: "https://images.unsplash.com/photo-1518770660439-4636190af475?auto=format&fit=crop&w=1200&q=80"
coverCaption: "Pengalihan penjadwalan proses CPU ke program eBPF dinamis tanpa kompilasi ulang kernel"
coverSource: "https://unsplash.com"
description: "Fitur sched_ext resmi digabung ke kernel Linux utama. Pengembang kini dapat mengganti algoritma penjadwalan CPU default dengan program eBPF dari userspace secara aman dan dinamis."
series: ""
series_part: 0
---

<span class="drop-cap">S</span>elama tiga dekade, penjadwal CPU (*CPU scheduler*) adalah wilayah paling sakral dan paling berisiko untuk disentuh di dalam kernel Linux.

Menulis atau menguji algoritma penjadwalan baru biasanya berarti Anda harus mengedit kode C di dalam pohon direktori `kernel/sched/`, mengompilasi ulang kernel secara utuh, melakukan *reboot* mesin, dan bersiap menghadapi *kernel panic* jika ada *lock* yang macet atau logika kuantum waktu yang keliru di jalur kritis eksekusi prosesor.

Dari era O(1) rancangan Ingo Molnár, Completely Fair Scheduler (CFS) yang bertahan hampir dua dekade, hingga EEVDF (Earliest Eligible Virtual Deadline First) yang diperkenalkan pada Linux 6.6, paradigma penjadwalan Linux selalu bersifat monolitik: satu algoritma umum yang dipaksa melayani semua jenis beban kerja sekaligus.

Penggabungan resmi **sched_ext** (*Sched Extensible*) ke dalam pohon kernel Linux utama mengubah batasan sejarah tersebut.

## Dilema Kompromi "Satu Ukuran untuk Semua"

Sebuah klaster komputasi komputasi saintifik (*batch processing*) menginginkan proses berjalan selama mungkin pada satu core CPU untuk memaksimalkan *cache locality* instruksi L2/L3. Sebaliknya, server web mikro I/O-bound atau basis data transaksional membutuhkan pergantian thread instan demi menekan latensi *tail* p99.

Di lingkungan desktop dan konsol portabel seperti Steam Deck, kebutuhan ini makin kontras: thread antarmuka grafis dan render game harus diprioritaskan di atas thread kompilasi shader latar belakang agar tidak memicu *stutter* visual.

Karena penjadwal kernel standar harus menjadi penengah bagi semua skenario tersebut, keputusannya selalu merupakan kompromi: cukup baik untuk semua orang, tetapi tidak pernah optimal untuk beban kerja spesifik.

Perusahaan teknologi skala besar selama bertahun-tahun terpaksa memelihara pohon kernel internal sendiri hanya untuk memodifikasi parameter penjadwalan, sebuah beban pemeliharaan (*maintenance burden*) yang luar biasa mahal saat kernel harus di-upgrade.

## Arsitektur sched_ext: Menaruh eBPF di Jalur Kritis

`sched_ext` memperkenalkan kelas penjadwal baru bernama `ext_sched_class`. Di dalam hierarki prioritas kernel Linux, kelas ini disematkan tepat di antara kelas Real-Time (`rt_sched_class`) dan kelas Completely Fair / EEVDF (`fair_sched_class`).

```text
+---------------------------------------------+
|             stop_sched_class                | (Prioritas Tertinggi)
+---------------------------------------------+
                       |
+---------------------------------------------+
|              dl_sched_class                 | (SCHED_DEADLINE)
+---------------------------------------------+
                       |
+---------------------------------------------+
|              rt_sched_class                 | (SCHED_FIFO / SCHED_RR)
+---------------------------------------------+
                       |
+---------------------------------------------+
|           ext_sched_class (sched_ext)       | <--- eBPF Dynamic Scheduler
+---------------------------------------------+
                       |
+---------------------------------------------+
|          fair_sched_class (EEVDF / CFS)     | (Fallback Otomatis)
+---------------------------------------------+
                       |
+---------------------------------------------+
|             idle_sched_class                | (Prioritas Terendah)
+---------------------------------------------+
```

Alih-alih mengunci logika di kode mesin kernel statis, `ext_sched_class` mengekspos serangkaian kait (*callbacks*) operasi yang dapat diimplementasikan oleh program eBPF:

- `ops.select_cpu`: Memilih core CPU target saat sebuah task terbangun, memperhitungkan topologi NUMA atau kondisi *idle*.
- `ops.enqueue`: Menempatkan task yang siap jalan ke dalam antrean eksekusi (Dispatch Queue / DSQ).
- `ops.dispatch`: Mengambil task dari DSQ dan menyerahkannya ke core CPU lokal yang siap mengeksekusi.

Program pengguna (*user space daemon*) cukup memuat program eBPF tersebut ke dalam kernel melalui system call `bpf()`. Dalam hitungan milidetik, kernel mengalihkan seluruh proses normal ke penjadwal eBPF Anda tanpa me-restart satu pun aplikasi yang sedang aktif.

## Garansi Keamanan: Watchdog dan Fallback Otomatis

Kekhawatiran pertama para insinyur keandalan sistem (*SRE*) saat mendengar konsep ini adalah stabilitas: *bagaimana jika kode eBPF scheduler kita mengalami bug, deadlock logika, atau starvation?*

Apakah seluruh server produksi akan hang seketika?

Jawabannya adalah tidak. Arsitektur `sched_ext` dirancang dengan jaring pengaman berlapis:

1. **eBPF Verifier**: Menolak program jika terdapat dereferensi pointer memori liar atau potensi pembacaan di luar batas memori yang diizinkan.
2. **Built-in Watchdog**: Kernel memonitor kemajuan eksekusi seluruh thread secara berkala. Jika ada task yang tidak kunjung dijadwalkan melampaui ambang batas waktu tertentu (*starvation*), watchdog kernel akan otomatis memutuskan sambungan (*detach*) program eBPF tersebut.
3. **Seamless Fallback**: Detik itu juga, seluruh task dipindahkan kembali ke `fair_sched_class` (EEVDF) bawaan kernel. Sistem tetap beroperasi normal tanpa *kernel panic*, dan pesan diagnostik kerusakan tersimpan rapi di buffer `dmesg`.

## Implementasi Nyata: Dari Meta hingga Handheld Gaming

`sched_ext` bukan sekadar proyek eksperimental di laboratorium akademis. Implementasinya sudah membuktikan efisiensi di berbagai medan produksi:

- **Meta**: Mengembangkan varian scheduler seperti `scx_rusty` dan `scx_layered` untuk klaster server mereka. Dengan memprioritaskan paket web berlatensi kritis di atas proses background batch, mereka mencatat penurunan variansi latensi p99 secara signifikan.
- **Valve & Komunitas Gaming Linux**: Memanfaatkan scheduler seperti `scx_lavd` (*Latency-criticality Aware Virtual Deadline*) pada perangkat Linux gaming. Hasilnya, responsivitas input dan stabilitas frame rate meningkat drastis bahkan saat proses berat berjalan simultan di belakang layar.

## Mencoba Langsung di Terminal

Bagi pengguna Linux modern dengan kernel yang mendukung `CONFIG_SCHED_CLASS_EXT=y`, mencoba penjadwal ini tidak membutuhkan penulisan ribuan baris kode C dari nol. Toolchain komunitas `scx` menyediakan serangkaian binary siap pakai:

```bash
# Menjalankan scheduler scx_rusty (implementasi penjadwal berbasis Rust)
sudo scx_rusty

# Di terminal lain, pantau distribusi task pada tiap core CPU
scx_stats
```

Saat Anda menekan `Ctrl+C` pada proses `scx_rusty`, daemon akan meng-unhook program eBPF dan kernel Linux secara mulus kembali menjalankan EEVDF tanpa kedipan sedikit pun pada sistem.

## Era Baru Penjadwalan Sistem

Hadirnya `sched_ext` menandai pergeseran fundamental dalam rekayasa sistem operasi. Selama puluhan tahun, kernel diposisikan sebagai kotak hitam yang hanya bisa dikonfigurasi melalui segelintir kenop `sysctl`.

Dengan eBPF yang kini merambah hingga lapisan penjadwalan CPU paling inti, dinding antara kebutuhan spesifik aplikasi di *user space* dan kendali perangkat keras di *kernel space* semakin tipis. Kita resmi memasuki era di mana sebuah daemon dapat menentukan irama detak prosesor sesuai kebutuhan sistem yang ia layani.

```references
- title: "sched_ext: Extensible BPF-based CPU Scheduler"
  author: "Linux Kernel Organization"
  url: "https://docs.kernel.org/scheduler/sched-ext.html"

- title: "The sched_ext Core Architecture and Upstream Patches"
  author: "Tejun Heo, David Vernet"
  year: 2024
  url: "https://lwn.net/Articles/922258/"

- title: "SCX: Sched-ext Schedulers and Tools Repository"
  author: "sched-ext Community"
  year: 2024
  url: "https://github.com/sched-ext/scx"
```
