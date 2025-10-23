# FITS Backend - Deployment Guide

## Production Deployment

### 1. Vorbereitung

#### Systemanforderungen

- **OS:** Linux (Ubuntu 20.04+ empfohlen)
- **Go:** 1.21 oder höher
- **PostgreSQL:** 15 oder höher
- **RAM:** Mindestens 2GB
- **Disk:** Mindestens 10GB

#### Sicherheits-Checkliste

Vor dem Deployment **MÜSSEN** folgende Änderungen vorgenommen werden:

- [ ] JWT Secret ändern (min. 32 Zeichen, kryptografisch sicher)
- [ ] Database Passwort ändern
- [ ] SSL/TLS für Database aktivieren (`ssl_mode = "require"`)
- [ ] Alle Secrets in `configs/config.toml` ändern
- [ ] HTTPS erzwingen (Reverse Proxy)
- [ ] CSP Headers konfigurieren
- [ ] CORS richtig konfigurieren

### 2. Binary erstellen

```bash
# Production Build mit Optimierungen
go build -ldflags="-s -w" -o fits-server cmd/server/main.go

# Für spezifische Plattform
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o fits-server cmd/server/main.go
```

### 3. Server Setup

#### User erstellen

```bash
sudo useradd -r -s /bin/false fits
sudo mkdir -p /opt/fits-backend
sudo chown fits:fits /opt/fits-backend
```

#### Dateien kopieren

```bash
sudo cp fits-server /opt/fits-backend/
sudo cp -r configs /opt/fits-backend/
sudo chown -R fits:fits /opt/fits-backend
sudo chmod 755 /opt/fits-backend/fits-server
```

### 4. PostgreSQL Setup

```bash
# PostgreSQL installieren
sudo apt update
sudo apt install postgresql postgresql-contrib

# Datenbank-User erstellen
sudo -u postgres psql
```

```sql
CREATE USER fits_prod WITH PASSWORD 'sehr-sicheres-passwort';
CREATE DATABASE fits_prod OWNER fits_prod;
GRANT ALL PRIVILEGES ON DATABASE fits_prod TO fits_prod;
\q
```

#### SSL/TLS aktivieren

```bash
# SSL Zertifikat generieren
sudo -u postgres openssl req -new -x509 -days 365 -nodes \
  -text -out /etc/postgresql/15/main/server.crt \
  -keyout /etc/postgresql/15/main/server.key

sudo chmod og-rwx /etc/postgresql/15/main/server.key
```

In `/etc/postgresql/15/main/postgresql.conf`:
```conf
ssl = on
ssl_cert_file = '/etc/postgresql/15/main/server.crt'
ssl_key_file = '/etc/postgresql/15/main/server.key'
```

### 5. Systemd Service

```bash
sudo nano /etc/systemd/system/fits-backend.service
```

```ini
[Unit]
Description=FITS Backend API Server
After=network.target postgresql.service
Requires=postgresql.service

[Service]
Type=simple
User=fits
Group=fits
WorkingDirectory=/opt/fits-backend
ExecStart=/opt/fits-backend/fits-server
Restart=always
RestartSec=10

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/fits-backend/uploads

# Environment
Environment="ENVIRONMENT=production"

[Install]
WantedBy=multi-user.target
```

```bash
# Service aktivieren und starten
sudo systemctl daemon-reload
sudo systemctl enable fits-backend
sudo systemctl start fits-backend
sudo systemctl status fits-backend
```

### 6. Reverse Proxy (Nginx)

```bash
sudo apt install nginx
sudo nano /etc/nginx/sites-available/fits-backend
```

```nginx
upstream fits_backend {
    server 127.0.0.1:8080;
}

server {
    listen 80;
    server_name api.fits.example.com;

    # Redirect HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.fits.example.com;

    # SSL Konfiguration
    ssl_certificate /etc/letsencrypt/live/api.fits.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.fits.example.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # Security Headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # Client max body size (für File Uploads)
    client_max_body_size 100M;

    location / {
        proxy_pass http://fits_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;

        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Health Check Endpoint (nicht gecached)
    location /health {
        proxy_pass http://fits_backend;
        proxy_no_cache 1;
        proxy_cache_bypass 1;
    }

    # Swagger UI
    location /swagger/ {
        proxy_pass http://fits_backend;
    }
}
```

```bash
sudo ln -s /etc/nginx/sites-available/fits-backend /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### 7. SSL/TLS mit Let's Encrypt

```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d api.fits.example.com
sudo certbot renew --dry-run  # Test auto-renewal
```

### 8. Firewall

```bash
sudo ufw allow 22/tcp   # SSH
sudo ufw allow 80/tcp   # HTTP
sudo ufw allow 443/tcp  # HTTPS
sudo ufw enable
sudo ufw status
```

### 9. Monitoring & Logging

#### Logs ansehen

```bash
# Service Logs
sudo journalctl -u fits-backend -f

# Nginx Logs
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log
```

#### Log Rotation

```bash
sudo nano /etc/logrotate.d/fits-backend
```

```
/var/log/fits-backend/*.log {
    daily
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 fits fits
    sharedscripts
    postrotate
        systemctl reload fits-backend
    endscript
}
```

### 10. Backup Strategy

#### Database Backup

```bash
#!/bin/bash
# /opt/fits-backend/backup.sh

BACKUP_DIR="/var/backups/fits"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR

# PostgreSQL Dump
pg_dump -U fits_prod -h localhost fits_prod | \
    gzip > $BACKUP_DIR/db_backup_$TIMESTAMP.sql.gz

# Alte Backups löschen (älter als 30 Tage)
find $BACKUP_DIR -name "db_backup_*.sql.gz" -mtime +30 -delete
```

```bash
# Cron Job für tägliches Backup
sudo crontab -e
```

```
0 2 * * * /opt/fits-backend/backup.sh
```

### 11. Health Checks & Monitoring

#### Prometheus (optional)

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'fits-backend'
    metrics_path: '/metrics'
    params:
      secret: ['your-metrics-secret']
    static_configs:
      - targets: ['localhost:8080']
```

#### Uptime Monitoring

```bash
# Simple Health Check Script
#!/bin/bash
# /opt/fits-backend/healthcheck.sh

HEALTH_URL="http://localhost:8080/health"
RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" $HEALTH_URL)

if [ $RESPONSE -ne 200 ]; then
    echo "Health check failed: HTTP $RESPONSE"
    systemctl restart fits-backend
fi
```

### 12. Updates & Maintenance

#### Update Prozess

```bash
# 1. Neues Binary bauen
go build -ldflags="-s -w" -o fits-server cmd/server/main.go

# 2. Auf Server kopieren
scp fits-server user@server:/tmp/

# 3. Service stoppen und ersetzen
sudo systemctl stop fits-backend
sudo cp /tmp/fits-server /opt/fits-backend/
sudo chown fits:fits /opt/fits-backend/fits-server
sudo chmod 755 /opt/fits-backend/fits-server
sudo systemctl start fits-backend

# 4. Verifizieren
sudo systemctl status fits-backend
curl https://api.fits.example.com/health
```

---

## Docker Deployment (Alternative)

### Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o fits-server cmd/server/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/fits-server .
COPY configs configs

EXPOSE 8080

CMD ["./fits-server"]
```

### docker-compose.yml

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: fits_db
      POSTGRES_USER: fits_user
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  fits-backend:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - postgres
    environment:
      DB_HOST: postgres
      DB_PASSWORD: ${DB_PASSWORD}
    volumes:
      - ./uploads:/root/uploads

volumes:
  postgres_data:
```

```bash
# Deployment
docker-compose up -d
docker-compose logs -f
```

---

## Performance Tuning

### PostgreSQL

```conf
# /etc/postgresql/15/main/postgresql.conf

shared_buffers = 256MB
effective_cache_size = 1GB
maintenance_work_mem = 64MB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1
effective_io_concurrency = 200
work_mem = 16MB
min_wal_size = 1GB
max_wal_size = 4GB
```

### Go Application

In `configs/config.toml`:

```toml
[database]
max_conns = 50  # Erhöhen für mehr Last
min_conns = 10
```

---

## Troubleshooting

### Service startet nicht

```bash
# Logs prüfen
sudo journalctl -u fits-backend -n 50 --no-pager

# Manueller Test
sudo -u fits /opt/fits-backend/fits-server
```

### Datenbank-Verbindung schlägt fehl

```bash
# PostgreSQL Status
sudo systemctl status postgresql

# Connection Test
psql -h localhost -U fits_prod -d fits_prod
```

### Hoher Memory-Verbrauch

```bash
# Memory Profiling aktivieren
go tool pprof http://localhost:8080/debug/pprof/heap
```

---

**Version:** 1.0
**Zuletzt aktualisiert:** 2025-10-18
