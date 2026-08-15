---
title: "Membangun Server DERP Mandiri untuk Jaringan Tailscale"
slug: c9d4e2a8
aliases: [tailscale-derp-custom-relay]
date: 2026-05-18
tags: [networking, tools, devops]
lang: id
draft: false
type: post
---

Tailscale mengandalkan jaringan mesh terenkripsi peer-to-peer (P2P) berbasis WireGuard. Namun, ketika topologi NAT simetris atau firewall ketat memblokir koneksi UDP langsung, Tailscale beralih ke relay DERP (Designated Encrypted Relay for Packets) untuk memastikan komunikasi data tetap terhubung.

## Fun Fact

**Fact 1.** Protokol WireGuard standar membutuhkan setidaknya satu node dengan IP publik statis dan port UDP terbuka untuk mengawali pertukaran kunci.

**Fact 2.** Tailscale memanfaatkan teknik Interactive Connectivity Establishment (ICE) dan probing paket STUN untuk melewati NAT tanpa port forwarding manual.

**Fact 3.** Paket data yang melewati node DERP tetap terenkripsi end-to-end dengan kunci WireGuard klien, sehingga pengelola server DERP tidak dapat mengintip isi lalu lintas data.

---

## Tips dan Trik

### 1. Mekanisme NAT Traversal dan Kasus Kegagalan P2P Direct

Tailscale menggunakan STUN (Session Traversal Utilities for NAT) untuk mengidentifikasi IP dan port publik dari setiap node. Jika kedua peer berada di belakang NAT moderat, koneksi P2P UDP dapat terbentuk secara langsung.

Namun, jika salah satu peer berada di bawah NAT simetris (di mana port publik berubah untuk setiap endpoint tujuan), koneksi UDP P2P gagal terbentuk. Pada skenario ini, paket disalurkan melalui HTTP/2 terenkripsi via server DERP.

### 2. Kompilasi dan Instalasi Binary derper di VPS Linux

Jalankan perintah berikut di VPS Linux dengan IP publik statis untuk mengompilasi binary `derper`:

```bash
# Menginstal derper dari pustaka resmi Tailscale
go install tailscale.com/cmd/derper@latest

# Memverifikasi hasil instalasi binary
/home/user/go/bin/derper -help
```

### 3. Konfigurasi Layanan Systemd untuk Server DERP Mandiri

Buat unit file `/etc/systemd/system/derper.service` untuk mengelola proses relay:

```ini
[Unit]
Description=Tailscale Custom DERP Server
After=network.target

[Service]
User=root
ExecStart=/root/go/bin/derper \
  -hostname=derp.example.com \
  -a=:443 \
  -http-port=80 \
  -stun-port=3478 \
  -verify-clients=true
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Jalankan daemon service dengan perintah:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now derper
```

### 4. Menambahkan Node DERP Pribadi pada ACL Map Admin Console

Buka Tailscale Admin Console pada menu Access Control, kemudian tambahkan blok `derpMap` berikut:

```json
{
  "derpMap": {
    "OmitDefaultRegions": false,
    "Regions": {
      "901": {
        "RegionID": 901,
        "RegionCode": "custom-vps",
        "RegionName": "Custom VPS Relay",
        "Nodes": [
          {
            "Name": "derp-node-1",
            "RegionID": 901,
            "HostName": "derp.example.com",
            "STUNPort": 3478,
            "DERPPort": 443
          }
        ]
      }
    }
  }
}
```

### 5. Verifikasi Status Koneksi dan Pengujian Latensi Node

Gunakan CLI Tailscale untuk memeriksa ketersediaan node DERP kustom dan latensi antar titik:

```bash
# Memeriksa latensi ke seluruh node STUN dan DERP
tailscale netcheck

# Menampilkan status jalur koneksi peer aktif
tailscale status
```

Output `tailscale status` akan menampilkan label `relay` jika koneksi melewati server DERP atau `direct` jika koneksi P2P UDP berhasil ditegakkan.
