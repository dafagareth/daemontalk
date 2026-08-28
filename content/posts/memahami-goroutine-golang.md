---
title: "Memahami Goroutine dalam Bahasa Go"
slug: memahami-goroutine-golang
aliases: []
date: 2026-08-28
author: "Daemontalk Editorial"
tags: ["Golang", "Concurrency"]
lang: id
draft: false
type: post
cover: "https://images.unsplash.com/photo-1555066931-4365d14bab8c?auto=format&fit=crop&q=80&w=1600"
coverCaption: "Code on screen"
coverSource: "https://unsplash.com"
readTime: 5
description: "Goroutine adalah inti dari concurrency di Go. Pelajari cara kerjanya, manajemen sinkronisasi dengan WaitGroup, dan komunikasi via Channel."
---

Goroutine adalah fungsi atau metode yang dieksekusi secara konkuren bersama dengan goroutine lain dalam *address space* yang sama.

Berbeda dengan *thread* sistem operasi (OS) tradisional yang bisa memakan memori hingga hitungan megabyte per thread, goroutine sangat ringan (membutuhkan memori awal sekitar 2KB). Go runtime secara otomatis mengelola goroutine (multiplexing) di atas thread OS.

## Konsep Dasar dan Cara Kerja

Saat Anda menggunakan keyword `go`, Go *scheduler* (yang berjalan di *background*) akan mengambil fungsi tersebut dan menugaskannya ke salah satu *logical processor* yang tersedia.

```go
package main

import (
	"fmt"
	"time"
)

func cetakPesan(pesan string) {
	for i := 0; i < 3; i++ {
		fmt.Println(pesan) // [!code hl]
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	// Fungsi dieksekusi di background
	go cetakPesan("Dari Goroutine") // [!code ++]
	
	// Dieksekusi secara sinkron
	cetakPesan("Dari Main")
}
```

> [!NOTE]
> Fungsi `main` pada sebuah program Go juga berjalan di dalam goroutine khusus, yang sering disebut sebagai **main goroutine**. Jika *main goroutine* selesai, program akan berhenti tanpa peduli goroutine lain sudah selesai atau belum.

## Sinkronisasi Menggunakan WaitGroup

Untuk mengatasi masalah *main goroutine* yang selesai terlalu cepat, kita harus menunggu semua goroutine anak (*child goroutines*) selesai. Cara terbaik adalah menggunakan `sync.WaitGroup`.

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func pekerja(id int, wg *sync.WaitGroup) {
	defer wg.Done() // Memberitahu WaitGroup bahwa tugas selesai
	fmt.Printf("Pekerja %d mulai...\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Pekerja %d selesai.\n", id)
}

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1) // Menambah counter goroutine
		go pekerja(i, &wg)
	}

	wg.Wait() // Memblokir eksekusi sampai counter kembali ke 0
	fmt.Println("Semua pekerja telah selesai.")
}
```

## Komunikasi Melalui Channel

Prinsip utama concurrency di Go adalah: *"Jangan berkomunikasi dengan berbagi memori (shared memory); berbagilah memori dengan berkomunikasi."*

**Channel** adalah pipa yang menghubungkan antar goroutine.

Ada dua jenis channel:
1. **Unbuffered Channel**: Pengirim akan diblokir sampai ada penerima yang mengambil datanya.
2. **Buffered Channel**: Memiliki antrean. Pengirim hanya diblokir jika antrean penuh.

```go
// Membuat channel string
pesanCh := make(chan string)

// Mengirim data ke dalam channel
go func() {
    pesanCh <- "Data siap"
}()

// Menerima data dari channel
hasil := <-pesanCh
fmt.Println(hasil)
```

## Praktik Terbaik & Pencegahan Kebocoran

1. **Hindari Goroutine Leak**: Pastikan setiap goroutine yang dibuat memiliki cara untuk berhenti atau selesai dieksekusi, misalnya menggunakan `context.Context` untuk pembatalan (*cancellation*).
2. **Batasi Concurrency Berlebihan**: Meski ringan, menjalankan jutaan goroutine sekaligus untuk mengakses database bisa membuat database Anda *crash*. Gunakan pola *worker pool*.

> [!WARNING]
> Jangan biarkan goroutine berjalan tanpa henti di background (infinite loop) tanpa mekanisme kontrol, karena ini akan menyebabkan *memory leak* (kebocoran memori) yang menumpuk seiring waktu jalannya aplikasi.

## Referensi Terverifikasi

```references
- title: "Effective Go: Concurrency"
  author: "The Go Authors"
  year: 2024
  publisher: "golang.org"
  url: "https://go.dev/doc/effective_go#concurrency"
  
- title: "Go Concurrency Patterns"
  author: "Rob Pike"
  year: 2012
  publisher: "Go Blog"
  url: "https://go.dev/blog/io2012-videos"
```
