---
title: "Evolusi Runtime Kontainer: Membedah Komponen Internal runc dan containerd"
slug: 0e1f2a3b
aliases: []
date: 2026-08-11
tags: [docker]
lang: id
draft: false
type: post
cover: ""
---

1. PENGANTAR: DEKONSTRUKSI MONOLITIK RUNTIME KONTAINER
Evolusi teknologi virtualisasi tingkat sistem operasi telah bertransformasi pesat sejak era awal kepopuleran Docker. Docker Engine pada awalnya direkayasa sebagai sistem monolitik. Semua kapabilitas yang berpusat pada penarikan image (image pulling), penyimpanan jaringan, eksekusi proses, dan pemantauan digabungkan dalam sebuah daemon tunggal. Seiring dengan peningkatan skala penggunaan di lingkungan produksi, desain arsitektur yang terpusat ini memunculkan keterbatasan dalam hal stabilitas, auditabilitas, dan isolasi kegagalan.

Sebagai respons terhadap tantangan ini, komunitas terbuka memelopori dekomposisi daemon monolitik menjadi tumpukan (stack) komponen terdistribusi yang lebih spesifik. Open Container Initiative (OCI) kemudian dibentuk untuk membakukan spesifikasi format citra (image format) dan spesifikasi waktu berjalan (runtime specification). Artikel ini mendiskusikan mekanisme internal dari dua entitas utama yang terlahir dari pemisahan ini, yaitu runc sebagai low-level runtime dan containerd sebagai high-level runtime.

2. RUNC SEBAGAI FONDASI EKSEKUSI TINGKAT RENDAH
Secara konseptual, runc adalah implementasi referensi resmi dari spesifikasi runtime OCI. Peran utamanya sangat spesifik. runc bertanggung jawab murni untuk meluncurkan proses kontainer yang dikonfigurasi berdasarkan parameter lingkungan dan instruksi filesystem yang diberikan kepadanya. runc tidak mengerti mengenai pemanggilan API, transfer protokol HTTP dari registry, atau manajemen metadata image.

2.1. Orkestrasi Namespace dan Cgroups
Ketika diinstruksikan untuk menjalankan sebuah kontainer, runc berinteraksi langsung dengan primitif kernel Linux. runc menginisialisasi namespaces (seperti PID, Mount, Network, UTS, dan IPC) untuk menciptakan pandangan realitas sistem operasi terisolasi untuk proses anak (child process). Secara simultan, runc berkoordinasi dengan utilitas Control Groups (cgroups) pada sistem operasi (khususnya cgroupfs atau systemd init) untuk mengonfigurasi batasan agregat penggunaan komputasi CPU, alokasi memori acak, dan pembatasan bandwidth I/O. 

2.2. Sistem Transien Tanpa Status
Karakteristik desain yang paling membedakan runc adalah sifatnya yang ephemeric. Setelah kontainer diluncurkan secara sukses, dan terminal utama diserahkan pada proses yang diinisiasi, proses runc itu sendiri berakhir dan keluar dari memori eksekusi. Oleh sebab itu, runc tidak berjalan sebagai service daemon yang selalu aktif di latar belakang (background), menjadikannya sangat efisien dan tahan banting terhadap kegagalan komponen sistem lainnya.

3. CONTAINERD: MANAJEMEN SIKLUS HIDUP KONTAINER TINGKAT TINGGI
Mengingat runc beroperasi pada abstraksi kernel yang sangat fundamental, sebuah lapisan orkestrasi perantara diperlukan untuk menangani kerumitan infrastruktur yang tidak relevan bagi sistem proses dasar. Komponen perantara inilah yang diperankan secara gemilang oleh containerd.

Containerd merupakan manajer runtime tingkat tinggi. Entitas ini didesain sebagai sebuah daemon persisten. Berbeda dengan runc, peran containerd mencakup manajemen metadata keseluruhan objek kontainer dalam lingkungan host.

3.1. Penanganan Citra (Image Management)
Saat sistem meminta penyebaran beban kerja (workload) baru, containerd menangani pemrosesan awal. Containerd mengunduh image dari repositori OCI (contohnya Docker Hub atau Google Container Registry), kemudian mengekstrak lapisan arsip tarball image tersebut (unpacking), lalu merangkai sistem file akar sementara menggunakan overlay filesystem, misalnya OverlayFS atau Btrfs. 

3.2. Shim Architecture
Untuk memantau umur panjang proses yang dijalankan oleh utilitas sementara seperti runc, containerd memanfaatkan lapisan pembantu berukuran mikroskopis yang disebut containerd-shim. Setiap kontainer yang beroperasi memiliki proses shim terdedikasi miliknya sendiri. Shim ini secara konstan menyalurkan output standard (stdout dan stderr) ke subsistem logging host. Penting untuk dicatat, eksistensi shim memungkinkan pemeliharaan daemon containerd untuk dihentikan paksa (restart atau upgrade) tanpa membunuh atau menjatuhkan proses kontainer yang sedang melayani pengguna akhir. Hal ini memungkinkan apa yang disebut pembaruan tanpa waktu henti (zero-downtime updates) untuk manajer kontainer.

4. KESIMPULAN ARSITEKTURAL
Perpisahan ekologis fungsional antara utilitas eksekusi tingkat rendah dan layanan manajerial tingkat tinggi menyederhanakan siklus pengembangan inovasi komputasi awan. Dengan mendelegasikan tugas-tugas pengelolaan sumber daya memori dan peluncuran kernel ke runc, serta sentralisasi administrasi penyimpanan dan jaringan pada containerd, kerangka arsitektur kontainer masa kini menunjukkan ketangguhan, prediktabilitas, dan keamanan yang memadai. Integrasi mendalam dengan orchestrator semacam Kubernetes pun menjadi sangat efisien, yang secara langsung berinteraksi dengan antarmuka Container Runtime Interface (CRI) milik containerd dalam memuluskan operasi skala raksasa.
[^1][^2]

## Referensi

[^1]: Open Container Initiative. "OCI Runtime Specification." OCI Working Group, 2020.
[^2]: Crosby, S., et al. "Understanding containerd: Industry-standard container runtime." CNCF Publications, 2019.