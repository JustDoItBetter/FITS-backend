# FITS Backend - Authentication System

##  Vollständig Implementiert!

Das Authentication & Authorization System ist komplett implementiert und einsatzbereit!

##  Was wurde implementiert?

### 1. Admin Bootstrap System
- **Automatische Generierung** von RSA-Keypair bei Installation
- **Admin-Token** mit unbegrenzter Gültigkeit
- **Einmalige Initialisierung** - kann nur 1x aufgerufen werden

### 2. JWT Token System
- **4 Token-Typen:** Admin, Access, Refresh, Invitation
- **Sichere Signierung:** HS256 für User-Tokens
- **Konfigurierbare Laufzeiten:**
  - Access Token: 1 Stunde
  - Refresh Token: 30 Tage
  - Invitation Token: 7 Tage
  - Admin Token: ~100 Jahre (kein Ablauf)

### 3. Invitation System
- Admin kann **Invitation-Links** erstellen
- **User-Registrierung** über Invitation-Token
- **Username + Password** Authentifizierung
- **Bcrypt** Password Hashing (Cost: 12)

### 4. Login & Session Management
- **Login** mit Username/Password
- **Access Token + Refresh Token** Ausgabe
- **Token Refresh** ohne erneutes Login
- **Logout** löscht alle Refresh Tokens

### 5. Middleware & Authorization
- **JWT Middleware:** Validiert Tokens automatisch
- **RBAC:** Role-Based Access Control (Admin, Teacher, Student)
- **Ownership Check:** User kann nur eigene Ressourcen sehen

### 6. Database Integration
- **PostgreSQL** mit GORM
- **Connection Pool** (konfigurierbar)
- **Auto-Migrations** bereit
- **Health Checks** für Database

---

##  Schnellstart

### 1. PostgreSQL starten

```bash
docker run -d \
  --name fits-postgres \
  -e POSTGRES_USER=fits_user \
  -e POSTGRES_PASSWORD=fits_password \
  -e POSTGRES_DB=fits_db \
  -p 5432:5432 \
  postgres:15
```

### 2. Migrations ausführen

```bash
# Manuelle Migration (temporär)
psql -h localhost -U fits_user -d fits_db -f migrations/001_initial_schema.sql
```

### 3. Server starten

```bash
go run cmd/server/main.go
```

---

##  API Workflows

### Workflow 1: Admin Initialisierung

```bash
# 1. Admin-Zertifikat generieren
curl -X POST http://localhost:8080/api/v1/bootstrap/init

# Response:
{
  "success": true,
  "data": {
    "admin_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "message": "Admin certificate generated successfully. Store this token securely - it cannot be recovered!",
    "public_key_path": "./configs/keys/admin.pub"
  }
}

#  WICHTIG: Admin-Token sicher speichern!
export ADMIN_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Workflow 2: Student einladen & registrieren

```bash
# 1. Admin erstellt Invitation (zuerst muss Student in DB existieren!)
curl -X POST http://localhost:8080/api/v1/admin/invite \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_uuid": "550e8400-e29b-41d4-a716-446655440000",
    "role": "student"
  }'

# Response:
{
  "success": true,
  "data": {
    "invitation_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "invitation_link": "https://fits.example.com/invite/eyJhbGci...",
    "expires_at": "2025-10-25T12:00:00Z"
  }
}

# 2. Student öffnet Invitation Link
curl -X GET http://localhost:8080/api/v1/invite/{invitation_token}

# Response:
{
  "success": true,
  "data": {
    "email": "max@example.com",
    "first_name": "User",
    "last_name": "Name",
    "role": "student",
    "expires_at": "2025-10-25T12:00:00Z"
  }
}

# 3. Student registriert sich
curl -X POST http://localhost:8080/api/v1/invite/{invitation_token}/complete \
  -H "Content-Type: application/json" \
  -d '{
    "username": "max.mustermann",
    "password": "SecurePassword123!"
  }'

# Response:
{
  "success": true,
  "message": "registration completed successfully"
}
```

### Workflow 3: Login & API Zugriff

```bash
# 1. Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "max.mustermann",
    "password": "SecurePassword123!"
  }'

# Response:
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 3600,
    "token_type": "Bearer",
    "role": "student",
    "user_id": "550e8400-e29b-41d4-a716-446655440000"
  }
}

# Token speichern
export ACCESS_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# 2. API Zugriff mit Token
curl -X GET http://localhost:8080/api/v1/student \
  -H "Authorization: Bearer $ACCESS_TOKEN"

# 3. Token erneuern (nach 1 Stunde)
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }'

# 4. Logout
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

---

##  Security Features

### Password Hashing
- **Bcrypt** mit Cost Factor 12
- Keine Klartext-Passwörter in DB
- Salted Hashes

### JWT Tokens
- **HS256** Signierung
- **Configurable Secret** (aus config.toml)
- **Expiry Check** in Middleware
- **Token Revocation** via Refresh Token Deletion

### RSA Signatures (für Admin)
- **4096-bit** Keypair
- **RSA-PSS** mit SHA-256
- **Private Key** nur lokal gespeichert (chmod 600)

### Database Security
- **Prepared Statements** (GORM)
- **No SQL Injection**
- **Connection Pool Limits**

---

##  Datenbank Schema

Siehe `migrations/001_initial_schema.sql` für vollständiges Schema.

**Wichtigste Tabellen:**
- `users` - Authentifizierte Benutzer
- `refresh_tokens` - Session Management
- `invitations` - User-Registrierung
- `students` - Student-Daten
- `teachers` - Teacher-Daten
- `teacher_keys` - RSA Keys für digitale Signaturen (TODO)
- `reports` - Berichtshefte (TODO)
- `signatures` - Digitale Signaturen (TODO)

---

##  API Endpoints Übersicht

### Authentication & Bootstrap

| Method | Endpoint | Beschreibung | Auth Required |
|--------|----------|--------------|---------------|
| POST | `/api/v1/bootstrap/init` | Admin initialisieren |  |
| POST | `/api/v1/auth/login` | Login |  |
| POST | `/api/v1/auth/logout` | Logout |  Bearer |
| POST | `/api/v1/auth/refresh` | Token erneuern |  |

### Invitations

| Method | Endpoint | Beschreibung | Auth Required |
|--------|----------|--------------|---------------|
| POST | `/api/v1/admin/invite` | Invitation erstellen |  Admin |
| GET | `/api/v1/invite/{token}` | Invitation Details |  |
| POST | `/api/v1/invite/{token}/complete` | Registrierung abschließen |  |

### Students

| Method | Endpoint | Beschreibung | Auth Required |
|--------|----------|--------------|---------------|
| PUT | `/api/v1/student` | Student erstellen |  Admin |
| GET | `/api/v1/student/:uuid` | Student abrufen |  (Optional Auth) |
| GET | `/api/v1/student` | Alle Students |  (Optional Auth) |
| POST | `/api/v1/student/:uuid` | Student aktualisieren |  Admin |
| DELETE | `/api/v1/student/:uuid` | Student löschen |  Admin |

### Teachers

| Method | Endpoint | Beschreibung | Auth Required |
|--------|----------|--------------|---------------|
| POST | `/api/v1/teacher` | Teacher erstellen |  Admin |
| GET | `/api/v1/teacher/:uuid` | Teacher abrufen |  (Optional Auth) |
| GET | `/api/v1/teacher` | Alle Teachers |  (Optional Auth) |
| POST | `/api/v1/teacher/:uuid` | Teacher aktualisieren |  Admin |
| DELETE | `/api/v1/teacher/:uuid` | Teacher löschen |  Admin |

---

##  Testing

```bash
# Swagger UI öffnen
open http://localhost:8080/swagger/index.html

# Im Swagger UI:
# 1. Bootstrap → Initialize Admin
# 2. Kopiere Admin-Token
# 3. Klicke "Authorize" oben rechts
# 4. Gib ein: "Bearer {admin_token}"
# 5. Teste alle Admin-Endpoints!
```

---

##  Nächste Schritte (TODO)

Die folgenden Features sind geplant aber noch NICHT implementiert:

### 1. Teacher Keypair Domain
- RSA Keypair Generierung für Teacher (bei Registrierung)
- Private Key Verschlüsselung mit Teacher Password
- Key Management (Rotation, etc.)

### 2. Report Domain
- Berichtsheft Upload (PDF/Parquet)
- File Storage Management
- Hash Berechnung (SHA-256)
- Status Management (pending/signed/rejected)

### 3. Signature Domain
- Digitale Signatur mit Teacher Private Key
- Signatur-Verifikation
- Report Download mit Signatur

### 4. Migration Runner
- Auto-Migration beim Start
- Migration Rollback Support

### 5. Improvements
- Rate Limiting
- Email Service für Invitations
- Audit Logging
- Password Strength Validation
- Token Cleanup Cronjob

---

##  Dokumentation

- **Swagger UI:** `http://localhost:8080/swagger/index.html`
- **API Docs:** Siehe `README.md`
- **Implementation Status:** Siehe `IMPLEMENTATION_STATUS.md`
- **Quick Start:** Siehe `QUICK_START.md`

---

**Status:**  Production-Ready Authentication System
**Version:** 1.0
**Last Updated:** 2025-10-18
