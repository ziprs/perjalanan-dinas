# Socket MCP Integration Guide

## Apa itu Socket MCP?

Socket MCP (Model Context Protocol) memungkinkan Claude Desktop untuk berinteraksi langsung dengan Socket Security API, memberikan real-time security insights saat development.

## Setup Socket MCP di Claude Desktop

### 1. Install Claude Desktop

Download dari: https://claude.ai/download

### 2. Add Socket MCP Server

Buka terminal dan jalankan:

```bash
claude mcp add --transport http socket https://mcp.socket.dev/
```

Atau edit manual config file:

**macOS/Linux:**
```bash
nano ~/Library/Application\ Support/Claude/claude_desktop_config.json
```

**Windows:**
```powershell
notepad %APPDATA%\Claude\claude_desktop_config.json
```

Tambahkan:
```json
{
  "mcpServers": {
    "socket": {
      "url": "https://mcp.socket.dev/",
      "transport": "http"
    }
  }
}
```

### 3. Restart Claude Desktop

Setelah config disimpan, restart Claude Desktop.

## Menggunakan Socket MCP

### Di Claude Desktop, Anda bisa bertanya:

**Contoh 1 - Scan Dependencies:**
```
Check my project dependencies for security issues
Project path: /Users/unclejss/Documents/Project SPD/perjalanan-dinas/frontend
```

**Contoh 2 - Analyze Package:**
```
Analyze security of package: react-hook-form
```

**Contoh 3 - Get Recommendations:**
```
What are safer alternatives for lodash?
```

**Contoh 4 - Check for CVEs:**
```
Check for known vulnerabilities in my package.json
```

## Local Socket CLI Usage

### Install Socket CLI

```bash
cd /Users/unclejss/Documents/Project\ SPD/perjalanan-dinas
./scripts/setup-socket-security.sh
```

### Get API Token

1. Sign up: https://socket.dev (free for open source)
2. Get token: https://socket.dev/settings/api
3. Set environment variable:

```bash
# Add to ~/.zshrc or ~/.bashrc
export SOCKET_SECURITY_API_KEY="your_api_key_here"
```

### Run Security Scans

```bash
# Scan all dependencies
npm run security:scan

# Scan frontend only
npm run security:frontend

# Generate detailed report
cd frontend
socket report

# Watch mode (continuous monitoring)
socket ci --watch
```

## Integration with VS Code

### Install Socket Extension

1. Open VS Code
2. Go to Extensions (Cmd+Shift+X)
3. Search "Socket Security"
4. Install extension

### Features:
- ✅ Real-time vulnerability warnings
- ✅ Inline security suggestions
- ✅ Package safety scores
- ✅ Auto-scan on npm install

## Pre-commit Hook (Optional)

Tambahkan security check sebelum commit:

```bash
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash

echo "Running Socket Security scan..."
cd frontend && npx socket-security ci

if [ $? -ne 0 ]; then
    echo "❌ Security issues found! Commit blocked."
    echo "Run 'npm run security:frontend' for details"
    exit 1
fi

echo "✅ Security check passed"
exit 0
EOF

chmod +x .git/hooks/pre-commit
```

## What Socket Checks:

### 🔴 Critical Issues:
- Malware detection
- Install scripts (potential backdoors)
- Typosquatting packages
- Known CVEs

### 🟡 Medium Issues:
- Deprecated packages
- Unmaintained dependencies
- License compliance
- Minified code

### 🟢 Best Practices:
- Package popularity
- Maintainer reputation
- Update frequency
- Community trust

## Troubleshooting

### Socket CLI not found
```bash
npm install -g @socketsecurity/cli
```

### API Key not working
```bash
# Verify key is set
echo $SOCKET_SECURITY_API_KEY

# Re-export if needed
export SOCKET_SECURITY_API_KEY="your_key"
```

### MCP not connecting
1. Check Claude Desktop version (latest)
2. Verify config syntax in JSON
3. Check network connection
4. Restart Claude Desktop

## Resources

- Socket Documentation: https://docs.socket.dev
- Socket MCP Docs: https://github.com/socketsecurity/mcp-server
- Claude MCP Guide: https://docs.anthropic.com/claude/docs/mcp
- Socket Dashboard: https://socket.dev/dashboard
