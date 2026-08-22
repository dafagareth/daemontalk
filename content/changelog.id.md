### v1.0.3 (Agustus 23, 2026) · Perombakan Sidebar Ponsel & Pembersihan UI

- **Desain Ulang Menu Ponsel**: Merombak total antarmuka *drawer* navigasi pada perangkat *mobile*. Memindahkan fitur pencarian ke atas sebagai bilah *input* khusus, mengangkat daftar topik (*Topic Streams*), dan menggunakan *font sans-serif* murni agar lebih ramah sentuh.
- **Penghapusan Jam Live**: Mencopot elemen jam *real-time* yang berlebihan dari *header* desktop dan ponsel untuk menjaga tampilan navigasi yang lebih bersih dan bebas distraksi.
- **Ekosistem Sosial**: Menambahkan deretan ikon media sosial khusus (rata kiri) di bagian paling bawah menu ponsel. Menyertakan ikon Facebook baru dan mengganti ikon Threads menjadi format garis luar (*outline*) SVG agar selaras secara estetika. Menyelaraskan seluruh tautan sosial (*handle*) agar langsung mengarah ke `daemontalk`.
- **Standarisasi Dokumentasi**: Membersihkan markup yang berantakan (*AI slop*) di seluruh file dokumentasi `.md`. Mengubah blok HTML kotor menjadi sintaks *native* ` ```faq ` milik proyek, serta membasmi penomoran *heading* yang kaku pada halaman legal dan pedoman kontribusi.

### v1.0.2 (Agustus 22, 2026) · Pembaruan Jalur CI/CD

- **Sinkronisasi Deployment**: Memperbaiki *race condition* pada GitHub Actions di mana VPS menarik *Docker image* lawas sebelum *image* terbaru selesai diproses oleh GHCR. *Deployment* kini diatur agar wajib menunggu hingga proses *build* selesai sepenuhnya.

### v1.0.1 (Agustus 21, 2026) · Peningkatan UX & Infrastruktur

- **Pencarian Instan**: Pembaruan mesin pencari menggunakan HTMX untuk menampilkan saran *real-time* saat mengetik. Memperbaiki *bug highlight* kuning agresif pada kueri satu huruf.
- **Penyimpanan Otomatis**: Editor Admin Markdown kini terintegrasi dengan fungsi *auto-save* ke `localStorage` untuk mencegah draf hilang secara tak sengaja.
- **Struktur Tipografi**: Memaksa penggunaan huruf sans-serif khusus untuk judul (*headings*) guna menjaga ketegasan visual saat pengguna mengaktifkan mode Serif.
- **Perbaikan CSS**: Mengatasi *bug* pewarisan (*inheritance*) di mana pemilih global `[id]` tanpa sengaja memaksa seluruh isi artikel menjadi huruf sans-serif dan merusak fungsi mode baca.
- **Standarisasi WWW**: Domain utama (`daemontalk.com`) kini dialihkan sepenuhnya (*301 redirect*) ke `www.daemontalk.com` melalui konfigurasi Caddy.
- **Dokumentasi Diperbarui**: Menyederhanakan instruksi `README.md` dan `VPS_SETUP.md` serta menyertakan bagan arsitektur vektor (SVG) murni yang dapat diskalakan.

### v1.0.0 (Agustus 2026) · Rilis Perdana

- **Arsitektur Inti**: Berkas biner tunggal mandiri berbasis Go menggunakan router `chi`, *server-side rendering* `a-h/templ`, dan Tailwind CSS v4. Tanpa framework JavaScript berat, tanpa pelacak.
- **Dukungan Multibahasa**: Rute terlokalisasi penuh untuk Bahasa Indonesia (`/id`) dan English dengan antarmuka serta metadata yang adaptif.
- **Terminal & CLI Interaktif**: Emulator shell UNIX virtual di browser (`/terminal`) dengan riwayat perintah, serta endpoint ramah `curl` (`/daily`, `/recipes`, `/p/:slug`).
- **Komentar Anonim & Buku Tamu**: Penyimpanan SQLite ringan dengan identitas pengunjung deterministik yang melekat pada cookie (`anonym_<hex>`) dan pengelompokan pesan berturut-turut.
- **Ekstensi Markdown Editorial**: Penyorotan sintaks kode di server, carousel responsif (` ```carousel `), galeri (` ```gallery `), akordeon FAQ (` ```faq `), kartu profil penulis (` ```author `), catatan kaki, dan daftar isi otomatis.
- **Pengalaman Membaca**: Pengaturan ukuran font (A+/A-), mode Serif, pemilih palet tema (Light, Sepia, Dark), markah bacaan (*reading list*), dan pencarian artikel *in-memory* yang cepat.
- **Feeds & SEO**: Generator otomatis kartu pratinjau OpenGraph, RSS 2.0 feed, JSON Feed, Sitemap XML, dan *Content Security Policy* (CSP) ketat.
