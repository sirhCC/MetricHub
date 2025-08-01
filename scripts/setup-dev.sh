#!/bin/bash

# Development setup script for MetricHub
set -e

echo "🚀 Setting up MetricHub development environment..."

# Check dependencies
check_dependency() {
    if ! command -v $1 &> /dev/null; then
        echo "❌ $1 is not installed. Please install it first."
        exit 1
    fi
    echo "✅ $1 is available"
}

echo "📋 Checking dependencies..."
check_dependency "docker"
check_dependency "docker-compose"
check_dependency "go"
check_dependency "node"

# Create necessary directories
echo "📁 Creating directories..."
mkdir -p backend/tmp
mkdir -p logs

# Install Go dependencies
echo "📦 Installing Go dependencies..."
cd backend
go mod tidy
cd ..

# Install Node.js dependencies
echo "📦 Installing Node.js dependencies..."
cd frontend
npm install
cd ..

# Create environment file
echo "🔧 Setting up environment..."
if [ ! -f .env ]; then
    cat > .env << EOF
# Database
DATABASE_URL=postgres://metrichub:metrichub_dev@localhost:5432/metrichub?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379

# NATS
NATS_URL=nats://localhost:4222

# API Configuration
API_PORT=8080
API_HOST=0.0.0.0

# Frontend
VITE_API_URL=http://localhost:8080

# Logging
LOG_LEVEL=debug

# JWT Secret (change in production)
JWT_SECRET=your-secret-key-change-in-production
EOF
    echo "✅ Created .env file"
else
    echo "✅ .env file already exists"
fi

echo "🎉 Development environment setup complete!"
echo ""
echo "To start development:"
echo "  1. Start services: docker-compose -f docker-compose.dev.yml up"
echo "  2. Access frontend: http://localhost:3000"
echo "  3. Access backend API: http://localhost:8080"
echo ""
echo "For hot reloading:"
echo "  Backend: cd backend && go run cmd/server/main.go"
echo "  Frontend: cd frontend && npm run dev"
