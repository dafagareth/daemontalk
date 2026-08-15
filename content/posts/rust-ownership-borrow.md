---
title: "Memahami Sistem Ownership dan Borrow Checker Rust"
slug: d5b8e2f1
aliases: rust-ownership-borrow
date: 2026-07-08
tags: [rust, backend]
lang: id
draft: false
---

Borrow checker adalah komponen kompilator Rust yang memvalidasi aturan kepemilikan memori pada waktu kompilasi. Bagi pemula, pesan error dari borrow checker sering terasa tidak intuitif. Namun memahami tiga aturan dasar ownership akan menjelaskan hampir semua kasus yang membingungkan tersebut.

## Fakta Menarik

**Fakta 1.** Rust tidak memiliki garbage collector dan tidak memerlukan alokasi memori manual seperti C. Memori dibebaskan secara deterministik saat variabel keluar dari scope, melalui mekanisme yang disebut RAII (Resource Acquisition Is Initialization).

**Fakta 2.** Borrow checker beroperasi seluruhnya pada waktu kompilasi. Program Rust yang berhasil dikompilasi bebas dari data race, dangling pointer, dan use-after-free tanpa overhead runtime.

**Fakta 3.** Fitur `non-lexical lifetimes` (NLL) yang diperkenalkan di Rust 2018 Edition membuat borrow checker lebih cerdas: borrow dianggap berakhir saat penggunaan terakhirnya, bukan saat akhir blok.

---

## Tips dan Trik

### 1. Tiga Aturan Ownership

Seluruh sistem ownership Rust didasarkan pada tiga aturan:
- Setiap nilai memiliki tepat satu pemilik (owner).
- Saat pemilik keluar dari scope, nilai tersebut di-drop (memorinya dibebaskan).
- Hanya boleh ada satu mutable reference, atau banyak immutable reference, tetapi tidak keduanya secara bersamaan.

```rust
fn main() {
    // s1 adalah pemilik String ini
    let s1 = String::from("halo");

    // Kepemilikan berpindah (move) ke s2
    // s1 tidak lagi valid setelah baris ini
    let s2 = s1;

    // Baris berikut akan gagal kompilasi:
    // println!("{}", s1); // error: value borrowed here after move

    println!("{}", s2); // OK

    // s2 di-drop di sini karena keluar dari scope
}
```

### 2. Move vs Copy

Tipe yang mengimplementasikan trait `Copy` disalin secara bitwise saat diassign, sehingga variabel asal tetap valid. Tipe yang mengalokasikan heap (seperti `String`, `Vec`) menggunakan move.

```rust
fn main() {
    // i32 mengimplementasikan Copy: kedua variabel valid
    let x: i32 = 42;
    let y = x;
    println!("x={}, y={}", x, y); // OK

    // String tidak Copy: terjadi move
    let s1 = String::from("teks");
    let s2 = s1.clone(); // Clone membuat salinan dalam (deep copy)
    println!("s1={}, s2={}", s1, s2); // OK karena clone

    // Tipe custom dapat mengimplementasikan Copy jika semua field-nya Copy
    #[derive(Debug, Copy, Clone)]
    struct Titik {
        x: f64,
        y: f64,
    }

    let p1 = Titik { x: 1.0, y: 2.0 };
    let p2 = p1; // Copy, bukan move
    println!("{:?} dan {:?}", p1, p2); // Keduanya valid
}
```

### 3. Borrow dan Lifetime yang Implisit

Reference meminjam nilai tanpa mengambil kepemilikan. Kompilator menggunakan lifetime untuk memastikan reference tidak hidup lebih lama dari data yang dirujuknya.

```rust
// Lifetime elision: kompilator menyimpulkan lifetime secara otomatis
fn panjang_terpanjang<'a>(s1: &'a str, s2: &'a str) -> &'a str {
    if s1.len() >= s2.len() {
        s1
    } else {
        s2
    }
}

fn main() {
    let string1 = String::from("kalimat yang lebih panjang");
    let hasil;
    {
        let string2 = String::from("xy");
        hasil = panjang_terpanjang(&string1, &string2);
        println!("Terpanjang: {}", hasil); // OK: hasil digunakan di sini
    }
    // Menggunakan hasil di luar blok ini akan gagal jika string2
    // memiliki lifetime yang lebih pendek dari string1.
}
```

```rust
// Aturan borrow dalam praktik
fn main() {
    let mut data = vec![1, 2, 3];

    // Immutable borrow
    let pertama = &data[0];
    println!("Elemen pertama: {}", pertama);
    // pertama tidak digunakan lagi di bawah ini (NLL)

    // Mutable borrow sekarang aman karena immutable borrow sudah selesai
    data.push(4);
    println!("Data: {:?}", data);
}
```

### 4. Kapan Menggunakan `Arc<Mutex<T>>`

Saat data perlu dibagi antar thread, ownership tunggal tidak mencukupi. `Arc` (Atomic Reference Counted) memungkinkan kepemilikan bersama antar thread, dan `Mutex` menjamin akses eksklusif saat modifikasi.

```rust
use std::sync::{Arc, Mutex};
use std::thread;

fn main() {
    let penghitung = Arc::new(Mutex::new(0u64));
    let mut handles = vec![];

    for _ in 0..8 {
        let penghitung_clone = Arc::clone(&penghitung);
        let h = thread::spawn(move || {
            let mut nilai = penghitung_clone.lock().unwrap();
            *nilai += 1;
            // MutexGuard di-drop di sini, lock dilepas
        });
        handles.push(h);
    }

    for h in handles {
        h.join().unwrap();
    }

    println!("Hasil akhir: {}", *penghitung.lock().unwrap());
    // Output: Hasil akhir: 8
}
```

Gunakan `Arc<Mutex<T>>` hanya ketika mutasi lintas thread memang diperlukan. Untuk data read-only yang dibagi antar thread, cukup gunakan `Arc<T>`.

### 5. Mengapa Rust Cocok untuk Kode Systems Level

Rust menjamin keamanan memori tanpa garbage collector, sehingga latensi bersifat deterministik. Tidak ada jeda GC yang tidak terduga, tidak ada dangling pointer, dan tidak ada data race.

```rust
// Contoh: implementasi buffer ring sederhana tanpa unsafe
struct RingBuffer<T, const N: usize> {
    data: [Option<T>; N],
    kepala: usize,
    ekor: usize,
}

impl<T: Copy, const N: usize> RingBuffer<T, N> {
    fn baru() -> Self {
        RingBuffer {
            data: [None; N],
            kepala: 0,
            ekor: 0,
        }
    }

    fn push(&mut self, item: T) -> bool {
        let berikutnya = (self.ekor + 1) % N;
        if berikutnya == self.kepala {
            return false; // Buffer penuh
        }
        self.data[self.ekor] = Some(item);
        self.ekor = berikutnya;
        true
    }

    fn pop(&mut self) -> Option<T> {
        if self.kepala == self.ekor {
            return None; // Buffer kosong
        }
        let item = self.data[self.kepala].take();
        self.kepala = (self.kepala + 1) % N;
        item
    }
}

fn main() {
    let mut buf: RingBuffer<u32, 4> = RingBuffer::baru();
    buf.push(10);
    buf.push(20);
    println!("{:?}", buf.pop()); // Some(10)
    println!("{:?}", buf.pop()); // Some(20)
    println!("{:?}", buf.pop()); // None
}
```

Kompilator Rust memverifikasi bahwa tidak ada akses di luar batas, tidak ada use-after-free, dan tidak ada kondisi balapan, semuanya tanpa biaya runtime.
