# Panduan Cleanup Database Production - SPD System

## ⚠️ PERINGATAN PENTING

**MEMBACA PANDUAN INI SEPENUHNYA SEBELUM EKSEKUSI!**

Script ini akan **MENGHAPUS SEMUA DATA TRANSAKSI** dari database production, namun:
- ✅ **Data Master TETAP** (Employees, Positions, Admins, Representative Configs)
- ✅ **Nomor Sequence SPD = 324** (SPD berikutnya = 325)
- ❌ **Semua Travel Requests DIHAPUS**
- ❌ **Semua Travel Reports DIHAPUS**
- ❌ **Semua Visit Proofs DIHAPUS**

## 📋 Checklist Persiapan

Sebelum eksekusi, pastikan:

- [ ] Sudah inform semua user bahwa sistem akan di-maintenance
- [ ] Sudah koordinasi dengan team terkait
- [ ] Sudah download/archive dokumen PDF/Excel jika diperlukan
- [ ] Sudah siap untuk mulai dari nomor SPD 325
- [ ] Punya akses SSH ke production server
- [ ] Punya akses PostgreSQL dengan privilege DELETE

## 🗂️ Files yang Dibutuhkan

1. **backup-full-database.sh** - Full database backup script
2. **clean-production-database.sql** - Main cleanup script
3. **verify-cleanup.sql** - Verification script
4. **CLEANUP-DATABASE-GUIDE.md** - Panduan ini

## 📖 Step-by-Step Guide

### Step 1: Koneksi ke Production Server

```bash
# SSH ke VPS production
ssh admin@103.160.37.195

# Navigate ke directory scripts
cd /path/to/perjalanan-dinas/backend/scripts
```

### Step 2: Backup Database (WAJIB!)

#### Option A: Menggunakan Backup Script (Recommended)

```bash
# Beri permission execute
chmod +x backup-full-database.sh

# Jalankan backup script
sudo ./backup-full-database.sh
```

Output yang diharapkan:
```
============================================================================
Database Backup Script - Perjalanan Dinas SPD
============================================================================

[1/5] Creating backup directory...
✓ Backup directory ready: /var/backups/perjalanan-dinas

[2/5] Checking database connection...
✓ Database connection successful

[3/5] Database statistics before backup:
  table_name        | total_records
--------------------+--------------
  travel_requests   | X
  travel_reports    | X
  employees         | X
  positions         | X
  admins            | X

[4/5] Creating full database backup...
✓ Backup created successfully

[5/5] Verifying backup file...
✓ Backup file exists
   Location: /var/backups/perjalanan-dinas/perjalanan_dinas_before_cleanup_YYYYMMDD_HHMMSS.sql
   Size: XXX MB

============================================================================
BACKUP COMPLETED SUCCESSFULLY
============================================================================
```

#### Option B: Manual Backup

```bash
# Create backup directory
sudo mkdir -p /var/backups/perjalanan-dinas

# Run pg_dump
sudo pg_dump -U postgres perjalanan_dinas > /var/backups/perjalanan-dinas/backup_$(date +%Y%m%d_%H%M%S).sql

# Verify backup file
ls -lh /var/backups/perjalanan-dinas/
```

**⚠️ JANGAN LANJUT jika backup GAGAL!**

### Step 3: Stop Backend Service (Optional, tapi Recommended)

```bash
# Cek status service
sudo systemctl status perjalanan-dinas-backend

# Stop service untuk mencegah write operations selama cleanup
sudo systemctl stop perjalanan-dinas-backend

# Verify stopped
sudo systemctl status perjalanan-dinas-backend
```

### Step 4: Koneksi ke Database

```bash
# Koneksi ke PostgreSQL
sudo -u postgres psql perjalanan_dinas

# Atau dengan user lain
psql -U postgres -d perjalanan_dinas
```

### Step 5: Eksekusi Cleanup Script

Di dalam psql console:

```sql
-- Load dan jalankan cleanup script
\i /path/to/scripts/clean-production-database.sql

-- Atau copy-paste isi file script ke console
```

Script akan:
1. ✅ Backup data ke tabel `*_backup_cleanup`
2. ✅ Tampilkan data sebelum cleanup
3. ✅ Hapus semua data transaksi
4. ✅ Reset auto-increment sequences
5. ✅ Set SPD sequence ke 324
6. ✅ Tampilkan data setelah cleanup
7. ✅ Tampilkan summary

### Step 6: Verifikasi Hasil Cleanup

```sql
-- Jalankan verification script
\i /path/to/scripts/verify-cleanup.sql

-- Atau manual check:
-- 1. Check transactional data (harus 0)
SELECT COUNT(*) FROM travel_requests;    -- Expected: 0
SELECT COUNT(*) FROM travel_reports;     -- Expected: 0

-- 2. Check master data (harus ada)
SELECT COUNT(*) FROM employees;          -- Expected: > 0
SELECT COUNT(*) FROM positions;          -- Expected: > 0

-- 3. Check SPD sequence (harus 324)
SELECT last_request_sequence FROM numbering_configs WHERE id = 1;  -- Expected: 324
```

Expected Output dari verify-cleanup.sql:

```
============================================================================
VERIFIKASI DATABASE SETELAH CLEANUP
============================================================================

1. VERIFIKASI DATA TRANSAKSI (Expected: 0 records)
------------------------------------------------------------
        table_name         | total_records |          status
---------------------------+---------------+---------------------------
 visit_proofs              |             0 | ✓ OK
 travel_reports            |             0 | ✓ OK
 travel_request_employees  |             0 | ✓ OK
 travel_requests           |             0 | ✓ OK

2. VERIFIKASI DATA MASTER (Must have data)
------------------------------------------------------------
        table_name         | total_records |          status
---------------------------+---------------+---------------------------
 employees                 |            XX | ✓ OK - Data preserved
 positions                 |            XX | ✓ OK - Data preserved
 admins                    |             1 | ✓ OK - Data preserved
 representative_configs    |             1 | ✓ OK - Data preserved

3. VERIFIKASI NOMOR SEQUENCE SPD (Expected: 324)
------------------------------------------------------------
 id | current_sequence |        status        |           next_spd_format
----+------------------+----------------------+-------------------------------------
  1 |              324 | ✓ OK - Set to 324    | Next SPD will be: 064/0325/DIB/{CODE}/NOTA

============================================================================
CLEANUP SUMMARY REPORT
============================================================================
      category      | count |    status
--------------------+-------+--------------
 Transactional Data |     0 | ✓ CLEANED
 Master Data        |    XX | ✓ PRESERVED
 SPD Sequence       |   324 | ✓ SET TO 324
```

### Step 7: Exit Database dan Start Backend Service

```bash
# Exit dari psql
\q

# Start backend service kembali
sudo systemctl start perjalanan-dinas-backend

# Verify service running
sudo systemctl status perjalanan-dinas-backend

# Check logs untuk memastikan tidak ada error
sudo journalctl -u perjalanan-dinas-backend -f
```

### Step 8: Test di Aplikasi

#### Test 1: Login ke Frontend

```
1. Buka browser: http://103.160.37.195:3000
2. Login dengan admin credentials
3. Pastikan bisa login dan tidak ada error
```

#### Test 2: Buat SPD Baru

```
1. Buka menu "Buat Perjalanan Dinas"
2. Isi form dengan data test:
   - Pilih Employee
   - Maksud: "Test SPD 325"
   - Tujuan: "Jakarta"
   - Jenis: "Luar Provinsi"
   - Tanggal: [pilih tanggal]
   - Transportasi: "Pesawat"
3. Submit
4. PASTIKAN nomor SPD = "064/0325/DIB/{CODE}/NOTA"
```

#### Test 3: Verify via API

```bash
# Get last travel request
curl -X GET http://103.160.37.195:8080/api/v1/travel-requests \
  -H "Authorization: Bearer YOUR_TOKEN" | jq

# Check request_number field
# Expected: "064/0325/DIB/XXXX/NOTA"
```

## 🔄 Rollback Procedure

Jika terjadi masalah dan perlu restore:

### Option 1: Restore dari Full Backup

```bash
# Stop backend service
sudo systemctl stop perjalanan-dinas-backend

# Drop dan recreate database
sudo -u postgres psql -c "DROP DATABASE perjalanan_dinas;"
sudo -u postgres psql -c "CREATE DATABASE perjalanan_dinas;"

# Restore dari backup
sudo -u postgres psql perjalanan_dinas < /var/backups/perjalanan-dinas/backup_TIMESTAMP.sql

# Start service
sudo systemctl start perjalanan-dinas-backend
```

### Option 2: Restore Hanya Numbering Config

```sql
-- Jika hanya perlu restore sequence number
UPDATE numbering_configs nc
SET
    last_request_sequence = ncb.last_request_sequence,
    last_report_sequence = ncb.last_report_sequence,
    updated_at = ncb.updated_at
FROM numbering_configs_backup_cleanup ncb
WHERE nc.id = ncb.id;
```

## 🐛 Troubleshooting

### Problem 1: Backup Script Permission Denied

```bash
# Solusi: Jalankan dengan sudo
sudo ./backup-full-database.sh

# Atau beri permission
chmod +x backup-full-database.sh
sudo chown $USER:$USER backup-full-database.sh
```

### Problem 2: Cannot Connect to Database

```bash
# Check PostgreSQL status
sudo systemctl status postgresql

# Check if database exists
sudo -u postgres psql -l | grep perjalanan

# Restart PostgreSQL jika perlu
sudo systemctl restart postgresql
```

### Problem 3: Foreign Key Constraint Error

Jika muncul error saat DELETE:

```sql
-- Disable foreign key checks temporarily
SET session_replication_role = 'replica';

-- Run delete commands
DELETE FROM visit_proofs;
DELETE FROM travel_reports;
DELETE FROM travel_request_employees;
DELETE FROM travel_requests;

-- Enable foreign key checks
SET session_replication_role = 'origin';
```

### Problem 4: Sequence Not Updated

```sql
-- Force update sequence
UPDATE numbering_configs
SET last_request_sequence = 324, updated_at = NOW()
WHERE id = 1;

-- Verify
SELECT * FROM numbering_configs;
```

### Problem 5: Backend Service Won't Start

```bash
# Check logs
sudo journalctl -u perjalanan-dinas-backend -n 50

# Check config
cat /etc/systemd/system/perjalanan-dinas-backend.service

# Restart
sudo systemctl daemon-reload
sudo systemctl restart perjalanan-dinas-backend
```

## ✅ Post-Cleanup Checklist

- [ ] Backup file tersimpan dan terverifikasi
- [ ] Data transaksi = 0 records
- [ ] Data master tetap ada
- [ ] SPD sequence = 324
- [ ] Backend service running
- [ ] Login frontend berhasil
- [ ] SPD baru menggunakan nomor 325
- [ ] API response correct
- [ ] Log tidak ada error
- [ ] Inform user bahwa maintenance selesai

## 📊 Expected Results Summary

| Item | Before Cleanup | After Cleanup |
|------|----------------|---------------|
| Travel Requests | X records | 0 records |
| Travel Reports | X records | 0 records |
| Visit Proofs | X records | 0 records |
| Employees | X records | X records (SAME) |
| Positions | X records | X records (SAME) |
| Admins | X records | X records (SAME) |
| SPD Sequence | X | 324 |
| Next SPD Number | - | 325 |

## 📝 Execution Log Template

```
===============================================
DATABASE CLEANUP EXECUTION LOG
===============================================
Date: _______________
Time: _______________
Executed by: _______________
Server: 103.160.37.195

Pre-Cleanup:
  - Travel Requests: _____ records
  - SPD Sequence: _____
  - Backup File: _____________________
  - Backup Size: _____ MB

Cleanup Execution:
  - Start Time: _____
  - End Time: _____
  - Duration: _____ minutes
  - Status: [ ] Success  [ ] Failed

Post-Cleanup:
  - Travel Requests: 0 records ✓
  - SPD Sequence: 324 ✓
  - Next SPD: 325 ✓
  - Test SPD Created: _______________

Issues/Notes:
_________________________________________
_________________________________________

Sign-off: _______________
===============================================
```

## 📞 Support

Jika mengalami masalah:
1. Check troubleshooting section di atas
2. Review logs di `/var/log/`
3. Contact developer team
4. Keep backup file sampai confirm sistem stabil

---

**REMINDER:**
- Selalu backup sebelum cleanup!
- Koordinasi dengan team sebelum eksekusi
- Test semua functionality setelah cleanup
- Document semua yang dilakukan

**Good luck! 🚀**
