---
title: "USB-C-nya AI Tools"
slug: "mcp-protocol-ai-tools"
aliases: []
author: "daemontalk team"
contributors: []
tags: ["radar", "ai", "agents", "llm"]
lang: "id"
status: "published"
type: "dispatch"
readTime: 4
cover: "https://images.unsplash.com/photo-1518770660439-4636190af475?auto=format&fit=crop&w=1200&q=80"
coverCaption: "Satu protokol untuk semua integrasi"
coverSource: "https://unsplash.com"
description: "Model Context Protocol (MCP) adalah standar terbuka agar LLM bisa berinteraksi dengan tools dan data eksternal tanpa integrasi custom per platform."
series: ""
series_part: 0
---

<span class="drop-cap">S</span>etiap penyedia model AI selama ini memiliki mekanisme pemanggilan fungsi (tool calling) dengan aturan dan skema JSON tersendiri.

Jika Anda mengelola backend dan ingin mengekspos fitur agar bisa dipakai agen AI, Anda dihadapkan pada masalah klasik: rute khusus untuk OpenAI, wrapper tersendiri untuk Claude, serta penyesuaian untuk framework seperti LangChain.

Ini adalah masalah integrasi $N \times M$. Menjaga kompatibilitas antara $N$ model AI dan $M$ layanan backend menguras energi engineering hanya untuk mengurus format payload yang berbeda-beda.

Model Context Protocol (MCP), standar terbuka yang diinisiasi oleh Anthropic, hadir untuk menyelesaikan fragmentasi ini. Singkatnya, MCP adalah USB-C untuk perkakas AI.

## Arsitektur Berbasis JSON-RPC 2.0

Alih-alih merancang protokol biner baru yang rumit, spesifikasi MCP memilih fondasi yang sudah matang: JSON-RPC 2.0.

Dalam arsitektur MCP, aplikasi antarmuka seperti IDE atau Claude Desktop bertindak sebagai Client. Layanan backend atau basis data Anda bertindak sebagai Server.

Keduanya berkomunikasi melalui dua jenis transport:
- `stdio` untuk menjalankan CLI lokal langsung dari mesin pengguna.
- HTTP dengan Server-Sent Events (SSE) untuk layanan jarak jauh di server cloud.

Pemisahan transport ini memberi fleksibilitas tinggi. Server MCP bisa berupa skrip Python sederhana di laptop, atau microservice Go berkemampuan tinggi di dalam klaster Kubernetes.

## Tiga Primitif Utama

Server MCP dapat mengekspos tiga komponen utama kepada model AI:

1. **Resources**: Data read-only untuk konteks percakapan. Resources menggunakan pengalamatan URI seperti `postgres://db/schema/users` sehingga agen bisa membaca struktur data sistem Anda secara langsung.
2. **Tools**: Fungsi executable yang membawa perubahan status (side effects). Anda menentukan nama fungsi, deskripsi semantik kegunaannya, dan skema parameter input. Model AI akan memutuskan kapan harus memanggil fungsi ini.
3. **Prompts**: Template percakapan terstruktur yang disiapkan oleh server, sehingga pengguna tidak perlu mengetik instruksi panjang secara manual.

## Apa Artinya bagi Backend Engineer?

Bagi pengembang backend, MCP mengubah cara kita membangun integrasi AI.

Anda tidak perlu lagi menulis adapter khusus untuk setiap vendor model. Anda cukup menulis satu server MCP di depan basis data, agregator log, atau API pembayaran internal.

Ketika aplikasi klien terhubung, model AI membaca skema tools yang tersedia secara dinamis. Jika tahun depan tim Anda memutuskan berganti penyedia LLM, server MCP Anda tetap bekerja tanpa perlu diubah satu baris pun.

HTTP menyatukan pertukaran dokumen web, TCP/IP menyatukan jaringan, dan USB menyatukan kabel perangkat keras. MCP membawa pendekatan serupa untuk menghubungkan model AI non-deterministik dengan sistem deterministik yang kita rawat setiap hari.

```references
- title: "Model Context Protocol Specification"
  author: "Anthropic"
  year: 2024
  url: "https://modelcontextprotocol.io/"
```
