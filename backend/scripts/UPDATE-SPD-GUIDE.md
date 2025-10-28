# Panduan Update Nomor SPD ke 324 (Production)

## Tujuan
Mengatur nomor SPD terakhir di production menjadi **324**, sehingga SPD berikutnya yang dibuat akan memiliki nomor **325**.

## Informasi Sistem
- **Tabel:** `numbering_configs`
- **Field:** `last_request_sequence`
- **Format Nomor SPD:** `064/{seq}/DIB/{code}/NOTA`
- **Contoh:** `064/0325/DIB/DPEB/NOTA` (untuk SPD ke-325)

## Langkah-langkah Eksekusi

### Persiapan

1. **Login ke Production Server**
   ```bash
   ssh admin@103.160.37.195
   ```

2. **Koneksi ke Database**
   ```bash
   psql -U postgres -d perjalanan_dinas
   ```

### Eksekusi Update

#### Step 1: Backup Data (WAJIB!)

```sql
-- Jalankan backup script
\i /path/to/scripts/backup-before-update.sql

-- Atau manual:
CREATE TABLE IF NOT EXISTS numbering_configs_backup AS
SELECT * FROM numbering_configs;

-- Verifikasi backup
SELECT * FROM numbering_configs_backup;
```

#### Step 2: Cek Nomor Saat Ini

```sql
-- Lihat nomor SPD terakhir saat ini
SELECT
    id,
    last_request_sequence AS nomor_spd_saat_ini,
    last_report_sequence,
    updated_at
FROM numbering_configs
WHERE id = 1;
```

**Catat nomor yang muncul untuk referensi!**

#### Step 3: Update Nomor ke 324

```sql
-- Update nomor SPD terakhir
UPDATE numbering_configs
SET
    last_request_sequence = 324,
    updated_at = NOW()
WHERE id = 1;
```

#### Step 4: Verifikasi Update

```sql
-- Verifikasi hasil update
SELECT
    id,
    last_request_sequence AS nomor_spd_terakhir,
    last_report_sequence,
    updated_at
FROM numbering_configs
WHERE id = 1;
```

**Expected Result:**
```
 id | nomor_spd_terakhir | last_report_sequence | updated_at
----+-------------------+----------------------+------------
  1 |               324 |        (unchanged)   | 2025-10-27...
```

### Verifikasi di Aplikasi

1. **Test Create SPD Baru:**
   - Login ke aplikasi frontend
   - Buat travel request baru
   - Pastikan nomor SPD yang dibuat adalah: `064/0325/DIB/{code}/NOTA`

2. **Cek via API:**
   ```bash
   curl -X POST http://103.160.37.195:8080/api/v1/travel-requests \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "employee_ids": [1],
       "purpose": "Test nomor SPD",
       "destination": "Jakarta",
       "destination_type": "outside_province",
       "departure_date": "2025-11-01",
       "return_date": "2025-11-03",
       "transportation": "pesawat"
     }'
   ```

3. **Cek response:**
   ```json
   {
     "message": "Travel request created successfully",
     "travel_request": {
       "request_number": "064/0325/DIB/DPEB/NOTA",
       ...
     }
   }
   ```

## Rollback (Jika Diperlukan)

Jika terjadi kesalahan dan perlu mengembalikan nomor:

```sql
-- Restore dari backup
UPDATE numbering_configs nc
SET
    last_request_sequence = ncb.last_request_sequence,
    last_report_sequence = ncb.last_report_sequence,
    updated_at = ncb.updated_at
FROM numbering_configs_backup ncb
WHERE nc.id = ncb.id;

-- Verifikasi restore
SELECT * FROM numbering_configs;
```

## Troubleshooting

### Problem 1: Nomor tidak berubah setelah update

**Solusi:**
```sql
-- Cek apakah ada multiple rows
SELECT COUNT(*) FROM numbering_configs;

-- Jika lebih dari 1, update semua
UPDATE numbering_configs SET last_request_sequence = 324;
```

### Problem 2: SPD baru masih menggunakan nomor lama

**Penyebab:** Aplikasi mungkin caching atau belum restart

**Solusi:**
```bash
# Restart backend service
sudo systemctl restart perjalanan-dinas-backend

# Atau jika menggunakan docker
docker restart perjalanan-dinas-backend
```

### Problem 3: Database connection timeout

**Solusi:**
```bash
# Cek status database
sudo systemctl status postgresql

# Restart jika perlu
sudo systemctl restart postgresql
```

## Checklist Eksekusi

- [ ] Backup database dibuat
- [ ] Nomor SPD saat ini dicatat
- [ ] Update query dijalankan
- [ ] Verifikasi query menunjukkan 324
- [ ] Test create SPD baru
- [ ] SPD baru menunjukkan nomor 325
- [ ] Dokumentasi update dicatat

## Log Eksekusi

```
Date: _____________
Executed by: _____________
Previous sequence number: _____________
New sequence number: 324
First SPD after update: _____________
Status: [ ] Success  [ ] Failed  [ ] Rolled back
Notes: _____________________________________________
```

## Kontak

Jika ada masalah, hubungi:
- Developer: [Your Contact]
- Database Admin: [Your Contact]

---

**CATATAN PENTING:**
- Selalu backup sebelum update production!
- Update ini hanya perlu dilakukan SEKALI
- Setelah ini sistem akan auto-increment dari 324 → 325 → 326 dst.
