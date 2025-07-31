# API Reference

## Table of Contents

- [Overview](#overview)
- [Authentication](#authentication)
- [Response Format](#response-format)
- [Error Handling](#error-handling)
- [User Management](#user-management)
- [URL Management](#url-management)
- [URL Redirection](#url-redirection)
- [Analytics](#analytics)
- [Health Check](#health-check)
- [Status Codes](#status-codes)

## Overview

The Shortly API is a RESTful HTTP API built with Go and Chi router. All endpoints use JSON for request and response payloads, with consistent error handling and response formats.

**Base URL**: `http://localhost:8080`  
**API Prefix**: `/api` (for most endpoints)

### Interactive Documentation

Swagger UI is available at `/swagger/index.html` when the server is running.

## Authentication

Shortly uses JWT (JSON Web Token) authentication for secure user authentication.

### How Authentication Works

1. **Register** or **Login** to receive a JWT token
2. Include the token in the `Authorization` header: `Bearer <token>`
3. Tokens expire after 24 hours (configurable)
4. Tokens are also set as HTTP-only cookies

### Authentication Header Format

```http
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...
```

## Response Format

All API responses follow a consistent JSON structure:

### Success Response

```json
{
  "message": "Operation successful",
  "data": {
    // Response payload (varies by endpoint)
  }
}
```

### Error Response

```json
{
  "message": "Error description",
  "data": null
}
```

## Error Handling

### Common HTTP Status Codes

| Status Code | Description           | When It Occurs                                    |
| ----------- | --------------------- | ------------------------------------------------- |
| 200         | OK                    | Successful GET, PATCH, DELETE operations          |
| 201         | Created               | Successful resource creation (POST)               |
| 400         | Bad Request           | Invalid request format or missing required fields |
| 401         | Unauthorized          | Missing, invalid, or expired JWT token            |
| 404         | Not Found             | Resource doesn't exist (URL, user, endpoint)      |
| 409         | Conflict              | Resource already exists (duplicate email)         |
| 422         | Unprocessable Entity  | Input validation failed                           |
| 500         | Internal Server Error | Unexpected server error                           |

## User Management

### Register User

Create a new user account and receive a JWT token.

**Endpoint**: `POST /api/user/register`

**Request Body**:

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePass123"
}
```

### Login User

Authenticate an existing user and receive a JWT token.

**Endpoint**: `POST /api/user/login`

**Request Body**:

```json
{
  "email": "john@example.com",
  "password": "SecurePass123"
}
```

### Logout User

Invalidate the current session by clearing the authentication cookie.

**Endpoint**: `GET /api/user/logout`

## URL Management

### Create Short URL

Generate a short URL from a long URL.

**Endpoint**: `POST /api/url/create`

**Authentication**: Required

**Request Body**:

```json
{
  "long_url": "https://www.example.com/very/long/path/to/resource"
}
```

### Get User URLs

Retrieve a paginated list of URLs created by the authenticated user.

**Endpoint**: `GET /api/urls`

**Authentication**: Required

### Update URL

Update the destination URL for an existing short URL. Only the URL owner can update it.

**Endpoint**: `PATCH /api/url/update`

**Authentication**: Required

### Delete URL

Permanently delete a short URL. Only the URL owner can delete it.

**Endpoint**: `DELETE /api/url/{urlId}`

**Authentication**: Required

## URL Redirection

### Redirect to Original URL

Redirect from a short URL to its original destination.

**Endpoint**: `GET /{shortCode}`

**Authentication**: Not required

### Get Original URL (Without Redirect)

Retrieve the original URL without performing a redirect.

**Endpoint**: `GET /api/{shortCode}`

**Authentication**: Not required

## Analytics

### Get URL Analytics

Retrieve click analytics for a specific short URL. Only the URL owner can access analytics.

**Endpoint**: `GET /api/url/analytics/{shortCode}`

**Authentication**: Required

## Health Check

### Application Health Status

Check the health status of the application and its dependencies.

**Endpoint**: `GET /api/health`

**Authentication**: Not required

## Status Codes

### Success Codes

- **200 OK**: Request successful (GET, PATCH, DELETE)
- **201 Created**: Resource created successfully (POST)
- **302 Found**: Redirect response (short URL redirection)

### Client Error Codes

- **400 Bad Request**: Invalid request format or structure
- **401 Unauthorized**: Authentication required or token invalid
- **404 Not Found**: Resource not found or not accessible
- **409 Conflict**: Resource conflict (e.g., email already exists)
- **422 Unprocessable Entity**: Input validation failed

### Server Error Codes

- **500 Internal Server Error**: Unexpected server error

---

This API reference accurately reflects the current implementation of the Shortly URL shortener service. For additional support, please refer to the [troubleshooting guide](TROUBLESHOOTING.md) or [open an issue](https://github.com/PraveenGongada/shortly/issues).
