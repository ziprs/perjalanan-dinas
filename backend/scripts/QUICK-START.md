# Quick Start - Database Cleanup (SPD Sequence = 324)

## ⚡ Eksekusi Cepat (Untuk yang Sudah Familiar)

### 1. SSH ke Server
```bash
ssh admin@103.160.37.195
cd /path/to/perjalanan-dinas/backend/scripts
```

### 2. Backup Database
```bash
sudo ./backup-full-database.sh
```

### 3. Stop Backend
```bash
sudo systemctl stop perjalanan-dinas-backend
```

### 4. Run Cleanup
```bash
sudo -u postgres psql perjalanan_dinas -f clean-production-database.sql
```

### 5. Verify
```bash
sudo -u postgres psql perjalanan_dinas -f verify-cleanup.sql
```

### 6. Start Backend
```bash
sudo systemctl start perjalanan-dinas-backend
```

### 7. Test
```bash
# Create test SPD via frontend
# Expected: 064/0325/DIB/{CODE}/NOTA
```

---

## 📋 Manual Commands (Alternative)

```bash
# 1. Backup
sudo pg_dump -U postgres perjalanan_dinas > /var/backups/backup_$(date +%Y%m%d).sql

# 2. Cleanup & Set Sequence
sudo -u postgres psql perjalanan_dinas <<EOF
DELETE FROM visit_proofs;
DELETE FROM travel_reports;
DELETE FROM travel_request_employees;
DELETE FROM travel_requests;
UPDATE numbering_configs SET last_request_sequence = 324 WHERE id = 1;
SELECT last_request_sequence FROM numbering_configs WHERE id = 1;
EOF

# 3. Restart
sudo systemctl restart perjalanan-dinas-backend
```

---

## ✅ Expected Results
- Travel Requests: **0 records**
- SPD Sequence: **324**
- Next SPD: **325** (format: 064/0325/DIB/XXXX/NOTA)
- Master Data: **PRESERVED**

---

## 🆘 Rollback
```bash
sudo systemctl stop perjalanan-dinas-backend
sudo -u postgres psql perjalanan_dinas < /var/backups/backup_YYYYMMDD.sql
sudo systemctl start perjalanan-dinas-backend
```

---

## 📚 Detailed Guide
Lihat **CLEANUP-DATABASE-GUIDE.md** untuk panduan lengkap.

## 📁 All Files
1. `backup-full-database.sh` - Backup script
2. `clean-production-database.sql` - Main cleanup
3. `verify-cleanup.sql` - Verification
4. `CLEANUP-DATABASE-GUIDE.md` - Full guide
5. `QUICK-START.md` - This file
