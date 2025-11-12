#!/bin/bash

# HSB Backend Server Runner
# This script starts the HSB backend server

echo "🚀 Starting HSB Backend Server..."
echo ""

# Check if .env file exists
if [ ! -f .env ]; then
    echo "⚠️  .env file not found!"
    echo "📝 Creating .env from .env.example..."
    cp .env.example .env
    echo "✅ .env file created. Please update it with your configuration."
    echo ""
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

# Check if MongoDB is running
if ! pgrep -x "mongod" > /dev/null; then
    echo "⚠️  MongoDB doesn't appear to be running."
    echo "💡 Start MongoDB with: mongod"
    echo "   Or using Docker: docker run -d -p 27017:27017 --name mongodb mongo:latest"
    echo ""
fi

# Download dependencies if needed
if [ ! -d "vendor" ]; then
    echo "📦 Downloading dependencies..."
    go mod download
    echo ""
fi

# Run the server
echo "🎯 Starting server..."
echo ""
go run cmd/server/main.go
