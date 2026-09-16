---
title: "WASM Akhirnya Dewasa"
slug: "wasm-component-model-dewasa"
aliases: []
author: "daemontalk team"
contributors: []
tags: ["radar", "wasm", "systems"]
lang: "id"
status: "published"
type: "dispatch"
readTime: 4
cover: "https://images.unsplash.com/photo-1558494949-ef010cbdcc31?auto=format&fit=crop&w=1200&q=80"
coverCaption: "Modular dan portabel, tanpa glue code manual"
coverSource: "https://unsplash.com"
description: "WebAssembly Component Model dan WASI Preview 2 akhirnya stabil. WebAssembly kini siap jadi fondasi plugin lintas bahasa dan runtime di luar peramban."
series: ""
series_part: 0
---

<span class="drop-cap">S</span>elama bertahun-tahun, WebAssembly (WASM) terjebak dalam siklus janji manis. Portabilitas mutlak, eksekusi dalam sandbox aman, dan performa mendekati native selalu jadi jargon utama.

Bagi backend engineer yang membangun sistem di atas Linux, kenyataannya berbeda. WASM sering terasa seperti mainan.

Anda bisa mengompilasi kode Rust atau C++ menjadi biner ringkas. Namun begitu modul berjalan, ia terisolasi. Mengirim data teks atau struktur data dari host Go ke modul WASM memaksa Anda mengelola pointer memori linear dan menulis glue code yang rapuh. 

WASM cepat dalam komputasi angka mentah, tetapi kaku dalam interoperabilitas. Kehadiran Component Model dan WASI Preview 2 akhirnya mengubah keadaan ini.

## Keterbatasan Modul Klasik

WASM versi awal (MVP) hanya mengenali empat tipe data dasar: i32, i64, f32, dan f64. Tidak ada konsep native untuk string, struct, maupun list.

Jika sebuah fungsi menerima string, fungsi tersebut pada tingkat biner hanya menerima angka pointer ke blok memori linear. Runtime host harus menyalin byte string ke lokasi memori tersebut secara manual sebelum modul bisa membacanya.

Masalah makin rumit saat dua modul dari bahasa berbeda ingin berkomunikasi. Rust dan Go memiliki alokator serta representasi memori internal yang berlainan. Tanpa standar bersama, menyatukan keduanya dalam satu alur kerja membutuhkan lapisan pembungkus yang tebal dan lambat.

## WIT sebagai Kontrak Bersama

Component Model menyelesaikan kebuntuan ini dengan memperkenalkan tipe data tingkat tinggi: string, record, variant, dan list.

Antarmuka antarkomponen didefinisikan secara deklaratif menggunakan WebAssembly Interface Types (WIT). Formatnya ringkas dan berfungsi mirip IDL pada gRPC:

```wit
package daemontalk:example;

interface logger {
    log: func(level: string, message: string);
}

world service {
    import logger;
    export process: func(data: list<u8>) -> result<string, string>;
}
```

Tooling otomatis menghasilkan kode binding untuk bahasa target Anda. Jika komponen ditulis dalam Rust dan dijalankan di atas host Go via Wasmtime, transfer data ditangani oleh runtime secara aman. 

Kedua pihak tidak perlu saling tahu manajemen memori internal lawan bicaranya. Anda bisa merangkai router Go dengan middleware Rust dalam satu proses tanpa overhead jaringan.

## WASI Preview 2 Lepas dari POSIX

WASI Preview 1 mencoba membuat WASM bertingkah seperti proses Linux biasa. Masalahnya, memaksakan abstraksi POSIX ke dalam sandbox berbasis kapabilitas menimbulkan banyak gesekan, terutama pada I/O asinkron dan soket jaringan.

WASI Preview 2 membuang pendekatan tersebut dan beralih penuh ke antarmuka WIT. Operasi file, stream, HTTP, dan soket kini bertipe data ketat dan modular.

Keamanan sistem menjadi lebih transparan. Jika sebuah komponen hanya diberi izin membaca file konfigurasi, runtime memastikan modul tersebut tidak bisa membuka koneksi keluar.

## Dampak Nyata untuk Sistem

Bagi arsitek sistem, ini adalah standar plugin yang sesungguhnya. Anda tidak perlu lagi mendikte bahasa ekstensi pada basis data atau API gateway Anda. Cukup sematkan runtime WASM, lalu biarkan pengguna mengirimkan biner dari bahasa pilihan mereka.

Untuk beban kerja serverless, cold start kontainer yang berkisar puluhan hingga ratusan milidetik dapat dipangkas mendekati nol. Modul ringan bisa dimuat dan dieksekusi dalam hitungan mikrodetik.

WASM bukan lagi sekadar eksperimen peramban. Dengan Component Model dan WASI Preview 2, fondasi komputasi modular lintas bahasa akhirnya siap digunakan di lingkungan produksi.
