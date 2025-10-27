# 🚀 PANDUAN DEPLOYMENT KE JETORBIT

## Persiapan

### 1. Beli Hosting Jetorbit
1. Buka https://jetorbit.com
2. Pilih **VPS Hosting** (Recommended)
   - Paket: VPS SSD 1 (Rp 100.000/bulan)
   - Spesifikasi: 1 CPU, 1GB RAM, 25GB SSD
   - OS: Ubuntu 22.04 LTS
3. Selesaikan pembayaran

### 2. Domain Gratis

#### Opsi A: Freenom (Gratis 1 Tahun)
1. Buka https://freenom.com
2. Cari domain yang diinginkan (contoh: `spd-bankjatim`)
3. Pilih ekstensi: `.tk`, `.ml`, `.ga`, `.cf`, atau `.gq`
4. Checkout dengan durasi 12 bulan (GRATIS)
5. Register/Login
6. Selesai! Domain gratis Anda: `spd-bankjatim.tk`

#### Opsi B: Subdomain Jetorbit (Gratis)
- Jetorbit memberikan subdomain gratis saat beli VPS
- Format: `namaanda.jetorbit.net`
- Sudah include SSL gratis

#### Opsi C: Cloudflare Pages (Gratis untuk Frontend)
1. Deploy frontend ke Cloudflare Pages (GRATIS)
2. Dapatkan domain: `spd-bankjatim.pages.dev`
3. Backend tetap di Jetorbit

---

## 🔧 SETUP VPS JETORBIT

### 1. Koneksi ke VPS
```bash
# Gunakan SSH client (Terminal/PuTTY)
ssh root@IP_VPS_ANDA

# Password dikirim ke email oleh Jetorbit
```

### 2. Install Dependencies
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install PostgreSQL
sudo apt install postgresql postgresql-contrib -y

# Install Go (Backend)
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/bin/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install Node.js (Frontend build)
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt install -y nodejs

# Install Nginx (Reverse Proxy)
sudo apt install nginx -y

# Install Git
sudo apt install git -y
```

### 3. Setup PostgreSQL
```bash
# Login ke PostgreSQL
sudo -u postgres psql

# Buat database dan user
CREATE DATABASE perjalanan_dinas;
CREATE USER spduser WITH ENCRYPTED PASSWORD 'password_kuat_anda';
GRANT ALL PRIVILEGES ON DATABASE perjalanan_dinas TO spduser;
\q

# Allow remote connection (opsional)
sudo nano /etc/postgresql/14/main/pg_hba.conf
# Tambahkan: host all all 0.0.0.0/0 md5

sudo systemctl restart postgresql
```

---

## 📦 DEPLOY APLIKASI

### 1. Upload Code ke VPS

#### Cara A: Via Git (Recommended)
```bash
# Di VPS, buat folder
mkdir -p /var/www/perjalanan-dinas
cd /var/www/perjalanan-dinas

# Clone dari repository Anda
git clone https://github.com/username/perjalanan-dinas.git .

# Atau upload manual dari komputer lokal
```

#### Cara B: Via SCP/SFTP
```bash
# Dari komputer lokal
scp -r /Users/unclejss/Documents/Project\ SPD/perjalanan-dinas root@IP_VPS:/var/www/
```

### 2. Setup Backend (Go)
```bash
cd /var/www/perjalanan-dinas/backend

# Copy environment file
cp .env.example .env

# Edit .env
nano .env
```

**Isi .env:**
```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=spduser
DB_PASSWORD=password_kuat_anda
DB_NAME=perjalanan_dinas
JWT_SECRET=secret_key_yang_sangat_kuat_dan_random_123456
```

```bash
# Install dependencies
go mod download

# Build aplikasi
go build -o spd-backend cmd/api/main.go

# Test run
./spd-backend
```

### 3. Setup Frontend (Next.js)
```bash
cd /var/www/perjalanan-dinas/frontend

# Install dependencies
npm install

# Edit .env.local
nano .env.local
```

**Isi .env.local:**
```env
NEXT_PUBLIC_API_URL=http://IP_VPS_ANDA:8080
# Atau jika sudah setup domain:
# NEXT_PUBLIC_API_URL=https://api.spd-bankjatim.tk
```

```bash
# Build production
npm run build

# Test production server
npm start
```

---

## 🔄 SETUP SYSTEMD SERVICE (Auto-Restart)

### 1. Backend Service
```bash
sudo nano /etc/systemd/system/spd-backend.service
```

**Isi file:**
```ini
[Unit]
Description=SPD Backend Service
After=network.target postgresql.service

[Service]
Type=simple
User=root
WorkingDirectory=/var/www/perjalanan-dinas/backend
ExecStart=/var/www/perjalanan-dinas/backend/spd-backend
Restart=always
RestartSec=5
StandardOutput=append:/var/log/spd-backend.log
StandardError=append:/var/log/spd-backend-error.log

Environment=PORT=8080
EnvironmentFile=/var/www/perjalanan-dinas/backend/.env

[Install]
WantedBy=multi-user.target
```

```bash
# Enable dan start service
sudo systemctl daemon-reload
sudo systemctl enable spd-backend
sudo systemctl start spd-backend
sudo systemctl status spd-backend
```

### 2. Frontend Service
```bash
sudo nano /etc/systemd/system/spd-frontend.service
```

**Isi file:**
```ini
[Unit]
Description=SPD Frontend Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/var/www/perjalanan-dinas/frontend
ExecStart=/usr/bin/npm start
Restart=always
RestartSec=5
StandardOutput=append:/var/log/spd-frontend.log
StandardError=append:/var/log/spd-frontend-error.log

Environment=NODE_ENV=production
Environment=PORT=3000

[Install]
WantedBy=multi-user.target
```

```bash
# Enable dan start service
sudo systemctl daemon-reload
sudo systemctl enable spd-frontend
sudo systemctl start spd-frontend
sudo systemctl status spd-frontend
```

---

## 🌐 SETUP NGINX (Reverse Proxy)

```bash
sudo nano /etc/nginx/sites-available/spd
```

**Isi file:**
```nginx
# Backend API
server {
    listen 80;
    server_name api.spd-bankjatim.tk;  # Ganti dengan domain Anda

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# Frontend
server {
    listen 80;
    server_name spd-bankjatim.tk www.spd-bankjatim.tk;  # Ganti dengan domain Anda

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

```bash
# Enable site
sudo ln -s /etc/nginx/sites-available/spd /etc/nginx/sites-enabled/

# Test konfigurasi
sudo nginx -t

# Restart nginx
sudo systemctl restart nginx
```

---

## 🔒 SETUP SSL GRATIS (Let's Encrypt)

```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx -y

# Generate SSL certificate
sudo certbot --nginx -d spd-bankjatim.tk -d www.spd-bankjatim.tk -d api.spd-bankjatim.tk

# Email untuk notifikasi
# Pilih: Agree to terms
# Pilih: Yes untuk redirect HTTP ke HTTPS

# Auto-renewal sudah aktif otomatis!
# Test renewal:
sudo certbot renew --dry-run
```

---

## 🔧 SETUP DOMAIN

### Jika Pakai Freenom:
1. Login ke Freenom
2. Go to **Services → My Domains**
3. Click **Manage Domain**
4. Pilih **Management Tools → Nameservers**
5. Pilih **Use custom nameservers**
6. Masukkan nameserver Jetorbit:
   ```
   ns1.jetorbit.net
   ns2.jetorbit.net
   ```

### Setup DNS di Jetorbit:
1. Login ke panel Jetorbit
2. Go to **DNS Management**
3. Tambahkan A Records:
   ```
   Type: A
   Name: @
   Value: IP_VPS_ANDA
   TTL: 3600

   Type: A
   Name: www
   Value: IP_VPS_ANDA
   TTL: 3600

   Type: A
   Name: api
   Value: IP_VPS_ANDA
   TTL: 3600
   ```

---

## ✅ VERIFIKASI DEPLOYMENT

### 1. Cek Services
```bash
# Cek backend
sudo systemctl status spd-backend
curl http://localhost:8080/api

# Cek frontend
sudo systemctl status spd-frontend
curl http://localhost:3000

# Cek nginx
sudo systemctl status nginx
```

### 2. Cek Logs
```bash
# Backend logs
tail -f /var/log/spd-backend.log
tail -f /var/log/spd-backend-error.log

# Frontend logs
tail -f /var/log/spd-frontend.log

# Nginx logs
tail -f /var/nginx/access.log
tail -f /var/nginx/error.log
```

### 3. Test dari Browser
1. Buka: https://spd-bankjatim.tk
2. Buka: https://api.spd-bankjatim.tk/api
3. Login admin dan test fitur

---

## 🔄 UPDATE APLIKASI

```bash
# Pull latest code
cd /var/www/perjalanan-dinas
git pull

# Update backend
cd backend
go build -o spd-backend cmd/api/main.go
sudo systemctl restart spd-backend

# Update frontend
cd ../frontend
npm install
npm run build
sudo systemctl restart spd-frontend
```

---

## 📊 MONITORING

### Setup Basic Monitoring
```bash
# Install htop
sudo apt install htop -y

# Monitor processes
htop

# Check disk space
df -h

# Check memory
free -h

# Check database connections
sudo -u postgres psql -c "SELECT count(*) FROM pg_stat_activity;"
```

---

## 🆘 TROUBLESHOOTING

### Backend tidak jalan:
```bash
sudo systemctl status spd-backend
journalctl -u spd-backend -n 50
```

### Frontend tidak jalan:
```bash
sudo systemctl status spd-frontend
journalctl -u spd-frontend -n 50
```

### Database error:
```bash
sudo systemctl status postgresql
sudo tail -f /var/log/postgresql/postgresql-14-main.log
```

### Nginx error:
```bash
sudo nginx -t
sudo tail -f /var/log/nginx/error.log
```

---

## 💰 ESTIMASI BIAYA

- VPS Jetorbit: **Rp 100.000/bulan**
- Domain Freenom: **GRATIS** (1 tahun)
- SSL Let's Encrypt: **GRATIS**
- **TOTAL: Rp 100.000/bulan = Rp 1.2 juta/tahun**

---

## 🎯 NEXT STEPS

1. ✅ Setup backup otomatis database
2. ✅ Setup monitoring (Uptime Robot gratis)
3. ✅ Setup firewall (UFW)
4. ✅ Hardening security

Aplikasi Anda sudah LIVE! 🎉
