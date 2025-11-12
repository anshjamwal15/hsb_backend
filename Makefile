.PHONY: help build run stop clean logs docker-build docker-run docker-stop docker-clean

# Default target
help:
	@echo "HSB Backend - Available Commands:"
	@echo ""
	@echo "Local Development:"
	@echo "  make run              - Run the server locally"
	@echo "  make build            - Build the binary"
	@echo "  make clean            - Clean build artifacts"
	@echo ""
	@echo "Docker Commands:"
	@echo "  make docker-build     - Build Docker image"
	@echo "  make docker-run       - Run with Docker Compose"
	@echo "  make docker-stop      - Stop Docker containers"
	@echo "  make docker-clean     - Remove Docker containers and volumes"
	@echo "  make docker-logs      - View Docker logs"
	@echo "  make docker-restart   - Restart Docker containers"
	@echo ""
	@echo "Utilities:"
	@echo "  make test             - Run tests"
	@echo "  make lint             - Run linter"
	@echo ""

# Local development
run:
	@echo "🚀 Starting HSB Backend locally..."
	go run cmd/server/main.go

build:
	@echo "🔨 Building HSB Backend..."
	go build -o bin/server cmd/server/main.go
	@echo "✅ Build complete: bin/server"

clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	go clean
	@echo "✅ Clean complete"

test:
	@echo "🧪 Running tests..."
	go test -v ./...

lint:
	@echo "🔍 Running linter..."
	golangci-lint run

# Docker commands
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t hsb-backend:latest .
	@echo "✅ Docker image built successfully"

docker-run:
	@echo "🐳 Starting Docker containers..."
	@if [ ! -f .env ]; then \
		echo "📝 Creating .env from .env.docker..."; \
		cp .env.docker .env; \
	fi
	docker-compose up -d
	@echo ""
	@echo "✅ Docker containers started!"
	@echo ""
	@echo "📡 Services:"
	@echo "   Backend:        http://localhost:8080"
	@echo "   Swagger UI:     http://localhost:8080/swagger"
	@echo "   Health Check:   http://localhost:8080/health"
	@echo "   Mongo Express:  http://localhost:8081 (admin/admin123)"
	@echo ""
	@echo "📊 View logs: make docker-logs"

docker-stop:
	@echo "🛑 Stopping Docker containers..."
	docker-compose down
	@echo "✅ Containers stopped"

docker-clean:
	@echo "🧹 Removing Docker containers and volumes..."
	docker-compose down -v
	@echo "✅ Cleanup complete"

docker-logs:
	@echo "📋 Viewing Docker logs..."
	docker-compose logs -f hsb_backend

docker-restart:
	@echo "🔄 Restarting Docker containers..."
	docker-compose restart
	@echo "✅ Containers restarted"

# Combined commands
docker-rebuild: docker-stop docker-build docker-run
	@echo "✅ Docker rebuild complete"

docker-fresh: docker-clean docker-build docker-run
	@echo "✅ Fresh Docker setup complete"
