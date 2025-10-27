# 🔄 PANDUAN UPDATE APLIKASI SPD

## 📝 Kapan Perlu Update?

Update aplikasi ketika ada:
- ✏️ Perubahan code (fitur baru, bug fix, perbaikan tampilan)
- 🔧 Update konfigurasi
- 📦 Update dependencies (package.json, go.mod)

---

## 🚀 CARA 1: UPDATE OTOMATIS (RECOMMENDED)

Cara termudah dan tercepat - hanya 1 command!

### Step-by-step:

**1. Edit code di komputer lokal Anda**
```bash
cd "/Users/unclejss/Documents/Project SPD/perjalanan-dinas"

# Edit file yang perlu diubah
# Contoh: edit frontend/src/app/page.tsx
```

**2. Commit dan push ke GitHub**
```bash
git add .
git commit -m "Update: deskripsi perubahan Anda"
git push
```

**3. SSH ke VPS dan jalankan update script**
```bash
ssh -i ~/.ssh/digitalbanking root@103.235.75.196

# Jalankan update script
cd /var/www/perjalanan-dinas
./update.sh
```

**Done!** ✅ Aplikasi sudah ter-update!

---

## 🛠️ CARA 2: UPDATE MANUAL (Step by Step)

Jika ingin lebih kontrol atau troubleshooting.

### A. Update Backend Only

Jika hanya mengubah code Go backend:

```bash
# SSH ke VPS
ssh -i ~/.ssh/digitalbanking root@103.235.75.196

# Masuk ke folder project
cd /var/www/perjalanan-dinas

# Pull code terbaru
git pull

# Rebuild backend
cd backend
go build -o spd-backend cmd/api/main.go

# Restart service
systemctl restart spd-backend

# Cek status
systemctl status spd-backend --no-pager

# Test backend
curl http://localhost:8080/api/positions
```

### B. Update Frontend Only

Jika hanya mengubah code Next.js frontend:

```bash
# SSH ke VPS
ssh -i ~/.ssh/digitalbanking root@103.235.75.196

# Masuk ke folder project
cd /var/www/perjalanan-dinas

# Pull code terbaru
git pull

# Rebuild frontend
cd frontend
npm install  # Hanya jika ada update package.json
npm run build

# Restart service
systemctl restart spd-frontend

# Cek status
systemctl status spd-frontend --no-pager

# Test frontend
curl -I http://localhost:3000
```

### C. Update Both (Backend + Frontend)

```bash
# SSH ke VPS
ssh -i ~/.ssh/digitalbanking root@103.235.75.196

# Masuk ke folder project
cd /var/www/perjalanan-dinas

# Pull code terbaru
git pull

# Update backend
cd backend
go build -o spd-backend cmd/api/main.go
systemctl restart spd-backend

# Update frontend
cd ../frontend
npm install  # Hanya jika ada update package.json
npm run build
systemctl restart spd-frontend

# Cek semua services
systemctl status spd-backend spd-frontend nginx --no-pager

# Test aplikasi
curl -I http://103.235.75.196/spd
```

---

## ⚡ WORKFLOW LENGKAP (Best Practice)

### Dari Komputer Lokal:

```bash
# 1. Pastikan di folder project
cd "/Users/unclejss/Documents/Project SPD/perjalanan-dinas"

# 2. Pull update terbaru (jika ada)
git pull

# 3. Buat branch baru untuk fitur/fix (optional tapi recommended)
git checkout -b feature/nama-fitur

# 4. Edit code Anda
# ... edit files ...

# 5. Test di local (optional)
# Backend:
cd backend
go run cmd/api/main.go

# Frontend (terminal baru):
cd frontend
npm run dev

# 6. Commit changes
git add .
git commit -m "Update: deskripsi lengkap perubahan"

# 7. Push ke GitHub
git push origin feature/nama-fitur
# Atau jika langsung ke main:
git checkout main
git merge feature/nama-fitur
git push origin main
```

### Di VPS (Production):

```bash
# SSH ke VPS
ssh -i ~/.ssh/digitalbanking root@103.235.75.196

# Run update script
cd /var/www/perjalanan-dinas
./update.sh

# Atau manual:
git pull
cd backend && go build -o spd-backend cmd/api/main.go && systemctl restart spd-backend
cd ../frontend && npm run build && systemctl restart spd-frontend
```

---

## 🔍 TROUBLESHOOTING UPDATE

### Problem 1: Git Pull Conflict

```bash
# Jika ada conflict saat git pull:
git stash              # Simpan perubahan lokal
git pull               # Pull update
git stash pop          # Restore perubahan lokal
# Atau jika ingin reset ke versi GitHub:
git reset --hard origin/main
git pull
```

### Problem 2: Build Failed

```bash
# Backend build error:
cd /var/www/perjalanan-dinas/backend
go mod tidy            # Fix dependencies
go build -o spd-backend cmd/api/main.go

# Frontend build error:
cd /var/www/perjalanan-dinas/frontend
rm -rf node_modules .next
npm install
npm run build
```

### Problem 3: Service Won't Start

```bash
# Cek log error:
journalctl -u spd-backend -n 50 --no-pager
journalctl -u spd-frontend -n 50 --no-pager

# Atau cek file log:
tail -50 /var/log/spd-backend-error.log
tail -50 /var/log/spd-frontend.log

# Restart service:
systemctl restart spd-backend
systemctl restart spd-frontend
systemctl restart nginx
```

### Problem 4: Changes Not Visible

```bash
# Clear browser cache:
# - Chrome: Ctrl+Shift+Delete
# - Firefox: Ctrl+Shift+Delete
# - Safari: Command+Option+E

# Atau force reload:
# - Chrome/Firefox: Ctrl+Shift+R
# - Safari: Command+Shift+R

# Clear Next.js cache di VPS:
cd /var/www/perjalanan-dinas/frontend
rm -rf .next
npm run build
systemctl restart spd-frontend
```

---

## 📊 MONITORING SETELAH UPDATE

```bash
# Cek semua services running:
systemctl is-active spd-backend spd-frontend nginx

# Cek aplikasi bisa diakses:
curl -I http://103.235.75.196/spd
curl -I http://103.235.75.196/spd/admin/login

# Cek API backend:
curl http://103.235.75.196/spd/api/positions

# Monitor logs real-time:
tail -f /var/log/spd-backend.log
tail -f /var/log/spd-frontend.log
tail -f /var/log/nginx/access.log
```

---

## 🔐 UPDATE ENVIRONMENT VARIABLES

Jika perlu update konfigurasi (database, API keys, etc):

```bash
# SSH ke VPS
ssh -i ~/.ssh/digitalbanking root@103.235.75.196

# Update backend .env
nano /var/www/perjalanan-dinas/backend/.env
# Edit, save (Ctrl+O, Enter, Ctrl+X)
systemctl restart spd-backend

# Update frontend .env.local
nano /var/www/perjalanan-dinas/frontend/.env.local
# Edit, save
npm run build
systemctl restart spd-frontend
```

---

## 📦 UPDATE DEPENDENCIES

### Update Go Dependencies (Backend)

```bash
cd /var/www/perjalanan-dinas/backend

# Update specific package:
go get -u github.com/gin-gonic/gin

# Update all packages:
go get -u ./...

# Clean up:
go mod tidy

# Rebuild:
go build -o spd-backend cmd/api/main.go
systemctl restart spd-backend
```

### Update Node Dependencies (Frontend)

```bash
cd /var/www/perjalanan-dinas/frontend

# Update specific package:
npm update next

# Update all packages:
npm update

# Or use npm-check-updates:
npx npm-check-updates -u
npm install

# Rebuild:
npm run build
systemctl restart spd-frontend
```

---

## 🗄️ UPDATE DATABASE SCHEMA

Jika ada perubahan struktur database:

```bash
# SSH ke VPS
ssh -i ~/.ssh/digitalbanking root@103.235.75.196

# Backup database dulu!
pg_dump -U spduser perjalanan_dinas > /tmp/backup_$(date +%Y%m%d_%H%M%S).sql

# Cara 1: Auto migration (sudah di code backend)
# Backend akan auto-migrate saat startup
systemctl restart spd-backend

# Cara 2: Manual SQL
psql -U spduser -d perjalanan_dinas
# Jalankan SQL commands...
# \q untuk keluar

# Verify:
psql -U spduser -d perjalanan_dinas -c "\dt"
```

---

## 🎯 CHECKLIST UPDATE

Gunakan checklist ini setiap kali update:

```
SEBELUM UPDATE:
□ Backup database (jika ada perubahan schema)
□ Catat versi aplikasi saat ini
□ Test di local dulu (optional)
□ Commit dan push ke GitHub

SAAT UPDATE:
□ SSH ke VPS
□ Pull latest code
□ Rebuild backend (jika ada perubahan)
□ Rebuild frontend (jika ada perubahan)
□ Restart services
□ Cek status services

SETELAH UPDATE:
□ Test akses aplikasi via browser
□ Test login admin
□ Test fitur yang diubah
□ Monitor logs untuk error
□ Verifikasi di berbagai device (desktop, mobile)
```

---

## 📞 BANTUAN

Jika mengalami masalah saat update:

1. **Cek logs** untuk lihat error message
2. **Screenshot error** dan simpan
3. **Rollback** jika perlu:
   ```bash
   cd /var/www/perjalanan-dinas
   git log --oneline -5
   git checkout [commit-hash-sebelumnya]
   ./update.sh
   ```

---

## 🚦 TIPS & BEST PRACTICES

1. **Selalu backup** sebelum update besar
2. **Test di local** sebelum push ke production
3. **Update di jam sepi** untuk minimize downtime
4. **Monitor logs** setelah update
5. **Dokumentasikan** setiap perubahan di commit message
6. **Gunakan branch** untuk fitur besar (git flow)
7. **Test di berbagai browser** setelah update frontend

---

**Generated by: Claude Code**
**Last Updated: 27 Oktober 2025**
