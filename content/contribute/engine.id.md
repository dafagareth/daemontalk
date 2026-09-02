## Pengembangan Core Engine Daemontalk

Daemontalk dibangun dari nol dengan arsitektur modular monolitik di Go, mengutamakan performa tinggi dan ketergantungan minimal.

### Komponen Utama Arsitektur

- **Backend & Router**: Go 1.23+, Chi HTTP Router, driver SQLite murni (`modernc.org/sqlite`), parser Goldmark Markdown.
- **Frontend & Rendering**: A-h Templ (type-safe template generation), Tailwind CSS v4 CLI, dan HTMX.
- **Terminal UI (TUI)**: Charmbracelet Bubble Tea & Lip Gloss, dengan server SSH daemon (`/tuisrv`).
- **Autentikasi & Forum**: GitHub OAuth 2.0, session cookies HMAC-SHA256, dan forum `/socket`.

---

## Alur Kerja Menjalankan Proyek Lokal

```bash
# 1. Clone repositori dan unduh dependensi
git clone https://github.com/dafagareth/daemontalk.git
cd daemontalk
go mod download
npm install

# 2. Kompilasi template Templ & Tailwind CSS
make build

# 3. Jalankan unit test
go test -count=1 ./...

# 4. Jalankan server lokal (akses http://localhost:8080)
./daemontalk
```

---

## Pedoman Pengajuan Kode (*Code Guidelines*)

1. **Format & Linter**: Pastikan kode Go terformat rapi menggunakan `go fmt ./...` atau `golangci-lint`.
2. **Uji Verifikasi**: Selalu sertakan unit test untuk setiap handler HTTP, fungsi storage, atau parser baru di `*_test.go`.
3. **Cabang Git**: Gunakan nama cabang deskriptif seperti `fix/deskripsi-bug` atau `feat/nama-fitur`.
