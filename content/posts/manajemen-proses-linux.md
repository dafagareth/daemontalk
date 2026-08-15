---
title: "Manajemen Proses di Linux: ps, kill, dan Teman-temannya"
slug: d43f0be5
aliases: [manajemen-proses-linux]
date: 2025-12-12
tags: [linux, cli, tools]
lang: id
draft: false
---

Sebagian besar developer tahu `kill` dan `ps aux`, tapi berhenti di sana. Akibatnya, ketika ada proses yang bermasalah atau perlu dikelola lebih teliti, mereka hanya bisa menebak-nebak atau langsung restart server.

Memahami alat manajemen proses Linux secara menyeluruh membuat debugging jauh lebih cepat.

## ps: Melihat Proses yang Berjalan

```bash
ps aux
```

Output `ps aux` berisi banyak kolom. Yang paling relevan:

| Kolom | Arti |
|---|---|
| `USER` | Pemilik proses |
| `PID` | Process ID |
| `%CPU` | Penggunaan CPU |
| `%MEM` | Penggunaan memori |
| `STAT` | Status proses |
| `COMMAND` | Perintah yang dijalankan |

Kolom `STAT` memberi informasi tambahan:
- `S` = sleeping (menunggu event)
- `R` = running
- `D` = uninterruptible sleep (biasanya I/O, tidak bisa di-kill)
- `Z` = zombie (sudah selesai tapi parent belum mengambil exit code-nya)
- `T` = stopped

Cari proses spesifik tanpa `grep`:

```bash
pgrep -la nginx
```

`pgrep` langsung mengembalikan PID berdasarkan nama, tanpa perlu `ps aux | grep nginx | grep -v grep`.

## kill: Mengirim Sinyal

`kill` bukan hanya untuk mematikan proses. Ia mengirim sinyal, dan sinyal paling umum adalah:

| Sinyal | Nomor | Arti |
|---|---|---|
| `SIGTERM` | 15 | Minta proses berhenti dengan baik (default) |
| `SIGKILL` | 9 | Paksa hentikan, tidak bisa diabaikan |
| `SIGHUP` | 1 | Reload konfigurasi (banyak daemon menggunakannya) |
| `SIGSTOP` | 19 | Pause proses |
| `SIGCONT` | 18 | Lanjutkan proses yang di-pause |

```bash
kill 1234          # kirim SIGTERM ke PID 1234
kill -9 1234       # kirim SIGKILL (paksa)
kill -HUP 1234     # reload konfigurasi
pkill nginx        # kill berdasarkan nama
killall python3    # kill semua proses bernama python3
```

Selalu coba `SIGTERM` dulu sebelum `SIGKILL`. `SIGTERM` memberi kesempatan proses untuk membersihkan resource (tutup file, simpan state). `SIGKILL` memotong semuanya paksa.

## Menjalankan Proses di Background

```bash
./server &          # jalankan di background
jobs               # lihat proses background di sesi ini
fg %1              # bawa kembali ke foreground
bg %1              # lanjutkan di background
```

Masalah: proses yang dijalankan dengan `&` akan ikut mati ketika sesi terminal ditutup.

Gunakan `nohup` supaya proses tetap hidup setelah logout:

```bash
nohup ./server > server.log 2>&1 &
```

`nohup` mengabaikan sinyal `SIGHUP` yang dikirim ketika terminal ditutup. Output diarahkan ke `server.log`.

## nice dan renice: Prioritas CPU

Setiap proses punya nilai "nice" antara -20 (prioritas tertinggi) sampai 19 (prioritas terendah). Nilai default adalah 0.

Jalankan proses dengan prioritas rendah supaya tidak mengganggu proses lain:

```bash
nice -n 10 ./heavy-task
```

Ubah prioritas proses yang sudah berjalan:

```bash
renice -n 15 -p 1234
```

Berguna ketika menjalankan proses berat seperti build kompilasi atau backup di server yang juga melayani traffic.

## lsof: Lihat File yang Dibuka

`lsof` (list open files) menunjukkan file apa saja yang sedang dibuka oleh proses, termasuk koneksi jaringan.

```bash
lsof -p 1234         # semua file yang dibuka oleh PID 1234
lsof -i :8080        # proses mana yang mendengarkan di port 8080
lsof -u deploy       # semua file yang dibuka user deploy
lsof /var/log/app.log  # proses mana yang membuka file ini
```

Berguna ketika kamu ingin tahu mengapa sebuah file tidak bisa dihapus (ada proses yang masih memegangnya).

## Memantau Proses Secara Real-time

```bash
top          # monitor standar
htop         # versi lebih nyaman, perlu diinstall
```

Alternatif modern yang lebih informatif:

```bash
# Arch Linux
pacman -S htop btop

# Debian/Ubuntu
apt install htop
```

`btop` menampilkan grafik CPU, memori, disk, dan jaringan sekaligus dalam satu tampilan terminal.

## Zombie Process

Proses zombie muncul di output `ps` dengan status `Z`. Zombie sudah selesai berjalan tapi entry-nya masih ada di proses tabel karena parent belum memanggil `wait()` untuk mengambil exit code-nya.

Kamu tidak bisa membunuh zombie langsung. Yang bisa dilakukan adalah kill parent-nya:

```bash
ps aux | grep 'Z'           # cari zombie
ps -o ppid= -p <zombie_pid> # cari parent PID
kill <parent_pid>           # kill parent
```

Jika parent tidak bisa di-kill, zombie akan hilang sendiri saat sistem di-restart. Sedikit zombie tidak berbahaya, tapi banyak zombie menandakan ada bug pada program yang tidak menangani child process dengan benar.
