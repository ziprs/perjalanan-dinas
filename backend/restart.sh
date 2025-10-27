#!/bin/bash

echo "======================================"
echo "  Restarting Backend Server"
echo "======================================"
echo ""

# Stop backend server on port 8080
echo "1. Stopping backend server on port 8080..."
lsof -ti:8080 | xargs kill -9 2>/dev/null

if [ $? -eq 0 ]; then
    echo "   ✓ Server stopped successfully"
else
    echo "   ℹ No server running on port 8080"
fi

echo ""
echo "2. Starting backend server..."
echo ""

# Start backend server
cd /Users/unclejss/Documents/Project\ SPD/perjalanan-dinas/backend
go run cmd/api/main.go
