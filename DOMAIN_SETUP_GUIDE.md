# 🌐 PANDUAN SETUP DOMAIN SPD
# Domain: www.digitalbankingbankjatim.my.id/spd

---

## ✅ YANG SUDAH DIKONFIGURASI DI VPS:

### 1. Nginx Configuration
- ✅ Server name: `digitalbankingbankjatim.my.id` dan `www.digitalbankingbankjatim.my.id`
- ✅ Subdirectory: `/spd` untuk aplikasi
- ✅ Proxy backend API: `/spd/api` → `http://localhost:8080/api`
- ✅ Proxy frontend: `/spd` → `http://localhost:3000`
- ✅ Root domain redirect: `/` → `/spd`

### 2. Frontend Configuration (Next.js)
- ✅ basePath: `/spd`
- ✅ assetPrefix: `/spd`
- ✅ API URL: `http://digitalbankingbankjatim.my.id/spd`
- ✅ Frontend rebuilt dengan config baru
- ✅ Service restarted

### 3. SSL/HTTPS Tools
- ✅ Certbot terinstall
- ✅ Python3-certbot-nginx terinstall
- ⏳ **MENUNGGU DNS**: SSL certificate belum diinstall (butuh DNS aktif dulu)

---

## 🔴 YANG HARUS ANDA LAKUKAN SEKARANG:

### STEP 1: Setting DNS di Panel Jetorbit

**Login ke Panel Domain:**
1. Buka https://my.jetorbit.com atau https://panel.jetorbit.com
2. Login dengan akun Jetorbit Anda
3. Masuk ke menu **Domain** atau **Kelola Domain**
4. Pilih domain: **digitalbankingbankjatim.my.id**
5. Klik **DNS Management** atau **DNS Zone Editor**

**Tambahkan/Update DNS Records:**

Tambahkan 2 record ini:

| Type | Name/Host | Value/Points To | TTL |
|------|-----------|-----------------|-----|
| **A** | `@` | `103.235.75.196` | 3600 (atau 1 Hour) |
| **A** | `www` | `103.235.75.196` | 3600 (atau 1 Hour) |

**Penjelasan:**
- Record **@** = untuk akses `digitalbankingbankjatim.my.id` (tanpa www)
- Record **www** = untuk akses `www.digitalbankingbankjatim.my.id` (dengan www)
- IP: **103.235.75.196** = IP VPS Jetorbit Anda
- TTL: 3600 detik = 1 jam (waktu cache DNS)

**Hapus record lama (jika ada):**
- Hapus record A yang mengarah ke IP lain
- Hapus CNAME yang konflik

**Screenshot contoh:**
```
Type: A
Name: @
Value: 103.235.75.196
TTL: 3600

Type: A
Name: www
Value: 103.235.75.196
TTL: 3600
```

---

### STEP 2: Tunggu DNS Propagasi

**Waktu tunggu:**
- Minimal: 5-15 menit
- Normal: 30 menit - 2 jam
- Maksimal: 24 jam (jarang)

**Cara cek DNS sudah aktif:**

**Opsi 1 - Via Terminal/CMD:**
```bash
# Mac/Linux
nslookup www.digitalbankingbankjatim.my.id

# Atau
ping www.digitalbankingbankjatim.my.id

# Windows CMD
nslookup www.digitalbankingbankjatim.my.id
```

**Opsi 2 - Via Website:**
- Buka: https://dnschecker.org
- Masukkan: `www.digitalbankingbankjatim.my.id`
- Klik **Search**
- Pastikan IP menunjukkan: **103.235.75.196**
- Tunggu sampai semua region hijau (global propagation)

**Opsi 3 - Via Browser:**
- Buka: `http://www.digitalbankingbankjatim.my.id/spd`
- Jika muncul aplikasi SPD = DNS sudah aktif! ✅
- Jika muncul error/timeout = DNS belum propagasi, tunggu lagi

---

### STEP 3: Install SSL Certificate (HTTPS)

**⚠️ HANYA LAKUKAN SETELAH DNS SUDAH AKTIF!**

Setelah DNS propagasi selesai (bisa diakses via browser), jalankan command ini untuk install SSL:

```bash
# SSH ke VPS
ssh -i ~/.ssh/digitalbanking root@103.235.75.196

# Install SSL Certificate
certbot --nginx -d digitalbankingbankjatim.my.id -d www.digitalbankingbankjatim.my.id

# Ikuti instruksi:
# 1. Masukkan email untuk renewal notifications
# 2. Agree to Terms of Service (Y)
# 3. Share email with EFF (Y/N - terserah)
# 4. Redirect HTTP to HTTPS? Pilih: 2 (Yes, redirect)

# Certbot akan otomatis:
# - Generate SSL certificate
# - Update Nginx config
# - Setup auto-renewal
# - Restart Nginx
```

**Jika berhasil, Anda akan lihat:**
```
Successfully received certificate.
Certificate is saved at: /etc/letsencrypt/live/digitalbankingbankjatim.my.id/fullchain.pem
Key is saved at:         /etc/letsencrypt/live/digitalbankingbankjatim.my.id/privkey.pem
```

---

### STEP 4: Update Frontend untuk HTTPS

Setelah SSL aktif, update config frontend:

```bash
# SSH ke VPS (jika belum)
ssh -i ~/.ssh/digitalbanking root@103.235.75.196

# Update Next.js config untuk HTTPS
cd /var/www/perjalanan-dinas/frontend
nano next.config.mjs

# Ubah baris ini:
NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL || 'http://digitalbankingbankjatim.my.id/spd',

# Menjadi:
NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL || 'https://digitalbankingbankjatim.my.id/spd',

# Save: Ctrl+O, Enter, Ctrl+X

# Rebuild frontend
npm run build

# Restart service
systemctl restart spd-frontend

# Test
curl -s https://www.digitalbankingbankjatim.my.id/spd | grep "Sistem Perjalanan Dinas"
```

---

## 🎉 SETELAH SEMUA SELESAI:

### Akses Aplikasi:

**Frontend (Form SPD):**
- HTTP: `http://www.digitalbankingbankjatim.my.id/spd`
- HTTPS: `https://www.digitalbankingbankjatim.my.id/spd` ✅ (setelah SSL)
- Atau tanpa www: `https://digitalbankingbankjatim.my.id/spd`

**Admin Panel:**
- HTTPS: `https://www.digitalbankingbankjatim.my.id/spd/admin/login` ✅
- Username: `admin`
- Password: `admin123`
- **⚠️ GANTI PASSWORD SEGERA!**

**Backend API:**
- Endpoint: `https://www.digitalbankingbankjatim.my.id/spd/api/*`
- Health: `https://www.digitalbankingbankjatim.my.id/spd/api/positions`
- Login: `https://www.digitalbankingbankjatim.my.id/spd/api/auth/login`

---

## 🔧 TROUBLESHOOTING:

### Problem 1: DNS tidak resolve setelah 24 jam
**Solusi:**
- Cek apakah DNS record sudah benar di panel Jetorbit
- Pastikan tidak ada typo di nama domain
- Coba flush DNS cache:
  - Mac: `sudo dscacheutil -flushcache; sudo killall -HUP mDNSResponder`
  - Windows: `ipconfig /flushdns`
  - Linux: `sudo systemd-resolve --flush-caches`

### Problem 2: Certbot gagal generate certificate
**Error: "Fetching http://.../.well-known/acme-challenge/... Connection refused"**

**Solusi:**
```bash
# Pastikan Nginx running
systemctl status nginx

# Pastikan firewall allow port 80 dan 443
ufw status
ufw allow 80/tcp
ufw allow 443/tcp

# Restart Nginx
systemctl restart nginx

# Coba lagi certbot
certbot --nginx -d digitalbankingbankjatim.my.id -d www.digitalbankingbankjatim.my.id
```

### Problem 3: Aplikasi tidak muncul di /spd
**Solusi:**
```bash
# Check Nginx config
nginx -t

# Check Nginx logs
tail -50 /var/log/nginx/error.log

# Check frontend service
systemctl status spd-frontend

# Restart semua
systemctl restart spd-frontend
systemctl restart nginx
```

### Problem 4: CSS/JS tidak load (404 errors)
**Solusi:**
```bash
# Pastikan basePath sudah benar
cd /var/www/perjalanan-dinas/frontend
cat next.config.mjs | grep basePath

# Harus menunjukkan: basePath: '/spd'

# Rebuild frontend
npm run build
systemctl restart spd-frontend
```

### Problem 5: API tidak bisa diakses (404)
**Solusi:**
```bash
# Check backend service
systemctl status spd-backend

# Check backend logs
tail -50 /var/log/spd-backend.log

# Test direct backend
curl http://localhost:8080/api/positions

# Restart backend
systemctl restart spd-backend
```

---

## 📊 MONITORING & MAINTENANCE:

### Check Services Status:
```bash
ssh -i ~/.ssh/digitalbanking root@103.235.75.196 'systemctl status spd-backend spd-frontend nginx --no-pager'
```

### Check Logs:
```bash
# Backend logs
ssh -i ~/.ssh/digitalbanking root@103.235.75.196 'tail -50 /var/log/spd-backend.log'

# Frontend logs
ssh -i ~/.ssh/digitalbanking root@103.235.75.196 'tail -50 /var/log/spd-frontend.log'

# Nginx access logs
ssh -i ~/.ssh/digitalbanking root@103.235.75.196 'tail -50 /var/log/nginx/access.log'

# Nginx error logs
ssh -i ~/.ssh/digitalbanking root@103.235.75.196 'tail -50 /var/log/nginx/error.log'
```

### SSL Certificate Auto-Renewal:
```bash
# Check renewal timer
systemctl status certbot.timer

# Test renewal (dry-run)
certbot renew --dry-run

# Manual renewal (jika perlu)
certbot renew
```

Certificate akan auto-renew setiap 60 hari via systemd timer. Tidak perlu manual renewal!

---

## 📝 SUMMARY KONFIGURASI:

### File yang sudah dimodifikasi:

1. **`/etc/nginx/sites-available/spd`** (Nginx config)
   - Server name dengan domain
   - Location blocks untuk /spd
   - Proxy ke backend dan frontend

2. **`/var/www/perjalanan-dinas/frontend/next.config.mjs`**
   - basePath: `/spd`
   - assetPrefix: `/spd`
   - API URL dengan domain

### Services yang running:

- ✅ `spd-backend.service` - Go backend API (port 8080)
- ✅ `spd-frontend.service` - Next.js frontend (port 3000)
- ✅ `nginx.service` - Reverse proxy (port 80, 443 setelah SSL)
- ✅ `postgresql.service` - Database
- ✅ `certbot.timer` - Auto SSL renewal

### Ports yang digunakan:

- Port 80: HTTP (Nginx) → akan redirect ke 443 setelah SSL
- Port 443: HTTPS (Nginx dengan SSL)
- Port 8080: Backend API (internal only)
- Port 3000: Frontend Next.js (internal only)
- Port 5432: PostgreSQL (internal only)

---

## 🔐 SECURITY CHECKLIST:

- [ ] DNS sudah pointing ke VPS
- [ ] SSL certificate terinstall (HTTPS)
- [ ] Firewall aktif (ufw)
- [ ] Port 22 (SSH) hanya dari IP tertentu (opsional)
- [ ] Admin password sudah diganti dari default
- [ ] Database credentials aman (tidak di-commit ke git)
- [ ] SSH key-based authentication (sudah aktif ✅)

---

## 📞 BUTUH BANTUAN?

Jika ada masalah:

1. **Cek status services:**
   ```bash
   ssh -i ~/.ssh/digitalbanking root@103.235.75.196 'systemctl status spd-backend spd-frontend nginx'
   ```

2. **Cek logs untuk error:**
   ```bash
   ssh -i ~/.ssh/digitalbanking root@103.235.75.196 'tail -100 /var/log/spd-backend-error.log'
   ```

3. **Screenshot error** dan kirim ke saya

4. **Restart semua services:**
   ```bash
   ssh -i ~/.ssh/digitalbanking root@103.235.75.196 'systemctl restart spd-backend spd-frontend nginx'
   ```

---

**🎊 SELAMAT! Setup domain hampir selesai!**

**Langkah terakhir Anda:**
1. ✅ Setting DNS di panel Jetorbit (A record @ dan www → 103.235.75.196)
2. ⏳ Tunggu DNS propagasi (5-30 menit)
3. ✅ Install SSL dengan certbot (setelah DNS aktif)
4. ✅ Update frontend config ke HTTPS
5. 🎉 Aplikasi siap di `https://www.digitalbankingbankjatim.my.id/spd`

Total waktu: **~30-60 menit** (termasuk DNS propagasi)

---

Generated by: **Claude Code**
Tanggal: **27 Oktober 2025**
VPS IP: **103.235.75.196**
Domain: **www.digitalbankingbankjatim.my.id**
