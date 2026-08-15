---
title: "Generate OG Image di Go Tanpa Headless Browser"
slug: b5a4194b
aliases: [og-image-generator-go]
date: 2026-06-06
tags: [go, web, image]
lang: id
draft: false
---

Setiap kali sebuah link blog dibagikan ke WhatsApp, X, atau LinkedIn, yang muncul adalah sebuah kartu preview dengan gambar besar. Gambar itu namanya **OG image** (Open Graph), dan kalau kamu tidak menyediakannya, link-mu cuma jadi teks polos yang gampang dilewati.

Banyak orang membuatnya pakai headless browser (Puppeteer, Playwright) yang me-render HTML jadi gambar. Cara itu jalan, tapi berat: butuh Chromium ratusan MB, RAM besar, dan startup yang lambat. Untuk situs Go satu binary, itu terasa berlebihan.

Saya pilih jalan lain: **gambar digambar langsung di Go**, murni pakai standard library + satu paket font.

## Idenya

Satu fungsi menerima judul post, lalu menggambar sebuah PNG 1200×630 yang memuat gradient gelap, accent bar, judul, subtitle, dan nama situs. Tidak ada browser, tidak ada file font eksternal (Go punya font bawaan yang bisa di-embed).

```go
type Card struct {
	Title    string
	Subtitle string
	Site     string
}

func Render(w io.Writer, c Card) error {
	img := image.NewRGBA(image.Rect(0, 0, 1200, 630))
	drawVerticalGradient(img, bgTop, bgBot)
	// ... gambar accent bar, judul, subtitle
	return png.Encode(w, img)
}
```

## Menggambar teks

Bagian paling tricky adalah teks. Paket `golang.org/x/image/font` menyediakan `font.Drawer` untuk menulis string ke kanvas:

```go
d := &font.Drawer{
	Dst:  img,
	Src:  image.NewUniform(color.White),
	Face: face,
	Dot:  fixed.P(x, y),
}
d.DrawString("Halo dunia")
```

Fontnya dari `gofont/gobold`, yaitu TTF yang sudah ter-embed di dalam paket, jadi tidak perlu menyertakan file `.ttf` apa pun.

## Word wrap manual

Judul yang panjang harus dipecah jadi beberapa baris agar tidak terpotong. Tidak ada CSS di sini, jadi saya hitung lebar tiap kata sendiri dengan `font.MeasureString`, lalu greedy-wrap:

```go
for _, word := range words[1:] {
	test := line + " " + word
	if textWidth(face, test) > maxW {
		lines = append(lines, line)
		line = word
	} else {
		line = test
	}
}
```

Di sinilah saya kena bug yang menarik. Kalau judul terlalu panjang, baris terakhir saya potong dengan elipsis `…`. Versi pertama saya tulis begini:

```go
for len(s) > 1 && textWidth(face, s) > maxW {
	s = s[:len(s)-2] + "…"
}
```

Kelihatan benar, tapi `…` itu **3 byte** dalam UTF-8. Memotong 2 byte lalu menambah `…` lagi malah merusak karakter dan membuat string makin panjang sehingga terjadi **infinite loop**. Kalau ada satu saja post berjudul super panjang, server-nya hang.

Perbaikannya: bekerja per-rune, bukan per-byte, dan pastikan ada kondisi berhenti:

```go
func truncate(face font.Face, s string, maxW int) string {
	runes := []rune(strings.TrimSuffix(s, "…"))
	for len(runes) > 0 {
		candidate := string(runes) + "…"
		if textWidth(face, candidate) <= maxW {
			return candidate
		}
		runes = runes[:len(runes)-1]
	}
	return "…"
}
```

Pelajaran lama yang selalu relevan: di Go, `string` itu byte, bukan karakter. Begitu menyentuh teks non-ASCII, pakai `[]rune`.

## Caching

Render PNG butuh beberapa milidetik, dan isinya tidak berubah sampai post di-update. Jadi hasilnya saya simpan di map in-memory dengan kunci slug:

```go
ogMu.RLock()
cached, hit := ogCache[slug]
ogMu.RUnlock()
if !hit {
	// render sekali, lalu simpan
}
```

Request kedua dan seterusnya langsung melayani byte dari memori. Ditambah header `Cache-Control: public, max-age=86400`, gambar yang sama tidak pernah dirender dua kali.

## Hasilnya

Setiap post sekarang punya OG image otomatis di `/blog/{slug}/og.png`. Post yang tidak menyediakan cover sendiri tetap dapat kartu yang rapi untuk dibagikan tanpa Chromium, tanpa file font, dan tanpa dependency berat. Cukup standard library dan sedikit matematika.

Bug infinite-loop tadi? Ketahuan oleh sebuah unit test sederhana yang menguji judul ekstra panjang. Itu pengingat kenapa test bukan formalitas: ia menangkap hal yang tidak terpikirkan saat menulis kode.
