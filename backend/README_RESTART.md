# Panduan Restart Backend Server

## Metode 1: Menggunakan Script Restart (Paling Mudah)

```bash
cd /Users/unclejss/Documents/Project\ SPD/perjalanan-dinas/backend
./restart.sh
```

Script ini akan:
1. Stop server yang sedang berjalan di port 8080
2. Start server kembali

## Metode 2: Menggunakan Air (Auto-Reload - Recommended)

### Install Air (Hanya sekali)
```bash
go install github.com/air-verse/air@latest
```

### Jalankan dengan Air
```bash
cd /Users/unclejss/Documents/Project\ SPD/perjalanan-dinas/backend
/Users/unclejss/go/bin/air
```

**Keuntungan menggunakan Air:**
- Otomatis restart setiap kali ada perubahan file `.go`
- Tidak perlu restart manual lagi
- Seperti hot-reload di Next.js

### Stop Air
Tekan `CTRL + C` di terminal

## Metode 3: Manual (Basic)

### Stop Server
```bash
lsof -ti:8080 | xargs kill -9
```

### Start Server
```bash
cd /Users/unclejss/Documents/Project\ SPD/perjalanan-dinas/backend
go run cmd/api/main.go
```

## Tips

- **Untuk development sehari-hari**: Gunakan Air (Metode 2)
- **Untuk restart cepat**: Gunakan script (Metode 1)
- **Jika ada masalah**: Gunakan manual (Metode 3)

## Troubleshooting

### Error: "permission denied" saat menjalankan restart.sh
```bash
chmod +x restart.sh
```

### Error: "air: command not found"
Gunakan full path:
```bash
/Users/unclejss/go/bin/air
```

Atau tambahkan Go bin ke PATH (opsional):
```bash
export PATH=$PATH:$(go env GOPATH)/bin
# Tambahkan ke ~/.zshrc atau ~/.bash_profile agar permanen
```

### Port 8080 masih terpakai
```bash
lsof -ti:8080 | xargs kill -9
```
