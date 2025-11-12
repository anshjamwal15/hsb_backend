# 🐳 Docker Setup Guide

Complete guide for running HSB Backend in Docker containers.

## 📋 Prerequisites

- Docker 20.10 or higher
- Docker Compose 2.0 or higher

## 🚀 Quick Start with Docker

### Option 1: Using Make (Recommended)

```bash
# Start everything with one command
make docker-run
```

### Option 2: Using Docker Compose

```bash
# Copy environment file
cp .env.docker .env

# Start all services
docker-compose up -d
```

### Option 3: Build and Run Manually

```bash
# Build the image
docker build -t hsb-backend:latest .

# Run with docker-compose
docker-compose up -d
```

## 📦 What Gets Started

When you run `docker-compose up`, three services start:

### 1. **MongoDB** (Port 27017)
- Database for the application
- Persistent data storage
- Health checks enabled

### 2. **HSB Backend** (Port 8080)
- Main application server
- Auto-connects to MongoDB
- Swagger UI available

### 3. **Mongo Express** (Port 8081) - Optional
- Web-based MongoDB admin interface
- Username: `admin`
- Password: `admin123`

## 🌐 Access Your Services

After starting with Docker:

| Service | URL | Description |
|---------|-----|-------------|
| **Backend API** | http://localhost:8080 | Main API server |
| **Swagger UI** | http://localhost:8080/swagger | Interactive API docs |
| **Health Check** | http://localhost:8080/health | Server health status |
| **Mongo Express** | http://localhost:8081 | Database management UI |

## 🛠️ Common Commands

### Start Services
```bash
# Start all services
make docker-run
# or
docker-compose up -d
```

### View Logs
```bash
# View backend logs
make docker-logs
# or
docker-compose logs -f hsb_backend

# View all logs
docker-compose logs -f
```

### Stop Services
```bash
# Stop all services
make docker-stop
# or
docker-compose down
```

### Restart Services
```bash
# Restart all services
make docker-restart
# or
docker-compose restart
```

### Rebuild After Code Changes
```bash
# Rebuild and restart
make docker-rebuild
# or
docker-compose up -d --build
```

### Clean Everything (Including Data)
```bash
# Remove containers and volumes
make docker-clean
# or
docker-compose down -v
```

## 🔧 Configuration

### Environment Variables

Edit `.env` file or set in `docker-compose.yml`:

```env
# JWT Secret
JWT_SECRET=your-secret-key

# Razorpay
RAZORPAY_KEY_ID=your_key_id
RAZORPAY_KEY_SECRET=your_key_secret
```

### MongoDB Configuration

Default credentials (change in production):
- Username: `admin`
- Password: `admin123`
- Database: `hsb_backend`

To change MongoDB credentials, edit `docker-compose.yml`:
```yaml
environment:
  MONGO_INITDB_ROOT_USERNAME: your_username
  MONGO_INITDB_ROOT_PASSWORD: your_password
```

### Port Configuration

To change ports, edit `docker-compose.yml`:
```yaml
ports:
  - "8080:8080"  # Change first number for host port
```

## 📊 Health Checks

All services have health checks:

```bash
# Check container health
docker-compose ps

# Manual health check
curl http://localhost:8080/health
```

## 🔍 Debugging

### View Container Status
```bash
docker-compose ps
```

### View Backend Logs
```bash
docker-compose logs hsb_backend
```

### View MongoDB Logs
```bash
docker-compose logs mongodb
```

### Enter Container Shell
```bash
# Backend container
docker exec -it hsb_backend sh

# MongoDB container
docker exec -it hsb_mongodb mongosh
```

### Check MongoDB Connection
```bash
# Connect to MongoDB
docker exec -it hsb_mongodb mongosh -u admin -p admin123

# List databases
show dbs

# Use hsb_backend database
use hsb_backend

# Show collections
show collections
```

## 🚀 Production Deployment

### Build Production Image
```bash
docker build -t hsb-backend:v1.0.0 .
```

### Push to Registry
```bash
# Tag for registry
docker tag hsb-backend:v1.0.0 your-registry/hsb-backend:v1.0.0

# Push to registry
docker push your-registry/hsb-backend:v1.0.0
```

### Production docker-compose.yml
```yaml
version: '3.8'

services:
  hsb_backend:
    image: your-registry/hsb-backend:v1.0.0
    restart: always
    ports:
      - "8080:8080"
    environment:
      MONGODB_URI: ${MONGODB_URI}
      JWT_SECRET: ${JWT_SECRET}
      RAZORPAY_KEY_ID: ${RAZORPAY_KEY_ID}
      RAZORPAY_KEY_SECRET: ${RAZORPAY_KEY_SECRET}
```

## 🔒 Security Best Practices

### 1. Change Default Passwords
```yaml
# In docker-compose.yml
MONGO_INITDB_ROOT_PASSWORD: strong_password_here
```

### 2. Use Secrets for Sensitive Data
```yaml
secrets:
  jwt_secret:
    file: ./secrets/jwt_secret.txt
```

### 3. Run as Non-Root User
Already configured in Dockerfile:
```dockerfile
USER appuser
```

### 4. Use Environment Variables
Never hardcode secrets in docker-compose.yml:
```yaml
environment:
  JWT_SECRET: ${JWT_SECRET}
```

## 📈 Monitoring

### Container Stats
```bash
docker stats hsb_backend hsb_mongodb
```

### Disk Usage
```bash
docker system df
```

### Clean Unused Resources
```bash
docker system prune -a
```

## 🐛 Troubleshooting

### Container Won't Start
```bash
# Check logs
docker-compose logs hsb_backend

# Check if port is in use
lsof -i :8080
```

### MongoDB Connection Issues
```bash
# Check MongoDB is running
docker-compose ps mongodb

# Check MongoDB logs
docker-compose logs mongodb

# Test connection
docker exec -it hsb_mongodb mongosh --eval "db.adminCommand('ping')"
```

### Backend Can't Connect to MongoDB
```bash
# Check network
docker network ls
docker network inspect hsb_backend_hsb_network

# Restart services
docker-compose restart
```

### Out of Disk Space
```bash
# Clean up
docker system prune -a --volumes

# Remove old images
docker image prune -a
```

## 📚 Additional Resources

- **Dockerfile**: Multi-stage build for optimal image size
- **.dockerignore**: Excludes unnecessary files
- **docker-compose.yml**: Complete service orchestration
- **Makefile**: Convenient command shortcuts

## 🎯 Complete Workflow Example

```bash
# 1. Clone and setup
cd /Users/ansh/volvrit/production/hsb_backend
cp .env.docker .env

# 2. Start services
make docker-run

# 3. Check health
curl http://localhost:8080/health

# 4. Open Swagger UI
open http://localhost:8080/swagger

# 5. View logs
make docker-logs

# 6. Stop when done
make docker-stop
```

## ✅ Verification Checklist

After starting Docker:

- [ ] All containers are running: `docker-compose ps`
- [ ] Backend is healthy: `curl http://localhost:8080/health`
- [ ] Swagger UI loads: http://localhost:8080/swagger
- [ ] MongoDB is accessible: http://localhost:8081
- [ ] Logs show no errors: `docker-compose logs`

---

**Need Help?** Check the logs with `make docker-logs` or `docker-compose logs -f`
