---
title: "Pentingnya Menulis Clean Code dalam Perangkat Lunak"
slug: pentingnya-menulis-clean-code-software
aliases: []
date: 2026-08-28
author: "Daemontalk Editorial"
tags: ["Software", "Engineering"]
lang: id
draft: false
type: post
cover: "https://images.unsplash.com/photo-1542831371-29b0f74f9713?auto=format&fit=crop&q=80&w=1600"
coverCaption: "Kode pemrograman di monitor"
coverSource: "https://unsplash.com"
readTime: 5
description: "Kode ditulis untuk dibaca oleh manusia, bukan hanya dieksekusi oleh mesin. Pelajari prinsip dasar penamaan yang baik, single responsibility, dan refactoring."
---

Dalam industri rekayasa perangkat lunak (*software engineering*), waktu yang dihabiskan seorang *programmer* untuk **membaca** kode jauh lebih banyak daripada waktu yang dihabiskan untuk **menulis** kode baru.

Rasio membaca berbanding menulis bisa mencapai 10:1. Oleh karena itu, menulis kode yang bisa bekerja (*working code*) saja tidak cukup. Kode tersebut harus bersih, terstruktur, dan mudah dipahami oleh anggota tim lain (*clean code*).

## Prinsip Penamaan yang Bermakna

Nama variabel, fungsi, atau kelas harus secara eksplisit memberi tahu pembaca mengapa ia ada, apa yang dilakukannya, dan bagaimana ia digunakan.

```go
// BURUK:
// Apa itu d? Hari? Jarak?
var d int 

// BAIK:
var elapsedTimeInDays int
var fileSizeInBytes int
```

> [!TIP]
> Jangan takut menggunakan nama variabel yang agak panjang asalkan deskriptif. Editor modern sudah memiliki fitur *autocomplete*, sehingga Anda tidak perlu mengetik panjang lebar.

## Aturan Fungsi yang Baik

Robert C. Martin (Uncle Bob) dalam bukunya *Clean Code* menyatakan dua aturan utama untuk membuat fungsi:
1. Fungsi harus kecil.
2. Fungsi harus lebih kecil lagi.

Fungsi hanya boleh melakukan **satu hal** (Single Responsibility Principle) dan harus melakukannya dengan baik.

```go
// BURUK: Satu fungsi melakukan segalanya (Validasi, Database, Email)
func registerUser(req Request) error { // [!code --]
	if req.Email == "" { return errors.New("empty email") } // [!code --]
	db.Insert(req) // [!code --]
	sendWelcomeEmail(req.Email) // [!code --]
	return nil // [!code --]
} // [!code --]

// BAIK: Delegasi tugas ke fungsi spesifik
func registerUser(req Request) error { // [!code ++]
	if err := validateRequest(req); err != nil { return err } // [!code ++]
	if err := saveToDatabase(req); err != nil { return err } // [!code ++]
	go sendWelcomeEmail(req.Email) // [!code ++]
	return nil // [!code ++]
} // [!code ++]
```

## Bahaya "Dead Code" dan Komentar Berlebihan

Seringkali *programmer* meninggalkan baris kode yang di-*comment out* karena takut suatu saat akan dibutuhkan lagi. Ini adalah kebiasaan buruk yang mengotori basis kode.

Percayakan riwayat kode pada sistem *Version Control* (seperti Git). Jika suatu saat Anda butuh kode lama tersebut, Anda bisa mencarinya di histori Git. Jangan biarkan *dead code* membusuk di *file* produksi.

> [!WARNING]
> Komentar kode tidak boleh digunakan untuk menutupi kode yang berantakan. Jika kode Anda butuh paragraf panjang agar bisa dipahami, itu pertanda kode tersebut harus di-*refactor*. Komentar hanya digunakan untuk menjelaskan **"mengapa"** (keputusan bisnis), bukan **"bagaimana"** kode itu bekerja.

## Referensi Terverifikasi

```references
- title: "Clean Code: A Handbook of Agile Software Craftsmanship"
  author: "Robert C. Martin"
  year: 2008
  publisher: "Prentice Hall"

- title: "The Pragmatic Programmer"
  author: "Andrew Hunt & David Thomas"
  year: 1999
  publisher: "Addison-Wesley"
```
