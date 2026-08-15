### v1.0.0 (Agustus 2026) · Rilis Perdana

*   **Arsitektur Inti**: Berkas biner tunggal mandiri berbasis Go menggunakan router `chi`, *server-side rendering* `a-h/templ`, dan Tailwind CSS v4. Tanpa framework JavaScript berat, tanpa pelacak.
*   **Dukungan Multibahasa**: Rute terlokalisasi penuh untuk Bahasa Indonesia (`/id`) dan English dengan antarmuka serta metadata yang adaptif.
*   **Terminal & CLI Interaktif**: Emulator shell UNIX virtual di browser (`/terminal`) dengan riwayat perintah, serta endpoint ramah `curl` (`/daily`, `/recipes`, `/p/:slug`).
*   **Komentar Anonim & Buku Tamu**: Penyimpanan SQLite ringan dengan identitas pengunjung deterministik yang melekat pada cookie (`anonym_<hex>`) dan pengelompokan pesan berturut-turut.
*   **Ekstensi Markdown Editorial**: Penyorotan sintaks kode di server, carousel responsif (` ```carousel `), galeri (` ```gallery `), akordeon FAQ (` ```faq `), kartu profil penulis (` ```author `), catatan kaki, dan daftar isi otomatis.
*   **Pengalaman Membaca**: Pengaturan ukuran font (A+/A-), mode Serif, pemilih palet tema (Light, Sepia, Dark), markah bacaan (*reading list*), dan pencarian artikel *in-memory* yang cepat.
*   **Feeds & SEO**: Generator otomatis kartu pratinjau OpenGraph, RSS 2.0 feed, JSON Feed, Sitemap XML, dan *Content Security Policy* (CSP) ketat.
