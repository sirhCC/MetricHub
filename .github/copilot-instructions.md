# MetricHub AI Agent Instructions

<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

## Project Overview

MetricHub is a professional-grade, open-source Universal DevOps Metrics Collector that aggregates DORA metrics from multiple tools into unified dashboards with community benchmarking capabilities. This project follows enterprise-level standards and best practices throughout.

## Architecture Philosophy

### Clean Architecture Principles
- **Hexagonal Architecture**: Core business logic isolated from external concerns
- **Domain-Driven Design**: Rich domain models with clear boundaries
- **Dependency Inversion**: Abstractions don't depend on details
- **Single Responsibility**: Each component has one reason to change

### Technology Stack Rationale
```
Backend:    Go 1.21+ (performance, concurrency, single binary distribution)
Frontend:   React 18 + TypeScript + Vite (type safety, modern tooling)
Database:   PostgreSQL + Redis (ACID compliance + high-performance caching)
Queue:      NATS (lightweight, high-throughput messaging)
Observability: OpenTelemetry + Prometheus + Jaeger (industry standards)
```

## Code Quality Standards

### Go Backend Conventions
```go
// Package structure follows Go standards
package collector

// Interfaces are small and focused
type MetricsCollector interface {
    Collect(ctx context.Context, config Config) (*MetricsData, error)
    HealthCheck(ctx context.Context) error
}

// Error handling is explicit and contextual
func (c *Collector) Process(ctx context.Context) error {
    if err := c.validate(); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    // Implementation...
}

// Structured logging with context
log.Info("processing metrics",
    "plugin", pluginName,
    "duration", processingTime,
    "records", recordCount)
```

### TypeScript Frontend Conventions
```typescript
// Props interfaces are explicit and well-documented
interface MetricsDashboardProps {
  timeRange: TimeRange;
  services: string[];
  onServiceSelect: (service: string) => void;
}

// Custom hooks for business logic
function useDoraMetrics(filters: MetricsFilters) {
  const { data, loading, error } = useQuery(DORA_METRICS_QUERY, {
    variables: filters,
    pollInterval: 30000, // Real-time updates
  });
  
  return { metrics: data?.doraMetrics, loading, error };
}

// Error boundaries for fault tolerance
class MetricsErrorBoundary extends Component<Props, State> {
  // Implementation with proper error reporting
}
```

## Plugin System Architecture

### Plugin Interface Contract
```go
// All plugins must implement this interface
type MetricsPlugin interface {
    // Metadata
    Name() string
    Version() string
    Description() string
    
    // Lifecycle
    Initialize(config PluginConfig) error
    Shutdown(ctx context.Context) error
    
    // Core functionality
    Collect(ctx context.Context) (*MetricsData, error)
    Validate(config PluginConfig) error
    HealthCheck(ctx context.Context) error
}

// Plugin configuration is type-safe
type GitHubConfig struct {
    Token      string        `json:"token" validate:"required"`
    Owner      string        `json:"owner" validate:"required"`
    Repository string        `json:"repository"`
    RateLimit  time.Duration `json:"rate_limit" default:"1s"`
}
```

### Plugin Development Patterns
- **Circuit Breaker**: Prevent cascade failures when external APIs are down
- **Exponential Backoff**: Graceful retry logic for transient failures
- **Rate Limiting**: Respect external API quotas and limits
- **Graceful Degradation**: Continue working when some plugins fail
- **Health Monitoring**: Continuous plugin health assessment

## DORA Metrics Implementation

### Calculation Accuracy
```go
// Deployment Frequency: Deployments per day/week/month
func (c *Calculator) DeploymentFrequency(deployments []Deployment, timeRange TimeRange) float64 {
    successful := filter(deployments, func(d Deployment) bool {
        return d.Status == StatusSuccess && d.Timestamp.In(timeRange)
    })
    return float64(len(successful)) / timeRange.Days()
}

// Lead Time: Time from commit to production
func (c *Calculator) LeadTimeForChanges(deployments []Deployment) time.Duration {
    var totalLeadTime time.Duration
    for _, deployment := range deployments {
        leadTime := deployment.ProductionTime.Sub(deployment.CommitTime)
        totalLeadTime += leadTime
    }
    return totalLeadTime / time.Duration(len(deployments))
}
```

### Data Quality Assurance
- **Idempotent Processing**: Same input always produces same output
- **Data Validation**: Multi-layer validation (input, business rules, output)
- **Audit Trails**: Event sourcing for complete change history
- **Data Lineage**: Track data flow from source to metrics

## API Design Principles

### RESTful Conventions
```
GET    /api/v1/metrics/dora                    # DORA metrics summary
GET    /api/v1/metrics/dora/deployment-frequency
GET    /api/v1/metrics/dora/lead-time
POST   /api/v1/collectors/{plugin}/webhook     # Plugin-specific webhooks
GET    /api/v1/health                          # System health
GET    /api/v1/plugins                         # Available plugins
```

### Response Formats
```json
{
  "data": {
    "deployment_frequency": {
      "value": 2.5,
      "unit": "per_day",
      "trend": "increasing",
      "benchmark_percentile": 75
    }
  },
  "metadata": {
    "time_range": "2024-01-01T00:00:00Z/2024-01-31T23:59:59Z",
    "data_quality": "high",
    "last_updated": "2024-01-31T12:00:00Z"
  }
}
```

## Testing Strategy

### Test Pyramid Implementation
```go
// Unit Tests: Fast, isolated, comprehensive
func TestDeploymentFrequencyCalculation(t *testing.T) {
    deployments := []Deployment{
        {Status: StatusSuccess, Timestamp: time.Now().Add(-24 * time.Hour)},
        {Status: StatusSuccess, Timestamp: time.Now().Add(-12 * time.Hour)},
    }
    
    frequency := calculator.DeploymentFrequency(deployments, Last7Days())
    assert.Equal(t, 2.0/7.0, frequency)
}

// Integration Tests: Real database, external services mocked
func TestGitHubPluginIntegration(t *testing.T) {
    // Use testcontainers for real database
    db := setupTestDatabase(t)
    defer db.Close()
    
    // Mock GitHub API responses
    server := httptest.NewServer(githubMockHandler)
    defer server.Close()
    
    plugin := &GitHubPlugin{client: github.NewClient(server.URL)}
    // Test integration...
}
```

### Frontend Testing Patterns
```typescript
// Component Testing with React Testing Library
test('displays DORA metrics correctly', async () => {
  const mockMetrics = { deploymentFrequency: 2.5, leadTime: '2h 30m' };
  
  render(<DoraMetricsCard metrics={mockMetrics} />);
  
  expect(screen.getByText('2.5 per day')).toBeInTheDocument();
  expect(screen.getByText('2h 30m')).toBeInTheDocument();
});

// E2E Testing with Playwright
test('user can filter metrics by time range', async ({ page }) => {
  await page.goto('/dashboard');
  await page.selectOption('[data-testid=time-range]', 'last-30-days');
  
  await expect(page.locator('[data-testid=metrics-chart]')).toContainText('Last 30 days');
});
```

## Performance Optimization

### Backend Performance
- **Database Indexing**: Optimize for time-series queries
- **Connection Pooling**: Efficient database connection management
- **Caching Strategy**: Redis for frequently accessed data
- **Batch Processing**: Process metrics in batches for efficiency

### Frontend Performance
- **Code Splitting**: Lazy load dashboard components
- **Memoization**: React.memo for expensive components
- **Virtual Scrolling**: Handle large datasets efficiently
- **Real-time Updates**: WebSocket connections with optimistic updates

## Security Implementation

### Authentication & Authorization
```go
// JWT-based authentication with proper validation
func (a *AuthMiddleware) ValidateToken(token string) (*Claims, error) {
    claims := &Claims{}
    tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
        return a.jwtKey, nil
    })
    
    if err != nil || !tkn.Valid {
        return nil, ErrInvalidToken
    }
    
    return claims, nil
}

// Role-based access control
func (h *Handler) requireRole(requiredRole Role) middleware.Func {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            claims := getClaimsFromContext(r.Context())
            if !claims.HasRole(requiredRole) {
                http.Error(w, "Forbidden", http.StatusForbidden)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

## Observability Requirements

### Metrics to Track
- **Business Metrics**: Deployments processed, APIs called, plugins active
- **Technical Metrics**: Response times, error rates, throughput
- **Infrastructure Metrics**: CPU, memory, disk usage, network I/O

### Logging Standards
```go
// Structured logging with consistent fields
logger.Info("plugin execution completed",
    zap.String("plugin", pluginName),
    zap.Duration("duration", executionTime),
    zap.Int("records_processed", recordCount),
    zap.String("trace_id", traceID))
```

## Documentation Requirements

### Code Documentation
- **Package Documentation**: Clear purpose and usage examples
- **Function Documentation**: Parameters, return values, error conditions
- **Architecture Decision Records**: Document significant design decisions
- **API Documentation**: OpenAPI 3.0 specification with examples

### User Documentation
- **Installation Guides**: Docker, Kubernetes, binary installation
- **Configuration Reference**: All configuration options documented
- **Plugin Development Guide**: How to create new plugins
- **Troubleshooting Guide**: Common issues and solutions

## Contribution Guidelines

### Code Review Standards
- All code changes require peer review
- Automated checks must pass (tests, linting, security)
- Documentation updates for user-facing changes
- Performance impact assessment for critical paths

### Release Process
- Semantic versioning (MAJOR.MINOR.PATCH)
- Automated release notes generation
- Security scanning before release
- Staged rollout strategy

## Success Metrics

### Technical Excellence
- **Test Coverage**: >95% for backend, >90% for frontend
- **Performance**: <200ms API response times, <2s page loads
- **Reliability**: 99.9% uptime, graceful degradation
- **Security**: Zero critical vulnerabilities, regular audits

### Community Growth
- **Adoption**: 1000+ GitHub stars, 100+ organizations using
- **Contributions**: 50+ contributors, 20+ community plugins
- **Documentation**: Comprehensive guides, video tutorials
- **Support**: Active community forums, responsive issue resolution

Follow these guidelines to ensure MetricHub becomes the industry standard for DevOps metrics collection while maintaining the highest standards of code quality and community engagement.
