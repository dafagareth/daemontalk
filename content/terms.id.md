# Syarat & Ketentuan Penggunaan

**Terakhir diperbarui: 5 September 2026** · **Versi: 2.3**

Selamat datang di Daemontalk. Dengan mengakses situs web, membaca publikasi artikel, berpartisipasi dalam diskusi Socket (`/socket`), memberikan komentar, memanfaatkan antarmuka gateway SSH TUI (`ssh daemontalk.com -p 2222`), utilitas baris perintah (`curl`), maupun berlangganan umpan sindikasi (RSS/JSON Feed) ("Layanan"), Anda menyetujui untuk terikat dengan Syarat & Ketentuan Penggunaan ini serta [Kebijakan Privasi](/id/privacy) kami. Jika Anda tidak menyetujui bagian mana pun dari ketentuan ini, silakan hentikan penggunaan platform.

---

## 1. Ruang Lingkup Layanan & Definisi Platform

Daemontalk adalah platform publikasi teknologi independen, buku catatan rekayasa perangkat lunak (*open tech notebook*), dan ruang komunitas pengembang yang berfokus pada sistem operasi, arsitektur backend, kernel Linux, dan kultur rekayasa perangkat lunak. Layanan ini mencakup:

- Situs web publik (`daemontalk.com` dan jalur bahasa `/id`) beserta seluruh artikel teknis, esai editorial, dan studi kasus arsitektur.
- Forum Komunitas & Tanya Jawab Teknis **Socket** (`/socket`).
- Sistem komentar artikel dengan autentikasi resmi GitHub atau sesi tamu.
- Antarmuka pembaca terminal melalui gateway **SSH TUI Publik** (`ssh daemontalk.com -p 2222`).
- *Endpoint* baris perintah ramah `curl` (seperti `curl daemontalk.com/daily`, `/recipes`, `/p/:slug`).
- Umpan sindikasi konten publik (RSS 2.0 di `/rss.xml` dan JSON Feed v1.1 di `/feed.json`).

*(Catatan arsitektur: Modul emulator terminal di dalam peramban web `/terminal` telah resmi dihentikan dan dihapus dari arsitektur platform demi fokus pada performa server-side rendering biner Go yang ringan, keamanan terisolasi, dan pengalaman terminal murni via SSH dan curl).*

## 2. Hak Kekayaan Intelektual & Lisensi Konten

Ketentuan hak cipta dan lisensi di seluruh platform Daemontalk diatur secara transparan untuk mendukung ekosistem edukasi terbuka:

**Artikel Teknis & Esai Editorial**: Seluruh artikel teknologi, esai arsitektur, dan materi editorial yang dipublikasikan di bawah direktori `content/posts/` dilisensikan di bawah **[Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International (CC BY-NC-SA 4.0)](https://creativecommons.org/licenses/by-nc-sa/4.0/)**. Anda bebas menyalin, membagikan, dan mengadaptasi materi ini untuk keperluan non-komersial, dengan kewajiban menyertakan atribusi kredit yang jelas (Dafa Gareth / Daemontalk) serta tautan aktif ke artikel sumber asli.

**Cuplikan Kode & Cetak Biru Konfigurasi**: Seluruh cuplikan kode (*code snippets*), parameter sysctl/kernel Linux, skrip shell, dan skema konfigurasi yang tercantum di dalam artikel disediakan di bawah lisensi terbuka permisif (MIT / Unlicense). Anda bebas menyalin, mengadaptasi, dan mengintegrasikannya ke dalam proyek pribadi maupun komersial Anda tanpa royalti.

**Kode Sumber Platform & Biner Peladen**: Seluruh kode sumber perangkat lunak platform Daemontalk (backend Go, templat terkompilasi Templ, stylesheet Tailwind, SSH server Wish, dan engine TUI Bubble Tea) dilisensikan di bawah **[PolyForm Noncommercial License 1.0.0](https://polyformproject.org/licenses/noncommercial/1.0.0)**. Anda bebas membaca, mengaudit, mem-fork, dan mengembangkannya untuk kebutuhan belajar mandiri, riset akademik, atau proyek non-komersial. Dilarang keras mengomersialisasikan, menjual ulang, atau menyediakannya sebagai layanan SaaS berbayar tanpa izin tertulis resmi dari pengelola.

**Merek Dagang & Identitas Visual**: Nama **Daemontalk**, domain `daemontalk.com`, logo mekanik daemon, token desain, tata letak antarmuka, dan seluruh aset identitas visual dilindungi hak cipta eksklusif (Hak Cipta © 2026 Dafa Gareth. *All Rights Reserved*). Dilarang keras mengkloning atau menggunakan identitas merek ini untuk tujuan yang menyesatkan publik.

## 3. Konten Buatan Pengguna, Komentar & Forum Socket

**Pemberian Lisensi Kontribusi**: Saat Anda mengirimkan pertanyaan teknis, solusi arsitektur, cuplikan kode, atau tanggapan pada forum Socket (`/socket`) maupun kolom komentar artikel, Anda memberikan lisensi non-eksklusif, berlaku permanen di seluruh dunia, dan bebas royalti kepada Daemontalk untuk menampilkan, menyusun indeks pencarian, memformat, serta mendistribusikan kontribusi tersebut sebagai bagian dari basis pengetahuan teknologi publik.

**Orisinalitas & Tanggung Jawab Hukum**: Anda tetap memegang hak cipta atas karya orisinal Anda. Anda bertanggung jawab penuh untuk memastikan bahwa konten atau kode yang Anda kirimkan tidak melanggar hak kekayaan intelektual pihak ketiga, perjanjian kerahasiaan (*Non-Disclosure Agreement* / NDA), atau rahasia dagang tempat Anda bekerja.

**Preservasi Diskusi & Anonimisasi Akun**: Apabila Anda memilih untuk menghapus akun, seluruh data profil pribadi Anda akan dihapus permanen dari server kami. Demi menjaga kesinambungan dan manfaat arsip diskusi teknis bagi pembaca lain, konten pertanyaan dan solusi yang telah Anda publikasikan di forum akan tetap tersimpan secara teranonimkan (`[Deleted User]`).

## 4. Akun Pengguna & Autentikasi GitHub OAuth

**Autentikasi Pengembang Terpercaya**: Akses keanggotaan Daemontalk menggunakan integrasi resmi GitHub OAuth. Daemontalk hanya mengakses informasi profil publik dasar (ID GitHub, nama pengguna, avatar, dan email utama) untuk keperluan verifikasi identitas pengembang. Kami tidak pernah meminta, mengakses, atau menyimpan kata sandi akun GitHub Anda.

**Keamanan & Aktivitas Sesi**: Anda bertanggung jawab penuh untuk menjaga keamanan akses akun GitHub Anda serta atas seluruh interaksi yang dilakukan di bawah sesi login Anda. Jika Anda mencurigai adanya akses tidak sah, Anda dapat segera mencabut (*revoke*) izin aplikasi Daemontalk melalui pengaturan keamanan akun GitHub Anda.

**Integritas Kepemilikan Akun**: Akun pengguna bersifat personal, tidak dapat dipindahtangankan, dan dilarang diperjualbelikan atau disewakan kepada pihak lain.

## 5. Pedoman Perilaku Komunitas & Larangan Penggunaan

Saat berpartisipasi dalam diskusi Socket, membuat topik baru, atau mengirimkan komentar artikel:

**Diskusi Teknis Konstruktif**: Pengguna diharapkan mengedepankan diskusi yang rasional, sopan, berbasis fakta atau pengujian yang dapat direproduksi (*reproducible benchmarks*), dan berorientasi solusi. Berdebatlah mengenai arsitektur atau implementasi teknis secara objektif tanpa menyerang pribadi (*ad hominem*).

**Larangan Spam & Konten Berbahaya**: Dilarang keras mengirimkan promosi komersial tanpa izin, tautan afiliasi terselubung, kampanye bot otomatis, teks sampah hasil generate AI yang tidak berbobot, tautan *phishing*, atau *payload* kode eksploitasi/berbahaya.

**Integritas Voting & Anti-Sybil**: Dilarang memanipulasi dukungan voting (*upvote*), memanipulasi penandaan solusi (*Solved mark*), mendongkrak reputasi palsu, atau menggunakan multi-akun/skrip otomasi untuk memintas batas frekuensi (*rate limits*).

**Anti-Pelecehan & Perlindungan Privasi**: Dilarang melakukan perundungan (*bullying*), pelecehan, doxxing data pribadi individu lain, atau ujaran kebencian. Pelanggaran berat akan berakibat pada pemblokiran akun secara permanen.

## 6. Integritas Sistem, Akses Gateway SSH & Utilitas Baris Perintah (CLI)

Daemontalk menyediakan antarmuka akses terminal langsung untuk memberikan pengalaman membaca teknis yang murni dan efisien:

**Akses Gateway SSH TUI Publik (`ssh daemontalk.com -p 2222`)**: Layanan SSH TUI disediakan untuk membaca artikel, ringkasan harian, dan navigasi arsip langsung dari terminal lokal Anda. Anda setuju untuk menggunakan sesi SSH ini secara wajar. Dilarang keras mencoba meloloskan diri dari lingkungan sandbox TUI (*pty breakout*), mengeksploitasi terminal escape sequences, memicu *panic* pada proses SSH server, atau menyalahgunakan koneksi SSH sebagai *proxy/tunnel*.

**Akses Endpoint CLI (`curl`)**: Streaming teks murni melalui utilitas `curl` (seperti `/daily`, `/recipes`, `/p/:slug`) disediakan untuk kemudahan akses dari konsol. Akses berkala diperbolehkan selama mematuhi *rate limit* wajar dan tidak membebani server secara destruktif.

**Keamanan Infrastruktur**: Dilarang keras melakukan pemindaian celah keamanan tanpa izin (*unauthorized port/vulnerability scanning*), serangan penolakan layanan (DDoS/DoS), serangan *brute force* pada port SSH/HTTP, injeksi data berbahaya, atau *scraping* agresif yang melanggar batas penggunaan wajar sistem.

## 7. Penyangkalan Teknis (Technical Disclaimer) & Ketentuan "As-Is"

Seluruh artikel teknologi, catatan arsitektur, parameter kernel Linux (sysctl), skrip shell, konfigurasi basis data, dan hasil pengujian performa (*benchmarks*) disediakan sebagaimana adanya (*as-is*) dan berdasarkan ketersediaan (*as-available*) murni untuk tujuan informasi dan edukasi.

Teknologi perangkat lunak dan lingkungan sistem operasi berkembang sangat cepat. Suatu konfigurasi yang berjalan baik di lingkungan pengujian kami belum tentu cocok dengan topologi infrastruktur Anda. Anda memegang tanggung jawab penuh untuk mengaudit, menguji, dan memvalidasi konfigurasi pada lingkungan lab atau *staging* terisolasi sebelum menerapkannya pada server produksi yang sesungguhnya.

## 8. Batasan Tanggung Jawab Hukum (Limitation of Liability)

Sepanjang diperbolehkan oleh peraturan perundang-undangan Republik Indonesia, pengelola Daemontalk (Dafa Gareth) serta kontributor tidak bertanggung jawab atas kerugian langsung, tidak langsung, insidental, khusus, atau konsekuensial apa pun.

Hal ini mencakup, namun tidak terbatas pada: kegagalan perangkat keras, *kernel panic*, kehilangan atau kerusakan data, *downtime* sistem produksi, terhentinya operasional bisnis, celah keamanan pada implementasi mandiri, atau biaya mitigasi insiden yang timbul dari pemanfaatan informasi, skrip, atau instruksi dari Platform ini.

## 9. Moderasi Komunitas, Penangguhan & Penghapusan Akun

**Kewenangan Moderasi**: Pengelola berhak secara independen meninjau, menyunting kategori, mengunci, menyembunyikan, atau menghapus topik forum dan komentar artikel apa pun yang melanggar ketentuan ini atau berpotensi membahayakan komunitas tanpa pemberitahuan sebelumnya.

**Penutupan Akses Pengguna**: Kami berhak membekukan atau mencabut hak akses akun anggota yang terbukti melakukan pelanggaran berulang, spamming massal, atau ancaman keamanan terhadap platform.

**Penghapusan Akun Mandiri**: Anda berhak menutup dan menghapus akun Anda kapan saja secara mandiri melalui menu "Hapus akun" pada profil pengguna Anda. Penghapusan akun akan langsung membersihkan seluruh sesi aktif dan data identitas pribadi Anda.

## 10. Tautan Pihak Ketiga & Referensi Luar

Artikel dan diskusi di Daemontalk secara rutin merujuk ke sumber daya eksternal, seperti repositori GitHub, dokumen standar RFC/IETF, artikel ilmiah, dan dokumentasi resmi open-source. Daemontalk tidak mengendalikan, mendukung, atau bertanggung jawab atas keakuratan, ketersediaan, maupun kebijakan privasi dari situs web pihak ketiga tersebut.

## 11. Ketersediaan Layanan & Evolusi Platform

Pengelola berhak memperbarui, memodifikasi, menangguhkan, atau menghentikan fitur, *endpoint*, atau komponen layanan apa pun (termasuk forum Socket, gateway SSH, atau rute API) sewaktu-waktu tanpa kewajiban kompensasi, sebagai bagian dari pemeliharaan server berkala, perbaikan keamanan (*patching*), dan refaktorisasi arsitektur perangkat lunak.

## 12. Hukum yang Berlaku, Penyelesaian Sengketa & Kontak Resmi

Syarat & Ketentuan Penggunaan ini diatur dan ditafsirkan berdasarkan hukum Negara Kesatuan Republik Indonesia. Setiap perselisihan yang timbul sehubungan dengan ketentuan ini akan diselesaikan secara musyawarah untuk mufakat melalui konsultasi teknis yang beritikad baik.

Apabila Anda memiliki pertanyaan, klarifikasi hukum, pemberitahuan hak cipta, atau permohonan lisensi komersial seputar ketentuan ini, silakan hubungi kami langsung di:
- **Email Resmi**: [realdaemontalk@gmail.com](mailto:realdaemontalk@gmail.com)
- **Repositori & Diskusi Proyek**: [github.com/dafagareth/daemontalk](https://github.com/dafagareth/daemontalk)
