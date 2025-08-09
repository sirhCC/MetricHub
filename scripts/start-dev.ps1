<#!
.SYNOPSIS
  Starts MetricHub backend and frontend locally without Docker.
.DESCRIPTION
  Launches Go backend (port 8080) and Vite frontend (port 5173) with env variables,
  watches processes, and provides a simple dashboard for stopping.
#>

param(
  [int]$BackendPort = 8080,
  [int]$FrontendPort = 5173,
  [string]$ApiHost = '0.0.0.0'
)

Write-Host "🚀 Starting MetricHub local development (no Docker)" -ForegroundColor Green

# Ensure frontend env var for this session
$env:VITE_API_URL = "http://localhost:$BackendPort"
$env:API_PORT = $BackendPort
$env:API_HOST = $ApiHost

# Start backend
Write-Host "▶️  Starting backend on :$BackendPort" -ForegroundColor Cyan
$backend = Start-Process powershell -ArgumentList "-NoLogo","-NoExit","-Command","cd backend; go run ./cmd/server" -PassThru
Start-Sleep -Seconds 1

# Start frontend
Write-Host "▶️  Starting frontend on :$FrontendPort (VITE_API_URL=$env:VITE_API_URL)" -ForegroundColor Cyan
$frontend = Start-Process powershell -ArgumentList "-NoLogo","-NoExit","-Command","cd frontend; npm run dev" -PassThru

Write-Host ""; Write-Host "🌐 Frontend:  http://localhost:$FrontendPort" -ForegroundColor Yellow
Write-Host "🛠  Backend API: http://localhost:$BackendPort" -ForegroundColor Yellow
Write-Host "💓 Health:     http://localhost:$BackendPort/api/v1/health" -ForegroundColor Yellow
Write-Host ""; Write-Host "Press 'q' then Enter to stop both processes." -ForegroundColor Magenta

while ($true) {
  if ($backend.HasExited -or $frontend.HasExited) {
    Write-Host "❌ One of the processes exited. Shutting down..." -ForegroundColor Red
    break
  }
  if ($Host.UI.RawUI.KeyAvailable) {
    $key = [System.Console]::ReadLine()
    if ($key -eq 'q') { break }
  }
  Start-Sleep -Milliseconds 500
}

Write-Host "🧹 Stopping processes..." -ForegroundColor Cyan
foreach ($p in @($backend,$frontend)) { if ($p -and -not $p.HasExited) { Stop-Process -Id $p.Id -Force } }
Write-Host "✅ Stopped." -ForegroundColor Green
