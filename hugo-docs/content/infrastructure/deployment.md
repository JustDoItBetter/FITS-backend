---
title: "Deployment Guide"
weight: 1
description: "Production deployment guide for FITS Backend. Binary deployment, Docker containerization, Kubernetes orchestration, and best practices."
---

# Deployment Guide

## Production Deployment

### Prerequisites

- Go 1.25.1+
- PostgreSQL 12+
- TLS certificates
- Domain name
- Reverse proxy (Nginx/Caddy)

### Deployment Options

## 1. Binary Deployment

### Build

```bash
# Build for production
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags="-w -s" \
  -o fits-backend \
  cmd/server/main.go
```

### Configuration

Create production config at `/etc/fits/config.toml`:

```toml
[server]
host = "0.0.0.0"
port = 8080
read_timeout = "30s"
write_timeout = "30s"
rate_limit = 100
allowed_origins = "https://yourdomain.com"

[server.tls]
enabled = true
cert_file = "/etc/ssl/certs/fits.crt"
key_file = "/etc/ssl/private/fits.key"

[database]
host = "db.internal"
port = 5432
user = "fits_prod"
password_env = "DB_PASSWORD"
database = "fits_production"
sslmode = "require"

[jwt]
secret_env = "JWT_SECRET"
access_token_expiry = "15m"
refresh_token_expiry = "168h"

[logging]
level = "info"
format = "json"
```

### Systemd Service

Create `/etc/systemd/system/fits-backend.service`:

```ini
[Unit]
Description=FITS Backend API
After=network.target postgresql.service

[Service]
Type=simple
User=fits
Group=fits
WorkingDirectory=/opt/fits-backend
ExecStart=/opt/fits-backend/fits-backend
Restart=always
RestartSec=5

Environment="FITS_LOG_LEVEL=info"
EnvironmentFile=/etc/fits/environment

NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl enable fits-backend
sudo systemctl start fits-backend
sudo systemctl status fits-backend
```

## 2. Docker Deployment

### Dockerfile

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /fits-backend cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /fits-backend .
COPY configs/config.toml ./configs/
RUN addgroup -g 1000 fits && adduser -D -u 1000 -G fits fits
USER fits
EXPOSE 8080
CMD ["./fits-backend"]
```

### Docker Compose

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: fits_production
      POSTGRES_USER: fits
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - fits-network
    restart: always

  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      FITS_DB_HOST: postgres
      FITS_JWT_SECRET: ${JWT_SECRET}
    networks:
      - fits-network
    depends_on:
      - postgres
    restart: always

volumes:
  postgres_data:

networks:
  fits-network:
```

## 3. Kubernetes Deployment

### Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: fits-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: fits-backend
  template:
    metadata:
      labels:
        app: fits-backend
    spec:
      containers:
      - name: api
        image: your-registry/fits-backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: FITS_DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: fits-secrets
              key: db-password
```

## Monitoring

### Prometheus

```yaml
scrape_configs:
  - job_name: 'fits-backend'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
```

## Backups

```bash
# Database backup script
pg_dump -U fits fits_production | gzip > backup.sql.gz
```

## Next Steps

- [Security Considerations](/infrastructure/security-considerations/)
