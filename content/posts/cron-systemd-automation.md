---
title: "Automasi Script Linux: Best Practice Logging & Error Trap"
slug: 7a8b9c0d
aliases: [automasi-script-linux]
date: 2026-06-25
tags: [bash, linux, devops]
lang: id
draft: false
---

Menjalankan script otomatisasi di server (baik via Cron maupun Systemd Timers) sering kali menjadi mimpi buruk saat script gagal diam-diam tanpa meninggalkan jejak log yang jelas.

Dengan menerapkan prinsip defensive scripting dan logging terstruktur, kamu bisa tidur nyenyak mengetahui status setiap cron job.

## Fun Fact

**Nama `cron` berasal dari kata Yunani *Chronos* yang berarti Waktu.**
Program cron pertama kali diimplementasikan oleh Brian Kernighan di Version 7 Unix pada tahun 1979.

**Cron Environment sangat minimalis.**
Banyak script cron gagal bukan karena logic error, melainkan karena `PATH` cron default hanya mencakup `/usr/bin` dan `/bin`, tanpa variabel lingkungan shell user seperti `.bashrc`.

**Fitur `trap` di Bash memungkinkan pembersihan file temporary otomatis saat script crash.**
Sinyal `EXIT`, `SIGINT`, atau `SIGTERM` bisa dicegat untuk mengeksekusi fungsi cleanup.

---

## Tips dan Trik

### 1. Gunakan Header Bash Ketat: `set -euo pipefail`

Pastikan script langsung berhenti saat ada error, variabel kosong tak terdefinisi, atau kegagalan di dalam rantai pipeline pipa:

```bash
#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'
```

### 2. Trap Cleanup untuk Menghapus File Temporary

Hapus file sementara apa pun status keluar dari script:

```bash
TMP_DIR=$(mktemp -d)
cleanup() {
    rm -rf "$TMP_DIR"
    echo "[INFO] Temporary files cleaned up."
}
trap cleanup EXIT INT TERM
```

### 3. Log Output Langsung ke Systemd Journal / Syslog

Gunakan utilitas `logger` agar log script otomatis masuk ke journald dengan level severity yang tepat:

```bash
logger -t "backup-db" -p local0.info "Database backup started at $(date)"
```

### 4. Gunakan `flock` untuk Mencegah Eksekusi Ganda (Overlap)

Cegah script cron yang lama berjalan agar tidak bertubrukan dengan cron iterasi berikutnya:

```bash
# Di dalam crontab:
* * * * * flock -n /var/lock/myjob.lock /usr/local/bin/myjob.sh
```

### 5. Selalu Definisikan Path Absolut

Hindari ketergantungan pada current working directory saat automasi:

```bash
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_FILE="${SCRIPT_DIR}/config.json"
```
