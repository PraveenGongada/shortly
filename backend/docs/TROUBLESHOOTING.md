# Troubleshooting Guide

## Table of Contents

- [Common Issues](#common-issues)
- [Application Startup Problems](#application-startup-problems)
- [Database Connection Issues](#database-connection-issues)
- [Redis Cache Problems](#redis-cache-problems)
- [Authentication and JWT Issues](#authentication-and-jwt-issues)
- [Docker and Container Issues](#docker-and-container-issues)
- [Getting Help](#getting-help)

## Common Issues

### Application Won't Start

**Symptom**: Application exits immediately or fails to start

**Common Causes**:
1. Configuration file not found or invalid
2. Database connection failure
3. Port already in use
4. Missing RSA keys for JWT

**Solutions**:

```bash
# Check if port is in use
lsof -i :8080

# Kill process using the port
kill -9 $(lsof -t -i:8080)

# Verify configuration
cat configs/application.yaml

# Check RSA keys exist
ls -la keys/

# Run with verbose logging
LOG_LEVEL=debug go run cmd/shortly/main.go
```

### Database Connection Refused

**Symptom**: `connection refused` or `connection timeout` errors

**Diagnosis**:
```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Test connection manually
psql -h localhost -p 5432 -U postgres -d shortly
```

**Solutions**:
```bash
# Start PostgreSQL
docker-compose up -d postgres

# Check PostgreSQL logs
docker-compose logs postgres
```

### Short URL Not Found (404)

**Symptom**: Accessing short URLs returns 404 Not Found

**Diagnosis**:
```bash
# Check if URL exists in database
docker-compose exec postgres psql -U postgres -d shortly -c "SELECT * FROM url WHERE short_url = 'abc123';"

# Check Redis cache
docker-compose exec redis redis-cli get "url:abc123"

# Check application logs
docker-compose logs app | grep abc123
```

## Application Startup Problems

### Configuration Validation Errors

**Error**: `configuration validation failed`

**Solutions**:
- Check required fields in `config.yaml`

### Missing RSA Keys

**Error**: `failed to load RSA keys` or `no such file or directory`

**Solutions**:
```bash
# Generate RSA keys
mkdir -p keys
openssl genpkey -algorithm RSA -out keys/private.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in keys/private.pem -out keys/public.pem
```

### Port Already in Use

**Error**: `bind: address already in use`

**Solutions**:
```bash
# Find and kill process using port 8080
sudo lsof -ti:8080 | xargs kill -9

# Use different port
export APP_PORT=8081
```

## Database Connection Issues

### Database Migration Failures

**Error**: Migration files not found or migration failed

**Solutions**:
```bash
# Check migration files exist
ls -la migrations/

# Run migrations manually
make migrate-up
```

## Redis Cache Problems

### Redis Connection Failed

**Error**: `dial tcp: connection refused` for Redis

**Solutions**:
```bash
# Check Redis status
docker-compose ps redis

# Test Redis connection
docker-compose exec redis redis-cli ping
```

## Authentication and JWT Issues

### Invalid JWT Token

**Error**: `invalid token` or `token expired`

**Solutions**:
- Verify RSA keys match
- Check JWT configuration
- Generate new token for testing

### User Authentication Failed

**Error**: `invalid credentials` or `user not found`

**Solutions**:
- Verify user registration was successful
- Check password hashing implementation

## Docker and Container Issues

### Container Won't Start

**Error**: Container exits immediately

**Diagnosis**:
```bash
# Check container logs
docker-compose logs app
```

## Getting Help

### Before Asking for Help

1. **Check this troubleshooting guide** for common issues
2. **Review application logs** for error messages
3. **Verify configuration** against the examples

### Support Channels

- **GitHub Issues**: [Report bugs and feature requests](https://github.com/PraveenGongada/shortly/issues)

---

This troubleshooting guide covers most common issues you'll encounter with Shortly. For issues not covered here, please check the GitHub issues or start a discussion.