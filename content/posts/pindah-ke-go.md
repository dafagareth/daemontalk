---
title: "Pindah ke Go Setelah Bertahun-tahun di PHP dan Python"
slug: d41cabe9
aliases: [pindah-ke-go]
date: 2025-10-11
tags: [go, programming]
lang: id
draft: false
---

Sebagian besar developer di Indonesia memulai dengan PHP atau Python. PHP karena ekosistem Laravel dan banyaknya proyek web, Python karena dipakai di mata kuliah dan terasa mudah dibaca. Keduanya bahasa yang baik dan akan tetap relevan lama.

Tapi pada satu titik, banyak orang mulai mendengar tentang Go, melihat lowongan yang memintanya, dan bertanya-tanya apakah layak dipelajari. Setelah beberapa waktu memakainya, jawaban saya: ya, dan transisinya lebih mulus dari yang dibayangkan, asalkan kamu siap dengan beberapa perbedaan cara berpikir.

## Apa yang Berbeda Sejak Awal

Hal pertama yang terasa adalah Go itu **dikompilasi dan bertipe statis**. Setelah terbiasa dengan PHP dan Python yang langsung jalan tanpa kompilasi, ini terdengar seperti langkah mundur. Ternyata sebaliknya.

Kompilasi menangkap sebagian besar kesalahan sebelum program dijalankan. Salah ketik nama variabel, mengirim tipe data yang salah ke fungsi, lupa menangani sebuah kasus, semuanya ditolak compiler sebelum kode sampai ke pengguna. Di PHP dan Python, kesalahan seperti ini sering baru ketahuan saat program sudah berjalan, kadang di produksi.

Tipe statis juga membuat kode lebih mudah dibaca orang lain. Saat sebuah fungsi menyatakan ia menerima `string` dan mengembalikan `int`, tidak ada keraguan. Tidak perlu menelusuri seluruh fungsi untuk menebak tipe yang diharapkan.

## Penanganan Error yang Eksplisit

Inilah perbedaan yang paling mengejutkan pendatang baru, dan paling sering dikeluhkan di awal. Go tidak punya exception. Sebagai gantinya, fungsi mengembalikan error sebagai nilai biasa yang harus kamu periksa.

```go
file, err := os.Open("data.txt")
if err != nil {
    return fmt.Errorf("gagal membuka file: %w", err)
}
defer file.Close()
```

Pola `if err != nil` ini muncul di mana-mana. Pendatang dari Python yang terbiasa dengan `try/except` awalnya merasa ini bertele-tele. Kenapa harus menulis tiga baris untuk sesuatu yang di Python cukup satu blok?

Tapi setelah beberapa minggu, banyak orang berubah pikiran. Penanganan error yang eksplisit memaksa kamu memikirkan apa yang terjadi saat sesuatu gagal, tepat di tempat kegagalan itu mungkin terjadi. Tidak ada error yang diam-diam terlempar entah ke mana lalu menjatuhkan seluruh program. Kamu melihat setiap jalur kegagalan dengan jelas. Kode menjadi lebih membosankan untuk ditulis, tapi jauh lebih bisa diprediksi saat berjalan.

## Tidak Ada Keajaiban

PHP dengan Laravel dan Python dengan Django penuh dengan "keajaiban": kode yang bekerja di balik layar lewat konvensi, refleksi, dan abstraksi yang dalam. Praktis, tapi kadang membuat sulit memahami apa yang sebenarnya terjadi.

Go sengaja menghindari ini. Tidak ada konstruktor tersembunyi, tidak ada anotasi ajaib, tidak ada metaprogramming yang rumit. Kode Go cenderung melakukan persis seperti yang tertulis. Membaca kode Go orang lain umumnya lebih mudah karena tidak ada lapisan sihir yang harus dipahami lebih dulu.

Konsekuensinya, Go lebih verbose. Kamu menulis lebih banyak baris untuk hal yang di framework lain cukup satu pemanggilan. Sebagian orang membenci ini, sebagian menghargainya. Pandangan saya, kejelasan biasanya lebih berharga daripada keringkasan, terutama saat sebuah proyek dirawat bertahun-tahun oleh banyak orang.

## Konkurensi yang Terasa Alami

Salah satu alasan utama orang beralih ke Go adalah kemudahannya menjalankan banyak hal secara bersamaan. Di PHP, konkurensi sebenarnya bukan hal yang lazim. Di Python, ada keterbatasan terkenal yang membuat paralelisme sejati rumit.

Go dibangun dengan konkurensi sebagai bagian inti bahasanya. Menjalankan fungsi secara bersamaan cukup dengan kata kunci `go`:

```go
go kirimEmail(user)
go catatLog(aktivitas)
```

Untuk komunikasi antar proses yang berjalan bersamaan, Go memakai channel, sebuah cara yang aman dan terstruktur untuk mengoper data. Detailnya butuh pembelajaran tersendiri, tapi yang penting di sini: hal yang di bahasa lain memerlukan pustaka tambahan dan kehati-hatian ekstra, di Go terasa menyatu dengan bahasanya.

## Satu Binary, Tanpa Drama Deployment

Bagian yang langsung terasa manfaatnya saat deployment. Program Go dikompilasi menjadi satu file binary yang berdiri sendiri. Tidak ada interpreter yang harus dipasang di server, tidak ada `composer install` atau `pip install` dengan dependensi yang kadang bentrok, tidak ada perbedaan versi runtime antara mesin lokal dan server.

```bash
# Kompilasi untuk Linux dari mesin mana pun
GOOS=linux GOARCH=amd64 go build -o app

# Salin satu file ini ke server, lalu jalankan
./app
```

Bagi siapa pun yang pernah menghabiskan sore menyelesaikan masalah versi PHP atau dependensi Python yang tidak cocok di server, kesederhanaan ini terasa seperti kelegaan.

## Apa yang Akan Kamu Rindukan

Jujur soal kekurangannya. Ekosistem Go lebih kecil dibanding PHP atau Python. Untuk kebutuhan tertentu, kamu mungkin tidak menemukan pustaka yang sematang yang ada di ekosistem lain, dan harus menulisnya sendiri. Go juga tidak cocok untuk semua hal, untuk scripting cepat atau analisis data, Python tetap pilihan yang lebih masuk akal.

Dan verbositas itu nyata. Jika kamu menikmati keringkasan dan ekspresivitas, Go kadang terasa kaku. Ini trade-off yang disengaja oleh perancangnya, bukan kekurangan yang tidak disadari.

---

Pindah ke Go bukan soal meninggalkan PHP atau Python, keduanya tetap berharga dan punya tempatnya. Ini soal menambah alat yang unggul di area tertentu: layanan yang butuh performa, sistem yang menangani banyak hal sekaligus, dan apa pun yang deploymentnya ingin kamu buat sesederhana mungkin. Transisinya menuntut penyesuaian cara berpikir, terutama soal penanganan error dan menerima kode yang lebih panjang demi kejelasan. Tapi bagi banyak orang, termasuk saya, penyesuaian itu sepadan. Go tidak berusaha mengesankan dengan kepintaran. Ia berusaha agar mudah dipahami, dan dalam jangka panjang, itu yang lebih penting.
