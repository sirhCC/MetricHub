// Mock API Server for MetricHub Frontend Development
// This simulates the Go backend API responses

const express = require('express');
const cors = require('cors');

const app = express();
const port = 8080;

// Middleware
app.use(cors({
  origin: ['http://localhost:3000', 'http://localhost:5173', 'http://localhost:5174', 'http://localhost:5175'],
  credentials: true
}));
app.use(express.json());

// Mock data
const mockDoraMetrics = {
  deployment_frequency: 2.5,
  lead_time: "2h 30m",
  mttr: "45m",
  change_failure_rate: 0.15,
  time_range: "last-30-days",
  last_updated: new Date().toISOString()
};

const mockHealthData = {
  status: "healthy",
  timestamp: new Date().toISOString(),
  version: "0.1.0",
  checks: {
    database: "healthy",
    redis: "healthy"
  }
};

const mockPlugins = [
  {
    name: "github",
    version: "1.0.0",
    description: "GitHub Actions and repository metrics collector",
    status: "healthy",
    enabled: true
  },
  {
    name: "jenkins",
    version: "1.0.0", 
    description: "Jenkins CI/CD metrics collector",
    status: "healthy",
    enabled: false
  },
  {
    name: "gitlab",
    version: "1.0.0",
    description: "GitLab CI/CD metrics collector", 
    status: "unknown",
    enabled: false
  }
];

// Routes

// Health endpoints
app.get('/api/v1/health', (req, res) => {
  res.json(mockHealthData);
});

app.get('/api/v1/health/database', (req, res) => {
  res.json({
    status: "healthy",
    timestamp: new Date().toISOString(),
    checks: { database: "healthy" }
  });
});

app.get('/api/v1/health/redis', (req, res) => {
  res.json({
    status: "healthy", 
    timestamp: new Date().toISOString(),
    checks: { redis: "healthy" }
  });
});

// DORA metrics endpoints
app.get('/api/v1/metrics/dora', (req, res) => {
  res.json({
    data: mockDoraMetrics,
    metadata: {
      version: "v1",
      total_count: 1,
      time_range: "last-30-days",
      data_quality: "high",
      last_updated: new Date().toISOString()
    }
  });
});

app.get('/api/v1/metrics/dora/deployment-frequency', (req, res) => {
  res.json({
    data: {
      value: 2.5,
      unit: "per_day", 
      trend: "increasing",
      benchmark: {
        percentile: 75,
        industry: "technology"
      }
    },
    metadata: {
      time_range: "last-30-days",
      data_quality: "high", 
      last_updated: new Date().toISOString()
    }
  });
});

app.get('/api/v1/metrics/dora/lead-time', (req, res) => {
  res.json({
    data: {
      value: "2h 30m",
      value_ms: 9000000,
      trend: "decreasing",
      benchmark: {
        percentile: 80,
        industry: "technology"
      }
    },
    metadata: {
      time_range: "last-30-days",
      data_quality: "high",
      last_updated: new Date().toISOString()
    }
  });
});

app.get('/api/v1/metrics/dora/mttr', (req, res) => {
  res.json({
    data: {
      value: "45m",
      value_ms: 2700000,
      trend: "stable",
      benchmark: {
        percentile: 70,
        industry: "technology"
      }
    },
    metadata: {
      time_range: "last-30-days",
      data_quality: "medium",
      last_updated: new Date().toISOString()
    }
  });
});

app.get('/api/v1/metrics/dora/change-failure-rate', (req, res) => {
  res.json({
    data: {
      value: 0.15,
      percentage: "15%",
      trend: "decreasing",
      benchmark: {
        percentile: 65,
        industry: "technology"
      }
    },
    metadata: {
      time_range: "last-30-days",
      data_quality: "high",
      last_updated: new Date().toISOString()
    }
  });
});

// Plugin endpoints
app.get('/api/v1/plugins', (req, res) => {
  res.json({
    data: mockPlugins,
    metadata: {
      total_count: mockPlugins.length,
      enabled_count: mockPlugins.filter(p => p.enabled).length,
      disabled_count: mockPlugins.filter(p => !p.enabled).length
    }
  });
});

app.get('/api/v1/plugins/:name/health', (req, res) => {
  const pluginName = req.params.name;
  const plugin = mockPlugins.find(p => p.name === pluginName);
  
  if (!plugin) {
    return res.status(404).json({
      data: {
        plugin: pluginName,
        status: "unknown",
        message: "Plugin not found",
        last_check: new Date().toISOString()
      }
    });
  }

  res.json({
    data: {
      plugin: pluginName,
      status: plugin.status,
      message: plugin.status === "healthy" ? "Plugin is running normally" : "Plugin is experiencing issues",
      last_check: new Date().toISOString()
    }
  });
});

// Webhook endpoint
app.post('/api/v1/webhook/:plugin', (req, res) => {
  const pluginName = req.params.plugin;
  console.log(`Received webhook for plugin: ${pluginName}`, req.body);
  
  res.json({
    message: "Webhook received successfully",
    plugin: pluginName,
    processed_at: new Date().toISOString()
  });
});

// Start server
app.listen(port, () => {
  console.log(`🚀 MetricHub Mock API Server running at http://localhost:${port}`);
  console.log(`📊 DORA metrics available at http://localhost:${port}/api/v1/metrics/dora`);
  console.log(`💚 Health check at http://localhost:${port}/api/v1/health`);
  console.log('');
  console.log('🎯 This is a mock server for frontend development.');
  console.log('👉 Install Go and run the real backend when ready!');
});
