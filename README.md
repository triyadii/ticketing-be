# Backend IT Helpdesk Ticketing & Manajemen Data Koperasi

Aplikasi backend RESTful API berkinerja tinggi yang dibangun menggunakan **Go (Golang)** dengan web framework **Fiber v2**, ORM **GORM**, dan basis data **PostgreSQL**. Dilengkapi manajemen tiket kendala IT (upload foto max 3 berkas), manajemen staf petugas IT, pelapor, kategori layanan, manajemen & export data lama koperasi (Excel `.xlsx` / `.csv`), serta dokumentasi interaktif dengan **Swagger**.

---

## 🚀 Fitur Utama

1. **Autentikasi & Otorisasi Tim IT (JWT)**
   - Login petugas IT dengan token JWT dan enkripsi password Bcrypt.
   - Proteksi endpoint administratif menggunakan Middleware Bearer Token.

2. **Manajemen Staf Petugas IT (`users`)**
   - Tambah petugas, ubah profil/password, atur status (`aktif` / `nonaktif`), dan hapus akun.

3. **Manajemen Kategori Layanan IT (`layanan`)**
   - CRUD jenis gangguan/layanan IT (contoh: Jaringan & Wi-Fi, Hardware/PC, Aplikasi Internal, dll).

4. **Master Data Pelapor Terpisah (`pelapor`)**
   - Master data pelapor tanpa registrasi login akun (`nik` sebagai Primary Key).
   - Riwayat seluruh tiket yang pernah dilaporkan oleh NIK terkait otomatis terhubung.

5. **Manajemen Tiket Kendala (`tickets`)**
   - Pembuatan tiket publik melalui form `multipart/form-data`.
   - **Upload Lampiran Foto:** Mendukung upload gambar (JPG, PNG, WEBP) **maksimal 3 berkas** (limit ukuran 5MB/file) yang otomatis disimpan di storage lokal dan disajikan secara statis.
   - **Auto-Generate Nomor Tiket:** Penomoran otomatis format `IT-YYYYMM-XXXX` (contoh: `IT-202611-0001`).
   - **Workflow Penanganan Tiket:**
     - Status: `menunggu` ➔ `diproses` (diambil petugas IT) ➔ `selesai` (wajib isi `solusi`, waktu `finished_at` tercatat otomatis) / `dibatalkan`.
   - **Tracking Tiket:** Publik dapat melacak status tiket langsung menggunakan nomor tiket atau UUID.

6. **Tarik Data & Export Data Lama Koperasi (`legacy_koperasi`)**
   - Menampilkan data lama dengan pagination, pencarian multi-kolom, dan filter tanggal.
   - Export langsung ke format **Excel (`.xlsx`)** dengan styling header profesional (`excelize/v2`).
   - Export ke format **CSV (`.csv`)** untuk kompatibilitas data luas.

7. **Dokumentasi API Terintegrasi (Swagger UI)**
   - OpenAPI 2.0 / 3.0 via `gofiber/swagger` pada route `/swagger/index.html`.

---

## 🛠️ Tech Stack

* **Bahasa**: Go 1.22+
* **Web Framework**: [GoFiber v2](https://github.com/gofiber/fiber)
* **Database**: PostgreSQL (Database: `ticketing_db`)
* **ORM**: [GORM](https://gorm.io/) + `gorm.io/driver/postgres`
* **Dokumentasi API**: [Swaggo](https://github.com/swaggo/swag) + `gofiber/swagger`
* **Spreadsheet Generator**: [Excelize v2](https://github.com/xuri/excelize)
* **Keamanan**: JWT (`golang-jwt/jwt/v5`) & Bcrypt (`golang.org/x/crypto/bcrypt`)
* **Containerization**: Docker & Docker Compose

---

## 🗄️ Skema Database

```
1. users (Khusus Staf Tim IT)
   ├── id (UUID, PK)
   ├── nama (VARCHAR 100)
   ├── username (VARCHAR 50, Unique)
   ├── password (VARCHAR 255 - Bcrypt)
   ├── status (ENUM: 'aktif', 'nonaktif')
   └── created_at, updated_at

2. layanan (Kategori Kendala IT)
   ├── id (UUID, PK)
   ├── nama_layanan (VARCHAR 100)
   ├── status (ENUM: 'aktif', 'nonaktif')
   └── created_at, updated_at

3. pelapor (Master Data Pegawai / Pemohon)
   ├── nik (VARCHAR 16, PK)
   ├── nama (VARCHAR 100)
   ├── instansi (VARCHAR 150)
   ├── no_hp (VARCHAR 20)
   └── created_at, updated_at

4. tickets (Transaksi Kendala IT)
   ├── id (UUID, PK)
   ├── nomor_ticket (VARCHAR 50, Unique)
   ├── nik_pelapor (FK -> pelapor.nik)
   ├── id_pelayanan (FK -> layanan.id)
   ├── id_petugas (FK -> users.id, Nullable)
   ├── judul_kendala (VARCHAR 150)
   ├── keterangan (TEXT)
   ├── foto (JSON Array - Max 3 URL path)
   ├── solusi (TEXT, Nullable)
   ├── status ('menunggu', 'diproses', 'selesai', 'dibatalkan')
   ├── finished_at (TIMESTAMP, Nullable)
   └── created_at, updated_at

5. legacy_koperasi (Data Lama Koperasi)
   ├── id (BIGSERIAL, PK)
   ├── kode_anggota, nama_anggota, unit_kerja
   ├── jenis_transaksi, nominal, keterangan, tanggal
   └── created_at
```

---

## ⚙️ Menjalankan Aplikasi

### Opsi A: Menggunakan Docker Compose (Database di Host)

1. Pastikan berkas `.env` sudah sesuai (password db, port, dll).
2. Jalankan container:
   ```bash
   docker compose up --build -d
   ```
3. Memeriksa logs aplikasi:
   ```bash
   docker compose logs -f app
   ```

> **Catatan Koneksi ke PostgreSQL Host (Linux):**
> * File `docker-compose.yml` telah dikonfigurasi menggunakan `extra_hosts: - "host.docker.internal:host-gateway"`.
> * Jika tidak ingin mengubah konfigurasi listen IP PostgreSQL di host, Anda juga bisa mengaktifkan `network_mode: "host"` pada `docker-compose.yml` dan tetap menggunakan `DB_HOST=127.0.0.1`.

---

### Opsi B: Menjalankan Langsung di Host (Go / Air)

```bash
# 1. Menggunakan Live-Reload Air
air

# 2. Atau build binary langsung
go build -o server cmd/api/main.go
./server
```

> **Akun Admin IT Default:**
> * **Username:** `admin`
> * **Password:** `password123`

---

## 📖 Dokumentasi Swagger UI

Buka tautan berikut di browser:
👉 **[http://localhost:8086/swagger/index.html](http://localhost:8086/swagger/index.html)**

Jika Anda melakukan perubahan anotasi Swagger di kode, perbarui dokumentasi dengan perintah:
```bash
swag init -g cmd/api/main.go
```

---

## 📡 Daftar Endpoint API

### 1. Autentikasi (`/api/v1/auth`)
* `POST /api/v1/auth/login` - Login petugas IT mendapatkan JWT Token
* `GET /api/v1/auth/profile` - Profil petugas yang sedang login *(Protected)*

### 2. Petugas IT (`/api/v1/users` - Protected)
* `GET /api/v1/users` - Daftar petugas IT (pagination, filter `status`, search)
* `POST /api/v1/users` - Tambah petugas IT baru
* `GET /api/v1/users/:id` - Detail petugas IT
* `PUT /api/v1/users/:id` - Update petugas IT (nama, status, password)
* `DELETE /api/v1/users/:id` - Hapus akun petugas IT

### 3. Layanan IT (`/api/v1/layanan`)
* `GET /api/v1/layanan` - Daftar layanan IT aktif
* `GET /api/v1/layanan/:id` - Detail layanan
* `POST /api/v1/layanan` - Tambah layanan baru *(Protected)*
* `PUT /api/v1/layanan/:id` - Update layanan *(Protected)*
* `DELETE /api/v1/layanan/:id` - Hapus layanan *(Protected)*

### 4. Pelapor (`/api/v1/pelapor`)
* `POST /api/v1/pelapor` - Simpan data pelapor
* `GET /api/v1/pelapor/:nik` - Detail pelapor dan riwayat seluruh tiketnya
* `GET /api/v1/pelapor` - Daftar pelapor *(Protected)*
* `PUT /api/v1/pelapor/:nik` - Update profil pelapor *(Protected)*
* `DELETE /api/v1/pelapor/:nik` - Hapus pelapor *(Protected)*

### 5. Tiket Kendala (`/api/v1/tickets`)
* `POST /api/v1/tickets` - Buat tiket baru (Multipart form-data: upload max 3 foto, auto generate nomor tiket)
* `GET /api/v1/tickets` - Daftar tiket (filter status, layanan, petugas, pelapor, rentang tanggal)
* `GET /api/v1/tickets/:identifier` - Detail tiket berdasarkan UUID atau Nomor Tiket (misal `IT-202611-0001`)
* `PATCH /api/v1/tickets/:id/assign` - Ambil tiket oleh petugas IT *(Protected)*
* `PATCH /api/v1/tickets/:id/status` - Update status tiket (wajib mengisi `solusi` saat status `selesai`) *(Protected)*
* `DELETE /api/v1/tickets/:id` - Hapus tiket *(Protected)*

### 6. Data Lama Koperasi (`/api/v1/koperasi` - Protected)
* `GET /api/v1/koperasi/data-lama` - Menampilkan data lama koperasi (pagination & filter)
* `GET /api/v1/koperasi/data-lama/export?format=excel` - Download file Excel (`.xlsx`)
* `GET /api/v1/koperasi/data-lama/export?format=csv` - Download file CSV (`.csv`)
