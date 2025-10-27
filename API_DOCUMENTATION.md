# API Documentation - Sistem Perjalanan Dinas

Base URL: `http://localhost:8080/api`

## Authentication

Admin endpoints memerlukan JWT token di header:
```
Authorization: Bearer <token>
```

## Endpoints

### 1. Authentication

#### Login Admin
```http
POST /api/auth/login
```

**Request Body:**
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "username": "admin",
  "message": "Login successful"
}
```

**Response (401 Unauthorized):**
```json
{
  "error": "Invalid credentials"
}
```

---

### 2. Employees (Karyawan)

#### Get All Employees (Public)
```http
GET /api/employees
```

**Response (200 OK):**
```json
{
  "employees": [
    {
      "id": 1,
      "nip": "199001012015011001",
      "name": "Ahmad Budiman",
      "position": "Direktur Utama",
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    }
  ]
}
```

#### Get Employee by ID (Public)
```http
GET /api/employees/:id
```

**Response (200 OK):**
```json
{
  "employee": {
    "id": 1,
    "nip": "199001012015011001",
    "name": "Ahmad Budiman",
    "position": "Direktur Utama",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

#### Create Employee (Admin Only)
```http
POST /api/admin/employees
```

**Request Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "nip": "199506066015011006",
  "name": "John Doe",
  "position": "Manager"
}
```

**Response (201 Created):**
```json
{
  "message": "Employee created successfully",
  "employee": {
    "id": 6,
    "nip": "199506066015011006",
    "name": "John Doe",
    "position": "Manager",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

#### Update Employee (Admin Only)
```http
PUT /api/admin/employees/:id
```

**Request Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "nip": "199506066015011006",
  "name": "John Doe Updated",
  "position": "Senior Manager"
}
```

**Response (200 OK):**
```json
{
  "message": "Employee updated successfully",
  "employee": {
    "id": 6,
    "nip": "199506066015011006",
    "name": "John Doe Updated",
    "position": "Senior Manager",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:30:00Z"
  }
}
```

#### Delete Employee (Admin Only)
```http
DELETE /api/admin/employees/:id
```

**Request Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "message": "Employee deleted successfully"
}
```

---

### 3. Position Codes (Pengkodean Jabatan)

#### Get All Position Codes (Public)
```http
GET /api/position-codes
```

**Response (200 OK):**
```json
{
  "position_codes": [
    {
      "id": 1,
      "position": "Direktur Utama",
      "code": "DIRUT",
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    }
  ]
}
```

#### Create Position Code (Admin Only)
```http
POST /api/admin/position-codes
```

**Request Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "position": "Manager",
  "code": "MNG"
}
```

**Response (201 Created):**
```json
{
  "message": "Position code created successfully",
  "position_code": {
    "id": 7,
    "position": "Manager",
    "code": "MNG",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

#### Update Position Code (Admin Only)
```http
PUT /api/admin/position-codes/:id
```

**Request Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "position": "Manager",
  "code": "MGR"
}
```

**Response (200 OK):**
```json
{
  "message": "Position code updated successfully",
  "position_code": {
    "id": 7,
    "position": "Manager",
    "code": "MGR",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:30:00Z"
  }
}
```

#### Delete Position Code (Admin Only)
```http
DELETE /api/admin/position-codes/:id
```

**Request Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "message": "Position code deleted successfully"
}
```

---

### 4. Travel Requests (Perjalanan Dinas)

#### Create Travel Request (Public)
```http
POST /api/travel-requests
```

**Request Body:**
```json
{
  "employee_id": 1,
  "purpose": "Kunjungan kerja ke kantor cabang Jakarta",
  "departure_place": "Surabaya",
  "destination": "Jakarta",
  "departure_date": "2024-02-01",
  "return_date": "2024-02-03",
  "transportation": "Pesawat"
}
```

**Response (201 Created):**
```json
{
  "message": "Travel request created successfully",
  "travel_request": {
    "id": 1,
    "employee_id": 1,
    "employee": {
      "id": 1,
      "nip": "199001012015011001",
      "name": "Ahmad Budiman",
      "position": "Direktur Utama"
    },
    "purpose": "Kunjungan kerja ke kantor cabang Jakarta",
    "departure_place": "Surabaya",
    "destination": "Jakarta",
    "departure_date": "2024-02-01T00:00:00Z",
    "return_date": "2024-02-03T00:00:00Z",
    "duration_days": 3,
    "transportation": "Pesawat",
    "request_number": "064/0001/DIB/DIRUT/NOTA",
    "report_number": "064/    /DIB/DIRUT/NOTA",
    "status": "pending",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

**Validation:**
- `employee_id`: required
- `purpose`: required
- `destination`: required
- `departure_date`: required (format: YYYY-MM-DD)
- `return_date`: required (format: YYYY-MM-DD, must be >= departure_date)
- `transportation`: required (options: "Angkutan Umum", "Pesawat", "Kereta Api")
- `departure_place`: optional (default: "Surabaya")

#### Get Travel Request by ID (Public)
```http
GET /api/travel-requests/:id
```

**Response (200 OK):**
```json
{
  "travel_request": {
    "id": 1,
    "employee_id": 1,
    "employee": { ... },
    "purpose": "Kunjungan kerja ke kantor cabang Jakarta",
    "departure_place": "Surabaya",
    "destination": "Jakarta",
    "departure_date": "2024-02-01T00:00:00Z",
    "return_date": "2024-02-03T00:00:00Z",
    "duration_days": 3,
    "transportation": "Pesawat",
    "request_number": "064/0001/DIB/DIRUT/NOTA",
    "report_number": "064/    /DIB/DIRUT/NOTA",
    "status": "pending",
    "travel_report": null,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

#### Get All Travel Requests (Admin Only)
```http
GET /api/admin/travel-requests
```

**Request Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "travel_requests": [
    {
      "id": 1,
      "employee": { ... },
      "request_number": "064/0001/DIB/DIRUT/NOTA",
      "status": "pending",
      ...
    }
  ]
}
```

#### Delete Travel Request (Admin Only)
```http
DELETE /api/admin/travel-requests/:id
```

**Request Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "message": "Travel request deleted successfully"
}
```

---

### 5. Travel Reports (Berita Acara)

#### Create Travel Report (Admin Only)
```http
POST /api/admin/travel-reports
```

**Request Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "travel_request_id": 1,
  "representative_name": "Ahmad Budiman",
  "representative_position": "Direktur Utama",
  "visit_proofs": [
    {
      "date": "2024-02-01",
      "depart_from": "Surabaya",
      "stay_or_stop_at": "Hotel Jakarta",
      "arrive_at": "Jakarta",
      "signature_proof": "Stempel & TTD Kantor Cabang Jakarta"
    },
    {
      "date": "2024-02-02",
      "depart_from": "Jakarta",
      "arrive_at": "Bandung",
      "signature_proof": "Stempel & TTD Kantor Cabang Bandung"
    }
  ]
}
```

**Response (201 Created):**
```json
{
  "message": "Travel report created successfully",
  "travel_report": {
    "id": 1,
    "travel_request_id": 1,
    "report_number": "064/    /DIB/DIRUT/NOTA",
    "representative_name": "Ahmad Budiman",
    "representative_position": "Direktur Utama",
    "visit_proofs": [
      {
        "id": 1,
        "travel_report_id": 1,
        "date": "2024-02-01T00:00:00Z",
        "depart_from": "Surabaya",
        "stay_or_stop_at": "Hotel Jakarta",
        "arrive_at": "Jakarta",
        "signature_proof": "Stempel & TTD Kantor Cabang Jakarta",
        "created_at": "2024-01-01T10:00:00Z",
        "updated_at": "2024-01-01T10:00:00Z"
      }
    ],
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

#### Get Travel Report by Request ID (Admin Only)
```http
GET /api/admin/travel-reports/:request_id
```

**Request Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "travel_report": {
    "id": 1,
    "travel_request_id": 1,
    "report_number": "064/    /DIB/DIRUT/NOTA",
    "representative_name": "Ahmad Budiman",
    "representative_position": "Direktur Utama",
    "visit_proofs": [ ... ],
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

---

### 6. PDF Downloads (Public)

#### Download Nota Permintaan PDF
```http
GET /api/pdf/nota-permintaan/:id
```

**Response:**
- Content-Type: `application/pdf`
- Content-Disposition: `attachment; filename=nota_permintaan_064-0001-DIB-DIRUT-NOTA.pdf`

#### Download Berita Acara PDF
```http
GET /api/pdf/berita-acara/:id
```

**Response:**
- Content-Type: `application/pdf`
- Content-Disposition: `attachment; filename=berita_acara_064----DIB-DIRUT-NOTA.pdf`

**Note:** Berita Acara hanya tersedia jika travel report sudah dibuat.

#### Download Combined PDF (Both Documents)
```http
GET /api/pdf/combined/:id
```

**Response:**
- Content-Type: `application/pdf`
- Content-Disposition: `attachment; filename=perjalanan_dinas_064-0001-DIB-DIRUT-NOTA.pdf`

**Note:** Combined PDF hanya tersedia jika travel report sudah dibuat.

---

## Error Responses

### 400 Bad Request
```json
{
  "error": "Validation error message"
}
```

### 401 Unauthorized
```json
{
  "error": "Invalid or expired token"
}
```

### 404 Not Found
```json
{
  "error": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error"
}
```

---

## Request Examples (cURL)

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### Create Employee
```bash
curl -X POST http://localhost:8080/api/admin/employees \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "nip": "199506066015011006",
    "name": "John Doe",
    "position": "Manager"
  }'
```

### Create Travel Request
```bash
curl -X POST http://localhost:8080/api/travel-requests \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": 1,
    "purpose": "Kunjungan kerja",
    "destination": "Jakarta",
    "departure_date": "2024-02-01",
    "return_date": "2024-02-03",
    "transportation": "Pesawat"
  }'
```

### Download PDF
```bash
curl -O -J http://localhost:8080/api/pdf/nota-permintaan/1
```

---

## Postman Collection

Import the following JSON into Postman for easy testing:

```json
{
  "info": {
    "name": "Perjalanan Dinas API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "variable": [
    {
      "key": "baseUrl",
      "value": "http://localhost:8080/api"
    },
    {
      "key": "token",
      "value": ""
    }
  ]
}
```

---

## Rate Limiting

Currently, there is no rate limiting implemented. For production, consider adding rate limiting middleware.

## Versioning

API Version: v1 (current)

Future versions will be prefixed: `/api/v2/...`

---

## Support

For API issues or questions, contact the development team.
