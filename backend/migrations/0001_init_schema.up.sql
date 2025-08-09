-- Deployments table
CREATE TABLE IF NOT EXISTS deployments (
  id TEXT PRIMARY KEY,
  service TEXT NOT NULL,
  environment TEXT NOT NULL,
  version TEXT,
  status TEXT NOT NULL,
  start_time TIMESTAMPTZ NOT NULL,
  end_time TIMESTAMPTZ,
  commit_sha TEXT,
  commit_time TIMESTAMPTZ NOT NULL,
  author TEXT,
  repository TEXT,
  branch TEXT,
  build_url TEXT,
  tags JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_deployments_start_time ON deployments(start_time);
CREATE INDEX IF NOT EXISTS idx_deployments_service ON deployments(service);

-- Incidents table
CREATE TABLE IF NOT EXISTS incidents (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT,
  service TEXT NOT NULL,
  environment TEXT NOT NULL,
  severity TEXT NOT NULL,
  start_time TIMESTAMPTZ NOT NULL,
  resolved_time TIMESTAMPTZ,
  root_cause TEXT,
  assignee TEXT,
  tags JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_incidents_start_time ON incidents(start_time);
CREATE INDEX IF NOT EXISTS idx_incidents_service ON incidents(service);
