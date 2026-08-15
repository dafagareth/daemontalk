---
title: "Bash Scripting yang Benar: Tiga Baris yang Sering Dilupakan"
slug: febaa50e
aliases: [bash-scripting-benar]
date: 2024-09-23
tags: [bash, linux, scripting]
lang: id
draft: false
---

Sebagian besar bash script yang beredar di internet dimulai langsung dengan perintah, tanpa pengaman apa pun. Selama input yang diberikan persis seperti yang diharapkan, script berjalan baik. Begitu ada satu variabel kosong atau satu perintah gagal di tengah jalan, script tetap lanjut seolah tidak terjadi apa-apa, dan kerusakannya baru ketahuan jauh kemudian.

Ada tiga baris di awal script yang mengubah perilaku berbahaya ini. Sayangnya jarang diajarkan.

## Masalahnya

Perhatikan script sederhana ini:

```bash
#!/bin/bash
cd /tmp/build
rm -rf *
```

Kelihatan tidak berbahaya. Tapi bayangkan direktori `/tmp/build` tidak ada. Perintah `cd` gagal, namun bash tetap melanjutkan ke baris berikutnya. Sekarang kamu berada di direktori sebelumnya, mungkin home directory, dan `rm -rf *` menghapus isinya.

Bash secara default tidak berhenti saat sebuah perintah gagal. Ia juga tidak protes saat kamu memakai variabel yang belum pernah didefinisikan. Perilaku ini adalah sumber dari banyak bug yang sulit dilacak.

## Tiga Baris Itu

```bash
#!/bin/bash
set -euo pipefail
```

Mari pecah satu per satu.

**`set -e`** membuat script berhenti begitu ada perintah yang gagal (mengembalikan exit code bukan nol). Tidak ada lagi script yang terus berjalan setelah langkah penting gagal.

**`set -u`** membuat script error ketika kamu memakai variabel yang belum didefinisikan, alih-alih diam-diam menggantinya dengan string kosong. Ini menangkap salah ketik nama variabel sebelum menyebabkan kerusakan.

**`set -o pipefail`** mengubah cara bash menilai pipeline. Secara default, exit code sebuah pipeline diambil dari perintah terakhir saja. Dengan `pipefail`, jika perintah mana pun dalam pipe gagal, seluruh pipeline dianggap gagal.

```bash
# Tanpa pipefail, ini dianggap sukses meski curl gagal
curl -s http://server-yang-mati | grep "data"

# Dengan pipefail, kegagalan curl terdeteksi
```

Kombinasi ketiganya membuat script berperilaku seperti yang sebenarnya kamu harapkan: berhenti saat ada yang salah, bukan terus melaju ke jurang.

## Selalu Beri Tanda Kutip pada Variabel

Aturan kedua yang sering diabaikan: variabel hampir selalu perlu tanda kutip ganda.

```bash
file="laporan akhir.txt"

# Salah: bash memecah jadi dua argumen "laporan" dan "akhir.txt"
rm $file

# Benar: diperlakukan sebagai satu nama file
rm "$file"
```

Tanpa tanda kutip, bash memecah nilai variabel berdasarkan spasi. Nama file dengan spasi, path yang mengandung karakter khusus, atau input dari pengguna bisa menyebabkan perilaku tak terduga. Membiasakan diri menulis `"$var"` menghilangkan seluruh kategori bug ini.

## Cek Keberadaan Sebelum Bertindak

Operasi yang merusak sebaiknya didahului pengecekan. Bash punya operator test yang ringkas untuk ini.

```bash
# Pastikan direktori ada sebelum masuk
if [ ! -d "$build_dir" ]; then
    echo "Direktori $build_dir tidak ditemukan" >&2
    exit 1
fi

# Pastikan file ada sebelum diproses
if [ -f "$config" ]; then
    source "$config"
fi
```

Beberapa operator yang sering dipakai: `-f` untuk file biasa, `-d` untuk direktori, `-z` untuk string kosong, `-n` untuk string tidak kosong. Perhatikan juga `>&2` pada pesan error di atas, itu mengarahkan pesan ke standard error, bukan standard output, sehingga tidak tercampur dengan output normal script.

## Variabel dengan Nilai Default

`set -u` akan menghentikan script saat variabel tidak ada. Tapi kadang kamu memang ingin variabel opsional dengan nilai cadangan. Bash menyediakan sintaks untuk itu.

```bash
# Pakai nilai $1 jika ada, jika tidak pakai "production"
env="${1:-production}"

# Pakai $PORT jika diset, jika tidak pakai 8080
port="${PORT:-8080}"
```

Pola `${var:-default}` ini membuat script tetap fleksibel tanpa mengorbankan keamanan dari `set -u`.

## Fungsi untuk Bagian yang Berulang

Saat script mulai panjang, pecah menjadi fungsi. Selain lebih rapi, ini memudahkan penanganan error.

```bash
#!/bin/bash
set -euo pipefail

log() {
    echo "[$(date +'%H:%M:%S')] $*"
}

cleanup() {
    log "Membersihkan file sementara"
    rm -rf "$tmp_dir"
}

# Jalankan cleanup otomatis saat script keluar, baik sukses maupun gagal
tmp_dir="$(mktemp -d)"
trap cleanup EXIT

log "Memulai proses"
# ... pekerjaan utama
log "Selesai"
```

Perhatikan `trap cleanup EXIT`. Perintah ini memastikan fungsi `cleanup` dijalankan ketika script berakhir, apa pun penyebabnya. File sementara tidak akan tertinggal meski script berhenti di tengah karena error.

---

Bash punya reputasi sebagai bahasa yang rawan kesalahan, dan sebagian reputasi itu pantas. Tapi banyak masalahnya berasal dari kebiasaan menulis script tanpa pengaman, bukan dari bahasanya sendiri. Tiga baris `set -euo pipefail`, tanda kutip yang konsisten, dan pengecekan sebelum operasi berbahaya sudah cukup membuat script kamu jauh lebih bisa diandalkan daripada mayoritas yang beredar di luar sana.
