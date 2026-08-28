---
title: "Mengenal Retrieval-Augmented Generation (RAG)"
slug: pengenalan-rag-ai
aliases: []
date: 2026-08-28
author: "Daemontalk Editorial"
tags: ["AI", "Machine Learning"]
lang: id
draft: false
type: post
cover: "https://images.unsplash.com/photo-1677442136019-21780ecad995?auto=format&fit=crop&q=80&w=1600"
coverCaption: "Abstract digital neural network and nodes"
coverSource: "https://unsplash.com"
readTime: 6
description: "RAG adalah solusi atas halusinasi LLM. Pelajari konsep chunking, vector embeddings, dan arsitektur RAG yang menyuntikkan fakta ke dalam prompt AI."
---

Model Bahasa Besar (LLM) seperti GPT-4 atau Llama 3 memang mampu merangkai kalimat dengan sangat baik. Namun, mereka memiliki titik buta yang kritis: Pengetahuan mereka hanya sebatas data latih (*training data*) masa lalu, dan mereka sama sekali tidak mengetahui data privat perusahaan Anda.

Jika Anda bertanya soal laporan keuangan perusahaan minggu lalu, LLM murni akan menebak-nebak dan mengarang bebas (halusinasi). 

Di sinilah **Retrieval-Augmented Generation (RAG)** menjadi pahlawan.

## Apa itu RAG?

RAG adalah teknik arsitektur yang menjembatani AI generatif dengan pangkalan data eksternal Anda. Secara harfiah:
- **Retrieval**: Mencari dokumen relevan dari *database*.
- **Augmented**: Memperkaya *prompt* dasar dengan dokumen yang baru saja ditemukan.
- **Generation**: LLM menghasilkan jawaban akhir berdasarkan panduan dokumen spesifik tersebut.

```text
+-----------+      +-------------------+      +-------------+
| User      | ---> | 1. Vector DB      | ---> | Context     |
| Query     |      |    (Retrieval)    |      | (Documents) |
+-----------+      +-------------------+      +-------------+
                                                    |
                                                    v
+-----------+      +-------------------+      +-------------+
| Final     | <--- | 3. LLM Model      | <--- | 2. Prompt + |
| Response  |      |    (Generation)   |      | Context     |
+-----------+      +-------------------+      +-------------+
```

## Anatomi Persiapan Data: Chunking dan Embeddings

Sebelum sistem RAG bisa melayani pengguna, data Anda harus disiapkan (proses *Ingestion*):

1. **Chunking**: Dokumen besar (seperti PDF 500 halaman) dipotong-potong menjadi paragraf-paragraf kecil (misalnya 500 token per *chunk*). Ini karena LLM memiliki batas *context window* dan tidak bisa membaca semua PDF sekaligus.
2. **Vector Embeddings**: Setiap *chunk* teks dikonversi menjadi deretan angka (vektor multidimensi) menggunakan model *embedding* (seperti `text-embedding-ada-002`). Model ini menangkap "makna semantik" dari teks tersebut.
3. **Penyimpanan**: Vektor-vektor ini disimpan ke dalam *Vector Database* (seperti Pinecone, Weaviate, Qdrant, atau ekstensi `pgvector` pada PostgreSQL).

## Proses Eksekusi Waktu Nyata

Saat pengguna bertanya "Bagaimana prosedur cuti tahunan?":
1. Pertanyaan pengguna diubah menjadi vektor.
2. Terjadi pencarian semantik (*Semantic Search*). Vektor pertanyaan diukur kemiripan arah matematisnya (*cosine similarity*) dengan vektor di *database*.
3. Sistem mengambil 3 *chunk* paling mirip yang membahas kebijakan cuti perusahaan.
4. Sistem mengirim instruksi ke LLM: *"Berdasarkan konteks kebijakan cuti di bawah ini, jawab pertanyaan pengguna. Konteks: [Isi Chunk]. Pertanyaan: Bagaimana prosedur cuti?"*

> [!IMPORTANT]
> Kelemahan RAG terletak pada proses pencariannya. Jika algoritma pencarian gagal menemukan dokumen yang tepat, secanggih apa pun LLM-nya, jawaban yang dihasilkan tetap akan melenceng (*Garbage in, Garbage out*).

## Referensi Lanjutan

```references
- title: "Retrieval-Augmented Generation for Knowledge-Intensive NLP Tasks"
  author: "Patrick Lewis et al."
  year: 2020
  publisher: "NeurIPS"
  url: "https://arxiv.org/abs/2005.11401"
```
