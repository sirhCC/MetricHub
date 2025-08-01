# MetricHub

> **Universal DevOps Metrics Collector** - Aggregate DORA metrics from any tool into unified dashboards with community benchmarking.

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)
[![React](https://img.shields.io/badge/React-18+-61DAFB?logo=react)](https://reactjs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5+-3178C6?logo=typescript)](https://www.typescriptlang.org/)

## 🎯 What is MetricHub?

MetricHub is an open-source platform that collects DevOps metrics from multiple tools and presents them in unified dashboards. It calculates **DORA metrics** (Deployment Frequency, Lead Time for Changes, Mean Time to Recovery, Change Failure Rate) and provides **anonymous community benchmarking**.

### Key Features

- 🔌 **Universal Integrations** - Works with GitHub Actions, Jenkins, GitLab, Azure DevOps, and more
- 📊 **DORA Metrics** - Accurate calculation of industry-standard DevOps performance metrics
- 🏆 **Community Benchmarking** - Compare your team's performance anonymously
- 🚀 **Easy Deployment** - Single binary, Docker containers, or Kubernetes
- 🔒 **Privacy-First** - Your data stays yours, only anonymized benchmarks are shared
- 📈 **Beautiful Dashboards** - Real-time visualizations and historical trends

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Data Sources  │    │   MetricHub     │    │   Dashboards    │
│                 │    │     Core        │    │                 │
│ • GitHub Actions│◄──►│                 │◄──►│ • Web Dashboard │
│ • Jenkins       │    │ • Collectors    │    │ • API Clients   │
│ • GitLab CI     │    │ • Calculators   │    │ • Mobile Apps   │
│ • Azure DevOps  │    │ • Storage       │    │ • Integrations  │
│ • Custom APIs   │    │ • Benchmarking  │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🚀 Quick Start

### Prerequisites

- **Docker & Docker Compose** (recommended)
- **OR** Go 1.21+ and Node.js 18+ for local development

### 1. Using Docker (Recommended)

```bash
# Clone the repository
git clone https://github.com/sirhCC/MetricHub.git
cd MetricHub

# Start all services
docker-compose -f docker-compose.dev.yml up

# Access the dashboard
open http://localhost:3000
```

### 2. Local Development

```bash
# Start infrastructure services
docker-compose -f docker-compose.dev.yml up postgres redis nats

# Terminal 1: Start backend
cd backend
go run cmd/server/main.go

# Terminal 2: Start frontend
cd frontend
npm run dev
```

## 📊 DORA Metrics

MetricHub calculates the four key DORA metrics:

| Metric | Description | Elite Performance |
|--------|-------------|-------------------|
| **Deployment Frequency** | How often code is deployed to production | Multiple times per day |
| **Lead Time for Changes** | Time from commit to production | Less than 1 day |
| **Mean Time to Recovery** | Time to recover from incidents | Less than 1 hour |
| **Change Failure Rate** | Percentage of deployments causing failures | 0-15% |

## 🔌 Supported Integrations

### CI/CD Platforms
- ✅ GitHub Actions
- ✅ GitLab CI/CD
- ✅ Jenkins
- ✅ Azure DevOps
- 🚧 CircleCI (coming soon)
- 🚧 Tekton (coming soon)

### Incident Management
- ✅ PagerDuty
- ✅ Opsgenie
- 🚧 Incident.io (coming soon)

### Monitoring
- ✅ Prometheus
- ✅ Grafana
- 🚧 DataDog (coming soon)

### Version Control
- ✅ GitHub
- ✅ GitLab
- ✅ Bitbucket
- 🚧 Azure Repos (coming soon)

## 🛠️ Development

### Project Structure

```
MetricHub/
├── backend/                 # Go backend
│   ├── cmd/                # Application entry points
│   ├── internal/           # Private application code
│   │   ├── api/           # REST API handlers
│   │   ├── collector/     # Metrics collection engine
│   │   ├── plugins/       # Plugin system
│   │   ├── storage/       # Database layer
│   │   └── config/        # Configuration
│   └── pkg/               # Public packages
├── frontend/              # React TypeScript frontend
│   ├── src/
│   │   ├── components/    # UI components
│   │   ├── pages/        # Route components
│   │   ├── hooks/        # Custom hooks
│   │   └── services/     # API clients
├── plugins/              # External plugins
├── docs/                # Documentation
└── deployments/         # Deployment configs
```

### Technology Stack

- **Backend**: Go 1.21, Gin, PostgreSQL, Redis, NATS
- **Frontend**: React 18, TypeScript, Vite, TailwindCSS
- **Testing**: Go testing, React Testing Library, Playwright
- **Observability**: OpenTelemetry, Prometheus, Jaeger
- **Deployment**: Docker, Kubernetes, Helm

### Running Tests

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
npm test

# E2E tests
npm run test:e2e
```

## 📚 Documentation

- [📖 User Guide](docs/user-guide.md)
- [🔧 Configuration Reference](docs/configuration.md)
- [🔌 Plugin Development](docs/plugin-development.md)
- [🏗️ Architecture Guide](docs/architecture.md)
- [🤝 Contributing Guide](CONTRIBUTING.md)

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Quick Contribution Steps

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for your changes
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## 📈 Roadmap

- [ ] **Q1 2025**: Core DORA metrics with 5+ integrations
- [ ] **Q2 2025**: Community benchmarking and mobile app
- [ ] **Q3 2025**: Advanced analytics and ML insights
- [ ] **Q4 2025**: Enterprise features and marketplace

## 🎖️ Community

- 💬 [Discussions](https://github.com/sirhCC/MetricHub/discussions) - Ask questions and share ideas
- 🐛 [Issues](https://github.com/sirhCC/MetricHub/issues) - Report bugs and request features
- 📧 [Newsletter](https://metrichub.dev/newsletter) - Stay updated with releases

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [DORA Research Program](https://dora.dev/) for metrics definitions
- [OpenTelemetry](https://opentelemetry.io/) for observability standards
- All the amazing open-source projects that make this possible

---

**⭐ Star this repository if you find it useful!**
