---
title: "Go Pprof: Menemukan CPU & Memory Leak dalam 5 Menit"
slug: 1a2b3c4d
aliases: [go-profiling-pprof]
date: 2026-05-18
tags: [go, performance, backend]
lang: id
draft: false
---

Mengoptimalkan performa program tanpa profiling ibarat menembak dalam kegelapan. Go hadir dengan tool profiling kelas dunia bernama **pprof** yang sudah tertanam langsung di standard library.

Hanya dengan satu baris import, kamu bisa melihat flame graph konsumsi CPU, alokasi memori heap, dan goroutine blocking secara visual di browser.

## Fun Fact

**Pprof awalnya dikembangkan oleh tim Google Performance Tools.**
Tool ini kemudian diintegrasikan secara native ke dalam runtime Go oleh tim inti Go di Google.

**Pprof bekerja dengan teknik Statistical Sampling.**
Profiler CPU Go mengambil sample stack trace setiap 10 milidetik (100 Hz), sehingga overhead performa saat profiling aktif sangat rendah (biasanya < 2%), aman dijalankan di server production saat traffic normal.

**Mendukung visualisasi Flame Graph interaktif secara native.**
Mulai Go 1.11, command line tool `go tool pprof` sudah memiliki web server internal yang merender flame graph interaktif tanpa perlu dependensi Perl eksternal.

---

## Tips dan Trik

### 1. Pasang Profiler di HTTP Server dengan Satu Baris

Cukup tambahkan side-effect import `net/http/pprof`:

```go
package main

import (
    "net/http"
    _ "net/http/pprof" // Mendaftarkan endpoint /debug/pprof/
)

func main() {
    go func() {
        http.ListenAndServe("localhost:6060", nil)
    }()
    // Logika utama aplikasi kamu...
}
```

### 2. Analisis CPU Profile dengan Web UI Visual

Ambil snapshot CPU selama 30 detik dan buka visualisasi browser secara otomatis:

```bash
go tool pprof -http=:8081 http://localhost:6060/debug/pprof/profile?seconds=30
```

### 3. Cari Memory Leak pada Heap Allocation

Lihat fungsi mana yang mengalokasikan memori paling banyak tapi tidak dibersihkan oleh Garbage Collector:

```bash
go tool pprof -http=:8081 http://localhost:6060/debug/pprof/heap
# Di menu UI, pilih 'View' -> 'Flame Graph' atau 'Top'
```

### 4. Lacak Goroutine Leak

Jika jumlah goroutine terus merangkak naik dari waktu ke waktu:

```bash
curl http://localhost:6060/debug/pprof/goroutine?debug=1
```

### 5. Benchmark Kode Tertentu dengan Flag `-cpuprofile` dan `-memprofile`

Jalankan profiling langsung dari unit test atau benchmark:

```bash
go test -bench=. -benchmem -cpuprofile=cpu.pprof -memprofile=mem.pprof
go tool pprof -http=:8081 cpu.pprof
```
