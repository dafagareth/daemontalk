## Ringkasan

**svault** adalah vault rahasia terenkripsi lokal untuk command line. Ia menyimpan
API key, password, URL database, dan string sensitif lainnya di dalam satu file
terenkripsi di mesinmu sendiri. Tanpa sinkronisasi cloud, tanpa telemetri, tanpa
akun, tanpa server yang perlu berjalan.

Seluruh tool ini dikemas dalam satu binary statis yang berjalan di Linux, macOS,
dan Windows. Buka kuncinya sekali dengan master password, kerjakan sesi 30 menit,
lalu ia mengunci dirinya sendiri secara otomatis ketika sesi berakhir.

```
$ svault set DB_PASSWORD supersecret123
OK  DB_PASSWORD saved

$ svault get DB_PASSWORD
supersecret123
```

> **Catatan nama:** svault dulunya bernama `stash`. Diubah pada v2.0.0 untuk
> menghindari konflik dengan `git stash` dan sebuah paket AUR yang tidak
> berkaitan. Mulai v2.0.0, nama perintah, direktori vault (`~/.svault`), dan
> environment variable semuanya menggunakan awalan `svault`.

## Kenapa svault?

Kebanyakan developer mengelola puluhan rahasia di banyak proyek. File `.env`
dalam format teks polos berisiko masuk ke riwayat git. Password manager lengkap
terlalu berat untuk skrip shell. Vault berbasis cloud menambah latensi jaringan
dan ketergantungan pada akun eksternal.

| Masalah | Solusi svault |
|---|---|
| File `.env` teks polos bisa ter-commit ke Git | Setiap nilai dienkripsi dengan AES-256-GCM di disk |
| Secret manager enterprise terlalu berat | Satu binary, tanpa server, daemon, atau konfigurasi rumit |
| Rahasia tersebar di banyak proyek | Satu namespace terpisah per proyek dalam satu vault |
| Tidak ada jejak siapa mengakses apa | Audit log lokal append-only mencatat setiap operasi |
| Harus mengetik master password terus-menerus | Buka kunci sekali, bekerja hingga 30 menit |

## Instalasi

### Arch Linux (AUR)

```bash
yay -S svault
# atau
paru -S svault
```

### Instalasi cepat: Linux / macOS

```bash
curl -fsSL https://raw.githubusercontent.com/dafagareth/svault/main/install.sh | sh
```

Skrip ini mendeteksi OS dan arsitektur CPU, mengunduh binary rilis yang sesuai,
memverifikasi checksumnya, dan menempatkannya di `/usr/local/bin` (atau
`~/.local/bin` jika direktori itu tidak bisa ditulis). Tentukan versi tertentu
dengan `SVAULT_VERSION=v2.0.0`, atau lokasi instalasi dengan `SVAULT_BIN_DIR`.

### Instalasi cepat: Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/dafagareth/svault/main/install.ps1 | iex
```

Memasang ke `%LOCALAPPDATA%\svault` dan menambahkan path tersebut ke `PATH` user.

### Unduh binary rilis

Pilih file yang sesuai dengan platformmu dari halaman
[Releases](https://github.com/dafagareth/svault/releases):

| Platform | File |
|---|---|
| Linux x86_64 | `svault-linux-amd64` |
| Linux ARM64 | `svault-linux-arm64` |
| macOS Intel | `svault-darwin-amd64` |
| macOS Apple Silicon | `svault-darwin-arm64` |
| Windows x86_64 | `svault-windows-amd64.exe` |

```bash
chmod +x svault-linux-amd64
sudo mv svault-linux-amd64 /usr/local/bin/svault
```

### Build dari sumber

```bash
git clone git@github.com:dafagareth/svault.git
cd svault
sudo make install               # build dan pasang ke /usr/local/bin
sudo make install PREFIX=/usr   # prefix alternatif
sudo make uninstall             # hapus binary
```

**Kebutuhan:** Go 1.25 atau lebih baru. Manifest Homebrew (macOS) dan Scoop
(Windows) sudah ditulis tapi belum dipublikasikan ke registry masing-masing.
Gunakan skrip instalasi cepat di atas untuk sementara ini.

## Mulai Cepat

```bash
# Langkah 1: inisialisasi vault sekali saja. Kamu akan diminta master password.
$ svault init
Master password: ********
Confirm password: ********
Vault initialized at ~/.svault/vault.enc

# Langkah 2: buka kunci sesi 30 menit.
$ svault unlock
Vault unlocked. Session valid for 30 minutes.

# Langkah 3: simpan rahasia.
$ svault set DB_PASSWORD supersecret123
$ svault set JWT_SECRET=myjwtsecret          # sintaks KEY=VALUE juga bisa

# Langkah 4: gunakan rahasia tersebut.
$ svault get DB_PASSWORD                     # tampilkan satu nilai
$ svault list                                # daftar semua key (nilai tersembunyi)
$ svault export > .env                       # tulis file .env standar
$ svault exec -- npm start                   # suntikkan rahasia langsung ke perintah

# Langkah 5: kunci saat selesai (atau biarkan sesi berakhir sendiri).
$ svault lock
```

## Referensi Perintah

### Inisialisasi dan autentikasi

```bash
svault init      # buat vault baru, meminta master password
svault unlock    # buka kunci vault, mulai sesi 30 menit
svault lock      # kunci vault segera
svault status    # tampilkan apakah vault terkunci atau terbuka
```

`svault status --short` mencetak representasi ringkas (ikon gembok atau sisa
menit) yang cocok digunakan di dalam prompt shell.

### Manajemen rahasia

```bash
svault set KEY VALUE      # simpan rahasia (sintaks KEY=VALUE juga bisa)
svault set KEY --stdin    # baca nilai dari stdin, menjaga keluar dari riwayat shell
svault get KEY            # ambil nilai yang tersimpan
svault edit KEY           # buka nilai di $EDITOR, ideal untuk rahasia multibaris
svault delete KEY         # hapus rahasia secara permanen
svault list               # daftar semua key di namespace aktif
svault search PATTERN     # cari key berdasarkan nama; --all mencari semua namespace
svault rename OLD NEW     # ganti nama key sambil mempertahankan nilainya
svault move KEY --to NS   # pindahkan key ke namespace lain
```

Membaca dari stdin melindungi nilai sensitif dari riwayat shell:

```bash
echo "supersecret" | svault set DB_PASSWORD --stdin
pbpaste | svault set API_KEY --stdin       # pipe dari clipboard di macOS
```

### Clipboard dan pembuat password

```bash
svault copy KEY                        # salin nilai, clipboard terhapus otomatis setelah 30d
svault generate                        # buat password acak, salin ke clipboard
svault generate --length 32 --save DB_PASSWORD
svault generate --no-symbols           # hanya alfanumerik
svault open GITHUB                     # buka GITHUB_URL di browser dan salin GITHUB_PASS
```

`svault open KEY` mencari `KEY_URL` (atau menggunakan `KEY` langsung jika berisi
URL HTTP), membukanya di browser default, lalu menyalin `KEY_PASS`,
`KEY_PASSWORD`, atau `KEY_TOKEN` ke clipboard jika salah satunya ada.

### Integrasi shell

```bash
svault exec -- npm start               # jalankan perintah dengan rahasia sebagai env var
svault exec --ns production -- ./deploy.sh
eval $(svault env)                     # muat semua rahasia ke sesi shell saat ini
eval $(svault env --ns production)
```

### Manajemen namespace

```bash
svault use NAMESPACE           # set namespace aktif
svault ns list                 # daftar semua namespace beserta jumlah key
svault ns rename OLD NEW       # ganti nama namespace
svault ns delete NAMESPACE     # hapus namespace beserta seluruh rahasianya
svault diff staging production # bandingkan dua namespace berdampingan
```

Contoh keluaran `svault diff`:

```
= DB_PASSWORD                   same
< JWT_SECRET                    only in [staging]
> STRIPE_KEY                    only in [production]
~ REDIS_URL                     value differs

[staging] vs [production]: 1 same, 3 differ
```

### Impor, ekspor, dan verifikasi

```bash
svault export > .env                   # ekspor namespace aktif sebagai baris KEY=VALUE
svault export --ns production > .env.production
svault import .env.example             # impor rahasia dari file .env yang ada
svault check                           # verifikasi vault terhadap .env.example
svault check .env.production           # verifikasi terhadap file lain
```

Contoh keluaran `svault check`:

```
OK    DB_PASSWORD
OK    JWT_SECRET
MISS  STRIPE_KEY

2/3 keys present in namespace [default], 1 missing
```

### Perawatan dan diagnostik

```bash
svault rotate                  # ganti master password dan enkripsi ulang seluruh vault
svault backup ~/safe/vault.bak
svault backup                  # backup bertimestamp otomatis disimpan di ~/.svault/
svault restore ~/safe/vault.bak
svault info                    # versi vault, path, namespace aktif, jumlah key, status sesi
svault log                     # tampilkan audit log
svault doctor                  # cek kesehatan instalasi dan vault
svault version                 # cetak nomor versi (--short hanya angkanya)
```

Setelah `svault rotate`, password lama tidak lagi berlaku. Sesi aktif disegarkan
secara otomatis sehingga tidak perlu membuka kunci ulang.

Menjalankan `svault doctor` adalah langkah pertama yang disarankan ketika ada
yang terasa tidak beres. Ia memeriksa direktori dan file vault, izin file, file
config dan audit log, status sesi, ketersediaan git untuk deteksi namespace
otomatis, dan tool clipboard:

```
[OK  ] vault directory              /home/kamu/.svault
[OK  ] vault file                   412 bytes, mode 0600
[OK  ] config file                  /home/kamu/.svault/config.json
[OK  ] audit log                    /home/kamu/.svault/vault.log
[OK  ] session                      unlocked, 24m remaining
[OK  ] git (auto-namespace)         current namespace: grocyvo
[OK  ] clipboard                    wl-copy

All checks passed.
```

Perintah `doctor`, `init`, `version`, `completion`, dan `help` berjalan tanpa
vault yang sudah ada. Semua perintah lain keluar lebih awal dengan pesan
deskriptif jika vault tidak ditemukan.

## Namespace

Setiap proyek mendapat namespace terisolasi tersendiri, sehingga dua proyek
masing-masing dapat mendefinisikan `DB_PASSWORD` tanpa konflik.

### Deteksi otomatis

Di dalam repositori git, namespace diturunkan dari nama repo secara otomatis.
Tidak perlu `svault use` manual.

```bash
~/grocyvo$ svault set DB_URL=postgres://localhost/grocyvo
~/grocyvo$ svault info
Namespace  : grocyvo (from git, 2 total)

~/portfolio$ svault get DB_URL     # membaca dari namespace 'portfolio', bukan 'grocyvo'
```

Prioritas deteksi, dari tertinggi ke terendah:

1. Flag `--ns` pada perintah saat ini
2. Environment variable `SVAULT_NS`
3. Nama repositori git saat ini
4. Namespace aktif yang di-set dengan `svault use`
5. `default`

### Kontrol manual

```bash
svault use production               # ganti namespace aktif
svault set --ns staging DB_URL=...  # override sekali pakai tanpa berpindah namespace
svault ns list                      # lihat semua namespace dan jumlah key-nya
```

## Sesi

Saat kamu membuka kunci vault, svault menurunkan kunci enkripsi dari master
password dan menulisnya ke file sesi. Ini memungkinkan perintah berikutnya
berjalan tanpa mengetik ulang password.

- Panjang sesi default: **30 menit**, dikonfigurasi via `SVAULT_SESSION_TTL`.
- Sesi berakhir ketika TTL habis atau saat kamu menjalankan `svault lock`.
- File sesi disimpan di `/tmp/.svault_session` dengan izin `0600` dan dihapus
  secara otomatis saat kedaluwarsa.

## Tata Letak Penyimpanan

```
~/.svault/
  vault.enc       # rahasiamu, dienkripsi dengan AES-256-GCM (struktur JSON di dalamnya)
  vault.log       # audit log append-only
  config.json     # pengaturan namespace aktif

/tmp/
  .svault_session # kunci sesi turunan, mode 0600, dihapus otomatis setelah TTL
```

Format mentah `vault.enc` di disk:

```
[ 16 byte salt ][ 12 byte nonce ][ ciphertext ]
```

## Konfigurasi

| Environment variable | Default | Tujuan |
|---|---|---|
| `SVAULT_SESSION_TTL` | `30` | Durasi sesi dalam menit |
| `SVAULT_NS` | (kosong) | Paksa namespace tertentu, menimpa deteksi git |

## Autolengkap Shell

svault menghasilkan skrip autolengkap untuk Bash, Zsh, dan Fish:

```bash
# Bash
svault completion bash | sudo tee /etc/bash_completion.d/svault > /dev/null

# Zsh
svault completion zsh > "${fpath[1]}/_svault"

# Fish
svault completion fish > ~/.config/fish/completions/svault.fish
```

Restart shell atau source file yang dihasilkan secara langsung untuk
mengaktifkan autolengkap.

## Model Keamanan

svault dibangun di atas serangkaian pilihan kriptografi yang disengaja:

- **Enkripsi: AES-256-GCM.** Cipher terautentikasi. Setiap manipulasi pada file
  vault terdeteksi saat dekripsi sebelum data apa pun terekspos.
- **Penurunan kunci: Argon2id.** Algoritma yang direkomendasikan saat ini untuk
  penurunan kunci berbasis password. Setiap vault menerima salt 16-byte yang
  dibuat secara acak.
- **Keacakan: `crypto/rand`.** Semua salt, nonce, dan password yang dihasilkan
  menggunakan RNG yang aman secara kriptografis. `math/rand` tidak pernah
  digunakan untuk hal yang berkaitan dengan keamanan.

Salt disimpan dalam bentuk teks polos di dalam `vault.enc`. Ini adalah praktik
standar: salt bukan rahasia. Perannya adalah membuat serangan prakomputasi
(rainbow table) terhadap master password menjadi tidak praktis.

### Jaminan keamanan

- **Master password tidak dapat dipulihkan.** Kehilangannya berarti kehilangan
  akses ke vault secara permanen, sesuai desain. Simpan salinan password di
  lokasi yang aman.
- **Nilai** rahasia tidak pernah ditulis ke log, file temp, atau output standar
  kecuali melalui `get`, `copy`, `env`, dan `export`.
- Setiap operasi tulis lebih dulu membuat `vault.enc.bak` dan beberapa salinan
  rollback bertimestamp. Penulisan yang gagal tidak dapat menghancurkan status
  bersih sebelumnya.
- Proses svault yang bersamaan diserialisasi menggunakan **file lock eksklusif**.
  Dua penulisan simultan tidak dapat merusak satu sama lain.
- File sesi dibuat dengan izin `0600`, hanya bisa dibaca oleh user pemilik.
- Key rahasia harus sesuai dengan aturan penamaan variabel shell yang valid,
  memastikan keluaran `export` dan `env` selalu aman untuk di-`eval`.

### Trade-off yang diketahui: kunci sesi ada di /tmp

Selama vault terbuka, kunci enkripsi turunan disimpan sebagai teks polos di
`/tmp/.svault_session` (izin `0600`). Ini adalah trade-off kenyamanan yang
disengaja untuk menghindari pengetikan master password berulang-ulang.

- User biasa lain di mesin yang sama **tidak dapat** membaca file ini.
- **Root, atau proses dengan hak root, dapat membacanya.**
- File dihapus secara otomatis setelah TTL atau segera saat `svault lock`.

Pada server bersama atau sistem di mana pengguna lain memiliki akses root,
jalankan `svault lock` segera setelah selesai bekerja, jangan menunggu TTL.

## Skenario Umum

### Mengganti file .env proyek

```bash
cd ~/myproject             # namespace otomatis menjadi "myproject"
svault import .env         # impor rahasia yang ada dari file tersebut
rm .env                    # hapus file teks polos
svault exec -- npm run dev # jalankan dengan rahasia tersuntik, tanpa file di disk
```

### Onboarding ke proyek dengan .env.example

```bash
svault check               # key mana dari .env.example yang belum di-set?
# OK    DATABASE_URL
# MISS  STRIPE_KEY         <- set ini berikutnya
svault set STRIPE_KEY=sk_test_...
```

### Menjauhkan rahasia dari riwayat shell

```bash
echo "$GENERATED" | svault set API_KEY --stdin   # nilai tidak pernah muncul sebagai argumen
svault edit MULTILINE_CERT                        # buka $EDITOR untuk nilai yang panjang
```

### Mempromosikan rahasia dari staging ke production

```bash
svault diff staging production       # tinjau perbedaan sebelum mengubah apa pun
svault move STRIPE_KEY --to production
```

### CI dan deployment otomatis

```bash
export SVAULT_NS=production
svault exec -- ./deploy.sh           # rahasia hanya tercakup untuk proses itu saja
```

### Mengganti master password

```bash
svault rotate   # meminta password lama dan baru, mengenkripsi ulang seluruh vault
                # password lama langsung tidak berlaku; sesi disegarkan otomatis
```

### Migrasi dari stash versi lama (sebelum v2.0)

```bash
mv ~/.stash ~/.svault    # pindahkan direktori vault ke lokasi baru
svault unlock            # buka kunci seperti biasa dengan master password yang sama
```

## FAQ

```faq
Q: Saya lupa master password. Apakah rahasia saya bisa dipulihkan?
A: Tidak. Tidak ada jalur pemulihan, backdoor, atau mekanisme reset. Ini disengaja: opsi pemulihan akan melemahkan jaminan keamanan bahwa hanya seseorang yang tahu password yang dapat membaca vault. Satu-satunya opsi adalah memulihkan dari backup yang dibuat saat kamu masih tahu password, atau menginisialisasi ulang vault dan memasukkan kembali rahasiamu. Simpan selalu master password di lokasi aman yang terpisah.

Q: Apakah svault mengirim data ke internet?
A: Tidak pernah. Tidak ada kode jaringan di jalur eksekusi utama. Semua data tetap ada di `~/.svault` di mesin lokalmu. Satu-satunya aktivitas jaringan yang mungkin terjadi adalah saat skrip instalasi mengunduh binary dari GitHub Releases selama setup awal, atau saat kamu membuka URL tersimpan di vault dengan `svault open`.

Q: Bisakah dua proyek berbeda menggunakan nama rahasia yang sama?
A: Ya. Namespace memisahkannya sepenuhnya. `myproject/DB_PASSWORD` dan `other-project/DB_PASSWORD` adalah nilai yang independen. Saat kamu menjalankan svault di dalam repositori git, namespace diset secara otomatis dari nama repo, sehingga tidak ada konfigurasi khusus yang diperlukan.

Q: Apakah aman menggunakan svault di server bersama?
A: Dengan kehati-hatian. File vault dan file sesi keduanya dibuat dengan izin `0600`, yang mencegah user biasa lain membacanya. Namun, siapa pun dengan akses root di mesin dapat membaca kunci sesi dari `/tmp/.svault_session` selama vault terbuka. Di server bersama atau multi-user, selalu jalankan `svault lock` segera setelah selesai bekerja, jangan mengandalkan TTL sesi.

Q: Apakah svault membutuhkan daemon atau layanan latar belakang?
A: Tidak. svault adalah tool command-line biasa. Ia memulai, menjalankan operasinya, dan keluar. Satu-satunya state yang menetap adalah file vault terenkripsi, audit log, dan file sesi berumur pendek di `/tmp`. Tidak ada yang berjalan di latar belakang di antara perintah-perintah.

Q: Apa yang terjadi jika file sesi di /tmp terhapus?
A: svault memperlakukan file sesi yang hilang sama seperti sesi yang kedaluwarsa. Perintah berikutnya yang memerlukan vault terbuka akan menolak berjalan dan memintamu menjalankan `svault unlock` lagi. Vault dan semua rahasia yang tersimpan tetap utuh; hanya state sesi yang hilang.

Q: Bisakah svault digunakan di pipeline CI/CD atau container Docker?
A: Ya, dengan sedikit perencanaan. Pendekatan yang disarankan adalah membuka kunci vault di awal skrip pipeline dan menggunakan `svault exec` untuk menyuntikkan rahasia ke setiap perintah. Set `SVAULT_NS` ke namespace yang benar dan gunakan `SVAULT_SESSION_TTL` untuk memperpanjang sesi jika pipelinemu membutuhkan lebih dari 30 menit. File sesi akan tercakup dalam container atau runner, sehingga tidak akan bertahan antar run.

Q: Bagaimana cara berbagi rahasia dengan rekan tim?
A: svault adalah tool lokal dan tidak memiliki mekanisme berbagi bawaan. Untuk berbagi rahasia, ekspor ke file `.env` dengan `svault export` dan transfer file tersebut melalui saluran aman (pesan terenkripsi, secrets manager, atau berbagi file aman). Penerima kemudian dapat mengimpornya dengan `svault import`. Jangan pernah mengirim file `.env` melalui email atau chat yang tidak terenkripsi.

Q: Apa yang terjadi jika svault init dijalankan saat vault sudah ada?
A: svault mendeteksi file vault yang ada dan menolak menimpanya tanpa konfirmasi. Rahasia yang ada tidak terpengaruh kecuali kamu mengonfirmasi reinisialisasi secara eksplisit. Jika ingin mulai dari awal, hapus `~/.svault/vault.enc` secara manual terlebih dahulu, lalu jalankan `svault init`.

Q: Bisakah saya memiliki beberapa vault untuk tujuan berbeda?
A: Tidak secara langsung. svault menggunakan satu file `vault.enc` dan memisahkan rahasia berdasarkan namespace di dalamnya. Untuk sebagian besar kasus, namespace sudah cukup: setiap proyek mendapat namespace terisolasinya sendiri, dan kamu bisa membuat sebanyak yang diperlukan. Jika kamu benar-benar membutuhkan vault terpisah dengan master password berbeda, kamu perlu mengelola file vault secara manual.

Q: Tool clipboard apa yang didukung svault?
A: svault mendeteksi tool clipboard yang tersedia berdasarkan lingkunganmu: **Linux (Wayland):** `wl-copy` (dari paket `wl-clipboard`) **Linux (X11):** `xsel` atau `xclip` **macOS:** `pbcopy` (sudah termasuk dalam sistem) Jalankan `svault doctor` untuk mengonfirmasi tool mana yang terdeteksi. Jika tidak ada yang ditemukan, perintah clipboard (`copy`, `generate`, `open`) akan gagal dengan pesan error yang deskriptif.

Q: Apa yang terjadi jika dua proses svault berjalan bersamaan?
A: svault menggunakan file lock eksklusif pada setiap operasi tulis. Jika dua proses mencoba menulis secara bersamaan, yang kedua menunggu yang pertama selesai baru kemudian melanjutkan. Pembacaan tidak dikunci dan bisa berjalan secara bersamaan. Ini mencegah race condition di mana dua penulisan simultan merusak satu sama lain dan diam-diam menghilangkan data.

Q: Bagaimana cara membuat backup vault?
A: Gunakan perintah backup bawaan: `svault backup ~/safe/vault.bak` menyalin vault ke path tertentu. `svault backup` (tanpa path) membuat salinan bertimestamp di `~/.svault/` secara otomatis. svault juga membuat salinan rollback otomatis pada setiap penulisan, sehingga penulisan yang buruk tidak pernah menghancurkan status baik sebelumnya. Untuk pemulihan bencana jangka panjang, salin `~/.svault/vault.enc` ke penyimpanan eksternal atau lokasi backup terenkripsi yang terpisah.

Q: Bisakah svault digunakan tanpa git terinstal?
A: Ya. Git hanya digunakan untuk deteksi namespace otomatis. Jika git tidak terinstal atau kamu tidak berada di dalam repositori git, svault kembali ke namespace yang di-set dengan `svault use`, environment variable `SVAULT_NS`, atau namespace `default`. Semua fitur lain berjalan normal. `svault doctor` akan melaporkan peringatan tentang git yang tidak ada, tapi tidak ada perintah yang gagal karena itu.

Q: Bisakah svault diotomatisasi di shell script tanpa prompt password interaktif?
A: Alur kerja yang dimaksudkan adalah menjalankan `svault unlock` sekali secara interaktif di awal sesi, lalu menggunakan semua perintah lain secara non-interaktif dalam jendela 30 menit tersebut. Di dalam pipeline CI di mana input interaktif tidak memungkinkan, gunakan `SVAULT_SESSION_TTL` untuk menjaga sesi tetap aktif selama durasi pipeline. Untuk lingkungan yang sepenuhnya otomatis, pertimbangkan apakah secrets manager khusus dengan autentikasi identitas mesin lebih sesuai.

Q: Apa pengaturan Argon2id yang digunakan svault?
A: svault menggunakan Argon2id yang disetel untuk sekitar 200ms waktu penurunan pada laptop modern. Ini memberikan resistansi yang berarti terhadap serangan brute-force dan dictionary, sambil menjaga perintah unlock cukup cepat untuk penggunaan interaktif sehari-hari. Parameter yang tepat (memori, iterasi, paralelisme) disimpan bersama salt sehingga versi mendatang dapat mengubah default tanpa merusak vault yang ada.
```

## Tumpukan Teknologi

- **Go 1.25+**, dikompilasi ke satu binary statis tanpa CGO dan tanpa dependensi runtime
- **`golang.org/x/crypto`** untuk penurunan kunci Argon2id dan enkripsi AES-256-GCM
- **Cobra** untuk pohon perintah dan pembuatan autolengkap shell
- Di-cross-compile untuk Linux, macOS, dan Windows (amd64 dan arm64) di CI via GitHub Actions

## Status

Dirilis pada v2.0.0. Tersedia dan terkemas di AUR. Manifest Homebrew dan Scoop
sudah ditulis tapi belum disubmit ke registry masing-masing. Kode sumber, binary
rilis, kebijakan keamanan lengkap, dan changelog tersedia di
[repositori GitHub](https://github.com/dafagareth/svault).
