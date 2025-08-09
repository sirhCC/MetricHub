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
  [string]$ApiHost = '0.0.0.0',
  [switch]$AutoRestart
)

Write-Host "Starting MetricHub local development (no Docker)" -ForegroundColor Green

# Ensure frontend env var for this session
$env:VITE_API_URL = "http://localhost:$BackendPort"
$env:API_PORT = $BackendPort
$env:API_HOST = $ApiHost

function Start-Backend {
  Write-Host "Starting backend on :$BackendPort" -ForegroundColor Cyan
  return Start-Process -FilePath go -ArgumentList @('run','./cmd/server') -WorkingDirectory "backend" -PassThru
}
function Start-Frontend {
  Write-Host "Starting frontend on :$FrontendPort (VITE_API_URL=$env:VITE_API_URL)" -ForegroundColor Cyan
  return Start-Process -FilePath npm -ArgumentList @('run','dev') -WorkingDirectory "frontend" -PassThru
}

$backend = Start-Backend
Start-Sleep -Milliseconds 800
$frontend = Start-Frontend

Write-Host ""; Write-Host "Frontend:  http://localhost:$FrontendPort" -ForegroundColor Yellow
Write-Host "Backend API: http://localhost:$BackendPort" -ForegroundColor Yellow
Write-Host "Health:     http://localhost:$BackendPort/api/v1/health" -ForegroundColor Yellow
Write-Host ""; Write-Host "Press 'q' then Enter to stop both processes." -ForegroundColor Magenta

while ($true) {
  $backendExited = $backend.HasExited
  $frontendExited = $frontend.HasExited

  if ($backendExited -or $frontendExited) {
    if ($backendExited) { Write-Host "Backend exited (code $($backend.ExitCode))." -ForegroundColor Red }
    if ($frontendExited) { Write-Host "Frontend exited (code $($frontend.ExitCode))." -ForegroundColor Red }
    if ($AutoRestart) {
      Write-Host "AutoRestart ON: restarting exited processes..." -ForegroundColor Yellow
      if ($backendExited) { $backend = Start-Backend; Start-Sleep -Milliseconds 800 }
      if ($frontendExited) { $frontend = Start-Frontend }
    } else {
      Write-Host "One of the processes exited. Press Enter to terminate or 'r' to restart, 'q' to quit." -ForegroundColor Magenta
  $resp = Read-Host
  if ($resp -eq 'r') {
        if ($backendExited) { $backend = Start-Backend; Start-Sleep -Milliseconds 800 }
        if ($frontendExited) { $frontend = Start-Frontend }
        continue
      }
      break
    }
  }

  if ($Host.UI.RawUI.KeyAvailable) {
    $key = [System.Console]::ReadKey($true)
    if ($key.KeyChar -eq 'q') { break }
  }
  Start-Sleep -Milliseconds 700
}

Write-Host "Stopping processes..." -ForegroundColor Cyan
foreach ($p in @($backend,$frontend)) { if ($p -and -not $p.HasExited) { Stop-Process -Id $p.Id -Force } }
Write-Host "Stopped." -ForegroundColor Green
