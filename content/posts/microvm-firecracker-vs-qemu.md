---
title: "MicroVM Firecracker vs QEMU: Isolasi Multi-Tenant Berkecepatan Tinggi"
slug: f7c2d9a4
aliases: [microvm-firecracker-vs-qemu]
date: 2026-07-10
tags: [linux, devops, performance]
lang: id
draft: false
type: post
---

Isolasi beban kerja multi-tenant pada infrastruktur cloud memerlukan keseimbangan antara keamanan virtualisasi dan efisiensi kontainer. Firecracker hadir sebagai teknologi Virtual Machine Monitor (VMM) berbasis KVM yang mengorbankan kompatibilitas perangkat legacy demi kecepatan booting dan penggunaan memori yang minimal.

## Fun Fact

**Fact 1.** QEMU memiliki lebih dari 1,4 juta baris kode C untuk mendukung emulasi ratusan perangkat keras legacy seperti IDE controller, floppy disk, dan kartu VGA PCI.

**Fact 2.** Firecracker dikembangkan dalam bahasa Rust oleh AWS dengan hanya menyisakan sekitar 50.000 baris kode untuk meminimalkan permukaan serangan (attack surface).

**Fact 3.** Platform serverless seperti AWS Lambda dan Fly.io mengoperasikan jutaan fungsi mikro dalam microVM Firecracker yang dibuat dan dihancurkan dalam hitungan milidetik.

---

## Tips dan Trik

### 1. Arsitektur Virtio Minimalis vs Full Device Emulation

QEMU menyediakan emulasi BIOS komprehensif, bus PCI, serta periferal fisik. Hal ini memungkinkan running OS tanpa modifikasi, namun memberikan overhead inisialisasi yang besar.

Firecracker menanggalkan seluruh emulasi PCI dan BIOS tradisional. Mesin virtual dinyalakan secara langsung menggunakan uncompressed Linux kernel image dan sistem I/O berbasis VirtIO (virtio-net, virtio-block, virtio-vsock, serial console).

### 2. Konfigurasi Kernel untuk Boot Time Kurang dari 10ms

Untuk mencapai waktu boot di bawah 10 milidetik, kernel Linux harus dikompilasi tanpa driver PCI atau modul yang tidak diperlukan. Konfigurasi boot source dikirimkan via REST API socket Firecracker:

```bash
curl --unix-socket /tmp/firecracker.socket -i \
  -X PUT 'http://localhost/boot-source' \
  -H 'Accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "kernel_image_path": "/opt/vmlinux-5.10.bin",
    "boot_args": "console=ttyS0 reboot=k panic=1 pci=off quiet"
  }'
```

### 3. Pengaturan Root Filesystem dan Konsumsi Memori Minimal

Setiap instance Firecracker dapat dikonfigurasi untuk hanya menggunakan memori di bawah 5MB saat diam. Pengaturan drive dilakukan dengan payload JSON berikut:

```bash
curl --unix-socket /tmp/firecracker.socket -i \
  -X PUT 'http://localhost/drives/rootfs' \
  -H 'Accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "drive_id": "rootfs",
    "path_on_host": "/opt/rootfs.ext4",
    "is_root_device": true,
    "is_read_only": false
  }'
```

### 4. Menentukan Alokasi Sumber Daya vCPU dan RAM

Spesifikasi perangkat keras disesuaikan secara dinamis sebelum proses virtualisasi dimulai:

```bash
curl --unix-socket /tmp/firecracker.socket -i \
  -X PUT 'http://localhost/machine-config' \
  -H 'Accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "vcpu_count": 1,
    "mem_size_mib": 128,
    "ht_enabled": false
  }'
```

### 5. Memulai Virtual Machine dan Eksekusi Ekosistem Serverless

Setelah seluruh parameter terkonfigurasi, pemicuan perintah `InstanceStart` akan mengeksekusi kernel secara instan:

```bash
curl --unix-socket /tmp/firecracker.socket -i \
  -X PUT 'http://localhost/actions' \
  -H 'Accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "action_type": "InstanceStart"
  }'
```

Pola ini memungkinkan infrastruktur cloud seperti Fly.io menjalankan sandbox pengguna yang terisolasi ketat dengan batas keamanan tingkat kernel Linux (KVM) tanpa penalti kinerja kontainer konvensional.
