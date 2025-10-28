# Database Management Scripts - Perjalanan Dinas SPD

Koleksi script untuk management database production, termasuk cleanup data dan update sequence SPD.

## 📁 Struktur Files

```
scripts/
├── README.md                          # File ini
├── QUICK-START.md                     # Quick reference untuk cleanup
├── CLEANUP-DATABASE-GUIDE.md          # Panduan lengkap cleanup database
├── UPDATE-SPD-GUIDE.md                # Panduan update sequence SPD
├── backup-full-database.sh            # Backup script (executable)
├── clean-production-database.sql      # Main cleanup script
├── verify-cleanup.sql                 # Verification script
├── update-spd-sequence.sql            # Update sequence only
└── backup-before-update.sql           # Backup config only
```

## 🎯 Use Cases

### 1. Cleanup Database + Set Sequence ke 324

**Scenario:** Hapus semua data transaksi, tapi pertahankan data master dan set SPD sequence ke 324.

**Files:**
- `CLEANUP-DATABASE-GUIDE.md` - Panduan lengkap
- `backup-full-database.sh` - Backup
- `clean-production-database.sql` - Cleanup
- `verify-cleanup.sql` - Verifikasi

**Quick Command:**
```bash
# 1. Backup
./backup-full-database.sh

# 2. Cleanup
psql -U postgres -d perjalanan_dinas -f clean-production-database.sql

# 3. Verify
psql -U postgres -d perjalanan_dinas -f verify-cleanup.sql
```

**Result:**
- ✅ Semua travel requests dihapus
- ✅ Data master (employees, positions) tetap
- ✅ SPD sequence = 324
- ✅ SPD berikutnya = 325

---

### 2. Update Sequence SPD Saja (Tanpa Cleanup)

**Scenario:** Hanya update nomor sequence SPD ke 324, tidak hapus data.

**Files:**
- `UPDATE-SPD-GUIDE.md` - Panduan
- `update-spd-sequence.sql` - Update script
- `backup-before-update.sql` - Backup config

**Quick Command:**
```bash
psql -U postgres -d perjalanan_dinas -f update-spd-sequence.sql
```

**Result:**
- ✅ SPD sequence = 324
- ✅ Semua data tetap utuh
- ✅ SPD berikutnya = 325

---

## 📖 Documentation

### Panduan Lengkap

1. **CLEANUP-DATABASE-GUIDE.md**
   - Step-by-step cleanup database
   - Backup procedures
   - Verification steps
   - Rollback procedures
   - Troubleshooting
   - **Baca ini jika first time cleanup!**

2. **UPDATE-SPD-GUIDE.md**
   - Update sequence number only
   - Tidak hapus data
   - Simple dan cepat

3. **QUICK-START.md**
   - Quick reference
   - Untuk yang sudah familiar
   - One-page cheatsheet

### Scripts Detail

#### backup-full-database.sh
```bash
# Backup full database PostgreSQL
# Output: /var/backups/perjalanan-dinas/perjalanan_dinas_before_cleanup_TIMESTAMP.sql

sudo ./backup-full-database.sh
```

**Features:**
- Auto create backup directory
- Verify database connection
- Show database statistics
- Create backup file
- Generate restore script
- Colored output

---

#### clean-production-database.sql
```sql
-- Main cleanup script
-- 1. Backup data ke tabel *_backup_cleanup
-- 2. Delete transactional data
-- 3. Reset auto-increment sequences
-- 4. Set SPD sequence ke 324
-- 5. Show summary

\i clean-production-database.sql
```

**Actions:**
- ✅ DELETE visit_proofs
- ✅ DELETE travel_reports
- ✅ DELETE travel_request_employees
- ✅ DELETE travel_requests
- ✅ RESET sequences to 1
- ✅ UPDATE SPD sequence to 324
- ❌ PRESERVE employees
- ❌ PRESERVE positions
- ❌ PRESERVE admins
- ❌ PRESERVE representative_configs

---

#### verify-cleanup.sql
```sql
-- Comprehensive verification
-- Checks:
--   1. Transactional data = 0
--   2. Master data preserved
--   3. SPD sequence = 324
--   4. Auto-increment sequences
--   5. Backup tables exist
--   6. Summary report

\i verify-cleanup.sql
```

**Output:**
- Detailed verification results
- Status indicators (✓/✗)
- Summary report
- Next SPD format

---

#### update-spd-sequence.sql
```sql
-- Update sequence only (no cleanup)
-- 1. Show current sequence
-- 2. Update to 324
-- 3. Verify result

\i update-spd-sequence.sql
```

---

#### backup-before-update.sql
```sql
-- Backup config tables only
-- Lighter than full backup
-- Include restore commands

\i backup-before-update.sql
```

---

## 🚀 Quick Examples

### Example 1: Full Cleanup untuk Fresh Start

```bash
# SSH ke production
ssh admin@103.160.37.195

# Navigate
cd /path/to/scripts

# Backup
sudo ./backup-full-database.sh
# ✓ Backup created: /var/backups/perjalanan-dinas/...

# Stop service
sudo systemctl stop perjalanan-dinas-backend

# Cleanup
sudo -u postgres psql perjalanan_dinas -f clean-production-database.sql
# ✓ All transactional data deleted
# ✓ SPD sequence set to 324

# Verify
sudo -u postgres psql perjalanan_dinas -f verify-cleanup.sql
# ✓ All checks passed

# Start service
sudo systemctl start perjalanan-dinas-backend

# Test di browser
# Create SPD → Expected: 064/0325/DIB/XXXX/NOTA ✓
```

---

### Example 2: Update Sequence Tanpa Hapus Data

```bash
# SSH ke production
ssh admin@103.160.37.195

# Quick backup
sudo pg_dump -U postgres perjalanan_dinas > /tmp/backup.sql

# Update sequence
sudo -u postgres psql perjalanan_dinas -f update-spd-sequence.sql

# Test
# Create SPD → Next number will be 325
```

---

### Example 3: Rollback dari Backup

```bash
# Stop service
sudo systemctl stop perjalanan-dinas-backend

# Restore database
sudo -u postgres psql perjalanan_dinas < /var/backups/perjalanan-dinas/backup_TIMESTAMP.sql

# Start service
sudo systemctl start perjalanan-dinas-backend

# Verify
curl http://localhost:8080/api/v1/health
```

---

## ⚠️ Important Notes

### Sebelum Eksekusi:

1. **SELALU backup database terlebih dahulu**
2. **Inform users** bahwa akan ada maintenance
3. **Koordinasi dengan team**
4. **Siapkan rollback plan**
5. **Test di staging** dulu jika memungkinkan

### Setelah Eksekusi:

1. **Verify** semua checks passed
2. **Test** create SPD baru
3. **Monitor** logs untuk errors
4. **Inform users** maintenance selesai
5. **Keep backup** sampai sistem stabil

### Data yang Dihapus:

- ❌ Travel Requests (all)
- ❌ Travel Reports (all)
- ❌ Travel Request Employees (all)
- ❌ Visit Proofs (all)

### Data yang Tetap:

- ✅ Employees
- ✅ Positions
- ✅ Admins
- ✅ Representative Configs
- ✅ Numbering Config (updated to 324)

---

## 🐛 Troubleshooting

### Problem: Cannot Connect to Database

```bash
# Check PostgreSQL
sudo systemctl status postgresql
sudo systemctl restart postgresql

# Check database exists
sudo -u postgres psql -l | grep perjalanan
```

### Problem: Permission Denied

```bash
# Run with sudo
sudo -u postgres psql perjalanan_dinas -f script.sql

# Or login as postgres first
sudo su - postgres
psql perjalanan_dinas -f /path/to/script.sql
```

### Problem: Foreign Key Constraint

```sql
-- In psql, disable temporarily
SET session_replication_role = 'replica';
-- Run deletes
DELETE FROM ...;
-- Re-enable
SET session_replication_role = 'origin';
```

### Problem: Backup Failed

```bash
# Check disk space
df -h

# Check permissions
ls -la /var/backups/

# Create directory manually
sudo mkdir -p /var/backups/perjalanan-dinas
sudo chmod 755 /var/backups/perjalanan-dinas
```

---

## 📞 Support

Untuk bantuan lebih lanjut:

1. Review **CLEANUP-DATABASE-GUIDE.md** untuk troubleshooting lengkap
2. Check application logs: `journalctl -u perjalanan-dinas-backend -f`
3. Check PostgreSQL logs: `/var/log/postgresql/`
4. Contact developer team

---

## 📊 Checklist Template

```
[ ] Backup database selesai
[ ] Backup file verified
[ ] Service stopped
[ ] Cleanup script executed
[ ] Verification passed
[ ] Transactional data = 0
[ ] Master data intact
[ ] SPD sequence = 324
[ ] Service started
[ ] Login test OK
[ ] Create SPD test OK
[ ] SPD number = 325
[ ] No errors in logs
[ ] Users informed
```

---

## 🔖 Version History

- **2025-10-27**: Initial version
  - Full cleanup script
  - Verification script
  - Backup automation
  - Comprehensive documentation

---

**Need Help?** Start with **CLEANUP-DATABASE-GUIDE.md** untuk panduan step-by-step lengkap.

**Quick Reference?** Check **QUICK-START.md** untuk commands cepat.

**Good luck! 🚀**
