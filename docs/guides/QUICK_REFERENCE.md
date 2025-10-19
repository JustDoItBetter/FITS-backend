# FITS Backend - Quick Reference

##  Schnellstart

```bash
make docker-up     # PostgreSQL starten
make run           # Server starten (Terminal 1)
make bootstrap     # Admin initialisieren (Terminal 2)
```

##  Häufig verwendete Befehle

### Entwicklung
```bash
make run               # Server starten
make test              # Alle Tests
make test-coverage     # Tests + Coverage
make e2e-test          # End-to-End Tests
make clean             # Aufräumen
```

### Datenbank
```bash
make docker-up         # PostgreSQL starten
make docker-down       # PostgreSQL stoppen
make reset-db          # DB zurücksetzen (Daten löschen)
```

### Testing
```bash
# Unit Tests
go test ./...                           # Alle Tests
go test ./pkg/crypto -v                 # Crypto Tests (verbose)
go test ./internal/middleware -cover    # Middleware mit Coverage

# Benchmarks
go test ./pkg/crypto -bench=. -benchmem

# E2E Test
./scripts/test_full_flow.sh
```

### Build
```bash
make build             # Binary erstellen (bin/fits-server)
./bin/fits-server      # Binary ausführen
```

##  API Quick Reference

### Admin Bootstrap
```bash
curl -X POST http://localhost:8080/api/v1/bootstrap/init | jq .
# Speichere: export ADMIN_TOKEN="..."
```

### Student einladen
```bash
curl -X POST http://localhost:8080/api/v1/admin/invite \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student@example.com",
    "first_name": "Max",
    "last_name": "Mustermann",
    "role": "student"
  }' | jq .
```

### Teacher einladen
```bash
curl -X POST http://localhost:8080/api/v1/admin/invite \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "teacher@example.com",
    "first_name": "Anna",
    "last_name": "Schmidt",
    "role": "teacher",
    "department": "Informatik"
  }' | jq .
```

### Registrieren
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "invitation_token": "TOKEN_HIER",
    "username": "max.mustermann",
    "password": "SecurePass123!"
  }' | jq .
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "max.mustermann",
    "password": "SecurePass123!"
  }' | jq .
# Speichere: export USER_TOKEN="..."
```

### Token Refresh
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "REFRESH_TOKEN_HIER"
  }' | jq .
```

### Logout
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer $USER_TOKEN" | jq .
```

##  Wichtige Dateien

```
├── Makefile                      # Haupt-Befehle
├── docker-compose.yml            # PostgreSQL Setup
├── .env.example                  # Umgebungsvariablen
├── configs/config.toml           # Konfiguration
│
├── TESTING_GUIDE.md              # Ausführlicher Testing-Guide
├── CHANGELOG.md                  # Alle Änderungen
├── README.md                     # Hauptdokumentation
│
├── docs/
│   ├── API.md                    # API-Dokumentation
│   ├── TESTING.md                # Test-Dokumentation
│   ├── DEPLOYMENT.md             # Deployment-Guide
│   └── ARCHITECTURE.md           # Architektur
│
└── scripts/
    ├── test_full_flow.sh         # E2E-Test
    ├── reset.sh                  # Komplett-Reset
    └── reset_db_only.sh          # DB-Reset
```

##  Troubleshooting

### Server startet nicht
```bash
# PostgreSQL prüfen
pg_isready -h localhost -p 5432

# PostgreSQL neu starten
make docker-down
make docker-up
```

### "Admin already initialized"
```bash
make reset-db      # DB zurücksetzen
make bootstrap     # Neu initialisieren
```

### Tests schlagen fehl
```bash
make clean         # Aufräumen
go mod tidy        # Dependencies aktualisieren
make test          # Tests erneut ausführen
```

### Port 8080 bereits belegt
```bash
# Anderen Prozess finden und beenden
lsof -ti:8080 | xargs kill -9
```

##  Test Coverage

```
pkg/crypto              73.7%  
internal/middleware     88.2%  
internal/domain/auth    18.1%  
Gesamt                  20.8%
```

##  Token-Gültigkeitszeiten

| Token-Typ | Gültigkeit |
|-----------|------------|
| Admin | ~100 Jahre |
| Access | 24 Stunden |
| Refresh | 7 Tage |
| Invitation | 7 Tage |

##  Rollen-Hierarchie

```
Admin
  ├─ Hat Zugriff auf ALLES
  │
Teacher
  ├─ Hat Zugriff auf Teacher-Endpunkte
  ├─ Hat Zugriff auf Student-Endpunkte
  │
Student
  └─ Hat nur Zugriff auf Student-Endpunkte
```

##  Weitere Ressourcen

- **Ausführlicher Testing-Guide:** [TESTING_GUIDE.md](TESTING_GUIDE.md)
- **API-Dokumentation:** [docs/API.md](docs/API.md)
- **Deployment:** [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md)
- **Architektur:** [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)

---

**Tipp:** Nutze `make help` für alle verfügbaren Befehle!
