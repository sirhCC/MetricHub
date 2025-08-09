package storage

import (
    "context"
    "database/sql"
    "fmt"
    "io/fs"
    "os"
    "path/filepath"
    "sort"
    "strings"
    "time"
)

// ApplyMigrations scans a directory for *.up.sql files and applies any that have not yet run.
// It records applied migrations in schema_migrations. This is a lightweight interim solution
// until a full migration tool is adopted.
func ApplyMigrations(ctx context.Context, db *sql.DB, dir string, logger func(msg string, kv ...interface{})) error {
    if logger == nil { logger = func(string, ...interface{}){} }
    if _, err := os.Stat(dir); err != nil { return fmt.Errorf("migrations directory: %w", err) }
    if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT now())`); err != nil {
        return fmt.Errorf("create schema_migrations: %w", err)
    }
    applied, err := loadApplied(ctx, db)
    if err != nil { return err }
    var files []string
    err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, walkErr error) error {
        if walkErr != nil { return walkErr }
        if d.IsDir() { return nil }
        if strings.HasSuffix(d.Name(), ".up.sql") { files = append(files, path) }
        return nil
    })
    if err != nil { return fmt.Errorf("scan migrations: %w", err) }
    sort.Strings(files)
    for _, f := range files {
        base := filepath.Base(f)
        versionKey := strings.TrimSuffix(base, ".up.sql")
        if applied[versionKey] { continue }
        sqlBytes, readErr := os.ReadFile(f)
        if readErr != nil { return fmt.Errorf("read %s: %w", f, readErr) }
        logger("applying migration", "file", base)
        tx, txErr := db.BeginTx(ctx, nil)
        if txErr != nil { return fmt.Errorf("begin tx %s: %w", base, txErr) }
        if _, execErr := tx.ExecContext(ctx, string(sqlBytes)); execErr != nil { _ = tx.Rollback(); return fmt.Errorf("execute %s: %w", base, execErr) }
        if _, insErr := tx.ExecContext(ctx, `INSERT INTO schema_migrations (version, applied_at) VALUES ($1,$2)`, versionKey, time.Now()); insErr != nil { _ = tx.Rollback(); return fmt.Errorf("record %s: %w", base, insErr) }
        if commitErr := tx.Commit(); commitErr != nil { return fmt.Errorf("commit %s: %w", base, commitErr) }
        logger("migration applied", "file", base)
    }
    return nil
}

func loadApplied(ctx context.Context, db *sql.DB) (map[string]bool, error) {
    rows, err := db.QueryContext(ctx, `SELECT version FROM schema_migrations`)
    if err != nil { return nil, fmt.Errorf("select schema_migrations: %w", err) }
    defer rows.Close()
    applied := make(map[string]bool)
    for rows.Next() { var v string; if err := rows.Scan(&v); err != nil { return nil, err }; applied[v] = true }
    return applied, rows.Err()
}
