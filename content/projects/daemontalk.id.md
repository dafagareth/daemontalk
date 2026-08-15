## Ikhtisar

Situs ini (**daemontalk.com**) adalah buku catatan rekayasa dan portofolio teknik pribadi yang dibangun menggunakan Go. Tanpa framework JavaScript kompleks atau pipeline build yang rumit. Mengandalkan rendering HTML di sisi server (SSR), sedikit HTMX untuk fitur interaktif, dan SQLite untuk komentar serta analitik.

## Arsitektur

Situs ini berupa satu berkas biner Go yang melayani:

- **Halaman Statis**: Beranda, tentang, alat/uses, now, resume, changelog, dan tautan
- **Blog**: Tulisan Markdown dengan frontmatter, penyorotan sintaks, daftar isi, dan navigasi seri
- **Proyek**: Daftar etalase proyek dengan halaman detail
- **Buku Tamu**: Papan pesan komunitas
- **TIL**: Catatan ringkas pembelajaran teknis
- **Admin**: Dasbor ringan untuk moderasi dan analitik dasar

## Tumpukan Teknologi

- **Go** dengan pustaka standar dan router Chi
- **templ** untuk pembuatan templat HTML di sisi server yang aman secara tipe data
- **HTMX** untuk interaksi dinamis tanpa berkas JS besar di klien
- **Tailwind CSS v4** untuk penataan gaya utilitas yang bersih
- **goldmark** untuk rendering Markdown dengan Chroma syntax highlighter
- **modernc.org/sqlite** (driver SQLite Go murni) untuk persistensi data
- Diterapkan pada satu VPS Linux mandiri dengan layanan systemd

## Tujuan Desain

- **Bebas ketergantungan framework JavaScript**: Semua halaman tetap dapat dibaca dan bekerja dengan baik tanpa JavaScript.
- **Waktu mulai dingin cepat**: Semua artikel diindeks ke dalam memori saat aplikasi dimulai untuk kecepatan akses instan.
- **Operasional sederhana**: Satu berkas biner, satu berkas SQLite, dan satu layanan systemd.
- **Multibahasa**: Mendukung penuh Bahasa Indonesia (`/id/`) dan Bahasa Inggris (`/`).
