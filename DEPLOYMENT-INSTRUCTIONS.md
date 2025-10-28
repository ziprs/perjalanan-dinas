# 🚀 Panduan Deployment ke Production VPS

## Status Update Terbaru
✅ **Code berhasil di-push ke GitHub repository**
✅ **Script deployment otomatis sudah dibuat**
✅ **Siap untuk di-deploy ke production VPS**

---

## 📋 Fitur Baru yang Akan Di-Deploy

### 1. **Form At-Cost di Halaman Public**
- User bisa membuat klaim At-Cost tanpa login sebagai admin
- Tab navigation untuk memisahkan form SPD dan At-Cost
- Dropdown pilihan perjalanan dinas dan pegawai

### 2. **Public API Endpoints**
- Endpoint At-Cost tidak memerlukan authentication
- Upload receipt dan create claim bisa dilakukan tanpa login
- API routes yang sudah dibuat public:
  - `POST /api/at-cost/upload-receipt`
  - `POST /api/at-cost/claims`
  - `GET /api/at-cost/claims`
  - `GET /api/travel-requests`

### 3. **Travel Agency Receipt Support**
- Support parsing receipt dari Sonyloka/tiketqta.com
- Format tanggal DD/MM/YYYY
- Format harga "Rp. 485.230"
- Deteksi hotel booking voucher

### 4. **Download Buttons**
- Tombol download PDF muncul otomatis setelah submit klaim
- Download Nota At-Cost
- Download PDF Lengkap (Nota + semua Receipt)
- Tombol "Buat Klaim Baru"

### 5. **PDF Formatting Improvements**
- Nama representative sekarang **bold dan underline**
- Berlaku untuk Nota Permintaan dan Berita Acara

### 6. **Admin Panel Updates**
- Status approval dihapus dari UI admin At-Cost
- Tetap ada fungsi view, detail, download, dan delete

---

## 🔧 Cara Deploy ke Production VPS

### Opsi 1: Menggunakan Script Otomatis (Recommended)

1. **Login ke VPS** (gunakan method yang tersedia - SSH console, cPanel terminal, dll)

2. **Navigate ke project directory:**
   ```bash
   cd /home/admin/perjalanan-dinas
   ```

3. **Jalankan deployment script:**
   ```bash
   bash deploy-production.sh
   ```

   Script akan otomatis:
   - ✅ Pull code terbaru dari Git
   - ✅ Install dependencies (Go & npm)
   - ✅ Build backend (Go binary)
   - ✅ Build frontend (Next.js production)
   - ✅ Restart services
   - ✅ Check service status

---

### Opsi 2: Manual Step-by-Step

Jika script otomatis tidak bekerja, ikuti langkah manual berikut:

#### Step 1: Update Code
```bash
cd /home/admin/perjalanan-dinas
git pull origin main
```

#### Step 2: Deploy Backend
```bash
cd backend

# Install dependencies
go mod tidy

# Build aplikasi
go build -o main cmd/api/main.go

# Restart service
sudo systemctl restart spd-backend

# Check status
sudo systemctl status spd-backend
```

#### Step 3: Deploy Frontend
```bash
cd ../frontend

# Install dependencies
npm install

# Build production
npm run build

# Restart service
sudo systemctl restart spd-frontend

# Check status
sudo systemctl status spd-frontend
```

#### Step 4: Verify Deployment
```bash
# Check both services
sudo systemctl status spd-backend
sudo systemctl status spd-frontend

# Check logs if needed
journalctl -u spd-backend -n 50
journalctl -u spd-frontend -n 50
```

---

## 🌐 URL Production

Setelah deployment selesai, aplikasi akan berjalan di:

- **Frontend (Public)**: http://103.160.37.195:3000
- **Backend API**: http://103.160.37.195:8080
- **Admin Panel**: http://103.160.37.195:3000/admin/login

---

## ✅ Testing Checklist

Setelah deployment, test fitur-fitur berikut:

### Halaman Public (http://103.160.37.195:3000)
- [ ] Tab "Buat Nota Perjalanan Dinas" berfungsi
- [ ] Tab "Buat Klaim At-Cost" muncul
- [ ] Dropdown "Pilih Perjalanan Dinas" terisi
- [ ] Upload receipt PDF berhasil
- [ ] Parsing receipt menampilkan data yang benar
- [ ] Submit klaim At-Cost berhasil
- [ ] Tombol download PDF muncul setelah submit
- [ ] Download Nota At-Cost berfungsi
- [ ] Download PDF Lengkap berfungsi

### Admin Panel (http://103.160.37.195:3000/admin)
- [ ] Login admin berhasil
- [ ] Menu "Klaim At-Cost" ada di sidebar
- [ ] List klaim At-Cost tampil
- [ ] Detail klaim menampilkan semua receipt
- [ ] Download buttons di detail page berfungsi
- [ ] Tidak ada status approval di UI
- [ ] Delete klaim berfungsi

### PDF Documents
- [ ] Nama representative **bold dan underline**
- [ ] Nota Permintaan format benar
- [ ] Berita Acara format benar
- [ ] Nota At-Cost format benar
- [ ] Combined PDF berisi nota + semua receipt

---

## 🔍 Troubleshooting

### Problem: Service tidak start
```bash
# Check logs
journalctl -u spd-backend -n 100
journalctl -u spd-frontend -n 100

# Check if port sudah digunakan
sudo lsof -i :8080  # Backend
sudo lsof -i :3000  # Frontend

# Restart services
sudo systemctl restart spd-backend
sudo systemctl restart spd-frontend
```

### Problem: Build frontend gagal
```bash
cd /home/admin/perjalanan-dinas/frontend

# Clear cache
rm -rf .next node_modules package-lock.json

# Reinstall
npm install
npm run build
```

### Problem: Database connection error
```bash
# Check PostgreSQL
sudo systemctl status postgresql

# Test connection
psql -U admin -d perjalanan_dinas -c "SELECT 1;"
```

---

## 📞 Support

Jika ada masalah saat deployment, cek:
1. Logs di journalctl
2. File .env di backend dan frontend
3. Database connection
4. Port availability

---

**Last Updated**: 28 Oktober 2025
**Commit Hash**: 75be2ff
**GitHub Repo**: https://github.com/ziprs/perjalanan-dinas
