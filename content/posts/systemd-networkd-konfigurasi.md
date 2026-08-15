---
title: "Menggunakan systemd-networkd sebagai Pengganti NetworkManager di Server Linux"
slug: b8d2f6a3
aliases: systemd-networkd-konfigurasi
date: 2026-05-14
tags: [linux, networking, devops]
lang: id
draft: false
---

NetworkManager dirancang untuk lingkungan desktop dengan kebutuhan konektivitas yang dinamis: berpindah antara Wi-Fi, VPN, dan hotspot. Di server Linux, kebutuhan tersebut hampir tidak pernah ada. `systemd-networkd` menawarkan alternatif yang lebih ringan, berbasis file konfigurasi, dan terintegrasi langsung dengan ekosistem systemd.

## Fakta Menarik

**Fakta 1.** `systemd-networkd` tidak memiliki daemon tambahan selain unit systemd itu sendiri. NetworkManager sebaliknya menjalankan beberapa proses dan bergantung pada D-Bus untuk hampir semua operasinya.

**Fakta 2.** File konfigurasi `.network` dan `.netdev` bersifat deklaratif dan dapat di-version control. Ini membuat konfigurasi jaringan server menjadi dapat direproduksi dan mudah diaudit.

**Fakta 3.** `networkctl` adalah alat diagnostik bawaan yang memberikan informasi status antarmuka secara real-time tanpa memerlukan pemasangan paket tambahan.

---

## Tips dan Trik

### 1. Migrasi dari NetworkManager ke systemd-networkd

Sebelum beralih, catat konfigurasi jaringan yang aktif, lalu matikan NetworkManager dan aktifkan `systemd-networkd`.

```bash
# Catat konfigurasi saat ini
ip addr show
ip route show
cat /etc/resolv.conf

# Nonaktifkan NetworkManager
sudo systemctl disable --now NetworkManager

# Aktifkan systemd-networkd dan systemd-resolved
sudo systemctl enable --now systemd-networkd
sudo systemctl enable --now systemd-resolved

# Arahkan resolv.conf ke stub resolver systemd-resolved
sudo ln -sf /run/systemd/resolve/stub-resolv.conf /etc/resolv.conf
```

### 2. Konfigurasi Antarmuka Statis dengan File .network

Semua file konfigurasi jaringan disimpan di `/etc/systemd/network/`. Urutan pemuatan ditentukan oleh nama file secara leksikografis.

```ini
# /etc/systemd/network/10-eth0.network

[Match]
Name=eth0

[Network]
Address=192.168.1.10/24
Gateway=192.168.1.1
DNS=1.1.1.1
DNS=8.8.8.8

[Link]
MTUBytes=1500
```

```bash
# Terapkan konfigurasi tanpa reboot
sudo networkctl reload

# Verifikasi status antarmuka
networkctl status eth0
```

### 3. Konfigurasi VLAN

VLAN memerlukan dua file: satu `.netdev` untuk membuat antarmuka virtual, dan satu `.network` untuk mengonfigurasinya.

```ini
# /etc/systemd/network/20-vlan100.netdev

[NetDev]
Name=eth0.100
Kind=vlan

[VLAN]
Id=100
```

```ini
# /etc/systemd/network/20-vlan100.network

[Match]
Name=eth0.100

[Network]
Address=10.100.0.5/24
Gateway=10.100.0.1
```

```ini
# /etc/systemd/network/10-eth0.network
# Antarmuka induk harus mendeklarasikan VLAN yang dimilikinya

[Match]
Name=eth0

[Network]
VLAN=eth0.100
```

```bash
sudo networkctl reload
networkctl list
```

### 4. Bonding Dua Antarmuka Jaringan

Bonding menggabungkan dua atau lebih antarmuka fisik menjadi satu antarmuka logis untuk redundansi atau peningkatan bandwidth.

```ini
# /etc/systemd/network/30-bond0.netdev

[NetDev]
Name=bond0
Kind=bond

[Bond]
Mode=active-backup
MIIMonitorSec=100ms
```

```ini
# /etc/systemd/network/30-bond0.network

[Match]
Name=bond0

[Network]
Address=192.168.1.20/24
Gateway=192.168.1.1
DNS=1.1.1.1
```

```ini
# /etc/systemd/network/31-eth1.network

[Match]
Name=eth1

[Network]
Bond=bond0
```

```ini
# /etc/systemd/network/32-eth2.network

[Match]
Name=eth2

[Network]
Bond=bond0
```

```bash
sudo networkctl reload

# Verifikasi status bonding
cat /proc/net/bonding/bond0
networkctl status bond0
```

### 5. Debugging dengan networkctl

`networkctl` memberikan informasi terstruktur tentang status setiap antarmuka dan dapat mengidentifikasi masalah konfigurasi.

```bash
# Tampilkan semua antarmuka dan statusnya
networkctl list

# Detail lengkap satu antarmuka
networkctl status eth0

# Ikuti log systemd-networkd secara real-time
journalctl -u systemd-networkd -f

# Paksa reload semua konfigurasi tanpa restart
sudo networkctl reload

# Aktifkan atau nonaktifkan antarmuka secara manual
sudo networkctl up eth0
sudo networkctl down eth0

# Verifikasi konfigurasi DNS dari systemd-resolved
resolvectl status
resolvectl query archlinux.org
```

Jika sebuah antarmuka berstatus `configuring` dalam waktu lama, periksa apakah ada konflik nama file `.network` atau kesalahan sintaks dengan `systemd-analyze verify /etc/systemd/network/*.network`.
