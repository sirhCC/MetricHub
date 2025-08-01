# Development setup script for MetricHub (PowerShell)
param(
    [switch]$SkipDependencyCheck
)

Write-Host "🚀 Setting up MetricHub development environment..." -ForegroundColor Green

# Function to check if a command exists
function Test-Command {
    param($Command)
    try {
        Get-Command $Command -ErrorAction Stop | Out-Null
        Write-Host "✅ $Command is available" -ForegroundColor Green
        return $true
    }
    catch {
        Write-Host "❌ $Command is not installed. Please install it first." -ForegroundColor Red
        return $false
    }
}

# Check dependencies
if (-not $SkipDependencyCheck) {
    Write-Host "📋 Checking dependencies..." -ForegroundColor Cyan
    
    $dependencies = @("docker", "go", "node")
    $allDepsAvailable = $true
    
    foreach ($dep in $dependencies) {
        if (-not (Test-Command $dep)) {
            $allDepsAvailable = $false
        }
    }
    
    if (-not $allDepsAvailable) {
        Write-Host "Please install missing dependencies and run again." -ForegroundColor Red
        exit 1
    }
}

# Create necessary directories
Write-Host "📁 Creating directories..." -ForegroundColor Cyan
New-Item -ItemType Directory -Force -Path "backend/tmp", "logs" | Out-Null

# Install Go dependencies
Write-Host "📦 Installing Go dependencies..." -ForegroundColor Cyan
Push-Location backend
& go mod tidy
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Failed to install Go dependencies" -ForegroundColor Red
    Pop-Location
    exit 1
}
Pop-Location

# Install Node.js dependencies
Write-Host "📦 Installing Node.js dependencies..." -ForegroundColor Cyan
Push-Location frontend
& npm install
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Failed to install Node.js dependencies" -ForegroundColor Red
    Pop-Location
    exit 1
}
Pop-Location

# Create environment file
Write-Host "Setting up environment..." -ForegroundColor Cyan
if (-not (Test-Path ".env")) {
    $envContent = @"
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
"@
    Set-Content -Path ".env" -Value $envContent
    Write-Host "Created .env file" -ForegroundColor Green
} else {
    Write-Host ".env file already exists" -ForegroundColor Green
}

Write-Host "🎉 Development environment setup complete!" -ForegroundColor Green
Write-Host ""
Write-Host "To start development:" -ForegroundColor Cyan
Write-Host "  1. Start services: docker-compose -f docker-compose.dev.yml up" -ForegroundColor White
Write-Host "  2. Access frontend: http://localhost:3000" -ForegroundColor White
Write-Host "  3. Access backend API: http://localhost:8080" -ForegroundColor White
Write-Host ""
Write-Host "For hot reloading:" -ForegroundColor Cyan
Write-Host "  Backend: cd backend && go run cmd/server/main.go" -ForegroundColor White
Write-Host "  Frontend: cd frontend && npm run dev" -ForegroundColor White
