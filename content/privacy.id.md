# Kebijakan Privasi

**Terakhir diperbarui: 27 Agustus 2026**

Daemontalk beroperasi dengan prinsip fundamental bebas pelacak (*zero-tracker*). Kami meyakini bahwa publikasi teknis harus mengedepankan kedaulatan pembaca, transparansi sistem, dan minimalisme komputasi. Kami tidak melacak, memprofilkan, memonetisasi, atau mengumpulkan data pribadi Anda.

## 1. Bebas Pelacak & Infrastruktur Mandiri

- **Tanpa Skrip Pelacak Perilaku**: Situs ini beroperasi tanpa Google Analytics, Meta/Facebook Pixel, suar pelacak (*tracking beacons*), kuki iklan, atau teknik sidik jari peramban (*browser fingerprinting*).
- **Tanpa Jaringan Iklan Komersial**: Kami tidak menampilkan banner iklan, tautan afiliasi, atau menjual data kepada pialang data (*data broker*).
- **Aset Mandiri (*Self-Hosted*)**: Seluruh lembar gaya (*CSS*), skrip (*JavaScript*), dan aset media dikirimkan langsung dari peladen binari Go mandiri kami, mengeliminasi paparan pengawasan pihak ketiga.

## 2. Penyimpanan Sisi Klien & Preferensi Lokal

Kami menggunakan mekanisme penyimpanan bawaan peramban (`localStorage` dan `sessionStorage`) secara eksklusif pada perangkat lokal Anda untuk menyimpan preferensi kenyamanan membaca:

- **Daftar Bacaan Tersimpan (`bookmarks`)**: Daftar artikel yang Anda simpan di buku besar bacaan (*Reading List*).
- **Pilihan Tema Tampilan (`theme`)**: Pilihan tema warna aktif (*Mode Terang, Gelap, atau Sepia*).
- **Tipografi & Aksesibilitas**: Skala ukuran fon, pilihan jenis fon (*Serif atau Sans*), serta intensitas kehangatan layar (*warm screen tint*).
- **Penanda Riwayat Artikel Dibaca (`readPosts`)**: Larik (*array*) lokal berisi *slug* artikel yang telah dibaca (dibatasi maksimal 200 data) semata-mata untuk membedakan artikel baru dan yang sudah dibaca pada halaman indeks.
- **Animasi Sesi (`visited`)**: Penanda sementara pada `sessionStorage` untuk mencegah pengulangan animasi pembuka saat berpindah halaman dalam satu sesi.

**Kedaulatan Data**: Seluruh data preferensi lokal ini tidak pernah diunggah, disinkronkan, atau dikirimkan ke basis data peladen kami maupun pihak luar. Anda dapat menghapus data ini sewaktu-waktu dengan membersihkan penyimpanan situs (*site data/cache*) di peramban Anda.

## 3. Log Peladen Sementara & Telemetri

Saat Anda terhubung ke Daemontalk melalui HTTP atau SSH, peladen mandiri kami mencatat metadata koneksi standar:

- Alamat IP pengunjung.
- Jalur berkas (*path*), metode HTTP, dan waktu permintaan (*timestamp*).
- Informasi *User-Agent* dan *Referrer* (apabila dikirimkan oleh peramban Anda).
- Kode status respon HTTP dan latensi pemrosesan.

**Kebijakan Retensi Log**: Log koneksi ini disimpan murni untuk pemantauan keamanan secara *real-time*, mitigasi serangan *brute-force* / DDoS (*rate-limiting*), dan diagnostik infrastruktur. Seluruh log akan dirotasi dan dihapus permanen secara otomatis setelah **14 hari**. Kami tidak pernah membuat profil kepribadian lintas situs atau menghubungkan alamat IP ke identitas dunia nyata.

## 4. Partisipasi Interaktif (Komentar, Reaksi, & Buku Tamu)

Saat Anda memanfaatkan fitur interaktif komunitas:

- **Diskusi Publik**: Saat Anda mengirimkan komentar atau pesan pada Buku Tamu, nama alias (*Callsign*) dan isi pesan yang Anda masukkan disimpan di basis data SQLite lokal kami dan ditampilkan secara terbuka.
- **Reaksi Artikel**: Reaksi pada artikel (*Like*, *Insightful*, dsb.) hanya menambah jumlah akumulasi angka pada artikel tanpa merekam data identitas pribadi Anda.
- **Tanpa Kewajiban Akun**: Anda tidak diwajibkan mendaftarkan kata sandi, alamat email, nomor telepon, atau akun media sosial untuk berdiskusi.
- **Pencegahan Bot Tanpa Pengawasan**: Kami menggunakan kolom *honeypot* tak kasat mata untuk menyaring bot spam tanpa membebani pembaca dengan widget CAPTCHA komersial yang melacak aktivitas Anda.

## 5. Antarmuka Terminal & Shell

- **Terminal Web (`/terminal`)**: Shell UNIX virtual di peramban berjalan sepenuhnya di sisi klien; riwayat pengetikan perintah hanya berada di memori lokal dan akan terhapus otomatis saat *tab* ditutup.
- **Akses SSH Publik (`ssh daemontalk.com -p 2222`)**: Sesi koneksi SSH berjalan pada proses terisolasi tanpa perekaman penekanan tombol (*zero keystroke logging*).

## 6. Hak Subjek Data & Penghapusan Konten

Anda berhak mengajukan ralat atau penghapusan permanen atas komentar maupun pesan buku tamu yang pernah Anda kirimkan. Silakan ajukan permintaan penghapusan melalui email ke **realdaemontalk@gmail.com** dengan menyertakan tautan artikel dan nama alias yang Anda gunakan.

## 7. Tautan & Rujukan Eksternal

Artikel teknis kami sering kali menyertakan rujukan ke repositori kode sumber terbuka (GitHub), portal dokumentasi resmi, serta publikasi jurnal ilmiah (arXiv, IEEE). Saat Anda mengeklik tautan ke luar Daemontalk, aktivitas penjelajahan Anda tunduk pada kebijakan privasi masing-masing situs tujuan tersebut.

## 8. Kontak Perlindungan Privasi

Apabila Anda memiliki pertanyaan, saran audit keamanan, atau permohonan terkait arsitektur privasi situs ini, silakan hubungi kami langsung di: **realdaemontalk@gmail.com**.
