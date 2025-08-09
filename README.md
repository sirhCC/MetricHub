git clone https://github.com/sirhCC/MetricHub.git
docker-compose -f docker-compose.dev.yml up
go mod tidy
go run cmd/server/main.go
npm run dev
go test ./...
# MetricHub

> Universal DevOps Metrics Collector – Unified DORA & reliability metrics, plugin ecosystem, and (upcoming) anonymous benchmarking.

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![React](https://img.shields.io/badge/React-18+-61DAFB?logo=react)](https://react.dev/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5+-3178C6?logo=typescript)](https://www.typescriptlang.org/)
[![Status](https://img.shields.io/badge/Stage-Pre_Alpha-orange)](./DEVELOPMENT_PLAN.md)

## Overview

MetricHub is an open-source platform that ingests deployment + incident events, calculates the **four DORA metrics**, classifies performance (elite/high/medium/low), and presents them in a modern dashboard. The current build uses **in‑memory development storage** (perfect for fast iteration) with simulation endpoints so you can explore the experience before real integrations & persistence land.

| Metric | What it Measures | Example Output | Classification |
|--------|------------------|----------------|----------------|
| Deployment Frequency | Deploys per day (time‑window normalized) | 3.2 | elite/high/... |
| Lead Time for Changes | Commit → prod duration (avg) | 2h13m | high |
| MTTR | Incident start → resolve (avg) | 47m | elite |
| Change Failure Rate | Failed / total deployments | 11% | medium |

## Current Status (Pre‑Alpha)

✅ Implemented:

- Go API server (health, DORA metrics, classification)
- Standardized JSON response envelopes (success + error)
- Request ID + structured logging + per‑request timeout middleware
- Lightweight SQL migrations + optional auto‑migrate
- Postgres repository layer (deployments, incidents) with in‑memory fallback
- React + Tailwind dashboard (dark mode, live polling)
- Incident & deployment simulation (UI widget + buttons)
- Real DORA calculation logic (frequency, lead time, MTTR, CFR)
- Performance classification + overall badge

🚧 In Progress / Planned (short term):

- Redis cache layer & query optimization
- Plugin system skeleton & first collectors (GitHub / GitLab / Jenkins)
- AuthN/Z (JWT + roles)
- Time‑range filtering & historical trend API

🛣️ Mid-Term:

- Anonymous benchmarking exchange
- Streaming ingestion (NATS) & WebSocket updates
- Advanced analytics (service breakdown, SLO overlays)

## Quick Start

### 1. Fast Dev (Recommended for exploring)

Windows PowerShell (uses orchestration script):

```powershell
npm run dev
```

macOS / Linux:

```bash
./scripts/start-dev.sh
```

This launches:

- Backend API: <http://localhost:8080>
- Frontend UI: <http://localhost:5173>

### 2. Manual (backend & frontend separately)

Backend:

```bash
cd backend
go run ./cmd/server
```

Frontend:

```bash
cd frontend
npm install
npm run dev
```

### 3. (Optional) Docker Compose (coming soon)

Dev compose file exists; integration will follow once persistence is wired.

## Simulating Data (In-Memory Mode)

From the dashboard Overview tab:

- "Simulate Incident" → adds active incident
- "Resolve" → closes incident (affects MTTR)
- "Simulate Deployment" → adds success/failure (affects Frequency & CFR)

You can also call endpoints directly (examples via curl):

```bash
# Create deployment
curl -X POST localhost:8080/api/v1/deployments -H "Content-Type: application/json" -d '{"id":"dep-1","service":"api","environment":"prod","status":"success"}'

# Create incident
curl -X POST localhost:8080/api/v1/incidents -H "Content-Type: application/json" -d '{"id":"inc-1","title":"Synthetic Outage","service":"api","environment":"prod","severity":"high"}'

# Resolve incident
curl -X POST localhost:8080/api/v1/incidents/inc-1/resolve

# View current state
curl localhost:8080/api/v1/state | jq
```

## API Surface (Current)

Base path: `/api/v1`

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | System health summary |
| `/health/database` | GET | (Stub) database health |
| `/health/redis` | GET | (Stub) redis health |
| `/metrics/dora` | GET | All DORA metrics + classification |
| `/deployments` | POST | Ingest (simulate) a deployment |
| `/incidents` | POST | Ingest (simulate) an incident |
| `/incidents/:id/resolve` | POST | Resolve an incident |
| `/state` | GET | Snapshot of in‑memory deployments & incidents |
| `/plugins` | GET | Stub plugin listing |
| `/plugins/:name/health` | GET | Stub plugin health |
| `/webhook/:plugin` | POST | Generic plugin webhook receiver (stub) |
### Standard Response Format

All API responses are wrapped to provide consistency and traceability.

Success:

```json
{
	"data": { /* endpoint-specific payload */ },
	"trace_id": "d3f5c8a2d0ad4e33"
}
```

Error:

```json
{
	"error": {
		"code": "validation_error",
		"message": "field 'service' is required",
		"details": { "field": "service" },
		"trace_id": "d3f5c8a2d0ad4e33"
	}
}
```

Error codes currently emitted:

| Code | HTTP | Meaning |
|------|------|---------|
| validation_error | 400 | Invalid input / failed constraints |
| unauthorized | 401 | (Reserved) auth required / invalid token |
| forbidden | 403 | (Reserved) caller lacks permission |
| not_found | 404 | Resource doesn't exist |
| conflict | 409 | State conflict (duplicate, version mismatch) |
| timeout | 504 | Server timed out processing request |
| internal_error | 500 | Unclassified server error |

Each request receives a request_id (trace_id) injected via middleware. Propagate this into logs / tracing for correlation.

### Request Timeouts

The server enforces a per‑request timeout (default 5s). Configure via environment variable:

```
REQUEST_TIMEOUT_SECONDS=8
```

Timed out requests return `504` with `code: "timeout"`.

### Migrations & Persistence

Set `AUTO_MIGRATE=true` (or corresponding config) to apply SQL files in `./migrations` on startup. When `DATABASE_URL` is unset the server transparently falls back to in‑memory slices (useful for quick local demos). Mixing modes is supported: you can start with memory then add a DB without code changes.

Time range query support for `/metrics/dora` (`?days=7|30|90`) will expand as historical persistence arrives.

## Architecture Principles

Clean layered approach influenced by Hexagonal + DDD:

- Domain package (`pkg/metrics`) holds core models & calculations
- API layer (Gin) performs orchestration & response shaping
- In-memory slices currently stand in for repository interfaces
- Frontend consumes typed responses via a small service layer (React Query for caching/polling)

Planned evolutions:

- Replace slices with repository interfaces (Postgres + Redis caching)
- Event-style ingestion (NATS) feeding calculator pipelines
- Plugin runtime applying rate limiting, retries, circuit breakers

## Project Structure (Active Portions)

```text
backend/
	cmd/server        # Entry point
	internal/api      # Routers & handlers
	internal/storage  # (Stubs) future persistence
	pkg/metrics       # Domain models & calculator
frontend/
	src/components    # Dashboard + Incident widget
	src/services      # API client
	src/types         # Shared TS interfaces
scripts/            # Dev orchestration scripts
```

## Development Workflow

Run locally, simulate data, iterate quickly before persistence adds complexity.

Tests (initial calculator tests exist; broaden soon):

```bash
cd backend && go test ./...
```

Frontend tests placeholder (enable after component set stabilizes).

## Roadmap (Snapshot)

- [x] Core DORA calculator
- [x] Performance classification (overall + per metric)
- [x] Incident + deployment ingestion (in-memory)
- [x] Interactive dashboard + dark mode
- [ ] Persistence layer (Postgres + migrations)
- [ ] Plugin baseline (GitHub → deployments, incidents)
- [ ] Auth (JWT + RBAC)
- [ ] Historical trend endpoints & charts
- [ ] Benchmark exchange (anonymous percentile data)
- [ ] Observability: traces + metrics + structured logs
- [ ] Export / CSV / API tokens

Extended multi-phase plan lives in `DEVELOPMENT_PLAN.md`.

## Contributing

Pre‑alpha: architecture still evolving. Early feedback & lightweight PRs welcome.

1. Fork & branch (`feat/xyz`)
2. Keep changes focused (small PRs merge faster)
3. Add / update tests where logic changes
4. Submit PR & include context / rationale

Roadmap issues will be labeled as they open—watch the repo to get notified.

## License

Apache 2.0 – see [LICENSE](LICENSE).

## Acknowledgements

- DORA Research (<https://dora.dev>)
- OpenTelemetry & wider CNCF ecosystem
- All upstream OSS dependencies

---
If this direction resonates, ⭐ star & follow along — more coming soon.
