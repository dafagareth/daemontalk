---
title: "Go Concurrency: Pola Worker Pool dan Channel yang Aman"
slug: 7e8f9a0b
aliases: [go-concurrency-patterns]
date: 2026-07-15
tags: [go, backend, tips]
lang: id
draft: false
---

Model konkurensi Go berbasis CSP (*Communicating Sequential Processes*) dengan Goroutine dan Channel adalah salah satu fitur terbaik bahasa ini. Namun, tanpa struktur yang benar, goroutine leak dan race condition bisa merusak stabilitas aplikasi.

Berikut konsep mendasar dan pola implementasi worker pool yang teruji di skala produksi.

## Fun Fact

**Goroutine hanya butuh ~2KB memori saat pertama kali dibuat.**
Bandingkan dengan thread OS standar di C/Java yang biasanya mengalokasikan 1MB hingga 8MB stack memori. Go runtime bisa dengan mudah menangani ratusan ribu goroutine aktif secara simultan.

**Motto resmi Go: "Do not communicate by sharing memory; instead, share memory by communicating."**
Filosofi ini berasal langsung dari karya riset Tony Hoare pada tahun 1978 tentang CSP.

**Race Detector Go (`go run -race`) ditenagai oleh ThreadSanitizer Google.**
Fitur deteksi race condition Go tidak menggunakan analisis statis, melainkan instrumentasi kode saat runtime yang memantau akses memori setiap goroutine secara akurat.

---

## Tips dan Trik

### 1. Pola Standar Worker Pool dengan `sync.WaitGroup`

Gunakan sejumlah worker tetap untuk memproses antrean tugas agar tidak membebani database backend:

```go
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    for j := range jobs {
        results <- j * 2
    }
}

func main() {
    const numJobs = 100
    const numWorkers = 5
    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)
    var wg sync.WaitGroup

    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go worker(w, jobs, results, &wg)
    }

    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs)

    wg.Wait()
    close(results)
}
```

### 2. Selalu Tangani Cancelation dengan `context.Context`

Jangan biarkan goroutine background berjalan selamanya saat HTTP request dibatalkan oleh client:

```go
func processTask(ctx context.Context) error {
    select {
    case <-time.After(5 * time.Second):
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

### 3. Hati-hati dengan Closure Variable di Loop

Sebelum Go 1.22, variabel loop dibagikan ke semua iterasi. Di versi Go modern ini sudah diperbaiki, namun membiasakan passing argumen secara eksplisit tetap merupakan praktik yang baik:

```go
for _, item := range items {
    go func(val string) {
        process(val)
    }(item)
}
```

### 4. Gunakan `sync.Once` untuk Inisialisasi Singleton

Hindari double locking mutex yang kompleks untuk lazy initialization:

```go
var (
    dbConn *sql.DB
    once   sync.Once
)

func GetDB() *sql.DB {
    once.Do(func() {
        dbConn = initDatabase()
    })
    return dbConn
}
```

### 5. Jalankan Race Detector di CI/CD Pipeline

Jadikan flag `-race` sebagai bagian wajib dari unit test sebelum deploy:

```bash
go test -race -v ./...
```
