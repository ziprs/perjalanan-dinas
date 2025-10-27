#!/bin/bash

echo "================================="
echo "Socket Security Setup"
echo "================================="

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if npm is installed
if ! command -v npm &> /dev/null; then
    echo "Error: npm is not installed"
    exit 1
fi

echo -e "${GREEN}Installing Socket Security CLI...${NC}"
npm install -g @socketsecurity/cli

echo ""
echo -e "${GREEN}Verifying installation...${NC}"
socket --version

echo ""
echo -e "${YELLOW}To use Socket Security:${NC}"
echo "1. Sign up at https://socket.dev (free for open source)"
echo "2. Get your API token from https://socket.dev/settings/api"
echo "3. Set environment variable:"
echo "   export SOCKET_SECURITY_API_KEY=your_api_key_here"
echo ""
echo -e "${GREEN}Available commands:${NC}"
echo "  npm run security:scan       - Scan all dependencies"
echo "  npm run security:frontend   - Scan frontend only"
echo "  socket ci                   - Run in CI mode"
echo "  socket report               - Generate security report"
echo ""
echo -e "${GREEN}Setup complete!${NC}"
