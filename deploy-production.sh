#!/bin/bash

# Deploy Script untuk Production VPS
# Jalankan script ini di VPS dengan perintah: bash deploy-production.sh

set -e  # Exit on error

echo "========================================="
echo "  DEPLOYING PERJALANAN DINAS TO PRODUCTION"
echo "========================================="
echo ""

# Navigate to project directory
cd /home/admin/perjalanan-dinas || { echo "Error: Project directory not found"; exit 1; }

echo "=== Step 1: Pulling latest code from Git ==="
git pull origin main
echo "✓ Code updated successfully"
echo ""

# Backend deployment
echo "=== Step 2: Deploying Backend ==="
cd backend

echo "  - Installing Go dependencies..."
go mod tidy

echo "  - Building backend application..."
go build -o main cmd/api/main.go

echo "✓ Backend built successfully"
echo ""

# Frontend deployment
echo "=== Step 3: Deploying Frontend ==="
cd ../frontend

echo "  - Installing npm dependencies..."
npm install

echo "  - Building frontend for production..."
npm run build

echo "✓ Frontend built successfully"
echo ""

# Restart services
echo "=== Step 4: Restarting Services ==="

echo "  - Stopping backend service..."
sudo systemctl stop spd-backend || echo "Backend service was not running"

echo "  - Starting backend service..."
sudo systemctl start spd-backend

echo "  - Stopping frontend service..."
sudo systemctl stop spd-frontend || echo "Frontend service was not running"

echo "  - Starting frontend service..."
sudo systemctl start spd-frontend

echo "✓ Services restarted successfully"
echo ""

# Check service status
echo "=== Step 5: Checking Service Status ==="
echo ""
echo "Backend Service:"
sudo systemctl status spd-backend --no-pager | head -10
echo ""
echo "Frontend Service:"
sudo systemctl status spd-frontend --no-pager | head -10
echo ""

echo "========================================="
echo "  DEPLOYMENT COMPLETED SUCCESSFULLY!"
echo "========================================="
echo ""
echo "Application is now running at:"
echo "  - Frontend: http://103.160.37.195:3000"
echo "  - Backend:  http://103.160.37.195:8080"
echo ""
