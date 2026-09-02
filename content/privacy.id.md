# Kebijakan Privasi

**Terakhir diperbarui: 30 Agustus 2026** · **Versi: 2.1**

Daemontalk beroperasi dengan prinsip fundamental bebas pelacak (*zero-tracker*). Kami meyakini bahwa publikasi teknis dan komunitas rekayasa perangkat lunak harus mengedepankan kedaulatan pembaca, transparansi sistem, dan minimalisme komputasi. Kami tidak melacak, memprofilkan, memonetisasi, atau mengumpulkan data pribadi Anda secara terselubung.

---

## 1. Bebas Pelacak & Infrastruktur Mandiri

**Tanpa Skrip Pelacak Perilaku**: Situs ini beroperasi tanpa Google Analytics, Meta/Facebook Pixel, suar pelacak (*tracking beacons*), kuki iklan pihak ketiga, atau teknik sidik jari peramban (*browser fingerprinting*).

**Tanpa Jaringan Iklan Komersial**: Kami tidak menampilkan banner iklan, tautan afiliasi, atau menjual data kepada pialang data (*data broker*).

**Aset Mandiri (*Self-Hosted*)**: Seluruh lembar gaya (*CSS*), skrip (*JavaScript*), fon tipografi, dan aset media dikirimkan langsung dari peladen binari Go mandiri kami, mengeliminasi paparan pengawasan pihak ketiga.

## 2. Basis Hukum Pemrosesan Data (UU PDP & GDPR)

Sesuai dengan Undang-Undang No. 27 Tahun 2022 tentang Pelindungan Data Pribadi (UU PDP) di Indonesia dan *General Data Protection Regulation* (GDPR), kami memproses data minimal semata-mata atas dasar hukum berikut:

**Persetujuan Eksplisit (*Consent*)**: Saat Anda secara sadar memilih masuk (*login*) via GitHub OAuth, membuat topik diskusi, membalas pertanyaan forum, atau mengirim formulir kontak.

**Kepentingan yang Sah (*Legitimate Interest*)**: Untuk pemeliharaan keandalan server, mitigasi serangan DDoS dan *brute-force* via *rate-limiting*, serta diagnostik teknis menggunakan log koneksi sementara yang dirotasi otomatis dengan retensi ketat 14 hari.

**Kewajiban Kontraktual / Operasional Layanan**: Untuk memvalidasi sesi aktif pengguna dan menghubungkan kepemilikan solusi teknis yang Anda terbitkan di forum.

## 3. Autentikasi Opsional GitHub OAuth & Profil Pengguna

Daemontalk menyediakan fitur masuk (*login*) opsional berbasis GitHub OAuth untuk memfasilitasi identitas terverifikasi pada forum diskusi, tanya-jawab sistem, dan komentar artikel teknis.

**Data Minimal yang Dikumpulkan**: Saat Anda masuk menggunakan akun GitHub, kami hanya meminta izin baca minimal (`read:user`, `user:email`). Kami hanya menyimpan ID Pengguna GitHub (*Provider ID*), *username* / *handle* publik GitHub (misal `octocat`), nama tampilan (*Display Name*), URL Avatar publik, URL profil GitHub publik, dan alamat email utama yang terverifikasi (digunakan murni untuk validasi sesi dan tidak pernah dijual atau dikirimi email pemasaran spam).

**Data yang Tidak Pernah Kami Akses**: Kami tidak pernah meminta, mengakses, atau menyimpan repositori privat, kode sumber proyek Anda, kunci SSH, data pembayaran/kartu kredit, token organisasi, maupun izin tulis (*write access*) ke akun GitHub Anda.

**Pengelolaan Sesi (*Session Management*)**: Sesi login dikelola menggunakan token acak berkode hash kriptografis SHA-256 yang disimpan dalam kuki peramban berstatus *HTTP-only* dan *SameSite=Lax*. Kami tidak menggunakan pelacak sesi pihak ketiga.

## 4. Penyimpanan Sisi Klien & Preferensi Lokal

Kami menggunakan mekanisme penyimpanan bawaan peramban (`localStorage` dan `sessionStorage`) secara eksklusif pada perangkat lokal Anda untuk menyimpan preferensi kenyamanan membaca:

**Daftar Bacaan Tersimpan (`bookmarks`)**: Daftar artikel yang Anda simpan di buku besar bacaan (*Reading List*).

**Pilihan Tema Tampilan (`theme`)**: Pilihan tema warna aktif (*Mode Terang atau Gelap*).

**Tipografi & Aksesibilitas**: Skala ukuran fon, pilihan jenis fon (*Serif atau Sans*), serta intensitas kehangatan layar (*warm screen tint*).

**Penanda Riwayat Artikel Dibaca (`readPosts`)**: Larik (*array*) lokal berisi *slug* artikel yang telah dibaca (dibatasi maksimal 200 data) semata-mata untuk membedakan artikel baru dan yang sudah dibaca pada halaman indeks.

**Animasi Sesi (`visited`)**: Penanda sementara pada `sessionStorage` untuk mencegah pengulangan animasi pembuka saat berpindah halaman dalam satu sesi.

**Kedaulatan Data**: Seluruh data preferensi lokal ini tidak pernah diunggah, disinkronkan, atau dikirimkan ke basis data peladen kami maupun pihak luar tanpa tindakan eksplisit Anda. Anda dapat menghapus data ini sewaktu-waktu dengan membersihkan penyimpanan situs (*site data/cache*) di peramban Anda.

## 5. Log Peladen Sementara & Telemetri

Saat Anda terhubung ke Daemontalk melalui HTTP atau SSH, peladen mandiri kami mencatat metadata koneksi standar termasuk alamat IP pengunjung, jalur berkas (*path*), metode HTTP, waktu permintaan (*timestamp*), informasi *User-Agent* dan *Referrer*, kode status respon HTTP, serta latensi pemrosesan.

**Kebijakan Retensi Log**: Log koneksi ini disimpan murni untuk pemantauan keamanan secara *real-time*, mitigasi serangan *brute-force* / DDoS (*rate-limiting*), dan diagnostik infrastruktur. Seluruh log akan dirotasi dan dihapus permanen secara otomatis setelah 14 hari. Kami tidak pernah membuat profil kepribadian lintas situs atau menghubungkan alamat IP ke identitas dunia nyata.

## 6. Partisipasi Komunitas Interaktif (Forum Diskusi, Komentar, & Reaksi)

Saat Anda memanfaatkan fitur interaktif komunitas:

**Forum & Tanya Jawab (`/discussions`)**: Pengguna terautentikasi dapat membuat topik diskusi teknis baru, mengirimkan solusi jawaban, membalas komentar bertingkat, dan memberikan dukungan (*upvote*). Konten diskusi berformat Markdown, status penyelesaian (*Solved*), dan stempel waktu disimpan di basis data SQLite lokal kami dan ditampilkan secara terbuka.

**Komentar Artikel**: Anda dapat mengirimkan komentar pada artikel menggunakan profil terverifikasi GitHub atau secara anonim dengan nama panggilan sementara.

**Reaksi Artikel**: Reaksi pada artikel (*Like*, *Insightful*, dsb.) hanya menambah jumlah akumulasi angka pada artikel tanpa merekam data identitas pribadi Anda.

**Pencegahan Bot Tanpa Pengawasan**: Kami menggunakan kolom *honeypot* tak kasat mata dan pembatasan frekuensi (*in-memory rate limiting*) untuk menyaring bot spam tanpa membebani pembaca dengan widget CAPTCHA komersial yang melacak aktivitas Anda.

## 7. Antarmuka Terminal & Shell

**Terminal Web (`/terminal`)**: Shell UNIX virtual di peramban berjalan sepenuhnya di sisi klien; riwayat pengetikan perintah hanya berada di memori lokal dan akan terhapus otomatis saat *tab* ditutup.

**Akses SSH Publik (`ssh daemontalk.com -p 2222`)**: Sesi koneksi SSH berjalan pada proses terisolasi tanpa perekaman penekanan tombol (*zero keystroke logging*).

## 8. Standar Keamanan Data & Kriptografi

Kami menerapkan rekayasa pertahanan berlapis (*defense-in-depth*) untuk mengamankan data yang disimpan maupun yang ditransmisikan:

**Enkripsi Jalur Komunikasi**: Seluruh lalu lintas web diamankan menggunakan TLS 1.3 dengan *Perfect Forward Secrecy* (PFS) dan *HTTP Strict Transport Security* (HSTS).

**Hasing Sesi**: Token sesi login mentah tidak pernah disimpan dalam bentuk teks polos; hanya hash SHA-256 ber-garam (*salted*) yang disimpan di basis data.

**Isolasi Database**: Basis data SQLite disimpan pada sistem berkas terisolasi dengan hak akses Unix ketat (`0600`) dan kueri terparameterisasi (*prepared statements*) untuk mengeliminasi celah *SQL Injection*.

## 9. Hak Subjek Data & Portabilitas Mandiri (*Self-Service*)

Anda memegang kendali penuh atas kedaulatan data pribadi Anda (sesuai UU PDP No. 27/2022 & GDPR):

**Unduh Data Mandiri (*Self-Service Export*)**: Anda dapat mengunduh seluruh informasi akun, riwayat topik forum, balasan, dan komentar artikel Anda dalam format JSON terstruktur kapan saja dengan mengeklik menu "Unduh data saya (JSON)" di dropdown profil atau mengakses `/auth/export`.

**Penghapusan Akun Mandiri (*Self-Service Purge*)**: Anda dapat menghapus akun Anda secara permanen kapan saja melalui menu "Hapus akun" di dropdown profil. Setelah konfirmasi, profil pengguna dan sesi login Anda akan dihapus permanen dari basis data kami, dan kontribusi publik Anda di forum akan dianonimkan secara otomatis (`[Deleted User]`) demi menjaga kelangsungan arsip diskusi komunitas tanpa menghubungkan identitas pribadi Anda.

**Permohonan Manual**: Anda juga dapat mengajukan permohonan ekspor, perbaikan, atau penghapusan data secara manual melalui email ke **realdaemontalk@gmail.com** menggunakan alamat email yang terhubung dengan akun GitHub Anda.

## 10. Batasan Usia & Privasi Anak

Daemontalk adalah platform riset rekayasa komputer. Kami tidak secara sengaja mengumpulkan atau meminta data pribadi dari anak di bawah usia 13 tahun (atau di bawah 16 tahun di yurisdiksi tertentu). Jika kami mengetahui bahwa data anak di bawah umur telah tersimpan tanpa persetujuan wali yang terverifikasi, kami akan segera menghapus data tersebut.

## 11. Lokasi Peladen & Transfer Data Internasional

Infrastruktur utama kami di-host secara mandiri pada peladen *bare-metal* aman yang berlokasi di pusat data berkepatuhan tinggi. Kami tidak pernah memindahkan, merutekan, atau menyalin data pribadi Anda ke cloud pemasaran pihak ketiga atau broker data internasional.

## 12. Protokol Notifikasi Insiden Keamanan

Apabila terjadi insiden keamanan yang berpotensi memengaruhi integritas atau kerahasiaan data pribadi, Daemontalk berkomitmen untuk memberi tahu pengguna yang terdampak dan pihak berwenang terkait tanpa penundaan yang tidak semestinya (maksimal dalam waktu 3x24 jam sejak insiden diketahui) sesuai ketentuan UU PDP dan regulasi terkait.

## 13. Perubahan Kebijakan & Saluran Kontak

Kami dapat memperbarui Kebijakan Privasi ini secara berkala seiring evolusi arsitektur sistem, fitur baru, atau penyesuaian regulasi hukum. Perubahan material akan dicatat pada halaman `/changelog` dan ditandai dengan tanggal pembaruan di bagian atas halaman ini.

Apabila Anda memiliki pertanyaan, saran audit keamanan, atau permohonan terkait arsitektur privasi situs ini, silakan hubungi kami langsung di: **realdaemontalk@gmail.com**.
