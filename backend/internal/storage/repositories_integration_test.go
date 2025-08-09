//go:build integration
// +build integration

package storage_test

import (
    "context"
    "database/sql"
    "fmt"
    "os/exec"
    "testing"
    "time"

    _ "github.com/lib/pq"
    "github.com/stretchr/testify/require"
    tc "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"

    "github.com/sirhCC/MetricHub/internal/storage"
    "github.com/sirhCC/MetricHub/pkg/metrics"
)

// withTestPostgres spins up a Postgres container and returns a cleanup func.
func withTestPostgres(t *testing.T) (*sql.DB, func()) {
    t.Helper()
    // Pre-flight: ensure docker is available
    if err := exec.Command("docker", "version").Run(); err != nil {
        t.Skipf("docker not available: %v", err)
    }
    ctx := context.Background()
    pg, err := postgres.RunContainer(ctx,
        postgres.WithDatabase("metrichub"),
        postgres.WithUsername("metrichub"),
        postgres.WithPassword("password"),
        // Migration script path relative to module root (go test runs from module root)
        postgres.WithInitScripts("migrations/0001_init_schema.up.sql"),
        tc.WithImage("postgres:15-alpine"),
    )
    require.NoError(t, err)

    host, err := pg.Host(ctx)
    require.NoError(t, err)
    port, err := pg.MappedPort(ctx, "5432")
    require.NoError(t, err)

    dsn := fmt.Sprintf("postgres://metrichub:password@%s:%s/metrichub?sslmode=disable", host, port.Port())
    db, err := sql.Open("postgres", dsn)
    require.NoError(t, err)
    require.Eventually(t, func() bool { return db.Ping() == nil }, 10*time.Second, 500*time.Millisecond)

    cleanup := func() {
        _ = db.Close()
        _ = pg.Terminate(ctx)
    }
    return db, cleanup
}

func TestPostgresDeploymentRepository_CreateAndListRange(t *testing.T) {
    db, cleanup := withTestPostgres(t)
    defer cleanup()
    repo := storage.NewPostgresDeploymentRepo(db)

    now := time.Now().Add(-2 * time.Hour)
    dep := &metrics.Deployment{
        ID:          "dep-1",
        Service:     "api",
        Environment: "dev",
        Version:     "1.0.0",
        Status:      metrics.DeploymentStatusSuccess,
        StartTime:   now,
        EndTime:     ptrTime(now.Add(5 * time.Minute)),
        CommitSHA:   "abc123",
        CommitTime:  now.Add(-10 * time.Minute),
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    err := repo.Create(context.Background(), dep)
    require.NoError(t, err)

    // second deployment outside range
    dep2 := &metrics.Deployment{
        ID:          "dep-2",
        Service:     "api",
        Environment: "dev",
        Version:     "1.0.1",
        Status:      metrics.DeploymentStatusFailed,
        StartTime:   now.Add(-48 * time.Hour),
        CommitSHA:   "def456",
        CommitTime:  now.Add(-49 * time.Hour),
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    require.NoError(t, repo.Create(context.Background(), dep2))

    list, err := repo.ListRange(context.Background(), now.Add(-1*time.Hour), time.Now())
    require.NoError(t, err)
    require.Len(t, list, 1)
    require.Equal(t, "dep-1", list[0].ID)
}

func TestPostgresIncidentRepository_CreateResolveListRange(t *testing.T) {
    db, cleanup := withTestPostgres(t)
    defer cleanup()
    repo := storage.NewPostgresIncidentRepo(db)

    start := time.Now().Add(-3 * time.Hour)
    inc := &metrics.Incident{
        ID:          "inc-1",
        Title:       "Outage",
        Service:     "api",
        Environment: "dev",
        Severity:    metrics.SeverityHigh,
        StartTime:   start,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    require.NoError(t, repo.Create(context.Background(), inc))

    // Incident outside range
    oldInc := &metrics.Incident{
        ID:          "inc-2",
        Title:       "Historical",
        Service:     "api",
        Environment: "dev",
        Severity:    metrics.SeverityLow,
        StartTime:   start.Add(-72 * time.Hour),
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    require.NoError(t, repo.Create(context.Background(), oldInc))

    // Resolve first
    require.NoError(t, repo.Resolve(context.Background(), "inc-1", time.Now()))

    list, err := repo.ListRange(context.Background(), start.Add(-1*time.Hour), time.Now())
    require.NoError(t, err)
    require.Len(t, list, 1)
    require.Equal(t, "inc-1", list[0].ID)
    require.True(t, list[0].IsResolved())
}

func ptrTime(t time.Time) *time.Time { return &t }
