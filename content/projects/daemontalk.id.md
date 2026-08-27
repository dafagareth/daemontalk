## Ikhtisar

Situs ini (**daemontalk.com**) adalah publikasi rekayasa perangkat lunak independen dan basis pengetahuan sistem yang dibangun sepenuhnya menggunakan Go. Proyek ini secara sadar menghindari *framework* JavaScript yang kompleks dan proses *build* yang berat, lebih memilih HTML yang dirender di sisi server, interaksi presisi melalui HTMX, dan basis data SQLite yang sangat ringan.

## Arsitektur

Keseluruhan sistem dikompilasi menjadi satu berkas biner Go mandiri yang melayani:

- **Halaman Statis**: Rute inti untuk beranda, tentang, *uses*, *now*, resume, *changelog*, dan tautan.
- **Mesin Blog**: Pemroses Markdown dengan metadata *frontmatter*, penyorotan sintaks (*syntax highlighting*) Chroma di sisi server, dan daftar isi dinamis.
- **Etalase Proyek**: Portofolio proyek terperinci berbasis Markdown.
- **Buku Tamu Komunitas**: Papan pesan persisten menggunakan identitas anonim yang terikat pada *cookie*.
- **Admin Studio**: Dasbor aman untuk moderasi, penyusunan draf konten, dan analitik telemetri.

## Tumpukan Teknologi

- **Go 1.26+** menggunakan pustaka standar dan `go-chi/chi` untuk perutean.
- **a-h/templ** untuk pembuatan HTML sisi server yang aman terhadap tipe (*type-safe*).
- **HTMX** untuk interaksi dinamis dan *real-time* (seperti pencarian instan) tanpa beban JavaScript.
- **Tailwind CSS v4** untuk penataan gaya utilitas yang ketat.
- **goldmark** untuk pemrosesan Markdown yang kokoh beserta ekstensinya.
- **modernc.org/sqlite** (driver SQLite murni Go) untuk penyimpanan data persisten.
- Disebarkan pada VPS Debian ringan yang dikelola oleh `systemd` dan `Caddy`.

## Prinsip Desain

- **Nirketergantungan JS**: Fungsi inti membaca dan navigasi harus berjalan sempurna meskipun JavaScript dimatikan.
- **Waktu Mulai Instan**: Seluruh artikel dan metadata diindeks ke dalam memori saat aplikasi menyala demi waktu respons di bawah satu milidetik.
- **Kesederhanaan Operasional**: Satu berkas biner, satu basis data SQLite, satu layanan systemd. Tidak memerlukan lapisan *cache* atau basis data eksternal.
- **Lokalisasi Bawaan**: Dukungan perutean dwibahasa penuh untuk Bahasa Inggris (`/`) dan Bahasa Indonesia (`/id/`).

## Kemampuan Utama

- Kontrol tipografi granular (Tema Gelap/Terang, peralihan Serif/Sans, penskalaan ukuran huruf dinamis).
- Pembuatan otomatis kartu pratinjau OpenGraph untuk keperluan berbagi di media sosial.
- Sindikasi komprehensif melalui RSS 2.0 dan JSON Feed.
- Integrasi SEO esensial (Sitemap XML, `robots.txt`, skema terstruktur).
- Pembatasan laju (*rate limiting*) dan mekanisme *honeypot* untuk pencegahan spam.
- Penjadwalan publikasi menggunakan metadata *frontmatter* `publish_at`.
