# Development Guide

## Table of Contents

- [Development Environment Setup](#development-environment-setup)
- [Project Structure](#project-structure)
- [Code Standards and Conventions](#code-standards-and-conventions)
- [Testing Strategy](#testing-strategy)
- [Development Tools](#development-tools)

## Development Environment Setup

### Prerequisites

Ensure you have the following installed:

- **Go 1.23+**: [Download Go](https://golang.org/dl/)
- **Docker Desktop**: [Download Docker](https://www.docker.com/products/docker-desktop)
- **Make**: Usually pre-installed on macOS/Linux, [download for Windows](http://gnuwin32.sourceforge.net/packages/make.htm)
- **Git**: [Download Git](https://git-scm.com/downloads)
- **PostgreSQL Client** (optional): For direct database access
- **Redis CLI** (optional): For cache debugging

### Local Setup

1. **Clone the repository**:

```bash
git clone https://github.com/PraveenGongada/shortly.git
cd shortly/backend
```

2. **Install dependencies**:

```bash
go mod download
go mod tidy
```

3. **Install development tools**:

```bash
# Install Air for hot reload
go install github.com/air-verse/air@latest

# Install golangci-lint for code quality
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install Swagger for API documentation
go install github.com/swaggo/swag/cmd/swag@latest

# Install Wire for dependency injection
go install github.com/google/wire/cmd/wire@latest
```

4. **Set up configuration**:

```bash
# Copy example config (if needed)
cp config.yaml.example config.yaml

# Generate RSA keys for JWT
mkdir -p keys
openssl genpkey -algorithm RSA -out keys/private.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in keys/private.pem -out keys/public.pem
```

5. **Start development environment**:

```bash
# Start all services with hot reload
make dev
```

## Project Structure

Understanding the codebase organization:

```
shortly/backend/
├── cmd/                    # Application entry points
│   ├── shortly/           # Main application
│   └── wire-app/          # Wire code generation utility
├── internal/              # Private application code
│   ├── application/       # Application services (use cases)
│   ├── domain/           # Domain logic (entities, business rules)
│   └── infrastructure/   # Infrastructure implementations
├── api/                  # OpenAPI/Swagger documentation
├── configs/              # Configuration files
├── docs/                 # Project documentation
├── keys/                 # JWT signing keys
├── migrations/           # Database migration files
└── tools.go              # Go tool dependencies
```

### Layer Responsibilities

- **Domain Layer**: Pure business logic, no external dependencies
- **Application Layer**: Orchestrates domain logic, handles use cases
- **Infrastructure Layer**: Implements domain interfaces, handles external systems

## Code Standards and Conventions

### Go Style Guidelines

Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and [Effective Go](https://golang.org/doc/effective_go.html).

## Testing Strategy

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...
```

## Development Tools

### Makefile Commands

```bash
# Development
make dev              # Start with hot reload
make build           # Build binary

# Database
make migrate-up       # Apply migrations
make migrate-down     # Rollback migrations

# Code Quality
make test           # Run all tests
```

---

This development guide provides comprehensive instructions for setting up your development environment and contributing effectively to the Shortly project. For questions or suggestions, please open an issue or discussion on GitHub.
