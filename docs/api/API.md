# FITS Backend - API Dokumentation

## Interaktive Dokumentation

Die vollständige, interaktive API-Dokumentation ist verfügbar über Swagger UI:

**URL:** http://localhost:8080/swagger/index.html

## API-Übersicht

### Base URL

```
http://localhost:8080
```

### Authentication

Alle geschützten Endpunkte verwenden Bearer Token Authentication:

```http
Authorization: Bearer <access_token>
```

## Endpunkte

###  Public Endpoints

#### Health Check
```http
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "database": "connected",
  "time": "2025-10-18T12:00:00Z"
}
```

---

###  Bootstrap (Einmalig)

#### Admin Initialisierung
```http
POST /api/v1/bootstrap/init
```

Erstellt den ersten Admin-Account und generiert RSA-Schlüssel.

**Response:**
```json
{
  "success": true,
  "data": {
    "admin_token": "eyJhbGc...",
    "message": "Admin certificate generated successfully",
    "public_key_path": "./configs/keys/admin.pub"
  }
}
```

 **Kann nur einmal aufgerufen werden!**

---

###  Authentication

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json
```

**Request:**
```json
{
  "username": "max.mustermann",
  "password": "SecurePass123!"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "expires_in": 3600,
    "token_type": "Bearer",
    "role": "student",
    "user_id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

#### Token Refresh
```http
POST /api/v1/auth/refresh
Content-Type: application/json
```

**Request:**
```json
{
  "refresh_token": "eyJhbGc..."
}
```

#### Logout
```http
POST /api/v1/auth/logout
Authorization: Bearer <access_token>
```

Invalidiert alle Refresh-Tokens des Benutzers.

---

###  Invitations

#### Einladung erstellen (Admin only)
```http
POST /api/v1/admin/invite
Authorization: Bearer <admin_token>
Content-Type: application/json
```

**Request:**
```json
{
  "email": "student@example.com",
  "first_name": "Max",
  "last_name": "Mustermann",
  "role": "student",
  "department": "IT"  // Nur für teacher
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "invitation_token": "eyJhbGc...",
    "invitation_link": "https://fits.example.com/invite/eyJhbGc...",
    "expires_at": "2025-10-25T12:00:00Z"
  }
}
```

#### Einladung abrufen
```http
GET /api/v1/invite/:token
```

**Response:**
```json
{
  "success": true,
  "data": {
    "email": "student@example.com",
    "first_name": "Max",
    "last_name": "Mustermann",
    "role": "student",
    "expires_at": "2025-10-25T12:00:00Z"
  }
}
```

#### Registrierung abschließen
```http
POST /api/v1/invite/:token/complete
Content-Type: application/json
```

**Request:**
```json
{
  "username": "max.mustermann",
  "password": "SecurePass123!"
}
```

---

### ‍ Students

#### Student erstellen (Admin only)
```http
PUT /api/v1/student
Authorization: Bearer <admin_token>
Content-Type: application/json
```

#### Student abrufen
```http
GET /api/v1/student/:uuid
Authorization: Bearer <token> (optional)
```

#### Student aktualisieren (Admin only)
```http
POST /api/v1/student/:uuid
Authorization: Bearer <admin_token>
Content-Type: application/json
```

#### Student löschen (Admin only)
```http
DELETE /api/v1/student/:uuid
Authorization: Bearer <admin_token>
```

#### Alle Studenten auflisten
```http
GET /api/v1/student
Authorization: Bearer <token> (optional)
```

---

### ‍ Teachers

Analog zu Students mit `/api/v1/teacher` Endpunkten.

---

###  System

#### Metrics (Protected)
```http
GET /metrics?secret=<metrics_secret>
```

Prometheus-kompatible Metriken.

---

## Error Responses

### Standard Error Format

```json
{
  "success": false,
  "error": "Error message",
  "details": "Detailed error information"
}
```

### HTTP Status Codes

| Code | Bedeutung |
|------|-----------|
| 200 | OK - Erfolgreiche Anfrage |
| 201 | Created - Ressource erstellt |
| 204 | No Content - Erfolgreich, keine Daten |
| 400 | Bad Request - Ungültige Anfrage |
| 401 | Unauthorized - Authentifizierung erforderlich |
| 403 | Forbidden - Keine Berechtigung |
| 404 | Not Found - Ressource nicht gefunden |
| 409 | Conflict - Ressourcen-Konflikt |
| 422 | Unprocessable Entity - Validierungsfehler |
| 500 | Internal Server Error - Server-Fehler |

---

## Rate Limiting

**Status:** Noch nicht implementiert

**Geplant:**
- 100 Requests pro Minute für authentifizierte Benutzer
- 20 Requests pro Minute für nicht-authentifizierte Anfragen

---

## Pagination

**Status:** Noch nicht implementiert

**Geplant:**
```http
GET /api/v1/student?page=1&limit=20
```

---
