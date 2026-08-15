---
title: "io_uring: Antarmuka I/O Asinkron Generasi Baru di Linux"
slug: d1a8c6f2
aliases: [kernel-io-uring-async]
date: 2026-07-20
tags: [linux, performance, backend]
lang: id
cover: /static/images/io_uring.png
draft: false
---
`io_uring` adalah antarmuka I/O asinkron yang diperkenalkan di Linux 5.1 (2019) oleh Jens Axboe. Desainnya mengatasi keterbatasan mendasar pada `epoll` dan `aio` lama, dan kini diadopsi oleh server-server produksi berperforma tinggi seperti Nginx, tokio, dan libSQL.

## Fakta Menarik

**Fakta 1.** `aio` Linux yang lama hanya mendukung operasi asinkron pada file yang dibuka dengan flag `O_DIRECT`. Operasi pada socket atau file reguler tanpa `O_DIRECT` akan kembali ke mode sinkron secara diam-diam, membuat banyak kode yang mengira dirinya asinkron sebenarnya berjalan secara blokir.

**Fakta 2.** `io_uring` dapat beroperasi dalam mode `SQPOLL` di mana kernel menjalankan thread polling khusus untuk menguras submission queue. Dalam mode ini, aplikasi tidak perlu melakukan syscall sama sekali untuk mengirimkan I/O, sehingga overhead context switch menjadi nol untuk beban kerja yang sibuk.

**Fakta 3.** Pada benchmark jaringan dengan koneksi yang sangat banyak, `io_uring` dapat mencapai throughput 20-30% lebih tinggi dari `epoll` karena pengurangan jumlah syscall dan penghapusan salinan data yang tidak perlu antar ruang kernel dan pengguna.

---

## Tips dan Trik

### 1. Memahami Masalah dengan epoll dan aio Lama

`epoll` bekerja dengan model reaktif: aplikasi menunggu kernel memberi tahu bahwa file descriptor siap, lalu melakukan operasi I/O (yang bisa saja memblokir). Setiap operasi memerlukan setidaknya dua syscall: `epoll_wait` dan kemudian `read`/`write`.

`aio` lama memiliki API yang rumit, dukungan operasi yang terbatas, dan perilaku fallback ke sinkron yang tidak dapat diprediksi. Kedua mekanisme tersebut tidak dirancang untuk menggabungkan berbagai jenis operasi (file, socket, timer) dalam satu event loop yang benar-benar asinkron.

### 2. Arsitektur Submission Queue dan Completion Queue

`io_uring` memperkenalkan dua ring buffer yang dipetakan ke memori bersama antara kernel dan aplikasi:

```
Aplikasi                     Kernel
---------                    ------
SQ (Submission Queue)  --->  Kernel memproses SQE
                             dan mengisi CQE
CQ (Completion Queue)  <---  Aplikasi membaca CQE
```

- **SQE (Submission Queue Entry):** Struktur yang mendeskripsikan operasi yang diminta (baca, tulis, accept, dll).
- **CQE (Completion Queue Entry):** Hasil operasi yang dikembalikan kernel, berisi kode hasil dan data pengguna.

Karena menggunakan memori bersama, pengiriman dan penerimaan hasil tidak memerlukan salinan data antara ruang kernel dan pengguna.

### 3. Penggunaan Dasar dengan liburing dari C

`liburing` adalah pustaka tipis yang menyederhanakan penggunaan `io_uring`. Instalasi:

```bash
# Debian/Ubuntu
sudo apt install liburing-dev

# Fedora/RHEL
sudo dnf install liburing-devel
```

Contoh membaca file secara asinkron:

```c
#include <fcntl.h>
#include <stdio.h>
#include <string.h>
#include <liburing.h>

#define BUF_SIZE 4096

int main(void) {
    struct io_uring ring;
    struct io_uring_sqe *sqe;
    struct io_uring_cqe *cqe;
    char buf[BUF_SIZE];
    int fd, ret;

    // Inisialisasi io_uring dengan kedalaman queue 8
    ret = io_uring_queue_init(8, &ring, 0);
    if (ret < 0) { perror("io_uring_queue_init"); return 1; }

    fd = open("/etc/hostname", O_RDONLY);
    if (fd < 0) { perror("open"); return 1; }

    // Ambil SQE dan isi dengan operasi baca
    sqe = io_uring_get_sqe(&ring);
    io_uring_prep_read(sqe, fd, buf, BUF_SIZE, 0);
    sqe->user_data = 42; // penanda bebas untuk identifikasi

    // Kirim ke kernel
    io_uring_submit(&ring);

    // Tunggu satu completion
    ret = io_uring_wait_cqe(&ring, &cqe);
    if (ret < 0) { perror("io_uring_wait_cqe"); return 1; }

    printf("Terbaca %d byte: %.*s\n", cqe->res, cqe->res, buf);
    io_uring_cqe_seen(&ring, cqe);

    io_uring_queue_exit(&ring);
    close(fd);
    return 0;
}
```

Kompilasi dan jalankan:

```bash
gcc -O2 -o read_example read_example.c -luring
./read_example
```

### 4. Mengirimkan Beberapa Operasi Sekaligus (Batching)

Kekuatan utama `io_uring` terletak pada kemampuan mengirimkan banyak operasi dalam satu syscall:

```c
// Kirim 4 operasi baca sekaligus tanpa syscall tambahan
for (int i = 0; i < 4; i++) {
    sqe = io_uring_get_sqe(&ring);
    io_uring_prep_read(sqe, fds[i], bufs[i], BUF_SIZE, 0);
    sqe->user_data = i;
}

// Satu syscall untuk semua
io_uring_submit(&ring);

// Tunggu semua completion
int completed = 0;
while (completed < 4) {
    ret = io_uring_wait_cqe(&ring, &cqe);
    printf("Operasi %llu selesai: %d byte\n",
           (unsigned long long)cqe->user_data, cqe->res);
    io_uring_cqe_seen(&ring, cqe);
    completed++;
}
```

### 5. Perbandingan Throughput vs epoll dan Mengapa Aplikasi Modern Mengadopsinya

Pengukuran sederhana dengan `fio` untuk membandingkan metode I/O asinkron pada NVMe lokal:

```bash
# Test io_uring
fio --name=iouring_test --ioengine=io_uring --rw=randread \
    --bs=4k --numjobs=4 --iodepth=128 --runtime=30 \
    --filename=/tmp/testfile --size=2G --group_reporting

# Test libaio (aio lama) untuk perbandingan
fio --name=aio_test --ioengine=libaio --rw=randread \
    --bs=4k --numjobs=4 --iodepth=128 --runtime=30 \
    --filename=/tmp/testfile --size=2G --group_reporting
```

Runtime modern seperti `tokio` (Rust) dan `glommio` sudah menggunakan `io_uring` sebagai backend I/O utama pada kernel yang mendukungnya. PostgreSQL 16+ memiliki opsi eksperimental untuk menggunakan `io_uring` pada operasi WAL, dan NGINX telah memperkenalkan dukungan `io_uring` melalui thread pool untuk operasi file statis.
