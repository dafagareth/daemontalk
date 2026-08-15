---
title: "Penghapusan Global Interpreter Lock (GIL) di Python 3.13+"
slug: e1c9d4b7
aliases: python-gil-penghapusan
date: 2026-06-18
tags: [python, performance, backend]
lang: id
draft: false
---

Selama tiga dekade, Global Interpreter Lock (GIL) menjadi batasan fundamental pada program Python multithread. Python 3.13 memperkenalkan build eksperimental tanpa GIL sebagai hasil dari PEP 703, membuka kemungkinan paralelisme sejati pada thread CPython untuk pertama kalinya.

## Fakta Menarik

**Fakta 1.** GIL diperkenalkan pada awal 1990-an sebagai solusi sederhana untuk keamanan thread pada manajemen memori CPython. Menghapusnya terbukti sangat sulit karena hampir setiap bagian interpreter bergantung padanya.

**Fakta 2.** PEP 703 ditulis oleh Sam Gross dan diterima pada tahun 2023. Implementasinya menggunakan teknik "biased reference counting" untuk mengurangi overhead sinkronisasi pada objek yang hanya diakses satu thread.

**Fakta 3.** Build free-threaded Python 3.13 tersedia secara paralel dengan build standar. Keduanya dapat diinstal bersamaan, dan ekstensi C harus dikompilasi ulang secara eksplisit untuk mendukung mode ini.

---

## Tips dan Trik

### 1. Apa Itu GIL dan Mengapa Dibuat

GIL adalah mutex yang melindungi akses ke objek Python sehingga hanya satu thread yang dapat mengeksekusi bytecode Python pada satu waktu. Ini menyederhanakan implementasi CPython karena struktur data internal tidak perlu thread-safe secara individual.

Akibatnya, program Python yang intensif CPU tidak mendapatkan manfaat dari multiple core jika menggunakan thread Python biasa. Solusi selama ini adalah `multiprocessing` (proses terpisah dengan overhead fork) atau ekstensi C yang melepas GIL secara manual seperti NumPy.

```python
import threading
import time

def hitung_intensif(n):
    total = 0
    for i in range(n):
        total += i * i
    return total

# Dengan GIL: dua thread tidak berjalan secara paralel untuk kode Python murni
mulai = time.perf_counter()
t1 = threading.Thread(target=hitung_intensif, args=(10_000_000,))
t2 = threading.Thread(target=hitung_intensif, args=(10_000_000,))
t1.start(); t2.start()
t1.join(); t2.join()
print(f"Waktu dengan thread: {time.perf_counter() - mulai:.2f}s")
```

### 2. Mengaktifkan Mode Free-Threaded di Python 3.13

Build free-threaded tersedia sebagai varian terpisah. Pada sistem dengan `pyenv` atau distribusi yang menyediakan paket `python3.13t`.

```bash
# Instal Python 3.13 free-threaded dengan pyenv
PYTHON_CONFIGURE_OPTS="--disable-gil" pyenv install 3.13.0

# Atau pada Arch Linux (jika tersedia di repositori)
sudo pacman -S python-freethreaded

# Verifikasi bahwa GIL dinonaktifkan
python3.13t -c "import sys; print(sys._is_gil_enabled())"
# Output: False

# Atau periksa pada runtime
python3.13t -X gil=0 -c "import sys; print(sys._is_gil_enabled())"
```

```python
# Deteksi mode free-threaded dari dalam program
import sys

if hasattr(sys, '_is_gil_enabled'):
    if not sys._is_gil_enabled():
        print("Berjalan dalam mode free-threaded")
    else:
        print("GIL aktif")
else:
    print("Python < 3.13, GIL selalu aktif")
```

### 3. Mengukur Peningkatan Performa pada Kode CPU-Bound

```python
import threading
import time
import sys

def hitung_prima(batas):
    """Fungsi intensif CPU murni Python."""
    prima = []
    for n in range(2, batas):
        if all(n % i != 0 for i in range(2, int(n**0.5) + 1)):
            prima.append(n)
    return prima

def benchmark(n_thread, batas=50_000):
    threads = [threading.Thread(target=hitung_prima, args=(batas,))
               for _ in range(n_thread)]
    mulai = time.perf_counter()
    for t in threads:
        t.start()
    for t in threads:
        t.join()
    return time.perf_counter() - mulai

# Jalankan dengan python3.13 (GIL aktif) dan python3.13t (GIL nonaktif)
print(f"GIL aktif: {sys._is_gil_enabled()}")
print(f"1 thread : {benchmark(1):.2f}s")
print(f"4 thread : {benchmark(4):.2f}s")
```

### 4. Implikasi untuk Ekstensi C dan Library Pihak Ketiga

Tidak semua library langsung kompatibel dengan mode free-threaded. Ekstensi C yang mengandalkan GIL sebagai mekanisme keamanan internal perlu diperbarui.

```python
# Periksa apakah modul ekstensi mendukung free-threaded
import importlib.util

def cek_kompatibilitas(nama_modul):
    spec = importlib.util.find_spec(nama_modul)
    if spec and spec.origin:
        # Ekstensi free-threaded biasanya memiliki suffix .cpython-313t
        return "313t" in (spec.origin or "")
    return False

for modul in ["numpy", "pandas", "cryptography"]:
    try:
        status = cek_kompatibilitas(modul)
        print(f"{modul}: {'kompatibel' if status else 'perlu verifikasi'}")
    except Exception as e:
        print(f"{modul}: tidak ditemukan ({e})")
```

### 5. Kapan Harus Beralih ke Free-Threaded

Mode free-threaded paling berguna untuk:
- Server web yang menangani banyak permintaan secara bersamaan dengan logika Python murni.
- Pipeline pemrosesan data yang sebelumnya memerlukan `multiprocessing`.
- Aplikasi yang menggabungkan I/O dan komputasi pada thread yang sama.

Program yang sudah menggunakan `asyncio` atau `multiprocessing` tidak serta-merta mendapat manfaat langsung. Untuk beban kerja I/O-bound, `asyncio` masih merupakan pilihan yang lebih efisien karena overhead sinkronisasi free-threaded tidak diperlukan.

```bash
# Jalankan program eksisting dengan mode free-threaded untuk pengujian
python3.13t -X gil=0 program_saya.py

# Aktifkan GIL kembali sementara untuk kompatibilitas mundur
python3.13t -X gil=1 program_saya.py
```
