# Deployment Guide

This guide covers deploying Anexis Server in various environments.

## Prerequisites

- Docker 24.0+
- Docker Compose 2.0+
- Domain name (for production)
- SSL certificate (for production)
- Backblaze B2 account

## Development Deployment

### Quick Start

```bash
# Start PostgreSQL and API server
make dev

# View logs
make logs

# Stop
make stop
```

### Local Development (without Docker)

```bash
# Start PostgreSQL manually or use Docker
docker run -d \
  --name anexis-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=anexis \
  -p 5432:5432 \
  postgres:16-alpine

# Run API server locally
make dev-local
```

## Production Deployment

### Environment Setup

1. Create `.env` file:
```bash
cp .env.example .env
```

2. Configure production values:
```env
ENVIRONMENT=production
DEBUG=false

# Strong JWT secret
JWT_SECRET=<generate-a-secure-random-string>

# Database (use strong password)
DB_HOST=postgres
DB_PASSWORD=<strong-password>

# Backblaze B2
B2_APPLICATION_KEY_ID=<your-key-id>
B2_APPLICATION_KEY=<your-application-key>
B2_BUCKET_NAME=<your-bucket>
```

### SSL Certificate Setup

Place your SSL certificates in `nginx/ssl/`:
```
nginx/ssl/
├── cert.pem
└── key.pem
```

Or use Let's Encrypt:
```bash
certbot certonly --standalone -d your-domain.com
```

### Enable HTTPS in Nginx

Edit `nginx/nginx.conf`, uncomment the HTTPS server block:
```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    # ... rest of config
}
```

### Start Production

```bash
# Start with 3 API replicas
make prod

# Scale to 5 replicas
make prod-scale N=5

# View logs
make logs-prod
```

## Docker Swarm Deployment

For larger deployments, use Docker Swarm:

```bash
# Initialize swarm
docker swarm init

# Deploy stack
docker stack deploy -c docker-compose.prod.yml anexis

# Scale service
docker service scale anexis_api=5

# View services
docker service ls
```

## Kubernetes Deployment

Basic Kubernetes manifests:

### Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: anexis-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: anexis-api
  template:
    metadata:
      labels:
        app: anexis-api
    spec:
      containers:
      - name: api
        image: anexis-api:latest
        ports:
        - containerPort: 8080
        envFrom:
        - secretRef:
            name: anexis-secrets
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
```

### Service
```yaml
apiVersion: v1
kind: Service
metadata:
  name: anexis-api
spec:
  selector:
    app: anexis-api
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

## Database Migrations

```bash
# Check migration status
make migrate-status

# Run pending migrations
make migrate

# Create new migration
make migrate-new NAME=add_feature
```

## Monitoring

### Health Check

```bash
# Check API health
curl http://localhost:8080/health
```

### Logs

```bash
# Docker Compose
docker-compose logs -f api

# Docker Swarm
docker service logs -f anexis_api
```

### Metrics

Consider adding:
- Prometheus for metrics collection
- Grafana for visualization
- Jaeger for distributed tracing

## Backup Strategy

### Database Backup

```bash
# Backup PostgreSQL
docker-compose exec postgres pg_dump -U postgres anexis > backup.sql

# Restore
cat backup.sql | docker-compose exec -T postgres psql -U postgres anexis
```

### Application Backups

Users can export their data via the Backup API:
- `POST /api/v1/backup/export`
- `GET /api/v1/backup/:id/download`

## Troubleshooting

### API Not Starting

```bash
# Check logs
docker-compose logs api

# Common issues:
# - Database not ready (check postgres health)
# - Invalid environment variables
# - Port already in use
```

### Database Connection Issues

```bash
# Test connection
docker-compose exec postgres pg_isready -U postgres

# Check connection string in logs
docker-compose logs api | grep "database"
```

### Storage Issues

```bash
# Verify B2 credentials
curl -X GET \
  "https://api.backblazeb2.com/b2api/v2/b2_authorize_account" \
  -u "$B2_APPLICATION_KEY_ID:$B2_APPLICATION_KEY"
```
