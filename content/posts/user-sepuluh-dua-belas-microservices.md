---
title: "User Baru Sepuluh, Sudah Bikin Dua Belas Microservices"
slug: "user-sepuluh-dua-belas-microservices"
aliases: []
author: "daemontalk team"
contributors: []
tags: ["essays", "opinion", "industry"]
lang: "id"
status: "published"
type: "dispatch"
readTime: 5
cover: "https://images.unsplash.com/photo-1558494949-ef010cbdcc31?auto=format&fit=crop&w=1200&q=80"
coverCaption: "Infrastruktur terdistribusi rumit untuk aplikasi yang belum punya trafik"
coverSource: "https://unsplash.com"
description: "Meniru arsitektur terdistribusi Netflix di startup tahap awal adalah resep bunuh diri operasional. Microservices adalah solusi masalah skala organisasi, bukan medali gengsi."
series: ""
series_part: 0
---

<span class="drop-cap">A</span>da pemandangan ganjil yang sering terjadi di industri startup teknologi hari ini.

Perusahaannya baru berumur tiga bulan, produknya baru diuji oleh lingkaran teman kantor pendiri, dan jumlah pengguna aktif hariannya mungkin belum genap dua puluh orang.

Namun jika Anda membuka dasbor penyedia cloud mereka, Anda akan menemukan klaster Kubernetes megah dengan dua belas microservices: auth service, payment service, notification service, catalog service, sampai analytics service.

Setiap kali ingin merilis perubahan kecil, tiga orang developer harus menyelaraskan kontrak API, alur deployment memakan waktu setengah jam, dan tagihan server bulanan membengkak sebelum ada sepeser pun pemasukan.

## Meniru Solusi untuk Masalah yang Belum Ada

Raksasa teknologi seperti Netflix, Amazon, atau Uber memecah sistem mereka menjadi ribuan layanan kecil bukan karena arsitekturnya terlihat keren di presentasi konferensi.

Mereka memecah sistem karena memiliki ribuan insinyur. Jika ribuan orang harus mengedit dan merilis satu repositori yang sama setiap hari, antrean koordinasi dan potensi konflik kode akan melumpuhkan seluruh perusahaan.

Microservices pada dasarnya adalah solusi untuk masalah skala organisasi, bukan semata-mata masalah keunggulan performa teknis.

Ketika tim Anda hanya terdiri dari empat orang, menerapkan arsitektur microservices sama saja dengan membebankan kerumitan koordinasi korporasi raksasa ke atas pundak segelintir orang.

## Pajak Jaringan yang Menguras Tenaga

Di dalam arsitektur monolit, komunikasi antarfungsi terjadi di dalam memori RAM yang sama dan selesai dalam hitungan nanodetik.

Begitu fungsi tersebut Anda pisahkan ke server berbeda sebagai microservice, setiap panggilan antarkomponen berubah menjadi panggilan jaringan melalui protokol HTTP atau gRPC.

Jaringan internet tidak pernah seratus persen dapat diandalkan. Anda mendadak harus menambahkan circuit breaker, mekanisme retry, tracing terdistribusi, dan pusing memikirkan konsistensi data transaksi antar-database yang terpisah.

Komputer server Anda kini menghabiskan lebih banyak waktu untuk serialisasi JSON dan jabat tangan TLS daripada menjalankan logika bisnis yang sebenarnya dibutuhkan pengguna.

## Nasihat Monolith First dan Pelajaran Prime Video

Dalam bukunya mengenai arsitektur software, Martin Fowler merangkum aturan dasar ini dalam satu kalimat lugas:

> "The first rule of distributed systems is don't distribute your system."
> (Martin Fowler)

Simon Brown, pencipta model arsitektur C4, mempertegas peringatan tersebut dengan pertanyaan sederhana: jika sebuah tim belum mampu membangun monolit yang rapi, apa yang membuat mereka yakin bahwa memecahnya menjadi lusinan bagian terpisah akan menyelesaikan masalah?

Sistem yang dibangun langsung sebagai kumpulan microservices sejak hari pertama hampir selalu menemui jalan buntu operasional.

Bahkan tim teknik Amazon Prime Video memberikan contoh nyata pada tahun 2023. Mereka mendokumentasikan migrasi layanan pemantauan audio-video dari arsitektur serverless terdistribusi kembali ke satu monolit terkonsolidasi, sebuah langkah yang sukses memangkas biaya komputasi mereka hingga sembilan puluh persen sekaligus menyederhanakan pengelolaan sistem.

Bangunlah monolit yang rapi terlebih dahulu. Jika suatu hari nanti bisnis Anda beruntung bisa tumbuh hingga jutaan pengguna dan database Anda mulai kewalahan, itu adalah masalah menyenangkan yang layak dirayakan.

Sebelum hari itu tiba, hemat waktu dan energi Anda untuk mencari pengguna pertama, bukan sibuk mengoleksi kontainer di awan.

```references
- title: "MonolithFirst"
  author: "Martin Fowler"
  year: 2015
  publisher: "martinfowler.com"
  url: "https://martinfowler.com/bliki/MonolithFirst.html"

- title: "Scaling up the Prime Video audio/video monitoring service and reducing costs by 90%"
  author: "Prime Video Tech Team"
  year: 2023
  publisher: "Amazon Prime Video Tech"
  url: "https://www.primevideotech.com/video-streaming/scaling-up-the-prime-video-audio-video-monitoring-service-and-reducing-costs-by-90"
```
