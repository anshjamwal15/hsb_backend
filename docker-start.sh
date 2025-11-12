#!/bin/bash

# HSB Backend Docker Startup Script
# This script starts the entire HSB backend stack using Docker

set -e

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "🐳 HSB Backend - Docker Startup"
echo "═══════════════════════════════════════════════════════════════"
echo ""

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed!"
    echo "📥 Please install Docker from: https://docs.docker.com/get-docker/"
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed!"
    echo "📥 Please install Docker Compose from: https://docs.docker.com/compose/install/"
    exit 1
fi

# Check if Docker daemon is running
if ! docker info &> /dev/null; then
    echo "❌ Docker daemon is not running!"
    echo "🔧 Please start Docker Desktop or Docker daemon"
    exit 1
fi

echo "✅ Docker is installed and running"
echo ""

# Check if .env file exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file from .env.docker..."
    cp .env.docker .env
    echo "✅ .env file created"
    echo "⚠️  Please update .env with your Razorpay credentials if needed"
    echo ""
fi

# Stop any existing containers
echo "🛑 Stopping any existing containers..."
docker-compose down 2>/dev/null || true
echo ""

# Build and start containers
echo "🔨 Building Docker images..."
docker-compose build
echo ""

echo "🚀 Starting Docker containers..."
docker-compose up -d
echo ""

# Wait for services to be healthy
echo "⏳ Waiting for services to be ready..."
sleep 5

# Check if containers are running
if docker-compose ps | grep -q "Up"; then
    echo "✅ All services started successfully!"
    echo ""
    
    echo "═══════════════════════════════════════════════════════════════"
    echo "🎉 HSB Backend is now running!"
    echo "═══════════════════════════════════════════════════════════════"
    echo ""
    echo "📡 Access your services:"
    echo ""
    echo "   🌐 Backend API:       http://localhost:8080"
    echo "   📚 Swagger UI:        http://localhost:8080/swagger"
    echo "   📄 Swagger YAML:      http://localhost:8080/swagger.yaml"
    echo "   ❤️  Health Check:      http://localhost:8080/health"
    echo "   🗄️  Mongo Express:     http://localhost:8081"
    echo "      └─ Username: admin"
    echo "      └─ Password: admin123"
    echo ""
    echo "═══════════════════════════════════════════════════════════════"
    echo "📋 Useful Commands:"
    echo "═══════════════════════════════════════════════════════════════"
    echo ""
    echo "   View logs:           docker-compose logs -f"
    echo "   View backend logs:   docker-compose logs -f hsb_backend"
    echo "   Stop services:       docker-compose down"
    echo "   Restart services:    docker-compose restart"
    echo "   Clean everything:    docker-compose down -v"
    echo ""
    echo "   Or use Make commands:"
    echo "   View logs:           make docker-logs"
    echo "   Stop services:       make docker-stop"
    echo "   Clean everything:    make docker-clean"
    echo ""
    echo "═══════════════════════════════════════════════════════════════"
    echo "💡 Quick Start:"
    echo "═══════════════════════════════════════════════════════════════"
    echo ""
    echo "   1. Open Swagger UI: http://localhost:8080/swagger"
    echo "   2. Register a user using POST /user/register"
    echo "   3. Login using POST /user/login"
    echo "   4. Copy the JWT token from response"
    echo "   5. Click 'Authorize' button and paste token"
    echo "   6. Test protected endpoints!"
    echo ""
    echo "═══════════════════════════════════════════════════════════════"
    echo "🎯 Container Status:"
    echo "═══════════════════════════════════════════════════════════════"
    echo ""
    docker-compose ps
    echo ""
    echo "═══════════════════════════════════════════════════════════════"
    echo "✨ Ready to accept requests!"
    echo "═══════════════════════════════════════════════════════════════"
    echo ""
    
    # Offer to show logs
    echo "📊 Would you like to view the logs? (y/n)"
    read -t 5 -n 1 answer || answer="n"
    echo ""
    if [ "$answer" = "y" ] || [ "$answer" = "Y" ]; then
        echo "Showing logs (Press Ctrl+C to exit)..."
        docker-compose logs -f
    fi
else
    echo "❌ Failed to start services!"
    echo "📋 Checking logs..."
    docker-compose logs
    exit 1
fi
