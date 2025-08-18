# Shortly - URL Shortener Service

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16+-336791?style=flat-square&logo=postgresql)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7+-DC382D?style=flat-square&logo=redis)](https://redis.io/)
[![License](https://img.shields.io/github/license/PraveenGongada/shortly?style=flat-square)](LICENSE)

A secure, multi-user URL shortener service built with Go following clean architecture principles. Create short links, track click analytics, and manage URLs through a REST API with JWT authentication.

## ✨ Features

### Core Functionality

- **URL Shortening**: Generate 7-character alphanumeric short codes automatically
- **User Authentication**: Secure JWT-based authentication with RSA signing
- **URL Management**: Create, update, delete, and list your short URLs
- **Click Analytics**: Track redirect counts for each short URL
- **Multi-User Support**: Each user manages their own URLs independently

### Security & Quality

- **Secure Authentication**: JWT tokens with RSA-256 signing
- **Password Security**: Bcrypt hashing with automatic salting
- **Input Validation**: Email format, password strength, and URL validation
- **Clean Architecture**: Domain-driven design with clear separation of concerns
- **Structured Logging**: Comprehensive logging with zerolog

### Technical Features

- **PostgreSQL Database**: Reliable data persistence with connection pooling
- **Redis Caching**: Redis is used to cache URL's for faster lookups
- **Graceful Shutdown**: Proper cleanup on application termination
- **Health Checks**: Built-in health monitoring endpoint
- **API Documentation**: Swagger/OpenAPI specification

## 🚀 Quick Start

### Prerequisites

- [Go 1.23+](https://golang.org/dl/)
- [Docker & Docker Compose](https://docs.docker.com/get-docker/)
- [Make](https://www.gnu.org/software/make/) (optional, for convenience commands)

### Installation

```bash
# Clone the repository
git clone https://github.com/PraveenGongada/shortly.git
cd shortly/backend

# Start all services (PostgreSQL, Redis, and the app)
make run
```

The service will be available at `http://localhost:8080`.

### Manual Setup

```bash
# Start infrastructure services
make start-services

# Run database migrations
make migrate-up

# Start the application
go run cmd/shortly/main.go
```

## 📊 API Endpoints

### Public Endpoints

- `GET /{shortCode}` - Redirect to original URL
- `GET /api/{shortCode}` - Get original URL without redirect
- `POST /api/user/register` - Register new user
- `POST /api/user/login` - User login
- `GET /api/user/logout` - User logout
- `GET /api/health` - Health check

### Authenticated Endpoints

- `POST /api/url/create` - Create short URL
- `PATCH /api/url/update` - Update existing URL
- `DELETE /api/url/{urlId}` - Delete URL
- `GET /api/url/analytics/{shortCode}` - Get click analytics
- `GET /api/urls` - List user's URLs (paginated)

Interactive API documentation available at `/swagger/index.html` when running.

## 🏗️ Architecture

Shortly follows Clean Architecture principles:

```
internal/
├── domain/           # Business entities and rules
│   ├── url/         # URL shortening domain
│   ├── user/        # User management domain
│   └── shared/      # Shared domain components
├── application/     # Use case implementations
│   └── service/     # Application services
└── infrastructure/  # External concerns
    ├── http/        # REST API handlers
    ├── persistence/ # PostgreSQL repositories
    ├── cache/       # Redis client
    ├── auth/        # JWT authentication
    └── config/      # Configuration management
```

## 🛠️ Development

### Common Commands

```bash
# Development with hot reload
make dev

# Run database migrations
make migrate-up
make migrate-down

# Create new migration
make migrate-create MIGRATION_NAME=add_new_feature

# Database access
make shell-db

# Build application
make build

# Generate API documentation
make swagger
```

### Configuration

The application uses YAML configuration files in the `configs/` directory:

- `application.yaml` - Server and app settings
- `database.yaml` - PostgreSQL and Redis configuration
- `auth.yaml` - JWT authentication settings

For development, copy and modify the example configuration:

```bash
cp config.yaml.example config.yaml
```

### JWT Keys Setup

```bash
# Generate RSA keys for JWT signing
mkdir -p keys
openssl genpkey -algorithm RSA -out keys/private.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in keys/private.pem -out keys/public.pem
```

## 🔧 Database Schema

### Users Table

- User accounts with email/password authentication
- Bcrypt password hashing
- UUID primary keys

### URLs Table

- Short URL mappings linked to users
- 7-character alphanumeric short codes
- Redirect counter for analytics
- User ownership enforcement

See [Database Documentation](docs/DATABASE.md) for complete schema details.

## 🚢 Deployment

### Docker Deployment

```bash
# Production deployment with Docker Compose
make run
```

### Manual Deployment

1. Set up PostgreSQL and Redis
2. Configure environment variables
3. Run migrations: `make migrate-up`
4. Build and run: `make build && ./bin/main`

See [Deployment Guide](docs/DEPLOYMENT.md) for production deployment instructions.

## 📚 Documentation

- **[API Reference](docs/API.md)** - Complete REST API documentation
- **[Architecture Guide](docs/ARCHITECTURE.md)** - System design and components
- **[Database Schema](docs/DATABASE.md)** - Data models and relationships
- **[Deployment Guide](docs/DEPLOYMENT.md)** - Production deployment
- **[Development Setup](docs/DEVELOPMENT.md)** - Local development guide
- **[Configuration](docs/CONFIGURATION.md)** - Settings and environment variables
- **[Troubleshooting](docs/TROUBLESHOOTING.md)** - Common issues and solutions

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Follow the [Development Guidelines](docs/DEVELOPMENT.md)
4. Submit a pull request

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

Third-party library attributions can be found in the [NOTICE](NOTICE) file.

## 🆘 Support

- **Issues**: [GitHub Issues](https://github.com/PraveenGongada/shortly/issues)
- **Documentation**: [Project Documentation](docs/)

---

<div align="center">
Built with ❤️ by <a href="https://praveengongada.com">Praveen Kumar</a>
</div>
