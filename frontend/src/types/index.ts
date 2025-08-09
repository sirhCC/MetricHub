// DORA Metrics Types - TypeScript definitions for MetricHub frontend
// Matching the Go backend domain models

export interface TimeRange {
  start: string;
  end: string;
}

export interface DORAMetrics {
  deployment_frequency: number;
  lead_time: string;           // duration string from backend
  mttr: string;                // duration string
  change_failure_rate: number; // decimal (0.15 = 15%)
  classification?: PerformanceClassification;
  overall_performance?: PerformanceLevel;
}

export interface MetricData {
  value: number | string;
  unit?: string;
  trend?: 'increasing' | 'decreasing' | 'stable';
  benchmark?: {
    percentile: number;
    industry: string;
  };
}

export interface MetricResponse<T = MetricData> {
  data: T;
  metadata: {
    time_range: string;
    data_quality: 'high' | 'medium' | 'low';
    last_updated: string;
    version?: string;
    total_count?: number;
  };
}

export interface Plugin {
  name: string;
  version: string;
  description: string;
  status: 'healthy' | 'unhealthy' | 'unknown';
  enabled: boolean;
}

export interface PluginsResponse {
  data: Plugin[];
  metadata: {
    total_count: number;
    enabled_count: number;
    disabled_count: number;
  };
}

export interface HealthCheck {
  status: 'healthy' | 'unhealthy';
  timestamp: string;
  version?: string;
  checks?: Record<string, string>;
}

export type PerformanceLevel = 'elite' | 'high' | 'medium' | 'low';

export interface PerformanceClassification {
  deployment_frequency: PerformanceLevel;
  lead_time: PerformanceLevel;
  mttr: PerformanceLevel;
  change_failure_rate: PerformanceLevel;
  overall: PerformanceLevel;
}

// Deployment & Incident domain (subset of backend models for UI needs)
export interface Deployment {
  id: string;
  service: string;
  environment: string;
  status: string;
  start_time: string;
  end_time?: string;
  commit_sha?: string;
  commit_time?: string;
  created_at?: string;
  updated_at?: string;
}

export interface Incident {
  id: string;
  title: string;
  description?: string;
  service: string;
  environment: string;
  severity: string; // low|medium|high|critical
  start_time: string;
  resolved_time?: string;
  created_at?: string;
  updated_at?: string;
}

export interface SystemState {
  deployments: Deployment[];
  incidents: Incident[];
}

// UI State Types
export interface DashboardFilters {
  timeRange: 'last-7-days' | 'last-30-days' | 'last-90-days' | 'this-month' | 'last-month' | 'custom';
  services: string[];
  environments: string[];
  customTimeRange?: TimeRange;
}

export interface ChartDataPoint {
  date: string;
  value: number;
  label?: string;
}

export interface TrendData {
  current: number;
  previous: number;
  change: number;
  trend: 'up' | 'down' | 'stable';
}

// API Error Types
export interface APIError {
  error: string;
  details?: string;
  code?: string;
}

// Webhook Types
export interface WebhookPayload {
  plugin: string;
  event_type: string;
  data: Record<string, unknown>;
  timestamp: string;
}

// Configuration Types
export interface AppConfig {
  apiUrl: string;
  version: string;
  environment: 'development' | 'production';
  features: {
    realTimeUpdates: boolean;
    benchmarking: boolean;
    exportData: boolean;
  };
}
