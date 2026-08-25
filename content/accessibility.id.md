# Pernyataan Aksesibilitas

**Terakhir diperbarui: Agustus 2026**

Daemontalk dibangun agar dapat diakses dengan cepat, nyaman, dan inklusif di berbagai perangkat, *screen reader*, maupun terminal teks murni, dengan mengacu pada panduan **WCAG 2.1 Level AA**.

## Navigasi Keyboard & Pintasan

Seluruh elemen interaktif dapat dioperasikan penuh menggunakan keyboard dengan indikator fokus yang tegas:

- `?` : Membuka daftar pintasan keyboard.
- `/` : Fokus langsung ke kotak pencarian.
- `t` : Mengganti tema tampilan (Light, Dark, Sepia).
- `j` / `k` : Berpindah ke artikel sebelum atau sesudah.
- `Esc` : Menutup menu, pop-up, atau jendela pencarian.
- `Tab` : Berpindah antar tautan dan tombol dengan garis fokus kontras.

## Tipografi & Kontras Teks

- **Pilihan Font Ergonomis**: Isi artikel menggunakan **Source Serif 4** untuk kenyamanan membaca esai panjang, antarmuka menggunakan **Plus Jakarta Sans**, dan blok kode menggunakan **JetBrains Mono**.
- **Skalabilitas Teks**: Tersedia tombol pengatur ukuran font (`A-` dan `A+`) pada setiap artikel tanpa merusak tata letak halaman.
- **Rasio Kontras Tinggi**: Kontras teks terhadap latar belakang selalu di atas 4.5:1 pada mode terang maupun gelap.
- **Filter Layar Hangat (Warm Tint)**: Slider kehangatan warna layar untuk mengurangi ketegangan mata di ruangan minim cahaya.

## Pembaca Layar (Screen Reader) & Struktur Semantik

- Menggunakan elemen semantik HTML5 murni (`<main>`, `<nav>`, `<article>`, `<header>`, `<footer>`, `<aside>`).
- Seluruh tombol dan ikon interaktif dilengkapi atribut ARIA (`aria-label`, `aria-hidden`).
- Gambar dan diagram sistem dilengkapi teks alternatif (`alt`) deskriptif.
- Menghormati pengaturan sistem `prefers-reduced-motion` untuk mematikan animasi transisi.

## Akses Terminal & Mode Teks Murni

Bagi pengguna yang menggunakan *Braille display*, lingkungan konsol, atau terminal tanpa antarmuka grafis, seluruh artikel dapat dibaca langsung melalui SSH:

```bash
$ ssh ssh.daemontalk.com -p 2222
```

## Bantuan & Umpan Balik

Jika Anda menemukan kendala aksesibilitas saat mengakses website ini, silakan hubungi kami melalui email di **realdaemontalk@gmail.com**. Laporan aksesibilitas akan kami tangani sebagai prioritas utama.
