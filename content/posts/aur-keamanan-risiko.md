---
title: "Kerentanan dan Risiko Keamanan AUR (Arch User Repository)"
slug: c4e7a1f9
aliases: aur-keamanan-risiko
date: 2026-04-22
tags: [linux, security]
lang: id
draft: false
---

AUR (Arch User Repository) adalah repositori berbasis komunitas yang memungkinkan pengguna Arch Linux mendistribusikan skrip build perangkat lunak yang tidak tersedia di repositori resmi. Model distribusinya yang terbuka memberikan fleksibilitas besar, tetapi juga menghadirkan risiko keamanan yang tidak boleh diabaikan.

## Fakta Menarik

**Fakta 1.** AUR tidak memiliki proses peninjauan kode otomatis. Siapa pun dapat membuat akun dan mengunggah PKGBUILD, sehingga tanggung jawab verifikasi sepenuhnya ada di tangan pengguna.

**Fakta 2.** Pada tahun 2021, paket `python-codecov` di AUR terdeteksi mengandung kode berbahaya yang mencuri variabel lingkungan dan token CI/CD dari mesin pengembang. Paket tersebut telah diunduh ribuan kali sebelum dilaporkan.

**Fakta 3.** AUR Helper seperti `yay` dan `paru` secara teknis melakukan `makepkg` atas nama pengguna, tetapi keduanya tidak dapat menggantikan penilaian manusia terhadap isi PKGBUILD sebelum eksekusi.

---

## Tips dan Trik

### 1. Pahami Model Trust AUR

AUR beroperasi berdasarkan kepercayaan komunitas. Setiap paket dimiliki oleh satu atau beberapa "Trusted User" atau kontributor biasa. Paket yang belum lama diunggah, memiliki sedikit komentar, atau berganti kepemilikan secara tiba-tiba patut dicurigai.

Periksa metadata paket sebelum menginstal:

```bash
# Lihat informasi paket termasuk tanggal modifikasi terakhir dan pemilik
paru -Si nama-paket

# Kunjungi halaman AUR secara langsung
xdg-open "https://aur.archlinux.org/packages/nama-paket"
```

### 2. Baca PKGBUILD Sebelum Build

Ini adalah langkah paling penting yang sering dilewati pengguna. PKGBUILD adalah skrip bash biasa yang dieksekusi dengan hak istimewa pengguna saat proses build.

```bash
# Clone repositori AUR secara manual
git clone https://aur.archlinux.org/nama-paket.git
cd nama-paket

# Baca setiap baris PKGBUILD
cat PKGBUILD

# Perhatikan fungsi-fungsi ini secara khusus:
# - prepare()   : sering digunakan untuk memodifikasi sumber
# - build()     : tempat kompilasi, bisa menjalankan perintah arbitrer
# - package()   : instalasi file ke direktori staging
# - install=    : file .install yang dieksekusi saat pasca-instalasi

# Periksa juga file .install jika ada
cat nama-paket.install 2>/dev/null
```

Waspadai pola seperti `curl | bash`, pengunduhan biner dari URL eksternal yang tidak diverifikasi, atau penggunaan `eval` pada string yang diambil dari internet.

### 3. Verifikasi Checksum dan Tanda Tangan

PKGBUILD yang baik mencantumkan array `sha256sums` atau `b2sums` untuk setiap sumber. Pastikan nilai tersebut konsisten dengan rilis resmi upstream.

```bash
# Unduh sumber dan verifikasi secara manual
cd nama-paket
makepkg --nobuild --skipchecksums

# Hitung checksum sendiri dan bandingkan
sha256sum src/nama-tarball.tar.gz

# Jika paket menggunakan GPG, impor kunci penulis terlebih dahulu
gpg --recv-keys FINGERPRINT_KUNCI_PENULIS
makepkg --verifysource
```

### 4. Gunakan paru dengan Konfigurasi yang Aman

`paru` lebih direkomendasikan daripada `yay` karena ditulis dalam Rust dan memiliki opsi audit yang lebih baik. Konfigurasikan `paru` agar selalu menampilkan PKGBUILD sebelum build.

```ini
# /etc/paru.conf atau ~/.config/paru/paru.conf

[options]
# Selalu tampilkan diff PKGBUILD sebelum build
AurOnly
NewsOnUpgrade

# Gunakan editor untuk meninjau file sebelum konfirmasi
# Pastikan variabel $EDITOR sudah diset
```

```bash
# Install paru dari AUR (bootstrap manual diperlukan)
git clone https://aur.archlinux.org/paru.git
cd paru
makepkg -si

# Saat upgrade, paru akan menampilkan perubahan PKGBUILD
paru -Syu

# Untuk melihat diff PKGBUILD secara eksplisit
paru --fm bat nama-paket
```

### 5. Sandboxing Build dengan Container

Cara paling aman adalah mengisolasi proses build di dalam container. `makepkg` dapat dijalankan di dalam container Arch Linux yang bersih sehingga sistem host tidak terekspos.

```bash
# Gunakan image Arch Linux resmi
docker pull archlinux:latest

# Jalankan build di dalam container yang terisolasi
docker run --rm -it \
  -v "$(pwd)/output:/output" \
  archlinux:latest \
  bash -c "
    pacman -Syu --noconfirm base-devel git &&
    useradd -m builder &&
    su builder -c '
      git clone https://aur.archlinux.org/nama-paket.git /tmp/pkg &&
      cd /tmp/pkg &&
      makepkg --noconfirm
    ' &&
    cp /home/builder/tmp/pkg/*.pkg.tar.zst /output/
  "

# Instal paket yang telah di-build dari host
sudo pacman -U output/nama-paket-*.pkg.tar.zst
```

Pendekatan ini memastikan bahwa meskipun PKGBUILD mengandung kode berbahaya, dampaknya terbatas pada lingkungan container yang sekali pakai.
