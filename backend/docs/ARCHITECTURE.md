# Architecture Guide

## Table of Contents

- [Overview](#overview)
- [Design Principles](#design-principles)
- [System Architecture](#system-architecture)
- [Layer Structure](#layer-structure)
- [Technology Stack](#technology-stack)

## Overview

Shortly is built using Clean Architecture principles with Domain-Driven Design (DDD) patterns. The system is designed for high performance, scalability, and maintainability.

### Core Design Goals

- **Separation of Concerns**: Clear boundaries between business logic and infrastructure
- **Testability**: Easy to unit test business logic without external dependencies
- **Maintainability**: Code organized around business domains for easier understanding
- **Scalability**: Stateless design supporting horizontal scaling
- **Performance**: Optimized for high throughput with caching and efficient database access

## Design Principles

### 1. Clean Architecture

The application follows Uncle Bob's Clean Architecture pattern:

```
┌──────────────────────────────────┐
│           Infrastructure         │
│  ┌──────────────────────────┐    │
│  │        Application       │    │
│  │  ┌──────────────────┐    │    │
│  │  │      Domain      │    │    │
│  │  │   (Entities &    │    │    │
│  │  │  Business Rules) │    │    │
│  │  └──────────────────┘    │    │
│  └──────────────────────────┘    │
└──────────────────────────────────┘
```

### 2. Domain-Driven Design

Business logic is organized around two main domains:

- **URL Domain**: URL shortening, validation, analytics
- **User Domain**: Authentication, user management

### 3. Dependency Inversion

All dependencies point inward toward the domain layer. Infrastructure depends on domain interfaces, not implementations.

## System Architecture

### High-Level Component Diagram

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Client Apps   │────│   Load Balancer │────│   API Gateway   │
│ (Web, Mobile,   │    │   (nginx/ALB)   │    │   (Optional)    │
│  CLI, etc.)     │    └─────────────────┘    └─────────────────┘
└─────────────────┘                                     │
                                                        │
                    ┌───────────────────────────────────┘
                    │
            ┌───────▼────────┐
            │  Shortly API   │
            │   (Go Server)  │
            └───────┬────────┘
                    │
        ┌───────────┼───────────┐
        │           │           │
   ┌────▼────-┐ ┌───▼────┐ ┌───▼─────┐
   │PostgreSQL│ │ Redis  │ │ Monitor │
   │Database  │ │ Cache  │ │  Stack  │
   └────────-─┘ └────────┘ └─────────┘
```

## Layer Structure

### Domain Layer (`internal/domain/`)

The innermost layer containing pure business logic:

```
domain/
├── interfaces/           # Domain service interfaces
│   ├── url.go           # URL repository contracts
│   └── user.go          # User repository contracts
├── shared/              # Cross-domain shared components
│   ├── config/          # Configuration interfaces
│   ├── errors/          # Domain-specific errors
│   └── logger/          # Logging abstractions
├── url/                 # URL shortening domain
│   ├── entity/          # URL business entities
│   ├── repository/      # Repository interfaces
│   ├── service/         # Domain services (validation, generation)
│   ├── cache/           # Caching interfaces
│   └── valueobject/     # DTOs and value objects
└── user/                # User management domain
    ├── entity/          # User entities
    ├── repository/      # Repository interfaces
    ├── service/         # User services (hashing, validation)
    └── valueobject/     # User DTOs
```

#### Key Domain Components

- **Entities**: Core business objects (User, URL)
- **Value Objects**: Immutable objects representing domain concepts
- **Domain Services**: Business logic that doesn't belong to entities
- **Repository Interfaces**: Data access contracts

### Application Layer (`internal/application/`)

Orchestrates domain logic to fulfill use cases:

```
application/
└── service/
    ├── url_service.go    # URL management use cases
    └── user_service.go   # User management use cases
```

#### Responsibilities

- Coordinate multiple domain services
- Handle application-specific business rules
- Manage transaction boundaries
- Convert between domain models and DTOs

### Infrastructure Layer (`internal/infrastructure/`)

Implements domain interfaces and handles external concerns:

```
infrastructure/
├── auth/                # JWT authentication
├── cache/redis/         # Redis caching implementation
├── config/              # Configuration management
├── http/                # REST API implementation
│   ├── handler/         # HTTP request handlers
│   ├── middleware/      # Cross-cutting concerns
│   ├── response/        # Response formatting
│   └── router/          # Route definitions
├── logging/             # Structured logging
├── persistence/postgres/ # PostgreSQL repositories
└── wire/                # Dependency injection
```

## Component Interactions

### Request Flow Diagram

```
┌─────────┐    ┌─────────┐    ┌─────────────┐    ┌─────────────┐
│ Client  │───▶│ Router  │───▶│  Handler    │───▶│ Application │
└─────────┘    └─────────┘    │(Controller) │    │  Service    │
                              └─────────────┘    └─────────────┘
                                     │                  │
                                     ▼                  ▼
┌─────────┐    ┌─────────┐    ┌─────────────┐    ┌─────────────┐
│Response │◀───│Formatter│◀───│ Middleware  │    │   Domain    │
└─────────┘    └─────────┘    │   Chain     │    │  Service    │
                              └─────────────┘    └─────────────┘
                                                        │
                                                        ▼
                               ┌─────────────┐    ┌─────────────┐
                               │    Cache    │    │ Repository  │
                               │  (Redis)    │    │(PostgreSQL) │
                               └─────────────┘    └─────────────┘
```

### Authentication Flow

```
┌─────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│ Client  │───▶│ Auth        │───▶│ JWT Token   │───▶│ Protected   │
│ Login   │    │ Middleware  │    │ Validation  │    │ Resource    │
└─────────┘    └─────────────┘    └─────────────┘    └─────────────┘
    │                 │                  │                  │
    ▼                 ▼                  ▼                  ▼
┌─────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│User Repo│    │ RSA Key     │    │ Claims      │    │ Business    │
│ Lookup  │    │ Validation  │    │ Extraction  │    │ Logic       │
└─────────┘    └─────────────┘    └─────────────┘    └─────────────┘
```

## Data Flow

### URL Creation Flow

1. **HTTP Request**: Client sends POST to `/api/url/create`
2. **Authentication**: JWT middleware validates token
3. **Request Validation**: Handler validates input format
4. **Business Logic**: Application service coordinates:
   - URL validation (domain service)
   - Short code generation (domain service)
   - Duplicate check (repository)
   - Persistence (repository)
   - Cache warming (cache service)
5. **Response**: Formatted response with short URL

### URL Resolution Flow

1. **HTTP Request**: Client accesses `/{shortCode}`
2. **Cache Lookup**: Check Redis for cached mapping
3. **Database Fallback**: Query PostgreSQL if cache miss
4. **Analytics Update**: Increment redirect counter
5. **Cache Update**: Store result in Redis
6. **HTTP Redirect**: Return 302 redirect to original URL

## Technology Stack

### Core Technologies

| Component          | Technology | Justification                                          |
| ------------------ | ---------- | ------------------------------------------------------ |
| **Language**       | Go         | High performance, excellent concurrency, strong typing |
| **Web Framework**  | Chi v5     | Lightweight, fast, middleware-friendly                 |
| **Database**       | PostgreSQL | ACID compliance, excellent performance, JSON support   |
| **Cache**          | Redis      | In-memory performance for URL lookups                  |
| **Authentication** | JWT        | Stateless, secure, scalable                            |

### Development & Operations

| Component                | Technology      | Purpose                             |
| ------------------------ | --------------- | ----------------------------------- |
| **Dependency Injection** | Google Wire     | Compile-time DI, better performance |
| **Configuration**        | Viper           | Flexible config management          |
| **Logging**              | Zerolog         | High-performance structured logging |
| **Database Migration**   | golang-migrate  | Version-controlled schema changes   |
| **API Documentation**    | Swagger/OpenAPI | Interactive API documentation       |
| **Hot Reload**           | Air             | Development productivity            |

---

This architecture supports Shortly's current needs while providing a clear path for future scaling.
