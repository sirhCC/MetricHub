package storage

import (
    "context"
    "database/sql"
    "fmt"
    "time"
    "github.com/sirhCC/MetricHub/pkg/metrics"
)

// DeploymentRepository defines persistence for deployments.
type DeploymentRepository interface {
    Create(ctx context.Context, d *metrics.Deployment) error
    ListRange(ctx context.Context, start, end time.Time) ([]metrics.Deployment, error)
}

// IncidentRepository defines persistence for incidents.
type IncidentRepository interface {
    Create(ctx context.Context, i *metrics.Incident) error
    Resolve(ctx context.Context, id string, resolvedAt time.Time) error
    ListRange(ctx context.Context, start, end time.Time) ([]metrics.Incident, error)
}

// PostgresDeploymentRepo implements DeploymentRepository.
type PostgresDeploymentRepo struct { db *sql.DB }

func NewPostgresDeploymentRepo(db *sql.DB) *PostgresDeploymentRepo { return &PostgresDeploymentRepo{db: db} }

func (r *PostgresDeploymentRepo) Create(ctx context.Context, d *metrics.Deployment) error {
    const q = `INSERT INTO deployments (id, service, environment, version, status, start_time, end_time, commit_sha, commit_time, author, repository, branch, build_url, tags, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`
    _, err := r.db.ExecContext(ctx, q,
        d.ID, d.Service, d.Environment, d.Version, d.Status, d.StartTime, d.EndTime, d.CommitSHA, d.CommitTime,
        d.Author, d.Repository, d.Branch, d.BuildURL, nil, d.CreatedAt, d.UpdatedAt,
    )
    return err
}

func (r *PostgresDeploymentRepo) ListRange(ctx context.Context, start, end time.Time) ([]metrics.Deployment, error) {
    const q = `SELECT id, service, environment, version, status, start_time, end_time, commit_sha, commit_time, author, repository, branch, build_url, created_at, updated_at FROM deployments
WHERE start_time BETWEEN $1 AND $2 ORDER BY start_time`
    rows, err := r.db.QueryContext(ctx, q, start, end)
    if err != nil { return nil, err }
    defer rows.Close()
    var out []metrics.Deployment
    for rows.Next() {
        var d metrics.Deployment
        if err := rows.Scan(&d.ID, &d.Service, &d.Environment, &d.Version, &d.Status, &d.StartTime, &d.EndTime, &d.CommitSHA, &d.CommitTime, &d.Author, &d.Repository, &d.Branch, &d.BuildURL, &d.CreatedAt, &d.UpdatedAt); err != nil { return nil, err }
        out = append(out, d)
    }
    return out, rows.Err()
}

// PostgresIncidentRepo implements IncidentRepository.
type PostgresIncidentRepo struct { db *sql.DB }

func NewPostgresIncidentRepo(db *sql.DB) *PostgresIncidentRepo { return &PostgresIncidentRepo{db: db} }

func (r *PostgresIncidentRepo) Create(ctx context.Context, i *metrics.Incident) error {
    const q = `INSERT INTO incidents (id, title, description, service, environment, severity, start_time, resolved_time, root_cause, assignee, tags, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`
    _, err := r.db.ExecContext(ctx, q, i.ID, i.Title, i.Description, i.Service, i.Environment, i.Severity, i.StartTime, i.ResolvedTime, i.RootCause, i.Assignee, nil, i.CreatedAt, i.UpdatedAt)
    return err
}

func (r *PostgresIncidentRepo) Resolve(ctx context.Context, id string, resolvedAt time.Time) error {
    const q = `UPDATE incidents SET resolved_time=$2, updated_at=$2 WHERE id=$1 AND resolved_time IS NULL`
    res, err := r.db.ExecContext(ctx, q, id, resolvedAt)
    if err != nil { return err }
    n, _ := res.RowsAffected()
    if n == 0 { return fmt.Errorf("incident not found or already resolved") }
    return nil
}

func (r *PostgresIncidentRepo) ListRange(ctx context.Context, start, end time.Time) ([]metrics.Incident, error) {
    const q = `SELECT id, title, description, service, environment, severity, start_time, resolved_time, root_cause, assignee, created_at, updated_at FROM incidents
WHERE start_time BETWEEN $1 AND $2 ORDER BY start_time`
    rows, err := r.db.QueryContext(ctx, q, start, end)
    if err != nil { return nil, err }
    defer rows.Close()
    var out []metrics.Incident
    for rows.Next() {
        var inc metrics.Incident
        if err := rows.Scan(&inc.ID, &inc.Title, &inc.Description, &inc.Service, &inc.Environment, &inc.Severity, &inc.StartTime, &inc.ResolvedTime, &inc.RootCause, &inc.Assignee, &inc.CreatedAt, &inc.UpdatedAt); err != nil { return nil, err }
        out = append(out, inc)
    }
    return out, rows.Err()
}
