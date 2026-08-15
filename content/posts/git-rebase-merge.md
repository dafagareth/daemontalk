---
title: "Rebase atau Merge? Memahami Perbedaannya Tanpa Takut"
slug: f8c4fc6d
aliases: [git-rebase-merge]
date: 2024-11-05
tags: [git, workflow]
lang: id
draft: false
---

Banyak developer memakai Git selama bertahun-tahun tanpa benar-benar paham perbedaan `merge` dan `rebase`. Mereka mengikuti satu pola yang pernah berhasil, lalu menghindari yang lain karena pernah membaca cerita horor soal "rebase yang merusak history". Akibatnya separuh kemampuan Git tidak pernah dipakai.

Padahal keduanya menyelesaikan masalah yang sama, hanya dengan hasil yang berbeda. Memahami perbedaan itu menghilangkan rasa takut.

## Masalah yang Sama

Bayangkan kamu membuat branch `fitur` dari `main`, lalu mengerjakannya selama beberapa hari. Sementara itu, rekan tim menambahkan commit baru ke `main`. Sekarang branch kamu tertinggal.

```
main:    A---B---C---D
                  \
fitur:             E---F
```

Branch `fitur` bercabang dari commit `B`, sementara `main` sudah maju ke `D`. Kamu perlu menggabungkan pekerjaan ini. Di sinilah `merge` dan `rebase` mengambil jalan berbeda.

## Merge: Menggabungkan dengan Jejak

`git merge` membuat sebuah commit baru yang menyatukan kedua riwayat.

```bash
git checkout fitur
git merge main
```

Hasilnya:

```
main:    A---B---C---D
                  \   \
fitur:             E---F---M
```

Commit `M` adalah merge commit. Ia punya dua "orang tua": commit terakhir branch `fitur` (`F`) dan commit terakhir `main` (`D`). Riwayat asli kedua branch tetap utuh persis seperti yang terjadi.

Kelebihannya jelas: merge tidak mengubah apa pun, ia hanya menambah. Commit yang sudah ada tidak disentuh. Ini aman dan jujur, history mencerminkan apa yang benar-benar terjadi, termasuk fakta bahwa cabang ini sempat berjalan paralel.

Kekurangannya juga jelas: setelah banyak merge, grafik riwayat menjadi penuh dengan commit `M` dan garis yang saling silang. Pada proyek dengan banyak kontributor, `git log --graph` bisa terlihat seperti diagram kabel yang kusut.

## Rebase: Menulis Ulang Cerita

`git rebase` mengambil pendekatan berbeda. Alih-alih menggabungkan, ia memindahkan commit branch kamu seolah-olah dibuat dari titik terbaru.

```bash
git checkout fitur
git rebase main
```

Hasilnya:

```
main:    A---B---C---D
                      \
fitur:                 E'---F'
```

Perhatikan commit `E` dan `F` berubah menjadi `E'` dan `F'`. Git membuat ulang commit tersebut di atas `D`. Isinya sama, tapi mereka adalah commit baru dengan hash berbeda. Riwayat sekarang menjadi garis lurus, seolah branch `fitur` memang dibuat setelah `main` mencapai `D`.

Kelebihannya: history bersih dan linear. `git log` terbaca seperti urutan kronologis yang rapi, tanpa percabangan yang membingungkan.

Kekurangannya: kamu menulis ulang commit. Dan di sinilah letak satu aturan yang tidak boleh dilanggar.

## Aturan Emas Rebase

**Jangan pernah rebase commit yang sudah dipush dan dipakai orang lain.**

Karena rebase membuat commit baru dengan hash berbeda, branch kamu menjadi tidak cocok dengan versi yang sudah ada di remote. Jika orang lain sudah menarik branch tersebut, mereka akan menghadapi konflik yang kacau karena Git melihat dua riwayat yang berbeda untuk pekerjaan yang sama.

Aturan praktisnya: rebase aman dilakukan pada commit lokal yang belum pernah kamu bagikan. Begitu sebuah commit sudah dipush dan kemungkinan dipakai tim, perlakukan ia sebagai sesuatu yang permanen. Gunakan merge.

## Kapan Memakai yang Mana

Pola yang banyak dipakai tim adalah kombinasi keduanya.

Sebelum membuat pull request, rebase branch lokal kamu ke `main` terbaru untuk merapikan riwayat dan menyelesaikan konflik lebih awal:

```bash
git checkout fitur
git rebase main
# selesaikan konflik jika ada
git push --force-with-lease
```

Perhatikan `--force-with-lease`, bukan `--force` biasa. Opsi ini menolak push jika ada perubahan di remote yang belum kamu lihat, mencegah kamu menimpa pekerjaan orang lain secara tidak sengaja.

Lalu saat menggabungkan ke `main`, banyak tim memilih merge agar tercatat kapan sebuah fitur masuk:

```bash
git checkout main
git merge --no-ff fitur
```

Opsi `--no-ff` memaksa pembuatan merge commit meski sebenarnya bisa fast-forward, sehingga jejak penggabungan fitur tetap terlihat di history.

## Rebase Interaktif untuk Merapikan

Salah satu kegunaan rebase yang paling berguna adalah membersihkan commit sebelum dibagikan.

```bash
git rebase -i HEAD~4
```

Perintah ini membuka editor berisi empat commit terakhir, di mana kamu bisa menggabungkan commit kecil ("fix typo", "fix typo lagi") menjadi satu, mengubah urutan, atau memperbaiki pesan commit. Hasilnya adalah riwayat yang menceritakan pekerjaan dengan jelas, bukan setiap langkah ragu-ragu di sepanjang jalan.

---

Merge dan rebase bukan pilihan benar atau salah, melainkan dua alat dengan tujuan berbeda. Merge menjaga riwayat apa adanya, cocok untuk menggabungkan pekerjaan ke branch bersama. Rebase merapikan riwayat, cocok untuk pekerjaan lokal sebelum dibagikan. Sekali kamu paham bahwa rebase menulis ulang commit dan karenanya tidak boleh menyentuh sesuatu yang sudah dipakai orang lain, rasa takut itu hilang, dan kamu mendapatkan kembali separuh kemampuan Git yang selama ini dihindari.
