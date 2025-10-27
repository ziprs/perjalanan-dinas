# 🚀 DEPLOYMENT KE RAILWAY.APP - 10 MENIT LIVE!

## 🎯 Mengapa Railway?

- ✅ **Deploy dalam 10 menit**
- ✅ **GRATIS $5 credit/bulan** (cukup untuk production kecil)
- ✅ **Auto SSL/HTTPS**
- ✅ **Domain gratis**: `*.railway.app`
- ✅ **PostgreSQL managed** (tidak perlu setup manual)
- ✅ **Support Go + Next.js**
- ✅ **Auto deploy dari GitHub**

---

## 📋 LANGKAH DEPLOYMENT

### STEP 1: Daftar Railway (2 menit)

1. Buka: https://railway.app
2. Klik **"Start a New Project"**
3. Login dengan **GitHub account** Anda
4. Authorize Railway akses ke GitHub

**Catatan**: Tidak perlu credit card untuk trial!

---

### STEP 2: Push Code ke GitHub (3 menit)

Jika belum punya GitHub repo, buat dulu:

```bash
# Di komputer Anda, masuk ke folder project
cd "/Users/unclejss/Documents/Project SPD/perjalanan-dinas"

# Initialize Git (jika belum)
git init

# Add all files
git add .

# Commit
git commit -m "Initial commit - SPD Bank Jatim"

# Create repo di GitHub
# 1. Buka: https://github.com/new
# 2. Nama repo: perjalanan-dinas
# 3. Pilih: Private
# 4. Klik: Create repository

# Push ke GitHub (ganti USERNAME dengan username GitHub Anda)
git remote add origin https://github.com/USERNAME/perjalanan-dinas.git
git branch -M main
git push -u origin main
```

---

### STEP 3: Deploy Backend ke Railway (2 menit)

1. **Di Railway Dashboard**, klik **"New Project"**
2. Pilih **"Deploy from GitHub repo"**
3. Pilih repo: `perjalanan-dinas`
4. Railway akan auto-detect Go backend

#### Konfigurasi Backend:

Setelah deploy dimulai, klik service backend dan tambahkan **Environment Variables**:

```env
PORT=8080
DB_HOST=${{Postgres.PGHOST}}
DB_PORT=${{Postgres.PGPORT}}
DB_USER=${{Postgres.PGUSER}}
DB_PASSWORD=${{Postgres.PGPASSWORD}}
DB_NAME=${{Postgres.PGDATABASE}}
JWT_SECRET=your_super_secret_jwt_key_here_change_this_123456789
```

**Root Directory**: Set ke `/backend`

**Build Command**: (otomatis detect)
```bash
go build -o main cmd/api/main.go
```

**Start Command**:
```bash
./main
```

---

### STEP 4: Tambah PostgreSQL Database (1 menit)

1. Di project Railway, klik **"+ New"**
2. Pilih **"Database"**
3. Pilih **"PostgreSQL"**
4. Railway akan otomatis create database dan set environment variables

**Database akan otomatis terkoneksi ke backend!**

---

### STEP 5: Deploy Frontend ke Railway (2 menit)

1. Klik **"+ New"** lagi
2. Pilih **"GitHub Repo"** → pilih repo yang sama
3. Tapi kali ini untuk **frontend**

#### Konfigurasi Frontend:

**Root Directory**: Set ke `/frontend`

**Environment Variables**:
```env
NEXT_PUBLIC_API_URL=${{Backend.RAILWAY_PUBLIC_DOMAIN}}
```

**Build Command**:
```bash
npm run build
```

**Start Command**:
```bash
npm start
```

---

### STEP 6: Generate Domain & SSL (Otomatis!)

Railway otomatis generate:

- **Backend**: `https://backend-xxx.railway.app`
- **Frontend**: `https://frontend-xxx.railway.app`
- **SSL**: Otomatis HTTPS ✅

---

## 🎉 SELESAI! APLIKASI SUDAH LIVE!

### Akses Aplikasi:

- **Frontend**: `https://spd-bankjatim-xxx.railway.app`
- **Admin Login**: `https://spd-bankjatim-xxx.railway.app/admin/login`
  - Username: `admin`
  - Password: `admin123`

---

## 🔧 MANAGEMENT

### View Logs:

1. Klik service (Backend/Frontend)
2. Tab **"Deployments"**
3. Klik deployment terbaru
4. Tab **"Logs"** untuk real-time logs

### Restart Service:

1. Klik service
2. Tab **"Settings"**
3. Klik **"Restart"**

### Update Code:

```bash
# Push ke GitHub
git add .
git commit -m "Update: feature baru"
git push

# Railway otomatis detect dan re-deploy! 🚀
```

---

## 💰 ESTIMASI BIAYA

### Free Tier:
- **$5 credit/bulan** (~ Rp 75.000)
- Cukup untuk:
  - 500 jam runtime
  - 1GB RAM
  - PostgreSQL 1GB storage
  - **Cocok untuk production kecil-menengah**

### Jika Credit Habis:
Bisa upgrade dengan bayar usage:
- **~$5-10/bulan** untuk traffic normal
- **~Rp 75.000 - 150.000/bulan**

---

## 🔒 SECURITY CHECKLIST

- ✅ HTTPS otomatis aktif
- ✅ PostgreSQL managed & secure
- ✅ Environment variables encrypted
- ⚠️ **WAJIB**: Ganti admin password setelah deploy
- ⚠️ **WAJIB**: Ganti JWT_SECRET ke random string

---

## 🎯 CUSTOM DOMAIN (Opsional)

Jika ingin pakai domain sendiri (contoh: `spd.bankjatim.co.id`):

1. **Beli domain** atau gunakan domain Bank Jatim
2. Di Railway service → **Settings** → **Domains**
3. Klik **"Custom Domain"**
4. Masukkan: `spd.bankjatim.co.id`
5. Update DNS di domain registrar:
   ```
   Type: CNAME
   Name: spd
   Value: [railway-provided-cname].railway.app
   ```

SSL otomatis provision untuk custom domain!

---

## 📊 MONITORING

Railway provides built-in monitoring:

1. **Metrics**: CPU, Memory, Network usage
2. **Logs**: Real-time application logs
3. **Deployments**: History & rollback capability

### External Monitoring (Gratis):

**UptimeRobot**: https://uptimerobot.com
- Monitor uptime
- Email alerts jika down
- Gratis untuk 50 monitors

---

## 🆘 TROUBLESHOOTING

### Backend tidak connect ke database?

Pastikan environment variables sudah set:
```env
DB_HOST=${{Postgres.PGHOST}}
DB_PORT=${{Postgres.PGPORT}}
DB_USER=${{Postgres.PGUSER}}
DB_PASSWORD=${{Postgres.PGPASSWORD}}
DB_NAME=${{Postgres.PGDATABASE}}
```

### Frontend tidak bisa hit backend API?

Update frontend env:
```env
NEXT_PUBLIC_API_URL=https://[your-backend-url].railway.app
```

### Build failed?

Cek logs di tab "Deployments" untuk error message.

---

## 🔄 MIGRASI KE JETORBIT NANTI

Jika VPS Jetorbit sudah siap, Anda bisa:

1. Export database dari Railway:
   ```bash
   # Get database URL dari Railway
   pg_dump DATABASE_URL > backup.sql
   ```

2. Import ke Jetorbit VPS:
   ```bash
   psql -h IP_VPS -U spduser -d perjalanan_dinas < backup.sql
   ```

3. Deploy code ke Jetorbit menggunakan `deploy.sh`

---

## 📞 SUPPORT

**Railway Documentation**: https://docs.railway.app
**Community Discord**: https://discord.gg/railway

---

## ✅ KEUNTUNGAN RAILWAY vs VPS

| Feature | Railway | Jetorbit VPS |
|---------|---------|--------------|
| Setup Time | 10 menit | 30-60 menit |
| SSL | Otomatis | Manual setup |
| Database | Managed | Manual setup |
| Auto Deploy | ✅ Git push | ❌ Manual |
| Monitoring | Built-in | Manual setup |
| Backup | Auto snapshot | Manual script |
| Scaling | 1-click | Manual |
| Maintenance | Zero | Manual updates |
| Harga | $5-10/bulan | Rp 100k/bulan |

**Rekomendasi**:
- **Railway** untuk development & production kecil
- **Jetorbit VPS** untuk full control & production besar

---

**🎊 Selamat! Aplikasi Anda akan live dalam 10 menit dengan Railway! 🎊**
