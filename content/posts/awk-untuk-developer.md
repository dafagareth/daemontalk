---
title: "awk untuk Developer: Filter dan Olah Teks di Terminal"
slug: 4f1343ae
aliases: [awk-untuk-developer]
date: 2025-06-15
tags: [linux, cli, tools]
lang: id
draft: false
---

`awk` sering muncul di skrip shell tapi jarang dipelajari secara serius. Kebanyakan orang hanya tahu satu polanya: `awk '{print $2}'` untuk mengambil kolom kedua, lalu berhenti di sana.

Padahal dengan memahami beberapa konsep inti `awk`, banyak pekerjaan pengolahan teks yang biasanya diselesaikan dengan Python bisa langsung ditangani di terminal.

## Cara Kerja Dasar

`awk` memproses teks baris per baris. Setiap baris dipecah menjadi kolom berdasarkan pemisah (defaultnya: spasi atau tab). Kolom pertama disimpan di `$1`, kedua di `$2`, dan seterusnya. `$0` adalah keseluruhan baris.

```bash
echo "Alice 27 engineer" | awk '{print $1, $3}'
# Alice engineer
```

## Memfilter dengan Kondisi

Tempatkan kondisi sebelum blok `{}` untuk memfilter baris:

```bash
ps aux | awk '$3 > 5.0 {print $1, $2, $3}'
```

Baris ini mencetak nama user, PID, dan pemakaian CPU untuk semua proses yang CPU-nya di atas 5%.

Kondisi bisa berupa regex:

```bash
awk '/error/ {print $0}' app.log
```

Setara dengan `grep "error" app.log`, tapi sekarang kamu bisa tambahkan pemrosesan lanjutan.

## Separator Kustom

Gunakan `-F` untuk mengubah pemisah kolom:

```bash
awk -F: '{print $1}' /etc/passwd
```

Mencetak semua nama user dari file passwd yang menggunakan `:` sebagai pemisah.

Untuk CSV sederhana:

```bash
awk -F, '{print $2, $4}' data.csv
```

## Blok BEGIN dan END

Kode di `BEGIN` dijalankan sebelum baris pertama diproses. Kode di `END` dijalankan setelah baris terakhir. Keduanya berguna untuk inisialisasi dan ringkasan.

```bash
awk 'BEGIN {total=0} {total += $3} END {print "Total:", total}' sales.txt
```

Menjumlahkan kolom ketiga dari seluruh file dan mencetak hasilnya di akhir.

## Menghitung Baris

```bash
awk 'END {print NR}' file.txt
```

`NR` adalah variabel built-in yang menyimpan nomor baris saat ini. Di blok `END`, nilainya sama dengan jumlah total baris.

## Contoh Nyata

**Lihat 10 file terbesar di direktori saat ini:**

```bash
ls -la | awk '{print $5, $9}' | sort -n | tail -10
```

**Hitung rata-rata waktu respons dari log Nginx:**

```bash
awk '{sum += $NF; count++} END {print sum/count "ms avg"}' access.log
```

`$NF` merujuk ke kolom terakhir, apapun jumlah kolomnya.

**Ambil semua IP unik dari log:**

```bash
awk '{print $1}' access.log | sort -u
```

**Format ulang output `df` jadi lebih ringkas:**

```bash
df -h | awk 'NR>1 {printf "%-20s %s used\n", $6, $5}'
```

`NR>1` melewati baris pertama (header).

## Gabungan dengan Perintah Lain

Kekuatan `awk` muncul ketika digabungkan dengan pipeline. Tidak perlu menulis skrip Python untuk tugas seperti ini:

```bash
# Proses yang memakan lebih dari 100MB memory
ps aux | awk '$6 > 100000 {print $2, $11}' | head -20

# Port yang sedang dipakai, tanpa duplikat
ss -tlnp | awk 'NR>1 {print $4}' | cut -d: -f2 | sort -un
```

`awk` bukan pengganti Python untuk logika yang kompleks, tapi untuk memfilter, mengubah format, dan merangkum data teks di terminal, tidak ada yang lebih cepat dan lebih portabel.
