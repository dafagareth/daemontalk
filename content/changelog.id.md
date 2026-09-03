### v1.4.0 (4 September 2026) · Deduplikasi Tayangan Unik, Migrasi Tag Socket, Keamanan Anti-XSS & Pemolesan UI

- **Deduplikasi Tayangan Unik (Human Views)**: Menghapus kenaikan angka views akibat refresh halaman (F5) dengan mengimplementasikan tabel pelacak penonton unik (`forum_topic_views` dan `post_views` di SQLite). Tayangan dihitung murni per pengunjung unik (`u:<userID>` untuk pengguna login, `v:<visitorID>` untuk tamu), serta secara otomatis mengecualikan pemilik postingan yang membuka artikelnya sendiri, admin, bot perayap web, dan request CLI.
- **Arsitektur Murni Tag pada Socket**: Menghapus dropdown kategori statis (`qna`, `architecture`, `kernel`) pada form pembuatan topik dan beralih sepenuhnya ke ekosistem dinamis berbasis `#tag`, lengkap dengan bilah filter tag aktif dan tampilan tag riil pada setiap baris diskusi.
- **Keamanan Markdown & Pencegahan Stored XSS**: Memperkuat pemrosesan konten pengguna di forum Socket dengan mencabut opsi `html.WithUnsafe()` dan menerapkan sanitasi ketat menggunakan `bluemonday.UGCPolicy()`, menangkal injeksi `<script>`, *event handler* berbahaya, dan manipulasi tautan `javascript:`.
- **Desain Ulang Antarmuka Socket & Diskusi**: Menyesuaikan tata letak kotak balasan dan form penulisan agar tampil *full-bleed* hingga ke tepi layar perangkat ponsel, memindahkan tombol kirim komentar ke sudut kanan bawah, memperbarui lencana `SOLVED` dan `SOLUSI` dengan gaya tombol kontras tinggi, memperbesar avatar profil, serta menyederhanakan tombol balas dengan fokus otomatis ke kolom komentar.
- **Penyempurnaan Header Artikel & Bilah Samping**: Menata ulang susunan metadata tajuk artikel dengan menempatkan foto profil di sisi kiri nama penulis, baris kontributor, serta tanggal dan durasi baca di bawahnya. Menyesuaikan tingkat transparansi bilah samping (Daftar Isi dan Kontrol Baca) menjadi 40% (`opacity-40`) saat posisi diam untuk kenyamanan membaca.

### v1.3.2 (2 September 2026) · Perbaikan Izin Eksekusi Skrip Deployment CI/CD

- **Izin Eksekusi Deployment VPS**: Menambahkan `chmod +x` otomatis dan pemanggilan `bash` eksplisit pada workflow GitHub Actions CD untuk mencegah galat *permission denied* saat eksekusi skrip di server VPS.

### v1.3.1 (2 September 2026) · Penanganan Metode Webhook & Dukungan Push Tag

- **Ketahanan Metode Webhook & GET**: Menangani request `GET` secara adaptif pada endpoint `/api/webhook/github` serta mengonfigurasi Caddy dengan redirect HTTP 308 untuk menjaga metode POST saat melewati reverse proxy.
- **Dukungan Push Tag Git**: Memperluas listener webhook GitHub agar memicu reload artikel secara otomatis saat terjadi push tag rilis (`refs/tags/*`) selain push branch utama.

### v1.3.0 (2 September 2026) · GitHub OAuth2, Forum Diskusi Socket & Lisensi Modular

- **Autentikasi GitHub OAuth2 & Profil Pengguna**: Menghadirkan sistem login GitHub resmi (`/auth/github`), sesi pengguna berbasis SQLite, badge profil navigasi desktop dan mobile, serta kartu profil publik (`/u/:username`).
- **Forum Diskusi Socket**: Membangun platform forum komunitas mandiri yang dirender di sisi server (`/socket`) dengan dukungan kategori topik, balasan bertingkat, sistem voting, dan penanda solusi terverifikasi.
- **Lisensi PolyForm Noncommercial & CC BY-NC-SA 4.0**: Mengadopsi lisensi PolyForm Noncommercial 1.0.0 untuk kode sumber, lisensi CC BY-NC-SA 4.0 untuk artikel riset, serta perlindungan hak cipta merek penuh (*All Rights Reserved*).
- **Publikasi Otomatis via GitHub Webhook**: Menambahkan endpoint webhook terverifikasi HMAC-SHA256 (`/api/webhook/github`) untuk me-reload artikel secara instan saat terjadi merge ke branch utama tanpa *downtime*.
- **Desain Ulang Portal Topik & Antarmuka**: Memperbarui portal arsip tag (`/blog/tag/:tag`) dengan kisi kartu 4 kolom minimalis, merapikan navigasi *drawer mobile*, dan menyederhanakan direktori kontributor.
- **Pembersihan Modul Lawas**: Menghapus modul usang seperti `/links`, web `/terminal`, `/guestbook`, dan engine eksekusi runner yang tidak terpakai.
- **Spesifikasi Colophon**: Menambahkan halaman arsitektur perangkat keras dan tumpukan perangkat lunak peladen (`/colophon`).

### v1.2.0 (28 Agustus 2026) · Penguatan Keamanan, Lokalisasi Penuh & Refaktor Arsitektur

- **Penguatan Keamanan & Otorisasi**: Memperkuat batasan akses, menerapkan validasi skema URL untuk mencegah XSS, membatasi ukuran muatan POST, serta menutup potensi kebocoran draf artikel pada endpoint feed dan CLI.
- **Pengalaman Dwibahasa Penuh**: Menuntaskan integrasi kamus UI multibahasa untuk seluruh modal interaktif, pencarian, tautan navigasi, dan buku tamu dalam Bahasa Indonesia (`/id`) dan English.
- **Perbaikan Pagination & Navigasi**: Memperbaiki celah pagination saat klik "Load More" agar tidak ada artikel yang terlewat, serta memastikan redirect tautan alias tetap mempertahankan bahasa aktif.
- **Penyederhanaan Tampilan Arsip**: Merapikan daftar *Chronological Archive River* agar berfokus murni pada judul berita, thumbnail, dan metadata tanpa deskripsi panjang yang berulang.
- **Penyederhanaan Palet Tema**: Menghapus mode lawas Sepia untuk menghadirkan mode Terang dan Gelap yang tajam dengan penyorotan kode yang presisi.
- **Stabilitas & Konkurensi SSH TUI**: Mengisolasi sesi tema terminal per pengguna, mencegah crash server akibat error pembacaan berkas, dan mengoptimalkan pencarian artikel menjadi O(1).

### v1.1.0 (27 Agustus 2026) · Optimisasi UI Mobile & Pembersihan Kode

- **Optimisasi Tampilan Mobile**: Menerapkan batasan layar penuh tanpa jarak (*edge-to-edge padding-free*) secara ketat dan rasio gambar 16:9 yang konsisten pada seluruh *thumbnail* di beranda untuk perangkat ponsel.
- **Desain Ulang Daftar Bacaan**: Memaksimalkan halaman 'Saved Dispatches' dengan mengubahnya menjadi tampilan buku besar (*ledger*) satu kolom bergaris zebra yang minimalis di layar ponsel.
- **Pembersihan Navigasi**: Merapikan laci navigasi *sidebar mobile*, menambahkan tautan esensial yang hilang, menyesuaikan label teks, serta menghapus ikon media sosial yang tidak relevan.
- **Pembersihan Source Code**: Melakukan pembersihan total pada fail *template*, memangkas baris komentar *styling* berlebihan, komentar lawas, serta teks sisa dari AI, menghasilkan kode yang lebih murni dan profesional.
- **Perbaikan Bug Pencarian**: Memperbaiki *bug* pada fitur pencarian berbasis HTMX, di mana menekan tombol Enter sebelumnya mengakibatkan respons HTML mentah tanpa desain (*deadcode*). Kini, halaman penuh dirender dengan sempurna.

### v1.0.6 (Agustus 24, 2026) · Pembersihan Fitur Menyeluruh

- **Pemusnahan Total TIL**: Menyisir sisa-sisa basis kode dan menghapus secara permanen seluruh dependensi "Today I Learned" (TIL) yang masih tertinggal, termasuk model data yang usang, kamus pelokalan bahasa, serta data uji tiruan (*mock*) pada *backend*.
- **Penghapusan Parser Author**: Memusnahkan *markdown parser* khusus untuk blok kode ` ```author `. Sistem kini tidak lagi mem-parsing dan menampilkan kartu profil penulis yang besar, demi menjaga pengalaman baca artikel yang bersih dan terfokus penuh pada konten teks.

### v1.0.5 (Agustus 23, 2026) · TUI LaTeX Graceful Degradation

- **Dukungan Matematika pada Terminal UI**: Memperbaiki *bug* di mana blok matematika LaTeX (`$$`) merusak format terminal akibat tabrakan dengan sintaks *markdown*. TUI kini dilengkapi *preprocessor regex* yang secara elegan mereduksi rumus LaTeX menjadi blok kode (*code blocks*), memastikan formula tetap utuh dan sangat nyaman dibaca lewat koneksi SSH.

### v1.0.4 (Agustus 23, 2026) · Pemangkasan Fitur & Optimasi

- **Penghapusan Fitur TIL**: Memusnahkan fitur *micro-blogging* "Today I Learned" (TIL) beserta seluruh jalur *routing*-nya secara permanen. Keputusan ini diambil untuk memfokuskan publikasi pada artikel teknis yang mendalam dan tajam (*long-form*), serta membersihkan antarmuka dari fitur yang berlebihan.
- **Perbaikan UI Buku Tamu**: Memindahkan tombol kirim pesan buku tamu ke pojok kanan bawah formulir input untuk tata letak yang lebih rapi dan konsisten.

### v1.0.3 (Agustus 23, 2026) · Desain Ulang Sidebar Mobile & Pembersihan UI

- **Perombakan Total Sidebar Mobile**: Mendesain ulang laci navigasi seluler secara penuh. Mengganti tautan *slide-in* generik dengan bilah pencarian atas khusus, topik utama yang menonjol, dan fon *sans-serif* asli demi kenyamanan sentuhan layar ponsel.
- **Penghapusan Jam Live**: Menghapus indikator jam waktu nyata dari *header* desktop dan seluler guna mempertahankan tampilan yang bersih dan bebas gangguan.
- **Ekosistem Sosial**: Menambahkan deretan ikon media sosial rata kiri di bagian bawah laci ponsel. Menghadirkan ikon Facebook dan mengganti ikon Threads dari gaya *solid fill* tebal menjadi SVG *outline* ramping sesuai panduan desain platform. Menyesuaikan tautan eksternal langsung ke `daemontalk`.
- **Standardisasi Dokumentasi**: Menghapus sintaks berlebihan di seluruh fail `.md`. Mengonversi blok HTML `details` yang berantakan ke format bawaan ` ```faq ` serta merapikan penomoran poin pada panduan hukum dan kontribusi.

### v1.0.2 (Agustus 22, 2026) · Penyesuaian Pipeline CI/CD

- **Sinkronisasi Deployment**: Memperbaiki *race condition* pada pipeline otomasi CI/CD di mana VPS menarik *container image* lawas sebelum *image* terbaru selesai diproses. *Deployment* kini diatur agar wajib menunggu hingga proses *build* selesai sepenuhnya.

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
- **Pengalaman Membaca**: Pengaturan ukuran font (A+/A-), mode Serif, pemilih tema (Light, Dark), markah bacaan (*reading list*), dan pencarian artikel *in-memory* yang cepat.
- **Feeds & SEO**: Generator otomatis kartu pratinjau OpenGraph, RSS 2.0 feed, JSON Feed, Sitemap XML, dan *Content Security Policy* (CSP) ketat.

### v1.3.3 (2026-09-02)
- **Fix**: Implementasi *fallback* otomatis menggunakan *Git Log History* untuk tanggal artikel. Artikel kini tidak wajib memiliki `date:` di *frontmatter*.
- **Fix**: Menghapus baris `date:` dari *template CLI* `new-post` untuk menyerahkan penentuan waktu 100% pada *Git commit*.
- **CI/CD**: Menambahkan jalur cepat (Fast-Track) sinkronisasi artikel via `.github/workflows/sync-content.yml`.
- **CI/CD**: Mengecualikan folder `content/**` dari antrean *build* Docker, memangkas waktu *publish* artikel dari 5 menit menjadi 5 detik.
- **Chore**: Mengubah *metadata OG image* dan deskripsi *JSON Feed* agar lebih relevan dengan visi platform baru.
