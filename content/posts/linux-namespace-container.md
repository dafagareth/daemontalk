---
title: "Linux Namespace: Fondasi dari Container Modern"
slug: b8e14d3f
aliases: linux-namespace-container
date: 2026-06-10
tags: [linux, devops, security]
lang: id
draft: false
---

Container bukan teknologi tersendiri, melainkan kumpulan fitur kernel Linux yang disusun bersama. Salah satu pilar utamanya adalah namespace, sebuah mekanisme yang memungkinkan proses melihat tampilan sistem yang terisolasi dari proses lainnya.

## Fakta Menarik

**Fakta 1.** Namespace pertama yang ditambahkan ke kernel Linux adalah `mnt` namespace, masuk pada kernel 2.4.19 tahun 2002. Namespace `pid`, `net`, dan `uts` menyusul bertahun-tahun kemudian, dan `user` namespace baru matang di kernel 3.8 (2013).

**Fakta 2.** Setiap proses di Linux selalu berada dalam satu set namespace. Bahkan proses pada sistem tanpa container tetap berada dalam "initial namespace" yang dibuat oleh kernel saat booting.

**Fakta 3.** Docker tidak mengimplementasikan isolasi sendiri. Ia memanggil syscall `clone(2)` dengan flag namespace yang sesuai, lalu mengonfigurasi cgroup untuk membatasi resource. Containerd dan runc adalah lapisan yang menangani detail ini.

---

## Tips dan Trik

### 1. Tujuh Tipe Namespace di Linux

Kernel Linux saat ini mendukung tujuh jenis namespace, masing-masing mengisolasi aspek sistem yang berbeda:

| Namespace | Flag Clone | Yang Diisolasi |
|---|---|---|
| `pid` | `CLONE_NEWPID` | Pohon proses (PID) |
| `net` | `CLONE_NEWNET` | Antarmuka jaringan, tabel routing, iptables |
| `mnt` | `CLONE_NEWNS` | Titik mount dan filesystem tree |
| `uts` | `CLONE_NEWUTS` | Hostname dan NIS domain name |
| `ipc` | `CLONE_NEWIPC` | Antrian pesan, semaphore, shared memory |
| `user` | `CLONE_NEWUSER` | UID/GID mapping (root di container, bukan di host) |
| `cgroup` | `CLONE_NEWCGROUP` | Tampilan hierarki cgroup |

Setiap baris di `/proc/<pid>/ns/` adalah symlink ke namespace yang saat ini dipakai proses tersebut.

```bash
ls -la /proc/$$/ns/
# lrwxrwxrwx 1 dd dd 0 Jun 10 09:00 net -> net:[4026531840]
# lrwxrwxrwx 1 dd dd 0 Jun 10 09:00 pid -> pid:[4026531836]
```

### 2. Membuat Namespace Manual dengan unshare

Perintah `unshare` menjalankan program dalam namespace baru tanpa memerlukan Docker atau runtime container apapun.

```bash
# Membuat UTS namespace baru dan mengubah hostname di dalamnya
sudo unshare --uts bash -c 'hostname container-test && hostname'
# container-test

# Verifikasi bahwa hostname host tidak berubah
hostname
# workstation (nama hostname asli tetap sama)
```

Untuk membuat PID namespace sekaligus me-mount /proc yang sesuai:

```bash
sudo unshare --pid --fork --mount-proc bash
echo $$
# 1    <- PID proses ini adalah 1 di dalam namespace
ps aux
# Hanya proses di dalam namespace ini yang terlihat
```

### 3. Isolasi Proses Sederhana di Shell

Contoh berikut menggabungkan beberapa namespace untuk membuat lingkungan yang lebih terisolasi, tanpa perlu image container:

```bash
sudo unshare \
  --pid \
  --fork \
  --mount-proc \
  --net \
  --uts \
  bash -c '
    hostname isolated-env
    echo "Hostname: $(hostname)"
    echo "PID shell ini: $$"
    echo "Antarmuka jaringan:"
    ip link show
    echo "Proses yang terlihat:"
    ps aux
  '
```

Output `ip link show` hanya akan menampilkan antarmuka `lo` karena namespace jaringan baru dibuat kosong. Proses dari host tidak akan terlihat di `ps aux`.

### 4. Memeriksa Namespace Suatu Proses

Untuk mengetahui namespace yang digunakan oleh container yang sedang berjalan:

```bash
# Temukan PID dari proses utama container
CONTAINER_PID=$(docker inspect --format '{{.State.Pid}}' nama-container)

# Lihat namespace-nya
ls -la /proc/${CONTAINER_PID}/ns/

# Bandingkan dengan namespace host
ls -la /proc/1/ns/
```

Jika inode pada symlink namespace sama antara container dan host, keduanya berbagi namespace tersebut. Docker secara default mengisolasi `pid`, `net`, `mnt`, `uts`, dan `ipc`, tetapi tidak mengisolasi `user` namespace kecuali dikonfigurasi secara eksplisit dengan `--userns-remap`.

### 5. Bagaimana Docker Menggunakan Namespace Secara Internal

Saat `docker run` dijalankan, runtime (runc) melakukan langkah-langkah berikut:

```
1. clone(CLONE_NEWPID | CLONE_NEWNET | CLONE_NEWNS | CLONE_NEWUTS | CLONE_NEWIPC)
   -> proses anak lahir di namespace baru

2. unshare(CLONE_NEWNS) di dalam anak
   -> pivot_root ke rootfs image (overlay filesystem)

3. Konfigurasi jaringan (veth pair, bridge docker0)
   -> antarmuka eth0 di container terhubung ke veth di host

4. execve("/bin/sh", ...) atau entrypoint yang ditentukan
   -> proses PID 1 di dalam container berjalan
```

Anda dapat memverifikasi veth pair yang dibuat Docker:

```bash
# Di host, lihat antarmuka yang namanya diawali "veth"
ip link show type veth

# Di dalam container
docker exec -it nama-container ip link show
# eth0 akan terlihat dengan IP yang dialokasikan Docker
```

Memahami layer namespace ini berguna saat melakukan debug network connectivity antar container, menganalisis masalah privilege escalation, atau merancang kebijakan keamanan yang tepat untuk workload container di lingkungan produksi.
