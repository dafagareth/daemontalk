## Koreksi Artikel & Perbaikan Cuplikan Kode

Seluruh artikel teknis di Daemontalk adalah dokumen sumber terbuka yang tersimpan di dalam direktori `content/posts/`. Kami menyambut perbaikan kesalahan ketik (*typo*), perbaikan *bug* pada cuplikan kode, atau pembaruan tautan rujukan.

---

## 2 Cara Melakukan Koreksi

### Jalur 1: 1-Klik Edit via GitHub Web
Pada artikel yang sedang Anda baca, klik icon pensil **Edit Suggest** di baris aksi bawah artikel. Tautan akan langsung membuka editor teks GitHub untuk mengajukan Pull Request secara instan.

### Jalur 2: Melalui Git Branch Lokal
```bash
git checkout -b docs/fix-nama-artikel
# Sunting berkas markdown di content/posts/nama-artikel.md
git commit -m "docs: fix typo in memory allocation section"
git push origin docs/fix-nama-artikel
```

---

## Atribusi Resmi Kontributor

Ketika Pull Request koreksi Anda digabungkan (*merged*), nama pengguna GitHub Anda akan dicantumkan secara resmi di *byline* header artikel sebagai **CONTRIBUTOR** lengkap dengan foto profil akun GitHub Anda.
