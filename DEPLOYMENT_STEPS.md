# 🚀 LANGKAH DEPLOYMENT KE RAILWAY - IKUTI INI!

## ✅ STATUS SAAT INI:
- ✅ Git repository sudah initialized
- ✅ Code sudah di-commit (86 files)
- ⏳ **NEXT**: Push ke GitHub → Deploy ke Railway

---

## 📝 STEP 2: BUAT GITHUB REPOSITORY & PUSH CODE

### A. Buat Repository di GitHub (3 menit)

1. **Buka browser**, pergi ke: https://github.com/new

2. **Isi form repository**:
   - **Repository name**: `perjalanan-dinas`
   - **Description**: `Sistem Perjalanan Dinas - Bank Jatim Digital Banking Division`
   - **Visibility**: Pilih **Private** (untuk keamanan)
   - **❌ JANGAN centang** "Add README" atau ".gitignore" atau "license"

3. **Klik "Create repository"**

4. **COPY URL** yang muncul, contoh:
   ```
   https://github.com/USERNAME/perjalanan-dinas.git
   ```
   (Ganti USERNAME dengan username GitHub Anda)

---

### B. Push Code ke GitHub (2 menit)

Setelah repository dibuat, jalankan command berikut di Terminal:

```bash
# Masuk ke folder project
cd "/Users/unclejss/Documents/Project SPD/perjalanan-dinas"

# Add remote GitHub (GANTI USERNAME!)
git remote add origin https://github.com/USERNAME/perjalanan-dinas.git

# Push ke GitHub
git branch -M main
git push -u origin main
```

**Jika diminta username & password:**
- Username: `[GitHub username Anda]`
- Password: **BUKAN password GitHub biasa**, tapi **Personal Access Token**

**Cara buat Personal Access Token** (jika belum punya):
1. GitHub → Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Generate new token (classic)
3. Centang: `repo` (full control)
4. Generate token
5. COPY token dan simpan (tidak bisa dilihat lagi!)
6. Pakai token ini sebagai password saat git push

---

## 📝 STEP 3: DEPLOY KE RAILWAY (5 menit)

### A. Daftar Railway

1. Buka: https://railway.app
2. Klik **"Login"** atau **"Start a New Project"**
3. Pilih **"Login with GitHub"**
4. Authorize Railway

---

### B. Create New Project

1. Di Railway dashboard, klik **"New Project"**
2. Pilih **"Deploy from GitHub repo"**
3. Pilih repository: **`perjalanan-dinas`**
4. Railway akan auto-detect dan mulai deploy

---

### C. Deploy Backend

Railway akan auto-detect Go application di folder `backend`.

**Setup Backend Service:**

1. Klik service **backend** yang muncul
2. Klik tab **"Variables"**
3. Klik **"+ New Variable"**
4. Tambahkan environment variables satu per satu:

```env
PORT=8080
JWT_SECRET=super_secret_jwt_key_change_this_in_production_12345678
```

**Catatan**: Database variables akan otomatis muncul setelah kita add PostgreSQL.

5. Klik tab **"Settings"**
6. Scroll ke **"Root Directory"**
7. Set ke: **`/backend`**
8. **Build Command**: Biarkan kosong (auto-detect)
9. **Start Command**: `./main`

---

### D. Add PostgreSQL Database

1. Klik **"+ New"** di project
2. Pilih **"Database"**
3. Pilih **"Add PostgreSQL"**
4. Railway akan auto-provision database

**Update Backend Variables:**

Setelah PostgreSQL dibuat, kembali ke backend service → Variables:

Tambahkan variables ini (Railway auto-fill dari PostgreSQL):

```env
DB_HOST=${{Postgres.PGHOST}}
DB_PORT=${{Postgres.PGPORT}}
DB_USER=${{Postgres.PGUSER}}
DB_PASSWORD=${{Postgres.PGPASSWORD}}
DB_NAME=${{Postgres.PGDATABASE}}
```

Cara add:
- Ketik `DB_HOST` di Name
- Ketik `${{Postgres.PGHOST}}` di Value (Railway auto-replace dengan nilai real)

Ulangi untuk semua DB variables.

---

### E. Deploy Frontend

1. Klik **"+ New"** lagi
2. Pilih **"GitHub Repo"**
3. Pilih repository yang sama: **`perjalanan-dinas`**
4. Railway akan create service baru

**Setup Frontend Service:**

1. Klik service **frontend**
2. Klik tab **"Variables"**
3. Tambahkan:

```env
NEXT_PUBLIC_API_URL=${{backend.RAILWAY_PUBLIC_DOMAIN}}
NODE_ENV=production
```

4. Klik tab **"Settings"**
5. **Root Directory**: Set ke **`/frontend`**
6. **Build Command**: `npm install && npm run build`
7. **Start Command**: `npm start`

---

### F. Generate Public URLs

1. Klik **backend service** → tab **"Settings"**
2. Scroll ke **"Networking"**
3. Klik **"Generate Domain"**
4. Copy URL yang muncul (contoh: `backend-production-xxx.railway.app`)

5. Klik **frontend service** → tab **"Settings"**
6. Scroll ke **"Networking"**
7. Klik **"Generate Domain"**
8. Copy URL yang muncul (contoh: `frontend-production-xxx.railway.app`)

---

## 🎉 SELESAI! APLIKASI SUDAH LIVE!

### Akses Aplikasi:

**Frontend**: `https://frontend-production-xxx.railway.app`
- Halaman utama form SPD

**Backend API**: `https://backend-production-xxx.railway.app/api/health`
- Test health endpoint

**Admin Panel**: `https://frontend-production-xxx.railway.app/admin/login`
- Username: `admin`
- Password: `admin123`
- **⚠️ SEGERA GANTI PASSWORD SETELAH LOGIN!**

---

## 🔍 VERIFIKASI DEPLOYMENT

### Check Backend:
```bash
# Health check
curl https://backend-production-xxx.railway.app/api/health

# Expected response:
# {"status":"healthy","timestamp":"..."}
```

### Check Frontend:
Buka browser: `https://frontend-production-xxx.railway.app`

### Check Logs:
1. Di Railway, klik service (backend/frontend)
2. Tab **"Deployments"**
3. Klik deployment terbaru
4. Tab **"Logs"** untuk real-time logs

---

## 🆘 TROUBLESHOOTING

### Backend tidak bisa connect ke database?
- Pastikan PostgreSQL service sudah running
- Pastikan environment variables sudah benar
- Cek logs backend untuk error message

### Frontend tidak bisa hit backend API?
- Pastikan `NEXT_PUBLIC_API_URL` pointing ke backend URL
- Pastikan backend sudah generate public domain
- Redeploy frontend setelah update env variable

### Build failed?
- Cek logs untuk error message
- Pastikan Root Directory sudah benar
- Pastikan dependencies di package.json/go.mod complete

---

## 📊 MONITORING

Railway provides:
- **Metrics**: CPU, Memory, Network usage
- **Logs**: Real-time application logs
- **Deployments**: History & rollback

---

## 💰 ESTIMASI BIAYA

**Free Tier**: $5 credit/bulan
- Cukup untuk ~500 jam runtime
- Backend + Frontend + PostgreSQL
- **Cocok untuk production kecil-menengah**

**Usage charges** (setelah credit habis):
- ~$5-10/bulan untuk traffic normal
- ~Rp 75.000 - 150.000/bulan

---

## 🔄 UPDATE APLIKASI NANTI

Setelah deploy, update sangat mudah:

```bash
# Di komputer lokal
cd "/Users/unclejss/Documents/Project SPD/perjalanan-dinas"

# Edit code Anda
# ...

# Commit & push
git add .
git commit -m "Update: fitur baru"
git push

# Railway otomatis detect dan re-deploy! 🚀
```

---

## ✅ CHECKLIST DEPLOYMENT

- [ ] GitHub repository created
- [ ] Code pushed to GitHub
- [ ] Railway account created (via GitHub)
- [ ] Railway project created
- [ ] Backend service deployed
- [ ] PostgreSQL database added
- [ ] Backend environment variables set
- [ ] Frontend service deployed
- [ ] Frontend environment variables set
- [ ] Public domains generated
- [ ] Application tested & working
- [ ] Admin password changed

---

**🎊 SELAMAT! APLIKASI ANDA AKAN LIVE DALAM BEBERAPA MENIT! 🎊**

Total waktu: **10-15 menit**
Total biaya: **GRATIS** (dengan $5 credit/bulan)

---

## 📞 BUTUH BANTUAN?

Jika ada error atau stuck di step manapun, screenshot error-nya dan share ke saya. Saya akan bantu troubleshoot!

Railway Documentation: https://docs.railway.app
Railway Discord: https://discord.gg/railway
