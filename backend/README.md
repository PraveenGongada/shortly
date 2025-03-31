# Shortly

<div align="center">
  <img src="https://raw.githubusercontent.com/PraveenGongada/shortly/refs/heads/main/docs/images/logo.svg" alt="Shortly Logo" width="200" />
  <p></p>

[![Release](https://img.shields.io/github/v/release/PraveenGongada/shortly?style=flat-square)](https://github.com/PraveenGongada/shortly/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/PraveenGongada/shortly/backend)](https://goreportcard.com/report/github.com/PraveenGongada/shortly)
[![Go Version](https://img.shields.io/badge/Go-1.23-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-latest-336791?style=flat-square&logo=postgresql)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-1.0.0-2496ED?style=flat-square&logo=docker)](https://www.docker.com/)
[![License](https://img.shields.io/github/license/PraveenGongada/shortly?style=flat-square)](LICENSE)

  <p></p>
  <p>Shortly is a fast, scalable URL shortener built with Go. Create short links, track analytics, and manage your URLs via a simple REST API.</p>
</div>

## ✨ Features

- 🔐 User authentication with JWT
- 🔗 Create short URLs
- 📊 Track redirect analytics
- 🧩 RESTful API design
- 📱 Clean architecture pattern
- 🐳 Docker and Docker Compose support
- 📝 Structured logging
- 🧪 Database migrations
- ⚡ High performance

## 🧱 Tech Stack

- **Backend**: Golang
- **Database**: PostgreSQL
- **Authentication**: JWT
- **API Documentation**: Swagger
- **Container**: Docker & Docker Compose
- **Logging**: Zerolog
- **Migration**: Golang-Migrate

## 🏗️ Architecture

Shortly follows the clean architecture pattern with a clear separation of concerns:

- **Domain Layer**: Core business logic and entities
- **Application Layer**: Use case implementations
- **Infrastructure Layer**: External services, repositories, and frameworks

```
├── internal/
│   ├── domain/          # Business entities and interfaces
│   │   ├── url/         # URL domain
│   │   └── user/        # User domain
│   ├── application/     # Application services
│   │   └── service/     # Application services implementation
│   └── infrastructure/  # Infrastructure implementations
│       ├── auth/        # Authentication
│       ├── config/      # Configuration
│       ├── http/        # HTTP server
│       ├── logging/     # Logging
│       └── persistence/ # Data persistence
└── cmd/                 # Applications entry points
```

## 📋 Prerequisites

- Go 1.20+
- Docker and Docker Compose
- Make (for running commands)

## 🚀 Getting Started

### Clone the Repository

```bash
git clone https://github.com/praveengongada/shortly.git
cd shortly/backend
```

### Configuration

1. Copy the example config file:

```bash
cp internal/infrastructure/config/config.yaml.example internal/infrastructure/config/config.yaml
```

2. Update the configuration values as needed, particularly the RSA keys for JWT authentication.

### Generate RSA Keys

```bash
# Generate private key
openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048

# Extract public key
openssl rsa -pubout -in private_key.pem -out public_key.pem

# Format for config.yaml
cat private_key.pem | sed 's/^/        /'
cat public_key.pem | sed 's/^/        /'
```

### Run with Docker

The easiest way to get started is using Docker Compose:

```bash
make run
```

This will:

1. Start the PostgreSQL database
2. Run migrations
3. Build and start the application

### Run for Development

```bash
make dev
```

This starts the application with hot reload enabled for easier development.

## 📚 API Documentation

Shortly comes with Swagger documentation available at `/swagger/index.html` when the server is running.

### Main Endpoints

| Method | Endpoint                        | Description                  | Auth Required |
| ------ | ------------------------------- | ---------------------------- | ------------- |
| POST   | `/api/user/register`            | Register a new user          | No            |
| POST   | `/api/user/login`               | Login a user                 | No            |
| GET    | `/api/user/logout`              | Logout a user                | No            |
| POST   | `/api/url/create`               | Create a short URL           | Yes           |
| GET    | `/api/urls`                     | Get all user URLs            | Yes           |
| PATCH  | `/api/url/update`               | Update a URL                 | Yes           |
| DELETE | `/api/url/{urlId}`              | Delete a URL                 | Yes           |
| GET    | `/api/url/analytics/{shortUrl}` | Get URL analytics            | Yes           |
| GET    | `/{shortUrl}`                   | Redirect to the original URL | No            |

## 🛠️ Makefile Commands

Shortly comes with helpful Make commands to streamline development:

```bash
# Start all services and the application
make run

# Run in development mode with hot reload
make dev

# Apply database migrations
make migrate-up

# Rollback migrations
make migrate-down

# Create a new migration
make migrate-create MIGRATION_NAME=your_migration_name

# Access the PostgreSQL shell
make shell-db

# Start only the services (Docker containers)
make start-services

# Stop all services
make stop-services

# Build the application
make build
```

## 📊 Database Schema

The application uses two main tables:

**User Table**

```sql
CREATE TABLE "user" (
    "id" character(36) NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "email" TEXT NOT NULL UNIQUE,
    "password" TEXT NOT NULL,
    "created_at" timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updated_at" timestamp with time zone
);
```

**URL Table**

```sql
CREATE TABLE url (
    "id" character(36) NOT NULL PRIMARY KEY,
    "user_id" character(36) NOT NULL REFERENCES "user"(id),
    "short_url" varchar(7) NOT NULL UNIQUE,
    "long_url" TEXT NOT NULL,
    "redirects" INT NOT NULL DEFAULT 0,
    "created_at" timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updated_at" timestamp with time zone
);
```

## 🌐 Frontend Integration

Shortly's API is designed to work with any frontend. See the CORS settings in `router.go` for allowed origins.

A sample React frontend can be found at [frontend](https://github.com/PraveenGongada/Shortly/blob/main/frontend).

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! Feel free to check [issues page](https://github.com/praveengongada/shortly/issues).

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

Third-party library attributions can be found in the [NOTICE](NOTICE) file.

---

<div align="center">
  <p>Made with ❤️ by <a href="https://github.com/PraveenGongada">Praveen Kumar</a></p>
  <p>
    <a href="https://linkedin.com/in/praveengongada">LinkedIn</a> •
    <a href="https://praveengongada.com">Website</a>
  </p>
</div>
