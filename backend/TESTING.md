# Backend Testing

## Test Types

- Unit tests: Fast, pure Go (default `go test ./...`).
- Integration tests: Require Docker; use build tag `integration`.

## Running Integration Tests

Ensure Docker is running, then execute:

```powershell
# Only storage integration tests
go test -tags=integration ./internal/storage -run TestPostgres -count=1

# All tests including integration
go test -tags=integration ./...
```

If Docker isn't available on the host (e.g., certain Windows environments), integration tests will be skipped automatically.

## Adding New Integration Tests

1. Place them in an appropriate package (e.g., `internal/storage`).
1. Add the build constraint at the top of the file:

```go
//go:build integration
// +build integration
```

1. Gate external dependencies (Docker, network) with skip checks where feasible.

## Roadmap

Future enhancements will introduce:

- Mockable interfaces for repositories (unit test coverage without containers)
- Testcontainers modules for Redis and NATS
- Metrics calculation end-to-end scenario tests
