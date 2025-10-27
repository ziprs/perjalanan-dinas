#!/bin/bash

# ============================================
# 🔄 AUTO UPDATE SCRIPT - SPD BANK JATIM
# ============================================

set -e  # Exit on error

echo "============================================"
echo "🔄 STARTING UPDATE PROCESS"
echo "============================================"
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to print colored messages
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ️  $1${NC}"
}

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    print_error "Please run as root (use sudo or login as root)"
    exit 1
fi

# Go to project directory
cd /var/www/perjalanan-dinas || {
    print_error "Project directory not found!"
    exit 1
}

print_success "In project directory: $(pwd)"
echo ""

# Step 1: Backup current version
echo "📦 Step 1: Creating backup..."
BACKUP_DIR="/tmp/spd_backup_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$BACKUP_DIR"
cp -r backend frontend "$BACKUP_DIR/" 2>/dev/null || true
print_success "Backup created at: $BACKUP_DIR"
echo ""

# Step 2: Get current commit hash
CURRENT_COMMIT=$(git rev-parse --short HEAD)
print_info "Current version: $CURRENT_COMMIT"
echo ""

# Step 3: Pull latest code
echo "📥 Step 2: Pulling latest code from GitHub..."
git fetch origin main

# Check if there are updates
UPDATES=$(git log HEAD..origin/main --oneline)
if [ -z "$UPDATES" ]; then
    print_warning "No updates available. Already up to date!"
    exit 0
fi

echo "New commits to pull:"
git log HEAD..origin/main --oneline --color=always
echo ""

# Pull the updates
git pull origin main || {
    print_error "Git pull failed! Rolling back..."
    git reset --hard "$CURRENT_COMMIT"
    exit 1
}

NEW_COMMIT=$(git rev-parse --short HEAD)
print_success "Updated from $CURRENT_COMMIT to $NEW_COMMIT"
echo ""

# Step 4: Check what changed
echo "🔍 Step 3: Analyzing changes..."
BACKEND_CHANGED=false
FRONTEND_CHANGED=false

if git diff --name-only "$CURRENT_COMMIT" "$NEW_COMMIT" | grep -q "^backend/"; then
    BACKEND_CHANGED=true
    print_info "Backend files changed - will rebuild backend"
fi

if git diff --name-only "$CURRENT_COMMIT" "$NEW_COMMIT" | grep -q "^frontend/"; then
    FRONTEND_CHANGED=true
    print_info "Frontend files changed - will rebuild frontend"
fi

echo ""

# Step 5: Update Backend if needed
if [ "$BACKEND_CHANGED" = true ]; then
    echo "🔨 Step 4a: Updating Backend..."
    
    cd backend
    
    # Check if go.mod changed (need to update dependencies)
    if git diff --name-only "$CURRENT_COMMIT" "$NEW_COMMIT" | grep -q "backend/go.mod"; then
        print_info "Dependencies changed - running go mod tidy..."
        go mod tidy
    fi
    
    # Build backend
    print_info "Building backend..."
    go build -o spd-backend cmd/api/main.go || {
        print_error "Backend build failed!"
        cd ..
        git reset --hard "$CURRENT_COMMIT"
        exit 1
    }
    
    # Restart backend service
    print_info "Restarting backend service..."
    systemctl restart spd-backend
    sleep 2
    
    # Check backend status
    if systemctl is-active --quiet spd-backend; then
        print_success "Backend updated and running"
    else
        print_error "Backend service failed to start!"
        systemctl status spd-backend --no-pager
        exit 1
    fi
    
    cd ..
    echo ""
else
    print_info "Backend unchanged - skipping backend update"
    echo ""
fi

# Step 6: Update Frontend if needed
if [ "$FRONTEND_CHANGED" = true ]; then
    echo "🔨 Step 4b: Updating Frontend..."
    
    cd frontend
    
    # Check if package.json changed (need to update dependencies)
    if git diff --name-only "$CURRENT_COMMIT" "$NEW_COMMIT" | grep -q "frontend/package.json"; then
        print_info "Dependencies changed - running npm install..."
        npm install
    fi
    
    # Clear Next.js cache
    print_info "Clearing Next.js cache..."
    rm -rf .next
    
    # Build frontend
    print_info "Building frontend..."
    npm run build || {
        print_error "Frontend build failed!"
        cd ..
        git reset --hard "$CURRENT_COMMIT"
        exit 1
    }
    
    # Restart frontend service
    print_info "Restarting frontend service..."
    systemctl restart spd-frontend
    sleep 2
    
    # Check frontend status
    if systemctl is-active --quiet spd-frontend; then
        print_success "Frontend updated and running"
    else
        print_error "Frontend service failed to start!"
        systemctl status spd-frontend --no-pager
        exit 1
    fi
    
    cd ..
    echo ""
else
    print_info "Frontend unchanged - skipping frontend update"
    echo ""
fi

# Step 7: Restart Nginx (just in case)
echo "🔄 Step 5: Restarting Nginx..."
systemctl restart nginx
sleep 1

if systemctl is-active --quiet nginx; then
    print_success "Nginx restarted"
else
    print_error "Nginx failed to restart!"
    systemctl status nginx --no-pager
    exit 1
fi

echo ""

# Step 8: Verify all services
echo "✅ Step 6: Verifying all services..."
ALL_OK=true

if systemctl is-active --quiet spd-backend; then
    print_success "Backend: Running"
else
    print_error "Backend: Not Running"
    ALL_OK=false
fi

if systemctl is-active --quiet spd-frontend; then
    print_success "Frontend: Running"
else
    print_error "Frontend: Not Running"
    ALL_OK=false
fi

if systemctl is-active --quiet nginx; then
    print_success "Nginx: Running"
else
    print_error "Nginx: Not Running"
    ALL_OK=false
fi

echo ""

# Step 9: Test endpoints
echo "🧪 Step 7: Testing endpoints..."

# Test main page
if curl -s -o /dev/null -w "%{http_code}" http://localhost/spd | grep -q "200"; then
    print_success "Main page: OK"
else
    print_warning "Main page: Failed"
    ALL_OK=false
fi

# Test admin login page
if curl -s -o /dev/null -w "%{http_code}" http://localhost/spd/admin/login | grep -q "200"; then
    print_success "Admin login: OK"
else
    print_warning "Admin login: Failed"
    ALL_OK=false
fi

# Test backend API
if curl -s -o /dev/null -w "%{http_code}" http://localhost/spd/api/positions | grep -q "200"; then
    print_success "Backend API: OK"
else
    print_warning "Backend API: Failed"
    ALL_OK=false
fi

echo ""
echo "============================================"

if [ "$ALL_OK" = true ]; then
    print_success "🎉 UPDATE COMPLETED SUCCESSFULLY!"
    echo ""
    echo "Updated from: $CURRENT_COMMIT"
    echo "         to: $NEW_COMMIT"
    echo ""
    echo "Backup saved at: $BACKUP_DIR"
    echo ""
    echo "Application URL: http://103.235.75.196/spd"
    echo "Admin Panel: http://103.235.75.196/spd/admin/login"
else
    print_error "⚠️  UPDATE COMPLETED WITH WARNINGS"
    echo ""
    echo "Some services or endpoints are not responding correctly."
    echo "Please check the logs:"
    echo "  - Backend: tail -50 /var/log/spd-backend-error.log"
    echo "  - Frontend: tail -50 /var/log/spd-frontend.log"
    echo "  - Nginx: tail -50 /var/log/nginx/error.log"
fi

echo "============================================"
