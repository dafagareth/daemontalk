---
# ==============================================================================
# DAEMONTALK ARTICLE FRONTMATTER SPECIFICATION
# ==============================================================================
# title: (Wajib) Judul artikel yang jelas, lugas, dan bebas karakter `#`.
title: "Panduan Lengkap Arsitektur Sistem Terdistribusi dan Fitur Markdown"

# slug: (Wajib) Identifikator URL unik (/blog/{slug}). Gunakan huruf kecil & strip.
slug: panduan-arsitektur-sistem

# aliases: (Opsional) URL lama yang otomatis dialihkan (redirect 301) ke artikel ini.
aliases: [arsitektur-sistem-modern, draft-rekayasa-kernel]

# date: (Wajib) Tanggal rilis artikel (Format: YYYY-MM-DD) untuk pengurutan linimasa.
date: 2026-08-14

# tags: (Wajib) Kategori artikel untuk indexing dan pencarian.
tags: [architecture, linux, go, performance]

# lang: (Wajib) Bahasa artikel ('id' untuk Indonesia, 'en' untuk English).
lang: id

# draft: (Opsional) Set 'false' untuk rilis publik, atau 'true' untuk mode draft internal.
draft: false

# type: (Opsional) 'post' untuk artikel standar atau 'til' untuk Today I Learned.
type: post

# cover: (Opsional) Gambar sampul utama artikel.
cover: "/static/logo/logo-dark.png"

# series: (Opsional) Mengelompokkan postingan ke dalam rangkaian seri topik.
series: "Distributed Systems Engineering"

# series_part: (Opsional) Urutan nomor bab dalam seri (1, 2, 3, dst).
series_part: 1

# summary: Ringkasan singkat untuk SEO dan kartu pratinjau.
summary: "Panduan komprehensif penulisan artikel rekayasa sistem, pemanfaatan carousel, galeri gambar, diagram ASCII, FAQ interaktif, dan kartu profil penulis."
---

Paragraf pembuka berfungsi sebagai intisari teknis dari keseluruhan dokumen[^1]. Kalimat pertama ini dirancang ringkas dan padat karena otomatis diekstraksi oleh parser sebagai *meta description* untuk mesin pencari (SEO) dan preview kartu media sosial.

---

## 1. Tipografi dan Penekanan Teks

Dokumen ini mendukung seluruh elemen tipografi standar dengan hierarki visual yang jelas:

- **Teks Tebal (*Bold*)**: Digunakan untuk menonjolkan terminologi penting seperti **Zero-Copy Memory**.
- *Teks Miring (*Italic*)**: Digunakan untuk istilah asing atau variabel matematis seperti *throughput limit* $O(1)$.
- ~~Teks Coret (*Strikethrough*)~~: Digunakan untuk menandai pendekatan lama yang sudah ditinggalkan.
- `Kode Sebaris (*Inline Code*)`: Digunakan untuk nama fungsi, flag kernel, atau perintah CLI seperti `sysctl net.core.somaxconn` dan `epoll_create1()`.
- [Tautan Eksternal (*Hyperlink*)](https://kernel.org): Tautan terisolasi yang aman dengan style kontras tinggi.

> **Prinsip Rekayasa:** Hindari alokasi memori dinamis di dalam *hot path* pemrosesan paket data. Selalu manfaatkan pool memori lokal untuk meminimalkan jeda *garbage collection*.

---

## 2. Penyisipan Media: Carousel & Galeri

### A. Carousel Gambar (Swipeable Slider)
Gunakan blok ` ```carousel ` murni untuk menyajikan rangkaian gambar atau diagram alur secara interaktif:

```carousel
![Arsitektur io_uring](/static/images/io_uring.png "Gambar 1: Topologi kernel async I/O ring buffer")
![Go Runtime Worker Engine](/static/images/golang.png "Gambar 2: Pipeline pemrosesan goroutine terdistribusi")
![Arch Linux Environment](/static/images/archlinux.png "Gambar 3: Lingkungan pengujian kernel 6.12 LTS")
```

### B. Galeri Gambar (Side-by-Side Grid)
Gunakan blok ` ```gallery ` untuk menyajikan perbandingan visual secara berdampingan:

```gallery
![Arch Linux Host](/static/images/archlinux.png "Arsitektur Host Kernel")
![Go Performance](/static/images/golang.png "Throughput Profiling Go")
```

---

## 3. Diagram Arsitektur Interaktif

Blok diagram ASCII atau box-drawing otomatis dideteksi oleh sistem antarmuka dan dilengkapi dengan tombol interaktif **Expand/Zoom** untuk kenyamanan membaca di layar mobile maupun desktop:

```text
┌─────────────────────────┐
│     Inbound Traffic     │
│   (HTTPS / TCP:443)     │
└────────────┬────────────┘
             │
             ▼
┌─────────────────────────┐         Zero-Copy Buffer        ┌─────────────────────────┐
│   eBPF XDP Layer        │ ──────────────────────────────► │   io_uring Worker Ring  │
│   (Kernel Bypass Drop)  │                                 │   (Fixed Registered FD) │
└────────────┬────────────┘                                 └────────────┬────────────┘
             │                                                           │
             │ Non-blocking Pipeline                                     ▼
             └─────────────────────────────────────────────► ┌─────────────────────────┐
                                                             │   Go Application Core   │
                                                             └─────────────────────────┘
```

---

## 4. Implementasi Kode dan Snippet Beraneka Bahasa

Seluruh blok kode disorot menggunakan *syntax highlighting* Chroma dengan font **JetBrains Mono** dan dilengkapi tombol salin kode (*one-click copy*):

### Contoh Implementasi Go (Worker Engine)

```go
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type Engine struct {
	workers int
	queue   chan []byte
}

func NewEngine(workers int, bufferSize int) *Engine {
	return &Engine{
		workers: workers,
		queue:   make(chan []byte, bufferSize),
	}
}

func (e *Engine) Start(ctx context.Context) {
	for i := 0; i < e.workers; i++ {
		go func(id int) {
			for {
				select {
				case <-ctx.Done():
					return
				case data := <-e.queue:
					fmt.Printf("[worker-%d] processed %d bytes\n", id, len(data))
				}
			}
		}(i)
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	engine := NewEngine(4, 1024)
	engine.Start(ctx)

	<-ctx.Done()
	fmt.Println("Graceful shutdown completed.")
}
```

### Contoh Konfigurasi Kernel Linux (`/etc/sysctl.d/99-network.conf`)

```ini
# Meningkatkan kapasitas antrean backlog koneksi masuk
net.core.somaxconn = 65535
net.ipv4.tcp_max_syn_backlog = 3240000

# Optimasi buffer TCP read/write untuk throughput tinggi
net.ipv4.tcp_rmem = 4096 87380 16777216
net.ipv4.tcp_wmem = 4096 65536 16777216
```

---

## 5. Matriks Perbandingan Performa

Gunakan tabel Markdown untuk menyajikan data kuantitatif atau benchmark perbandingan:

| Metode Sinkronisasi | Latensi P95 | Latensi P99 | Alokasi Memori / Operasi | Throughput Relatif |
| :--- | :--- | :--- | :--- | :--- |
| `sync.Mutex` Standar | 12.4 µs | 48.2 µs | 32 B/op | 1.0x (Baseline) |
| `sync/atomic.Value` | 3.1 µs | 8.6 µs | 0 B/op | 4.2x |
| Lock-Free Ring Buffer | 0.8 µs | 1.4 µs | 0 B/op | 14.8x |

---

## 6. Daftar Tugas Verifikasi Sistem (Task List)

- [x] Konfigurasi isolasi CPU core pada NUMA node 0.
- [x] Validasi driver antarmuka jaringan dengan `ethtool -k eth0`.
- [ ] Implementasi failover otomatis pada kluster multi-region.

---

## 7. Pertanyaan Operasional (Interactive FAQ)

Gunakan blok ` ```faq ` untuk membuat bagian tanya-jawab akordeon interaktif yang rapi dan mudah diakses:

```faq
Q: Kapan sebaiknya memilih channel vs mutex untuk sinkronisasi state di Go?
A: Gunakan channel ketika Anda mentransfer kepemilikan data antar goroutine (*passing data ownership*). Gunakan `sync.Mutex` atau operasi `sync/atomic` saat Anda hanya memproteksi state internal pada struktur data tunggal dengan durasi lock yang sangat singkat.

Q: Apakah arsitektur io_uring aman digunakan pada lingkungan multi-tenant?
A: Ya, dengan catatan kernel yang digunakan adalah versi 6.1 LTS ke atas dan pembatasan seccomp serta cgroups v2 telah diaktifkan untuk mengisolasi ring buffer per namespace.
```

---

## 8. Catatan Kaki (Footnotes)

Sistem secara otomatis menghubungkan nomor catatan kaki di teks dengan daftar referensi berikut:

[^1]: Dokumentasi ini disusun berdasarkan hasil pengujian beban di lingkungan server Linux x86_64 dengan kernel versi 6.12 LTS.

---

## 9. Bibliografi dan Referensi

1. **Gregg, Brendan.** (2020). *Systems Performance: Enterprise and the Cloud (2nd Edition)*. Addison-Wesley Professional.
2. **Love, Robert.** (2013). *Linux System Programming: Talking Directly to the Kernel and C Library*. O'Reilly Media.
3. **Axboe, Jens.** (2019). *Efficient IO with io_uring*. Kernel Documentation Archive.

---

## 10. Tentang Penulis (About Author)

Gunakan blok ` ```author ` untuk menampilkan identitas dan profil singkat penulis:

```author
name: Nama Penulis
role: Systems Engineer & Open Source Enthusiast
avatar: /static/logo/logo-dark.png
bio: Deskripsi singkat mengenai latar belakang teknis, fokus rekayasa, atau topik keahlian yang Anda bagikan.
github: https://github.com/username
email: author@example.com
```
