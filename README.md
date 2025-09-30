# FITS-api

## API Endpoints

### Health Check
```http
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "time": "2025-09-30T12:00:00Z"
}
```

---

### Prometheus Metrics
```http
GET /metrics
Authorization: Bearer <metrics_secret>
```

**Response:** Prometheus-Format
```
# HELP go_goroutines Number of goroutines
# TYPE go_goroutines gauge
go_goroutines 10
...
```

---

### Upload Parquet File
```http
POST /api/v1/upload
Content-Type: multipart/form-data

file: <parquet-file>
```

**Status:** 🚧 Not Implemented

---

### Get Sign Requests
```http
GET /api/v1/sign_requests
```

**Response:** Parquet file mit pending sign requests

**Status:** 🚧 Not Implemented

---

### Upload Signed Requests
```http
POST /api/v1/sign_uploads
Content-Type: multipart/form-data

file: <parquet-file>
```

**Status:** 🚧 Not Implemented

---

### Student Management

#### Register Student
```http
PUT /api/v1/student
Authorization: Bearer <registration_secret>
Content-Type: text/plain

uuid = "550e8400-e29b-41d4-a716-446655440000"
first_name = "Max"
last_name = "Mustermann"
email = "max@example.com"
teacher_id = "teacher-uuid"
```

**Status:** 🚧 Not Implemented

#### Update Student
```http
POST /api/v1/student
Authorization: Bearer <update_secret>
Content-Type: text/plain

uuid = "550e8400-e29b-41d4-a716-446655440000"
first_name = "Moritz"
```

**Status:** 🚧 Not Implemented

#### Delete Student
```http
DELETE /api/v1/student/:uuid
Authorization: Bearer <deletion_secret>
```

**Status:** 🚧 Not Implemented

---

### Teacher Management

#### Register Teacher
```http
POST /api/v1/teacher
Authorization: Bearer <registration_secret>
Content-Type: text/plain

uuid = "teacher-uuid"
first_name = "Anna"
last_name = "Schmidt"
email = "anna@example.com"
department = "Computer Science"
```

**Status:**  Not Implemented

#### Update Teacher
```http
POST /api/v1/teacher/update
Content-Type: text/plain
```

**Status:**  Not Implemented

#### Delete Teacher
```http
DELETE /api/v1/teacher/:uuid
```

**Status:**  Not Implemented
