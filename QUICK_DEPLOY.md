# ⚡ QUICK DEPLOYMENT GUIDE - 15 MENIT LIVE!

## 🎯 Ringkasan
Deploy aplikasi SPD Bank Jatim ke VPS Jetorbit dengan domain GRATIS dalam 15 menit!

---

## 💰 TOTAL BIAYA: Rp 100.000/BULAN

- ✅ VPS Jetorbit: Rp 100.000/bulan
- ✅ Domain Freenom: GRATIS
- ✅ SSL Let's Encrypt: GRATIS
- ✅ Total: **Rp 1.2 juta/tahun**

---

## 📝 LANGKAH DEPLOYMENT

### STEP 1: Beli VPS Jetorbit (5 menit)

1. **Buka**: https://jetorbit.com
2. **Pilih**: VPS SSD 1
   - CPU: 1 Core
   - RAM: 1 GB
   - Storage: 25 GB SSD
   - Harga: **Rp 100.000/bulan**
3. **OS**: Ubuntu 22.04 LTS
4. **Checkout** dan bayar
5. **Tunggu email** dari Jetorbit berisi:
   - IP VPS: `xxx.xxx.xxx.xxx`
   - Username: `root`
   - Password: `xxxxxx`

---

### STEP 2: Domain Gratis dari Freenom (3 menit)

1. **Buka**: https://freenom.com
2. **Cari domain**: ketik `spd-bankjatim`
3. **Pilih ekstensi GRATIS**:
   - `.tk` (Tokelau) ✅ Recommended
   - `.ml` (Mali)
   - `.ga` (Gabon)
   - `.cf` (Central African Republic)
   - `.gq` (Equatorial Guinea)
4. **Check availability** → Pilih domain
5. **Duration**: 12 Months @ FREE
6. **Get it now** → Checkout
7. **Register/Login**
8. **Complete Order**
9. Domain Anda: `spd-bankjatim.tk` ✅

---

### STEP 3: Upload Code ke VPS (2 menit)

#### Opsi A: Via SFTP/FileZilla (Recommended untuk Pemula)

1. **Download FileZilla**: https://filezilla-project.org
2. **Connect ke VPS**:
   - Host: `xxx.xxx.xxx.xxx` (IP VPS dari Jetorbit)
   - Username: `root`
   - Password: (dari email Jetorbit)
   - Port: `22`
3. **Upload folder** `perjalanan-dinas` ke `/var/www/`

#### Opsi B: Via SCP (Command Line)

```bash
# Dari komputer Anda (Mac/Linux)
scp -r /Users/unclejss/Documents/Project\ SPD/perjalanan-dinas root@IP_VPS:/var/www/

# Contoh:
# scp -r /Users/unclejss/Documents/Project\ SPD/perjalanan-dinas root@103.123.45.67:/var/www/
```

---

### STEP 4: Deploy Otomatis (5 menit)

1. **SSH ke VPS**:
   ```bash
   ssh root@IP_VPS_ANDA
   # Masukkan password dari email Jetorbit
   ```

2. **Jalankan script deployment**:
   ```bash
   cd /var/www/perjalanan-dinas
   chmod +x deploy.sh
   ./deploy.sh
   ```

3. **Ikuti prompt**:
   - Database name: `perjalanan_dinas` (Enter)
   - Database user: `spduser` (Enter)
   - Database password: `[masukkan password kuat]`
   - Domain name: `spd-bankjatim.tk`
   - Install SSL: `y`

4. **Tunggu proses selesai** (~5 menit)

---

### STEP 5: Setup DNS Domain (3 menit)

1. **Login ke Freenom**: https://my.freenom.com
2. **Services** → **My Domains** → **Manage Domain**
3. **Management Tools** → **Nameservers**
4. **Use custom nameservers**:
   ```
   ns1.jetorbit.net
   ns2.jetorbit.net
   ```
5. **Change Nameservers** → **Confirm**

6. **Login ke Jetorbit Panel**
7. **DNS Management** → **Add Domain**: `spd-bankjatim.tk`
8. **Add A Records**:
   ```
   Type: A    | Name: @   | Value: [IP_VPS_ANDA] | TTL: 3600
   Type: A    | Name: www | Value: [IP_VPS_ANDA] | TTL: 3600
   ```
9. **Save**

---

### STEP 6: Verifikasi (2 menit)

1. **Tunggu propagasi DNS** (5-30 menit)
2. **Test akses**:
   - Via IP: http://IP_VPS_ANDA
   - Via Domain: http://spd-bankjatim.tk (tunggu DNS propagasi)
   - Via HTTPS: https://spd-bankjatim.tk (setelah SSL terinstall)

3. **Login Admin**:
   - URL: https://spd-bankjatim.tk/admin/login
   - Username: `admin`
   - Password: `admin123`
   - ⚠️ **SEGERA GANTI PASSWORD!**

4. **Test fitur**:
   - ✅ Input SPD
   - ✅ Download PDF
   - ✅ Admin dashboard
   - ✅ Monitoring iuran

---

## 🎉 SELESAI! APLIKASI SUDAH LIVE!

### Akses Aplikasi:
- **Frontend**: https://spd-bankjatim.tk
- **Admin Panel**: https://spd-bankjatim.tk/admin/login
- **Backend API**: https://spd-bankjatim.tk/api

---

## 🔧 MANAGEMENT COMMANDS

### Cek Status:
```bash
# SSH ke VPS
ssh root@IP_VPS_ANDA

# Cek status services
systemctl status spd-backend
systemctl status spd-frontend
systemctl status nginx
```

### Restart Services:
```bash
systemctl restart spd-backend
systemctl restart spd-frontend
systemctl restart nginx
```

### View Logs:
```bash
# Backend logs
tail -f /var/log/spd-backend.log

# Frontend logs
tail -f /var/log/spd-frontend.log

# Nginx logs
tail -f /var/log/nginx/access.log
tail -f /var/log/nginx/error.log
```

### Update Aplikasi:
```bash
cd /var/www/perjalanan-dinas

# Pull latest code (jika pakai Git)
git pull

# Rebuild backend
cd backend
go build -o spd-backend cmd/api/main.go
systemctl restart spd-backend

# Rebuild frontend
cd ../frontend
npm install
npm run build
systemctl restart spd-frontend
```

---

## 🆘 TROUBLESHOOTING

### Aplikasi tidak bisa diakses?
```bash
# Cek firewall
ufw status

# Cek Nginx
systemctl status nginx
nginx -t

# Cek backend
systemctl status spd-backend
journalctl -u spd-backend -n 50

# Cek frontend
systemctl status spd-frontend
journalctl -u spd-frontend -n 50
```

### Database error?
```bash
# Cek PostgreSQL
systemctl status postgresql

# Login ke database
sudo -u postgres psql
\l  # List databases
\q  # Quit
```

### Domain tidak bisa diakses?
1. Tunggu DNS propagasi (hingga 24 jam, biasanya 30 menit)
2. Cek DNS: https://dnschecker.org
3. Pastikan nameserver sudah diubah di Freenom
4. Pastikan A record sudah dibuat di Jetorbit DNS

---

## 📊 MONITORING

### Setup Monitoring Gratis:

1. **UptimeRobot** (https://uptimerobot.com)
   - Monitor uptime website
   - Alert via email jika down
   - GRATIS untuk 50 monitors

2. **Cloudflare** (https://cloudflare.com)
   - Analytics traffic
   - DDoS protection
   - CDN gratis

---

## 🔒 SECURITY CHECKLIST

- ✅ SSL/HTTPS aktif
- ✅ Firewall (UFW) configured
- ✅ Database password kuat
- ✅ JWT secret random
- ⚠️ **WAJIB**: Ganti admin password default!
- ⚠️ **WAJIB**: Setup backup database

---

## 💾 SETUP BACKUP OTOMATIS

```bash
# Buat script backup
nano /root/backup-db.sh
```

Isi file:
```bash
#!/bin/bash
BACKUP_DIR="/root/backups"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p ${BACKUP_DIR}

# Backup database
sudo -u postgres pg_dump perjalanan_dinas > ${BACKUP_DIR}/db_${DATE}.sql

# Hapus backup >7 hari
find ${BACKUP_DIR} -name "db_*.sql" -mtime +7 -delete

echo "Backup completed: db_${DATE}.sql"
```

```bash
# Make executable
chmod +x /root/backup-db.sh

# Setup cron (backup setiap hari jam 2 pagi)
crontab -e

# Tambahkan:
0 2 * * * /root/backup-db.sh >> /var/log/backup.log 2>&1
```

---

## 📞 SUPPORT

### Jetorbit Support:
- Live Chat: https://jetorbit.com
- Email: support@jetorbit.com
- WhatsApp: (Ada di website)

### Dokumentasi Lengkap:
- Baca: `DEPLOYMENT_JETORBIT.md`

---

## 🎯 UPGRADE PATH

### Jika Traffic Meningkat:

**VPS SSD 2** (Rp 200.000/bulan):
- 2 CPU Cores
- 2 GB RAM
- 50 GB SSD
- Lebih cepat, handle lebih banyak user

**Domain Berbayar** (Rp 150.000/tahun):
- `.co.id` - Lebih profesional
- Kredibilitas lebih tinggi
- Support lokal Indonesia

---

## ✅ CHECKLIST DEPLOYMENT

- [ ] VPS Jetorbit sudah dibeli
- [ ] Email VPS diterima (IP, username, password)
- [ ] Domain gratis dari Freenom sudah didapat
- [ ] Code sudah diupload ke VPS
- [ ] Script deploy.sh sudah dijalankan
- [ ] DNS nameserver sudah diubah ke Jetorbit
- [ ] A records sudah dibuat di Jetorbit
- [ ] SSL Let's Encrypt sudah terinstall
- [ ] Aplikasi bisa diakses via domain
- [ ] Admin login berhasil
- [ ] Password admin sudah diganti
- [ ] Backup database sudah disetup
- [ ] Monitoring sudah aktif

---

**🎊 SELAMAT! APLIKASI ANDA SUDAH LIVE DI INTERNET! 🎊**

Total waktu: **15 menit**
Total biaya: **Rp 100.000/bulan**
Domain: **GRATIS**
SSL: **GRATIS**

**Aplikasi production-ready untuk Bank Jatim Digital Banking Division!** 🏦✨
