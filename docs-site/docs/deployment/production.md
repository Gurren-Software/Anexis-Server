---
sidebar_position: 3
slug: /deployment/production
---

# Production Deployment

Deploy Anexis Server for production with high availability.

## Architecture

```
                    ┌─────────────┐
                    │    Nginx    │
                    │ Load Balancer│
                    └──────┬──────┘
                           │
         ┌─────────────────┼─────────────────┐
         │                 │                 │
    ┌────▼────┐      ┌────▼────┐      ┌────▼────┐
    │  API 1  │      │  API 2  │      │  API 3  │
    └────┬────┘      └────┬────┘      └────┬────┘
         │                 │                 │
         └─────────────────┼─────────────────┘
                           │
                    ┌──────▼──────┐
                    │ PostgreSQL  │
                    │   Primary   │
                    └─────────────┘
```

## Using Docker Compose Production

```bash
# Clone repository
git clone https://github.com/Gurren-Software/Anexis-Server.git
cd anexis-server

# Edit production environment
cp .env.example .env.prod
nano .env.prod

# Start production stack
docker-compose -f docker-compose.prod.yml up -d
```

### Production Configuration

```bash
# .env.prod
SERVER_MODE=saas
ENVIRONMENT=production
DEBUG=false

# Storage (B2 example)
STORAGE_PROVIDER=b2
B2_APPLICATION_KEY_ID=your_key_id
B2_APPLICATION_KEY=your_key
B2_BUCKET_NAME=your_bucket

# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=anexis
DB_PASSWORD=strong_password
DB_NAME=anexis
DB_SSLMODE=enable

# Security
JWT_SECRET=very_long_random_secret
JWT_EXPIRATION_HOURS=24
```

---

## Manual Production Setup

### 1. Build Binary

```bash
cd apps/api
CGO_ENABLED=0 go build -a -installsuffix cgo -o anexis-server ./cmd/server
```

### 2. Systemd Service

Create `/etc/systemd/system/anexis.service`:

```ini
[Unit]
Description=Anexis Server
After=network.target postgresql.service

[Service]
Type=simple
User=anexis
WorkingDirectory=/opt/anexis
ExecStart=/opt/anexis/anexis-server
EnvironmentFile=/opt/anexis/.env
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

### 3. Configure Nginx

```nginx
upstream anexis_backend {
    server 127.0.0.1:8081;
    server 127.0.0.1:8082;
    server 127.0.0.1:8083;
}

server {
    listen 80;
    server_name yourdomain.com;
    
    location / {
        proxy_pass http://anexis_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

---

## Security Checklist

- [ ] Change default JWT_SECRET
- [ ] Enable SSL/TLS
- [ ] Configure firewall rules
- [ ] Set proper file permissions
- [ ] Enable rate limiting
- [ ] Configure log rotation
- [ ] Set up monitoring