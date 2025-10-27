#!/bin/bash

# 🚀 Script Deployment Otomatis untuk Jetorbit VPS
# Author: UncleJSS
# Description: Script untuk deploy aplikasi Perjalanan Dinas ke VPS Jetorbit

set -e  # Exit on error

echo "=========================================="
echo "🚀 DEPLOYMENT SCRIPT - SPD BANK JATIM"
echo "=========================================="
echo ""

# Warna untuk output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Fungsi helper
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ️  $1${NC}"
}

# Cek apakah running sebagai root
if [ "$EUID" -ne 0 ]; then
    print_error "Script ini harus dijalankan sebagai root!"
    echo "Gunakan: sudo ./deploy.sh"
    exit 1
fi

print_info "Script berjalan sebagai root"

# 1. Update System
echo ""
echo "📦 Step 1: Update System"
echo "----------------------------"
apt update && apt upgrade -y
print_success "System updated"

# 2. Install Dependencies
echo ""
echo "🔧 Step 2: Install Dependencies"
echo "----------------------------"

# PostgreSQL
if ! command -v psql &> /dev/null; then
    print_info "Installing PostgreSQL..."
    apt install postgresql postgresql-contrib -y
    systemctl enable postgresql
    systemctl start postgresql
    print_success "PostgreSQL installed"
else
    print_success "PostgreSQL already installed"
fi

# Go
if ! command -v go &> /dev/null; then
    print_info "Installing Go..."
    wget -q https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
    tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/bin/go/bin' >> ~/.bashrc
    export PATH=$PATH:/usr/local/bin/go/bin
    rm go1.21.5.linux-amd64.tar.gz
    print_success "Go installed"
else
    print_success "Go already installed"
fi

# Node.js
if ! command -v node &> /dev/null; then
    print_info "Installing Node.js..."
    curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
    apt install -y nodejs
    print_success "Node.js installed"
else
    print_success "Node.js already installed"
fi

# Nginx
if ! command -v nginx &> /dev/null; then
    print_info "Installing Nginx..."
    apt install nginx -y
    systemctl enable nginx
    systemctl start nginx
    print_success "Nginx installed"
else
    print_success "Nginx already installed"
fi

# Git
if ! command -v git &> /dev/null; then
    apt install git -y
    print_success "Git installed"
else
    print_success "Git already installed"
fi

# 3. Setup PostgreSQL
echo ""
echo "🗄️  Step 3: Setup PostgreSQL"
echo "----------------------------"

# Prompt for database credentials
read -p "Database name [perjalanan_dinas]: " DB_NAME
DB_NAME=${DB_NAME:-perjalanan_dinas}

read -p "Database user [spduser]: " DB_USER
DB_USER=${DB_USER:-spduser}

read -sp "Database password: " DB_PASSWORD
echo ""

if [ -z "$DB_PASSWORD" ]; then
    print_error "Password tidak boleh kosong!"
    exit 1
fi

# Create database and user
sudo -u postgres psql -c "CREATE DATABASE ${DB_NAME};" 2>/dev/null || print_info "Database already exists"
sudo -u postgres psql -c "CREATE USER ${DB_USER} WITH ENCRYPTED PASSWORD '${DB_PASSWORD}';" 2>/dev/null || print_info "User already exists"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE ${DB_NAME} TO ${DB_USER};"

print_success "Database setup completed"

# 4. Setup Application Directory
echo ""
echo "📁 Step 4: Setup Application"
echo "----------------------------"

APP_DIR="/var/www/perjalanan-dinas"
mkdir -p ${APP_DIR}

# Check if code exists
if [ ! -d "${APP_DIR}/backend" ]; then
    print_info "Code directory not found"
    read -p "Upload code ke ${APP_DIR} terlebih dahulu. Continue? (y/n): " continue
    if [ "$continue" != "y" ]; then
        exit 0
    fi
fi

cd ${APP_DIR}

# 5. Setup Backend
echo ""
echo "🔨 Step 5: Setup Backend"
echo "----------------------------"

cd ${APP_DIR}/backend

# Create .env file
cat > .env << EOF
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=${DB_USER}
DB_PASSWORD=${DB_PASSWORD}
DB_NAME=${DB_NAME}
JWT_SECRET=$(openssl rand -base64 32)
EOF

print_success ".env file created"

# Build backend
print_info "Building backend..."
/usr/local/go/bin/go mod download
/usr/local/go/bin/go build -o spd-backend cmd/api/main.go
chmod +x spd-backend

print_success "Backend built successfully"

# 6. Setup Frontend
echo ""
echo "🎨 Step 6: Setup Frontend"
echo "----------------------------"

cd ${APP_DIR}/frontend

# Get server IP
SERVER_IP=$(curl -s ifconfig.me)

# Create .env.local
cat > .env.local << EOF
NEXT_PUBLIC_API_URL=http://${SERVER_IP}:8080
EOF

print_success ".env.local created"

# Install and build
print_info "Installing frontend dependencies..."
npm install --silent

print_info "Building frontend..."
npm run build

print_success "Frontend built successfully"

# 7. Create Systemd Services
echo ""
echo "⚙️  Step 7: Create Systemd Services"
echo "----------------------------"

# Backend service
cat > /etc/systemd/system/spd-backend.service << EOF
[Unit]
Description=SPD Backend Service
After=network.target postgresql.service

[Service]
Type=simple
User=root
WorkingDirectory=${APP_DIR}/backend
ExecStart=${APP_DIR}/backend/spd-backend
Restart=always
RestartSec=5
StandardOutput=append:/var/log/spd-backend.log
StandardError=append:/var/log/spd-backend-error.log

EnvironmentFile=${APP_DIR}/backend/.env

[Install]
WantedBy=multi-user.target
EOF

# Frontend service
cat > /etc/systemd/system/spd-frontend.service << EOF
[Unit]
Description=SPD Frontend Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=${APP_DIR}/frontend
ExecStart=/usr/bin/npm start
Restart=always
RestartSec=5
StandardOutput=append:/var/log/spd-frontend.log
StandardError=append:/var/log/spd-frontend-error.log

Environment=NODE_ENV=production
Environment=PORT=3000

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable spd-backend spd-frontend
systemctl restart spd-backend spd-frontend

print_success "Services created and started"

# 8. Setup Nginx
echo ""
echo "🌐 Step 8: Setup Nginx"
echo "----------------------------"

read -p "Domain name (optional, tekan Enter untuk skip): " DOMAIN

if [ -z "$DOMAIN" ]; then
    # Tanpa domain, gunakan IP
    cat > /etc/nginx/sites-available/spd << EOF
server {
    listen 80 default_server;
    server_name ${SERVER_IP};

    location /api {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_cache_bypass \$http_upgrade;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
    }

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_cache_bypass \$http_upgrade;
    }
}
EOF
else
    # Dengan domain
    cat > /etc/nginx/sites-available/spd << EOF
server {
    listen 80;
    server_name ${DOMAIN} www.${DOMAIN};

    location /api {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_cache_bypass \$http_upgrade;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
    }

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_cache_bypass \$http_upgrade;
    }
}
EOF
fi

# Enable site
rm -f /etc/nginx/sites-enabled/default
ln -sf /etc/nginx/sites-available/spd /etc/nginx/sites-enabled/

# Test and reload nginx
nginx -t
systemctl reload nginx

print_success "Nginx configured"

# 9. Setup Firewall
echo ""
echo "🔒 Step 9: Setup Firewall"
echo "----------------------------"

if ! command -v ufw &> /dev/null; then
    apt install ufw -y
fi

ufw --force reset
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow 'Nginx Full'
ufw --force enable

print_success "Firewall configured"

# 10. Setup SSL (Optional)
echo ""
echo "🔐 Step 10: Setup SSL"
echo "----------------------------"

if [ ! -z "$DOMAIN" ]; then
    read -p "Install SSL dengan Let's Encrypt? (y/n): " install_ssl
    if [ "$install_ssl" == "y" ]; then
        apt install certbot python3-certbot-nginx -y
        certbot --nginx -d ${DOMAIN} -d www.${DOMAIN} --non-interactive --agree-tos --register-unsafely-without-email --redirect
        print_success "SSL installed"
    fi
else
    print_info "SSL tidak diinstall (butuh domain)"
fi

# Final Status Check
echo ""
echo "=========================================="
echo "✅ DEPLOYMENT SELESAI!"
echo "=========================================="
echo ""
echo "📊 Status Services:"
systemctl status spd-backend --no-pager -l | head -3
systemctl status spd-frontend --no-pager -l | head -3
systemctl status nginx --no-pager -l | head -3
echo ""
echo "🌐 Akses Aplikasi:"
if [ -z "$DOMAIN" ]; then
    echo "   Frontend: http://${SERVER_IP}"
    echo "   Backend:  http://${SERVER_IP}/api"
else
    if [ "$install_ssl" == "y" ]; then
        echo "   Frontend: https://${DOMAIN}"
        echo "   Backend:  https://${DOMAIN}/api"
    else
        echo "   Frontend: http://${DOMAIN}"
        echo "   Backend:  http://${DOMAIN}/api"
    fi
fi
echo ""
echo "📝 Default Admin Login:"
echo "   Username: admin"
echo "   Password: admin123"
echo "   ⚠️  PENTING: Segera ganti password setelah login!"
echo ""
echo "📋 Useful Commands:"
echo "   Restart Backend:  systemctl restart spd-backend"
echo "   Restart Frontend: systemctl restart spd-frontend"
echo "   View Logs Backend:  tail -f /var/log/spd-backend.log"
echo "   View Logs Frontend: tail -f /var/log/spd-frontend.log"
echo ""
echo "=========================================="
