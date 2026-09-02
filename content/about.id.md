# Tentang Daemontalk

Daemontalk adalah publikasi rekayasa sistem komputer independen, buku catatan riset terbuka, dan ruang eksplorasi komputasi tingkat rendah tanpa pelacak pihak ketiga.

---

## Misi & Pendekatan Eksplorasi

Alih-alih sekadar merangkum teori abstrak akademis, Daemontalk dibangun sebagai arsip catatan kerja sistem dan portofolio rekayasa sistem informasi. Fokus utama di sini adalah membedah dan memverifikasi langsung cara kerja sistem nyata: menganalisis perilaku runtime bahasa, menyelami mekanisme internal kernel Linux, model konkurensi Go, hingga arsitektur penyimpanan performa tinggi dengan tolok ukur yang dapat diuji dan direproduksi mandiri.

Setiap tulisan di sini mengutamakan reproduksibilitas: mulai dari perintah shell, skrip pengujian beban performa, diagram arsitektur, hingga kode sumber terbuka yang dapat dijalankan langsung di lingkungan lab Anda.

## Kurator & Komunitas

Situs ini diprakarsai dan dikurasi oleh **Dafa Gareth** sebagai catatan studi dan riset sistem informasi. Meskipun dikelola secara mandiri, Daemontalk terbuka bagi tulisan rekan-rekan engineer, arsitek sistem, dan peneliti yang ingin membagikan catatan teknisnya melalui kontribusi Pull Request di GitHub maupun forum diskusi komunitas.

## Prinsip & Standar Editorial

Untuk menjaga integritas teknis dan kenyamanan membaca, seluruh publikasi di Daemontalk berpegang teguh pada prinsip-prinsip berikut:

**Lugas & To The Point**: Tulisan langsung masuk ke inti persoalan teknis, arsitektur, dan kode. Menghindari basa-basi pembuka dan klise yang tidak memiliki nilai teknis nyata.

**Fakta & Rujukan Terverifikasi**: Setiap artikel mendalam (*deep-dive*) dilengkapi dengan blok referensi resmi, seperti RFC standar internet, repositori kode sumber kernel Linux, dokumen manual arsitektur prosesor, atau *paper* riset ilmiah.

**Zero Tracking & Kedaulatan Data**: Dihosting secara mandiri dalam satu binari Go tanpa Google Analytics, pixel pelacak iklan, *paywall*, maupun skrip pengawasan pihak ketiga.

## Fitur & Antarmuka Interaktif

Selain publikasi artikel teknis mendalam, Daemontalk dilengkapi dengan berbagai antarmuka komputasi interaktif:

**Forum Diskusi & Tanya Jawab (`/discussions`)**: Ruang kolaborasi teknis terbuka bagi anggota untuk mendiskusikan arsitektur sistem, memecahkan masalah error produksi, dan berbagi solusi kode dengan integrasi masuk resmi GitHub OAuth.

**Terminal UNIX di Peramban (`/terminal`)**: Antarmuka shell virtual berbasis web yang berjalan sepenuhnya di sisi klien untuk eksplorasi perintah sistem dan utilitas diagnostik secara interaktif.

**Akses Server SSH Publik (`ssh daemontalk.com -p 2222`)**: Akses antarmuka baris perintah (CLI/TUI) langsung melalui protokol SSH tanpa perlu mengunduh aplikasi tambahan.

## Teknologi di Balik Layar

Web ini dibangun sebagai satu berkas biner Go mandiri menggunakan router `chi`, mesin templat terkompilasi `templ`, dan TailwindCSS. Seluruh halaman dirender cepat di sisi peladen (*server-side rendering*) tanpa framework JavaScript berat agar tetap ringan, aman, dan efisien dalam penggunaan sumber daya memori. Anda dapat melihat rincian infrastruktur lengkap di halaman [Di Balik Layar (Colophon)](/id/colophon).

## Tautan & Hubungi Kami

Tertarik untuk berdiskusi, memberikan masukan, atau membagikan catatan lab teknis Anda? Silakan baca [Panduan Kontribusi](/id/contribute) atau hubungi kami langsung via email di [realdaemontalk@gmail.com](mailto:realdaemontalk@gmail.com).
