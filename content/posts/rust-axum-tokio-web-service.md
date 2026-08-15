---
title: "Membangun Microservice Berkinerja Tinggi dengan Axum dan Tokio"
slug: b3e81a5f
aliases: [rust-axum-tokio-web-service]
date: 2026-06-20
tags: [rust, backend, performance]
lang: id
draft: false
type: post
---

Axum adalah kerangka kerja web berkinerja tinggi untuk ekosistem Rust yang dibangun di atas runtime Tokio dan abstraksi Tower. Kombinasi ini memberikan garansi type safety saat penanganan HTTP request tanpa mengorbankan kecepatan eksekusi async I/O.

## Fun Fact

**Fact 1.** Axum dikembangkan langsung oleh tim pengembang utama Tokio, sehingga mengintegrasikan lapisan Hyper dan Tower secara native tanpa wrapper overhead.

**Fact 2.** Sistem extractor Axum memeriksa validitas tipe data header, query parameter, dan payload request pada saat waktu kompilasi (compile-time).

**Fact 3.** Aturan kepemilikan memori (ownership) dan peminjaman (borrowing) Rust menjamin operasi async I/O di Tokio bebas dari bahaya data race secara komputasional.

---

## Tips dan Trik

### 1. Desain Extractor Berbasis Tipe (Type-Safe Extraction)

Axum menggunakan trait `FromRequest` atau `FromRequestParts` untuk mengekstrak data dari HTTP request secara deklaratif. Jika tipe data tidak sesuai, Axum secara otomatis mengembalikan respon kesalahan tanpa masuk ke logika handler:

```rust
use axum::{extract::{Path, State}, Json, routing::get, Router};
use serde::{Deserialize, Serialize};
use std::sync::Arc;

struct AppState {
    db_name: String,
}

#[derive(Serialize)]
struct UserResponse {
    id: u64,
    name: String,
}

async fn get_user_handler(
    State(state): State<Arc<AppState>>,
    Path(user_id): Path<u64>,
) -> Json<UserResponse> {
    Json(UserResponse {
        id: user_id,
        name: format!("User-{} dari database {}", user_id, state.db_name),
    })
}
```

### 2. Middleware Komposabel dengan Tower Services

Tower memungkinkan penggabungan layer middleware seperti logging, rate-limiting, dan CORS tanpa memodifikasi kode handler utama:

```rust
use tower_http::trace::TraceLayer;
use tower_http::cors::{CorsLayer, Any};
use std::net::SocketAddr;

#[tokio::main]
async fn main() {
    let shared_state = Arc::new(AppState {
        db_name: String::from("production_db"),
    });

    let app = Router::new()
        .route("/users/:id", get(get_user_handler))
        .layer(TraceLayer::new_for_http())
        .layer(CorsLayer::new().allow_origin(Any))
        .with_state(shared_state);

    let addr = SocketAddr::from(([127, 0, 0, 1], 3000));
    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
```

### 3. Penanganan Concurrency Async I/O Tanpa Data Race

Runtime Tokio mengelola alokasi thread pool berbasis work-stealing. Untuk membagi state antar thread secara aman, gunakan abstraksi `Arc<T>` atau `Arc<Mutex<T>>`.

Karena sistem tipe Rust mewajibkan syarat `Send + Sync` pada data yang dikirim antar thread async, potensi race condition pada shared state dapat dideteksi sebelum kode dikompilasi.

### 4. Benchmarking Throughput Request Per Detik Menggunakan wrk

Uji performa HTTP endpoint Axum menggunakan perkakas `wrk` dengan 12 thread dan 400 koneksi simultan:

```bash
wrk -t12 -c400 -d30s http://127.0.0.1:3000/users/42
```

Hasil pengujian pada perangkat keras modern umumnya menunjukkan throughput melebihi 100.000 request per detik dengan latensi rata-rata di bawah 1 milidetik.

### 5. Optimalisasi Profile Release pada Cargo.toml

Tambahkan konfigurasi kompilasi berikut di `Cargo.toml` untuk memaksimalkan performa binary produksi:

```toml
[profile.release]
opt-level = 3
lto = true
codegen-units = 1
panic = "abort"
strip = true
```
