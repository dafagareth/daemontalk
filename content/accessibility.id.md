# Pernyataan Aksesibilitas

**Terakhir diperbarui: 4 September 2026** · **Standar: WCAG 2.1 Level AA**

Daemontalk berkomitmen untuk menghadirkan pengalaman membaca dan eksplorasi teknologi yang inklusif, cepat, dan mudah diakses oleh semua orang, termasuk pengguna dengan disabilitas sensorik, motorik, kognitif, maupun pengguna di lingkungan konsol dan *screen reader*.

---

## Prinsip Desain Inklusif

Platform ini dirancang dengan prinsip *universal accessibility* yang mengacu pada standar internasional Web Content Accessibility Guidelines (**WCAG 2.1 Level AA**). Kami memastikan bahwa kemudahan navigasi, keterbacaan tipografi, dan kontras visual menjadi prioritas utama tanpa mengorbankan performa sistem.

## Navigasi Penuh Menggunakan Keyboard

Seluruh fitur dan elemen interaktif di Daemontalk dapat dioperasikan secara penuh tanpa menggunakan mouse atau perangkat penunjuk:

**Pintasan Global (`?`)**: Menampilkan modal daftar seluruh pintasan keyboard yang tersedia di situs.

**Pencarian Instan (`/`)**: Mengarahkan fokus kursor secara langsung ke kotak pencarian artikel.

**Ganti Tema Warna (`t`)**: Mengubah mode tema visual antara Terang (*Light*) dan Gelap (*Dark*).

**Navigasi Artikel (`j` dan `k`)**: Berpindah ke artikel sebelum atau sesudah secara berurutan.

**Menutup Jendela (`Esc`)**: Menutup modal aktif, menu pencarian, atau *drawer* navigasi.

**Navigasi Sekuensial (`Tab` & `Shift+Tab`)**: Berpindah antar tautan, tombol, dan elemen formulir dengan indikator cincin fokus (*focus ring*) berkontras tinggi yang jelas terlihat.

## Tipografi Ergonomis & Kontras Visual

**Kombinasi Fon Terbaca Tinggi**: Teks artikel utama menggunakan fon **Source Serif 4** dengan ritme spasi yang proporsional untuk kenyamanan membaca esai panjang, antarmuka menggunakan **Plus Jakarta Sans**, dan blok kode menggunakan **JetBrains Mono**.

**Skalabilitas Ukuran Fon**: Setiap artikel dilengkapi tombol kontrol ukuran teks (`A-` dan `A+`) yang mengubah ukuran teks secara vertikal tanpa memotong kata (*clipping*) atau merusak tata letak responsif.

**Rasio Kontras Terkalibrasi**: Seluruh kombinasi warna teks terhadap latar belakang memenuhi ambang batas kontras minimal 4.5:1 untuk teks normal dan 3:1 untuk teks berukuran besar, baik pada mode terang maupun gelap.

**Filter Layar Hangat (*Warm Screen Tint*)**: Slider kehangatan warna layar terintegrasi untuk mengurangi paparan cahaya biru (*blue light*) dan kelelahan mata saat membaca di lingkungan minim cahaya.

## Struktur Semantik & Kompatibilitas Pembaca Layar

**Elemen Semantik HTML5**: Struktur halaman dibangun menggunakan tag semantik murni (`<main>`, `<nav>`, `<article>`, `<header>`, `<footer>`, `<aside>`, `<section>`) yang memungkinkan perangkat lunak *screen reader* (NVDA, VoiceOver, JAWS, Orca) menavigasi struktur dokumen dengan presisi.

**Atribut ARIA yang Jelas**: Seluruh tombol interaktif, ikon, status pemuatan, dan modal dialog dilengkapi atribut ARIA yang tepat (`aria-label`, `aria-expanded`, `aria-hidden`, `role="dialog"`).

**Teks Alternatif Deskriptif**: Gambar diagram arsitektur, grafik performa, dan bagan alur sistem dilengkapi atribut `alt` deskriptif yang menjelaskan isi gambar secara kontekstual.

## Pengurangan Gerakan & Kenyamanan Sensorik

**Dukungan Reduced Motion**: Situs ini secara otomatis mendeteksi preferensi sistem `prefers-reduced-motion` dan mematikan animasi transisi non-esensial bagi pengguna dengan sensitivitas vestibular atau gangguan gerak.

**Bebas Elemen Mengganggu**: Tidak ada konten yang berkedip cepat (*zero flashing content*), tidak ada audio yang diputar otomatis (*no autoplay*), dan tidak ada *pop-up* iklan yang menghalangi fokus membaca.

## Aksesibilitas Terminal & Mode Teks Murni

Bagi pengguna yang bekerja di lingkungan konsol murni, terminal *Braille*, atau perangkat dengan keterbatasan grafis, seluruh arsip publikasi dan catatan teknis dapat diakses langsung melalui protokol SSH tanpa memerlukan peramban grafis:

```bash
$ ssh ssh.daemontalk.com -p 2222
```

## Umpan Balik & Saluran Bantuan

Kami terus memantau dan menguji kepatuhan aksesibilitas situs ini secara berkala. Jika Anda menemukan kendala aksesibilitas, elemen yang sulit dibaca, atau navigasi keyboard yang terhambat, silakan hubungi kami langsung via email di: **realdaemontalk@gmail.com**. Setiap laporan aksesibilitas akan kami tangani sebagai prioritas utama.
