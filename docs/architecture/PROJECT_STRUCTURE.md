# FITS Backend - Projekt-Struktur

## Übersicht

```
FITS-backend/
├── cmd/                        # Anwendungs-Entry Points
│   └── server/
│       └── main.go            # Haupt-Server Entry Point
│
├── configs/                    # Konfigurationsdateien
│   ├── config.toml            # Haupt-Konfiguration
│   └── keys/                  # RSA-Schlüssel (generiert)
│
├── docs/                       # Dokumentation
│   ├── API.md                 # API-Dokumentation
│   ├── ARCHITECTURE.md        # System-Architektur
│   ├── DEPLOYMENT.md          # Deployment-Guide
│   ├── TESTING.md             # Test-Dokumentation
│   ├── PROJECT_STRUCTURE.md   # Diese Datei
│   ├── docs.go                # Swagger Code-Gen
│   ├── swagger.json           # OpenAPI Spec (JSON)
│   └── swagger.yaml           # OpenAPI Spec (YAML)
│
├── internal/                   # Private Application Code
│   ├── common/                # Shared Code
│   │   ├── errors/            # Error Types & Handling
│   │   └── response/          # HTTP Response Wrapper
│   │
│   ├── config/                # Configuration Loader
│   │   └── config.go
│   │
│   ├── domain/                # Domain Logic (DDD)
│   │   ├── auth/              # Authentication Domain  TESTED
│   │   │   ├── model.go       # Data Models & DTOs
│   │   │   ├── repository.go  # Data Access Interface
│   │   │   ├── auth_service.go         # Login, Logout, Refresh
│   │   │   ├── auth_service_test.go    # Unit Tests
│   │   │   ├── bootstrap_service.go    # Admin Init
│   │   │   ├── invitation_service.go   # User Invitations
│   │   │   └── handler.go              # HTTP Handlers
│   │   │
│   │   ├── signing/           # Digital Signatures
│   │   ├── student/           # Student Management
│   │   └── teacher/           # Teacher Management
│   │
│   └── middleware/            # HTTP Middleware  TESTED
│       ├── jwt.go             # JWT Authentication
│       ├── jwt_test.go        # JWT Tests
│       ├── rbac.go            # Role-Based Access Control
│       └── rbac_test.go       # RBAC Tests
│
├── migrations/                 # SQL Migrationen (Legacy)
│   └── 001_initial_schema.sql # Initial Schema (für Referenz)
│
├── pkg/                        # Public Libraries
│   ├── crypto/                # Cryptography Package  TESTED
│   │   ├── jwt.go             # JWT Implementation
│   │   ├── jwt_test.go        # JWT Tests
│   │   ├── password.go        # Bcrypt Password Hashing
│   │   ├── password_test.go   # Password Tests
│   │   ├── rsa.go             # RSA Signatures
│   │   └── rsa_test.go        # RSA Tests
│   │
│   └── database/              # Database Layer
│       ├── database.go        # Connection & Auto-Setup
│       └── migrations.go      # Go-based Migrations
│
├── scripts/                    # Utility Scripts
│   ├── README.md              # Scripts Dokumentation
│   └── test_auth_flow.sh      # E2E Auth Test
│
├── uploads/                    # File Uploads (gitignored)
│
├── .gitignore                 # Git Ignore Rules
├── go.mod                     # Go Module Definition
├── go.sum                     # Go Module Checksums
├── coverage.out               # Test Coverage Data
└── README.md                  # Projekt Haupt-README
```

## Verzeichnis-Details

### `/cmd`
**Zweck:** Anwendungs-Entry Points

Enthält ausführbare Programme. Jedes Programm hat sein eigenes Unterverzeichnis.

**Konvention:** Package `main`, minimale Logik (nur Initialisierung und Routing)

### `/configs`
**Zweck:** Konfigurationsdateien

- `config.toml` - Haupt-Konfiguration (Server, DB, JWT, etc.)
- `keys/` - RSA-Schlüssel (automatisch generiert beim Bootstrap)

**Wichtig:**
- Niemals Secrets in Git committen!
- Verwende `*.local.toml` für lokale Overrides

### `/docs`
**Zweck:** Projekt-Dokumentation

| Datei | Beschreibung |
|-------|--------------|
| API.md | API-Endpunkte und Beispiele |
| ARCHITECTURE.md | System-Architektur und Design-Entscheidungen |
| DEPLOYMENT.md | Production Deployment Guide |
| TESTING.md | Test-Strategie und Coverage |
| PROJECT_STRUCTURE.md | Diese Datei |

### `/internal`
**Zweck:** Private Application Code (nicht importierbar von außen)

#### `/internal/common`
Shared Code für alle Domains:
- `errors/` - Custom Error Types
- `response/` - Standardisierte HTTP Responses

#### `/internal/config`
Configuration Loading und Validation

#### `/internal/domain`
**Domain-Driven Design (DDD)**

Jede Domain ist eigenständig und folgt dem 4-Layer Pattern:

1. **Model Layer** - Entities, Value Objects, DTOs
2. **Repository Layer** - Data Access Interface
3. **Service Layer** - Business Logic
4. **Handler Layer** - HTTP Endpoints

**Domains:**
- `auth/` - Authentifizierung & Autorisierung  **Komplett**
- `signing/` - Digitale Signaturen  **In Progress**
- `student/` - Studenten-Verwaltung  **In Progress**
- `teacher/` - Lehrer-Verwaltung  **In Progress**

#### `/internal/middleware`
HTTP Middleware (Fiber Framework):
- `jwt.go` - JWT Token Validation
- `rbac.go` - Role-Based Access Control

### `/migrations`
**Zweck:** Legacy SQL Migrations (nur für Referenz)

 **Nicht mehr aktiv verwendet!**

Migrationen laufen jetzt automatisch über Go-Code in `pkg/database/migrations.go`

### `/pkg`
**Zweck:** Public Libraries (kann von anderen Projekten importiert werden)

#### `/pkg/crypto`
Kryptografie-Funktionen:
- JWT Token Generation & Validation
- Bcrypt Password Hashing
- RSA 4096-bit Signatures

#### `/pkg/database`
Database Layer:
- Connection Management
- Auto-Initialization
- Go-based Migrations
- Schema Verification

### `/scripts`
**Zweck:** Utility Scripts

- `test_auth_flow.sh` - End-to-End Auth Testing
- `README.md` - Scripts Dokumentation

### `/uploads`
**Zweck:** File Upload Speicher

- Gitignored
- Für Report-PDFs und andere Uploads
- Konfigurierbar in `config.toml`

## Code-Organisation Prinzipien

### 1. Clean Architecture
- Dependency Inversion (Interfaces)
- Domain-zentrisch
- Framework-unabhängige Business Logic

### 2. Domain-Driven Design (DDD)
- Bounded Contexts (Domains)
- Entities & Value Objects
- Repository Pattern
- Service Layer

### 3. Go Best Practices
- `internal/` für private code
- `pkg/` für reusable libraries
- `cmd/` für executables
- Flat package structure

### 4. Test-Organisation
- Tests neben dem Code (`*_test.go`)
- Mock-Implementierungen im gleichen Package
- Table-Driven Tests
- Benchmark Tests für Performance-kritische Funktionen

## Naming Conventions

### Packages
- Lowercase, single word
- Beschreibend: `auth`, `crypto`, `database`

### Files
- Lowercase mit Underscores
- `<topic>_<type>.go`: `auth_service.go`, `jwt_test.go`

### Interfaces
- `-er` Suffix: `Repository`, `Handler`
- Im Package des Consumers definiert

### Tests
- `Test<Function>` für Unit Tests
- `Benchmark<Function>` für Benchmarks
- `Example<Function>` für Examples

## Import Paths

```go
import (
    "github.com/JustDoItBetter/FITS-backend/internal/domain/auth"
    "github.com/JustDoItBetter/FITS-backend/pkg/crypto"
    "github.com/JustDoItBetter/FITS-backend/pkg/database"
)
```

## Build & Run

```bash
# Development
go run cmd/server/main.go

# Build
go build -o bin/fits-server cmd/server/main.go

# Tests
go test ./...

# Coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Hinzufügen neuer Features

### Neue Domain hinzufügen

1. Erstelle Verzeichnis: `internal/domain/<name>/`
2. Erstelle Files:
   - `model.go` - Entities & DTOs
   - `repository.go` - Interface
   - `<name>_service.go` - Business Logic
   - `handler.go` - HTTP Handlers
3. Schreibe Tests: `*_test.go`
4. Registriere Routes in `cmd/server/main.go`
5. Update Swagger Docs

### Neue Middleware hinzufügen

1. Erstelle File: `internal/middleware/<name>.go`
2. Implementiere `fiber.Handler` Funktion
3. Schreibe Tests: `internal/middleware/<name>_test.go`
4. Verwende in `cmd/server/main.go`

### Neue Migration hinzufügen

1. Editiere: `pkg/database/migrations.go`
2. Füge zu `getAllMigrations()` hinzu:
```go
{
    Version: "002",
    Name:    "add_new_table",
    Up:      migration002AddNewTable,
}
```
3. Implementiere Migrations-Funktion
4. Server startet und führt Migration automatisch aus

## Umgebungsvariablen

```bash
# Optional - überschreibt config.toml
DB_HOST=localhost
DB_PORT=5432
DB_USER=fits_user
DB_PASSWORD=secret
DB_NAME=fits_db
JWT_SECRET=super-secret-key
ENVIRONMENT=production
```
