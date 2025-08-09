#!/usr/bin/env bash
set -euo pipefail

BACKEND_PORT=${BACKEND_PORT:-8080}
FRONTEND_PORT=${FRONTEND_PORT:-5173}
API_HOST=${API_HOST:-0.0.0.0}

export VITE_API_URL="http://localhost:${BACKEND_PORT}" API_PORT=${BACKEND_PORT} API_HOST

echo "🚀 Starting MetricHub local development (no Docker)"

echo "▶️  Starting backend on :$BACKEND_PORT"
(cd backend && go run ./cmd/server) &
BACKEND_PID=$!

sleep 1

echo "▶️  Starting frontend on :$FRONTEND_PORT (VITE_API_URL=$VITE_API_URL)"
(cd frontend && npm run dev) &
FRONTEND_PID=$!

cleanup() {
  echo "\n🧹 Stopping processes..."
  kill $BACKEND_PID $FRONTEND_PID 2>/dev/null || true
  wait $BACKEND_PID $FRONTEND_PID 2>/dev/null || true
  echo "✅ Stopped."
}
trap cleanup EXIT INT TERM

echo "\n🌐 Frontend:  http://localhost:$FRONTEND_PORT"
echo "🛠  Backend API: http://localhost:$BACKEND_PORT"
echo "💓 Health:     http://localhost:$BACKEND_PORT/api/v1/health"
echo "\nPress Ctrl+C to stop."

wait
