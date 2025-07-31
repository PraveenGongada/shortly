# Database Documentation

## Table of Contents

- [Overview](#overview)
- [Entity Relationship Diagram](#entity-relationship-diagram)
- [Schema Design](#schema-design)
- [Indexing Strategy](#indexing-strategy)
- [Migration Management](#migration-management)

## Overview

Shortly uses PostgreSQL as its primary database, chosen for its excellent performance, ACID compliance, and robust feature set. The database design follows normalization principles while maintaining optimal query performance.

### Database Configuration

- **Engine**: PostgreSQL
- **Connection Pool**: pgxpool
- **Migration Tool**: golang-migrate/migrate
- **Character Set**: UTF-8
- **Timezone**: UTC

## Entity Relationship Diagram

```
┌─────────────────────────────────────┐
│                USER                 │
├─────────────────────────────────────┤
│ id (PK)          │ char(36)         │
│ name             │ text             │
│ email            │ text (unique)    │
│ password         │ text             │
│ created_at       │ timestamptz      │
│ updated_at       │ timestamptz      │
└─────────────────┬───────────────────┘
                  │
                  │ 1:N
                  │
┌─────────────────▼───────────────────┐
│                URL                  │
├─────────────────────────────────────┤
│ id (PK)          │ char(36)         │
│ user_id (FK)     │ char(36)         │
│ short_url        │ varchar(7)       │
│ long_url         │ text             │
│ redirects        │ integer          │
│ created_at       │ timestamptz      │
│ updated_at       │ timestamptz      │
└─────────────────────────────────────┘
```

### Relationship Details

- **One-to-Many**: User → URLs
  - One user can create multiple URLs
  - Each URL belongs to exactly one user
  - Enforced by foreign key constraint

## Schema Design

### User Table

The `user` table stores user account information with secure password hashing.

```sql
CREATE TABLE IF NOT EXISTS "user" (
    "id" character(36) NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "email" TEXT NOT NULL UNIQUE,
    "password" TEXT NOT NULL,
    "created_at" timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updated_at" timestamp with time zone
);
```

#### Field Specifications

| Field        | Type        | Constraints             | Description                       |
| ------------ | ----------- | ----------------------- | --------------------------------- |
| `id`         | char(36)    | PRIMARY KEY, NOT NULL   | UUID v4 identifier                |
| `name`       | text        | NOT NULL                | User's display name (1-100 chars) |
| `email`      | text        | NOT NULL, UNIQUE        | User's email address              |
| `password`   | text        | NOT NULL                | Bcrypt hashed password            |
| `created_at` | timestamptz | NOT NULL, DEFAULT NOW() | Account creation timestamp        |
| `updated_at` | timestamptz | NULL                    | Last modification timestamp       |

#### Constraints and Validations

- **Email Uniqueness**: Enforced at database level
- **Password Security**: Bcrypt hash with cost factor 12
- **UUID Format**: 36-character UUID v4 strings
- **Timezone**: All timestamps stored in UTC

### URL Table

The `url` table stores URL mappings and analytics data.

```sql
CREATE TABLE IF NOT EXISTS url (
    "id" character(36) NOT NULL PRIMARY KEY,
    "user_id" character(36) NOT NULL REFERENCES "user"(id),
    "short_url" varchar(7) NOT NULL UNIQUE,
    "long_url" TEXT NOT NULL,
    "redirects" INT NOT NULL DEFAULT 0,
    "created_at" timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updated_at" timestamp with time zone
);
```

#### Field Specifications

| Field        | Type        | Constraints             | Description                 |
| ------------ | ----------- | ----------------------- | --------------------------- |
| `id`         | char(36)    | PRIMARY KEY, NOT NULL   | UUID v4 identifier          |
| `user_id`    | char(36)    | NOT NULL, FOREIGN KEY   | Reference to user.id        |
| `short_url`  | varchar(7)  | NOT NULL, UNIQUE        | Short URL code (6-7 chars)  |
| `long_url`   | text        | NOT NULL                | Original destination URL    |
| `redirects`  | integer     | NOT NULL, DEFAULT 0     | Click/redirect counter      |
| `created_at` | timestamptz | NOT NULL, DEFAULT NOW() | Creation timestamp          |
| `updated_at` | timestamptz | NULL                    | Last modification timestamp |

#### Constraints and Validations

- **Foreign Key**: `user_id` references `user(id)` with CASCADE delete
- **Short URL Uniqueness**: Enforced at database level
- **URL Length**: Long URLs can be up to 2048 characters
- **Short Code**: Alphanumeric characters, 6-7 length

## Indexing Strategy

Indexes are strategically placed to optimize common query patterns while minimizing storage overhead.

### Primary Indexes

```sql
-- Automatically created with PRIMARY KEY constraints
CREATE UNIQUE INDEX "user_pkey" ON "user" USING btree (id);
CREATE UNIQUE INDEX "url_pkey" ON url USING btree (id);
```

### Secondary Indexes

```sql
-- User table indexes
CREATE UNIQUE INDEX "user_email_idx" ON "user" USING btree (email);

-- URL table indexes
CREATE UNIQUE INDEX "url_short_url_idx" ON url USING btree (short_url);
CREATE INDEX "url_user_id_idx" ON url USING btree (user_id);
```

## Migration Management

Database schema changes are managed using golang-migrate with sequential versioning.

### Migration File Structure

```
migrations/
├── 000001_init_schema.up.sql      # Create initial tables
├── 000001_init_schema.down.sql    # Drop initial tables
└── ...
```

### Migration Commands

```bash
# Apply all pending migrations
make migrate-up

# Rollback last N migrations
make migrate-down ROLLBACK_COUNT=1

# Create new migration files
make migrate-create MIGRATION_NAME=migration_name
```

---

This database documentation provides comprehensive guidance for managing Shortly's PostgreSQL database. Regular monitoring and maintenance ensure optimal performance and data integrity.
