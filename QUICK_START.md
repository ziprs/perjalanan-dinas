# Quick Start Guide - Sistem Perjalanan Dinas

Panduan cepat untuk menjalankan aplikasi di local development.

## Prerequisites

Pastikan sudah terinstall:
- **PostgreSQL** 14+
- **Go** 1.21+
- **Node.js** 18+
- **npm** atau **yarn**

## Setup dalam 5 Menit

### 1. Clone/Download Project

```bash
cd /path/to/project/perjalanan-dinas
```

### 2. Setup Database

```bash
# Buat database
createdb perjalanan_dinas

# Atau dengan psql:
psql -U postgres
CREATE DATABASE perjalanan_dinas;
\q
```

### 3. Setup Backend

```bash
cd backend

# Copy environment file
cp .env.example .env

# Edit .env jika perlu (sesuaikan dengan konfigurasi PostgreSQL Anda)
# nano .env

# Install dependencies
go mod download

# Run migrations (opsional - akan auto migrate saat aplikasi start)
# psql -U postgres -d perjalanan_dinas -f migrations/001_initial_schema.sql
# psql -U postgres -d perjalanan_dinas -f migrations/002_seed_data.sql

# Jalankan backend
go run cmd/api/main.go
```

Backend akan berjalan di **http://localhost:8080**

### 4. Setup Frontend (Terminal Baru)

```bash
cd frontend

# Install dependencies
npm install

# Copy environment file
cp .env.local.example .env.local

# Jalankan frontend
npm run dev
```

Frontend akan berjalan di **http://localhost:3000**

## Testing Aplikasi

### 1. Buka Browser

Akses: **http://localhost:3000**

### 2. Login sebagai Admin

- Klik "Admin Login" di halaman utama
- Login dengan:
  - **Username:** `admin`
  - **Password:** `admin123`

### 3. Setup Data Master

#### a. Tambah Pengkodean Jabatan
1. Menu: **Pengkodean Jabatan**
2. Tambahkan beberapa kode jabatan:
   - Posisi: `Manager`, Kode: `MNG`
   - Posisi: `Staff`, Kode: `STF`
   - Posisi: `Direktur`, Kode: `DIR`

Atau gunakan data sample yang sudah ada dari migration.

#### b. Tambah Karyawan
1. Menu: **Kelola Karyawan**
2. Tambahkan karyawan baru atau gunakan data sample

### 4. Test Form Perjalanan Dinas

1. Logout dari admin (atau buka incognito window)
2. Buka **http://localhost:3000**
3. Pilih karyawan dari dropdown
4. Isi form perjalanan dinas:
   - Maksud: "Kunjungan kerja ke kantor cabang"
   - Tujuan: "Jakarta"
   - Tanggal berangkat: (pilih tanggal)
   - Tanggal kembali: (pilih tanggal)
   - Angkutan: "Pesawat"
5. Klik **"Buat Perjalanan Dinas"**
6. Download PDF yang dihasilkan

### 5. Cek Data di Admin Panel

1. Login kembali sebagai admin
2. Menu: **Daftar Perjalanan Dinas**
3. Lihat data yang baru dibuat
4. Download PDF dari admin panel

## Struktur URL

- **Frontend (Public):** http://localhost:3000
- **Admin Login:** http://localhost:3000/admin/login
- **Admin Dashboard:** http://localhost:3000/admin/dashboard
- **Backend API:** http://localhost:8080/api
- **API Documentation:** http://localhost:8080/api (lihat README.md untuk endpoint list)

## Default Credentials

### Admin
- Username: `admin`
- Password: `admin123`

**⚠️ PENTING:** Ubah password default setelah login pertama kali!

## Sample Data

Jika Anda menjalankan migration, berikut data sample yang tersedia:

### Position Codes
- Direktur Utama: DIRUT
- Direktur: DIR
- Manager: MNG
- Supervisor: SPV
- Staff: STF
- Admin: ADM

### Employees
- Ahmad Budiman (NIP: 199001012015011001) - Direktur Utama
- Siti Nurhaliza (NIP: 199102022015011002) - Manager
- Budi Santoso (NIP: 199203033015011003) - Supervisor
- Dewi Lestari (NIP: 199304044015011004) - Staff
- Eko Prasetyo (NIP: 199405055015011005) - Staff

## Troubleshooting

### Backend tidak bisa start

**Error: "failed to connect to database"**
```bash
# Cek PostgreSQL running
sudo systemctl status postgresql  # Linux
brew services list  # macOS

# Cek database exists
psql -U postgres -l | grep perjalanan_dinas

# Cek credentials di .env
cat backend/.env
```

**Error: "bind: address already in use"**
```bash
# Port 8080 sudah digunakan, matikan aplikasi yang menggunakan port tersebut
# Atau ubah PORT di .env
```

### Frontend tidak bisa start

**Error: "ECONNREFUSED"**
- Pastikan backend sudah running di port 8080
- Cek API_URL di `frontend/.env.local`

**Error: "Cannot find module"**
```bash
# Install ulang dependencies
rm -rf node_modules package-lock.json
npm install
```

### Database Migration Issues

**Jika ingin reset database:**
```bash
# Drop database
dropdb perjalanan_dinas

# Buat ulang
createdb perjalanan_dinas

# Run migrations
cd backend
psql -U postgres -d perjalanan_dinas -f migrations/001_initial_schema.sql
psql -U postgres -d perjalanan_dinas -f migrations/002_seed_data.sql
```

### Error saat membuat perjalanan dinas

**"Position code not found for employee's position"**
- Pastikan jabatan karyawan sudah didaftarkan di menu Pengkodean Jabatan
- Nama jabatan harus sama persis (case-sensitive)

## Useful Commands

### Backend
```bash
# Run backend
cd backend && go run cmd/api/main.go

# Build binary
cd backend && go build -o bin/server cmd/api/main.go

# Run binary
cd backend && ./bin/server

# Using Makefile
cd backend && make run
```

### Frontend
```bash
# Development mode
cd frontend && npm run dev

# Production build
cd frontend && npm run build && npm start

# Lint check
cd frontend && npm run lint
```

### Database
```bash
# Connect to database
psql -U postgres -d perjalanan_dinas

# Backup database
pg_dump perjalanan_dinas > backup.sql

# Restore database
psql -U postgres -d perjalanan_dinas < backup.sql

# View all tables
psql -U postgres -d perjalanan_dinas -c "\dt"
```

## Next Steps

1. ✅ Aplikasi running di local
2. ✅ Login sebagai admin berhasil
3. ✅ Data master sudah disetup
4. ✅ Test membuat perjalanan dinas berhasil
5. ✅ PDF berhasil di-generate

**Selanjutnya:**
- Customize tampilan sesuai kebutuhan
- Tambah fitur berita acara perjalanan dinas
- Setup deployment ke production (lihat DEPLOYMENT.md)

## Support

Jika menemukan masalah:
1. Cek section Troubleshooting di atas
2. Baca README.md untuk dokumentasi lengkap
3. Cek logs di terminal backend dan frontend
4. Hubungi tim development

## Resources

- **README.md** - Dokumentasi lengkap
- **DEPLOYMENT.md** - Panduan deployment ke production
- **backend/migrations/** - SQL migrations
- **API Endpoints** - Lihat di README.md

---

Happy Coding! 🚀
