# 🚀 Complete Deployment Guide

This guide covers all deployment options for HSB Backend.

## 📋 Table of Contents

1. [Docker Deployment (Recommended)](#docker-deployment)
2. [Local Development](#local-development)
3. [Production Deployment](#production-deployment)
4. [Cloud Deployment](#cloud-deployment)

---

## 🐳 Docker Deployment (Recommended)

### Quick Start

```bash
# One command to start everything
./docker-start.sh

# Or using Make
make docker-run

# Or using docker-compose
docker-compose up -d
```

### What You Get

- ✅ Backend API on port 8080
- ✅ MongoDB on port 27017
- ✅ Mongo Express UI on port 8081
- ✅ Swagger UI at http://localhost:8080/swagger
- ✅ Automatic health checks
- ✅ Persistent data storage

### Access Points

| Service | URL | Credentials |
|---------|-----|-------------|
| API | http://localhost:8080 | - |
| Swagger | http://localhost:8080/swagger | - |
| Mongo Express | http://localhost:8081 | admin/admin123 |

### Common Commands

```bash
# Start services
make docker-run

# View logs
make docker-logs

# Stop services
make docker-stop

# Clean everything
make docker-clean

# Rebuild after code changes
make docker-rebuild
```

### Configuration

Edit `.env` file:
```env
JWT_SECRET=your-secret-key
RAZORPAY_KEY_ID=your_key_id
RAZORPAY_KEY_SECRET=your_key_secret
```

---

## 💻 Local Development

### Prerequisites

- Go 1.21+
- MongoDB 4.4+

### Setup

```bash
# 1. Install dependencies
go mod download

# 2. Setup environment
cp .env.example .env
# Edit .env with your configuration

# 3. Start MongoDB
mongod

# 4. Run the server
go run cmd/server/main.go

# Or use the run script
./run.sh
```

### Development Workflow

```bash
# Run with hot reload (if you have air installed)
air

# Run tests
go test ./...

# Build binary
go build -o bin/server cmd/server/main.go

# Run binary
./bin/server
```

---

## 🏭 Production Deployment

### Option 1: Docker Production

#### 1. Build Production Image

```bash
# Build optimized image
docker build -t hsb-backend:v1.0.0 .

# Test locally
docker run -p 8080:8080 \
  -e MONGODB_URI=mongodb://your-mongo-host:27017 \
  -e JWT_SECRET=your-secret \
  hsb-backend:v1.0.0
```

#### 2. Push to Registry

```bash
# Tag for your registry
docker tag hsb-backend:v1.0.0 your-registry.com/hsb-backend:v1.0.0

# Push
docker push your-registry.com/hsb-backend:v1.0.0
```

#### 3. Deploy with Docker Compose

Create `docker-compose.prod.yml`:
```yaml
version: '3.8'

services:
  hsb_backend:
    image: your-registry.com/hsb-backend:v1.0.0
    restart: always
    ports:
      - "8080:8080"
    environment:
      MONGODB_URI: ${MONGODB_URI}
      JWT_SECRET: ${JWT_SECRET}
      RAZORPAY_KEY_ID: ${RAZORPAY_KEY_ID}
      RAZORPAY_KEY_SECRET: ${RAZORPAY_KEY_SECRET}
    healthcheck:
      test: ["CMD", "wget", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

Deploy:
```bash
docker-compose -f docker-compose.prod.yml up -d
```

### Option 2: Binary Deployment

#### 1. Build Binary

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o hsb-backend cmd/server/main.go

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o hsb-backend cmd/server/main.go

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o hsb-backend.exe cmd/server/main.go
```

#### 2. Deploy Binary

```bash
# Copy to server
scp hsb-backend user@server:/opt/hsb-backend/
scp swagger.yaml user@server:/opt/hsb-backend/
scp .env user@server:/opt/hsb-backend/

# SSH to server
ssh user@server

# Run
cd /opt/hsb-backend
./hsb-backend
```

#### 3. Create Systemd Service

Create `/etc/systemd/system/hsb-backend.service`:
```ini
[Unit]
Description=HSB Backend Service
After=network.target

[Service]
Type=simple
User=hsb
WorkingDirectory=/opt/hsb-backend
ExecStart=/opt/hsb-backend/hsb-backend
Restart=always
RestartSec=10
Environment="PORT=8080"
EnvironmentFile=/opt/hsb-backend/.env

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl daemon-reload
sudo systemctl enable hsb-backend
sudo systemctl start hsb-backend
sudo systemctl status hsb-backend
```

---

## ☁️ Cloud Deployment

### AWS ECS (Elastic Container Service)

#### 1. Push to ECR

```bash
# Login to ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin YOUR_ACCOUNT.dkr.ecr.us-east-1.amazonaws.com

# Tag and push
docker tag hsb-backend:latest YOUR_ACCOUNT.dkr.ecr.us-east-1.amazonaws.com/hsb-backend:latest
docker push YOUR_ACCOUNT.dkr.ecr.us-east-1.amazonaws.com/hsb-backend:latest
```

#### 2. Create Task Definition

```json
{
  "family": "hsb-backend",
  "containerDefinitions": [
    {
      "name": "hsb-backend",
      "image": "YOUR_ACCOUNT.dkr.ecr.us-east-1.amazonaws.com/hsb-backend:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "PORT",
          "value": "8080"
        }
      ],
      "secrets": [
        {
          "name": "MONGODB_URI",
          "valueFrom": "arn:aws:secretsmanager:region:account:secret:mongodb-uri"
        },
        {
          "name": "JWT_SECRET",
          "valueFrom": "arn:aws:secretsmanager:region:account:secret:jwt-secret"
        }
      ],
      "healthCheck": {
        "command": ["CMD-SHELL", "wget --spider http://localhost:8080/health || exit 1"],
        "interval": 30,
        "timeout": 5,
        "retries": 3
      }
    }
  ]
}
```

### Google Cloud Run

```bash
# Build and push to GCR
gcloud builds submit --tag gcr.io/PROJECT_ID/hsb-backend

# Deploy
gcloud run deploy hsb-backend \
  --image gcr.io/PROJECT_ID/hsb-backend \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars MONGODB_URI=your-mongo-uri \
  --set-env-vars JWT_SECRET=your-secret
```

### Azure Container Instances

```bash
# Login to Azure
az login

# Create resource group
az group create --name hsb-backend-rg --location eastus

# Deploy container
az container create \
  --resource-group hsb-backend-rg \
  --name hsb-backend \
  --image your-registry.azurecr.io/hsb-backend:latest \
  --dns-name-label hsb-backend \
  --ports 8080 \
  --environment-variables \
    PORT=8080 \
    MONGODB_URI=your-mongo-uri \
  --secure-environment-variables \
    JWT_SECRET=your-secret
```

### DigitalOcean App Platform

Create `app.yaml`:
```yaml
name: hsb-backend
services:
- name: api
  github:
    repo: your-username/hsb-backend
    branch: main
  dockerfile_path: Dockerfile
  http_port: 8080
  health_check:
    http_path: /health
  envs:
  - key: MONGODB_URI
    value: ${MONGODB_URI}
  - key: JWT_SECRET
    value: ${JWT_SECRET}
    type: SECRET
```

Deploy:
```bash
doctl apps create --spec app.yaml
```

### Kubernetes

Create `k8s-deployment.yaml`:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hsb-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: hsb-backend
  template:
    metadata:
      labels:
        app: hsb-backend
    spec:
      containers:
      - name: hsb-backend
        image: your-registry/hsb-backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        - name: MONGODB_URI
          valueFrom:
            secretKeyRef:
              name: hsb-secrets
              key: mongodb-uri
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: hsb-secrets
              key: jwt-secret
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: hsb-backend-service
spec:
  selector:
    app: hsb-backend
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer
```

Deploy:
```bash
kubectl apply -f k8s-deployment.yaml
```

---

## 🔒 Security Checklist

### Before Production

- [ ] Change all default passwords
- [ ] Use strong JWT secret (32+ characters)
- [ ] Enable HTTPS/TLS
- [ ] Set up firewall rules
- [ ] Use environment variables for secrets
- [ ] Enable MongoDB authentication
- [ ] Set up backup strategy
- [ ] Configure rate limiting
- [ ] Enable logging and monitoring
- [ ] Set up alerts

### Environment Variables

```bash
# Required
JWT_SECRET=<strong-random-string>
MONGODB_URI=mongodb://user:pass@host:27017/dbname
RAZORPAY_KEY_ID=<your-key>
RAZORPAY_KEY_SECRET=<your-secret>

# Optional
PORT=8080
```

---

## 📊 Monitoring

### Health Check Endpoint

```bash
curl http://localhost:8080/health
```

Response:
```json
{"status":"ok"}
```

### Logs

```bash
# Docker
docker-compose logs -f hsb_backend

# Systemd
journalctl -u hsb-backend -f

# File
tail -f /var/log/hsb-backend.log
```

### Metrics

Consider adding:
- Prometheus for metrics
- Grafana for dashboards
- ELK stack for log aggregation
- Sentry for error tracking

---

## 🆘 Troubleshooting

### Container Won't Start

```bash
# Check logs
docker-compose logs hsb_backend

# Check if port is in use
lsof -i :8080

# Restart
docker-compose restart hsb_backend
```

### MongoDB Connection Issues

```bash
# Test connection
docker exec -it hsb_mongodb mongosh --eval "db.adminCommand('ping')"

# Check network
docker network inspect hsb_backend_hsb_network
```

### High Memory Usage

```bash
# Check container stats
docker stats hsb_backend

# Restart container
docker-compose restart hsb_backend
```

---

## 📚 Additional Resources

- [Docker Documentation](DOCKER.md)
- [Quick Start Guide](QUICKSTART.md)
- [API Documentation](swagger.yaml)
- [Implementation Status](API_IMPLEMENTATION_STATUS.md)

---

**Need Help?** Check the logs or open an issue on GitHub.
