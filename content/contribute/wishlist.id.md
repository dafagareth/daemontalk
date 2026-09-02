## Daftar Topik Riset yang Dicari (*Call for Dispatches*)

Jika Anda ingin menulis untuk Daemontalk namun belum menentukan tema, berikut adalah topik-topik arsitektur dan rekayasa sistem berprioritas tinggi yang sangat dinantikan oleh komunitas:

---

### 1. Sistem Operasi & Linux Kernel

- **`[WANTED]` Pembedahan Arsitektur `io_uring` vs `epoll` pada High-Throughput I/O**: Mengapa `io_uring` mengeliminasi *context switch* dan bagaimana ring buffers kernel beroperasi.
- **`[WANTED]` Memory Allocators Shootout: `jemalloc` vs `mimalloc` vs Go Runtime Allocator**: Analisis fragmentasi memori, arena allocation, dan cache-friendly chunking.
- **`[WANTED]` eBPF Security Profiling dengan Aya (Rust) & Cilium**: Mencegat syscall berbahaya langsung pada level kernel tanpa kernel module eksternal.

---

### 2. Konkurensi & Runtimes

- **`[WANTED]` Deep-Dive Lock-Free Queue: Pembedahan Algoritma Michael-Scott**: Implementasi CAS atomik pada struktur data antrean tanpa mutex.
- **`[WANTED]` Go Goroutine Work-Stealing Scheduler Internals**: Bagaimana `sysmon`, network poller, dan algoritma *stealing* mendistribusikan beban CPU di multi-core.

---

### 3. Database Terdistribusi & Penyimpanan

- **`[WANTED]` Google Spanner Architecture: TrueTime API dan Atomic Clocks**: Bagaimana GPS dan jam atom menyelesaikan masalah konsistensi lintas data center global.
- **`[WANTED]` Write-Ahead Logging (WAL) & Crash Recovery Internals pada SQLite/Postgres**: Membedah checkpointing, ARIES recovery, dan fsync cost.

---

### 4. Jaringan & Kriptografi

- **`[WANTED]` Anatomi BGP Hijacking & Mitigasi RPKI**: Bagaimana routing internet global bisa dibelokkan dan cara kerja validasi kriptografi RPKI.
- **`[WANTED]` WireGuard Protocol Internals vs OpenVPN / IPsec**: Enkripsi modern berbasis Curve25519 dan performa kernel-space VPN.

---

## Ingin Mengambil Salah Satu Topik?

1. Unduh template: `curl -s https://daemontalk.com/daemontalk-template.md -o nama-topik.md`
2. Buka diskusi atau kirimkan draf artikel Anda melalui Pull Request ke [GitHub Repositori](https://github.com/dafagareth/daemontalk).
3. Setelah digabungkan (*merged*), profil Anda otomatis mendapatkan atribusi resmi sebagai **CONTRIBUTOR**.
