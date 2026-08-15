---
title: "Memahami File Permission Linux: Angka di Balik chmod"
slug: 4dc48488
aliases: [file-permission-linux]
date: 2024-08-12
tags: [linux, tips]
lang: id
draft: false
---

Hampir semua orang yang baru pakai Linux pernah menemui perintah ini di tutorial:

```bash
chmod 755 script.sh
```

Lalu menyalinnya tanpa benar-benar tahu apa arti angka `755`. Selama scriptnya jalan, tidak ada yang bertanya lebih jauh. Masalahnya muncul belakangan, ketika sebuah file tiba-tiba tidak bisa diakses, atau lebih buruk, ketika kamu memberi izin terlalu longgar pada file yang seharusnya rahasia.

Permission di Linux sebenarnya sederhana kalau dipahami logikanya, bukan dihafal angkanya.

## Tiga Kelompok, Tiga Izin

Setiap file dan direktori punya izin untuk tiga kelompok:

- **owner**: pemilik file
- **group**: grup yang memiliki file
- **others**: semua orang lain

Untuk masing-masing kelompok, ada tiga jenis izin: **read (r)**, **write (w)**, dan **execute (x)**. Jalankan `ls -l` dan kamu akan melihatnya:

```bash
$ ls -l script.sh
-rwxr-xr-- 1 dafa staff 220 Aug 12 09:14 script.sh
```

Bagian `-rwxr-xr--` di depan itulah yang penting. Pecah menjadi empat:

```
-    rwx    r-x    r--
tipe owner  group  others
```

Karakter pertama menandakan tipe (`-` untuk file biasa, `d` untuk direktori, `l` untuk symbolic link). Sembilan karakter sisanya adalah izin untuk owner, group, dan others secara berurutan.

Dari contoh di atas: owner bisa baca-tulis-eksekusi, group bisa baca dan eksekusi, others hanya bisa baca.

## Dari Mana Angkanya

Angka seperti `755` hanyalah cara ringkas menuliskan izin yang sama. Setiap izin punya nilai:

- read = 4
- write = 2
- execute = 1

Jumlahkan nilai untuk setiap kelompok:

```
rwx = 4 + 2 + 1 = 7
r-x = 4 + 0 + 1 = 5
r-- = 4 + 0 + 0 = 4
```

Jadi `-rwxr-xr--` sama dengan `754`, bukan `755`. Sekarang angka itu tidak lagi terlihat seperti mantra. Ia hanya penjumlahan.

Beberapa kombinasi yang sering dipakai:

```
644  rw-r--r--   file biasa: owner bisa edit, lainnya baca saja
755  rwxr-xr-x   script/program: semua bisa jalankan, hanya owner edit
600  rw-------   file rahasia: hanya owner yang bisa akses
700  rwx------   direktori pribadi
```

## Cara Simbolik yang Lebih Mudah Dibaca

Selain angka, `chmod` menerima notasi simbolik yang sering lebih jelas maksudnya. Gunakan `u` (user/owner), `g` (group), `o` (others), `a` (all), lalu `+`, `-`, atau `=`.

```bash
# Tambahkan izin eksekusi untuk owner
chmod u+x script.sh

# Hapus izin tulis untuk group dan others
chmod go-w catatan.txt

# Set izin baca untuk semua orang, tanpa menyentuh izin lain
chmod a+r dokumen.pdf
```

Notasi ini punya keunggulan: ia menambah atau mengurangi izin tanpa harus menghitung ulang seluruh angka. Saat kamu hanya ingin membuat satu script bisa dieksekusi, `chmod +x script.sh` jauh lebih jelas daripada mengingat angka tiga digit.

## Permission pada Direktori Berbeda Maknanya

Ini bagian yang sering membingungkan. Pada direktori, izin `x` bukan berarti "eksekusi" melainkan izin untuk **masuk** ke direktori tersebut. Tanpa `x`, kamu tidak bisa `cd` ke dalamnya meskipun punya izin baca.

```bash
# Direktori bisa dilihat isinya tapi tidak bisa dimasuki
chmod 644 folder   # r-- tanpa x → akan menolak cd
```

Inilah kenapa direktori hampir selalu butuh izin `x`. Kombinasi `755` pada direktori berarti semua orang bisa masuk dan melihat isinya, sementara hanya owner yang bisa menambah atau menghapus file di dalamnya.

## Owner dan Group: chown

`chmod` mengatur izin, tapi tidak mengubah siapa pemiliknya. Untuk itu ada `chown`.

```bash
# Ganti owner menjadi dafa
sudo chown dafa file.txt

# Ganti owner dan group sekaligus
sudo chown dafa:developers file.txt

# Terapkan ke seluruh isi direktori secara rekursif
sudo chown -R dafa:developers /var/www/project
```

Kesalahan klasik saat deploy aplikasi web adalah file dimiliki `root` padahal web server berjalan sebagai user lain seperti `www-data`. Akibatnya server tidak bisa membaca file dan muncul error permission denied. `chown` ke user yang benar menyelesaikan ini.

## Soal 777 yang Berbahaya

Di banyak forum, solusi instan untuk masalah permission adalah `chmod 777`. Angka itu memberi izin penuh kepada siapa pun: baca, tulis, eksekusi untuk semua orang.

Pada mesin pribadi mungkin tidak terasa dampaknya. Tapi pada server, ini lubang keamanan serius. Siapa pun yang punya akses ke sistem, termasuk proses yang dikompromikan, bisa menulis ulang file tersebut. File konfigurasi, script yang dieksekusi otomatis, semuanya menjadi rawan.

Aturan praktisnya sederhana: beri izin sekecil mungkin yang masih membuat sesuatu berfungsi. Kalau web server hanya perlu membaca file, jangan beri izin tulis. Kalau file berisi password database, `600` sudah cukup, tidak perlu lebih.

---

Permission Linux bukan sistem yang rumit, hanya jarang dijelaskan dari prinsipnya. Begitu kamu paham bahwa angka itu adalah penjumlahan dari read-write-execute untuk tiga kelompok, kamu tidak perlu lagi menyalin perintah dari tutorial tanpa mengerti. Dan yang lebih penting, kamu berhenti memberi `777` pada hal yang seharusnya dijaga.
