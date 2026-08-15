---
title: "Penyelarasan Kecerdasan Buatan melalui Direct Preference Optimization (DPO)"
slug: 8a7b6c5d
aliases: []
date: 2026-08-11
tags: [ai]
lang: id
draft: false
type: post
cover: ""
---

## Pengantar Penyelarasan Model

Pelatihan awal Model Bahasa Besar (LLM) pada korpus teks berskala raksasa menghasilkan kapabilitas untuk memprediksi token selanjutnya secara sangat akurat. Namun, model pra-latih (pre-trained model) ini pada dasarnya hanya merupakan mesin probabilistik yang tidak memiliki pemahaman intrinsik mengenai etika, kebenaran faktual, atau instruksi manusia. Agar kecerdasan buatan ini aman dan bermanfaat, diperlukan tahap esensial yang disebut penyelarasan (alignment). Penyelarasan bertujuan untuk mengubah model probabilistik bebas menjadi asisten yang mematuhi preferensi, nilai, dan batasan moral kemanusiaan.

Secara historis, metodologi standar emas untuk tugas penyelarasan adalah Reinforcement Learning from Human Feedback (RLHF). Meskipun sangat sukses dalam memfasilitasi peluncuran berbagai asisten virtual komersial, RLHF terkenal sangat rumit secara arsitektural dan rentan terhadap ketidakstabilan numerik. Baru-baru ini, Direct Preference Optimization (DPO) diusulkan sebagai alternatif teoretis dan praktis yang mengatasi banyak kekurangan RLHF.

## Dekonstruksi Kompleksitas RLHF

Untuk mengapresiasi inovasi yang dibawa oleh DPO, kita harus membedah alur kerja RLHF terlebih dahulu. Proses RLHF terdiri dari beberapa fase terpisah. Pertama, anotator manusia mengurutkan berbagai luaran model berdasarkan kualitasnya. Kumpulan data preferensi ini digunakan untuk melatih jaringan saraf terpisah yang disebut model hadiah (reward model). Model hadiah belajar memberikan skor pada probabilitas bahwa sebuah teks akan disukai oleh manusia.

Fase akhir adalah optimasi model bahasa itu sendiri menggunakan algoritma Pembelajaran Penguatan, yang paling umum adalah Proximal Policy Optimization (PPO). Model kebijakan menghasilkan teks, model hadiah mengevaluasi kualitasnya, dan umpan balik tersebut memodulasi parameter kebijakan. Siklus umpan balik tertutup yang melibatkan dua model deep learning besar secara simultan ini membutuhkan kalibrasi hiperparameter yang ekstrem, sensitif terhadap fluktuasi pelatihan, dan menghabiskan sumber daya memori dalam jumlah masif.

## Konseptualisasi Matematika DPO

Direct Preference Optimization merestrukturisasi masalah pemodelan hadiah ini secara fundamental. Para peneliti menemukan equivalensi analitik eksplisit antara bentuk optimal dari model hadiah berbasis preferensi Bradley-Terry dengan fungsi kebijakan dari LLM itu sendiri. Melalui derivasi matematika, fungsi tujuan (objective function) berbasis hadiah dalam RLHF dapat dibalik posisinya.

Alih-alih melatih model perantara untuk mengestimasi skor, DPO secara langsung memanfaatkan jaringan saraf yang sedang diselaraskan sebagai representasi dari preferensi manusia. Objektif pelatihan DPO adalah sebuah fungsi kerugian klasifikasi silang-entropi yang sederhana. Fungsi ini meningkatkan probabilitas logaritmik model dalam menghasilkan respons yang disukai (chosen), sementara secara simultan menurunkan probabilitas logaritmik terhadap respons yang ditolak (rejected), semua ini dengan mempertimbangkan penyimpangan dinamis (KL divergence) dari model referensi orisinal.

## Signifikansi Algoritmik dan Komputasi

Penyederhanaan matematis dalam kerangka DPO berdampak sangat transformatif terhadap siklus rekayasa AI. 

### Stabilitas dan Skalabilitas

Pertama, pelatihan model menjadi proses klasifikasi standar alih-alih proses Reinforcement Learning yang berosilasi. Algoritma optimasi seperti AdamW dapat diterapkan secara langsung dengan perilaku konvergensi yang dapat diprediksi secara empiris. Penurunan gradien deterministik ini mencegah keruntuhan mode (mode collapse) yang sering terjadi pada pelatihan dengan umpan balik RL yang bising.

### Efisiensi Sumber Daya Memori

Kedua, jejak operasional secara dramatis berkurang. Dalam skema RLHF tradisional, memori GPU harus menyimpan parameter dan turunan gradien untuk model referensi, model kebijakan aktif, model hadiah, dan seringkali fungsi penilaian nilai (value function). Dengan DPO, persyaratan beban memori direduksi secara radikal. Pelatihan hanya mengelola status model referensi yang dikunci (frozen) dan model yang sedang dioptimalkan. Berkurangnya batas sumber daya ini mendemokratisasi proses fine-tuning preferensi, memungkinkan institusi kecil atau perorangan melakukan penyelarasan LLM pada infrastruktur yang terjangkau.

## Arah Penelitian Kontemporer

Validasi komprehensif membuktikan performa kompetitif DPO berbanding metode konvensional dalam menekan halusinasi faktual dan menyaring ujaran agresif. Meski begitu, wacana epistemik terkini masih menyelidiki batasan DPO. Fokus analisis beralih pada kapasitas generalisasi preferensi di luar domain teks pelatihan (out-of-distribution generalization). Masalah penolakan berlebih (over-refusal) di mana model menjadi terlampau konservatif juga masih menjadi tantangan empiris. Selain itu, varian metodologi optimasi preferensi lain yang menangani anomali dalam penilaian anotator terus disintesis.

## Kesimpulan

Penemuan Direct Preference Optimization menandai pergeseran substansial dalam rekayasa penyelarasan AI. Dengan memintas arsitektur model hadiah dan formulasi RL yang berbelit-belit, DPO mereduksi penyelarasan manusia menjadi fungsi kerugian optimasi sederhana yang elegan secara matematis dan tangguh secara komputasi. Inovasi metodologis ini menyederhanakan saluran distribusi, mempercepat siklus penelitian, dan secara fundamental memastikan bahwa pengerahan teknologi kognitif generatif berskala besar dapat berjalan lebih aman dan akomodatif terhadap batasan etika kolektif masyarakat.[^1][^2]

## Referensi

[^1]: Rafailov, R., et al. "Direct Preference Optimization: Your Language Model is Secretly a Reward Model." NeurIPS, 2023.
[^2]: Ouyang, L., et al. "Training language models to follow instructions with human feedback." (RLHF Baseline). NeurIPS, 2022.