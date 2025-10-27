#!/bin/bash

# VPS Setup Script for Perjalanan Dinas Application
# This script should be run as root on a fresh Ubuntu 22.04 VPS

set -e

echo "================================="
echo "VPS Setup for Perjalanan Dinas"
echo "================================="

# Update system
echo "Updating system packages..."
apt update && apt upgrade -y

# Install basic utilities
echo "Installing utilities..."
apt install -y curl wget git ufw fail2ban

# Setup firewall
echo "Configuring firewall..."
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow http
ufw allow https
ufw --force enable

# Install Go
echo "Installing Go..."
GO_VERSION="1.21.0"
wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
rm -rf /usr/local/go
tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
rm go${GO_VERSION}.linux-amd64.tar.gz

# Add Go to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
source /etc/profile

# Install PostgreSQL
echo "Installing PostgreSQL..."
apt install -y postgresql postgresql-contrib

# Setup PostgreSQL
echo "Configuring PostgreSQL..."
sudo -u postgres psql << EOF
CREATE DATABASE perjalanan_dinas;
CREATE USER perjalanan_user WITH ENCRYPTED PASSWORD 'CHANGE_THIS_PASSWORD';
GRANT ALL PRIVILEGES ON DATABASE perjalanan_dinas TO perjalanan_user;
\q
EOF

# Install Nginx
echo "Installing Nginx..."
apt install -y nginx

# Install Certbot for SSL
echo "Installing Certbot..."
apt install -y certbot python3-certbot-nginx

# Create application directory
echo "Creating application directory..."
mkdir -p /var/www/perjalanan-dinas/backend
chown -R $USER:$USER /var/www/perjalanan-dinas

# Create systemd service
echo "Creating systemd service..."
cat > /etc/systemd/system/perjalanan-dinas.service << 'EOF'
[Unit]
Description=Perjalanan Dinas API
After=network.target postgresql.service

[Service]
Type=simple
User=root
WorkingDirectory=/var/www/perjalanan-dinas/backend
ExecStart=/var/www/perjalanan-dinas/backend/app
Restart=on-failure
RestartSec=5s

# Environment variables
Environment="DB_HOST=localhost"
Environment="DB_PORT=5432"
Environment="DB_USER=perjalanan_user"
Environment="DB_PASSWORD=CHANGE_THIS_PASSWORD"
Environment="DB_NAME=perjalanan_dinas"
Environment="JWT_SECRET=CHANGE_THIS_SECRET"
Environment="PORT=8080"

[Install]
WantedBy=multi-user.target
EOF

# Create Nginx configuration
echo "Creating Nginx configuration..."
cat > /etc/nginx/sites-available/perjalanan-dinas << 'EOF'
server {
    listen 80;
    server_name api.yourdomain.com;

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

    location /health {
        proxy_pass http://localhost:8080/health;
        access_log off;
    }
}
EOF

# Enable Nginx site
ln -sf /etc/nginx/sites-available/perjalanan-dinas /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default

# Test Nginx configuration
nginx -t

# Reload systemd and start services
systemctl daemon-reload
systemctl enable perjalanan-dinas
systemctl restart nginx

echo ""
echo "================================="
echo "Setup Complete!"
echo "================================="
echo ""
echo "Next steps:"
echo "1. Update /etc/systemd/system/perjalanan-dinas.service with your actual credentials"
echo "2. Update /etc/nginx/sites-available/perjalanan-dinas with your actual domain"
echo "3. Setup SSL certificate: certbot --nginx -d api.yourdomain.com"
echo "4. Deploy your application files to /var/www/perjalanan-dinas/backend/"
echo "5. Start the service: systemctl start perjalanan-dinas"
echo ""
echo "Check service status: systemctl status perjalanan-dinas"
echo "Check logs: journalctl -u perjalanan-dinas -f"
echo ""
