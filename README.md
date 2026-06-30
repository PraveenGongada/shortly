# Shortly

<div align="center">
  <img src="https://raw.githubusercontent.com/PraveenGongada/shortly/refs/heads/main/docs/images/logo.svg" alt="Shortly Logo" width="200" />

[![Backend Build](https://img.shields.io/github/actions/workflow/status/PraveenGongada/Shortly/backend-build.yaml?style=flat-square&logo=github&label=backend)](https://github.com/PraveenGongada/Shortly/actions)
[![Frontend Build](https://img.shields.io/github/actions/workflow/status/PraveenGongada/Shortly/frontend-build.yaml?style=flat-square&logo=github&label=frontend)](https://github.com/PraveenGongada/Shortly/actions)
[![License](https://img.shields.io/badge/License-Apache%202.0-green?style=flat-square)](LICENSE)

  <p></p>
  <p>A simple and efficient URL shortening service that allows users to convert long URLs into short, shareable links</p>
</div>

## 🌟 Overview

Shortly is a complete URL shortening solution combining a beautiful, responsive frontend with a powerful backend. This monorepo contains both components as separate modules.

## 🖥️ Screenshots

<div align="center">
  <img src="https://raw.githubusercontent.com/PraveenGongada/shortly/refs/heads/main/frontend/docs/images/dashboard.png" alt="Dashboard" width="80%" style="border-radius: 12px;"/>
  <p><em>Dashboard - Manage all your shortened URLs</em></p>
  
  <br />
  
  <div style="display: flex; justify-content: space-between;">
    <img src="https://raw.githubusercontent.com/PraveenGongada/shortly/refs/heads/main/frontend/docs/images/home.png" alt="Home Page" width="48%" style="border-radius: 12px;"/>
    <img src="https://raw.githubusercontent.com/PraveenGongada/shortly/refs/heads/main/frontend/docs/images/qr.png" alt="Create URL" width="48%" style="border-radius: 12px;"/>
  </div>
  <p><em>Home Page and QR Code Display for Short URL</em></p>
</div>

## ✨ Features

- 🔗 Instant URL shortening
- 📊 Comprehensive analytics for tracking link performance
- 🔐 Secure user authentication with JWT
- 📱 Responsive design optimized for all devices
- 📈 Dashboard to manage all shortened URLs
- 🔄 RESTful API design with clean architecture
- 🔍 Centralized logging and monitoring for operational insights

## 🏗️ Repository Structure

```
shortly/
├── frontend/               # React + Vite frontend (static SPA, served by nginx)
│   ├── src/                # Pages, components, hooks, lib
│   ├── nginx.conf          # Production web server config
│   ├── README.md           # Frontend-specific documentation
│   └── ...
├── backend/                # Go backend service
│   ├── cmd/                # Application entry points
│   ├── internal/           # Clean architecture implementation
│   ├── README.md           # Backend-specific documentation
│   └── ...
├── infra/                  # Kubernetes configurations (submodule)
│   ├── namespaces.yaml     # Kubernetes namespace definitions
│   ├── istio.yaml          # Service mesh configuration
│   ├── cert-manager.yaml   # TLS certificate management
│   ├── deployment.yaml     # Application deployments
│   ├── elastic-search.yaml # Elasticsearch configuration
│   ├── kibana.yaml         # Kibana dashboard configuration
│   ├── logstash.yaml       # Log processing pipeline
│   ├── filebeat.yaml       # Log collection agent
│   └── monitoring/         # Prometheus and Grafana configs
└── docs/                   # Project documentation and assets
```

## 🧱 Tech Stack

### Frontend

- **Build Tool**: Vite
- **Library**: React 18
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **Data Fetching**: TanStack Query
- **Forms**: React Hook Form + Zod
- **UI Components**: Radix UI/shadcn-style
- **Icons**: Lucide React
- **Serving**: nginx (static)

### Backend

- **Language**: Golang
- **Database**: PostgreSQL
- **Authentication**: JWT
- **API Documentation**: Swagger
- **Containerization**: Docker & Docker Compose
- **Logging**: Zerolog
- **Migration**: Golang-Migrate

### Monitoring & Observability

- **Metrics**: Prometheus & Grafana
- **Logging**: ELK Stack (Elasticsearch, Logstash, Kibana)
- **Log Collection**: Filebeat
- **Service Mesh**: Istio

## 📚 Detailed Documentation

- [Frontend Documentation](https://github.com/PraveenGongada/Shortly/blob/main/frontend/README.md)
- [Backend Documentation](https://github.com/PraveenGongada/Shortly/blob/main/backend/README.md)

## 🚀 Deployment & Infrastructure

Shortly uses Kubernetes for deployment and infrastructure management. The Kubernetes configuration files are maintained in a separate private repository and linked to this repository using Git submodules.

### Infrastructure Setup

```bash
# Clone the complete repository with infrastructure configs
git clone --recursive git@github.com:PraveenGongada/shortly.git

# Or initialize submodules after cloning
git submodule init
git submodule update

# To fetch latest changes
git submodule update --remote
```

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! Feel free to check [issues page](https://github.com/praveengongada/shortly/issues).

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

---

<div align="center">
  <p>Made with ❤️ by <a href="https://github.com/PraveenGongada">Praveen Kumar</a></p>
  <p>
    <a href="https://linkedin.com/in/praveengongada">LinkedIn</a> •
    <a href="https://praveengongada.com">Website</a>
  </p>
</div>
