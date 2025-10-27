# Project Summary - Sistem Perjalanan Dinas

## Overview

Aplikasi web full-stack untuk digitalisasi proses penginputan dan pengelolaan surat perjalanan dinas karyawan. Menggantikan proses manual menggunakan Microsoft Word dengan sistem otomatis yang menghasilkan PDF terstandar.

## Problem Statement

### Kondisi Saat Ini (Manual)
- Karyawan harus mengisi dokumen Word secara manual
- Input berulang untuk data yang sama (NIP, Nama, Jabatan)
- Penomoran manual rawan kesalahan
- Format dokumen tidak konsisten
- Sulit tracking dan manajemen data perjalanan dinas
- Penyimpanan file tidak terorganisir

### Solusi yang Diberikan
- Form web yang user-friendly
- Auto-fill data karyawan dari database
- Sistem penomoran otomatis dan terstandar
- Generate PDF dengan format konsisten
- Database terpusat untuk semua perjalanan dinas
- Admin panel untuk manajemen data

## Features

### Public Features (Tanpa Login)
1. **Form Perjalanan Dinas**
   - Pilih karyawan dari dropdown
   - Auto-fill NIP, Nama, Jabatan
   - Input maksud perjalanan
   - Input tempat tujuan
   - Pilih tanggal berangkat dan kembali
   - Kalkulasi otomatis lama perjalanan
   - Pilih jenis angkutan
   - Submit dan download PDF Nota Permintaan

### Admin Features (Dengan Login)
1. **Dashboard**
   - Statistik perjalanan dinas
   - Quick access menu
   - Informasi sistem

2. **Kelola Karyawan**
   - Tambah karyawan baru
   - Edit data karyawan
   - Hapus karyawan
   - Lihat daftar semua karyawan

3. **Pengkodean Jabatan**
   - Tambah kode untuk setiap jabatan
   - Edit kode jabatan
   - Hapus kode jabatan
   - Preview format nomor

4. **Daftar Perjalanan Dinas**
   - Lihat semua perjalanan dinas
   - Filter berdasarkan status
   - Download PDF (Nota Permintaan, Berita Acara, atau Combined)
   - Hapus perjalanan dinas

## Technical Stack

### Backend
- **Language:** Go 1.21+
- **Framework:** Gin (HTTP web framework)
- **Database:** PostgreSQL 14+
- **ORM:** GORM
- **Authentication:** JWT
- **PDF Generation:** gofpdf

### Frontend
- **Framework:** Next.js 14 (App Router)
- **Language:** TypeScript
- **UI Library:** React 18
- **Styling:** Tailwind CSS
- **Form Management:** React Hook Form
- **HTTP Client:** Axios
- **Date Handling:** date-fns

### DevOps
- **Database Migrations:** SQL scripts
- **Process Manager:** systemd (backend), PM2 (frontend)
- **Reverse Proxy:** Nginx
- **SSL:** Let's Encrypt
- **Containerization:** Docker (optional)

## Project Structure

```
perjalanan-dinas/
├── backend/
│   ├── cmd/api/              # Main application entry point
│   ├── config/               # Configuration management
│   ├── internal/
│   │   ├── database/         # Database connection & migrations
│   │   ├── handlers/         # HTTP request handlers
│   │   │   ├── auth_handler.go
│   │   │   ├── employee_handler.go
│   │   │   ├── position_code_handler.go
│   │   │   ├── travel_request_handler.go
│   │   │   ├── travel_report_handler.go
│   │   │   └── pdf_handler.go
│   │   ├── middleware/       # Authentication & CORS
│   │   ├── models/           # Database models
│   │   ├── repository/       # Data access layer
│   │   ├── services/         # Business logic (PDF generator)
│   │   └── utils/            # Utility functions
│   ├── migrations/           # SQL migration files
│   ├── .env.example
│   ├── .gitignore
│   ├── go.mod
│   ├── go.sum
│   └── Makefile
│
├── frontend/
│   ├── src/
│   │   ├── app/
│   │   │   ├── admin/        # Admin panel pages
│   │   │   │   ├── login/
│   │   │   │   ├── dashboard/
│   │   │   │   ├── employees/
│   │   │   │   ├── position-codes/
│   │   │   │   └── travel-requests/
│   │   │   ├── layout.tsx
│   │   │   ├── page.tsx      # Public form page
│   │   │   └── globals.css
│   │   ├── components/       # React components
│   │   │   └── AdminLayout.tsx
│   │   ├── lib/              # API client & utilities
│   │   │   └── api.ts
│   │   └── types/            # TypeScript type definitions
│   │       └── index.ts
│   ├── public/               # Static assets
│   ├── .env.local.example
│   ├── .gitignore
│   ├── next.config.js
│   ├── package.json
│   ├── tailwind.config.ts
│   ├── tsconfig.json
│   └── postcss.config.js
│
├── README.md                 # Main documentation
├── QUICK_START.md           # Quick start guide
├── API_DOCUMENTATION.md     # API reference
├── DEPLOYMENT.md            # Deployment guide
└── PROJECT_SUMMARY.md       # This file
```

## Database Schema

### Tables

1. **admins**
   - Admin user credentials
   - Fields: id, username, password (hashed)

2. **employees**
   - Karyawan data
   - Fields: id, nip, name, position

3. **position_codes**
   - Kode untuk setiap jabatan
   - Fields: id, position, code

4. **numbering_configs**
   - Auto-increment sequences
   - Fields: id, last_request_sequence, last_report_sequence

5. **travel_requests**
   - Nota Permintaan Surat Tugas
   - Fields: employee_id, purpose, destination, dates, transportation, request_number, status

6. **travel_reports**
   - Berita Acara Perjalanan Dinas
   - Fields: travel_request_id, report_number, representative info

7. **visit_proofs**
   - Bukti kunjungan dalam berita acara
   - Fields: travel_report_id, date, locations, signature proof

### Relationships
- Employee → Travel Requests (1:N)
- Travel Request → Travel Report (1:1)
- Travel Report → Visit Proofs (1:N)

## Document Numbering System

### Nota Permintaan Surat Tugas Perjalanan Dinas
**Format:** `064/[sequence]/DIB/[position_code]/NOTA`

**Example:** `064/0001/DIB/DIRUT/NOTA`

**Components:**
- `064` - Fixed prefix
- `[sequence]` - Auto-increment (0001, 0002, ...)
- `DIB` - Fixed code
- `[position_code]` - Dynamic based on employee position
- `NOTA` - Fixed suffix

### Berita Acara Perjalanan Dinas
**Format:** `064/    /DIB/[position_code]/NOTA`

**Example:** `064/    /DIB/DIRUT/NOTA`

**Note:** Space instead of sequence number

## User Flows

### Flow 1: Karyawan Membuat Perjalanan Dinas
1. Buka aplikasi (no login required)
2. Pilih nama dari dropdown
3. Sistem auto-fill NIP, Nama, Jabatan
4. Isi maksud perjalanan dinas
5. Pilih tempat tujuan
6. Pilih tanggal berangkat dan kembali
7. Sistem kalkulasi otomatis durasi
8. Pilih jenis angkutan
9. Submit form
10. Download PDF Nota Permintaan

### Flow 2: Admin Setup Master Data
1. Login ke admin panel
2. Tambah kode jabatan:
   - Input nama jabatan
   - Input kode (3-4 huruf)
   - Save
3. Tambah karyawan:
   - Input NIP
   - Input Nama
   - Pilih Jabatan (harus sudah ada di kode jabatan)
   - Save

### Flow 3: Admin Melihat Data
1. Login ke admin panel
2. Menu "Daftar Perjalanan Dinas"
3. Lihat semua perjalanan dinas
4. Download PDF yang diperlukan
5. Hapus jika perlu

## API Endpoints

### Public Endpoints
- `POST /api/auth/login` - Login admin
- `GET /api/employees` - List employees
- `GET /api/position-codes` - List position codes
- `POST /api/travel-requests` - Create travel request
- `GET /api/pdf/*` - Download PDFs

### Admin Endpoints (Require JWT)
- `POST /api/admin/employees` - Create employee
- `PUT /api/admin/employees/:id` - Update employee
- `DELETE /api/admin/employees/:id` - Delete employee
- `POST /api/admin/position-codes` - Create position code
- `PUT /api/admin/position-codes/:id` - Update position code
- `DELETE /api/admin/position-codes/:id` - Delete position code
- `GET /api/admin/travel-requests` - List all travel requests
- `DELETE /api/admin/travel-requests/:id` - Delete travel request
- `POST /api/admin/travel-reports` - Create travel report
- `GET /api/admin/travel-reports/:id` - Get travel report

## Security Features

1. **Authentication**
   - JWT-based authentication for admin
   - Bcrypt password hashing
   - Token expiration (24 hours)

2. **Authorization**
   - Role-based access (admin vs public)
   - Protected routes with middleware

3. **Data Validation**
   - Input validation on both frontend and backend
   - SQL injection prevention via ORM
   - XSS protection

4. **Best Practices**
   - Environment variables for sensitive data
   - HTTPS recommended for production
   - CORS configuration
   - Soft deletes for data integrity

## Performance Considerations

1. **Database**
   - Indexed columns for faster queries
   - Connection pooling
   - Optimized queries with GORM

2. **Frontend**
   - Code splitting with Next.js
   - Static generation where possible
   - Lazy loading components
   - Optimized images

3. **API**
   - Efficient JSON serialization
   - Minimal data transfer
   - Caching headers

## Testing Strategy

### Backend Testing
- Unit tests for business logic
- Integration tests for API endpoints
- Database transaction tests

### Frontend Testing
- Component testing with React Testing Library
- E2E testing with Cypress (recommended)
- Form validation testing

## Deployment Options

### Option 1: Traditional (systemd + PM2)
- Backend: systemd service
- Frontend: PM2 process manager
- Proxy: Nginx
- Database: PostgreSQL instance

### Option 2: Docker Compose
- All services containerized
- Easy scaling
- Environment isolation

### Option 3: Cloud Platform
- Backend: Cloud Run / AWS ECS
- Frontend: Vercel / Netlify
- Database: Cloud SQL / RDS

## Monitoring & Logging

### Logs Location
- Backend: systemd journal (`journalctl -u perjalanan-dinas-api`)
- Frontend: PM2 logs (`pm2 logs`)
- Nginx: `/var/log/nginx/`

### Metrics to Monitor
- API response times
- Database query performance
- Error rates
- PDF generation success rate
- User activity

## Future Enhancements

### Phase 2 Features
1. **Email Notifications**
   - Notifikasi ke admin saat ada perjalanan dinas baru
   - Email approval workflow

2. **Advanced Reporting**
   - Export to Excel
   - Charts and analytics
   - Monthly/yearly reports

3. **Berita Acara Form**
   - Form untuk mengisi berita acara
   - Upload bukti tanda tangan/stempel

4. **Approval Workflow**
   - Multi-level approval
   - Status tracking
   - Comments/notes

5. **Mobile App**
   - React Native / Flutter
   - Mobile-optimized UI
   - Push notifications

6. **Audit Trail**
   - Track all changes
   - User activity log
   - Data version history

## Maintenance

### Regular Tasks
- [ ] Database backup (daily)
- [ ] Log rotation (weekly)
- [ ] Security updates (monthly)
- [ ] Performance review (monthly)
- [ ] User feedback review (monthly)

### Backup Strategy
- Database: Daily automated backup
- Files: Weekly backup
- Retention: 30 days
- Test restore: Monthly

## Documentation

- ✅ README.md - Installation & usage
- ✅ QUICK_START.md - Quick setup guide
- ✅ API_DOCUMENTATION.md - API reference
- ✅ DEPLOYMENT.md - Deployment guide
- ✅ PROJECT_SUMMARY.md - Project overview
- ✅ Code comments - Inline documentation

## Support & Contact

### For Developers
- Review code comments
- Check API documentation
- Follow coding standards

### For Users
- User manual (to be created)
- Training materials (to be created)
- Support email/ticket system

### For Admins
- System administration guide
- Troubleshooting guide
- Backup/restore procedures

## License

Proprietary - Internal use only

## Credits

**Developed by:** [Your Team Name]
**Version:** 1.0.0
**Last Updated:** January 2024

---

## Quick Reference

### Default Credentials
- Username: `admin`
- Password: `admin123`

### URLs (Development)
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080/api
- Admin Panel: http://localhost:3000/admin/login

### Database
- Name: `perjalanan_dinas`
- Default user: `postgres`

### Important Commands

**Start Backend:**
```bash
cd backend && go run cmd/api/main.go
```

**Start Frontend:**
```bash
cd frontend && npm run dev
```

**Database Migration:**
```bash
psql -U postgres -d perjalanan_dinas -f migrations/001_initial_schema.sql
```

---

**For detailed information, please refer to the respective documentation files.**
