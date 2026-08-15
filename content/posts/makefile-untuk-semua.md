---
title: "Makefile Bukan Cuma untuk C"
slug: fd41195b
aliases: [makefile-untuk-semua]
date: 2025-04-20
tags: [make, tools, automation]
lang: id
draft: false
---

Banyak orang mengasosiasikan `make` dengan compile program C era lawas. Begitu mendengar "Makefile", terbayang proyek lama penuh konfigurasi rumit. Padahal `make` adalah salah satu alat otomasi paling sederhana dan tersedia di mana saja, dan ia sama bergunanya untuk proyek web, Python, Go, atau apa pun yang punya perintah berulang.

Inti `make` cuma satu: kamu memberi nama pada sekumpulan perintah, lalu menjalankannya dengan nama itu. Tidak lebih rumit dari itu.

## Masalah yang Dipecahkan

Setiap proyek punya perintah yang sering diketik. Menjalankan server, menjalankan test, membangun aset, membersihkan file sementara. Perintahnya panjang dan mudah lupa:

```bash
docker compose -f docker-compose.dev.yml up --build
./node_modules/.bin/tailwindcss -i src/input.css -o dist/output.css --minify
go test ./... -count=1 -race
```

Tanpa pencatatan, perintah ini tersebar di README, di riwayat shell, atau di ingatan. Anggota tim baru harus bertanya "cara jalanin ini gimana?". Makefile menjadikannya satu tempat yang konsisten.

## Makefile Pertama

Buat file bernama `Makefile` di akar proyek. Setiap blok perintah disebut **target**.

```makefile
dev:
	docker compose -f docker-compose.dev.yml up --build

test:
	go test ./... -count=1 -race

css:
	./node_modules/.bin/tailwindcss -i src/input.css -o dist/output.css --minify
```

Sekarang perintah panjang tadi cukup dijalankan dengan:

```bash
make dev
make test
make css
```

Ada satu hal yang wajib diperhatikan: **indentasi di Makefile harus berupa tab, bukan spasi.** Ini sumber kebingungan paling umum bagi pemula. Jika kamu memakai spasi, `make` akan mengeluarkan error seperti "missing separator". Pastikan editormu memakai tab untuk file ini.

## Target yang Saling Bergantung

Target bisa bergantung pada target lain. Tuliskan dependensinya setelah titik dua, dan `make` akan menjalankannya lebih dulu secara berurutan.

```makefile
build: css test
	go build -o app .
```

Menjalankan `make build` akan menjalankan `css`, lalu `test`, baru kemudian perintah build. Jika test gagal, proses berhenti dan build tidak dijalankan. Ini cara ringkas merangkai langkah yang berurutan tanpa menulis script terpisah.

## Target .PHONY

Secara historis, `make` dirancang untuk membuat file. Target bernama `build` membuatnya berasumsi akan ada file bernama `build` sebagai hasilnya. Jika kebetulan ada file atau direktori dengan nama yang sama dengan target, `make` bisa bingung dan mengira pekerjaan sudah selesai.

Solusinya adalah mendeklarasikan target sebagai `.PHONY`, menandakan bahwa ini perintah, bukan file.

```makefile
.PHONY: dev test css build clean

dev:
	docker compose up --build

clean:
	rm -rf dist/ app
```

Kebiasaan baik adalah selalu mencantumkan `.PHONY` untuk target yang bukan menghasilkan file dengan nama tersebut. Pada proyek non-C, ini berarti hampir semua target.

## Variabel dan Argumen

Makefile mendukung variabel, berguna untuk nilai yang dipakai berulang.

```makefile
BINARY := app
PORT := 8080

run: build
	PORT=$(PORT) ./$(BINARY)

build:
	go build -o $(BINARY) .
```

Variabel ditulis dengan `$(NAMA)`. Mengubah nama binary atau port cukup dilakukan di satu tempat di atas, dan seluruh target yang memakainya ikut berubah.

## Bantuan yang Mendokumentasikan Dirinya

Saat Makefile sudah punya banyak target, akan berguna jika ada perintah yang menampilkan daftarnya. Pola ini umum dipakai:

```makefile
.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "%-12s %s\n", $$1, $$2}'

dev: ## Jalankan server pengembangan
	docker compose up --build

test: ## Jalankan seluruh test
	go test ./...

clean: ## Hapus artefak build
	rm -rf dist/ app
```

Dengan menambahkan komentar berawalan `##` di setiap target, `make help` akan menampilkan daftar rapi beserta penjelasannya. Tanda `@` di depan perintah membuat `make` tidak mencetak perintahnya sendiri ke layar, hanya menjalankannya.

```bash
$ make help
dev          Jalankan server pengembangan
test         Jalankan seluruh test
clean        Hapus artefak build
```

Makefile berubah menjadi dokumentasi yang hidup. Siapa pun yang membuka proyek cukup mengetik `make help` untuk tahu apa yang bisa dilakukan.

---

`make` sudah berusia puluhan tahun dan tersedia di hampir setiap sistem Unix tanpa instalasi tambahan. Ia tidak butuh konfigurasi, tidak butuh dependency, dan sintaksnya untuk kebutuhan sederhana bisa dipelajari dalam sepuluh menit. Untuk proyek apa pun yang punya sekumpulan perintah berulang, sebuah Makefile kecil menggantikan README yang panjang, mengakhiri pertanyaan "cara jalanin ini gimana", dan memberi seluruh tim satu cara yang sama untuk bekerja. Tidak perlu menulis satu baris C pun untuk memanfaatkannya.
