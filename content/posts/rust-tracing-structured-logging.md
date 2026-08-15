---
title: "Observabilitas dan Log Terstruktur di Rust dengan Crate Tracing"
slug: e5a1b4f9
aliases: [rust-tracing-structured-logging]
date: 2026-05-22
tags: [rust, backend, debugging]
lang: id
draft: false
type: post
---

Crate `tracing` di Rust menyediakan kerangka kerja untuk mengumpulkan data observabilitas berbasis rentang waktu (span) dan kejadian (event). Berbeda dari logging tradisional yang bersifat teks linier, `tracing` melacak aliran eksekusi konkuren pada aplikasi asinkron secara terstruktur. Tulisan ini membahas perbedaan logging konvensional dengan pemrosesan kontekstual, penggunaan makro `#[instrument]`, pengiriman trace ke OpenTelemetry, serta analisis efisiensi performanya.

## Fakta Menarik

**Fakta 1.** Crate `tracing` dikembangkan oleh tim Tokio project untuk menyelesaikan masalah pelacakan eksekusi asynchronous task yang saling tumpang tindih dalam runtime Tokio.

**Fakta 2.** Span pada `tracing` memiliki siklus hidup yang eksplisit (enter dan exit), yang secara otomatis menangani konteks eksekusi ketika task asinkron di-pause atau di-resume oleh executor.

**Fakta 3.** Penilaian kondisi penyaringan log pada `tracing` dilakukan secara terpusat oleh `Subscriber` sebelum bidang log dievaluasi, sehingga meminimalkan alokasi memori untuk log yang diabaikan.

---

## Tips dan Trik

### 1. Bedakan Penggunaan Span dan Event

Gunakan `span` untuk mencakup konteks rentang eksekusi (seperti penanganan HTTP request) dan `event` untuk mencatat kejadian instan di dalam span tersebut.

```rust
use tracing::{info, info_span};

fn process_user_order(user_id: u64) {
    let span = info_span!("process_order", user_id = user_id);
    let _guard = span.enter();

    /* Event ini secara otomatis mewarisi bidang user_id dari span aktif */
    info!(order_id = 9842, "Memproses transaksi pembayaran");
}
```

### 2. Gunakan Makro Instrument untuk Fungsi Asinkron

Makro `#[tracing::instrument]` secara otomatis membuat dan memasuki span setiap kali fungsi dipanggil, mencatat parameter fungsi sebagai field terstruktur.

```rust
use tracing::instrument;

#[instrument(skip(db_pool), fields(db.vendor = "postgres"))]
async fn fetch_user_profile(user_id: u64, db_pool: &sqlx::PgPool) -> Result<String, sqlx::Error> {
    /* Parameter db_pool diabaikan dari log untuk menghindari output verbose */
    Ok(format!("User_{}", user_id))
}
```

### 3. Konfigurasikan Format Output JSON dengan tracing-subscriber

Ubah output log teks mentah menjadi objek JSON yang siap diindeks oleh sistem manajemen log seperti Elasticsearch atau Loki.

```rust
use tracing_subscriber::fmt;
use tracing_subscriber::layer::SubscriberExt;
use tracing_subscriber::util::SubscriberInitExt;

fn init_logging() {
    tracing_subscriber::registry()
        .with(fmt::layer().json().with_current_span(true))
        .init();
}
```

### 4. Hubungkan Tracing ke Jaeger via OpenTelemetry Layer

Ekspor trace terstruktur ke collector OpenTelemetry untuk visualisasi rantai pemanggilan microservices.

```rust
use tracing_opentelemetry::OpenTelemetryLayer;
use tracing_subscriber::layer::SubscriberExt;

fn init_tracer() {
    let tracer = opentelemetry_jaeger::new_agent_pipeline()
        .with_service_name("user-service")
        .install_simple()
        .unwrap();

    let telemetry = OpenTelemetryLayer::new(tracer);
    tracing_subscriber::registry().with(telemetry).init();
}
```

### 5. Evaluasi Implikasi Performa: Tracing vs Log Crate

Crate `log` hanya memformat string secara sekuensial. `tracing` menambahkan abstraksi span, namun mengoptimalkannya dengan menonaktifkan pembuatan string jika level log tidak memenuhi kriteria penyaringan.

```rust
/* Crate log tradisional: pembuatan string terjadi sebelum filter jika tidak hati-hati */
/* log::debug!("User: {}", complex_formatting_function()); */

/* Crate tracing: argumen dinilai secara lazy oleh subscriber */
tracing::debug!(user = %complex_formatting_function(), "Detail pengguna");
```
